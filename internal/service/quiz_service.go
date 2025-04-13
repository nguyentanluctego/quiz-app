package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/gorilla/websocket"

	"quiz-app/internal/constants"
	"quiz-app/internal/domain"
)

// QuizService implements the domain.QuizService interface
type QuizService struct {
	repo domain.QuizRepository
}

func NewQuizService(repo domain.QuizRepository) *QuizService {
	return &QuizService{
		repo: repo,
	}
}

// GetQuiz retrieves a quiz by ID
func (s *QuizService) GetQuiz(id string) (*domain.Quiz, error) {
	quiz, err := s.repo.GetQuiz(id)
	if err != nil {
		return nil, fmt.Errorf("quiz with ID %s not found", id)
	}
	return quiz, nil
}

func (s *QuizService) GetQuizForClient(id string) (*domain.QuizResponse, error) {
	quiz, err := s.repo.GetQuiz(id)
	if err != nil {
		return nil, fmt.Errorf("quiz with ID %s not found", id)
	}

	quizCopy := &domain.QuizResponse{
		ID:        quiz.ID,
		Questions: make([]domain.Question, len(quiz.Questions)),
	}

	for i, q := range quiz.Questions {
		quizCopy.Questions[i] = q
		quizCopy.Questions[i].Answer = -1 // Hide the answer
	}

	return quizCopy, nil
}

func (s *QuizService) JoinQuiz(quizID string, userName string, conn *websocket.Conn) (*domain.User, *domain.Quiz, error) {
	quiz, err := s.repo.GetQuiz(quizID)
	if err != nil {
		return nil, nil, err
	}

	existingUser, exists := quiz.Users[userName]
	if exists {
		quiz.Mutex.Lock()
		defer quiz.Mutex.Unlock()
		existingUser.Conn = conn
		existingUser.LastActive = time.Now()
		return existingUser, quiz, nil
	}

	user := &domain.User{
		ID:         userName,
		JoinedAt:   time.Now(),
		Name:       userName,
		Score:      0,
		Conn:       conn,
		LastActive: time.Now(),
	}

	_, err = s.repo.AddUserToQuiz(quizID, user)
	if err != nil {
		return nil, nil, err
	}

	s.UpdateLeaderboard(quiz)

	return user, quiz, nil
}

func (s *QuizService) SubmitAnswer(payload domain.AnswerPayload) (bool, int, *domain.Quiz, error) {
	quiz, err := s.repo.GetQuiz(payload.QuizID)
	if err != nil {
		return false, 0, nil, fmt.Errorf("quiz with ID %s not found", payload.QuizID)
	}

	quiz.Mutex.Lock()
	defer quiz.Mutex.Unlock()

	user, exists := quiz.Users[payload.UserID]
	if !exists {
		return false, 0, nil, fmt.Errorf("user with ID %s not found in quiz", payload.UserID)
	}

	var question *domain.Question
	for i := range quiz.Questions {
		if quiz.Questions[i].ID == payload.QuestionID {
			question = &quiz.Questions[i]
			break
		}
	}

	if question == nil {
		return false, 0, nil, fmt.Errorf("question with ID %d not found in quiz", payload.QuestionID)
	}

	correct := question.Answer == payload.Answer
	if correct {
		user.Score += constants.PointsPerCorrectAnswer
	}

	s.UpdateLeaderboard(quiz)

	return correct, user.Score, quiz, nil
}

func (s *QuizService) UpdateLeaderboard(quiz *domain.Quiz) {

	leaderboard := make([]domain.LeaderboardEntry, 0, len(quiz.Users))

	for _, user := range quiz.Users {
		leaderboard = append(leaderboard, domain.LeaderboardEntry{
			UserId: user.ID,
			Score:  user.Score,
		})
	}
	// Sort leaderboard by score
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	quiz.Leaderboard = leaderboard
}
