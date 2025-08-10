package events

type GameEventPayloadGameBegins struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

func NewGameEventPayloadGameBegins(gameID string, playerID string) GameEventPayloadGameBegins {
	return GameEventPayloadGameBegins{
		GameID:   gameID,
		PlayerID: playerID,
	}
}
