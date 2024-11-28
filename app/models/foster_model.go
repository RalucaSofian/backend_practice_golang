package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Foster DB Model
type Foster struct {
	bun.BaseModel `bun:"table:foster"`

	Id           int        `bun:"id,pk,autoincrement"`
	UserId       *int       `bun:"user_id"`
	User         *AuthUser  `bun:"rel:belongs-to,join:user_id=id"`
	Description  *string    `bun:"description"`
	PetId        *int       `bun:"pet_id"`
	Pet          *Pet       `bun:"rel:belongs-to,join:pet_id=id"`
	StartDate    time.Time  `bun:"start_date,notnull"`
	EndDate      *time.Time `bun:"end_date"`
	SearchVector string     `bun:"search_vector,notnull,scanonly"`
}

func (foster Foster) IsModel() {}
