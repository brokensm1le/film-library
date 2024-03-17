package usecase

import (
	"crypto/sha256"
	"film_library/internal/auth"
	"film_library/internal/cconstant"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthUsecase struct {
	repo auth.Repository
}

func NewAuthUsecase(repo auth.Repository) auth.Usecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) CreateUser(user *auth.User) error {
	user.Password = u.generatePasswordHash(user.Password)
	return u.repo.CreateUser(user)
}

func (u *AuthUsecase) GenerateToken(params *auth.SignInParams) (string, error) {
	params.Password = u.generatePasswordHash(params.Password)
	user, err := u.repo.GetUser(params)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cconstant.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:   user.Id,
		Role: user.Role,
	})
	return token.SignedString([]byte(cconstant.SignedKey))
}

func (u *AuthUsecase) ParseToken(accessToken string) (*auth.TokenData, error) {
	token, err := jwt.ParseWithClaims(accessToken, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid singing method")
		}

		return []byte(cconstant.SignedKey), nil
	})
	if err != nil {
		return &auth.TokenData{}, err
	}

	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok {
		return &auth.TokenData{}, fmt.Errorf("invalid claims type")
	}

	return &auth.TokenData{Id: claims.Id, Role: claims.Role}, nil
}

func (u *AuthUsecase) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(cconstant.Salt)))
}
