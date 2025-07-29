# Neural Decks

Neural Decks is a real-time online card game you can play with friends or random people. Our AI generates unique, hilarious card decks that keep the game fresh every time you play.

## 🚀 Quick Start (Recommended)

The fastest way to get started is using Docker Compose, which sets up everything automatically:

### Using Docker (Recommended for new developers)

```bash
# Clone the repository
git clone <your-repo-url>
cd neural-decks

# Start the entire application stack
docker-compose -f docker-compose.dev.yml up --build
```

This will start:
- **Frontend**: http://localhost:5173 (React with hot reloading)
- **Backend**: http://localhost:8080 (Go API)
- **Database**: PostgreSQL on port 5432
- **Cache**: Redis on port 6379
- **Redis Management**: http://localhost:8081 (Redis Commander)

### Production Mode

```bash
# Build and start production containers
docker-compose up --build
```

**Production URLs:**
- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080

## 🔐 OAuth Authentication Setup

Neural Decks supports OAuth authentication with Google and Discord. You'll need to set up OAuth applications for both providers.

### Google OAuth Setup

1. **Create Google OAuth Application:**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one
   - Enable the Google+ API
   - Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client IDs"
   - Set Application Type to "Web application"
   - Add authorized redirect URIs:
     - Development: `http://localhost:8080/auth/google/callback`
     - Production: `https://yourdomain.com/auth/google/callback`

2. **Get OAuth Credentials:**
   - Copy the Client ID and Client Secret
   - Add them to your environment variables

### Discord OAuth Setup

1. **Create Discord OAuth Application:**
   - Go to [Discord Developer Portal](https://discord.com/developers/applications)
   - Click "New Application"
   - Go to "OAuth2" → "General"
   - Add redirect URIs:
     - Development: `http://localhost:8080/auth/discord/callback`
     - Production: `https://yourdomain.com/auth/discord/callback`
   - Copy the Client ID and Client Secret

### Environment Variables

Create a `.env` file in the `golang-backend/web` directory with the following variables:

```bash
# OAuth Configuration
GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
GOOGLE_OAUTH_REDIRECT_URI=http://localhost:8080/auth/google/callback

DISCORD_OAUTH_CLIENT_ID=your-discord-client-id
DISCORD_OAUTH_CLIENT_SECRET=your-discord-client-secret
DISCORD_OAUTH_REDIRECT_URI=http://localhost:8080/auth/discord/callback

# JWT Configuration
JWT_VERIFY_SECRET=your-super-secret-jwt-key-at-least-32-characters

# Database Configuration
DATABASE_DSN=postgresql://postgres:password@localhost:5432/neural_decks

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# AI Configuration
CHATGPT_API_KEY=your-openai-api-key

# Application Environment
APP_ENV=development
```

### OAuth Flow

1. **User clicks "Login with Google/Discord"**
2. **Redirected to OAuth provider** (Google/Discord)
3. **User authorizes the application**
4. **Callback to backend** with authorization code
5. **Backend exchanges code for tokens**
6. **Backend creates/updates user in database**
7. **Backend creates JWT access and refresh tokens**
8. **User redirected to game lobby**

### Session Management

- **Access Tokens**: 7-day lifetime, stored in HTTP-only cookies
- **Refresh Tokens**: 30-day lifetime, stored in HTTP-only cookies
- **Automatic Refresh**: Tokens are automatically refreshed when about to expire
- **Secure Logout**: All tokens are invalidated on logout

## 🛠️ Local Development Setup

If you prefer to run services locally without Docker:

### Prerequisites

- **Go 1.24+**: [Install Go](https://go.dev/doc/install)
- **Node.js 20+**: [Install Node.js](https://nodejs.org/)
- **PostgreSQL 15+**: [Install PostgreSQL](https://www.postgresql.org/download/)
- **Redis 7+**: [Install Redis](https://redis.io/download)

### Backend Setup

```bash
# Navigate to backend directory
cd golang-backend/web

# Install dependencies
go mod tidy

# Set up environment variables (create .env file)
cp .env.example .env
# Edit .env with your OAuth credentials and database settings

# Run the backend
go run cmd/http/main.go
```

### Frontend Setup

```bash
# Navigate to frontend directory
cd react-frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## 📁 Project Structure

```
neural-decks/
├── golang-backend/          # Go backend API
│   ├── web/                # Main backend code
│   │   ├── api/           # API routes and controllers
│   │   │   ├── controller/ # OAuth and game controllers
│   │   │   ├── middleware/ # JWT authentication middleware
│   │   │   └── route/     # API route definitions
│   │   ├── bootstrap/     # Application initialization
│   │   ├── cmd/http/      # Main entry point
│   │   ├── config/        # OAuth configuration
│   │   ├── domain/        # Data models and business logic
│   │   ├── repository/    # Database access layer
│   │   └── services/      # JWT and business logic services
│   ├── Dockerfile         # Production Docker image
│   └── Dockerfile.dev     # Development Docker image
├── react-frontend/         # React frontend
│   ├── src/               # Source code
│   │   ├── components/    # React components
│   │   │   ├── login/     # OAuth login components
│   │   │   └── shared/    # Shared UI components
│   │   ├── context/       # React context providers
│   │   ├── hooks/         # Custom React hooks
│   │   └── types/         # TypeScript type definitions
│   ├── Dockerfile         # Production Docker image
│   └── Dockerfile.dev     # Development Docker image
├── docker-compose.yml      # Production Docker setup
├── docker-compose.dev.yml  # Development Docker setup
├── OAUTH2_SESSION_PERSISTENCE.md # OAuth session management guide
└── DOCKER_README.md       # Detailed Docker documentation
```

## 🎮 Game Features

- **Real-time multiplayer**: Play with friends or random players
- **AI-generated cards**: Unique, hilarious cards every game
- **WebSocket communication**: Real-time game updates
- **OAuth authentication**: Secure login with Google and Discord
- **Session persistence**: Extended login sessions with refresh tokens
- **Game rooms**: Create or join game sessions
- **Responsive design**: Works on desktop and mobile

## 🛠️ Development Workflow

### Making Changes

1. **Frontend changes**: Edit files in `react-frontend/src/` - changes auto-reload
2. **Backend changes**: Edit files in `golang-backend/web/` - Air will restart the server
3. **Database changes**: Use the PostgreSQL connection or Redis Commander for data management

### Useful Commands

```bash
# View all container logs
docker-compose -f docker-compose.dev.yml logs -f

# View specific service logs
docker-compose -f docker-compose.dev.yml logs -f backend
docker-compose -f docker-compose.dev.yml logs -f frontend

# Access containers for debugging
docker-compose -f docker-compose.dev.yml exec backend sh
docker-compose -f docker-compose.dev.yml exec frontend sh

# Rebuild specific service
docker-compose -f docker-compose.dev.yml build backend
docker-compose -f docker-compose.dev.yml build frontend

# Stop all services
docker-compose -f docker-compose.dev.yml down
```

### Database Management

```bash
# Access PostgreSQL
docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d neural_decks

# Access Redis CLI
docker-compose -f docker-compose.dev.yml exec redis redis-cli

# View Redis data in browser
# Open http://localhost:8081 (Redis Commander)
```

## 🔧 Configuration

### Environment Variables

The application uses these environment variables:

**OAuth Configuration:**
- `GOOGLE_OAUTH_CLIENT_ID`: Google OAuth client ID
- `GOOGLE_OAUTH_CLIENT_SECRET`: Google OAuth client secret
- `GOOGLE_OAUTH_REDIRECT_URI`: Google OAuth redirect URI
- `DISCORD_OAUTH_CLIENT_ID`: Discord OAuth client ID
- `DISCORD_OAUTH_CLIENT_SECRET`: Discord OAuth client secret
- `DISCORD_OAUTH_REDIRECT_URI`: Discord OAuth redirect URI

**JWT Configuration:**
- `JWT_VERIFY_SECRET`: Secret key for JWT token signing

**Database Configuration:**
- `DATABASE_DSN`: PostgreSQL connection string
- `REDIS_HOST`: Redis server hostname
- `REDIS_PORT`: Redis server port
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_DB`: Redis database number

**AI Configuration:**
- `CHATGPT_API_KEY`: OpenAI API key for card generation

### Default Values (Docker)

- **Database**: `postgresql://postgres:password@localhost:5432/neural_decks`
- **Redis**: `redis://localhost:6379`
- **Backend API**: `http://localhost:8080`
- **OAuth Redirects**: `http://localhost:8080/auth/{provider}/callback`

## 🐛 Troubleshooting

### Common Issues

1. **OAuth Setup Issues:**
   - Verify OAuth client IDs and secrets are correct
   - Ensure redirect URIs match exactly (including protocol)
   - Check that OAuth applications are properly configured
   - Verify environment variables are loaded correctly

2. **Port conflicts**: Make sure ports 3000, 5173, 8080, 5432, 6379, and 8081 are available
3. **Docker build failures**: Check logs with `docker-compose logs [service-name]`
4. **Database connection issues**: Ensure PostgreSQL and Redis are running
5. **Hot reload not working**: Check that volumes are properly mounted in development

### OAuth Troubleshooting

1. **"State Mismatch" Error:**
   - Clear browser cookies and try again
   - Check that OAuth state is being properly managed

2. **"Code-Token Exchange Failed":**
   - Verify OAuth client credentials
   - Check redirect URI configuration
   - Ensure OAuth application is properly set up

3. **"User Data Fetch Failed":**
   - Check OAuth scopes are correctly configured
   - Verify API permissions for the OAuth application

### Reset Everything

```bash
# Stop all containers and remove volumes
docker-compose down -v
docker-compose -f docker-compose.dev.yml down -v

# Remove all images
docker system prune -a

# Start fresh
docker-compose -f docker-compose.dev.yml up --build
```

### Database Reset

```bash
# Stop services
docker-compose down

# Remove database volume
docker volume rm neural-decks_postgres_data_dev

# Restart services
docker-compose -f docker-compose.dev.yml up --build
```

## 📚 Tech Stack

- **Frontend**: React 19, TypeScript, Vite, Tailwind CSS
- **Backend**: Go 1.24, Fiber, GORM, WebSocket
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **AI**: OpenAI GPT integration
- **Real-time**: WebSocket communication
- **Authentication**: OAuth 2.0 (Google, Discord) with JWT tokens
- **Session Management**: Refresh tokens with automatic renewal

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Test thoroughly
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 📖 Additional Documentation

- [OAuth2 Session Persistence Guide](OAUTH2_SESSION_PERSISTENCE.md) - Detailed OAuth session management
- [Docker Setup Guide](DOCKER_README.md) - Detailed Docker documentation
- [API Documentation](docs/api.md) - Backend API reference
- [Frontend Components](docs/frontend.md) - React component documentation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
