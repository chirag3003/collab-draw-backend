package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chirag3003/collab-draw-backend/graph"
	"github.com/chirag3003/collab-draw-backend/graph/resolvers"
	"github.com/chirag3003/collab-draw-backend/internal/auth"
	"github.com/chirag3003/collab-draw-backend/internal/db"
	"github.com/chirag3003/collab-draw-backend/internal/repository"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkHttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	//Loading Env Variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Connect to MongoDB
	db.ConnectMongo()

	// Setting up repositories
	repo := repository.Setup()

	//setting up Clerk
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers.NewResolver(repo)}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			// Extract authorization from connection params
			authHeader := initPayload.Authorization()
			if authHeader == "" {
				// Try to get from other params
				if auth, ok := initPayload["authorization"].(string); ok {
					authHeader = auth
				}
			}

			log.Printf("WebSocket InitFunc - Auth header: %v", authHeader != "")

			// If we have authorization, validate it
			if authHeader != "" {
				// Create a fake request to validate the token
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", authHeader)

				// Use Clerk to verify the session
				clerkClient := clerkHttp.RequireHeaderAuthorization()
				var validatedCtx context.Context
				var authOk bool

				// Create a test handler to capture the context
				testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					claims, ok := clerk.SessionClaimsFromContext(r.Context())
					if ok {
						validatedCtx = context.WithValue(ctx, auth.UserContextKey, claims)
						authOk = true
						log.Printf("WebSocket auth successful for user: %v", claims.Subject)
					}
				})

				// Wrap with Clerk validation
				clerkClient(testHandler).ServeHTTP(nil, req.WithContext(ctx))

				if authOk {
					return validatedCtx, &initPayload, nil
				}
			}

			log.Printf("WebSocket auth failed or no auth provided")
			return ctx, &initPayload, nil
		},

		// The `github.com/gorilla/websocket.Upgrader` is used to handle the transition
		// from an HTTP connection to a WebSocket connection. Among other options, here
		// you must check the origin of the request to prevent cross-site request forgery
		// attacks.
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow exact match on host.
				origin := r.Header.Get("Origin")
				if origin == "" || origin == r.Header.Get("Host") {
					return true
				}

				// Match on allow-listed origins.
				// return slices.Contains([]string{":3000", "https://ui.mysite.com"}, origin)
				return true
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://collab.chirag.codes"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// Custom middleware that allows WebSocket upgrades to bypass auth middleware
	router.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a WebSocket upgrade request
		if r.Header.Get("Upgrade") == "websocket" {
			log.Printf("WebSocket upgrade request detected, bypassing auth middleware")
			srv.ServeHTTP(w, r)
			return
		}
		// For regular HTTP requests, use auth middleware
		auth.Middleware()(srv).ServeHTTP(w, r)
	}))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
