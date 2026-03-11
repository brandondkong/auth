package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brandondkong/auth/internal/auth"
	"github.com/brandondkong/auth/internal/config"
	"github.com/brandondkong/auth/internal/jwt"
	"github.com/brandondkong/auth/internal/token"
	"github.com/brandondkong/auth/internal/user"
	"github.com/brandondkong/auth/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-co-op/gocron/v2"
)

const PORT int = 5000

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
	err = db.AutoMigrate(&user.User{}, &token.MagicLinkToken{}, &jwt.RefreshToken{}, &user.OAuthAccount{})
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*", "0.0.0.0", "127.0.0.1"},
	})

	r.Use(cors.Handler)

	auth.Routes(r)
	
	// Run cron jobs here
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(3, 0, 0),
			),
		),
		gocron.NewTask(func() {
			token.CleanupStaleTokens()
			jwt.CleanupStaleRefreshTokens()
		}),
		)

	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	log.Println("Starting cron worker")
	scheduler.Start()

	log.Printf("Starting server on port %d\n", PORT)
	go http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)

	<- sigChan
	log.Println("Shutting down cron jobs")
	err = scheduler.Shutdown()
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	cancel()
	log.Println("Server shutdown complete")
}
