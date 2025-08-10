package services

import (
	"cardgame/internal/domain/aggregates"
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/repositories"
	"cardgame/internal/domain/services"
	"cardgame/internal/domain/valueobjects"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Publisher interface {
	PublishToRoom(roomID string, eventType string, payload any) error
}

type GameCoordinator struct {
	gameRepository      repositories.GameRepository
	eventRepository     repositories.EventRepository
	deckCreationService services.DeckCreationService
	publisher           Publisher
}

func NewGameCoordinator(
	gameRepository repositories.GameRepository,
	eventRepository repositories.EventRepository,
	deckCreationService services.DeckCreationService,
	publisher Publisher,
) *GameCoordinator {
	return &GameCoordinator{
		gameRepository:      gameRepository,
		eventRepository:     eventRepository,
		deckCreationService: deckCreationService,
		publisher:           publisher,
	}
}

func (gc *GameCoordinator) Create(name string, deckSubject string, winnerCount int, maxPlayerCount int, claim *entities.CustomClaim) (*aggregates.Game, error) {
	collection, err := gc.deckCreationService.GenerateDeck(deckSubject)

	if err != nil {
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	player, err := aggregates.NewPlayer(claim)

	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	players := []*aggregates.Player{
		player,
	}

	gameID := uuid.New().String()

	game := aggregates.NewGame(
		gameID,
		name,
		collection,
		winnerCount,
		maxPlayerCount,
		valueobjects.Setup,
		players,
		[]*entities.Card{},
		nil,
		valueobjects.Waiting,
		0,
		nil,
		time.Now(),
		time.Now(),
		time.Now(),
		time.Now(),
		time.Now(),
		time.Now(),
	)

	game, err = gc.gameRepository.Create(game)

	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

func (gc *GameCoordinator) Join(gameId string, claim *entities.CustomClaim) (*aggregates.Game, error) {
	game, err := gc.gameRepository.GetByID(gameId)

	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	player, err := game.FindPlayerByUserId(claim.UserID)

	if err != nil {
		return nil, fmt.Errorf("failed to find player: %w", err)
	}

	if player != nil {
		gc.publisher.PublishToRoom(game.ID, string(aggregates.GameUpdate), game)
		return game, nil
	}

	player, err = aggregates.NewPlayer(claim)

	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	game.AddPlayer(player)

	_, err = gc.gameRepository.Update(game)

	if err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	eventPayload, err := json.Marshal(aggregates.NewGameEventPayloadJoinedGame(game.ID, claim.UserID, claim))

	if err != nil {
		return nil, fmt.Errorf("failed to marshal game begins payload: %w", err)
	}

	event := aggregates.NewGameEvent(
		game.ID,
		aggregates.EventJoinedGame,
		eventPayload,
	)

	err = gc.eventRepository.AppendEvent(event)

	if err != nil {
		return nil, fmt.Errorf("failed to append joined game event: %w", err)
	}

	gc.publisher.PublishToRoom(game.ID, string(aggregates.GameUpdate), game)

	return game, nil
}

func (gc *GameCoordinator) BeginGame(gameId string, claim *entities.CustomClaim) {
	// g, err := gc.gameRepository.GetByID(gameId)

	// if err != nil {
	// 	return fmt.Errorf("failed to get game: %w", err)
	// }

	// if !g.IsInSetup() {
	// 	return fmt.Errorf("game %s is not in setup status", gameId)
	// }

	// payload, err := json.Marshal(aggregates.NewGameEventPayloadGameBegins(gameId, claim.UserID))

	// if err != nil {
	// 	return fmt.Errorf("failed to marshal game begins payload: %w", err)
	// }

	// event := events.NewGameEvent(
	// 	gameId,
	// 	events.EventGameBegins,
	// 	payload,
	// )

	// err = g.ApplyEvent(event)

	// if err != nil {
	// 	return fmt.Errorf("failed to apply game begins event: %w", err)
	// }

	// err = gc.eventRepository.AppendEvent(event)

	// if err != nil {
	// 	return fmt.Errorf("failed to append game begins event: %w", err)
	// }

	// return nil
}

func (gc *GameCoordinator) Leave(gameId string, claim *entities.CustomClaim) {
	// Ensure the websocket client is removed from the room to stop broadcasts
	// gc.broadcaster.Leave(gameId, client)

	// Close the client to stop any pending operations
	// client.Close()

	// TODO: Implement game state updates when player leaves
	// game, err := gc.gameRepository.GetByID(gameId)
	// if err != nil {
	// 	return fmt.Errorf("failed to get game: %w", err)
	// }

	// player, err := game.FindPlayerByUserId(claim.UserID)
	// if err != nil {
	// 	return fmt.Errorf("failed to find player: %w", err)
	// }

	// err = game.RemovePlayer(player)
	// if err != nil {
	// 	return fmt.Errorf("failed to remove player: %w", err)
	// }

	// _, err = gc.gameRepository.Update(game)
	// if err != nil {
	// 	return fmt.Errorf("failed to update game: %w", err)
	// }
}
