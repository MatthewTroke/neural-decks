package valueobjects

import "fmt"

type RoundStatus string

const (
	Waiting                    RoundStatus = "Waiting"
	PlayersPickingCard         RoundStatus = "PlayersPickingCard"
	CardCzarPickingWinningCard RoundStatus = "CardCzarPickingWinningCard"
	CardCzarChoseWinningCard   RoundStatus = "CardCzarChoseWinningCard"
	GameOver                   RoundStatus = "GameOver"
)

func (r RoundStatus) String() string {
	return string(r)
}

func (r RoundStatus) IsValid() bool {
	switch r {
	case Waiting, PlayersPickingCard, CardCzarPickingWinningCard, CardCzarChoseWinningCard, GameOver:
		return true
	default:
		return false
	}
}

func NewRoundStatus(status string) (RoundStatus, error) {
	roundStatus := RoundStatus(status)

	if !roundStatus.IsValid() {
		return "", fmt.Errorf("invalid round status: %s", status)
	}
	return roundStatus, nil
}

func (r RoundStatus) CanTransitionTo(target RoundStatus) bool {
	switch r {
	case Waiting:
		return target == PlayersPickingCard
	case PlayersPickingCard:
		return target == CardCzarPickingWinningCard
	case CardCzarPickingWinningCard:
		return target == CardCzarChoseWinningCard
	case CardCzarChoseWinningCard:
		return target == PlayersPickingCard || target == GameOver
	case GameOver:
		return false
	default:
		return false
	}
}

func (r RoundStatus) IsTerminal() bool {
	return r == GameOver
}

func (r RoundStatus) IsActive() bool {
	switch r {
	case PlayersPickingCard, CardCzarPickingWinningCard, CardCzarChoseWinningCard:
		return true
	default:
		return false
	}
}

func (r RoundStatus) RequiresPlayerAction() bool {
	return r == PlayersPickingCard
}

func (r RoundStatus) RequiresCardCzarAction() bool {
	return r == CardCzarPickingWinningCard
}

func (r RoundStatus) CanPlayCards() bool {
	return r == PlayersPickingCard
}

func (r RoundStatus) CanPickWinningCard() bool {
	return r == CardCzarPickingWinningCard
}

func (r RoundStatus) CanContinueRound() bool {
	return r == CardCzarChoseWinningCard
}

func (r RoundStatus) GetValidTransitions() []RoundStatus {
	switch r {
	case Waiting:
		return []RoundStatus{PlayersPickingCard}
	case PlayersPickingCard:
		return []RoundStatus{CardCzarPickingWinningCard}
	case CardCzarPickingWinningCard:
		return []RoundStatus{CardCzarChoseWinningCard}
	case CardCzarChoseWinningCard:
		return []RoundStatus{PlayersPickingCard, GameOver}
	case GameOver:
		return []RoundStatus{} // No transitions from terminal state
	default:
		return []RoundStatus{}
	}
}

func (r RoundStatus) ValidateTransition(target RoundStatus) error {
	if !r.CanTransitionTo(target) {
		return fmt.Errorf("invalid transition from %s to %s", r, target)
	}
	return nil
}

func (r RoundStatus) GetDescription() string {
	switch r {
	case Waiting:
		return "Waiting for players to join"
	case PlayersPickingCard:
		return "Players are picking their cards"
	case CardCzarPickingWinningCard:
		return "Card Czar is picking the winning card"
	case CardCzarChoseWinningCard:
		return "Card Czar has chosen the winning card"
	case GameOver:
		return "Game is over"
	default:
		return "Unknown round status"
	}
}
