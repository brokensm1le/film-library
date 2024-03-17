package http

import (
	"context"
	"film_library/internal/cconstant"
	"fmt"
	"net/http"
	"strings"
)

func (s *ServiceHandler) userIdentity(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(cconstant.AuthHeader)
		if header == "" {
			http.Error(rw, fmt.Sprintf("empty auth header"), http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(rw, fmt.Sprintf("invalid auth header"), http.StatusUnauthorized)
			return
		}

		tokenData, err := s.authUC.ParseToken(headerParts[1])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "tokenData", tokenData)

		h.ServeHTTP(rw, r.WithContext(ctx))
	})
}
