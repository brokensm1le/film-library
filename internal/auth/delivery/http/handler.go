package http

import (
	"encoding/json"
	"film_library/internal/auth"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type AuthHandler struct {
	authUC auth.Usecase
}

func NewAuthHandler(authUC auth.Usecase) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

// @Summary      SignUp
// @Description  Create account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input	body	auth.User  true  "user data"
// @Success      200  {object}	auth.ResponseModel
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /auth/signUp [post]
func (h *AuthHandler) SignUp(rw http.ResponseWriter, r *http.Request) {
	var (
		data auth.User
		resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	)

	log.Printf("Request: SignUP")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: SignUp. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	pattern, _ := regexp.Compile("[A-Za-z0-9@.]+")
	if !pattern.MatchString(data.Login) {
		log.Printf("Request: SignUp. Error: %s", "Uncorrect data")
		http.Error(rw, fmt.Sprintf("login must contain the characters a-z, A-z, 0-9, @ and ."), http.StatusBadRequest)
		return
	}

	err := h.authUC.CreateUser(&data)
	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
	rawResponse, _ := json.Marshal(resp)
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}

// @Summary      SignIn
// @Description  Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input	body	auth.SignInParams  true  "login and password"
// @Success      200  {object}	auth.SignInResponse
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /auth/signIn [post]
func (h *AuthHandler) SignIn(rw http.ResponseWriter, r *http.Request) {
	var (
		data  auth.SignInParams
		token string
		err   error
	)

	log.Printf("Request: SignIn")

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Request: SignIn. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	pattern, _ := regexp.Compile("[A-Za-z0-9@.]+")
	if !pattern.MatchString(data.Login) {
		log.Printf("Request: SignIn. Error: %s", "Uncorrect data")
		http.Error(rw, fmt.Sprintf("login must contain the characters a-z, A-z, 0-9, @ and ."), http.StatusBadRequest)
		return
	}

	token, err = h.authUC.GenerateToken(&data)
	if err != nil {
		log.Printf("Request: SignIn. Error: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rawResponse, _ := json.Marshal(auth.SignInResponse{Token: token})
	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(rawResponse)
}
