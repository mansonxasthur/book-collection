package main

type Command interface {
	Signature() string
	Description() string
	Execute(*Config, []string) error
}
