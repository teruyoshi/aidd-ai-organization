#!/bin/bash

# AIDD TODO Sample - Production Deployment Script
# Handles production deployment with health checks and rollback

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
COMPOSE_FILE="sample-docker-compose.yml"
COMPOSE_PROD_FILE="sample-docker-compose.prod.yml"
PROJECT_NAME="aidd-todo-sample"
BACKUP_DIR="backups"
DEPLOY_TIMEOUT=300

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

# Check prerequisites
check_prerequisites() {
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi

    # Check environment variables
    if [ -z "$MYSQL_ROOT_PASSWORD" ]; then
        log_error "MYSQL_ROOT_PASSWORD environment variable is required"
        exit 1
    fi

    if [ -z "$JWT_SECRET" ]; then
        log_error "JWT_SECRET environment variable is required"
        exit 1
    fi
}

# Create backup directory
setup_backup_dir() {
    if [ ! -d "$BACKUP_DIR" ]; then
        mkdir -p "$BACKUP_DIR"
        log_info "Created backup directory: $BACKUP_DIR"
    fi
}

# Create database backup before deployment
backup_database() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$BACKUP_DIR/pre_deploy_$timestamp.sql"

    log_info "Creating database backup before deployment..."

    if docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps mysql | grep -q "Up"; then
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec -T mysql \
            mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" todo_db > "$backup_file" 2>/dev/null

        if [ -f "$backup_file" ] && [ -s "$backup_file" ]; then
            log_success "Database backup created: $backup_file"
            echo "$backup_file"
        else
            log_error "Failed to create database backup"
            exit 1
        fi
    else
        log_warning "MySQL service is not running, skipping backup"
        echo ""
    fi
}

# Health check for service
health_check() {
    local service=$1
    local url=$2
    local max_attempts=30
    local attempt=0

    log_info "Performing health check for $service..."

    while [ $attempt -lt $max_attempts ]; do
        if curl -f -s "$url" > /dev/null 2>&1; then
            log_success "$service health check passed"
            return 0
        fi

        attempt=$((attempt + 1))
        echo -n "."
        sleep 5
    done

    log_error "$service health check failed after $max_attempts attempts"
    return 1
}

# Deploy application
deploy_application() {
    local backup_file=$1

    log_info "Starting production deployment..."

    # Build and deploy with production configuration
    if [ -f "$COMPOSE_PROD_FILE" ]; then
        log_info "Using production compose file..."
        docker-compose -f "$COMPOSE_FILE" -f "$COMPOSE_PROD_FILE" -p "$PROJECT_NAME" up -d --build
    else
        log_info "Using default compose file..."
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d --build
    fi

    # Wait for services to be ready
    sleep 10

    # Health checks
    log_info "Performing health checks..."

    # Backend health check
    if ! health_check "backend" "http://localhost:8080/health"; then
        log_error "Backend health check failed"
        rollback_deployment "$backup_file"
        exit 1
    fi

    # Frontend health check
    if ! health_check "frontend" "http://localhost:3000"; then
        log_error "Frontend health check failed"
        rollback_deployment "$backup_file"
        exit 1
    fi

    # Database health check
    if ! docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec mysql mysqladmin ping -h localhost --silent; then
        log_error "Database health check failed"
        rollback_deployment "$backup_file"
        exit 1
    fi

    log_success "All health checks passed!"
}

# Rollback deployment
rollback_deployment() {
    local backup_file=$1

    log_warning "Rolling back deployment..."

    # Stop current containers
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down

    # Restore database if backup exists
    if [ -n "$backup_file" ] && [ -f "$backup_file" ]; then
        log_info "Restoring database from backup: $backup_file"

        # Start only MySQL for restore
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d mysql

        # Wait for MySQL to be ready
        sleep 30

        # Restore database
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec -T mysql \
            mysql -u root -p"$MYSQL_ROOT_PASSWORD" todo_db < "$backup_file"

        log_success "Database restored from backup"
    fi

    log_error "Deployment rollback completed"
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."

    # Wait for backend to be ready
    sleep 5

    if docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec backend test -f cmd/migrate/main.go; then
        docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" exec backend go run cmd/migrate/main.go
        log_success "Database migrations completed"
    else
        log_warning "Migration script not found, skipping..."
    fi
}

# Cleanup old containers and images
cleanup() {
    log_info "Cleaning up old containers and images..."

    # Remove old containers
    docker container prune -f

    # Remove old images (keep last 3 versions)
    docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedSince}}" | \
        grep "$PROJECT_NAME" | tail -n +4 | awk '{print $3}' | xargs -r docker rmi -f

    log_success "Cleanup completed"
}

# Show deployment status
show_status() {
    log_info "Deployment Status:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Service status
    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Application URLs
    echo "🌐 Application URLs:"
    echo "   Frontend: http://localhost:3000"
    echo "   Backend:  http://localhost:8080"
    echo "   API Docs: http://localhost:8080/docs"

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Full deployment process
full_deploy() {
    log_info "Starting full production deployment process..."

    # Check prerequisites
    check_prerequisites

    # Setup backup directory
    setup_backup_dir

    # Create backup
    local backup_file=$(backup_database)

    # Deploy application
    deploy_application "$backup_file"

    # Run migrations
    run_migrations

    # Cleanup
    cleanup

    # Show status
    show_status

    log_success "🚀 Production deployment completed successfully!"
}

# Quick deployment (without backup)
quick_deploy() {
    log_info "Starting quick deployment (no backup)..."

    check_prerequisites

    # Deploy without backup
    deploy_application ""

    # Show status
    show_status

    log_success "🚀 Quick deployment completed!"
}

# Stop production environment
stop_production() {
    log_info "Stopping production environment..."

    docker-compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down

    log_success "Production environment stopped"
}

# Show help
show_help() {
    echo "AIDD TODO Sample - Production Deployment Script"
    echo
    echo "Usage: $0 [COMMAND]"
    echo
    echo "Commands:"
    echo "  deploy       Full production deployment with backup and health checks"
    echo "  quick        Quick deployment without backup"
    echo "  stop         Stop production environment"
    echo "  status       Show deployment status"
    echo "  backup       Create database backup only"
    echo "  cleanup      Cleanup old containers and images"
    echo "  help         Show this help message"
    echo
    echo "Required Environment Variables:"
    echo "  MYSQL_ROOT_PASSWORD    MySQL root password"
    echo "  JWT_SECRET            JWT signing secret"
    echo
    echo "Optional Environment Variables:"
    echo "  MYSQL_DATABASE        Database name (default: todo_db)"
    echo "  MYSQL_USER           Database user (default: todo_user)"
    echo "  MYSQL_PASSWORD       Database password (default: todo_password)"
    echo
    echo "Examples:"
    echo "  export MYSQL_ROOT_PASSWORD=secure_password"
    echo "  export JWT_SECRET=your_jwt_secret_here"
    echo "  $0 deploy"
}

# Main execution
main() {
    local command=${1:-help}

    case $command in
        "deploy"|"full")
            full_deploy
            ;;
        "quick")
            quick_deploy
            ;;
        "stop")
            stop_production
            ;;
        "status")
            show_status
            ;;
        "backup")
            setup_backup_dir
            backup_database
            ;;
        "cleanup")
            cleanup
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

main "$@"