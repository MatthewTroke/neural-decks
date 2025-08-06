package repositories

import (
	"cardgame/domain/aggregates"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisGameRepository struct {
	client *redis.Client
}

func NewRedisGameRepository(client *redis.Client) *RedisGameRepository {
	return &RedisGameRepository{
		client: client,
	}
}

func (r *RedisGameRepository) Create(game *aggregates.Game) (*aggregates.Game, error) {
	ctx := context.Background()

	// Set creation time if not provided
	if game.CreatedAt.IsZero() {
		game.CreatedAt = time.Now()
	}

	// Set updated time
	game.UpdatedAt = time.Now()

	// Serialize the game
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game: %w", err)
	}

	gameKey := fmt.Sprintf("game:%s", game.ID)
	err = r.client.Set(ctx, gameKey, gameJSON, 0).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store game in Redis: %w", err)
	}

	return game, nil
}

func (r *RedisGameRepository) GetByID(id string) (*aggregates.Game, error) {
	ctx := context.Background()

	gameKey := fmt.Sprintf("game:%s", id)
	gameJSON, err := r.client.Get(ctx, gameKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get game from Redis: %w", err)
	}

	var game aggregates.Game
	if err := json.Unmarshal([]byte(gameJSON), &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game: %w", err)
	}

	return &game, nil
}

func (r *RedisGameRepository) Update(game *aggregates.Game) (*aggregates.Game, error) {
	ctx := context.Background()

	// Set updated time
	game.UpdatedAt = time.Now()

	// Serialize the game
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game: %w", err)
	}

	// Store updated game in Redis
	gameKey := fmt.Sprintf("game:%s", game.ID)
	err = r.client.Set(ctx, gameKey, gameJSON, 0).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to update game in Redis: %w", err)
	}

	return game, nil
}

func (r *RedisGameRepository) Delete(id string) error {
	ctx := context.Background()

	// Delete the game
	gameKey := fmt.Sprintf("game:%s", id)
	err := r.client.Del(ctx, gameKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete game from Redis: %w", err)
	}

	return nil
}

func (r *RedisGameRepository) GetAllGames() ([]*aggregates.Game, error) {
	ctx := context.Background()

	var games []*aggregates.Game
	var cursor uint64
	var err error

	for {
		var keys []string
		keys, cursor, err = r.client.Scan(ctx, cursor, "game:*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan Redis keys: %w", err)
		}

		for _, key := range keys {
			gameJSON, err := r.client.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue // Skip if key doesn't exist
				}
				continue // Skip games that can't be loaded
			}

			var game aggregates.Game
			if err := json.Unmarshal([]byte(gameJSON), &game); err != nil {
				continue // Skip games that can't be unmarshaled
			}

			games = append(games, &game)
		}

		// If cursor is 0, we've scanned all keys
		if cursor == 0 {
			break
		}
	}

	return games, nil
}
