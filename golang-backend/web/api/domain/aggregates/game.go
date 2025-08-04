package aggregates

import (
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/events"
	"cardgame/internal/domain/valueobjects"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"encoding/json"

	"gorm.io/gorm"
)

type GameRepository interface {
	Create(game Game) (Game, error)
}

type GameService interface {
	AddGame(game *Game) *Game
	GetAllGames() []*Game
	GetGameById(gameId string) (*Game, error)
}

type InboundWebsocketGameType string

const (
	JoinGame        InboundWebsocketGameType = "JOIN_GAME"
	LeaveGame       InboundWebsocketGameType = "LEAVE_GAME"
	PlayCard        InboundWebsocketGameType = "PLAY_CARD"
	PickWinningCard InboundWebsocketGameType = "PICK_WINNING_CARD"
	ContinueRound   InboundWebsocketGameType = "CONTINUE_ROUND"
	BeginGame       InboundWebsocketGameType = "BEGIN_GAME"
	EmojiClicked    InboundWebsocketGameType = "EMOJI_CLICKED"
)

type OutboundWebsocketGameType string

const (
	GameUpdate          OutboundWebsocketGameType = "GAME_UPDATE"
	ChatMessage         OutboundWebsocketGameType = "CHAT_MESSAGE"
	EmojiClickedMessage OutboundWebsocketGameType = "EMOJI_CLICKED"
)

type GameEventPayloadCardCzarChoseWinningCard struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

func NewGameEventPayloadCardCzarChoseWinningCard(gameID string, cardID string) GameEventPayloadCardCzarChoseWinningCard {
	return GameEventPayloadCardCzarChoseWinningCard{
		GameID: gameID,
		CardID: cardID,
	}
}

type GameEventPayloadGameBegins struct {
	GameID string `json:"game_id"`
	UserID string `json:"user_id"`
}

func NewGameEventPayloadGameBegins(gameID string, userID string) GameEventPayloadGameBegins {
	return GameEventPayloadGameBegins{
		GameID: gameID,
		UserID: userID,
	}
}

type GameEventPayloadJoinedGame struct {
	GameID string                `json:"game_id"`
	UserID string                `json:"user_id"`
	Claim  *entities.CustomClaim `json:"claim"`
}

func NewGameEventPayloadJoinedGame(gameID string, userID string, claim *entities.CustomClaim) GameEventPayloadJoinedGame {
	return GameEventPayloadJoinedGame{
		GameID: gameID,
		UserID: userID,
		Claim:  claim,
	}
}

type GameEventPayloadGameRoundContinuedWithCards struct {
	GameID      string            `json:"game_id"`
	UserID      string            `json:"user_id"`
	PlayerCards map[string]string `json:"player_cards"`  // playerID -> cardID
	BlackCardID string            `json:"black_card_id"` // cardID for the new black card
}

func NewGameEventPayloadGameRoundContinuedWithCards(gameID string, userID string, playerCards map[string]string, blackCardID string) GameEventPayloadGameRoundContinuedWithCards {
	return GameEventPayloadGameRoundContinuedWithCards{
		GameID:      gameID,
		UserID:      userID,
		PlayerCards: playerCards,
		BlackCardID: blackCardID,
	}
}

type GameEventPayloadPlayCard struct {
	GameID string                `json:"game_id"`
	CardID string                `json:"card_id"`
	Claim  *entities.CustomClaim `json:"claim"`
}

func NewGameEventPayloadPlayCard(gameID string, cardID string, claim *entities.CustomClaim) GameEventPayloadPlayCard {
	return GameEventPayloadPlayCard{
		GameID: gameID,
		CardID: cardID,
		Claim:  claim,
	}
}

type GameEventPayloadShuffle struct {
	GameID    string `json:"game_id"`
	Seed      int64  `json:"seed"`
	ShuffleID string `json:"shuffle_id"`
}

func NewGameEventPayloadShuffle(gameID string, seed int64, shuffleID string) GameEventPayloadShuffle {
	return GameEventPayloadShuffle{
		GameID:    gameID,
		Seed:      seed,
		ShuffleID: shuffleID,
	}
}

type GameEventPayloadDealCards struct {
	GameID   string   `json:"game_id"`
	PlayerID string   `json:"player_id"`
	CardIDs  []string `json:"card_ids"`
}

func NewGameEventPayloadDealCards(gameID string, playerID string, cardIDs []string) GameEventPayloadDealCards {
	return GameEventPayloadDealCards{
		GameID:   gameID,
		PlayerID: playerID,
		CardIDs:  cardIDs,
	}
}

type GameEventPayloadDrawBlackCard struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

func NewGameEventPayloadDrawBlackCard(gameID string, cardID string) GameEventPayloadDrawBlackCard {
	return GameEventPayloadDrawBlackCard{
		GameID: gameID,
		CardID: cardID,
	}
}

type GameEventPayloadSetCardCzar struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

func NewGameEventPayloadSetCardCzar(gameID string, playerID string) GameEventPayloadSetCardCzar {
	return GameEventPayloadSetCardCzar{
		GameID:   gameID,
		PlayerID: playerID,
	}
}

type GameEventPayloadClockUpdate struct {
	GameID             string    `json:"game_id"`
	NextAutoProgressAt time.Time `json:"next_auto_progress_at"`
}

func NewGameEventPayloadClockUpdate(gameID string, nextAutoProgressAt time.Time) GameEventPayloadClockUpdate {
	return GameEventPayloadClockUpdate{
		GameID:             gameID,
		NextAutoProgressAt: nextAutoProgressAt,
	}
}

type GameEventPayloadGameWinner struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
	Score    int    `json:"score"`
}

func NewGameEventPayloadGameWinner(gameID string, playerID string, score int) GameEventPayloadGameWinner {
	return GameEventPayloadGameWinner{
		GameID:   gameID,
		PlayerID: playerID,
		Score:    score,
	}
}

type GameEvent struct {
	ID        string               `json:"id"`
	GameID    string               `json:"game_id"`
	Type      events.GameEventType `json:"type"`
	Payload   json.RawMessage      `json:"payload"`
	CreatedAt time.Time            `json:"created_at"`
}

type Game struct {
	Mutex              sync.RWMutex             `json:"-"`
	ID                 string                   `json:"id"`
	Name               string                   `json:"name"`
	Collection         *Collection              `json:"collection"`
	WinnerCount        int                      `json:"winner_count"`
	MaxPlayerCount     int                      `json:"max_player_count"`
	Status             valueobjects.GameStatus  `json:"status"`
	Players            []*Player                `json:"players"`
	WhiteCards         []*entities.Card         `json:"white_cards"`
	BlackCard          *entities.Card           `json:"black_card"`
	RoundStatus        valueobjects.RoundStatus `json:"round_status"`
	CurrentGameRound   int                      `json:"current_game_round"`
	RoundWinner        *Player                  `json:"round_winner"`
	LastVacatedAt      *time.Time               `json:"last_vacated_at"`
	LastEventAt        *time.Time               `json:"last_event_at"`
	NextAutoProgressAt *time.Time               `json:"next_auto_progress_at"` // Timestamp when next auto-progress will happen
	Vacated            bool                     `json:"vacated"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
	DeletedAt          gorm.DeletedAt           `json:"-"`
}

func (g *Game) Lock() {
	g.Mutex.RLock()
}

func (g *Game) Unlock() {
	g.Mutex.RUnlock()
}

func (g *Game) SetStatus(status valueobjects.GameStatus) {
	g.Status = status
}

func (g *Game) SetRoundStatus(status valueobjects.RoundStatus) {
	g.RoundStatus = status
}

func (g *Game) SetRoundWinner(player *Player) {
	g.RoundWinner = player
}

func (g *Game) SetLastEventAt(t time.Time) {
	g.LastEventAt = &t
}

func (g *Game) SetNextAutoProgressAt(t time.Time) {
	g.NextAutoProgressAt = &t
}

func (g *Game) AddPlayer(player *Player) error {
	if g.Players == nil {
		return fmt.Errorf("could not add player, players on game are nil")
	}

	if player == nil {
		return fmt.Errorf("could not add nil player to game")
	}

	g.Players = append(g.Players, player)

	return nil
}

func (g *Game) RemoveWasCardCzarFromAllPlayers() error {
	if len(g.Players) == 0 {
		return fmt.Errorf("could not remove card czar from all players, no players exist")
	}

	for _, player := range g.Players {
		player.WasCardCzar = false
	}

	return nil
}

func (g *Game) ClearBoard() {
	g.WhiteCards = []*entities.Card{}
	g.BlackCard = nil
}

func (g *Game) PickNewBlackCard() error {
	if g.Collection == nil {
		return fmt.Errorf("could not pick new black card because collection is nil")
	}

	g.BlackCard = g.Collection.DrawCards(1, entities.Black)[0]

	return nil
}

func (g *Game) IncrementGameRound() {
	g.CurrentGameRound++
}

func (g *Game) FindNewCardCzar() (*Player, error) {
	if len(g.Players) == 0 {
		return nil, fmt.Errorf("could not find new card czar because no players in game %s", g.ID)
	}

	for _, player := range g.Players {
		if !player.WasCardCzar && !player.IsCardCzar {
			return player, nil
		}
	}

	return nil, errors.New("no eligible player found to be the new card czar")
}

func (g *Game) FindCurrentCardCzar() (*Player, error) {
	if len(g.Players) == 0 {
		return nil, fmt.Errorf("could not find current card czar because no players in game %s", g.ID)
	}

	for _, player := range g.Players {
		if player.IsCardCzar {
			return player, nil
		}
	}

	return nil, fmt.Errorf("could not find current card czar in game %s", g.ID)
}

func (g *Game) PickNewCardCzar() error {
	if len(g.Players) == 0 {
		return fmt.Errorf("could not pick a new card czar, no player length")
	}

	canPromoteNewCardCzar := false

	for _, player := range g.Players {
		if !player.WasCardCzar {
			canPromoteNewCardCzar = true
			break
		}
	}

	if !canPromoteNewCardCzar {
		err := g.RemoveWasCardCzarFromAllPlayers()

		if err != nil {
			return fmt.Errorf("error picking new card czar: %w", err)
		}
	}

	player, err := g.FindNewCardCzar()

	if err != nil {
		return fmt.Errorf("error picking new card czar: %w", err)
	}

	player.SetIsCardCzar(true)

	return nil
}

// THERE IS A BUG HERE. THIS SHOULD BE FIND PLAYER BY USER ID AND GAME ID.
func (g *Game) FindPlayerByUserId(userId string) (*Player, error) {
	if g.Players == nil {
		return nil, fmt.Errorf("could not find player by user id, players are nil")
	}

	if len(g.Players) == 0 {
		return nil, fmt.Errorf("could not find player by user id because no players in game %s", g.ID)
	}

	if userId == "" {
		return nil, fmt.Errorf("could not find player by user id because user id is empty")
	}

	for i := range g.Players {
		if g.Players[i].UserID == userId {
			return g.Players[i], nil
		}
	}

	return nil, fmt.Errorf("could not find player by user id because user id %s does not exist in game %s", userId, g.ID)
}

func (g *Game) FindCardByPlayerId(playerId string, cardId string) (*entities.Card, error) {
	player, err := g.FindPlayerByUserId(playerId)

	if err != nil {
		log.Printf("could not find card by player id: %v", err)
		return nil, err
	}

	for _, card := range player.Deck {
		if card.ID == cardId {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card not found in player's deck")
}

func (g *Game) FindWhiteCardByCardId(cardId string) (*entities.Card, error) {
	for _, card := range g.WhiteCards {
		if card.ID == cardId {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card not found in cards")
}

func (g *Game) AddWhiteCardToGameBoard(card *entities.Card) error {
	if card == nil {
		return fmt.Errorf("could not add nil card to game board")
	}

	g.WhiteCards = append(g.WhiteCards, card)

	return nil
}

func (g *Game) HasAllPlayersPlayedWhiteCard() (bool, error) {
	if g.Players == nil {
		return false, fmt.Errorf("could not check if all players have played white card because players are nil")
	}

	if len(g.Players) == 0 {
		return false, nil
	}

	for _, player := range g.Players {
		if player.PlacedCard == nil && !player.IsCardCzar {
			return false, nil
		}
	}

	return true, nil
}

func (g *Game) FindWhiteCardOwner(card *entities.Card) (*Player, error) {
	if card == nil {
		return nil, fmt.Errorf("could not find white card owner because card is nil")
	}

	for _, player := range g.Players {
		if player.IsCardCzar {
			continue
		}

		if player.PlacedCard.ID == card.ID {
			return player, nil
		}
	}

	return nil, fmt.Errorf("could not find white card owner")
}

func (g *Game) CheckForWinner() *Player {
	for _, player := range g.Players {
		if player.Score >= g.WinnerCount {
			return player
		}
	}
	return nil
}

func (g *Game) Clone() *Game {
	g.Lock()
	defer g.Unlock()

	cloned := &Game{
		ID:               g.ID,
		Name:             g.Name,
		WinnerCount:      g.WinnerCount,
		MaxPlayerCount:   g.MaxPlayerCount,
		Status:           g.Status,
		RoundStatus:      g.RoundStatus,
		CurrentGameRound: g.CurrentGameRound,
		Vacated:          g.Vacated,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		DeletedAt:        g.DeletedAt,
	}

	// Clone collection
	if g.Collection != nil {
		cloned.Collection = g.Collection.Clone()
	}

	// Clone players
	if g.Players != nil {
		cloned.Players = make([]*Player, len(g.Players))
		for i, player := range g.Players {
			cloned.Players[i] = player.Clone()
		}
	}

	// Clone white cards
	if g.WhiteCards != nil {
		cloned.WhiteCards = make([]*entities.Card, len(g.WhiteCards))
		for i, card := range g.WhiteCards {
			cloned.WhiteCards[i] = card.Clone()
		}
	}

	// Clone black card
	if g.BlackCard != nil {
		cloned.BlackCard = g.BlackCard.Clone()
	}

	// Clone round winner
	if g.RoundWinner != nil {
		cloned.RoundWinner = g.RoundWinner.Clone()
	}

	// Clone last vacated time
	if g.LastVacatedAt != nil {
		lastVacated := *g.LastVacatedAt
		cloned.LastVacatedAt = &lastVacated
	}

	// Clone last event time
	if g.LastEventAt != nil {
		lastEvent := *g.LastEventAt
		cloned.LastEventAt = &lastEvent
	}

	// Clone next auto progress time
	if g.NextAutoProgressAt != nil {
		nextAutoProgress := *g.NextAutoProgressAt
		cloned.NextAutoProgressAt = &nextAutoProgress
	}

	return cloned
}

func (g *Game) ApplyEvent(event GameEvent) error {
	g.Lock()
	defer g.Unlock()

	g.SetLastEventAt(event.CreatedAt)

	switch event.Type {
	case events.EventGameBegins:
		var payload GameEventPayloadGameBegins

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameBegins payload: %w", err)
		}

		g.SetRoundStatus(valueobjects.PlayersPickingCard)
		g.SetStatus(valueobjects.InProgress)
	case events.EventShuffle:
		var payload GameEventPayloadShuffle

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventShuffle payload: %w", err)
		}

		// Use the stored seed for deterministic shuffling
		g.Collection.ShuffleWithSeed(payload.Seed)
	case events.EventDealCards:
		var payload GameEventPayloadDealCards

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventDealCards payload: %w", err)
		}

		// Find the player and give them the specific cards
		for _, player := range g.Players {
			if player.UserID == payload.PlayerID {
				// Clear existing deck and add the specific cards
				player.Deck = []*entities.Card{}
				for _, cardID := range payload.CardIDs {
					card := g.Collection.FindCardByID(cardID)
					if card != nil {
						player.Deck = append(player.Deck, card)
					}
				}
				break
			}
		}
	case events.EventDrawBlackCard:
		var payload GameEventPayloadDrawBlackCard

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventDrawBlackCard payload: %w", err)
		}

		// Set the specific black card that was drawn
		card := g.Collection.FindCardByID(payload.CardID)
		if card != nil {
			g.BlackCard = card
		}
	case events.EventSetCardCzar:
		var payload GameEventPayloadSetCardCzar

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventSetCardCzar payload: %w", err)
		}

		// Find the player and set them as card czar
		for _, player := range g.Players {
			if player.UserID == payload.PlayerID {
				// Remove card czar from all players first
				for _, p := range g.Players {
					p.SetIsCardCzar(false)
				}
				// Set this player as card czar
				player.SetIsCardCzar(true)
				break
			}
		}
	case events.EventJoinedGame:
		var payload GameEventPayloadJoinedGame

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameBegins payload: %w", err)
		}

		player := NewPlayer(payload.Claim)

		if g.Players == nil {
			return fmt.Errorf("could not join game ID %s, game players is nil", payload.GameID)
		}

		g.Players = append(g.Players, player)
	case events.EventRoundContinued:
		var payload GameEventPayloadGameRoundContinuedWithCards

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameBegins payload: %w", err)
		}

		// Remove placed cards and give new cards based on payload
		for _, player := range g.Players {
			if player.IsCardCzar {
				continue
			}

			player.RemovePlacedCard()

			// Check if this player should get a new card
			if cardID, exists := payload.PlayerCards[player.UserID]; exists {
				// Find the card in the collection
				for _, card := range g.Collection.Cards {
					if card.ID == cardID {
						player.Deck = append(player.Deck, card)
						break
					}
				}
			}
		}

		currentCardCzar, err := g.FindCurrentCardCzar()

		if err != nil {
			return fmt.Errorf("could not continue round: %w", err)
		}

		currentCardCzar.SetIsCardCzar(false)
		currentCardCzar.SetWasCardCzar(true)

		g.SetRoundWinner(nil)
		g.ClearBoard()

		err = g.PickNewCardCzar()

		if err != nil {
			return fmt.Errorf("could not pick new card czar: %w", err)
		}

		// Set the new black card based on the payload
		if payload.BlackCardID != "" {
			// Find the black card in the collection
			for _, card := range g.Collection.Cards {
				if card.ID == payload.BlackCardID {
					g.BlackCard = card
					break
				}
			}
		}

		g.IncrementGameRound()
		g.SetRoundStatus(valueobjects.PlayersPickingCard)

	case events.EventCardCzarChoseWinningCard:
		var payload GameEventPayloadCardCzarChoseWinningCard

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventCardCzarChoseWinningCard payload: %w", err)
		}
		winningCard, _ := g.FindWhiteCardByCardId(payload.CardID)
		winner, _ := g.FindWhiteCardOwner(winningCard)

		winner.IncrementScore()
		g.SetRoundStatus(valueobjects.CardCzarChoseWinningCard)
		g.SetRoundWinner(winner)

	case events.EventCardPlayed:
		var payload GameEventPayloadPlayCard

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventCardPlayed payload: %w", err)
		}

		fmt.Printf("EventCardPlayed: UserID=%s, CardID=%s\n", payload.Claim.UserID, payload.CardID)

		player, _ := g.FindPlayerByUserId(payload.Claim.UserID)
		card, _ := g.FindCardByPlayerId(payload.Claim.UserID, payload.CardID)

		fmt.Printf("Found player: %s, Found card: %s\n", player.UserID, card.ID)

		err := player.RemoveCardFromDeck(card.ID)

		if err != nil {
			return fmt.Errorf("unable to play white card: %w", err)
		}

		err = player.SetCardAsPlacedCard(card)

		if err != nil {
			return fmt.Errorf("unable to play white card: %w", err)
		}

		err = g.AddWhiteCardToGameBoard(card)

		if err != nil {
			return fmt.Errorf("unable to play white card: %w", err)
		}

		fmt.Printf("Added card to game board. WhiteCards count: %d\n", len(g.WhiteCards))

		hasPlayersPlayedWhiteCard, err := g.HasAllPlayersPlayedWhiteCard()

		if err != nil {
			return fmt.Errorf("unable to play white card: %w", err)
		}

		if hasPlayersPlayedWhiteCard {
			g.SetRoundStatus(valueobjects.CardCzarPickingWinningCard)
		}
	case events.EventGameWinner:
		var payload GameEventPayloadGameWinner

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameWinner payload: %w", err)
		}

		// Find the winning player and mark them as game winner
		for _, player := range g.Players {
			if player.UserID == payload.PlayerID {
				player.IsGameWinner = true
				break
			}
		}

		// Set game status to finished
		g.SetRoundStatus(valueobjects.GameOver)
		g.SetStatus(valueobjects.Finished)

	case events.EventClockUpdate:
		var payload GameEventPayloadClockUpdate

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventClockUpdate payload: %w", err)
		}

		g.NextAutoProgressAt = &payload.NextAutoProgressAt
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	return nil
}
