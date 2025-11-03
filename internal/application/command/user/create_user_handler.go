package user

import (
	"context"

	domainuser "github.com/mansonxasthur/book-collection/internal/domain/user"
	"github.com/mansonxasthur/book-collection/internal/ports/output"
	"github.com/mansonxasthur/book-collection/pkg/cqrs"
)

type CreateUserHandler struct {
	repo output.UserRepository
}

func NewCreateUserHandler(repo output.UserRepository) cqrs.CommandHandler[CreateUserCommand] {
	return &CreateUserHandler{
		repo: repo,
	}
}

func (h CreateUserHandler) Handle(ctx context.Context, command CreateUserCommand) error {
	user, err := domainuser.NewUser(
		command.Name,
		command.Email,
		command.Password,
	)

	if err != nil {
		return err
	}

	return h.repo.Save(ctx, user)
}
