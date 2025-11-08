package migration

import (
	"database/sql"
	"fmt"
	"strings"
)

const migrationsTableName = "migrations"

type DomainMigration interface {
	Migrations() []Migration
}

type Migration interface {
	TableName() string
	Description() string
	CreatedAt() string
	Name() string
	Up(*sql.DB) error
	Down(*sql.DB) error
}

type BaseMigration struct {
	TableNameValue   string
	CreatedAtValue   string
	DescriptionValue string
}

func (b BaseMigration) TableName() string {
	return b.TableNameValue
}

func (b BaseMigration) Description() string {
	return b.DescriptionValue
}

func (b BaseMigration) CreatedAt() string {
	return b.CreatedAtValue
}

func (b BaseMigration) Name() string {
	createdAt := strings.ReplaceAll(strings.ToLower(b.CreatedAt()), " ", "_")
	createdAt = strings.ReplaceAll(createdAt, "-", "")
	createdAt = strings.ReplaceAll(createdAt, ":", "")
	description := strings.ReplaceAll(strings.ToLower(b.Description()), " ", "_")
	return fmt.Sprintf("%s_%s", createdAt, description)
}
