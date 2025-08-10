package valueobjects

import "errors"

type GameStatus string

const (
	Setup      GameStatus = "Setup"
	InProgress GameStatus = "InProgress"
	Finished   GameStatus = "Finished"
)

func (g GameStatus) String() string {
	return string(g)
}

func (g GameStatus) IsValid() bool {
	switch g {
	case Setup, InProgress, Finished:
		return true
	default:
		return false
	}
}

func NewGameStatus(status string) (GameStatus, error) {
	gameStatus := GameStatus(status)

	if !gameStatus.IsValid() {
		return "", errors.New("invalid game status")
	}

	return gameStatus, nil
}
