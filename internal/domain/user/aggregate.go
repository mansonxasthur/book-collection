package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/mansonxasthur/book-collection/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	domain.BaseAggregate
	id        string
	name      string
	email     string
	password  Password
	createdAt string
	updatedAt string
}

func createPasswordHash(password string) (Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return Password(hash), nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password.Value()), []byte(password))
	return err == nil
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	if email == "" {
		return nil, ErrEmailRequired
	}

	if password == "" {
		return nil, ErrPasswordRequired
	}

	id := uuid.New().String()
	hashedPassword, err := createPasswordHash(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dateTime := now.Format(time.DateTime)
	user := &User{
		id:        id,
		name:      name,
		email:     email,
		password:  hashedPassword,
		createdAt: dateTime,
		updatedAt: dateTime,
	}

	user.ApplyEvent(&CreateEvent{
		userID: id,
		BaseEvent: domain.BaseEvent{
			NameVal:       "user.create",
			OccurredAtVal: now,
			AggregateVal:  id,
			VersionVal:    "1",
		},
	})

	return user, nil
}

func (u *User) ID() string        { return u.id }
func (u *User) Name() string      { return u.name }
func (u *User) Email() string     { return u.email }
func (u *User) CreatedAt() string { return u.createdAt }
func (u *User) UpdatedAt() string { return u.updatedAt }
