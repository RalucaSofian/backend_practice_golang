package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[up migration] Add search_vector Column to auth_users Table")

		_, err := db.NewRaw(`
			ALTER TABLE auth_users
				ADD IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS (
					to_tsvector('simple', coalesce(email, '')) || ' ' ||
					to_tsvector('simple', coalesce(name, '')) || ' ' ||
					to_tsvector('simple', coalesce(address, '')) || ' ' ||
					to_tsvector('simple', coalesce(phone, ''))
				) STORED;
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, errIdx := db.NewRaw(`CREATE INDEX IF NOT EXISTS idx_search ON auth_users USING GIN(search_vector);`).Exec(ctx)
		if errIdx != nil {
			fmt.Println(errIdx)
			return errIdx
		}

		return nil

	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[down migration] Drop search_vector Column from auth_users Table")

		_, errIdx := db.NewRaw(`DROP INDEX IF EXISTS idx_search;`).Exec(ctx)
		if errIdx != nil {
			fmt.Println(errIdx)
			return errIdx
		}
		_, err := db.NewRaw(`ALTER TABLE auth_users DROP COLUMN IF EXISTS search_vector`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
