package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"quiz-app/internal/domain"
)

type HTTPHandler struct {
	quizService domain.QuizService
	router      *gin.Engine
}

func NewHTTPHandler(quizService domain.QuizService) *HTTPHandler {
	router := gin.Default()

	handler := &HTTPHandler{
		quizService: quizService,
		router:      router,
	}

	// Register routes
	handler.registerRoutes()

	return handler
}

func (h *HTTPHandler) GetQuiz(c *gin.Context) {
	id := c.Param("id")

	quizResponse, err := h.quizService.GetQuizForClient(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quizResponse)
}

// registerRoutes registers HTTP routes
func (h *HTTPHandler) registerRoutes() {
	api := h.router.Group("/api")
	{
		api.GET("/quiz/:id", h.GetQuiz)
	}
}

func (h *HTTPHandler) GetRouter() *gin.Engine {
	return h.router
}
