package database

import (
	"awesomeProject/structs"
)

func GetPermissionByID(permissionID int) (structs.Permission, error) {
	var permission structs.Permission
	query := "SELECT id, name, description FROM permissions WHERE id = ?"
	row := db.QueryRow(query, permissionID)
	err := row.Scan(&permission.ID, &permission.Name, &permission.Description)
	if err != nil {
		return permission, err
	}
	return permission, nil
}

func GetAllPermissions() ([]structs.Permission, error) {
	var permissions []structs.Permission
	query := "SELECT id, name, description FROM permissions"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission structs.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
