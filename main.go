package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mosesbenjamin/rss-feed-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envVars := map[string]string{
		"SERVER_PORT":   os.Getenv("SERVER_PORT"),
		"PSQL_HOST":     os.Getenv("PSQL_HOST"),
		"PSQL_PORT":     os.Getenv("PSQL_PORT"),
		"PSQL_USER":     os.Getenv("PSQL_USER"),
		"PSQL_PASSWORD": os.Getenv("PSQL_PASSWORD"),
		"PSQL_DATABASE": os.Getenv("PSQL_DATABASE"),
		"PSQL_SSLMODE":  os.Getenv("PSQL_SSLMODE"),
	}

	for k, val := range envVars {
		if val == "" {
			log.Fatalf("%s environment variable is missing.", k)
		}
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envVars["PSQL_HOST"], envVars["PSQL_PORT"], envVars["PSQL_USER"], envVars["PSQL_PASSWORD"], envVars["PSQL_DATABASE"], envVars["PSQL_SSLMODE"],
	))
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)

	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + envVars["SERVER_PORT"],
		Handler: router,
	}

	log.Printf("Starting server on port %s", envVars["SERVER_PORT"])
	log.Fatal(srv.ListenAndServe())
}
