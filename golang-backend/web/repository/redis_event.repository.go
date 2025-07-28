package repository

import (
	"cardgame/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisEventRepository struct {
	client *redis.Client
}

func NewRedisEventRepository(client *redis.Client) *RedisEventRepository {
	return &RedisEventRepository{
		client: client,
	}
}

// AppendEvent adds a new event to the game's event stream
func (r *RedisEventRepository) AppendEvent(event domain.GameEvent) error {
	ctx := context.Background()

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	// Set creation time if not provided
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	// Serialize the event
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Store event in Redis list (game:events:{gameID})
	eventKey := fmt.Sprintf("game:events:%s", event.GameID)
	err = r.client.RPush(ctx, eventKey, eventJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to append event to Redis: %w", err)
	}

	// Also store event by ID for quick lookup
	eventIDKey := fmt.Sprintf("event:%s", event.ID)
	err = r.client.Set(ctx, eventIDKey, eventJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to store event by ID: %w", err)
	}

	// Update the game's last event timestamp
	lastEventKey := fmt.Sprintf("game:last_event:%s", event.GameID)
	err = r.client.Set(ctx, lastEventKey, event.CreatedAt.Unix(), 0).Err()
	if err != nil {
		log.Printf("Warning: failed to update last event timestamp: %v", err)
	}

	return nil
}

func (r *RedisEventRepository) GetEventsForGame(gameID string) ([]domain.GameEvent, error) {
	ctx := context.Background()

	eventKey := fmt.Sprintf("game:events:%s", gameID)
	eventJSONs, err := r.client.LRange(ctx, eventKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get events from Redis: %w", err)
	}

	var events []domain.GameEvent
	for _, eventJSON := range eventJSONs {
		var event domain.GameEvent
		if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
			log.Printf("Warning: failed to unmarshal event: %v", err)
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventsSince retrieves events since a specific timestamp
func (r *RedisEventRepository) GetEventsSince(gameID string, since time.Time) ([]domain.GameEvent, error) {
	allEvents, err := r.GetEventsForGame(gameID)
	if err != nil {
		return nil, err
	}

	var filteredEvents []domain.GameEvent
	for _, event := range allEvents {
		if event.CreatedAt.After(since) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

// GetEventByID retrieves a specific event by its ID
func (r *RedisEventRepository) GetEventByID(eventID string) (*domain.GameEvent, error) {
	ctx := context.Background()

	eventIDKey := fmt.Sprintf("event:%s", eventID)
	eventJSON, err := r.client.Get(ctx, eventIDKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("event not found: %s", eventID)
		}
		return nil, fmt.Errorf("failed to get event from Redis: %w", err)
	}

	var event domain.GameEvent
	if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return &event, nil
}

// DeleteGameEvents deletes all events for a game (useful for cleanup)
func (r *RedisEventRepository) DeleteGameEvents(gameID string) error {
	ctx := context.Background()

	// Get all events first to delete them by ID
	events, err := r.GetEventsForGame(gameID)
	if err != nil {
		return err
	}

	// Delete events by ID
	for _, event := range events {
		eventIDKey := fmt.Sprintf("event:%s", event.ID)
		r.client.Del(ctx, eventIDKey)
	}

	// Delete the event list
	eventKey := fmt.Sprintf("game:events:%s", gameID)
	r.client.Del(ctx, eventKey)

	// Delete last event timestamp
	lastEventKey := fmt.Sprintf("game:last_event:%s", gameID)
	r.client.Del(ctx, lastEventKey)

	return nil
}
