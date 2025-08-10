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

	"github.com/google/uuid"
)

type GameEventType string

const (
	EventGameBegins            GameEventType = "GameBegins"
	EventJoinedGame            GameEventType = "JoinedGame"
	EventCardPlayed            GameEventType = "CardPlayed"
	EventRoundContinued        GameEventType = "RoundContinued"
	EventJudgeChoseWinningCard GameEventType = "JudgeChoseWinningCard"
	EventShuffle               GameEventType = "Shuffle"
	EventDealCards             GameEventType = "DealCards"
	EventDrawBlackCard         GameEventType = "DrawBlackCard"
	EventSetJudge              GameEventType = "SetJudge"
	EventTimerUpdate           GameEventType = "TimerUpdate"
	EventGameWinner            GameEventType = "GameWinner"
	EventClockUpdate           GameEventType = "ClockUpdate"
	EventEmojiClicked          GameEventType = "EmojiClicked"
)

type GameEvent struct {
	ID        string          `json:"id"`
	GameID    string          `json:"game_id"`
	Type      GameEventType   `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func NewGameEvent(gameID string, eventType GameEventType, payload json.RawMessage) *GameEvent {
	return &GameEvent{
		ID:        uuid.New().String(),
		GameID:    gameID,
		Type:      eventType,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
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

type GameEventPayloadJudgeChoseWinningCard struct {
	GameID string `json:"game_id"`
	CardID string `json:"card_id"`
}

func NewGameEventPayloadJudgeChoseWinningCard(gameID string, cardID string) GameEventPayloadJudgeChoseWinningCard {
	return GameEventPayloadJudgeChoseWinningCard{
		GameID: gameID,
		CardID: cardID,
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

type GameEventPayloadSetJudge struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

func NewGameEventPayloadSetJudge(gameID string, playerID string) GameEventPayloadSetJudge {
	return GameEventPayloadSetJudge{
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
	UsedCards          []*entities.Card         `json:"used_cards"`
	BlackCard          *entities.Card           `json:"black_card"`
	RoundStatus        valueobjects.RoundStatus `json:"round_status"`
	CurrentGameRound   int                      `json:"current_game_round"`
	RoundWinner        *Player                  `json:"round_winner"`
	LastVacatedAt      time.Time                `json:"last_vacated_at"`
	LastEventAt        time.Time                `json:"last_event_at"`
	NextAutoProgressAt time.Time                `json:"next_auto_progress_at"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
	DeletedAt          time.Time                `json:"-"`
}

func NewGame(
	id string,
	name string,
	collection *Collection,
	winnerCount int,
	maxPlayerCount int,
	status valueobjects.GameStatus,
	players []*Player,
	whiteCards []*entities.Card,
	usedCards []*entities.Card,
	blackCard *entities.Card,
	roundStatus valueobjects.RoundStatus,
	currentGameRound int,
	roundWinner *Player,
	lastVacatedAt time.Time,
	lastEventAt time.Time,
	nextAutoProgressAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt time.Time,
) *Game {
	return &Game{
		ID:                 id,
		Name:               name,
		Collection:         collection,
		WinnerCount:        winnerCount,
		MaxPlayerCount:     maxPlayerCount,
		Status:             status,
		Players:            players,
		WhiteCards:         whiteCards,
		UsedCards:          usedCards,
		BlackCard:          blackCard,
		RoundStatus:        roundStatus,
		CurrentGameRound:   currentGameRound,
		RoundWinner:        roundWinner,
		LastVacatedAt:      lastVacatedAt,
		LastEventAt:        lastEventAt,
		NextAutoProgressAt: nextAutoProgressAt,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		DeletedAt:          deletedAt,
	}
}

func (g *Game) Lock() {
	g.Mutex.RLock()
}

func (g *Game) Unlock() {
	g.Mutex.RUnlock()
}

func (g *Game) ClearUsedCards() {
	g.Lock()
	defer g.Unlock()

	g.UsedCards = []*entities.Card{}
}

func (g *Game) ShouldShuffle() bool {
	g.Lock()
	defer g.Unlock()

	availableWhiteCards := len(g.GetUnplayedWhiteCards())
	availableBlackCards := len(g.GetUnplayedBlackCards())
	playersNeedingCards := len(g.GetNonJudgePlayers())

	return availableWhiteCards < playersNeedingCards || availableBlackCards < 1
}

func (g *Game) GetNonJudgePlayers() []*Player {
	g.Lock()
	defer g.Unlock()

	nonJudgePlayers := []*Player{}

	for _, player := range g.Players {
		if !player.IsJudge {
			nonJudgePlayers = append(nonJudgePlayers, player)
		}
	}

	return nonJudgePlayers
}

func (g *Game) GetUnplayedBlackCards() []*entities.Card {
	g.Lock()
	defer g.Unlock()

	usedCardsMap := make(map[string]*entities.Card)

	for _, card := range g.UsedCards {
		usedCardsMap[card.ID] = card
	}

	unusedBlackCards := []*entities.Card{}

	for _, card := range g.Collection.Cards {
		if card.Type == "Black" && usedCardsMap[card.ID] == nil {
			unusedBlackCards = append(unusedBlackCards, card)
		}
	}

	return unusedBlackCards
}

func (g *Game) GetUnplayedWhiteCards() []*entities.Card {
	g.Lock()
	defer g.Unlock()

	usedCardsMap := make(map[string]*entities.Card)

	for _, card := range g.UsedCards {
		usedCardsMap[card.ID] = card
	}

	unusedWhiteCards := []*entities.Card{}

	for _, card := range g.Collection.Cards {
		if card.Type == "White" && usedCardsMap[card.ID] == nil {
			unusedWhiteCards = append(unusedWhiteCards, card)
		}
	}

	return unusedWhiteCards
}

func (g *Game) SetStatus(status valueobjects.GameStatus) {
	g.Lock()
	defer g.Unlock()

	g.Status = status
}

func (g *Game) SetRoundStatus(status valueobjects.RoundStatus) {
	g.Lock()
	defer g.Unlock()

	g.RoundStatus = status
}

func (g *Game) SetRoundWinner(player *Player) {
	g.Lock()
	defer g.Unlock()

	g.RoundWinner = player
}

func (g *Game) SetLastEventAt(t time.Time) {
	g.Lock()
	defer g.Unlock()

	g.LastEventAt = t
}

func (g *Game) SetNextAutoProgressAt(t time.Time) {
	g.Lock()
	defer g.Unlock()

	g.NextAutoProgressAt = t
}

func (g *Game) SetBlackCard(card *entities.Card) {
	g.Lock()
	defer g.Unlock()

	g.BlackCard = card
}

func (g *Game) SetWhiteCards(cards []*entities.Card) {
	g.Lock()
	defer g.Unlock()

	g.WhiteCards = cards
}

func (g *Game) IncrementGameRound() {
	g.Lock()
	defer g.Unlock()

	g.CurrentGameRound++
}

func (g *Game) HasPlayers() bool {
	g.Lock()
	defer g.Unlock()

	if g.Players == nil {
		return false
	}

	return len(g.Players) > 0
}

func (g *Game) HasCollection() bool {
	g.Lock()
	defer g.Unlock()

	return g.Collection != nil
}

func (g *Game) IsInProgress() bool {
	g.Lock()
	defer g.Unlock()

	return g.Status == valueobjects.InProgress
}

func (g *Game) IsInSetup() bool {
	g.Lock()
	defer g.Unlock()

	return g.Status == valueobjects.Setup
}

func (g *Game) AddPlayer(player *Player) error {
	g.Lock()
	defer g.Unlock()

	if player == nil {
		return fmt.Errorf("could not add nil player to game")
	}

	if !g.HasPlayers() {
		return fmt.Errorf("could not add player, players on game are nil")
	}

	g.Players = append(g.Players, player)

	return nil
}

func (g *Game) RemovePlayer(player *Player) error {
	g.Lock()
	defer g.Unlock()

	if !g.HasPlayers() {
		return fmt.Errorf("could not remove player, no players exist")
	}

	for i, p := range g.Players {
		if p.UserID == player.UserID {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("could not remove player, player not found")
}

func (g *Game) RemoveWasJudgeFromAllPlayers() error {
	g.Lock()
	defer g.Unlock()

	if !g.HasPlayers() {
		return fmt.Errorf("could not remove judge from all players, no players exist")
	}

	for _, player := range g.Players {
		player.WasJudge = false
	}

	return nil
}

func (g *Game) ClearBoard() {
	g.Lock()
	defer g.Unlock()

	g.SetWhiteCards([]*entities.Card{})
	g.SetBlackCard(nil)
	g.SetRoundStatus(valueobjects.Waiting)
	g.SetRoundWinner(nil)
}

func (g *Game) PickNewBlackCard() error {
	g.Lock()
	defer g.Unlock()

	if !g.HasCollection() {
		return fmt.Errorf("could not pick new black card because collection is nil")
	}

	g.SetBlackCard(g.Collection.DrawCards(1, entities.Black)[0])

	return nil
}

func (g *Game) FindNewJudge() (*Player, error) {
	if !g.HasPlayers() {
		return nil, fmt.Errorf("could not find new judge because no players in game %s", g.ID)
	}

	for _, player := range g.Players {
		if !player.WasJudge && !player.IsJudge {
			return player, nil
		}
	}

	return nil, errors.New("no eligible player found to be the new judge")
}

func (g *Game) FindCurrentJudge() (*Player, error) {
	if !g.HasPlayers() {
		return nil, fmt.Errorf("could not find current judge because no players in game %s", g.ID)
	}

	for _, player := range g.Players {
		if player.IsJudge {
			return player, nil
		}
	}

	return nil, fmt.Errorf("could not find current judge in game %s", g.ID)
}

func (g *Game) WasAllPlayersJudge() bool {
	g.Lock()
	defer g.Unlock()

	for _, player := range g.Players {
		if !player.WasJudge {
			return false
		}
	}

	return true
}

func (g *Game) PickNewJudge() error {
	g.Lock()
	defer g.Unlock()

	if !g.HasPlayers() {
		return fmt.Errorf("could not pick a new judge, no player length")
	}

	if g.WasAllPlayersJudge() {
		err := g.RemoveWasJudgeFromAllPlayers()

		if err != nil {
			return fmt.Errorf("error removing judge from all players: %w", err)
		}
	}

	player, err := g.FindNewJudge()

	if err != nil {
		return fmt.Errorf("error picking new judge: %w", err)
	}

	player.SetIsJudge(true)

	return nil
}

func (g *Game) FindPlayerByUserId(userId string) (*Player, error) {
	g.Lock()
	defer g.Unlock()

	if !g.HasPlayers() {
		return nil, fmt.Errorf("could not find player by user id in game %s, no players in game", g.ID)
	}

	for _, player := range g.Players {
		if player.UserID == userId {
			return player, nil
		}
	}

	return nil, fmt.Errorf("could not find player by user id because user id %s does not exist in game %s", userId, g.ID)
}

func (g *Game) FindCardByPlayerId(playerId string, cardId string) (*entities.Card, error) {
	g.Lock()
	defer g.Unlock()

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
	g.Lock()
	defer g.Unlock()

	for _, card := range g.WhiteCards {
		if card.ID == cardId {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card not found in cards")
}

func (g *Game) AddWhiteCardToGameBoard(card *entities.Card) error {
	g.Lock()
	defer g.Unlock()

	if card == nil {
		return fmt.Errorf("could not add nil card to game board")
	}

	g.WhiteCards = append(g.WhiteCards, card)

	return nil
}

func (g *Game) HasAllPlayersPlayedWhiteCard() (bool, error) {
	if !g.HasPlayers() {
		return false, fmt.Errorf("could not check if all players have played white card because no players in game")
	}

	for _, player := range g.Players {
		if player.PlacedCard == nil && !player.IsJudge {
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
		if player.IsJudge {
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

func (g *Game) ApplyEvent(event *GameEvent) error {
	g.Lock()
	defer g.Unlock()

	g.SetLastEventAt(event.CreatedAt)

	switch event.Type {
	case EventGameBegins:
		var payload events.GameEventPayloadGameBegins

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameBegins payload: %w", err)
		}

		g.SetRoundStatus(valueobjects.PlayersPickingCard)
		g.SetStatus(valueobjects.InProgress)
	case EventShuffle:
		var payload GameEventPayloadShuffle

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventShuffle payload: %w", err)
		}

		g.ClearUsedCards()
		g.Collection.ShuffleWithSeed(payload.Seed)
	case EventDealCards:
		var payload GameEventPayloadDealCards

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventDealCards payload: %w", err)
		}

		for _, player := range g.Players {
			if player.UserID == payload.PlayerID {
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
	case EventDrawBlackCard:
		var payload GameEventPayloadDrawBlackCard

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventDrawBlackCard payload: %w", err)
		}

		card := g.Collection.FindCardByID(payload.CardID)

		if card != nil {
			g.SetBlackCard(card)
		}
	case EventSetJudge:
		var payload GameEventPayloadSetJudge

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventSetJudge payload: %w", err)
		}

		for _, player := range g.Players {
			if player.UserID == payload.PlayerID {
				// Remove judge from all players first
				for _, p := range g.Players {
					p.SetIsJudge(false)
				}
				// Set this player as judge
				player.SetIsJudge(true)
				break
			}
		}
	case EventJoinedGame:
		var payload GameEventPayloadJoinedGame

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameBegins payload: %w", err)
		}

		player, err := NewPlayer(payload.Claim)

		if err != nil {
			return fmt.Errorf("could not join game ID %s, error creating player: %w", payload.GameID, err)
		}

		if g.Players == nil {
			return fmt.Errorf("could not join game ID %s, game players is nil", payload.GameID)
		}

		g.Players = append(g.Players, player)
	case EventRoundContinued:
		var payload GameEventPayloadGameRoundContinuedWithCards

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameBegins payload: %w", err)
		}

		for _, player := range g.Players {
			if player.IsJudge {
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

		judge, err := g.FindCurrentJudge()

		if err != nil {
			return fmt.Errorf("could not continue round: %w", err)
		}

		judge.SetIsJudge(false)
		judge.SetWasJudge(true)

		g.SetRoundWinner(nil)
		g.ClearBoard()

		err = g.PickNewJudge()

		if err != nil {
			return fmt.Errorf("could not pick new judge: %w", err)
		}

		if payload.BlackCardID != "" {
			for _, card := range g.Collection.Cards {
				if card.ID == payload.BlackCardID {
					g.SetBlackCard(card)
					break
				}
			}
		}

		g.IncrementGameRound()
		g.SetRoundStatus(valueobjects.PlayersPickingCard)

	case EventJudgeChoseWinningCard:
		var payload GameEventPayloadJudgeChoseWinningCard

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventJudgeChoseWinningCard payload: %w", err)
		}

		winningCard, err := g.FindWhiteCardByCardId(payload.CardID)

		if err != nil {
			return fmt.Errorf("could not find winning card: %w", err)
		}

		winner, err := g.FindWhiteCardOwner(winningCard)

		if err != nil {
			return fmt.Errorf("could not find white card owner: %w", err)
		}

		winner.IncrementScore()
		g.SetRoundStatus(valueobjects.JudgeChoseWinningCard)
		g.SetRoundWinner(winner)

	case EventCardPlayed:
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
			g.SetRoundStatus(valueobjects.JudgePickingWinningCard)
		}
	case EventGameWinner:
		var payload GameEventPayloadGameWinner

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventGameWinner payload: %w", err)
		}

		// Find the winning player and mark them as game winner
		for _, player := range g.Players {
			if player.UserID == payload.PlayerID {
				player.SetIsGameWinner(true)
				break
			}
		}

		g.SetRoundStatus(valueobjects.GameOver)
		g.SetStatus(valueobjects.Finished)

	case EventClockUpdate:
		var payload GameEventPayloadClockUpdate

		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal EventClockUpdate payload: %w", err)
		}

		g.NextAutoProgressAt = payload.NextAutoProgressAt
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	return nil
}
