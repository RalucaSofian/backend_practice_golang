package main

import (
	"app/db"
	"app/migrations"
	"context"
	"fmt"
	"os"

	"github.com/uptrace/bun/migrate"
)

type cliCommand struct {
	Description string
	Action      func(params []string)
	Params      string
}

var cliCommands = make(map[string]cliCommand)

func init() {
	cliCommands["create-migration"] = cliCommand{Description: "Create Migration with given name", Action: func(params []string) { handleCreateMigration(params[0]) }, Params: "migration name: string"}
	cliCommands["apply-migrations"] = cliCommand{Description: "Apply all unapplied Migrations", Action: func(params []string) { handleApplyMigrations() }, Params: "none"}
	cliCommands["revert-migrations"] = cliCommand{Description: "Revert most recently applied Migrations", Action: func(params []string) { handleRevertMigrations() }, Params: "none"}
	cliCommands["list-migrations"] = cliCommand{Description: "List all Migrations", Action: func(params []string) { handleListMigrations() }, Params: "none"}
	cliCommands["help"] = cliCommand{Description: "CLI Migrations Help", Action: func(params []string) { handleHelp() }, Params: "none"}
}

func getMigrator() *migrate.Migrator {
	migrator := migrate.NewMigrator(db.GetConn(), migrations.Migrations, migrate.WithMarkAppliedOnSuccess(true))
	migrator.Init(context.Background())
	return migrator
}

func handleCLI() {
	fmt.Println("-- [CLI] -- ")
	fmt.Print("\n")

	args := os.Args[1:]
	command, ok := cliCommands[args[0]]
	if !ok {
		fmt.Println("[CLI] Invalid Command Line Argument")
	} else {
		fmt.Println("[CLI] Command:", command.Description)
		command.Action(args[1:])
	}

	fmt.Print("\n")
	fmt.Println("-- [END] --")
}

func handleCreateMigration(migrationName string) {
	fmt.Println("[CLI] Handle Create Migration")

	migrator := getMigrator()
	migrationFile, err := migrator.CreateGoMigration(context.Background(), migrationName)
	if err != nil {
		fmt.Println("[CLI] Error occurred when Creating Migration:", err)
		return
	}
	fmt.Println("[CLI] Created Migration:", migrationFile.Name)
}

func handleApplyMigrations() {
	fmt.Println("[CLI] Handle Apply Migrations")

	migrator := getMigrator()
	err := migrator.Lock(context.Background())
	if err != nil {
		fmt.Println("[CLI]", err)
		return
	}
	defer migrator.Unlock(context.Background())

	group, err := migrator.Migrate(context.Background())
	if err != nil {
		fmt.Println("[CLI] Error occurred when Applying Migrations:", err)
		return
	}

	if group.IsZero() {
		fmt.Println("[CLI] No Migrations to Apply. DB is up to date")
	} else {
		fmt.Println("[CLI] Applied Migrations:", group)
	}
}

func handleRevertMigrations() {
	fmt.Println("[CLI] Handle Revert Migrations")

	migrator := getMigrator()
	err := migrator.Lock(context.Background())
	if err != nil {
		fmt.Println("[CLI]", err)
		return
	}
	defer migrator.Unlock(context.Background())

	group, err := migrator.Rollback(context.Background())
	if err != nil {
		fmt.Println("[CLI] Error occurred when Reverting Migrations:", err)
		return
	}

	if group.IsZero() {
		fmt.Println("[CLI] No Migrations to Revert")
	} else {
		fmt.Println("[CLI] Reverted Migrations:", group)
	}
}

func handleListMigrations() {
	fmt.Println("[CLI] Handle List Migrations")

	migrator := getMigrator()

	slice, err := migrator.MigrationsWithStatus(context.Background())
	if err != nil {
		fmt.Println("[CLI] Error occurred when retrieving Migrations List:", err)
		return
	}

	fmt.Println("[CLI] List of Applied Migrations is:", slice.Applied())
	fmt.Println("[CLI] List of Unapplied Migrations is:", slice.Unapplied())
}

func handleHelp() {
	fmt.Println("[CLI] Handle Help")

	fmt.Print("[CLI] Available CLI commands\n\n")
	for i, cmd := range cliCommands {
		fmt.Println("[CLI]", i, ":", cmd.Description, "; Parameters:", cmd.Params)
	}
}
