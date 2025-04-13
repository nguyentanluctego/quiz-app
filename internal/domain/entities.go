package domain

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Quiz struct {
	ID          string             `json:"id"`
	Questions   []Question         `json:"questions"`
	Users       map[string]*User   `json:"users"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Mutex       sync.RWMutex       `json:"-"`
}

type Question struct {
	ID        int      `json:"id"`
	Text      string   `json:"text"`
	Options   []string `json:"options"`
	Answer    int      `json:"answer"`
	TimeLimit int      `json:"timeLimit"` // in seconds
}

type User struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Score      int             `json:"score"`
	Conn       *websocket.Conn `json:"-"` // Currently not support multiple connections
	JoinedAt   time.Time       `json:"joinedAt"`
	LastActive time.Time       `json:"lastActive"`
}

type LeaderboardEntry struct {
	UserId string `json:"userId"`
	Score  int    `json:"score"`
}
