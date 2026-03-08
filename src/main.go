package main

import (
	"log"

	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/router"
	"simon-weij/wayland-recorder-backend/src/router/auth"
	"simon-weij/wayland-recorder-backend/src/router/videos"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1000 * 1024 * 1024,
	})
	app.Use(logger.New())

	app.Get("/", auth.Middleware, router.HelloWorld)

	authGroup := app.Group("/auth")
	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/refresh", auth.RefreshToken)

	videosGroup := app.Group("/videos")
	videosGroup.Post("/upload", auth.Middleware, videos.UploadVideo)
	videosGroup.Get("/get/:id", auth.Middleware, videos.ServeVideoById)

	database.InitialiseDatabase()

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
