package main

import (
	"awesomeProject/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	database.InitDatabase()

	app.Listen(":3000")
}
