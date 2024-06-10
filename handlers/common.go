package handlers

import (
	"awesomeProject/structs"
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

func InitHandlers(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Post("/api/user/login", userLogin)
	app.Post("/api/user", userRegister)
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
	app.Post("/api/role/:uid", addRoleToUserHandler)
	app.Delete("/api/role/:uid", removeRoleFromUserHandler)
	app.Get("/api/user", getAllUsers)
	app.Post("/api/roles/:uid", addRolesToUserHandler)
}

var JwtSecret = "RTRT"

func validatePermission(ctx *fiber.Ctx) bool {
	userLocal := ctx.Locals("user").(*jwt.Token)
	claims := userLocal.Claims.(jwt.MapClaims)
	rolesInterface := claims["roles"].([]interface{})

	var userType []structs.Role
	for _, roleInterface := range rolesInterface {
		roleMap := roleInterface.(map[string]interface{})
		role := structs.Role{
			ID:   int(roleMap["id"].(float64)), // 注意：JSON数字在反序列化后是float64类型
			Name: roleMap["name"].(string),
		}

		permissionsInterface := roleMap["permissions"].([]interface{})
		for _, permissionInterface := range permissionsInterface {
			permissionMap := permissionInterface.(map[string]interface{})
			permission := structs.Permission{
				ID:   int(permissionMap["id"].(float64)), // 同样的，JSON数字在反序列化后是float64类型
				Name: permissionMap["name"].(string),
			}
			role.Permissions = append(role.Permissions, permission)
		}

		userType = append(userType, role)
	}

	fmt.Println(userType)
	return userType != nil
}

func getSessionUser(ctx *fiber.Ctx) structs.User {
	userLocal := ctx.Locals("user").(*jwt.Token)
	claims := userLocal.Claims.(jwt.MapClaims)
	user := structs.User{}
	user.ID = int(claims["id"].(float64))
	user.Username = claims["username"].(string)

	var userType []structs.Role
	rolesInterface := claims["roles"].([]interface{})
	for _, roleInterface := range rolesInterface {
		roleMap := roleInterface.(map[string]interface{})
		role := structs.Role{
			ID:   int(roleMap["id"].(float64)), // 注意：JSON数字在反序列化后是float64类型
			Name: roleMap["name"].(string),
		}

		permissionsInterface := roleMap["permissions"].([]interface{})
		for _, permissionInterface := range permissionsInterface {
			permissionMap := permissionInterface.(map[string]interface{})
			permission := structs.Permission{
				ID:   int(permissionMap["id"].(float64)), // 同样的，JSON数字在反序列化后是float64类型
				Name: permissionMap["name"].(string),
			}
			role.Permissions = append(role.Permissions, permission)
		}

		userType = append(userType, role)
	}

	user.Roles = userType
	return user
}
