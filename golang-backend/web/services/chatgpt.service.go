package services

import (
	"cardgame/bootstrap"
	"cardgame/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGPTDeckResponse struct {
	DeckName string `json:"deck_name"`
	Cards    []struct {
		Type  domain.CardType `json:"type"`
		Value string          `json:"value"`
	} `json:"cards"`
}

type ChatGPTService struct {
	client *openai.Client
}

var (
	instance *ChatGPTService
	once     sync.Once
)

func NewChatGPTService(env *bootstrap.Env) *ChatGPTService {
	once.Do(func() {
		apiKey := env.ChatGPTAPIKey
		if apiKey == "" {
			log.Fatal("CHATGPT_API_KEY environment variable not set")
		}
		instance = &ChatGPTService{
			client: openai.NewClient(apiKey),
		}
	})
	return instance
}

func (s *ChatGPTService) GenerateDeck(subject string) (ChatGPTDeckResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: "system",
			Content: `You are a creative assistant that generates new Cards Against Humanity decks.
Your output must strictly follow this JSON format without any additional text:
{
  "deck_name": "string",
  "cards": [
    { "type": "string", "value": "string" }
  ]
}
- You must supply at least 30 White cards and 10 Black cards
- Try to make the cards as funny and creative as possible, but not too specific to the subject itself.
- Have some variety in card content lengths, with most of the White cards being between 1 and 5 words.
- The "deck_name" should be a string and should be the name of the given subject.
- Each card object must have:
  - "type": which can only be either "Black" or "White",
  - "value": a string containing the card's content.`,
		},
		{
			Role:    "user",
			Content: "Generate a new deck of cards using the following subject: " + subject,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
	}

	fmt.Printf("Beginning deck creation...")

	resp, err := s.client.CreateChatCompletion(context.Background(), req)

	fmt.Printf("Ending deck creation...")

	if err != nil {
		fmt.Errorf("Error generating deck: %v", err)
		return ChatGPTDeckResponse{}, err
	}

	var response ChatGPTDeckResponse

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &response)

	if err != nil {
		return ChatGPTDeckResponse{}, err
	}

	return response, nil
}
