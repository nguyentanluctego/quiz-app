package repository

import (
	"fmt"
	"sync"

	"quiz-app/internal/constants"
	"quiz-app/internal/domain"
)

type InMemoryQuizRepository struct {
	quizzes   map[string]*domain.Quiz
	quizMutex sync.RWMutex
}

func NewInMemoryQuizRepository() *InMemoryQuizRepository {
	return &InMemoryQuizRepository{
		quizzes: make(map[string]*domain.Quiz),
	}
}

// GetQuiz retrieves a quiz by ID
func (r *InMemoryQuizRepository) GetQuiz(id string) (*domain.Quiz, error) {
	r.quizMutex.RLock()
	defer r.quizMutex.RUnlock()

	quiz, exists := r.quizzes[id]
	if !exists {
		return nil, fmt.Errorf("quiz with ID %s not found", id)
	}
	return quiz, nil
}

func (r *InMemoryQuizRepository) SaveQuiz(quiz *domain.Quiz) {
	r.quizMutex.Lock()
	defer r.quizMutex.Unlock()

	r.quizzes[quiz.ID] = quiz
}

func (r *InMemoryQuizRepository) AddUserToQuiz(quizID string, user *domain.User) (*domain.Quiz, error) {
	quiz, err := r.GetQuiz(quizID)

	if err != nil {
		return nil, err
	}
	quiz.Mutex.Lock()
	defer quiz.Mutex.Unlock()

	quiz.Users[user.ID] = user
	return quiz, nil
}

// InitSampleData creates a sample quiz for testing
func (r *InMemoryQuizRepository) InitSampleData() {
	quiz := &domain.Quiz{
		ID: "bTaskee",
		Questions: []domain.Question{
			{
				ID:   1,
				Text: "Khi phát hiện một đoạn code không tối ưu trong hệ thống, một Senior Software Engineer nên làm gì?",
				Options: []string{
					"Viết lại ngay lập tức mà không cần hỏi ai để đảm bảo code sạch hơn",
					"Tạo một ticket trong backlog và đề xuất cải tiến, giải thích lý do cho team",
					"Bỏ qua vì hệ thống vẫn đang chạy tốt, không cần thay đổi nếu không có lỗi",
					"Chỉ sửa trong phạm vi của mình mà không cần thông báo cho ai"},
				Answer:    1,
				TimeLimit: constants.DefaultTimeLimit,
			},
			{
				ID:   2,
				Text: "Yếu tố nào sau đây phản ánh đúng nhất khi đánh giá một Senior Software Engineer?",
				Options: []string{
					"Có khả năng viết code phức tạp, sử dụng nhiều kỹ thuật cao cấp trong mọi tình huống",
					"Thường xuyên làm việc độc lập và tránh giao tiếp với các thành viên khác để tối ưu thời gian",
					"Có thể thiết kế kiến trúc hệ thống, mentor cho người khác và giải quyết vấn đề ở mức tổng thể",
					"Biết sử dụng nhiều framework và công cụ lập trình khác nhau"},
				Answer:    2,
				TimeLimit: constants.DefaultTimeLimit,
			},
			{
				ID:        3,
				Text:      "Bạn có muốn thành một Senior Software Engineer không?",
				Options:   []string{"Có", "Không"},
				Answer:    0,
				TimeLimit: constants.DefaultTimeLimit,
			},
		},
		Users:       make(map[string]*domain.User),
		Leaderboard: []domain.LeaderboardEntry{},
	}

	r.SaveQuiz(quiz)
}
