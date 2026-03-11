package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/brandondkong/auth/internal/jwt"
	"github.com/brandondkong/auth/internal/middleware"
	"github.com/brandondkong/auth/internal/token"
	"github.com/go-chi/chi/v5"
)

type CreateMagicLinkPayload struct {
	Email	string	`json:"email"`
}

type CreateMagicLinkResponse struct {
	Token	string `json:"token"`
}

func Routes(router chi.Router) {
	r := chi.NewRouter()
	r.Post("/magic-link", createMagicLink)
	r.Get("/magic-link/callback", consumeMagicLink)
	router.Mount("/api/auth", r)
}

func createMagicLink(w http.ResponseWriter, r *http.Request) {
	var payload CreateMagicLinkPayload

	err := middleware.DecodeJsonRequestBody(w, r, &payload)
	if err != nil {
		var mr *middleware.ErrMalformedRequest
		if errors.As(err, &mr) {
			middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
				Code:	mr.Status(),
				Error: &middleware.BAD_REQUEST_ERROR_CODE,
				Message: mr.Error(),	
			})
		} else {
			middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
				Code:	http.StatusInternalServerError,
				Error: &middleware.INTERNAL_SERVER_ERROR_CODE,
				Message: "Something went wrong",	
			})
		}
		return
	}

	// Verify email
	addr, err := mail.ParseAddress(payload.Email)
	if err != nil {
		middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
				Code:	http.StatusBadRequest,
				Error: &middleware.INVALID_EMAIL_ERROR_CODE,
				Message: "Invalid email address",
			})
		return
	}

	token, err := token.GenerateMagicLink(addr.Address, r)
	if err != nil {
		res := fmt.Sprintf("Error verifying token: %v\n", err)
		middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
				Code:	http.StatusBadRequest,
				Error: &middleware.BAD_REQUEST_ERROR_CODE,
				Message: res,	
			})

		return
	}

	middleware.WriteJsonResponse(w, middleware.ResponseOptions[CreateMagicLinkResponse]{
		Code: http.StatusOK,
		Data: CreateMagicLinkResponse{
			Token: token,
		},
	})
}


func consumeMagicLink(w http.ResponseWriter, r *http.Request) {
	urlToken := r.URL.Query().Get("token")
	user, err := token.ConsumeMagicLink(urlToken)

	if err != nil {
		res := fmt.Sprintf("Error verifying token: %v", err)
		middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
				Code:	http.StatusBadRequest,
				Error: &middleware.BAD_REQUEST_ERROR_CODE,
				Message: res,	
			})
		return
	}
	
	// Generate a TokenPair for the user
	tokenPair, err := jwt.CreateTokenPair(user)	
	if err != nil {
		res := fmt.Sprintf("Error creating token pair: %v", err)
		middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
				Code:	http.StatusBadRequest,
				Error: &middleware.BAD_REQUEST_ERROR_CODE,
				Message: res,	
			})
		return
	}

	// Set the cookie in the header
	http.SetCookie(w, &http.Cookie{
		Name: "refresh_token",
		Value: tokenPair.Refresh,
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteLaxMode,
		Expires: time.Now().Add(jwt.RefreshTokenLifeTime),
		Path: "/",
	})

	middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
			Code:	http.StatusOK,
			Data: tokenPair,
		})
}
