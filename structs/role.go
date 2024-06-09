package structs

type Role struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty"`
}
