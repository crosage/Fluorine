package handlers

import (
	"awesomeProject/database"
	"awesomeProject/structs"
	"github.com/gofiber/fiber/v2"
)

func getAllRoles(ctx *fiber.Ctx) error {
	var roles []structs.Role
	roles, err := database.GetAllRoles()
	print(database.GetAllRoles())
	print(roles)
	print(err)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"total": len(roles), "roles": roles})
}
