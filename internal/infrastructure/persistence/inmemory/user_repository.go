package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/mansonxasthur/book-collection/internal/domain/user"
	"github.com/mansonxasthur/book-collection/internal/ports/output"
)

var _ output.UserRepository = (*UserRepository)(nil)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserCanNotBeNil = errors.New("user can not be nil")
	ErrInvalidID       = errors.New("invalid id")
)

type UserRepository struct {
	users map[string]user.User
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]user.User),
	}
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if u == nil {
		return ErrUserCanNotBeNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID()] = *u

	fmt.Printf("User saved: %+v\n", u)
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if id == "" {
		return nil, ErrInvalidID
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	userCopy := u

	return &userCopy, nil
}
