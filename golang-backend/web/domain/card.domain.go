package domain

type CardType string

const (
	Black CardType = "Black"
	White CardType = "White"
)

type Card struct {
	ID        string   `json:"id"`
	Type      CardType `json:"type"`
	CardValue string   `json:"card_value"`
}

// Clone creates a copy of the card
func (c *Card) Clone() *Card {
	if c == nil {
		return nil
	}
	return &Card{
		ID:        c.ID,
		Type:      c.Type,
		CardValue: c.CardValue,
	}
}
