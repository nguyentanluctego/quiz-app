package constants

// WebSocket message types
const (
	MessageTypeJoin        = "join"
	MessageTypeAnswer      = "answer"
	MessageTypeJoined      = "joined"
	MessageTypeResult      = "result"
	MessageTypeError       = "error"
	MessageTypeLeaderboard = "leaderboard"
)

const (
	PointsPerCorrectAnswer = 10
	DefaultTimeLimit       = 30
)
