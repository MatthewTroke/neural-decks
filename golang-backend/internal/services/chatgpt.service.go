package services

import (
	"cardgame/internal/domain/aggregates"
	"cardgame/internal/domain/entities"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

type ChatGPTDeckResponse struct {
	DeckName string `json:"deck_name"`
	Cards    []struct {
		Type  entities.CardType `json:"type"`
		Value string            `json:"value"`
	} `json:"cards"`
}

type ChatGPTService struct {
	client *openai.Client
}

var (
	instance *ChatGPTService
	once     sync.Once
)

func NewChatGPTService(chatgptApiKey string) *ChatGPTService {
	once.Do(func() {
		instance = &ChatGPTService{
			client: openai.NewClient(chatgptApiKey),
		}
	})

	return instance
}

func (s *ChatGPTService) sanitizeSubject(subject string) (string, error) {
	subject = strings.TrimSpace(subject)

	if len(subject) > 100 {
		return "", fmt.Errorf("subject too long (max 100 characters)")
	}

	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,!?()]+$`)

	if !validPattern.MatchString(subject) {
		return "", fmt.Errorf("subject contains invalid characters")
	}

	subject = strings.ReplaceAll(subject, "\n", " ")
	subject = strings.ReplaceAll(subject, "\r", " ")

	subject = regexp.MustCompile(`\s+`).ReplaceAllString(subject, " ")

	return subject, nil
}

func (s *ChatGPTService) GenerateDeck(subject string) (*aggregates.Collection, error) {
	sanitizedSubject, err := s.sanitizeSubject(subject)

	if err != nil {
		return nil, fmt.Errorf("invalid subject: %w", err)
	}

	functionDefinition := openai.FunctionDefinition{
		Name:        "generate_cards_against_humanity_deck",
		Description: "Generate a Cards Against Humanity deck based on a subject",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"deck_name": map[string]interface{}{
					"type":        "string",
					"description": "The name of the deck based on the subject",
				},
				"cards": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type":        "string",
								"description": "The type of card, either 'Black' or 'White'",
							},
							"value": map[string]interface{}{
								"type":        "string",
								"description": "The content of the card",
							},
						},
						"required": []string{"type", "value"},
					},
				},
			},
			"required": []string{"deck_name", "cards"},
		},
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role: "system",
			Content: `You are a creative assistant that generates Cards Against Humanity decks.
IMPORTANT: You MUST generate BOTH Black cards and White cards.

Requirements:
- Generate at least 30 White cards and 10 Black cards
- Black cards are the prompt cards that contain blanks (like "_____ is the key to success")
- White cards are the answer cards that fill in the blanks
- Make the cards funny and creative, with most White cards being 1-5 words
- The deck name should be based on the given subject
- Each card must have:
  - "type": which can only be either "Black" or "White",
  - "value": a string containing the card's content

Example Black cards:
- "_____ is the key to success"
- "The best thing about _____ is _____"
- "I never thought I'd see _____ in my lifetime"

Example White cards:
- "A microwave"
- "My ex"
- "The internet"`,
		},
		{
			Role:    "user",
			Content: "Generate a new deck of cards using the following subject: " + sanitizedSubject,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4oMini,
		Messages:  messages,
		Functions: []openai.FunctionDefinition{functionDefinition},
	}

	resp, err := s.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error generating deck: %w", err)
	}

	// Extract function call arguments
	if len(resp.Choices) == 0 || resp.Choices[0].Message.FunctionCall == nil {
		return nil, fmt.Errorf("no function call in response")
	}

	functionCall := resp.Choices[0].Message.FunctionCall
	var response ChatGPTDeckResponse

	err = json.Unmarshal([]byte(functionCall.Arguments), &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse function call arguments: %w", err)
	}

	collection := aggregates.NewCollection()

	for _, card := range response.Cards {
		collection.AddCard(
			entities.NewCard(uuid.New().String(), card.Type, card.Value),
		)
	}

	return collection, nil
}
