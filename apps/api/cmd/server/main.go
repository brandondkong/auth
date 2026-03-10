package main

import (
	"fmt"
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

const PORT int = 5000

func main() {
	log.Println("Reading environment variables")
	configs, err := config.LoadConfigs()
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	log.Println("Starting database")
	db, err := database.StartDatabase(configs)
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	log.Println("Migrating database tables")
	err = db.AutoMigrate(&user.User{}, token.MagicLinkToken{})
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	auth.Routes(r)
	
	log.Printf("Starting server on port %d\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
}
