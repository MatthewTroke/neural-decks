package services

import (
	"cardgame/domain/aggregates"
	"cardgame/domain/entities"
	"cardgame/domain/events"
	"cardgame/domain/repositories"
	"cardgame/domain/valueobjects"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type GameCoordinator struct {
	gameRepository      repositories.GameRepository
	eventRepository     repositories.EventRepository
	deckCreationService DeckCreationService
}

func NewGameCoordinator(
	gameRepository repositories.GameRepository,
	eventRepository repositories.EventRepository,
	deckCreationService DeckCreationService,
) *GameCoordinator {
	return &GameCoordinator{
		gameRepository:      gameRepository,
		eventRepository:     eventRepository,
		deckCreationService: deckCreationService,
	}
}

func (gc *GameCoordinator) Create(name string, deckSubject string, winnerCount int, maxPlayerCount int, userID string) (*aggregates.Game, error) {
	collection, err := gc.deckCreationService.GenerateDeck(deckSubject)

	if err != nil {
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	collection.Shuffle()

	game := &aggregates.Game{
		ID:               uuid.New().String(),
		Name:             name,
		Collection:       collection,
		WinnerCount:      winnerCount,
		MaxPlayerCount:   maxPlayerCount,
		Status:           valueobjects.Setup,
		Players:          []*aggregates.Player{},
		WhiteCards:       []*entities.Card{},
		BlackCard:        nil,
		RoundStatus:      valueobjects.Waiting,
		CurrentGameRound: 0,
		LastVacatedAt:    nil,
		Vacated:          false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	gc.gameRepository.Create(game)

	payload, err := json.Marshal(aggregates.NewGameEventPayloadGameBegins(game.ID, userID))

	if err != nil {
		return nil, fmt.Errorf("failed to marshal game begins payload: %w", err)
	}

	err = gc.eventRepository.AppendEvent(
		events.NewGameEvent(
			game.ID,
			events.EventGameBegins,
			payload,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to append game begins event: %w", err)
	}

	gc.BroadcastGameUpdate(game)

	return game, nil
}

func (gc *GameCoordinator) Join(gameId string, claim *entities.CustomClaim, wsConnection *websocket.Conn) (*aggregates.Game, error) {
	game, err := gc.gameRepository.GetByID(gameId)

	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	player, err := aggregates.NewPlayer(claim, wsConnection)

	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	game.AddPlayer(player)

	_, err = gc.gameRepository.Update(game)

	if err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return game, nil
}

func (gc *GameCoordinator) BeginGame(gameId string, claim *entities.CustomClaim) error {
	g, err := gc.gameRepository.GetByID(gameId)

	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

	if !g.IsInSetup() {
		return fmt.Errorf("game %s is not in setup status", gameId)
	}

	payload, err := json.Marshal(aggregates.NewGameEventPayloadGameBegins(gameId, claim.UserID))

	if err != nil {
		return fmt.Errorf("failed to marshal game begins payload: %w", err)
	}

	event := events.NewGameEvent(
		gameId,
		events.EventGameBegins,
		payload,
	)

	err = g.ApplyEvent(event)

	if err != nil {
		return fmt.Errorf("failed to apply game begins event: %w", err)
	}

	err = gc.eventRepository.AppendEvent(event)

	if err != nil {
		return fmt.Errorf("failed to append game begins event: %w", err)
	}

	return nil
}

func (gc *GameCoordinator) Leave(gameId string, claim *entities.CustomClaim) error {
	game, err := gc.gameRepository.GetByID(gameId)

	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

	player, err := game.FindPlayerByUserId(claim.UserID)

	if err != nil {
		return fmt.Errorf("failed to find player: %w", err)
	}

	err = game.RemovePlayer(player)

	if err != nil {
		return fmt.Errorf("failed to remove player: %w", err)
	}

	_, err = gc.gameRepository.Update(game)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}

func (gc *GameCoordinator) BroadcastToGame(game *aggregates.Game, messageType string, payload interface{}) error {
	return game.BroadcastToAllPlayers(messageType, payload)
}

func (gc *GameCoordinator) BroadcastToGameExcept(game *aggregates.Game, messageType string, payload interface{}, exceptUserID string) error {
	return game.BroadcastToAllPlayersExcept(messageType, payload, exceptUserID)
}

func (gc *GameCoordinator) BroadcastGameUpdate(game *aggregates.Game) error {
	return gc.BroadcastToGame(game, string(aggregates.GameUpdate), game)
}

func (gc *GameCoordinator) BroadcastChatMessage(game *aggregates.Game, message string) error {
	return gc.BroadcastToGame(game, string(aggregates.ChatMessage), message)
}

func (gc *GameCoordinator) BroadcastEmojiClicked(game *aggregates.Game, emojiData interface{}) error {
	return gc.BroadcastToGame(game, string(aggregates.EmojiClickedMessage), emojiData)
}
