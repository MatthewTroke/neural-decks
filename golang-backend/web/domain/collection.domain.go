package domain

import (
	"math/rand"
	"time"
)

type Collection struct {
	Cards []*Card
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
