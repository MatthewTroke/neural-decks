package request

type JoinGameRequest struct {
	GameID string `validate:"required"`
}

// TODO set constraits to make sure these fields meet the frontend requirements.
type CreateGameRequest struct {
	Name           string `json:"name" validate:"required"`
	WinnerCount    int    `json:"winner_count" validate:"required,numeric"`
	MaxPlayerCount int    `json:"max_player_count" validate:"required,numeric"`
	Subject        string `json:"subject" validate:"required"`
}
