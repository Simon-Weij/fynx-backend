package main

import (
	"log"

	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/router"
	"simon-weij/wayland-recorder-backend/src/router/auth"
	"simon-weij/wayland-recorder-backend/src/router/user"
	"simon-weij/wayland-recorder-backend/src/router/videos"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1000 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
	}))

	api := app.Group("/api")

	api.Get("/", auth.Middleware, router.HelloWorld)

	authGroup := api.Group("/auth")
	authGroup.Post("/signup", auth.Signup)
	authGroup.Post("/login", auth.Login)
	authGroup.Post("/refresh", auth.RefreshToken)
	authGroup.Post("/logout", auth.Logout)

	videosGroup := api.Group("/videos")
	videosGroup.Post("/upload", auth.Middleware, videos.UploadVideo)
	videosGroup.Get("/get/:id", auth.Middleware, videos.ServeVideoById)

	userGroup := api.Group("/user")
	userGroup.Get("/data", auth.Middleware, user.UserData)

	database.InitialiseDatabase()

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
