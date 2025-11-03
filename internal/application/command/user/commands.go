package user

const CreateUserCommandName = "CreateUserCommand"

type CreateUserCommand struct {
	Name     string
	Email    string
	Password string
}

func (c CreateUserCommand) CommandName() string {
	return CreateUserCommandName
}
