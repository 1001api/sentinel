package dto

type CreateUserInput struct {
	Fullname        string `validate:"required,min=3,max=100"`
	Email           string `validate:"required,email,max=200"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
	Password        string `validate:"required,min=8,max=255,password"`
}
