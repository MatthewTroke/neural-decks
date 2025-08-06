package valueobjects

import "fmt"

type PlayerRole string

const (
	Participant PlayerRole = "Participant"
	Owner       PlayerRole = "Owner"
)

func NewPlayerRole(role string) (PlayerRole, error) {
	playerRole := PlayerRole(role)

	if !playerRole.IsValid() {
		return "", fmt.Errorf("invalid player role: %s", role)
	}

	return playerRole, nil
}

func (r PlayerRole) IsValid() bool {
	return r == Participant || r == Owner
}

func (r PlayerRole) String() string {
	return string(r)
}

func (r PlayerRole) IsParticipant() bool {
	return r == Participant
}

func (r PlayerRole) IsOwner() bool {
	return r == Owner
}
