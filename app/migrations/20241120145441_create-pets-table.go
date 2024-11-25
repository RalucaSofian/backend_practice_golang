package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[up migration] Create pets Table")

		_, err := db.NewRaw(`
			CREATE TABLE IF NOT EXISTS pets (
				id            BIGSERIAL PRIMARY KEY,
				name          VARCHAR(50) NOT NULL,
				species       VARCHAR(50),
				gender        VARCHAR(1),
				age           REAL,
				description   VARCHAR(100),
				search_vector tsvector GENERATED ALWAYS AS (
					to_tsvector('simple', coalesce(name, '')) || ' ' ||
					to_tsvector('simple', coalesce(species, '')) || ' ' ||
					to_tsvector('simple', coalesce(gender, '')) || ' ' ||
					to_tsvector('simple', coalesce(description, ''))
				) STORED
			);
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, errIdx := db.NewRaw(`CREATE INDEX IF NOT EXISTS idx_pets_search ON pets USING GIN(search_vector);`).Exec(ctx)
		if errIdx != nil {
			fmt.Println(errIdx)
			return errIdx
		}

		return nil

	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[down migration] Drop pets Table")

		_, err := db.NewRaw(`
			DROP TABLE IF EXISTS pets;
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
