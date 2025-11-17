package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// Split into scheme + credentials
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// Decode the Base64 portion
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid base64 credentials: %w", err))
				return
			}

			// Credentials come as "username:password"
			cred := strings.SplitN(string(decoded), ":", 2)
			if len(cred) != 2 {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("credentials format is invalid"))
				return
			}

			username, password := app.config.auth.basic.username, app.config.auth.basic.password
			if cred[0] != username || cred[1] != password {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("not authorized"))
				return
			}

			// Success → continue the handler chain
			next.ServeHTTP(w, r)
		})

	}
}
