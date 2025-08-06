package request

import (
	"cardgame/domain/events"
	"encoding/json"
)

type GameEventRequest struct {
	GameID  string               `json:"game_id"`
	Type    events.GameEventType `json:"type"`
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
