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
	GameID string `json:"game_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}

type GameEventPayloadJoinedGameRequest struct {
	GameID string `json:"game_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}

type GameEventPayloadGameRoundContinuedRequest struct {
	GameID      string            `json:"game_id" validate:"required"`
	UserID      string            `json:"user_id" validate:"required"`
	PlayerCards map[string]string `json:"player_cards" validate:"required"`
	BlackCardID string            `json:"black_card_id" validate:"required"`
}

type GameEventPayloadCardCzarChoseWinningCardRequest struct {
	GameID string `json:"game_id" validate:"required"`
	CardID string `json:"card_id" validate:"required"`
}

type GameEventPayloadPlayCardRequest struct {
	GameID string `json:"game_id" validate:"required"`
	CardID string `json:"card_id" validate:"required"`
}

type GameEventPayloadEmojiClickedRequest struct {
	GameID string `json:"game_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
	Emoji  string `json:"emoji" validate:"required"`
}
