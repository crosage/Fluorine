package structs

type Role struct {
	ID          int
	Name        string
	Permissions []Permission
}
