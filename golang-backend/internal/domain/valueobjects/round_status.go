package valueobjects

import "fmt"

type RoundStatus string

const (
	Waiting                 RoundStatus = "Waiting"
	PlayersPickingCard      RoundStatus = "PlayersPickingCard"
	JudgePickingWinningCard RoundStatus = "JudgePickingWinningCard"
	JudgeChoseWinningCard   RoundStatus = "JudgeChoseWinningCard"
	GameOver                RoundStatus = "GameOver"
)

func (r RoundStatus) String() string {
	return string(r)
}

func (r RoundStatus) IsValid() bool {
	switch r {
	case Waiting, PlayersPickingCard, JudgePickingWinningCard, JudgeChoseWinningCard, GameOver:
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
		return target == JudgePickingWinningCard
	case JudgePickingWinningCard:
		return target == JudgeChoseWinningCard
	case JudgeChoseWinningCard:
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
	case PlayersPickingCard, JudgePickingWinningCard, JudgeChoseWinningCard:
		return true
	default:
		return false
	}
}

func (r RoundStatus) RequiresPlayerAction() bool {
	return r == PlayersPickingCard
}

func (r RoundStatus) RequiresJudgeAction() bool {
	return r == JudgePickingWinningCard
}

func (r RoundStatus) CanPlayCards() bool {
	return r == PlayersPickingCard
}

func (r RoundStatus) CanPickWinningCard() bool {
	return r == JudgePickingWinningCard
}

func (r RoundStatus) CanContinueRound() bool {
	return r == JudgeChoseWinningCard
}

func (r RoundStatus) GetValidTransitions() []RoundStatus {
	switch r {
	case Waiting:
		return []RoundStatus{PlayersPickingCard}
	case PlayersPickingCard:
		return []RoundStatus{JudgePickingWinningCard}
	case JudgePickingWinningCard:
		return []RoundStatus{JudgeChoseWinningCard}
	case JudgeChoseWinningCard:
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
	case JudgePickingWinningCard:
		return "Judge is picking the winning card"
	case JudgeChoseWinningCard:
		return "Judge has chosen the winning card"
	case GameOver:
		return "Game is over"
	default:
		return "Unknown round status"
	}
}
