# Production Deployment Guide

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Domain name (for SSL)
- At least 2GB RAM
- 20GB disk space

## Quick Start

### 1. Clone and Setup

```bash
# Clone repository
git clone https://github.com/suk-chanthea/ezra.git
cd ezra

# Make deployment script executable
chmod +x deploy.sh
```

### 2. Configure Environment

```bash
# Copy example environment file
cp .env.production.example .env.production

# Edit with your settings
nano .env.production
```

**Important: Change these values:**
- `DB_PASSWORD` - Use a strong password (min 16 chars)
- `SECRET_KEY` - Use a random string (min 32 chars)
- `REDIS_PASSWORD` - Use a strong password

### 3. Generate Strong Secrets

```bash
# Generate SECRET_KEY
openssl rand -base64 32

# Generate DB_PASSWORD
openssl rand -base64 24

# Generate REDIS_PASSWORD
openssl rand -base64 16
```

### 4. Deploy

```bash
# Start all services
./deploy.sh start

# Check status
./deploy.sh health

# View logs
./deploy.sh logs
```

## SSL/HTTPS Setup (Recommended for Production)

### Option 1: Let's Encrypt (Free)

```bash
# Install certbot
sudo apt install certbot

# Get certificate
sudo certbot certonly --standalone -d your-domain.com

# Copy certificates
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem config/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem config/ssl/key.pem
sudo chmod 644 config/ssl/*.pem
```

### Option 2: Self-Signed (Development/Testing)

```bash
mkdir -p config/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout config/ssl/key.pem \
  -out config/ssl/cert.pem
```

### Enable HTTPS in Nginx

Edit `config/nginx.conf` and uncomment the HTTPS server block (lines starting with `#`).

```bash
# Restart nginx
docker-compose restart nginx
```

## Management Commands

```bash
# Start services
./deploy.sh start

# Stop services
./deploy.sh stop

# Restart services
./deploy.sh restart

# View all logs
./deploy.sh logs

# View specific service logs
./deploy.sh logs api
./deploy.sh logs postgres
./deploy.sh logs nginx

# Check health
./deploy.sh health

# Create database backup
./deploy.sh backup

# Restore database
./deploy.sh restore backups/ezra_db_20241015_120000.sql.gz
```

## Database Backups

### Automatic Backups (Recommended)

Add to crontab:

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * cd /path/to/ezra && ./deploy.sh backup >> /var/log/ezra-backup.log 2>&1
```

### Manual Backup

```bash
./deploy.sh backup
```

Backups are stored in `./backups/` directory.

## Monitoring

### View Container Status

```bash
docker-compose ps
```

### View Resource Usage

```bash
docker stats
```

### View Logs

```bash
# Real-time logs
docker-compose logs -f

# Last 100 lines
docker-compose logs --tail=100

# Specific service
docker-compose logs -f api
```

## Troubleshooting

### API not responding

```bash
# Check if containers are running
docker-compose ps

# Check API logs
./deploy.sh logs api

# Restart API
docker-compose restart api
```

### Database connection errors

```bash
# Check postgres logs
./deploy.sh logs postgres

# Restart postgres
docker-compose restart postgres

# Connect to postgres shell
docker exec -it ezra-postgres-prod psql -U postgres -d ezradb
```

### High memory usage

```bash
# Check resource usage
docker stats

# Restart all services
./deploy.sh restart
```

## Security Checklist

- [ ] Change default passwords in `.env.production`
- [ ] Use strong SECRET_KEY (32+ characters)
- [ ] Enable HTTPS/SSL
- [ ] Set up firewall (UFW/iptables)
- [ ] Enable rate limiting in nginx
- [ ] Regular database backups
- [ ] Keep Docker images updated
- [ ] Monitor logs regularly
- [ ] Use non-root user in containers
- [ ] Restrict database port access

## Firewall Setup (UFW)

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

## Updating the Application

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
./deploy.sh restart

# Check logs
./deploy.sh logs
```

## Performance Tuning

### Postgres Configuration

Edit `docker-compose.yml` and add under postgres service:

```yaml
command: postgres -c max_connections=200 -c shared_buffers=256MB
```

### Redis Configuration

Edit `docker-compose.yml` and modify redis command:

```yaml
command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru
```

## Production Checklist

Before going live:

- [ ] `.env.production` configured with strong passwords
- [ ] SSL certificates installed
- [ ] Database backups scheduled
- [ ] Firewall configured
- [ ] Monitoring set up
- [ ] Test all API endpoints
- [ ] Load testing completed
- [ ] Error logging configured
- [ ] Documentation updated

## Support

For issues or questions:
- GitHub Issues: https://github.com/suk-chanthea/ezra/issues
- Documentation: https://docs.example.com

## License

[Your License Here]