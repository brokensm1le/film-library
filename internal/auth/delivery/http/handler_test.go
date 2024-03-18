package http

import (
	"bytes"
	"encoding/json"
	"film_library/internal/auth"
	mock_auth "film_library/internal/auth/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {
	type mockBehavior func(s *mock_auth.MockUsecase, user auth.User)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)
	resp = &auth.ResponseModel{Status: "error", Error: "error"}
	Err, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           auth.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Login":"abc", "Password":"123"}`,
			inputUser: auth.User{
				Login:    "abc",
				Password: "123",
			},
			mockBehavior: func(s *mock_auth.MockUsecase, user auth.User) {
				s.EXPECT().CreateUser(&user).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: `{"Login":"abc", "Password":"123"}`,
			inputUser: auth.User{
				Login:    "abc",
				Password: "123",
			},
			mockBehavior: func(s *mock_auth.MockUsecase, user auth.User) {
				s.EXPECT().CreateUser(&user).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: Err,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockAuth, testCase.inputUser)

			handler := NewAuthHandler(mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/auth/signUp", bytes.NewBufferString(testCase.inputBody))
			handler.SignUp(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestSignIn(t *testing.T) {
	type mockBehavior func(s *mock_auth.MockUsecase, user auth.SignInParams)
	var token = "12344321"
	var resp *auth.SignInResponse = &auth.SignInResponse{Token: token}
	ans, _ := json.Marshal(resp)
	//Err, _ := json.Marshal(fmt.Errorf("error").Error())

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           auth.SignInParams
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Login":"abc", "Password":"123"}`,
			inputUser: auth.SignInParams{
				Login:    "abc",
				Password: "123",
			},
			mockBehavior: func(s *mock_auth.MockUsecase, user auth.SignInParams) {
				s.EXPECT().GenerateToken(&user).Return(token, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "OK",
			inputBody: `{"Login":"abc", "Password":"123"}`,
			inputUser: auth.SignInParams{
				Login:    "abc",
				Password: "123",
			},
			mockBehavior: func(s *mock_auth.MockUsecase, user auth.SignInParams) {
				s.EXPECT().GenerateToken(&user).Return("", fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockAuth, testCase.inputUser)

			handler := NewAuthHandler(mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/auth/signIn", bytes.NewBufferString(testCase.inputBody))
			handler.SignIn(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}
