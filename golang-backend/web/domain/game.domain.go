package domain

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

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
	ChatMessage     InboundWebsocketGameType = "CHAT_MESSAGE"
	BeginGame       InboundWebsocketGameType = "BEGIN_GAME"
)

type OutboundWebsocketGameType string

const (
	GameUpdate OutboundWebsocketGameType = "GAME_UPDATE"
)

type RoundStatus string

const (
	Waiting                    RoundStatus = "Waiting"
	PlayersPickingCard         RoundStatus = "PlayersPickingCard"
	CardCzarPickingWinningCard RoundStatus = "CardCzarPickingWinningCard"
	CardCzarChoseWinningCard   RoundStatus = "CardCzarChoseWinningCard"
)

type GameStatus string

const (
	Setup      GameStatus = "Setup"
	InProgress GameStatus = "InProgress"
	Finished   GameStatus = "Finished"
)

type Game struct {
	Mutex            sync.RWMutex
	ID               string
	Name             string
	Collection       *Collection
	WinnerCount      int
	MaxPlayerCount   int
	Status           GameStatus
	Players          []*Player
	WhiteCards       []*Card
	BlackCard        *Card
	RoundStatus      RoundStatus
	CurrentGameRound int
	RoundWinner      *Player
	LastVacatedAt    *time.Time
	Vacated          bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

func (g *Game) Lock() {
	g.Mutex.RLock()
}

func (g *Game) Unlock() {
	g.Mutex.RUnlock()
}

func (g *Game) SetStatus(status GameStatus) {
	g.Status = status
}

func (g *Game) SetRoundStatus(status RoundStatus) {
	g.RoundStatus = status
}

func (g *Game) SetRoundWinner(player *Player) {
	g.RoundWinner = player
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
	g.WhiteCards = []*Card{}
	g.BlackCard = nil
}

func (g *Game) PickNewBlackCard() error {
	if g.Collection == nil {
		return fmt.Errorf("could not pick new black card because collection is nil")
	}

	g.BlackCard = g.Collection.DrawCards(1, Black)[0]

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

func (g *Game) FindCardByPlayerId(playerId string, cardId string) (*Card, error) {
	player, err := g.FindPlayerByUserId(playerId)

	if err != nil {
		log.Println("could not find card by player id: %w", err)
		return nil, err
	}

	for _, card := range player.Deck {
		if card.ID == cardId {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card not found in player's deck")
}

func (g *Game) FindWhiteCardByCardId(cardId string) (*Card, error) {
	for _, card := range g.WhiteCards {
		if card.ID == cardId {
			return card, nil
		}
	}

	return nil, fmt.Errorf("card not found in cards")
}

func (g *Game) AddWhiteCardToGameBoard(card *Card) error {
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

func (g *Game) FindWhiteCardOwner(card *Card) (*Player, error) {
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
