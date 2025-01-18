package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	environment := os.Getenv("ENV")
	if environment != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	r := createRouterAndSetRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	slog.Info(fmt.Sprintf("Listening on port: %v", port))
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
