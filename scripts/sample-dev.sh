#!/bin/bash

# AIDD TODO Sample - Development Environment Management Script
# Manages the complete development environment with Docker Compose

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="sample-docker-compose.yml"
PROJECT_NAME="aidd-todo-sample"
SERVICES=("mysql" "redis" "backend" "frontend" "nginx")

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        log_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Check if docker-compose file exists
check_compose_file() {
    if [ ! -f "$COMPOSE_FILE" ]; then
        log_error "Docker Compose file '$COMPOSE_FILE' not found."
        log_info "Please run this script from the project root directory."
        exit 1
    fi
}

# Wait for service to be healthy
wait_for_service() {
    local service=$1
    local max_attempts=30
    local attempt=0

    log_info "Waiting for $service to be healthy..."

    while [ $attempt -lt $max_attempts ]; do
        if docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps "$service" | grep -q "Up (healthy)"; then
            log_success "$service is healthy"
            return 0
        fi

        attempt=$((attempt + 1))
        echo -n "."
        sleep 2
    done

    log_error "$service failed to become healthy within ${max_attempts} attempts"
    return 1
}

# Start all services
start_services() {
    log_info "Starting AIDD TODO Sample development environment..."

    # Build and start services
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d --build

    # Wait for critical services
    wait_for_service "mysql"
    wait_for_service "redis"
    wait_for_service "backend"
    wait_for_service "frontend"

    log_success "All services started successfully!"
    show_urls
}

# Stop all services
stop_services() {
    log_info "Stopping AIDD TODO Sample development environment..."
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down
    log_success "All services stopped successfully!"
}

# Restart all services
restart_services() {
    log_info "Restarting AIDD TODO Sample development environment..."
    stop_services
    sleep 2
    start_services
}

# Show service status
show_status() {
    log_info "Service Status:"
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps
}

# Show logs for all or specific service
show_logs() {
    local service=${1:-}

    if [ -n "$service" ]; then
        log_info "Showing logs for $service..."
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f "$service"
    else
        log_info "Showing logs for all services..."
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f
    fi
}

# Clean up everything (containers, volumes, images)
clean_all() {
    log_warning "This will remove all containers, volumes, and images for the project."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo

    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "Cleaning up all project resources..."

        # Stop and remove containers, networks, volumes
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down -v --remove-orphans

        # Remove project images
        docker images --filter "label=project=$PROJECT_NAME" -q | xargs -r docker rmi -f

        # Remove dangling images
        docker image prune -f

        log_success "Cleanup completed!"
    else
        log_info "Cleanup cancelled."
    fi
}

# Initialize the database with sample data
init_db() {
    log_info "Initializing database with sample data..."

    # Wait for backend to be ready
    wait_for_service "backend"

    # Run database migrations
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec backend go run cmd/migrate/main.go

    # Seed with sample data (if seed script exists)
    if docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec backend test -f cmd/seed/main.go; then
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec backend go run cmd/seed/main.go
        log_success "Database initialized with sample data!"
    else
        log_success "Database migrations completed!"
    fi
}

# Run tests
run_tests() {
    log_info "Running all tests..."

    # Backend tests
    log_info "Running backend tests..."
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec backend go test -v ./...

    # Frontend tests
    log_info "Running frontend tests..."
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec frontend npm run test

    log_success "All tests completed!"
}

# Show application URLs
show_urls() {
    echo
    log_success "🚀 AIDD TODO Sample is running!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📱 Frontend:          http://localhost:3000"
    echo "🔧 API Documentation: http://localhost:8080/docs"
    echo "📊 Prometheus:        http://localhost:9090"
    echo "📈 Grafana:          http://localhost:3001"
    echo "🗄️  MySQL:            localhost:3306 (user: todo_user, db: todo_db)"
    echo "🔴 Redis:            localhost:6379"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo
}

# Backup database
backup_db() {
    local backup_file="backup_$(date +%Y%m%d_%H%M%S).sql"
    log_info "Creating database backup: $backup_file"

    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec -T mysql \
        mysqldump -u todo_user -ptodo_password todo_db > "$backup_file"

    log_success "Database backup created: $backup_file"
}

# Show help
show_help() {
    echo "AIDD TODO Sample - Development Environment Manager"
    echo
    echo "Usage: $0 [COMMAND]"
    echo
    echo "Commands:"
    echo "  start     Start all services"
    echo "  stop      Stop all services"
    echo "  restart   Restart all services"
    echo "  status    Show service status"
    echo "  logs      Show logs for all services"
    echo "  logs [service]  Show logs for specific service"
    echo "  clean     Remove all containers, volumes, and images"
    echo "  init-db   Initialize database with migrations and sample data"
    echo "  test      Run all tests (backend and frontend)"
    echo "  backup    Create database backup"
    echo "  urls      Show application URLs"
    echo "  help      Show this help message"
    echo
    echo "Available services: ${SERVICES[*]}"
    echo
    echo "Examples:"
    echo "  $0 start              # Start development environment"
    echo "  $0 logs backend       # Show backend logs"
    echo "  $0 restart            # Restart all services"
    echo "  $0 clean              # Clean up everything"
}

# Main execution
main() {
    local command=${1:-help}

    # Check prerequisites
    check_docker

    if [ "$command" != "help" ]; then
        check_compose_file
    fi

    case $command in
        "start")
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "$2"
            ;;
        "clean")
            clean_all
            ;;
        "init-db")
            init_db
            ;;
        "test")
            run_tests
            ;;
        "backup")
            backup_db
            ;;
        "urls")
            show_urls
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            echo
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"