package main

import (
	"log"
	"net/http"

	"github.com/brandondkong/auth/internal/auth"
	"github.com/brandondkong/auth/internal/config"
	"github.com/brandondkong/auth/internal/token"
	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("Reading environment variables")
	configs, err := config.LoadConfigs()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Println("Starting database")
	db, err := database.StartDatabase(configs)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Println("Migrating database tables")
	err = db.AutoMigrate(&user.User{}, token.MagicLinkToken{})

	log.Println("Starting auth server...")
	
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	auth.Routes(r)

	http.ListenAndServe(":5000", r)
}
