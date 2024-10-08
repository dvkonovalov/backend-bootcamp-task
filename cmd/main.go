package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"main/internal/config"
	"main/internal/http_server/mailsender"
	"main/internal/http_server/urls/auth"
	"main/internal/http_server/urls/flat"
	"main/internal/http_server/urls/house"
	"main/internal/storage/db"
	"net/http"
	"os"
)

const (
	envLocal      = "local"
	envProduction = "production"
)

func main() {
	// Load config
	cnf := config.MustLoad()

	log := SetUpLogger(cnf.Env)
	log.Info("Logger start!", slog.String("env", cnf.Env))

	// Load Database
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	storage, err := db.NewStorage(dsn)
	if err != nil {
		log.Error("Failed to create storage", err)
		os.Exit(1)
	}

	mailSender := mailsender.New()

	// Config Router
	router := mux.NewRouter()
	router.HandleFunc("/house/create", house.Create(log, storage))
	router.HandleFunc("/flat/create", flat.Create(log, storage))
	router.HandleFunc("/flat/update", flat.Update(log, mailSender, storage, storage))
	router.HandleFunc("/house/{id}", house.GetFlats(log, storage))
	router.HandleFunc("/dummyLogin", auth.CreateToken(log))
	router.HandleFunc("/register", auth.CreateUser(log, storage))
	router.HandleFunc("/login", auth.LoginUser(log, storage))
	router.HandleFunc("/house/{id}/subscribe", house.Subscribe(log, storage))

	// Starting server
	log.Info("Starting server", slog.String("address", cnf.HttpServer.Address))
	srv := &http.Server{
		Addr:         cnf.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cnf.HttpServer.Timeout,
		WriteTimeout: cnf.HttpServer.Timeout,
		IdleTimeout:  cnf.HttpServer.IdleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("Failed to start server", err)
		os.Exit(1)
	}

}

func SetUpLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	}
	return logger
}
