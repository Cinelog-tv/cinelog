package main

import (
	"log"
	"os"

	"github.com/Cinelog-tv/backend/internal/cinelog"
	"github.com/Cinelog-tv/backend/internal/database"
)

// @title Swagger CineLog Swagger API
// @version 1.0
// @description Api for CineLog
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /api/v1
func main() {
	var db database.Database
	if os.Getenv("environment") == "" || os.Getenv("environment") == "development" {
		db = &database.SqliteDatabase{}
	}

	err := db.Connect()
	if err != nil {
		log.Fatalf("Cannot connect to database %+v", err)
	}

	tmdbApiKey := os.Getenv("TMDB_API_KEY")

	if tmdbApiKey == "" {
		log.Fatal("TMDB_API_KEY is not set")
	}

	cinelog.Execute(db, tmdbApiKey)
}
