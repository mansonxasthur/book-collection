package main

import (
	"net/http"

	usercommand "github.com/mansonxasthur/book-collection/internal/application/command/user"
	handlers "github.com/mansonxasthur/book-collection/internal/infrastructure/http"
	"github.com/mansonxasthur/book-collection/internal/infrastructure/persistence/inmemory"
	"github.com/mansonxasthur/book-collection/pkg/cqrs"
	"github.com/mansonxasthur/book-collection/pkg/grace"
)

type App struct {
	Address string
	srv     *http.Server
}

func NewApp(address string) *App {
	return &App{
		Address: address,
	}
}

func (a *App) Bootstrap() {
	router := http.NewServeMux()
	commandBus := cqrs.NewInMemoryBus()
	userRepo := inmemory.NewUserRepository()
	createUserCommandHandler := usercommand.NewCreateUserHandler(userRepo)
	commandBus.Register(usercommand.CreateUserCommandName, createUserCommandHandler)
	userHandler := handlers.NewUserHandler(commandBus)
	userHandler.Handle(router)

	a.srv = &http.Server{
		Addr:    a.Address,
		Handler: router,
	}
}

func (a *App) Run() {
	grace.HandleShutdown(a.srv)
}
