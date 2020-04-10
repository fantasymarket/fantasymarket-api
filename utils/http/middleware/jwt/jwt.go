package jwt

import (
	"context"
	"fantasymarket/utils/http/responses"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
)

type key string

// UserKey is the key of the user context value
const UserKey key = "user"

// Middleware implements a middleware handler for adding jwt auth to a route.
func Middleware(secret string, optional bool) func(next http.Handler) http.Handler {
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			jwtauth.Verifier(tokenAuth)(next)
			authenticator(next, optional)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func authenticator(next http.Handler, optional bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if !optional {
			if err != nil || token == nil || !token.Valid {
				responses.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if tokenType, ok := claims["type"]; !ok || tokenType != "auth" {
				responses.ErrorResponse(w, http.StatusUnauthorized, "you can't use this type of token on this method")
				return
			}
		}

		// Set userID context to the user_id
		userID, userOk := claims["user_id"]

		if userOk {
			ctx := context.WithValue(r.Context(), UserKey, Claims{
				UserID: userID.(uuid.UUID),
			})

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Error if something is wrong with the token
		if !optional {
			responses.ErrorResponse(w, http.StatusUnauthorized, "malformed jwt token")
		}
	})
}

// Claims is our custom claims type
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

// CreateToken issues a new jwt token
func CreateToken(secret string, username string, userID uuid.UUID) (string, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	// Create the Claims
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "fantasymarket-api",
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}
