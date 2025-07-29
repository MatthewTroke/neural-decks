# OAuth2 Session Persistence Guide

## Overview

This guide explains how to maintain OAuth2 login sessions for extended periods using refresh tokens and automatic token refresh mechanisms.

## Key Improvements Made

### 1. Extended Token Lifetimes
- **Access Token**: Extended from 8 hours to 7 days
- **Refresh Token**: 30 days lifetime
- **Cookies**: Aligned with token expiration times

### 2. Refresh Token Implementation
- Added refresh token creation and validation
- Automatic token refresh when access tokens are about to expire
- Secure storage of refresh tokens in HTTP-only cookies

### 3. Frontend Session Management
- Automatic token refresh on the client side
- Proper error handling for expired tokens
- Improved logout functionality

## Implementation Details

### Backend Changes

#### JWT Service (`golang-backend/web/services/jwt_auth.service.go`)
```go
// Extended access token lifetime (7 days)
ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))

// New refresh token functionality
func (jwtas *JWTAuthService) CreateRefreshToken(userId string) (string, error)
func (jwtas *JWTAuthService) RefreshAccessToken(refreshToken string) (string, error)
func (jwtas *JWTAuthService) HandleSetRefreshTokenInCookie(c *fiber.Ctx, refreshToken string) error
func (jwtas *JWTAuthService) ClearAuthCookies(c *fiber.Ctx) error
```

#### Auth Controller (`golang-backend/web/api/controller/auth.controller.go`)
- Updated OAuth callbacks to create both access and refresh tokens
- Added refresh token endpoint (`/auth/refresh`)
- Added logout endpoint (`/auth/logout`)

#### JWT Middleware (`golang-backend/web/api/middleware/jwt.middleware.go`)
- Automatic token refresh when tokens are about to expire
- Fallback to refresh tokens when access tokens are invalid
- Single middleware type: `RequireAuth` for all protected API routes

### Frontend Changes

#### AuthContext (`react-frontend/src/context/AuthContext.tsx`)
- Automatic token refresh on app initialization
- Periodic token refresh before expiration
- Proper error handling and cleanup
- Improved logout functionality

#### User Type (`react-frontend/src/types/types.d.ts`)
- Added optional `image` and `email_verified` fields to match JWT structure

## Session Persistence Strategies

### 1. Token Refresh Flow
```
User logs in → Access Token (7 days) + Refresh Token (30 days)
↓
Access Token expires → Use Refresh Token to get new Access Token
↓
Refresh Token expires → User must log in again
```

### 2. Automatic Refresh Triggers
- **Frontend**: Checks token expiration every minute, refreshes 1 minute before expiry
- **Backend**: Middleware checks token expiration on each request, refreshes if within 5 minutes of expiry

### 3. Security Considerations
- Refresh tokens stored in HTTP-only cookies (not accessible via JavaScript)
- Access tokens also HTTP-only for better security
- Proper token validation and error handling
- Automatic cleanup of expired tokens

## Usage Examples

### Backend Routes
```go
// OAuth login (creates both tokens)
GET /auth/google
GET /auth/discord

// Token management
POST /auth/refresh    // Refresh access token
POST /auth/logout     // Clear all tokens

// Protected API routes (use RequireAuth)
GET /games            // Get all games
POST /games/new       // Create new game
GET /ws/game/:id      // WebSocket connection
```

### Frontend Usage
```typescript
const { user, logout, refreshToken } = useAuth();

// Automatic token refresh happens in background
// Manual refresh if needed
await refreshToken();

// Logout
await logout();
```

## Configuration

### Environment Variables
Ensure these are set in your `.env` file:
```
JWT_VERIFY_SECRET=your-secret-key
GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
DISCORD_OAUTH_CLIENT_ID=your-discord-client-id
DISCORD_OAUTH_CLIENT_SECRET=your-discord-client-secret
```

### Cookie Settings
- **Access Token**: 7 days, HTTP-only, SameSite=Strict
- **Refresh Token**: 30 days, HTTP-only, SameSite=Strict
- **Secure**: Set to `true` in production (HTTPS required)

## Best Practices

### 1. Token Security
- Always use HTTP-only cookies for token storage
- Implement proper token validation
- Use secure random secrets for JWT signing
- Rotate secrets periodically in production

### 2. Error Handling
- Graceful degradation when refresh fails
- Clear user session on authentication errors
- Provide clear error messages to users

### 3. Performance
- Minimize token refresh frequency
- Use efficient token validation
- Cache user information when appropriate

### 4. User Experience
- Seamless token refresh (user shouldn't notice)
- Clear logout functionality
- Proper loading states during authentication

## Troubleshooting

### Common Issues

1. **Token Refresh Fails**
   - Check refresh token validity
   - Verify JWT secret is consistent
   - Ensure cookies are properly set

2. **Session Expires Unexpectedly**
   - Verify token expiration times
   - Check cookie settings
   - Ensure automatic refresh is working

3. **CORS Issues**
   - Configure CORS to allow credentials
   - Set proper cookie domains
   - Check SameSite cookie settings

### Debug Steps

1. Check browser cookies for token presence
2. Verify token expiration times
3. Monitor network requests for refresh calls
4. Check server logs for authentication errors

## Production Considerations

### Security
- Use HTTPS in production
- Set `Secure: true` for cookies in production
- Implement rate limiting on auth endpoints
- Monitor for suspicious authentication patterns

### Performance
- Consider Redis for token storage
- Implement token blacklisting for logout
- Use CDN for static assets
- Monitor authentication performance

### Monitoring
- Log authentication events
- Monitor token refresh patterns
- Track user session durations
- Alert on authentication failures

## Future Enhancements

1. **Remember Me Functionality**
   - Longer refresh token lifetime for "remember me"
   - Separate token storage for persistent sessions

2. **Multi-Device Support**
   - Track active sessions per user
   - Allow session management from dashboard

3. **Advanced Security**
   - Device fingerprinting
   - Location-based authentication
   - Two-factor authentication integration

4. **Analytics**
   - Session duration tracking
   - Authentication method preferences
   - User engagement metrics 