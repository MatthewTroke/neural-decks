package aggregates

import (
	"cardgame/domain/entities"
	"cardgame/domain/valueobjects"
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

type Player struct {
	Score         int                     `json:"score"`
	Role          valueobjects.PlayerRole `json:"role"`
	UserID        string                  `json:"user_id"`
	Name          string                  `json:"name"`
	Image         string                  `json:"image"`
	Deck          []*entities.Card        `json:"deck"`
	IsCardCzar    bool                    `json:"is_card_czar"`
	WasCardCzar   bool                    `json:"was_card_czar"`
	PlacedCard    *entities.Card          `json:"placed_card"`
	IsRoundWinner bool                    `json:"is_round_winner"`
	IsGameWinner  bool                    `json:"is_game_winner"`
	WSConnection  *websocket.Conn         `json:"-"`
}

func NewPlayer(claim *entities.CustomClaim, wsConnection *websocket.Conn) (*Player, error) {
	role, _ := valueobjects.NewPlayerRole(string(valueobjects.Participant))

	if wsConnection == nil {
		return nil, fmt.Errorf("could not create player because websocket connection is nil")
	}

	if !role.IsValid() {
		role = valueobjects.Participant
	}

	return &Player{
		Score:         0,
		Role:          role,
		UserID:        claim.UserID,
		Name:          claim.Name,
		Image:         claim.Image,
		Deck:          []*entities.Card{},
		IsCardCzar:    false,
		WasCardCzar:   false,
		PlacedCard:    nil,
		IsRoundWinner: false,
		IsGameWinner:  false,
		WSConnection:  wsConnection,
	}, nil
}

func (p *Player) SetIsCardCzar(isCardCzar bool) *Player {
	p.IsCardCzar = isCardCzar

	return p
}

func (p *Player) SetWasCardCzar(wasCardCzar bool) *Player {
	p.WasCardCzar = wasCardCzar

	return p
}

func (p *Player) SetPlacedCard(card *entities.Card) error {
	if card == nil {
		return fmt.Errorf("could not set nil card as placed card")
	}

	p.PlacedCard = card

	return nil
}

func (p *Player) RemovePlacedCard() error {
	if p.PlacedCard == nil {
		return fmt.Errorf("could not remove placed card because it is nil")
	}

	p.PlacedCard = nil

	return nil
}

func (p *Player) RemoveCardFromDeck(cardId string) error {
	for i, card := range p.Deck {
		if card.ID == cardId {
			p.Deck = append(p.Deck[:i], p.Deck[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("could not remove card from deck because card with id %s not found", cardId)
}

func (p *Player) HasAlreadyPlayedWhiteCard() bool {
	return p.PlacedCard != nil
}

func (p *Player) SetCardAsPlacedCard(card *entities.Card) error {
	if card == nil {
		return fmt.Errorf("could not set nil card as placed card")
	}

	p.PlacedCard = card

	return nil
}

func (p *Player) HasPlayedCard(card *entities.Card) (bool, error) {
	if card == nil {
		return false, fmt.Errorf("could not check if player has played card because card is nil")
	}

	return p.PlacedCard != nil && p.PlacedCard.ID == card.ID, nil
}

func (p *Player) IncrementScore() {
	p.Score++
}

// Clone creates a deep copy of the player
func (p *Player) Clone() *Player {
	if p == nil {
		return nil
	}

	cloned := &Player{
		Score:         p.Score,
		Role:          p.Role,
		UserID:        p.UserID,
		Name:          p.Name,
		Image:         p.Image,
		IsCardCzar:    p.IsCardCzar,
		WasCardCzar:   p.WasCardCzar,
		IsRoundWinner: p.IsRoundWinner,
		IsGameWinner:  p.IsGameWinner,
	}

	// Clone deck
	if p.Deck != nil {
		cloned.Deck = make([]*entities.Card, len(p.Deck))
		for i, card := range p.Deck {
			cloned.Deck[i] = card.Clone()
		}
	}

	// Clone placed card
	if p.PlacedCard != nil {
		cloned.PlacedCard = p.PlacedCard.Clone()
	}

	return cloned
}

func (p *Player) SendMessage(messageType string, payload interface{}) error {
	if p.WSConnection == nil {
		return fmt.Errorf("player %s has no websocket connection", p.UserID)
	}

	message := struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}{
		Type:    messageType,
		Payload: payload,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message for player %s: %w", p.UserID, err)
	}

	return p.WSConnection.WriteMessage(websocket.TextMessage, messageBytes)
}
