# Neural Decks - Event-Driven Architecture Design Schema

## Overview

Neural Decks implements an **Event Sourcing** pattern with **CQRS (Command Query Responsibility Segregation)**. The architecture is built around immutable events that represent all state changes, enabling complete audit trails and game state reconstruction.

## Core Architecture Components

### 1. Event-Driven Flow Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   WebSocket      │    │   Controller    │    │   Handler       │
│   (React)       │───▶│   Connection     │───▶│   (Game)        │───▶│   (Specific)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └─────────────────┘
                                                                              │
                                                                              ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Redis         │◀───│   Event Service  │◀───│   Game State    │◀───│   Game Domain   │
│   (Event Store) │    │                  │    │   Service       │    │   (ApplyEvent)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └─────────────────┘
```

### 2. Detailed Component Flow

#### A. WebSocket Message Processing
```
1. Frontend sends WebSocket message
   ↓
2. GameController.HandleGameRoomWebsocketInboundMessage()
   ↓
3. Message type determines handler creation:
   - BEGIN_GAME → BeginGameHandler
   - JOIN_GAME → JoinGameHandler  
   - PLAY_CARD → PlayCardHandler
   - PICK_WINNING_CARD → PickWinningCardHandler
   - CONTINUE_ROUND → ContinueRoundHandler
   ↓
4. Handler.Handle() executes
```

#### B. Event Creation & Persistence
```
1. Handler calls EventService.CreateGameEvent()
   ↓
2. EventService creates GameEvent with:
   - Unique ID (UUID)
   - Game ID
   - Event Type
   - JSON Payload
   - Timestamp
   ↓
3. Handler calls EventService.AppendEvent()
   ↓
4. RedisEventRepository stores in Redis:
   - game:events:{gameID} (list)
   - event:{eventID} (hash)
   - game:last_event:{gameID} (timestamp)
```

#### C. Game State Reconstruction
```
1. EventService.BuildGameByGameId()
   ↓
2. GameStateService.GetGameById() (in-memory)
   ↓
3. EventService.rebuildGameFromEvents():
   - Fetch all events from Redis
   - Clone initial game state
   - Apply each event sequentially
   - Return current game state
```

## 3. Game State Service - The Orchestrator

The `GameStateService` is the central orchestrator that manages:

### A. In-Memory Game Management
```go
type GameStateService struct {
    mu           sync.RWMutex
    games        map[string]*domain.Game    // Active games in memory
    eventService *EventService              // For event operations
    roomManager  *domain.RoomManager        // WebSocket broadcasting
    timerResets  map[string]chan bool       // Auto-progress timers
}
```

**Note**: Currently, initial game state is stored in-memory only. The Redis event store takes this initial state and reconstructs the current game state by applying all events. In the future, this will be replaced with a more persistent storage solution (e.g., PostgreSQL) for the initial game state.

### B. Auto-Progress System
The service implements intelligent auto-progression:

```go
func (s *GameStateService) autoContinueGame(gameID string) {
    for {
        select {
        case <-time.After(30 * time.Second):
            // Auto-progress game state
            s.autoProgress(game)
        case resetSignal := <-resetChan:
            // Timer reset by manual action
        }
    }
}
```

**Auto-Progress Logic:**
- **PlayersPickingCard**: Auto-plays cards for inactive players
- **CardCzarPickingWinningCard**: Auto-selects winning card
- **CardCzarChoseWinningCard**: Auto-continues to next round

### C. Timer Management
- Each game has a dedicated timer channel
- Manual actions reset the 30-second countdown
- Timers automatically clean up when games end

## 4. Event Sourcing Implementation

### A. Event Types
```go
const (
    EventGameBegins               GameEventType = "GameBegins"
    EventJoinedGame               GameEventType = "JoinedGame"
    EventCardPlayed               GameEventType = "CardPlayed"
    EventRoundContinued           GameEventType = "RoundContinued"
    EventCardCzarChoseWinningCard GameEventType = "CardCzarChoseWinningCard"
    EventShuffle                  GameEventType = "Shuffle"
    EventDealCards                GameEventType = "DealCards"
    EventDrawBlackCard            GameEventType = "DrawBlackCard"
    EventSetCardCzar              GameEventType = "SetCardCzar"
    EventTimerUpdate              GameEventType = "TimerUpdate"
    EventGameWinner               GameEventType = "GameWinner"
    EventClockUpdate              GameEventType = "ClockUpdate"
    EventEmojiClicked             GameEventType = "EmojiClicked"
)
```

### B. Event Application
```go
func (g *Game) ApplyEvent(event GameEvent) error {
    g.Lock()
    defer g.Unlock()
    
    g.SetLastEventAt(event.CreatedAt)
    
    switch event.Type {
    case EventGameBegins:
        g.SetRoundStatus(PlayersPickingCard)
        g.SetStatus(InProgress)
    case EventCardPlayed:
        // Apply card play logic
    case EventRoundContinued:
        // Apply round continuation logic
    // ... other event types
    }
}
```

## 5. Redis Integration

### A. Event Storage Strategy
```go
// Primary storage: game:events:{gameID} (Redis List)
err = r.client.RPush(ctx, eventKey, eventJSON).Err()

// Quick lookup: event:{eventID} (Redis Hash)
err = r.client.Set(ctx, eventIDKey, eventJSON, 0).Err()

// Timestamp tracking: game:last_event:{gameID}
err = r.client.Set(ctx, lastEventKey, event.CreatedAt.Unix(), 0).Err()
```

### B. Used Cards Tracking
```go
// Track used cards to prevent reshuffling
func (r *RedisEventRepository) AddUsedCard(gameID, cardID string) error {
    key := fmt.Sprintf("game:used_cards:%s", gameID)
    return r.client.SAdd(ctx, key, cardID).Err()
}
```

## 6. Handler Pattern Implementation

### A. Handler Interface
```go
type WebSocketHandler interface {
    Handle() error
}
```

### B. Handler Flow Example (PlayCard)
```go
func (h *PlayCardHandler) Handle() error {
    // 1. Rebuild current game state from events
    game, err := h.EventService.BuildGameByGameId(h.Payload.GameID)
    
    // 2. Create new event
    event, err := h.EventService.CreateGameEvent(
        h.Payload.GameID,
        domain.EventCardPlayed,
        domain.NewGameEventPayloadPlayCard(h.Payload.GameID, h.Payload.CardID, h.Claim),
    )
    
    // 3. Apply event to game state
    newGame := game.Clone()
    newGame.ApplyEvent(event)
    
    // 4. Persist event to Redis
    h.EventService.AppendEvent(event)
    
    // 5. Broadcast updated state
    message := domain.NewWebSocketMessage(domain.GameUpdate, newGame)
    h.Hub.Broadcast(jsonMessage)
    
    return nil
}
```

## 7. Key Benefits of This Architecture

### A. Event Sourcing Advantages
- **Complete Audit Trail**: Every state change is recorded
- **Debugging**: Easy to replay events to reproduce issues
- **Scalability**: Events can be processed asynchronously

### C. Real-time Capabilities
- **WebSocket Broadcasting**: Instant state updates to all players
- **Auto-progression**: Intelligent game flow management
- **Timer Management**: Automatic cleanup and progression

## 8. Data Flow Summary

```
1. User Action (Frontend)
   ↓
2. WebSocket Message
   ↓
3. Controller Route Selection
   ↓
4. Handler Creation & Execution
   ↓
5. Event Creation & Validation
   ↓
6. Game State Reconstruction (from Redis events)
   ↓
7. Event Application (immutable state change)
   ↓
8. Event Persistence (Redis)
   ↓
9. State Broadcast (WebSocket)
   ↓
10. Frontend Update
```

## 9. Current Architectural Concerns

While the current architecture provides a solid foundation, there are several areas that could be improved for better maintainability and scalability:

### A. Service Coupling Issues
**Problem**: Services are tightly coupled through direct dependency injection
```go
// Current tight coupling example
type GameStateService struct {
    eventService *EventService  // Direct dependency
    roomManager  *domain.RoomManager
}

type EventService struct {
    gameStateService *GameStateService  // Circular dependency
}
```

**Issues**:
- **Circular Dependencies**: Services reference each other directly
- **Tight Coupling**: Changes in one service require changes in others
- **Testing Complexity**: Hard to unit test services in isolation
- **Violation of Single Responsibility**: Services know too much about each other

**Potential Solutions**:
- **Event Bus/Message Queue**: Decouple services through events
- **Interface Segregation**: Define clear contracts between services
- **Dependency Inversion**: Depend on abstractions, not concretions

### B. Game State Service Over-Responsibility
**Problem**: The `GameStateService` has become a "god object" handling too many concerns

**Current Responsibilities**: (THIS IS TOO MUCH IMO)
- In-memory game state management
- Timer management and auto-progression
- WebSocket broadcasting coordination
- Event service coordination
- Game cleanup and lifecycle management

**Issues**:
- **Single Responsibility Violation**: One service doing too many things
- **Hard to Test**: Complex interactions make unit testing difficult
- **Scalability Concerns**: Timer management doesn't scale across multiple instances
- **Maintenance Burden**: Changes affect multiple concerns simultaneously

### C. Initial State Management
**Current State**: Initial game state is stored in-memory only
```go
// Current approach
games map[string]*domain.Game  // In-memory only
```

**Future Direction**: Move to persistent storage for initial game state
```go
// Future approach
type GameRepository interface {
    Create(game *Game) error
    GetByID(id string) (*Game, error)
    Update(game *Game) error
    Delete(id string) error
}
```