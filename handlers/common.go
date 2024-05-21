package handlers

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitHandlers(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Post("/api/register", userLogin)
	app.Post("/api/register", userRegister)
	app.Get("/api/role", getAllRoles)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(JwtSecret),
		},
		ContextKey: "user",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "未授权的访问"})
		},
	}))
}

var JwtSecret = "RTRT"
