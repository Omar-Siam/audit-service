package main

import (
	"canonicalAuditlog/internal/handlers"
	"canonicalAuditlog/internal/middleware"
	"canonicalAuditlog/internal/repository"
	"canonicalAuditlog/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Todo: Move connection URI to secure config
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://osiam002:...."))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	eventRepo := repository.NewMongoRepository(client)
	eventService := services.NewEventService(eventRepo)

	handler := &handlers.Handlers{
		Service: eventService,
	}

	mux := http.NewServeMux()
	protectedMux := http.NewServeMux()

	// Todo: Use gorilla mux to allow for method based routing
	protectedMux.HandleFunc("/events/create", handler.CreateEvent)
	protectedMux.HandleFunc("/events/query", handler.QueryEvents)
	mux.HandleFunc("/jwt", handler.GetJwt)

	// Todo: Secure jwt endpoint to avoid exploiters
	protectedHandler := middleware.JWTMiddleware(protectedMux)
	mux.Handle("/events/create", protectedHandler)
	mux.Handle("/events/query", protectedHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Service has been spun up")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
