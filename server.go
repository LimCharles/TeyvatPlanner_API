package main

import (
	"log"
	"net/http"
	"os"
	"teyvat_planner_api/auth"
	"teyvat_planner_api/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"	

func main() {
	err := godotenv.Load()
	
	if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
	
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	Database := graph.Connect()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: Database}}))

	graphqlHandlerWithMiddleware := auth.Middleware(Database)(srv)
	
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphqlHandlerWithMiddleware)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
