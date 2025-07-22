package domain

import (
	"fmt"
)

type PlayerRole string

const (
	Participant PlayerRole = "Participant"
	Owner       PlayerRole = "Owner"
)

type Player struct {
	Score         int
	Role          PlayerRole
	UserID        string
	Name          string
	Image         string
	Deck          []*Card
	IsCardCzar    bool
	WasCardCzar   bool
	PlacedCard    *Card
	IsRoundWinner bool
	IsGameWinner  bool
}

func NewPlayer(claim *CustomClaim) *Player {
	return &Player{
		Score:         0,
		Role:          "Participant",
		UserID:        claim.UserID,
		Name:          claim.Name,
		Image:         claim.Image,
		Deck:          []*Card{},
		IsCardCzar:    false,
		WasCardCzar:   false,
		PlacedCard:    nil,
		IsRoundWinner: false,
		IsGameWinner:  false,
	}
}

func (p *Player) SetIsCardCzar(isCardCzar bool) *Player {
	p.IsCardCzar = isCardCzar

	return p
}

func (p *Player) SetWasCardCzar(wasCardCzar bool) *Player {
	p.WasCardCzar = wasCardCzar

	return p
}

func (p *Player) SetPlacedCard(card *Card) error {
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

func (p *Player) SetCardAsPlacedCard(card *Card) error {
	if card == nil {
		return fmt.Errorf("could not set nil card as placed card")
	}

	p.PlacedCard = card

	return nil
}

func (p *Player) HasPlayedCard(card *Card) (bool, error) {
	if card == nil {
		return false, fmt.Errorf("could not check if player has played card because card is nil")
	}

	return p.PlacedCard != nil && p.PlacedCard.ID == card.ID, nil
}

func (p *Player) IncrementScore() {
	p.Score++
}
