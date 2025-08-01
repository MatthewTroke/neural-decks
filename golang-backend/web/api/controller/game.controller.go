package controller

import (
	"cardgame/api/handler"
	"cardgame/bootstrap"
	"cardgame/domain"
	"cardgame/request"
	"cardgame/services"
	"cardgame/utils"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type GameController struct {
	Env              *bootstrap.Env
	GameService      domain.GameService
	EventService     *services.EventService
	GameStateService *services.GameStateService
	RoomManager      *domain.RoomManager
	ChatGPTService   *services.ChatGPTService
}

// NewGameController creates a new GameController with initialized cache
func NewGameController(env *bootstrap.Env, gameService domain.GameService, eventService *services.EventService, gameStateService *services.GameStateService, roomManager *domain.RoomManager, chatGPTService *services.ChatGPTService) *GameController {
	return &GameController{
		Env:              env,
		GameService:      gameService,
		EventService:     eventService,
		GameStateService: gameStateService,
		RoomManager:      roomManager,
		ChatGPTService:   chatGPTService,
	}
}

func (gc *GameController) HandleGameRoomWebsocketInboundMessage(msg []byte, hub *domain.Hub, claim *domain.CustomClaim) error {
	var message request.GameEventRequest

	if err := json.Unmarshal(msg, &message); err != nil {
		return errors.New("unable to unmarshal WebSocket message into a GameEventRequest")
	}

	var wsHandler handler.WebSocketHandler

	switch message.Type {
	case domain.EventGameBegins:
		var payload request.GameEventPayloadGameBeginsRequest

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal BeginGame payload")
		}

		wsHandler = handler.NewBeginGameHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case domain.EventJoinedGame:
		var payload request.GameEventPayloadJoinedGameRequest

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal JoinGamePayload")
		}

		wsHandler = handler.NewJoinGameHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case domain.EventCardPlayed:
		var payload request.GameEventPayloadPlayCardRequest

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal PlayCardPayload")
		}

		wsHandler = handler.NewPlayCardHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case domain.EventCardCzarChoseWinningCard:
		var payload request.GameEventPayloadCardCzarChoseWinningCardRequest
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal PickWinningCardPayload")
		}
		wsHandler = handler.NewPickWinningCardHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case domain.EventRoundContinued:
		var payload request.GameEventPayloadGameRoundContinuedRequest
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal ContinueRoundPayload")
		}

		wsHandler = handler.NewContinueRoundHandler(
			payload,
			gc.EventService,
			gc.GameStateService,
			claim,
			hub,
		)
	case domain.EventEmojiClicked:
		var payload request.GameEventPayloadEmojiClickedRequest
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal EmojiClickedPayload")
		}
		wsHandler = handler.NewEmojiClickedHandler(
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

	response, err := gc.ChatGPTService.GenerateDeckWithFunctionCalling(request.Subject)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	collection := domain.NewCollection()

	for _, card := range response.Cards {
		cardId, _ := uuid.NewRandom()

		collection.AddCard(&domain.Card{
			ID:        cardId.String(),
			Type:      card.Type,
			CardValue: card.Value,
		})
	}

	gameId, _ := uuid.NewRandom()

	game := &domain.Game{
		ID:               gameId.String(),
		Name:             request.Name,
		Collection:       collection,
		WinnerCount:      request.WinnerCount,
		MaxPlayerCount:   request.MaxPlayerCount,
		Status:           "Setup",
		Players:          []*domain.Player{},
		WhiteCards:       []*domain.Card{},
		BlackCard:        nil,
		RoundStatus:      "Waiting",
		CurrentGameRound: 0,
		LastVacatedAt:    nil,
		Vacated:          false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	gc.GameService.AddGame(game)

	return c.Status(fiber.StatusCreated).JSON(game)
}

func (gc *GameController) HandleJoinWebsocketGameRoom(c *websocket.Conn) {
	roomId := c.Params("id")
	claim := c.Locals("user").(*domain.CustomClaim)

	hub := gc.RoomManager.GetRoom(roomId)

	hub.RegisterClient(c)

	game, _ := gc.EventService.GetGameById(roomId)

	websocketMessage := domain.NewWebSocketMessage(domain.GameUpdate, game)
	log := claim.Name + " has joined the room."

	websocketChatMessage := domain.NewWebSocketMessage(domain.ChatMessage, log)

	gameState, _ := json.Marshal(websocketMessage)
	chatMessage, _ := json.Marshal(websocketChatMessage)

	c.WriteMessage(websocket.TextMessage, gameState)
	hub.Broadcast(chatMessage)

	// Capture the user's name for the leave message
	userName := claim.Name

	defer func() {
		hub.UnregisterClient(c)

		// Broadcast leave message
		leaveLog := userName + " has left the room."
		websocketLeaveMessage := domain.NewWebSocketMessage(domain.ChatMessage, leaveLog)
		leaveMessage, _ := json.Marshal(websocketLeaveMessage)
		hub.Broadcast(leaveMessage)

		if len(hub.Clients) == 0 {
			gc.RoomManager.RemoveRoom(roomId)
		}
	}()

	for {
		_, msg, err := c.ReadMessage()

		if err != nil {
			break
		}

		err = gc.HandleGameRoomWebsocketInboundMessage(msg, hub, claim)

		if err != nil {
			fmt.Println("Error handling message:", err)
			continue
		}
	}
}
