package output

import (
	"context"

	"github.com/mansonxasthur/book-collection/internal/domain/user"
)

type UserRepository interface {
	Save(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
}
