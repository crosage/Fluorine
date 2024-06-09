package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var db *sql.DB

func InitDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to open database")
	}
	createTables()
	initRolesandPermissions()
}
func createTables() {
	createPermissionsTableSQL := `
	CREATE TABLE IF NOT EXISTS permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	createRolesTableSQL := `
	CREATE TABLE IF NOT EXISTS roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	createRolePermissionsTableSQL := `
	CREATE TABLE IF NOT EXISTS role_permissions (
		role_id INTEGER NOT NULL,
		permission_id INTEGER NOT NULL,
		FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE ON UPDATE CASCADE,
		PRIMARY KEY (role_id, permission_id)
	);`

	createUserRolesTableSQL := `
	CREATE TABLE IF NOT EXISTS user_roles (
		user_id INTEGER NOT NULL,
		role_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
		PRIMARY KEY (user_id, role_id)
	);`

	createCheckInTableSQL := `
	CREATE TABLE IF NOT EXISTS check_ins (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		check_in_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
	);
	`

	_, err := db.Exec(createPermissionsTableSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create permissions table")
	}

	_, err = db.Exec(createRolesTableSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create roles table")
	}

	_, err = db.Exec(createUsersTableSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create users table")
	}

	_, err = db.Exec(createRolePermissionsTableSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create role_permissions table")
	}

	_, err = db.Exec(createUserRolesTableSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create user_roles table")
	}

	_, err = db.Exec(createCheckInTableSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create check_in table")
	}
}

func initRolesandPermissions() {
	roles := map[string][]string{
		"普通员工":  {"查看普通日志"},
		"审计人员":  {"查看usb连接", "查看安全日志", "查看用户登录记录"},
		"系统管理员": {"管理权限"},
	}

	for roleName, permissions := range roles {
		// 插入角色
		_, err := db.Exec("INSERT OR IGNORE INTO roles (name) VALUES (?)", roleName)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to insert role: %s", roleName)
		}

		// 获取插入的角色的ID
		var roleId int
		err = db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleId)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to get role ID for role: %s", roleName)
		}

		// 插入权限并与角色关联
		for _, permission := range permissions {
			_, err := db.Exec(`
				INSERT OR IGNORE INTO permissions (name) VALUES (?);
				INSERT OR IGNORE INTO role_permissions (role_id, permission_id) 
				SELECT ?, id FROM permissions WHERE name = ?;
			`, permission, roleId, permission)
			if err != nil {
				log.Fatal().Err(err).Msgf("Failed to insert permission: %s", permission)
			}
		}
	}
}
