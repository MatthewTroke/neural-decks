package aggregates

import (
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/valueobjects"
	"fmt"

	"github.com/google/uuid"
)

type Player struct {
	ID            string                  `json:"id"`
	Score         int                     `json:"score"`
	Role          valueobjects.PlayerRole `json:"role"`
	IsOwner       bool                    `json:"is_owner"`
	UserID        string                  `json:"user_id"`
	Name          string                  `json:"name"`
	Image         string                  `json:"image"`
	Deck          []*entities.Card        `json:"deck"`
	IsJudge       bool                    `json:"is_judge"`
	WasJudge      bool                    `json:"was_judge"`
	PlacedCard    *entities.Card          `json:"placed_card"`
	IsRoundWinner bool                    `json:"is_round_winner"`
	IsGameWinner  bool                    `json:"is_game_winner"`
}

func NewPlayer(claim *entities.CustomClaim) (*Player, error) {
	role, _ := valueobjects.NewPlayerRole(string(valueobjects.Participant))

	if !role.IsValid() {
		role = valueobjects.Participant
	}

	return &Player{
		ID:            uuid.NewString(),
		Score:         0,
		Role:          role,
		IsOwner:       false,
		UserID:        claim.UserID,
		Name:          claim.Name,
		Image:         claim.Image,
		Deck:          []*entities.Card{},
		IsJudge:       false,
		WasJudge:      false,
		PlacedCard:    nil,
		IsRoundWinner: false,
		IsGameWinner:  false,
	}, nil
}

func (p *Player) SetIsOwner(isOwner bool) *Player {
	p.IsOwner = isOwner

	return p
}

func (p *Player) SetIsGameWinner(isGameWinner bool) *Player {
	p.IsGameWinner = isGameWinner

	return p
}

func (p *Player) SetIsJudge(isJudge bool) *Player {
	p.IsJudge = isJudge

	return p
}

func (p *Player) SetWasJudge(wasJudge bool) *Player {
	p.WasJudge = wasJudge

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
		IsJudge:       p.IsJudge,
		WasJudge:      p.WasJudge,
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
