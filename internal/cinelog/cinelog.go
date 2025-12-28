package cinelog

import (
	"os"
	"time"

	_ "github.com/Cinelog-tv/backend/docs"
	"github.com/Cinelog-tv/backend/internal/controllers"
	"github.com/Cinelog-tv/backend/internal/database"
	"github.com/Cinelog-tv/backend/internal/models"
	"github.com/Cinelog-tv/backend/internal/services"
	"github.com/Cinelog-tv/backend/pkg/tmdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Execute(db database.Database, tmdbApikey string) error {

	if err := models.AutoMigrate(db.Database()); err != nil {
		return err
	}

	c := cache.New(5*time.Minute, 10*time.Minute)
	tmdbClient := tmdb.NewTmdbClient(tmdbApikey)

	mediaService := services.NewMediaService(tmdbClient, c)
	watchedService := services.NewWatchedService(db.Database())

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}

	api := e.Group("/api/v1")

	// Controllers
	controllers.RegisterAuthController(api, db.Database(), jwtSecret)
	controllers.RegisterMovieController(api, mediaService)
	controllers.RegisterTVController(api, mediaService)
	controllers.RegisterWatchedController(api, watchedService)
	controllers.RegisterGenreController(api, tmdbClient, c)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))

	return nil

}
