package user

import (
	"net/http"

	"github.com/brandondkong/auth/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/user", getUser)
	return r
}

func getUser(w http.ResponseWriter, req *http.Request) {
	userId, err := middleware.GetUserId(req)
	if err != nil {
		return
	}

	user, err := GetUserById(userId, nil)
	if err != nil {
		return
	}

	middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
		Code: http.StatusOK,
		Data: user,
	})
}
