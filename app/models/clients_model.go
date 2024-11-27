package models

import (
	"github.com/uptrace/bun"
)

// Client DB Model
type Client struct {
	bun.BaseModel `bun:"table:clients"`

	Id           int       `bun:"id,pk,autoincrement"`
	UserId       *int      `bun:"user_id"`
	User         *AuthUser `bun:"rel:belongs-to,join:user_id=id"`
	Description  *string   `bun:"description"`
	SearchVector string    `bun:"search_vector,notnull,scanonly"`
}

func (pet Client) IsModel() {}
