package controllers

import (
	"cardgame/bootstrap/environment"
	"cardgame/domain/entities"
	"cardgame/domain/events"
	"cardgame/domain/services"
	"cardgame/http/handlers"
	"cardgame/http/request"
	"cardgame/infra/external/ai"
	"cardgame/infra/websockets"
	"cardgame/utils"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type GameController struct {
	Env              *environment.Env
	EventService     *services.EventService
	GameStateService *services.GameStateService
	RoomManager      *websockets.RoomManager
	//TODO GET RID OF THIS FROM HERE.
	ChatGPTService  *ai.ChatGPTService
	GameCoordinator *services.GameCoordinator
}

func NewGameController(
	env *environment.Env,
	// eventService *services.EventService,
	// gameStateService *services.GameStateService,
	roomManager *websockets.RoomManager,
	chatGPTService *ai.ChatGPTService,
	gameCoordinator *services.GameCoordinator,
) *GameController {
	return &GameController{
		Env: env,
		// EventService:     eventService,
		// GameStateService: gameStateService,
		RoomManager:    roomManager,
		ChatGPTService: chatGPTService,
	}
}

func (gc *GameController) HandleGameRoomWebsocketInboundMessage(msg []byte, claim *entities.CustomClaim) error {
	var message request.GameEventRequest

	if err := json.Unmarshal(msg, &message); err != nil {
		return errors.New("unable to unmarshal WebSocket message into a GameEventRequest")
	}

	var wsHandler handlers.WebSocketHandler

	switch message.Type {
	case events.EventGameBegins:
		var payload request.GameEventPayloadGameBeginsRequest

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal BeginGame payload")
		}

		err := gc.GameCoordinator.BeginGame(payload.GameID, claim)

		if err != nil {
			return fmt.Errorf("failed to begin game: %w", err)
		}

		return nil
	case events.EventJoinedGame:
		var payload request.GameEventPayloadJoinedGameRequest

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal JoinGamePayload")
		}

		wsHandler = handlers.NewJoinGameHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case events.EventCardPlayed:
		var payload request.GameEventPayloadPlayCardRequest

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal PlayCardPayload")
		}

		wsHandler = handlers.NewPlayCardHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case events.EventCardCzarChoseWinningCard:
		var payload request.GameEventPayloadCardCzarChoseWinningCardRequest
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal PickWinningCardPayload")
		}
		wsHandler = handlers.NewPickWinningCardHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case events.EventRoundContinued:
		var payload request.GameEventPayloadGameRoundContinuedRequest
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal ContinueRoundPayload")
		}

		wsHandler = handlers.NewContinueRoundHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case events.EventEmojiClicked:
		var payload request.GameEventPayloadEmojiClickedRequest
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal EmojiClickedPayload")
		}
		wsHandler = handlers.NewEmojiClickedHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	}

	if wsHandler == nil {
		return fmt.Errorf("unable to handle inbound message, no handler found")
	}

	hasError := wsHandler.Validate()

	if hasError != nil {
		return hasError
	}

	wsHandler.Handle()

	return nil
}

func (gc *GameController) CreateGame(c *fiber.Ctx) error {
	var request request.CreateGameRequest

	if err := utils.BindAndValidate(c, &request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	claim := c.Locals("user").(*entities.CustomClaim)
	userID := claim.UserID

	game, err := gc.GameCoordinator.Create(request.Name, request.Subject, request.WinnerCount, request.MaxPlayerCount, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(game)
}

func (gc *GameController) HandleJoinWebsocketGameRoom(connection *websocket.Conn) error {
	gameId := connection.Params("id")
	claim := connection.Locals("user").(*entities.CustomClaim)

	game, err := gc.GameCoordinator.Join(gameId, claim, connection)

	if err != nil {
		return fmt.Errorf("Could not join game %s: %v\n", gameId, err)
	}

	defer func() {
		err := gc.GameCoordinator.Leave(game.ID, claim)

		if err != nil {
			fmt.Printf("Warning: Could not remove websocket connection for player %s: %v\n", claim.UserID, err)
		}
	}()

	for {
		_, msg, err := connection.ReadMessage()

		if err != nil {
			break
		}

		err = gc.HandleGameRoomWebsocketInboundMessage(msg, claim)

		if err != nil {
			fmt.Println("Error handling message:", err)
			continue
		}
	}
}

// var role string

// const (
// 	PlayerRolePlayer = "PLAYER"
// 	PlayerRoleCardCzar = "SPECTATOR"
// )

// type Player struct {
// 	user *entities.User
// 	role aggregates.PlayerRole
// 	conn *websocket.Conn
// }
