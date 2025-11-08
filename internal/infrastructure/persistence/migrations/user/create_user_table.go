package user

import (
	"database/sql"

	"github.com/mansonxasthur/book-collection/pkg/migration"
)

type CreateUserTable struct {
	migration.BaseMigration
}

func NewCreateUserTableMigration() *CreateUserTable {
	return &CreateUserTable{
		BaseMigration: migration.BaseMigration{
			TableNameValue:   usersTableName,
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
