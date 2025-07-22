package domain

type CardType string

const (
	Black CardType = "Black"
	White CardType = "White"
)

type Card struct {
	ID        string
	Type      CardType
	CardValue string
}
