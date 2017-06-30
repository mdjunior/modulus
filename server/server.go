package main

import (
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mdjunior/modulus/server/handlers"
	customMiddleware "github.com/mdjunior/modulus/server/middleware"
	"github.com/patrickmn/go-cache"
)

func main() {
	e := echo.New()
	c := cache.New(10*time.Minute, 5*time.Minute)

	// DB
	dbHandler := handlers.NewDBHandler(c)

	// Middlewares
	e.Use(middleware.RequestID())
	e.Use(customMiddleware.RequestLogDataMiddleware())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Healthcheck
	e.GET("/status", handlers.Status)

	// List games
	e.GET("/games", dbHandler.GamesList)
	e.POST("/games", dbHandler.GamesCreate)
	e.GET("/games/:id", dbHandler.GamesGET)
	e.POST("/games/:id/join", dbHandler.GamesJoin)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
