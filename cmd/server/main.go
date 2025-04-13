package main

import (
	"fmt"

	"quiz-app/internal/handler"
	"quiz-app/internal/repository"
	"quiz-app/internal/service"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	quizRepo := repository.NewInMemoryQuizRepository()

	// Seed data for testing
	quizRepo.InitSampleData()

	quizService := service.NewQuizService(quizRepo)

	httpHandler := handler.NewHTTPHandler(quizService)
	router := httpHandler.GetRouter()

	wsHandler := handler.NewWebSocketHandler(quizService)
	wsHandler.RegisterRoutes(router)

	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// Start the server
	fmt.Println("Server started on :8080")
	router.Run(":8080")
}
