// package models

// type Role string

// const (
// 	RoleAdmin Role = "admin"
// 	RoleUser  Role = "user"
// )

// type User struct {
// 	ID       string `json:"id" db:"id"`
// 	Name     string `json:"name" db:"name"`
// 	Email    string `json:"email" db:"email"`
// 	Password string `json:"-" db:"password"`
// 	Role     Role   `json:"role" db:"role"`
// }
package models

// import (
	// "gorm.io/gorm"
// )

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     Role   `json:"role" gorm:"not null;default:user"`
}