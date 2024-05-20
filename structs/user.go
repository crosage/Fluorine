package structs

type User struct {
	ID       int
	Username string
	Password string
	Roles    []Role
}
