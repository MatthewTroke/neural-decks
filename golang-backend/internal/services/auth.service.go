package services

import (
	"cardgame/internal/domain/entities"
	"cardgame/internal/domain/repositories"
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

type AuthService struct {
	userRepository repositories.UserRepository
	jwtService     *JWTAuthService
}

func NewAuthService(userRepository repositories.UserRepository, jwtService *JWTAuthService) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s *AuthService) AuthenticateUser(userData *entities.User, secret string) (*entities.AuthResult, error) {
	existingUser, err := s.userRepository.GetUserByEmail(userData.Email)
	if err == nil && existingUser != nil {
		if existingUser.Provider != userData.Provider {
			return nil, fmt.Errorf("account with email %s already exists via %s", userData.Email, existingUser.Provider)
		}
	}

	result, err := s.userRepository.UpsertUserByID(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}

	accessToken, err := s.jwtService.CreateAccessToken(
		result.Name,
		result.Email,
		result.ID,
		result.Image,
		result.EmailVerified,
		secret,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := s.jwtService.CreateRefreshToken(result.ID, secret)

	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &entities.AuthResult{
		User:         result,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RedirectURL:  "http://localhost:5173/games",
	}, nil
}

func (s *AuthService) ExchangeOAuthCode(code string, config oauth2.Config) (*oauth2.Token, error) {
	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}

func (s *AuthService) LogoutUser(refreshToken string) error {
	if refreshToken != "" {
		if err := s.jwtService.InvalidateRefreshToken(refreshToken); err != nil {
			return fmt.Errorf("failed to invalidate refresh token: %w", err)
		}
	}
	return nil
}

func (s *AuthService) InvalidateAllUserTokens(userID string) error {
	if err := s.jwtService.InvalidateAllRefreshTokensForUser(userID); err != nil {
		return fmt.Errorf("failed to invalidate all tokens: %w", err)
	}
	return nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	// TODO: Implement refresh token logic
	return "", fmt.Errorf("not implemented")
}
