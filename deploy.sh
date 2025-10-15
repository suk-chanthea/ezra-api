#!/bin/bash

# Production Deployment Script for Ezra API
# Usage: ./deploy.sh [start|stop|restart|logs|backup]

set -e

ENV_FILE=".env.production"
COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="./backups"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env.production exists
check_env() {
    if [ ! -f "$ENV_FILE" ]; then
        echo -e "${RED}Error: $ENV_FILE not found!${NC}"
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
    echo -e "${GREEN}Starting Ezra API (Production)...${NC}"
    check_env
    load_env
    docker-compose -f "$COMPOSE_FILE" up -d --build
    echo -e "${GREEN}Services started successfully!${NC}"
    echo "View logs: ./deploy.sh logs"
}

# Stop services
stop() {
    echo -e "${YELLOW}Stopping Ezra API...${NC}"
    docker-compose -f "$COMPOSE_FILE" down
    echo -e "${GREEN}Services stopped successfully!${NC}"
}

# Restart services
restart() {
    echo -e "${YELLOW}Restarting Ezra API...${NC}"
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
    echo -e "${GREEN}Creating database backup...${NC}"
    check_env
    load_env
    
    mkdir -p "$BACKUP_DIR"
    BACKUP_FILE="$BACKUP_DIR/ezra_db_$(date +%Y%m%d_%H%M%S).sql"
    
    docker exec ezra-postgres-prod pg_dump -U "${DB_USER:-postgres}" "${DB_NAME:-ezradb}" > "$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Backup created: $BACKUP_FILE${NC}"
        
        # Compress backup
        gzip "$BACKUP_FILE"
        echo -e "${GREEN}Backup compressed: ${BACKUP_FILE}.gz${NC}"
        
        # Keep only last 7 backups
        ls -t "$BACKUP_DIR"/*.gz | tail -n +8 | xargs -r rm
        echo -e "${GREEN}Old backups cleaned up (keeping last 7)${NC}"
    else
        echo -e "${RED}Backup failed!${NC}"
        exit 1
    fi
}

# Restore database
restore() {
    BACKUP_FILE=$1
    if [ -z "$BACKUP_FILE" ]; then
        echo -e "${RED}Error: Please specify backup file${NC}"
        echo "Usage: ./deploy.sh restore <backup_file.sql.gz>"
        exit 1
    fi
    
    if [ ! -f "$BACKUP_FILE" ]; then
        echo -e "${RED}Error: Backup file not found: $BACKUP_FILE${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Restoring database from: $BACKUP_FILE${NC}"
    check_env
    load_env
    
    # Decompress if gzipped
    if [[ "$BACKUP_FILE" == *.gz ]]; then
        gunzip -c "$BACKUP_FILE" | docker exec -i ezra-postgres-prod psql -U "${DB_USER:-postgres}" "${DB_NAME:-ezradb}"
    else
        docker exec -i ezra-postgres-prod psql -U "${DB_USER:-postgres}" "${DB_NAME:-ezradb}" < "$BACKUP_FILE"
    fi
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Database restored successfully!${NC}"
    else
        echo -e "${RED}Database restore failed!${NC}"
        exit 1
    fi
}

# Check service health
health() {
    echo -e "${GREEN}Checking service health...${NC}"
    docker-compose -f "$COMPOSE_FILE" ps
    echo ""
    echo -e "${GREEN}API Health:${NC}"
    curl -s http://localhost:8090/ping || echo -e "${RED}API not responding${NC}"
}

# Show usage
usage() {
    echo "Ezra API Deployment Script"
    echo ""
    echo "Usage: ./deploy.sh [command]"
    echo ""
    echo "Commands:"
    echo "  start          Start all services"
    echo "  stop           Stop all services"
    echo "  restart        Restart all services"
    echo "  logs [service] View logs (optional: specify service)"
    echo "  backup         Create database backup"
    echo "  restore <file> Restore database from backup"
    echo "  health         Check service health"
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