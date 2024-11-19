package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[up migration] Create auth_users Table")

		_, err := db.NewRaw(`
			CREATE TABLE IF NOT EXISTS auth_users (
				id       BIGSERIAL PRIMARY KEY,
				email    VARCHAR(50) NOT NULL,
				password VARCHAR(256) NOT NULL,
				name     VARCHAR(50),
				address  VARCHAR(100),
				phone    VARCHAR(30)
			);
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil

	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[down migration] Drop auth_users Table")

		_, err := db.NewRaw(`DROP TABLE IF EXISTS auth_users`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
