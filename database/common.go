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
        name TEXT NOT NULL UNIQUE,
        description TEXT
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
        FOREIGN KEY (role_id) REFERENCES roles(id),
        FOREIGN KEY (permission_id) REFERENCES permissions(id),
        PRIMARY KEY (role_id, permission_id)
    );`

	createUserRolesTableSQL := `
    CREATE TABLE IF NOT EXISTS user_roles (
        user_id INTEGER NOT NULL,
        role_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (role_id) REFERENCES roles(id),
        PRIMARY KEY (user_id, role_id)
    );`

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
}

func initRolesandPermissions() {
	roles := map[string][]string{
		"普通员工": {"查看薪水", "签到打卡", "申请请假"},
		"部门组长": {"补打卡", "编辑信息", "审批请假"},
		"人事管理": {"管理薪水", "编辑信息"},
		"总经理":  {"管理权限"},
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
