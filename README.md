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
# Edit .env with your database and Redis credentials

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
│   │   ├── bootstrap/     # Application initialization
│   │   ├── cmd/http/      # Main entry point
│   │   ├── domain/        # Data models and business logic
│   │   ├── repository/    # Database access layer
│   │   └── services/      # Business logic services
│   ├── Dockerfile         # Production Docker image
│   └── Dockerfile.dev     # Development Docker image
├── react-frontend/         # React frontend
│   ├── src/               # Source code
│   │   ├── components/    # React components
│   │   ├── context/       # React context providers
│   │   ├── hooks/         # Custom React hooks
│   │   └── types/         # TypeScript type definitions
│   ├── Dockerfile         # Production Docker image
│   └── Dockerfile.dev     # Development Docker image
├── docker-compose.yml      # Production Docker setup
├── docker-compose.dev.yml  # Development Docker setup
└── DOCKER_README.md       # Detailed Docker documentation
```

## 🎮 Game Features

- **Real-time multiplayer**: Play with friends or random players
- **AI-generated cards**: Unique, hilarious cards every game
- **WebSocket communication**: Real-time game updates
- **User authentication**: Secure login and registration
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

**Backend:**
- `REDIS_HOST`: Redis server hostname
- `REDIS_PORT`: Redis server port
- `POSTGRES_HOST`: PostgreSQL server hostname
- `POSTGRES_PORT`: PostgreSQL server port
- `POSTGRES_USER`: Database username
- `POSTGRES_PASSWORD`: Database password
- `POSTGRES_DB`: Database name

**Frontend:**
- `VITE_API_URL`: Backend API URL

### Default Values (Docker)

- **Database**: `postgresql://postgres:password@localhost:5432/neural_decks`
- **Redis**: `redis://localhost:6379`
- **Backend API**: `http://localhost:8080`

## 🐛 Troubleshooting

### Common Issues

1. **Port conflicts**: Make sure ports 3000, 5173, 8080, 5432, 6379, and 8081 are available
2. **Docker build failures**: Check logs with `docker-compose logs [service-name]`
3. **Database connection issues**: Ensure PostgreSQL and Redis are running
4. **Hot reload not working**: Check that volumes are properly mounted in development

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
- **Authentication**: JWT tokens

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Test thoroughly
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 📖 Additional Documentation

- [Docker Setup Guide](DOCKER_README.md) - Detailed Docker documentation
- [API Documentation](docs/api.md) - Backend API reference
- [Frontend Components](docs/frontend.md) - React component documentation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
