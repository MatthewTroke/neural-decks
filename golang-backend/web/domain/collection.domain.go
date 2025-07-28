package domain

import (
	"math/rand"
	"time"
)

type Collection struct {
	Cards []*Card `json:"cards"`
}

func NewCollection() *Collection {
	return &Collection{Cards: []*Card{}}
}

func (c *Collection) AddCard(card *Card) {
	c.Cards = append(c.Cards, card)
}

func (c *Collection) SetCards(cards []*Card) {
	c.Cards = cards
}

func (c *Collection) Shuffle() {
	rand.Seed(time.Now().UnixNano())

	n := len(c.Cards)

	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		c.Cards[i], c.Cards[j] = c.Cards[j], c.Cards[i]
	}
}

func (c *Collection) ShuffleWithSeed(seed int64) {
	rand.Seed(seed)

	n := len(c.Cards)

	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		c.Cards[i], c.Cards[j] = c.Cards[j], c.Cards[i]
	}
}

func (c *Collection) FindCardByID(cardID string) *Card {
	for _, card := range c.Cards {
		if card.ID == cardID {
			return card
		}
	}
	return nil
}

func (c *Collection) DrawCards(n int, cardType CardType) []*Card {
	var drawnCards []*Card
	var remainingCards []*Card

	count := 0

	for _, card := range c.Cards {
		if card.Type == cardType && count < n {
			drawnCards = append(drawnCards, card)
			count++
		} else {
			remainingCards = append(remainingCards, card)
		}
	}

	c.Cards = remainingCards

	return drawnCards
}

// Clone creates a deep copy of the collection
func (c *Collection) Clone() *Collection {
	if c == nil {
		return nil
	}

	cloned := &Collection{}

	if c.Cards != nil {
		cloned.Cards = make([]*Card, len(c.Cards))
		for i, card := range c.Cards {
			cloned.Cards[i] = card.Clone()
		}
	}

	return cloned
}
