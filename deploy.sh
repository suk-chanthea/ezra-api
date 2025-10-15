#!/bin/bash

# Production Deployment Script for Ezra API
# Usage: ./deploy.sh [start|stop|restart|logs|backup|restore|health]

set -e

ENV_FILE=".env.production"
COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="./backups"

# Colors - Using printf for better compatibility
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
BLUE=$(printf '\033[0;34m')
NC=$(printf '\033[0m') # No Color

# Function to print colored messages
print_error() {
    echo "${RED}Error: $1${NC}"
}

print_success() {
    echo "${GREEN}$1${NC}"
}

print_warning() {
    echo "${YELLOW}$1${NC}"
}

print_info() {
    echo "${BLUE}$1${NC}"
}

# Check if .env.production exists
check_env() {
    if [ ! -f "$ENV_FILE" ]; then
        print_error "$ENV_FILE not found!"
        echo "Please create $ENV_FILE from .env.production.example"
        exit 1
    fi
}

# Load environment variables
load_env() {
    if [ -f "$ENV_FILE" ]; then
        export $(cat "$ENV_FILE" | grep -v '^#' | xargs)
    fi
}

# Start services
start() {
    print_info "Starting Ezra API (Production)..."
    check_env
    load_env
    docker-compose -f "$COMPOSE_FILE" up -d --build
    print_success "Services started successfully!"
    echo ""
    echo "View logs: ./deploy.sh logs"
    echo "Check health: ./deploy.sh health"
}

# Stop services
stop() {
    print_warning "Stopping Ezra API..."
    docker-compose -f "$COMPOSE_FILE" down
    print_success "Services stopped successfully!"
}

# Restart services
restart() {
    print_warning "Restarting Ezra API..."
    stop
    sleep 2
    start
}

# View logs
logs() {
    SERVICE=${1:-}
    if [ -z "$SERVICE" ]; then
        docker-compose -f "$COMPOSE_FILE" logs -f
    else
        docker-compose -f "$COMPOSE_FILE" logs -f "$SERVICE"
    fi
}

# Backup database
backup() {
    print_info "Creating database backup..."
    check_env
    load_env
    
    mkdir -p "$BACKUP_DIR"
    BACKUP_FILE="$BACKUP_DIR/ezra_db_$(date +%Y%m%d_%H%M%S).sql"
    
    docker exec ezra-postgres-prod pg_dump -U "${DB_USER:-postgres}" "${DB_NAME:-ezradb}" > "$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        print_success "Backup created: $BACKUP_FILE"
        
        # Compress backup
        gzip "$BACKUP_FILE"
        print_success "Backup compressed: ${BACKUP_FILE}.gz"
        
        # Keep only last 7 backups
        ls -t "$BACKUP_DIR"/*.gz 2>/dev/null | tail -n +8 | xargs -r rm
        print_info "Old backups cleaned up (keeping last 7)"
    else
        print_error "Backup failed!"
        exit 1
    fi
}

# Restore database
restore() {
    BACKUP_FILE=$1
    if [ -z "$BACKUP_FILE" ]; then
        print_error "Please specify backup file"
        echo "Usage: ./deploy.sh restore <backup_file.sql.gz>"
        exit 1
    fi
    
    if [ ! -f "$BACKUP_FILE" ]; then
        print_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    
    print_warning "Restoring database from: $BACKUP_FILE"
    check_env
    load_env
    
    # Decompress if gzipped
    if [[ "$BACKUP_FILE" == *.gz ]]; then
        gunzip -c "$BACKUP_FILE" | docker exec -i ezra-postgres-prod psql -U "${DB_USER:-postgres}" "${DB_NAME:-ezradb}"
    else
        docker exec -i ezra-postgres-prod psql -U "${DB_USER:-postgres}" "${DB_NAME:-ezradb}" < "$BACKUP_FILE"
    fi
    
    if [ $? -eq 0 ]; then
        print_success "Database restored successfully!"
    else
        print_error "Database restore failed!"
        exit 1
    fi
}

# Check service health
health() {
    print_info "Checking service health..."
    echo ""
    docker-compose -f "$COMPOSE_FILE" ps
    echo ""
    print_info "API Health Check:"
    
    # Try to ping the API
    if curl -s http://localhost:80/ping > /dev/null 2>&1; then
        print_success "✅ API is responding"
        curl -s http://localhost:80/ping | jq . 2>/dev/null || curl -s http://localhost:80/ping
    else
        print_error "❌ API is not responding"
    fi
}

# Show usage
usage() {
    echo "${GREEN}Ezra API Deployment Script${NC}"
    echo ""
    echo "Usage: ./deploy.sh [command]"
    echo ""
    echo "Commands:"
    echo "  ${BLUE}start${NC}          Start all services"
    echo "  ${BLUE}stop${NC}           Stop all services"
    echo "  ${BLUE}restart${NC}        Restart all services"
    echo "  ${BLUE}logs${NC} [service] View logs (optional: specify service)"
    echo "  ${BLUE}backup${NC}         Create database backup"
    echo "  ${BLUE}restore${NC} <file> Restore database from backup"
    echo "  ${BLUE}health${NC}         Check service health"
    echo ""
    echo "Examples:"
    echo "  ./deploy.sh start"
    echo "  ./deploy.sh logs api"
    echo "  ./deploy.sh backup"
    echo "  ./deploy.sh restore backups/ezra_db_20241015_120000.sql.gz"
    echo ""
}

# Main script
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs "$2"
        ;;
    backup)
        backup
        ;;
    restore)
        restore "$2"
        ;;
    health)
        health
        ;;
    *)
        usage
        exit 1
        ;;
esac