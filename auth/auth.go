package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}
type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type User struct {
	ID    string
	Username string
	Email string
}

// Middleware decodes the session cookie and packs the userID into context.
func Middleware(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("tp-access-token")

			// Allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Decode the JWT
			claims := &jwt.StandardClaims{}
            accessToken, err := jwt.ParseWithClaims(c.Value, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(os.Getenv("TP_ACCESS_TOKEN_SECRET")), nil
            })

            if err != nil || !accessToken.Valid {
                http.Error(w, "Invalid token", http.StatusForbidden)
                return
            }

            userID := claims.Subject

			// Check if the userID exists in the database
			var user User
			err = db.QueryRow(r.Context(), "SELECT id, username, email FROM users WHERE userid = $1", userID).Scan(&user.ID, &user.Username, &user.Email)
            if err != nil {
                http.Error(w, "Invalid userID", http.StatusForbidden)
                return
            }

			// Put the userID in context
			ctx := context.WithValue(r.Context(), userCtxKey, userID)

			// Call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}