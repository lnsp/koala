package security

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type guardResponse struct {
	Security string `json:"security"`

	ClientID string `json:"clientId"`
	Provider string `json:"provider"`
}

type Guard func(next http.Handler) http.Handler

// None defines a no-security guard.
func None() Guard {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				resp, _ := json.Marshal(guardResponse{Security: "none"})
				w.Write(resp)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// JWT defines a JWT-based security guard that verifies the Authorization header token.
func JWT(secret string) Guard {
	return func(next http.Handler) http.Handler {
		secret := []byte(secret)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				resp, _ := json.Marshal(guardResponse{Security: "jwt"})
				w.Write(resp)
				return
			}
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
				"path":   r.URL.Path,
				"claims": token.Claims,
			}).Debug("Authenticated request")
			next.ServeHTTP(w, r)
		})
	}
}

// OIDC defines an OpenID-Connect security guard by verifying the id_token in the Authorization header.
func OIDC(clientID, identityProvider string) Guard {
	provider, err := oidc.NewProvider(context.Background(), identityProvider)
	if err != nil {
		logrus.WithError(err).Fatal("set up OIDC provider")
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				resp, _ := json.Marshal(guardResponse{
					Security: "oidc",
					ClientID: clientID,
					Provider: identityProvider,
				})
				w.Write(resp)
				return
			}
			// TODO: Implement token verification
			auth := r.Header["Authorization"]
			if auth == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(auth[0], "Bearer ")
			idToken, err := verifier.Verify(context.Background(), tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			logrus.WithFields(logrus.Fields{
				"path":    r.URL.Path,
				"subject": idToken.Subject,
			}).Debug("Authenticated request")
			next.ServeHTTP(w, r)
		})
	}
}
