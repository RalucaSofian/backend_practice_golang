package models

import "github.com/uptrace/bun"

// Pet DB Model
type Pet struct {
	bun.BaseModel `bun:"table:pets"`

	Id           int      `bun:"id,pk,autoincrement"`
	Name         string   `bun:"name,notnull"`
	Species      *string  `bun:"species"`
	Gender       *string  `bun:"gender"`
	Age          *float32 `bun:"age"`
	Description  *string  `bun:"description"`
	SearchVector string   `bun:"search_vector,notnull,scanonly"`
}

func (pet Pet) IsModel() {}
