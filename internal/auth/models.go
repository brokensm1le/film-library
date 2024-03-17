package auth

import "github.com/dgrijalva/jwt-go"

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	Role     int    `json:"-" db:"role"`
}

type SignInParams struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type TokenData struct {
	Id   int `json:"id"`
	Role int `json:"role"`
}

type CustomClaims struct {
	jwt.StandardClaims
	Id   int `json:"id"`
	Role int `json:"role"`
}

type ResponseModel struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
