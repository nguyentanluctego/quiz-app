package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"quiz-app/internal/constants"
	"quiz-app/internal/domain"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	quizService domain.QuizService
	upgrader    websocket.Upgrader
}

func NewWebSocketHandler(quizService domain.QuizService) *WebSocketHandler {
	return &WebSocketHandler{
		quizService: quizService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Handle WebSocket connection
	go h.HandleConnection(conn)
}

func (h *WebSocketHandler) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	var user *domain.User

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var message domain.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		log.Println("Received message:", message)
		// Handle message based on type
		switch message.Type {

		case constants.MessageTypeJoin:

			var payload domain.JoinPayload
			payloadBytes, _ := json.Marshal(message.Payload)
			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				log.Println("Payload parse error:", err)
				continue
			}
			newUser, quiz, err := h.quizService.JoinQuiz(payload.QuizID, payload.UserName, conn)
			if err != nil {
				log.Println("Error joining quiz:", err)
				h.SendError(conn, "Quiz not found")
				continue
			}

			user = newUser

			h.SendMessage(conn, constants.MessageTypeJoined, map[string]string{
				"quizId": quiz.ID,
				"userId": user.ID,
			})

			h.BroadcastLeaderboard(quiz)

		case constants.MessageTypeAnswer:
			var payload domain.AnswerPayload
			payloadBytes, _ := json.Marshal(message.Payload)
			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				log.Println("Payload parse error:", err)
				continue
			}

			correct, score, quiz, err := h.quizService.SubmitAnswer(payload)
			if err != nil {
				log.Println("Error submitting answer:", err)
				h.SendError(conn, "Failed to process answer")
				continue
			}
			// Response to user
			h.SendMessage(conn, constants.MessageTypeResult, map[string]interface{}{
				"correct": correct,
				"score":   score,
			})

			h.BroadcastLeaderboard(quiz)
		}
	}
}

func (h *WebSocketHandler) SendMessage(conn *websocket.Conn, messageType string, payload interface{}) {
	message := domain.Message{
		Type:    messageType,
		Payload: payload,
	}
	if err := conn.WriteJSON(message); err != nil {
		log.Println("Write error:", err)
	}
}

func (h *WebSocketHandler) SendError(conn *websocket.Conn, errorMsg string) {
	h.SendMessage(conn, constants.MessageTypeError, map[string]string{
		"message": errorMsg,
	})
}

func (h *WebSocketHandler) BroadcastLeaderboard(quiz *domain.Quiz) {
	quiz.Mutex.RLock()
	defer quiz.Mutex.RUnlock()

	// Convert leaderboard entries to DTOs with user details
	leaderboardDTOs := make([]domain.LeaderboardEntryDTO, 0, len(quiz.Leaderboard))
	for _, entry := range quiz.Leaderboard {
		user, ok := quiz.Users[entry.UserId]
		if !ok {
			continue
		}
		
		leaderboardDTOs = append(leaderboardDTOs, domain.LeaderboardEntryDTO{
			UserID:     entry.UserId,
			Score:      entry.Score,
			UserName:   user.Name,
			JoinedAt:   user.JoinedAt,
			LastActive: user.LastActive,
		})
	}

	// Broadcast the leaderboard to all users
	for _, user := range quiz.Users {
		h.SendMessage(user.Conn, constants.MessageTypeLeaderboard, leaderboardDTOs)
	}
}

// RegisterRoutes registers WebSocket
func (h *WebSocketHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/ws", h.HandleWebSocket)
}
