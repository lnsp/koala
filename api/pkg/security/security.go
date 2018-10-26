package security

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

// NewJWTGuard returns a new middleware instance that guards the server
// from unauthorized requests.
func NewJWTGuard(secret []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header["Authorization"]
		if auth == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(auth[0], "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unknown signing method")
			}
			return secret, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		logrus.WithFields(logrus.Fields{
			"path": r.URL.Path,
			"claims": token.Claims,
		}).Debug("Authenticated request")
		next.ServeHTTP(w, r)
	})
}
