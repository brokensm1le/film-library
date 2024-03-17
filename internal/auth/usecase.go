package auth

type Usecase interface {
	CreateUser(user *User) error
	GenerateToken(params *SignInParams) (string, error)
	ParseToken(token string) (*TokenData, error)
}
