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

	log.Println("Starting database")
	_, err := database.StartDatabase()
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	log.Println("Migrating database tables")
	err = database.Migrate(&user.User{}, &user.OAuthAccount{}, &token.MagicLinkToken{}, &jwt.RefreshToken{})
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"127.0.0.1"},
		AllowCredentials: true,
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
			err := token.CleanupStaleTokens()
			if err != nil {
				log.Printf("error cleaning tokens: %v\n", err)
			}
			err = jwt.CleanupStaleRefreshTokens()
			if err != nil {
				log.Printf("error cleaning refresh tokens: %v\n", err)
			}

		}),
		)

	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	log.Println("Starting cron worker")
	scheduler.Start()

	log.Printf("Starting server on port %d\n", PORT)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("server error: %v\n", err)
	case <-sigChan:
	}

	log.Println("Shutting down cron jobs")
	err = scheduler.Shutdown()
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

	cancel()
	log.Println("Server shutdown complete")
}
