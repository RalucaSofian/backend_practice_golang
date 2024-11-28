package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[up migration] Create foster Table")

		_, err := db.NewRaw(`
			CREATE TABLE IF NOT EXISTS foster (
				id            BIGSERIAL PRIMARY KEY,
				user_id       INTEGER REFERENCES auth_users(id) ON DELETE CASCADE,
				description   VARCHAR(100),
				pet_id        INTEGER REFERENCES pets(id) ON DELETE CASCADE,
				start_date    TIMESTAMPTZ NOT NULL,
				end_date      TIMESTAMPTZ,
				search_vector tsvector GENERATED ALWAYS AS (
					to_tsvector('simple', coalesce(user_id::text, '')) || ' ' ||
					to_tsvector('simple', coalesce(description, '')) || ' ' ||
					to_tsvector('simple', coalesce(pet_id::text, ''))
				) STORED
			)
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, errIdx := db.NewRaw(`CREATE INDEX IF NOT EXISTS idx_foster_search ON foster USING GIN(search_vector);`).Exec(ctx)
		if errIdx != nil {
			fmt.Println(errIdx)
			return errIdx
		}

		return nil

	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[down migration] Drop foster Table")

		_, err := db.NewRaw(`
			DROP TABLE IF EXISTS foster;
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
