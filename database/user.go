package database

import (
	"awesomeProject/structs"
	"database/sql"
)

func CreateUser(user structs.User) (int, error) {
	stmt, err := db.Prepare("INSERT INTO users (`username`,`password`) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func DeleteUser(userID int) error {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(user structs.User) error {
	stmt, err := db.Prepare("UPDATE users SET username = ?, password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Username, user.Password, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserRole(userID int, roles []structs.Role) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM user_roles WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, role := range roles {
		_, err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, role.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func AddRoleToUser(userID, roleID int) error {
	_, err := db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID)
	return err
}

func RemoveRoleFromUser(userID, roleID int) error {
	_, err := db.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID)
	return err
}

func GetRolesByUsername(username string) ([]structs.Role, error) {
	var roles []structs.Role
	query := `
		SELECT r.id, r.name
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		JOIN users u ON ur.user_id = u.id
		WHERE u.username = ?
	`
	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role structs.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}

		permissions, err := GetPermissionsByRole(role.Name)
		if err != nil {
			return nil, err
		}
		role.Permissions = permissions

		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func GetUserByID(userID int) (structs.User, error) {
	var user structs.User
	query := "SELECT id, username, password FROM users WHERE id = ?"
	row := db.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil // Return empty user object and nil error
		}
		return user, err
	}

	roles, err := GetRolesByUsername(user.Username)
	if err != nil {
		return user, err
	}
	user.Roles = roles

	return user, nil
}

func GetUserByUsername(username string) (structs.User, error) {
	var user structs.User
	query := "SELECT id,username,password FROM users WHERE username=?"
	row := db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, sql.ErrNoRows
		}
		return user, err
	}
	roles, err := GetRolesByUsername(user.Username)
	if err != nil {
		return user, err
	}
	user.Roles = roles

	return user, nil
}
func GetAllUsers() ([]structs.User, error) {
	var users []structs.User
	query := "SELECT id,username,password FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user structs.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		roles, err := GetRolesByUsername(user.Username)
		if err != nil {
			return nil, err
		}
		user.Roles = roles
		users = append(users, user)
	}
	return users, nil
}
