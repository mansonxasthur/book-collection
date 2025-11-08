package migrations

import (
	"database/sql"

	"github.com/mansonxasthur/book-collection/pkg/migration"
)

const tableName = "users"

type UserMigration struct{}

func NewUserMigration() *UserMigration {
	return &UserMigration{}
}

func (UserMigration) Migrations() []migration.Migration {
	return []migration.Migration{
		NewCreateUserTableMigration(),
	}
}

type CreateUserTable struct {
	migration.BaseMigration
}

func NewCreateUserTableMigration() *CreateUserTable {
	return &CreateUserTable{
		BaseMigration: migration.BaseMigration{
			TableNameValue:   tableName,
			DescriptionValue: "Create users table",
			CreatedAtValue:   "2025-011-08 20:18:00",
		},
	}
}

func (CreateUserTable) Up(dbConn *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    );`

	_, err := dbConn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (CreateUserTable) Down(dbConn *sql.DB) error {
	query := `DROP TABLE IF EXISTS users;`
	_, err := dbConn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
