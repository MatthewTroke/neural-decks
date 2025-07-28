package request

import (
	"cardgame/domain"
	"encoding/json"
)

// TODO set constraits to make sure these fields meet the frontend requirements.
type CreateGameRequest struct {
	Name           string `json:"name" validate:"required"`
	WinnerCount    int    `json:"winner_count" validate:"required,numeric"`
	MaxPlayerCount int    `json:"max_player_count" validate:"required,numeric"`
	Subject        string `json:"subject" validate:"required"`
}

type GameEventRequest struct {
	GameID  string               `json:"game_id"`
	Type    domain.GameEventType `json:"type"`
	Payload json.RawMessage      `json:"payload"`
}

type GameEventPayloadGameBeginsRequest struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type GameEventPayloadJoinedGameRequest struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type GameEventPayloadGameRoundContinuedRequest struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

type GameEventPayloadCardCzarChoseWinningCardRequest struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

type GameEventPayloadPlayCardRequest struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}
