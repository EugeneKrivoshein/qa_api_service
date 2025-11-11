package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/EugeneKrivoshein/qa_api_service/config"
	"github.com/EugeneKrivoshein/qa_api_service/internal/answer"
	"github.com/EugeneKrivoshein/qa_api_service/internal/question"
	"github.com/EugeneKrivoshein/qa_api_service/internal/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	qRepo := repo.NewQuestionRepo(db)
	aRepo := repo.NewAnswerRepo(db)

	qService := question.NewService(qRepo, aRepo)
	aService := answer.NewService(aRepo, qRepo)

	qHandler := &question.Handler{Service: qService}
	aHandler := &answer.Handler{Service: aService}

	r := NewRouter(qHandler, aHandler)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Server running on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server gracefully stopped")
}
