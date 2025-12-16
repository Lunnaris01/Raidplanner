package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	//dbQueries *database.Queries
	dbURL    string
	dbToken  string
	platform string
}

func main() {
	fmt.Println("Civ API started!")
	godotenv.Load()

	env_platform := os.Getenv("PLATFORM")
	env_dbURL := os.Getenv("TURSO_DATABASE_URL")
	env_dbToken := os.Getenv("TURSO_AUTH_TOKEN")
	dbCombinedURL := env_dbURL + "?authToken=" + env_dbToken
	db, err := sql.Open("libsql", dbCombinedURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbCombinedURL, err)
		os.Exit(1)
	} else {
		fmt.Printf("Connected to Database\n")
	}

	defer db.Close()

	apiCfg := apiConfig{
		dbURL:    env_dbURL,
		dbToken:  env_dbToken,
		platform: env_platform,
	}

	fmt.Println(apiCfg)

}
