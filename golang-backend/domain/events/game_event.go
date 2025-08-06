package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GameEvent struct {
	ID        string          `json:"id"`
	GameID    string          `json:"game_id"`
	Type      GameEventType   `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func NewGameEvent(gameID string, eventType GameEventType, payload json.RawMessage) *GameEvent {
	return &GameEvent{
		ID:        uuid.New().String(),
		GameID:    gameID,
		Type:      eventType,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}

type GameEventType string

const (
	EventGameBegins               GameEventType = "GameBegins"
	EventJoinedGame               GameEventType = "JoinedGame"
	EventCardPlayed               GameEventType = "CardPlayed"
	EventRoundContinued           GameEventType = "RoundContinued"
	EventCardCzarChoseWinningCard GameEventType = "CardCzarChoseWinningCard"
	EventShuffle                  GameEventType = "Shuffle"
	EventDealCards                GameEventType = "DealCards"
	EventDrawBlackCard            GameEventType = "DrawBlackCard"
	EventSetCardCzar              GameEventType = "SetCardCzar"
	EventTimerUpdate              GameEventType = "TimerUpdate"
	EventGameWinner               GameEventType = "GameWinner"
	EventClockUpdate              GameEventType = "ClockUpdate"
	EventEmojiClicked             GameEventType = "EmojiClicked"
)
