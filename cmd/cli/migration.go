package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/mansonxasthur/book-collection/internal/infrastructure/persistence/migrations"
	"github.com/mansonxasthur/book-collection/pkg/migration"
)

type MigrationCommand struct {
}

var _ Command = (*MigrationCommand)(nil)

func NewMigrationCommand() *MigrationCommand {
	return &MigrationCommand{}
}

func (c MigrationCommand) Signature() string {
	return "migration"
}

func (c MigrationCommand) Description() string {
	return "Run migration operations"
}

func (c MigrationCommand) Execute(config *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("required commands: migrate or rollback")
	}

	conn, err := c.getConnection(config)
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("error closing database connection: %v\n", err)
		}
	}()

	migrator := c.getMigrator(config, conn)
	migrator.Register(
		migrations.NewUserMigration(),
	)

	cmd := args[0]
	switch cmd {
	case "migrate":
		return migrator.Migrate()
	case "rollback":
		return migrator.Rollback()
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (c MigrationCommand) getMigrator(config *Config, conn *sql.DB) *migration.Migrator {
	if config.DB.Driver != postgresDriver {
		log.Fatalf("unsupported database driver: %s\n", config.DB.Driver)
	}

	return migration.NewMigrator(conn)
}

func (c MigrationCommand) getConnection(config *Config) (*sql.DB, error) {
	dbConfig := config.DB
	sslMode := "enable"
	if config.ENV == envDev {
		sslMode = "disable"
	}
	connection := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.Driver,
		dbConfig.User,
		dbConfig.Pass,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DB,
		sslMode,
	)

	return sql.Open(dbConfig.Driver, connection)
}
