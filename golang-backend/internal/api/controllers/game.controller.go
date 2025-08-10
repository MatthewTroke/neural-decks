package controllers

import (
	"cardgame/internal/api/request"
	"cardgame/internal/api/validation"
	"cardgame/internal/domain/aggregates"
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/repositories"
	"cardgame/internal/infra/environment"
	"cardgame/internal/infra/ws"
	"cardgame/internal/services"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type GameController struct {
	Env             *environment.Env
	gameCoordinator *services.GameCoordinator
	gameRepository  repositories.GameRepository
	hub             *ws.Hub
}

func NewGameController(
	env *environment.Env,
	gameCoordinator *services.GameCoordinator,
	gameRepository repositories.GameRepository,
	hub *ws.Hub,
) *GameController {
	return &GameController{
		Env:             env,
		gameCoordinator: gameCoordinator,
		gameRepository:  gameRepository,
		hub:             hub,
	}
}

func (gc *GameController) HandleJoinWebsocketGameRoom(connection *websocket.Conn) {
	gameId := connection.Params("id")
	claim := connection.Locals("user").(*entities.CustomClaim)

	client := &ws.Client{
		Conn:   connection,
		RoomID: gameId,
		Send:   make(chan []byte, 64),
	}

	gc.hub.Join(gameId, client)
	gc.gameCoordinator.Join(gameId, claim)

	defer func() {
		gc.hub.Leave(client)
		gc.gameCoordinator.Leave(gameId, claim)
	}()

	// Incoming message pump
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}

		var message request.GameEventRequest
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error unmarshaling WebSocket message: %v", err)
			continue // Continue processing other messages
		}

		switch message.Type {
		case aggregates.EventGameBegins:
			var payload request.GameEventPayloadGameBeginsRequest
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling GameBegins payload: %v", err)
				continue
			}
			gc.gameCoordinator.BeginGame(payload.GameID, claim)

		case aggregates.EventCardPlayed:
			var payload request.GameEventPayloadJoinedGameRequest
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling CardPlayed payload: %v", err)
				continue
			}
			// gc.gameCoordinator.PlayCard(payload.GameID, claim, payload.CardID)

		case aggregates.EventJudgeChoseWinningCard:
			var payload request.GameEventPayloadJudgeChoseWinningCardRequest
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling JudgeChoseWinningCard payload: %v", err)
				continue
			}
			// gc.gameCoordinator.PickWinningCard(payload.GameID, claim, payload.CardID)

		case aggregates.EventRoundContinued:
			var payload request.GameEventPayloadGameRoundContinuedRequest
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling RoundContinued payload: %v", err)
				continue
			}
			// gc.gameCoordinator.ContinueRound(payload.GameID, claim)

		case aggregates.EventEmojiClicked:
			var payload request.GameEventPayloadEmojiClickedRequest
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				log.Printf("Error unmarshaling EmojiClicked payload: %v", err)
				continue
			}
			// gc.gameCoordinator.ClickEmoji(payload.GameID, claim, payload.Emoji)

		default:
			log.Printf("Unknown WebSocket message type: %s", message.Type)
			continue
		}
	}
}

func (gc *GameController) HandleGetGames(c *fiber.Ctx) error {
	games, err := gc.gameRepository.GetAllGames()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(games)
}

func (gc *GameController) CreateGame(c *fiber.Ctx) error {
	var request request.CreateGameRequest

	if err := validation.BindAndValidate(c, &request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	claim := c.Locals("user").(*entities.CustomClaim)

	game, err := gc.gameCoordinator.Create(request.Name, request.Subject, request.WinnerCount, request.MaxPlayerCount, claim)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(game)
}
