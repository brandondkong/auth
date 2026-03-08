package main

import (
	"log"

	"github.com/brandondkong/auth/internal/config"
	"github.com/brandondkong/auth/internal/token"
	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/database"
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
}
