package models

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	Role     Role   `json:"role" db:"role"`
}
