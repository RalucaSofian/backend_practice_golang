package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[up migration] Add roles Column to auth_users Table")

		_, err := db.NewRaw(`
		CREATE TYPE user_roles AS ENUM ('ADMIN', 'USER');
		ALTER TABLE auth_users
				ADD IF NOT EXISTS role USER_ROLES NOT NULL DEFAULT 'USER';
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil

	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Println("[down migration] Drop roles Column from auth_users Table")

		_, err := db.NewRaw(`
		DROP TYPE IF EXISTS user_roles;
		ALTER TABLE auth_users DROP COLUMN IF EXISTS role;
		`).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
