package auth

import (
	"fmt"
	"net/http"

	"github.com/brandondkong/auth/internal/token"
	"github.com/go-chi/chi/v5"
)

func Routes(router chi.Router) {
	r := chi.NewRouter()
	r.Get("/magic-link", createMagicLink)
	r.Get("/magic-link/callback", consumeMagicLink)
	router.Mount("/api/auth", r)
}

func createMagicLink(w http.ResponseWriter, r *http.Request) {
	token, err := token.GenerateMagicLink("email", r)
	buf := []byte{}
	if err != nil {
		w.WriteHeader(400)
		res := fmt.Appendf(buf, "error verifying token: %v\n", err)
		w.Write(res)
		return
	}

	w.WriteHeader(200)
	res := fmt.Appendf(buf, "user id: %s\n", token)
	w.Write(res)
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
