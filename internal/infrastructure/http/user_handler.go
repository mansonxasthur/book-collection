package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mansonxasthur/book-collection/internal/application/command/user"
	userdomain "github.com/mansonxasthur/book-collection/internal/domain/user"
	"github.com/mansonxasthur/book-collection/internal/infrastructure/http/requests"
	"github.com/mansonxasthur/book-collection/internal/infrastructure/http/responses"
	"github.com/mansonxasthur/book-collection/pkg/cqrs"
	"github.com/mansonxasthur/book-collection/pkg/response"
)

type UserHandler struct {
	commandBus cqrs.CommandBus
}

func NewUserHandler(commandBus cqrs.CommandBus) *UserHandler {
	return &UserHandler{
		commandBus: commandBus,
	}
}

func (h *UserHandler) Handle(router *http.ServeMux) {
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("POST /", h.createUser)
	userGroup := http.NewServeMux()
	userGroup.Handle("/users/", http.StripPrefix("/users", userRouter))
	router.Handle("/", userGroup)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var userRequest requests.CreateUserRequest
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("error closing request body: %v\n", err)
		}
	}()

	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	command := user.CreateUserCommand{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	err = h.commandBus.Dispatch(r.Context(), command)
	if err != nil {
		var status response.Status
		switch {
		case errors.Is(err, userdomain.ErrNameRequired), errors.Is(err, userdomain.ErrEmailRequired), errors.Is(err, userdomain.ErrPasswordRequired):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		response.Error(w, fmt.Errorf("falied to create user: %w", err), status)
		return
	}

	response.Send(w, responses.CreateUserResponse{
		Success: true,
	}, http.StatusCreated)
}
