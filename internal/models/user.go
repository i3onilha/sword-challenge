package models

import "time"

type UserRole string

const (
	RoleManager    UserRole = "manager"
	RoleTechnician UserRole = "technician"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) IsManager() bool {
	return u.Role == RoleManager
}

func (u *User) IsTechnician() bool {
	return u.Role == RoleTechnician
}
