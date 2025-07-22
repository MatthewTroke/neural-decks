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
	Env         *bootstrap.Env
	GameService domain.GameService
	RoomManager *domain.RoomManager
	// Potentially interface with this incase you wanna swap to another gpt interface.
	ChatGPTService *services.ChatGPTService
}

func (gc *GameController) HandleGameRoomWebsocketInboundMessage(msg []byte, hub *domain.Hub, claim *domain.CustomClaim) error {
	var message domain.WebSocketMessage[domain.InboundWebsocketGameType, json.RawMessage]

	if err := json.Unmarshal(msg, &message); err != nil {
		return errors.New("unable to unmarshal WebSocket message")
	}

	var wsHandler handler.WebSocketHandler

	switch message.Type {
	case domain.BeginGame:
		var payload handler.BeginGamePayload

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal BeginGame payload")
		}

		wsHandler = handler.NewBeginGameHandler(
			payload,
			gc.GameService,
			claim,
			hub,
		)
	case domain.JoinGame:
		var payload handler.JoinGamePayload

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal JoinGamePayload")
		}

		wsHandler = handler.NewJoinGameHandler(
			payload,
			gc.GameService,
			claim,
			hub,
		)
	case domain.PlayCard:
		var payload handler.PlayCardPayload

		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal PlayCardPayload")
		}

		wsHandler = handler.NewPlayCardHandler(
			payload,
			gc.GameService,
			claim,
			hub,
		)
	case domain.PickWinningCard:
		var payload handler.PickWinningCardPayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal PickWinningCardPayload")
		}
		wsHandler = handler.NewPickWinningCardHandler(
			payload,
			gc.GameService,
			claim,
			hub,
		)
	case domain.ContinueRound:
		var payload handler.ContinueRoundPayload
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			return errors.New("unable to unmarshal ContinueRoundPayload")
		}

		wsHandler = handler.NewContinueRoundHandler(
			payload,
			gc.GameService,
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

	response, err := gc.ChatGPTService.GenerateDeck(request.Subject)

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
		WinnerCount:      5,
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

	return nil
}

func (gc *GameController) HandleGameRoomWebsocket(c *websocket.Conn) {
	roomId := c.Params("id")
	claim := c.Locals("user").(*domain.CustomClaim)

	hub := gc.RoomManager.GetRoom(roomId)

	hub.RegisterClient(c)

	game, _ := gc.GameService.GetGameById(roomId)

	websocketMessage := domain.NewWebSocketMessage(domain.GameUpdate, game)

	gameState, _ := json.Marshal(websocketMessage)

	c.WriteMessage(websocket.TextMessage, gameState)

	defer func() {
		hub.UnregisterClient(c)

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
