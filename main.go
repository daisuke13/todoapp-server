package main

import (
	"log"
	"os"

	"github.com/daisuke13/todo-app/server/src/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// instantiate echo
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routing
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(handler.Config))
	api.GET("/tasks", handler.GetTasks)
	api.POST("/tasks", handler.CreateTask)
	api.PUT("/tasks/:id/completed", handler.UpdateTask)
	api.DELETE("/tasks/:id", handler.DeleteTask)

	// launch server
	e.Start(":1313")
}
