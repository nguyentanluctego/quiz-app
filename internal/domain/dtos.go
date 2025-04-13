package domain

import (
	"time"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type AnswerPayload struct {
	QuizID     string `json:"quizId"`
	UserID     string `json:"userId"`
	QuestionID int    `json:"questionId"`
	Answer     int    `json:"answer"`
}

type JoinPayload struct {
	QuizID   string `json:"quizId"`
	UserName string `json:"userName"`
}

type QuizResponse struct {
	ID        string     `json:"id"`
	Questions []Question `json:"questions"`
}

type ResultPayload struct {
	Correct bool `json:"correct"`
	Score   int  `json:"score"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

// JoinedPayload is a DTO for successful join response
type JoinedPayload struct {
	QuizID string `json:"quizId"`
	UserID string `json:"userId"`
}

type LeaderboardEntryDTO struct {
	UserID     string    `json:"userId"`
	Score      int       `json:"score"`
	UserName   string    `json:"userName"`
	JoinedAt   time.Time `json:"joinedAt"`
	LastActive time.Time `json:"lastActive"`
}
