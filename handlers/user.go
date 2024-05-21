package handlers

import (
	"awesomeProject/database"
	"awesomeProject/structs"
	"awesomeProject/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

func userLogin(ctx *fiber.Ctx) error {
	var user structs.User
	if err := jsoniter.Unmarshal(ctx.Body(), &user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Failed to parse request body"})
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "用户密码处理失败"})
	}
	user.Password = hashedPassword
	datauser, err := database.GetUserByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("No user")
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "没有该用户"})
		} else {
			log.Error().Err(err).Msg("Error login")
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "登录失败"})
		}
	}
	if datauser.Password == hashedPassword {
		claims := jwt.MapClaims{
			"id":       datauser.ID,
			"username": datauser.Username,
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JwtSecret))
		if err != nil {
			log.Error().Err(err).Msg("Error token")
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token生成失败"})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"id":       datauser.ID,
			"username": datauser.Username,
			"token":    token,
			"roles":    datauser.Roles,
		})
	} else {
		log.Error().Err(err).Msg("No user")
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"error": "密码错误"})
	}
}

func userRegister(ctx *fiber.Ctx) error {
	var user structs.User
	err := jsoniter.Unmarshal(ctx.Body(), &user)
	if err != nil || len(user.Username) == 0 || len(user.Password) == 0 {
		return ctx.Status(403).JSON(nil)
	}
	id, err := database.CreateUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "创建用户失败"})
	}
	err = database.UpdateUserRole(id, []structs.Role{*structs.GetRoleByName("普通用户")})
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON("")
}
