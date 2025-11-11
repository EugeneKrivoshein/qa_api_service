package main

import (
	"github.com/EugeneKrivoshein/qa_api_service/internal/api"
)

// @title QA API Service
// @version 1.0
// @description API для вопросов и ответов
// @host localhost:8080
// @BasePath /
func main() {
	api.StartServer()
}
