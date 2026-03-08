package router

import "github.com/gofiber/fiber/v3"

// /
func HelloWorld(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
