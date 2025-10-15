#!/bin/bash

# Pre-deployment Setup Script
# Run this BEFORE ./deploy.sh start

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Ezra API Pre-Deployment Setup ===${NC}"

# 1. Create necessary directories
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p config/ssl
mkdir -p backups
mkdir -p logs

# 2. Set permissions
echo -e "${YELLOW}Setting permissions...${NC}"
chmod 755 config/ssl
chmod 755 backups

# 3. Create .gitkeep for empty directories
touch config/ssl/.gitkeep
touch backups/.gitkeep

# 4. Check if .env.production exists
if [ ! -f ".env.production" ]; then
    echo -e "${YELLOW}Creating .env.production from example...${NC}"
    cp .env.production.example .env.production
    echo -e "${GREEN}✅ Created .env.production${NC}"
    echo -e "${YELLOW}⚠️  Please edit .env.production and set your passwords!${NC}"
else
    echo -e "${GREEN}✅ .env.production already exists${NC}"
fi

# 5. Make deploy.sh executable
chmod +x deploy.sh

echo ""
echo -e "${GREEN}=== Setup Complete! ===${NC}"
echo ""
echo "Next steps:"
echo "1. Edit .env.production and set strong passwords:"
echo "   nano .env.production"
echo ""
echo "2. Deploy the application:"
echo "   ./deploy.sh start"
echo ""
echo "3. Check health:"
echo "   ./deploy.sh health"
echo ""
echo "4. View logs:"
echo "   ./deploy.sh logs"