package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Lunnaris01/Raidplanner/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	db       *database.Queries
	platform string
	port     string
}

func main() {
	fmt.Println("Civ API started!")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load necessary environment variables with err: %v", err)
	}

	env_platform := os.Getenv("PLATFORM")
	env_dbURL := os.Getenv("TURSO_DATABASE_URL")
	env_dbToken := os.Getenv("TURSO_AUTH_TOKEN")
	env_port := os.Getenv("PORT")
	dbCombinedURL := env_dbURL + "?authToken=" + env_dbToken
	log.Printf("Connecting to db at %s,", env_dbURL)

	sqlitedb, err := sql.Open("libsql", dbCombinedURL)
	if err != nil {
		log.Fatalf("Failed to connect to database with err: %v\n", err)
	}
	defer sqlitedb.Close()

	dbQueries := database.InitDB(dbCombinedURL)

	log.Println("Database connection successful!")

	apiCfg := apiConfig{
		db:       dbQueries,
		platform: env_platform,
		port:     env_port,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", apiCfg.handlerIndex)

	log.Printf("Server running and waiting for requests\n")
	http.ListenAndServe(":"+apiCfg.port, router)

	fmt.Println(apiCfg)

}

func (cfg apiConfig) handlerIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
