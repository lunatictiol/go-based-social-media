package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (a *application) basicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				a.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("autorization header missing"))
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Basic" {
				a.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("autorization header is malformed"))
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				a.unauthorizedBasicErrorResponse(w, r, err)
				return
			}
			adminUser := a.config.auth.basic.admin
			adminPassword := a.config.auth.basic.adminPassword
			creds := strings.SplitN(string(decoded), ":", 2)

			if len(creds) != 2 || creds[0] != adminUser || creds[1] != adminPassword {
				a.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials "))
				return
			}
			next.ServeHTTP(w, r)

		})
	}
}
