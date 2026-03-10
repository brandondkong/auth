package auth

import (
	"errors"
	"fmt"
	"net/http"

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
			http.Error(w, mr.Error(), mr.Status())
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	token, err := token.GenerateMagicLink(payload.Email, r)
	if err != nil {
		res := fmt.Sprintf("error verifying token: %v\n", err)
		http.Error(w, res, http.StatusBadRequest)
		return
	}

	middleware.WriteJsonResponse(w, http.StatusAccepted, "success", CreateMagicLinkResponse{
		Token: token,
	})
}


func consumeMagicLink(w http.ResponseWriter, r *http.Request) {
	urlToken := chi.URLParam(r, "token")
	user, err := token.ConsumeMagicLink(urlToken)
	buf := []byte{}

	if err != nil {
		w.WriteHeader(400)
		res := fmt.Appendf(buf, "error verifying token: %v\n", err)
		w.Write(res)
		return
	}
	w.WriteHeader(200)
	res := fmt.Appendf(buf, "user id: %s\n", user.ID.String())
	w.Write(res)
}
