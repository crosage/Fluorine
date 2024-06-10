package handlers

import (
	"awesomeProject/database"
	"awesomeProject/structs"
	"awesomeProject/utils"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"strconv"
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
			"roles":    datauser.Roles,
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
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "密码错误"})
	}
}

func getAllUsers(ctx *fiber.Ctx) error {
	hasPermission := validatePermission(ctx)
	if !hasPermission {
		return ctx.Status(403).JSON(fiber.Map{"error": "无权限"})
	}
	users, err := database.GetAllUsers()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching users")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "获取用户信息失败"})
	}

	var usersJSON []map[string]interface{}
	for _, user := range users {
		userJSON := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"roles":    user.Roles,
		}
		usersJSON = append(usersJSON, userJSON)
	}

	return ctx.Status(fiber.StatusOK).JSON(usersJSON)
}

func userRegister(ctx *fiber.Ctx) error {
	var user structs.User
	err := jsoniter.Unmarshal(ctx.Body(), &user)

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "用户密码处理失败"})
	}
	user.Password = hashedPassword

	if err != nil || len(user.Username) == 0 || len(user.Password) == 0 {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "无效的用户名或密码"})
	}

	id, err := database.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "创建用户失败"})
	}

	role, err := database.GetRoleByName("普通员工")
	fmt.Println(err)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "获取角色失败"})
	}

	err = database.AddRoleToUser(int(id), role.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "分配角色失败"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "用户注册成功"})
}

func AddRoleToUserHandler(c *fiber.Ctx) error {

	userID, err := strconv.Atoi(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var request struct {
		RoleID int `json:"rid"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	err = database.AddRoleToUser(userID, request.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add role to user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role added to user successfully",
	})
}
