# Neural Decks - Docker Setup

This project includes Docker Compose configurations for running the entire application stack locally.

## Architecture

- **Frontend**: React + Vite (port 3000 in production, 5173 in development)
- **Backend**: Go + Fiber (port 8080)
- **Database**: PostgreSQL (port 5432)
- **Cache**: Redis (port 6379)
- **Redis Management**: Redis Commander (port 8081, development only)

## Quick Start

### Production Mode

To run the application in production mode:

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up --build -d
```

### Development Mode

To run the application in development mode with hot reloading:

```bash
# Build and start all services with hot reloading
docker-compose -f docker-compose.dev.yml up --build

# Or run in detached mode
docker-compose -f docker-compose.dev.yml up --build -d
```

## Services

### Frontend (React)
- **Production**: Served via Nginx on port 3000
- **Development**: Vite dev server on port 5173 with hot reloading
- **URL**: http://localhost:3000 (production) or http://localhost:5173 (development)

### Backend (Go)
- **Port**: 8080
- **URL**: http://localhost:8080
- **Development**: Uses Air for hot reloading
- **Production**: Optimized multi-stage build

### Database (PostgreSQL)
- **Port**: 5432
- **Database**: neural_decks
- **User**: postgres
- **Password**: password
- **Connection**: postgresql://postgres:password@localhost:5432/neural_decks

### Cache (Redis)
- **Port**: 6379
- **Persistence**: AOF enabled
- **Connection**: redis://localhost:6379

### Redis Commander (Development Only)
- **Port**: 8081
- **URL**: http://localhost:8081
- **Purpose**: Web-based Redis management interface

## Environment Variables

The following environment variables are automatically set in the Docker containers:

### Backend
- `REDIS_HOST=redis`
- `REDIS_PORT=6379`
- `POSTGRES_HOST=postgres`
- `POSTGRES_PORT=5432`
- `POSTGRES_USER=postgres`
- `POSTGRES_PASSWORD=password`
- `POSTGRES_DB=neural_decks`

### Frontend (Development)
- `VITE_API_URL=http://localhost:8080`

## Useful Commands

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Stop services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Rebuild specific service
```bash
# Rebuild backend
docker-compose build backend

# Rebuild frontend
docker-compose build frontend
```

### Access containers
```bash
# Access backend container
docker-compose exec backend sh

# Access frontend container
docker-compose exec frontend sh

# Access database
docker-compose exec postgres psql -U postgres -d neural_decks
```

### Database Management
```bash
# Create database backup
docker-compose exec postgres pg_dump -U postgres neural_decks > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres neural_decks < backup.sql
```

## Development Workflow

1. **Start development environment**:
   ```bash
   docker-compose -f docker-compose.dev.yml up --build
   ```

2. **Make changes to your code** - the services will automatically reload

3. **View logs for debugging**:
   ```bash
   docker-compose -f docker-compose.dev.yml logs -f backend
   ```

4. **Stop development environment**:
   ```bash
   docker-compose -f docker-compose.dev.yml down
   ```

## Production Deployment

For production deployment:

1. **Build and start**:
   ```bash
   docker-compose up --build -d
   ```

2. **Check service health**:
   ```bash
   docker-compose ps
   ```

3. **View logs**:
   ```bash
   docker-compose logs -f
   ```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Make sure ports 3000, 8080, 5432, and 6379 are available
2. **Permission issues**: On Linux/Mac, you might need to use `sudo` for Docker commands
3. **Build failures**: Check the logs with `docker-compose logs [service-name]`

### Reset Everything

To completely reset the environment:

```bash
# Stop all containers and remove volumes
docker-compose down -v
docker-compose -f docker-compose.dev.yml down -v

# Remove all images
docker-compose down --rmi all
docker-compose -f docker-compose.dev.yml down --rmi all

# Clean up Docker system
docker system prune -a
```

### Database Reset

To reset the database:

```bash
# Stop services
docker-compose down

# Remove database volume
docker volume rm neural-decks_postgres_data

# Restart services
docker-compose up --build
```

## Network

All services are connected via the `neural-decks-network` bridge network, allowing them to communicate using service names as hostnames.

## Volumes

- `postgres_data`: PostgreSQL database persistence
- `redis_data`: Redis data persistence

These volumes ensure data persists between container restarts. 