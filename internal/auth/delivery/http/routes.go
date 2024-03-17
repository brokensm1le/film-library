package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MapRoutes(rtr *mux.Router, s *AuthHandler) {
	rtr.HandleFunc("/auth/signUp", s.SignUp).Methods(http.MethodPost)
	rtr.HandleFunc("/auth/signIn", s.SignIn).Methods(http.MethodPost)
}
