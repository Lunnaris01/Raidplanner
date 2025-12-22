package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Lunnaris01/Raidplanner/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed static/*
var staticFiles embed.FS

var contentTypes = map[string]string{
	".html": "text/html; charset=utf-8",
	".css":  "text/css; charset=utf-8",
	".js":   "text/javascript; charset=utf-8",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".json": "application/json",
}

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

	router.Get("/api/roasters", apiCfg.handlerGetRoaster)
	router.Post("/api/roasters", apiCfg.AddRoasterHandler)
	router.Delete("/api/roasters/*", apiCfg.handlerDeleteRoaster)
	router.Get("/roasters/*", apiCfg.handlerShowRoaster)
	router.Get("/*", apiCfg.handlerStatic)

	log.Printf("Server running and waiting for requests\n")
	http.ListenAndServe(":"+apiCfg.port, router)

	fmt.Println(apiCfg)

}

// Define the Roaster struct to match the sample data
type Roaster struct {
	Username  string    `json:"username"`
	Server    string    `json:"server"`
	CreatedAt time.Time `json:"created_at"`
}

func (cfg apiConfig) handlerGetRoaster(w http.ResponseWriter, r *http.Request) {

	roasters, err := cfg.db.GetRoasters(r.Context())
	if err != nil {
		fmt.Print(err)
		respondWithError(w, 500, "Failed to fetch Roasters!", err)
		return
	}

	// Respond with the JSON data
	respondWithJSON(w, http.StatusOK, roasters)

}

func (cfg apiConfig) handlerShowRoaster(w http.ResponseWriter, r *http.Request) {

	// Respond with the JSON data
	respondWithJSON(w, http.StatusOK, roasters)

}


func (cfg apiConfig) AddRoasterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody struct {
		Username string `json:"username"`
		Server   string `json:"server"`
	}

	// Decode the JSON payload from the body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	fmt.Print(err)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call AddRoaster to insert the new roaster into the database
	roasterID, err := cfg.db.AddRoaster(r.Context(), requestBody.Username, requestBody.Server)
	if err != nil {
		fmt.Print(err)
		respondWithError(w, 500, "Failed to add Roaster!", err)
		return
	}

	// Respond with a success message
	response := map[string]interface{}{
		"success":    true,
		"roaster_id": roasterID,
	}

	// Send the JSON response
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg apiConfig) handlerDeleteRoaster(w http.ResponseWriter, r *http.Request) {
	// Extract RoasterID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	roasterID, err := strconv.Atoi(parts[3]) // This assumes the URL is "/api/roasters/{id}"
	if err != nil {
		fmt.Print(err)
		respondWithError(w, 500, "Failed to identify ID to delete!", err)
		return
	}

	fmt.Print(roasterID)

	err = cfg.db.DeleteRoaster(r.Context(), roasterID)
	if err != nil {
		fmt.Print(err)
		respondWithError(w, 500, "Failed to delete Roaster!", err)
		return
	}

	// Send success response
	respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (cfg apiConfig) handlerStatic(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path
	log.Printf("Requested path: %s", filepath)
	if filepath == "/" {
		filepath = "/static/index.html"
	} else if !strings.HasPrefix(filepath, "/static/") {
		filepath = "/static" + filepath
	}
	log.Printf("Filepath to open: %s", strings.TrimPrefix(filepath, "/"))
	f, err := staticFiles.Open(strings.TrimPrefix(filepath, "/"))
	if err != nil {
		log.Printf("Error opening index.html: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	ext := strings.ToLower(path.Ext(filepath))
	w.Header().Set("Content-Type", contentTypes[ext])

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if _, err := io.Copy(w, f); err != nil {
		log.Printf("Error copying file to response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
