package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DMaryanskiy/gopherstatus/internal/bot"
	"github.com/DMaryanskiy/gopherstatus/internal/monitor"
	"github.com/DMaryanskiy/gopherstatus/internal/server"
	"github.com/DMaryanskiy/gopherstatus/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system vars")
	}

	db, err := storage.NewDB()
	if err != nil {
		log.Fatal("db error:", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	monitor := monitor.NewMonitor(db)
	monitor.Start()

	serv := server.NewServer(monitor, db)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: serv.Handler(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println("starting server at http://localhost:8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutdown signal received...")
		cancel()
	}()

	if botToken != "" {
		telegramBot, err := bot.NewBot(botToken, db)
		if err != nil {
			log.Fatalf("telegram init error: %v", err)
		}
		go telegramBot.Start(ctx)
	}

	<-ctx.Done()

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server Shutdown error: %v", err)
	}

	log.Println("Server exited. Bye!")
}
