package database

import (
	"awesomeProject/structs"
	"database/sql"
	"errors"
)

func GetPermissionsByRole(roleName string) ([]structs.Permission, error) {
	var Permissions []structs.Permission
	rows, err := db.Query(`
		SELECT p.id,p.name 
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN roles r ON rp.role_id = r.id
		WHERE r.name = ?
	`, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var permission structs.Permission
		if err := rows.Scan(&permission.ID, &permission.Name); err != nil {
			return nil, err
		}
		Permissions = append(Permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Permissions, nil
}

func GetRoleByName(roleName string) (structs.Role, error) {
	var role structs.Role
	err := db.QueryRow("SELECT id, name FROM roles WHERE name = ?", roleName).Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return role, errors.New("角色不存在")
		}
		return role, err
	}
	return role, nil
}

func GetRoleByID(roleID int) (structs.Role, error) {
	var role structs.Role
	query := "SELECT id, name FROM roles WHERE id = ?"
	row := db.QueryRow(query, roleID)
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return role, nil
		}
		return role, err
	}

	permissions, err := GetPermissionsByRole(role.Name)
	if err != nil {
		return role, err
	}
	role.Permissions = permissions

	return role, nil
}

func GetAllRoles() ([]structs.Role, error) {
	var roles []structs.Role
	query := "SELECT id, name FROM roles"
	rows, err := db.Query(query)
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
