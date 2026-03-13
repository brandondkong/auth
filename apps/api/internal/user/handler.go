package user

import (
	"net/http"

	"github.com/brandondkong/auth/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)
	r.Get("/", getUser)
	return r
}

func getUser(w http.ResponseWriter, req *http.Request) {
	userId, err := middleware.GetUserId(w, req)
	if err != nil {
		return
	}

	user, err := GetUserById(userId, nil)
	if err != nil {
		middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
			Code: http.StatusBadRequest,
			Error: &USER_NOT_FOUND_ERROR_CODE,
			Message: "User does not exist",
		})
		return
	}

	middleware.WriteJsonResponse(w, middleware.ResponseOptions[any]{
		Code: http.StatusOK,
		Data: user,
	})
}
