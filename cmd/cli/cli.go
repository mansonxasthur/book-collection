package main

type Command interface {
	Signature() string
	Description() string
	Execute(*Config, []string) error
}

type BaseCommand struct {
	SignatureValue   string
	DescriptionValue string
}

func (b BaseCommand) Signature() string {
	return b.SignatureValue
}

func (b BaseCommand) Description() string {
	return b.DescriptionValue
}
