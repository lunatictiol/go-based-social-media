package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lunatictiol/go-based-social-media/internal/store"
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
func (a *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.unauthorisedResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.unauthorisedResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := a.authenticator.ValidateToken(token)
		if err != nil {
			a.unauthorisedResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			a.unauthorisedResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := a.getUser(ctx, userID)
		if err != nil {
			a.unauthorisedResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if !a.config.redisConfig.enabled {
		return a.store.Users.GetUserByID(ctx, userID)
	}

	user, err := a.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = a.store.Users.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if err := a.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (a *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.config.rateLimiter.Enabled {
			if allow, retryAfter := a.ratelimiter.Allow(r.RemoteAddr); !allow {
				a.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserfromCtx(r)
		post := a.getPostfromCtx(r)

		if post.UserId == user.Id {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := a.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			a.WriteInternalServerError(w, r, err)
			return
		}

		if !allowed {
			a.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
