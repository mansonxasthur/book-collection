package user

type Password string

func (p Password) String() string {
	return "[REDACTED]"
}

func (p Password) GoString() string {
	return "[REDACTED]"
}

func (p Password) Value() string {
	return string(p)
}
