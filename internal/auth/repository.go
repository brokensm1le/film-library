package auth

type Repository interface {
	CreateUser(user *User) error
	GetUser(params *SignInParams) (*User, error)
}
