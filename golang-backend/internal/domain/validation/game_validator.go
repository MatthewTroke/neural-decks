package validation

type GameValidator interface {
	ValidatePlayCard(gameID, playerID, cardID string) ValidationResult
	ValidateVote(gameID, playerID, winningCardID string) ValidationResult
	ValidateStartRound(gameID string) ValidationResult
	ValidateJoinGame(gameID, playerID string) ValidationResult
	ValidateLeaveGame(gameID, playerID string) ValidationResult
	ValidateEndGame(gameID string) ValidationResult
}
