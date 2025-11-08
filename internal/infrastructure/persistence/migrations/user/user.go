package user

import (
	"github.com/mansonxasthur/book-collection/pkg/migration"
)

const usersTableName = "users"

var (
	_ migration.DomainMigration = (*Migration)(nil)
	_ migration.Migration       = (*CreateUserTable)(nil)
)

type Migration struct{}

func NewMigration() *Migration {
	return &Migration{}
}

func (Migration) Migrations() []migration.Migration {
	return []migration.Migration{
		NewCreateUserTableMigration(),
	}
}
