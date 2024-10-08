package dto

type CreateUserInput struct {
	Provider      string
	Fullname      string
	Email         string
	OAuthProvider string
	OAuthID       string
	ProfileURL    string
	PublicKey     string
}
