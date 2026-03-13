package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/brandondkong/auth/internal/config"
	"github.com/brandondkong/auth/pkg/jwtutil"
	"github.com/google/uuid"
)

func GetUserId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	userId, ok := r.Context().Value(UserIdKey).(string)
	if !ok {
		WriteJsonResponse(w, ResponseOptions[any]{
			Code: http.StatusUnauthorized,
			Error: &UNAUTHORIZED_ERROR_CODE,
			Message: "Failed to retrieve user ID",
		})

		return uuid.Nil, errors.New("could not find user ID")
	}

	parsed, err := uuid.Parse(userId)
	if err != nil {
		WriteJsonResponse(w, ResponseOptions[any]{
			Code: http.StatusUnauthorized,
			Error: &UNAUTHORIZED_ERROR_CODE,
			Message: "Failed to parse user ID",
		})
		return uuid.Nil, err
	}

	return parsed, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Extract the authorization header
		token := req.Header.Get("Authorization")
		if token == "" {
			// 401
			WriteJsonResponse(w, ResponseOptions[any]{
				Code:	http.StatusUnauthorized,
				Error: &UNAUTHORIZED_ERROR_CODE,
				Message: "Failed to retrieve authorization key in header",	
			})
			
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		// Get the access token key
		configs, err := config.LoadConfigs()
		if err != nil {
			// 401
			WriteJsonResponse(w, ResponseOptions[any]{
				Code:	http.StatusUnauthorized,
				Error: &UNAUTHORIZED_ERROR_CODE,
				Message: "Failed to load signing key",	
			})

			return
		}
		jwtToken, err := jwtutil.ParseToken(token, configs.JwtAccessSigningKey)
		if err != nil {
			// 401
			WriteJsonResponse(w, ResponseOptions[any]{
				Code:	http.StatusUnauthorized,
				Error: &UNAUTHORIZED_ERROR_CODE,
				Message: "Failed to parse token",	
			})

			return
		}

		userId, err := jwtToken.Claims.GetSubject()
		if err != nil {
			// 401
			WriteJsonResponse(w, ResponseOptions[any]{
				Code:	http.StatusUnauthorized,
				Error: &UNAUTHORIZED_ERROR_CODE,
				Message: "Failed to retrieve user ID from JWT claim",	
			})
			
			return
		}	

		req = req.WithContext(context.WithValue(req.Context(), UserIdKey, userId))
		next.ServeHTTP(w, req)
	})
}

