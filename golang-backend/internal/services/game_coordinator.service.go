package services

import (
	"cardgame/internal/domain/aggregates"
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/events"
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
	games               []*aggregates.Game
}

func NewGameCoordinator(
	gameRepository repositories.GameRepository,
	eventRepository repositories.EventRepository,
	deckCreationService services.DeckCreationService,
	publisher Publisher,
	games []*aggregates.Game,
) *GameCoordinator {
	return &GameCoordinator{
		gameRepository:      gameRepository,
		eventRepository:     eventRepository,
		deckCreationService: deckCreationService,
		publisher:           publisher,
		games:               games,
	}
}

func (gc *GameCoordinator) getGameByID(gameId string) *aggregates.Game {
	for _, game := range gc.games {
		if game.ID == gameId {
			return game
		}
	}

	return nil
}

func (gc *GameCoordinator) Create(name string, deckSubject string, winnerCount int, maxPlayerCount int, claim *entities.CustomClaim) (*aggregates.Game, error) {
	collection, err := gc.deckCreationService.GenerateDeck(deckSubject)

	if err != nil {
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	player, err := aggregates.NewPlayer(claim)

	player.SetIsOwner(true)

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

	gc.games = append(gc.games, game)

	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

func (gc *GameCoordinator) Join(gameId string, claim *entities.CustomClaim) (*aggregates.Game, error) {
	// validationResult := gc.gameValidator.ValidateJoinGame(game, gameId, claim.UserID)

	// if !validationResult.IsValid {
	// 	return nil, fmt.Errorf("failed to validate join game: %w", validationResult.Errors)
	// }

	game := gc.getGameByID(gameId)

	if game == nil {
		return nil, fmt.Errorf("failed to get game")
	}

	player, err := game.FindPlayerByUserId(claim.UserID)

	if err == nil && player != nil {
		gc.publisher.PublishToRoom(game.ID, string(aggregates.GameUpdate), game)
		return game, nil
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

	err = game.ApplyEvent(event)

	if err != nil {
		return nil, fmt.Errorf("failed to apply joined game event: %w", err)
	}

	gc.publisher.PublishToRoom(game.ID, string(aggregates.GameUpdate), game)

	err = gc.eventRepository.AppendEvent(event)

	if err != nil {
		return nil, fmt.Errorf("failed to append joined game event: %w", err)
	}

	_, err = gc.gameRepository.Update(game)

	if err != nil {
		return nil, fmt.Errorf("failed to persist game update: %w", err)
	}

	return game, nil
}

func (gc *GameCoordinator) BeginGame(gameId string, claim *entities.CustomClaim) error {
	g := gc.getGameByID(gameId)

	if g == nil {
		return fmt.Errorf("failed to get game")
	}

	if !g.IsInSetup() {
		return fmt.Errorf("game %s is not in setup status", gameId)
	}

	player, err := g.FindPlayerByUserId(claim.UserID)

	if err != nil {
		return fmt.Errorf("failed to find player: %w", err)
	}

	if !player.IsOwner {
		return fmt.Errorf("player %s is not the owner of the game", player.ID)
	}

	payload, err := json.Marshal(events.NewGameEventPayloadGameBegins(gameId, player.ID))

	if err != nil {
		return fmt.Errorf("failed to marshal game begins payload: %w", err)
	}

	event := aggregates.NewGameEvent(
		gameId,
		aggregates.EventGameBegins,
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

func (gc *GameCoordinator) ContinueRound(gameId string, claim *entities.CustomClaim) error {
	g := gc.getGameByID(gameId)

	if g == nil {
		return fmt.Errorf("failed to get game")
	}

	if !g.IsInProgress() {
		return fmt.Errorf("game %s is not in progress status", gameId)
	}

	if g.ShouldShuffle() {
		shufflePayload, err := json.Marshal(aggregates.NewGameEventPayloadShuffle(gameId, time.Now().UnixNano(), uuid.New().String()))

		if err != nil {
			return fmt.Errorf("failed to marshal shuffle payload: %w", err)
		}

		shuffleEvent := aggregates.NewGameEvent(
			gameId,
			aggregates.EventShuffle,
			shufflePayload,
		)

		err = g.ApplyEvent(shuffleEvent)

		if err != nil {
			return fmt.Errorf("failed to create shuffle event: %w", err)
		}
	}

	playerCards := make(map[string]string)

	unusedWhiteCards := g.GetUnplayedWhiteCards()
	unusedBlackCards := g.GetUnplayedBlackCards()

	// Give each player a white card
	for _, player := range g.Players {
		if player.IsJudge {
			continue
		}

		playerCards[player.UserID] = unusedWhiteCards[0].ID
		unusedWhiteCards = unusedWhiteCards[1:]
	}

	payload, err := json.Marshal(aggregates.NewGameEventPayloadGameRoundContinuedWithCards(gameId, claim.UserID, playerCards, unusedBlackCards[0].ID))

	if err != nil {
		return fmt.Errorf("failed to marshal game round continued with cards payload: %w", err)
	}

	event := aggregates.NewGameEvent(
		gameId,
		aggregates.EventRoundContinued,
		payload,
	)

	err = g.ApplyEvent(event)

	if err != nil {
		return fmt.Errorf("failed to apply game round continued event: %w", err)
	}

	err = gc.eventRepository.AppendEvent(event)

	if err != nil {
		return fmt.Errorf("failed to append game round continued event: %w", err)
	}

	_, err = gc.gameRepository.Update(g)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	gc.publisher.PublishToRoom(gameId, string(aggregates.GameUpdate), g)

	return nil
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
