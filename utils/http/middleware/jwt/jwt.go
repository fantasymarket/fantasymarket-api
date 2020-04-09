package jwt

import (
	"context"
	"fantasymarket/utils/http/responses"
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

		if !optional && (err != nil || token == nil || !token.Valid) {
			responses.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Set userID context to the user_id
		user, userOk := claims["user"].(UserClaims)

		if userOk && user.UserID != "" && user.Username != "" {
			ctx := context.WithValue(r.Context(), UserKey, user)

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
	jwt.StandardClaims
	User UserClaims `json:"user"`
}

// UserClaims are user informations
type UserClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// CreateToken issues a new jwt token
func CreateToken(secret string, username string, userID string) (string, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	// Create the Claims
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "fantasymarket-api",
		},
		User: UserClaims{
			Username: username,
			UserID:   userID,
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}
