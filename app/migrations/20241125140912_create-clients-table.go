package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[up migration] Create clients Table")

		_, err := db.NewRaw(`
			CREATE TABLE IF NOT EXISTS clients (
				id            BIGSERIAL PRIMARY KEY,
				user_id       INTEGER REFERENCES auth_users(id) ON DELETE CASCADE,
				description   VARCHAR(100),
				search_vector tsvector GENERATED ALWAYS AS (
					to_tsvector('simple', coalesce(user_id::text, '')) || ' ' ||
					to_tsvector('simple', coalesce(description, ''))
				) STORED
			);
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, errIdx := db.NewRaw(`CREATE INDEX IF NOT EXISTS idx_clients_search ON clients USING GIN(search_vector);`).Exec(ctx)
		if errIdx != nil {
			fmt.Println(errIdx)
			return errIdx
		}

		return nil

	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[down migration] Drop clients Table")

		_, err := db.NewRaw(`
			DROP TABLE IF EXISTS clients;
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
