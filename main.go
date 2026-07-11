package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	store    *Store
	handler  *Handler
	chiLambda *chiadapter.ChiLambda
)

func init() {
	store = NewStore()
	handler = NewHandler(store)
	router := setupRouter(handler)
	chiLambda = chiadapter.NewChiLambda(router)
}

// setupRouter configures the chi router with all routes.
func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Health check
	r.Get("/health", h.Health)

	// Address routes
	r.Route("/addresses", func(r chi.Router) {
		r.Post("/", h.CreateAddress)
		r.Get("/", h.ListAddresses)
		r.Get("/{id}", h.GetAddress)
		r.Put("/{id}", h.UpdateAddress)
		r.Delete("/{id}", h.DeleteAddress)
	})

	// V2 routes
	r.Route("/v2/addresses", func(r chi.Router) {
		r.Get("/search", h.SearchAddresses)
	})

	return r
}

// LambdaHandler is the Lambda handler for API Gateway events.
func LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

// LocalHandler starts a local HTTP server for testing.
func LocalHandler() {
	router := setupRouter(handler)
	fmt.Println("Starting Address API on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func main() {
	// Check if running in Lambda environment
	if lambdaTaskRootEnv := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRootEnv != "" {
		lambda.Start(LambdaHandler)
	} else {
		// Local development mode
		LocalHandler()
	}
}
