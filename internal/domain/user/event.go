package user

import "github.com/mansonxasthur/book-collection/internal/domain"

type CreateEvent struct {
	domain.BaseEvent
	userID string
}

func (e *CreateEvent) UserID() string { return e.userID }
