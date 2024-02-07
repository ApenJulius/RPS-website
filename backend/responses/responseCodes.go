package responses

type GameStatus int

const (
	GameFound     GameStatus = 10
	PlayerJoined  GameStatus = 20
	PlayerLeft    GameStatus = 30
	GameStarted   GameStatus = 40
	GameCountdown GameStatus = 41
	// other statuses...
)
