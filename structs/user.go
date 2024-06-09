package structs

type User struct {
	ID       int    `json:"uid,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
}
