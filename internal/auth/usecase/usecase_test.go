package usecase

import (
	"film_library/internal/auth"
	mock_auth "film_library/internal/auth/mocks"
	"film_library/internal/cconstant"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseToken(t *testing.T) {

	u := AuthUsecase{}

	var (
		id   = 999
		role = 1
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cconstant.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:   id,
		Role: role,
	})
	accessToken, err := token.SignedString([]byte(cconstant.SignedKey))
	require.NoError(t, err)

	encodeData, err := u.ParseToken(accessToken)
	require.NoError(t, err)
	require.Equal(t, id, encodeData.Id)
	require.Equal(t, role, encodeData.Role)
}

func TestCreateUser(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	repo := mock_auth.NewMockRepository(ctr)
	in := auth.User{Login: "123", Password: "123"}

	repo.EXPECT().CreateUser(&in).Return(nil).Times(1)
	useCase := NewAuthUsecase(repo)
	err := useCase.CreateUser(&in)
	require.NoError(t, err)
}

func TestGenerateToken(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	repo := mock_auth.NewMockRepository(ctr)
	u := AuthUsecase{}
	in := auth.SignInParams{Login: "123", Password: "123"}
	out := auth.User{Id: 1, Login: "123", Password: u.generatePasswordHash("123"), Role: 1}

	repo.EXPECT().GetUser(&in).Return(&out, nil).Times(1)
	useCase := NewAuthUsecase(repo)
	accessToken, err := useCase.GenerateToken(&in)
	require.NoError(t, err)

	data, err := useCase.ParseToken(accessToken)
	require.NoError(t, err)
	require.Equal(t, data.Id, 1)
	require.Equal(t, data.Role, 1)
}
