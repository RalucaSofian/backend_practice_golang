package models

import (
	"github.com/uptrace/bun"
)

// Possible Auth User Roles
type UserRole string

const (
	UserRole_User  UserRole = "USER"
	UserRole_Admin UserRole = "ADMIN"
)

// Auth User DB Model
type AuthUser struct {
	bun.BaseModel `bun:"table:auth_users"`

	Id           int      `bun:"id,pk,autoincrement"`
	Email        string   `bun:"email,notnull"`
	Password     string   `bun:"password,notnull"`
	Name         *string  `bun:"name"`
	Address      *string  `bun:"address"`
	Phone        *string  `bun:"phone"`
	Role         UserRole `bun:"role,notnull"`
	SearchVector string   `bun:"search_vector,notnull,scanonly"`
}

func (user AuthUser) IsModel() {}
