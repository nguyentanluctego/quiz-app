package domain

import "github.com/gorilla/websocket"

// Repository interfaces
type QuizRepository interface {
	GetQuiz(id string) (*Quiz, error)
	SaveQuiz(quiz *Quiz)
	AddUserToQuiz(quizID string, user *User) (*Quiz, error)
}

// Service interfaces
type QuizService interface {
	GetQuiz(id string) (*Quiz, error)
	GetQuizForClient(id string) (*QuizResponse, error)
	JoinQuiz(quizID string, userName string, conn *websocket.Conn) (*User, *Quiz, error)
	SubmitAnswer(payload AnswerPayload) (bool, int, *Quiz, error)
	UpdateLeaderboard(quiz *Quiz)
}

// Handler interfaces
type WebSocketHandler interface {
	HandleConnection(conn *websocket.Conn)
	SendMessage(conn *websocket.Conn, messageType string, payload interface{})
	SendError(conn *websocket.Conn, errorMsg string)
	BroadcastLeaderboard(quiz *Quiz)
}
