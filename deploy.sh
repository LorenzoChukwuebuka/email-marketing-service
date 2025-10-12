#!/bin/bash

# Deployment script for CrabMailer

set -e  # Exit on any error

echo "🚀 Starting CrabMailer deployment..."

# Configuration
APP_DIR="/opt/crabmailer"
BACKUP_DIR="/opt/crabmailer-backups"
DOCKER_COMPOSE_FILE="docker-compose.prod.yml"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Function to backup database
backup_database() {
    echo "📦 Creating database backup..."
    BACKUP_FILE="$BACKUP_DIR/db-backup-$(date +%Y%m%d-%H%M%S).sql"
    docker exec prod-crabmailer-postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB > $BACKUP_FILE
    echo "✅ Database backed up to: $BACKUP_FILE"
}

# Function to backup volumes
backup_volumes() {
    echo "📦 Creating volume backups..."
    BACKUP_DATE=$(date +%Y%m%d-%H%M%S)
    mkdir -p $BACKUP_DIR/volumes-$BACKUP_DATE
    
    # Backup important volumes
    docker run --rm -v crabmailer_uploads:/data -v $BACKUP_DIR/volumes-$BACKUP_DATE:/backup alpine tar czf /backup/uploads.tar.gz -C /data .
    docker run --rm -v crabmailer_templates:/data -v $BACKUP_DIR/volumes-$BACKUP_DATE:/backup alpine tar czf /backup/templates.tar.gz -C /data .
    docker run --rm -v crabmailer_smtp_settings:/data -v $BACKUP_DIR/volumes-$BACKUP_DATE:/backup alpine tar czf /backup/smtp_settings.tar.gz -C /data .
    
    echo "✅ Volumes backed up to: $BACKUP_DIR/volumes-$BACKUP_DATE"
}

# Check if we're in the right directory
if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
    echo "❌ docker-compose.prod.yml not found in current directory"
    echo "Please run this script from the application root directory"
    exit 1
fi

# Load environment variables
if [ -f ".env.production" ]; then
    source .env.production
else
    echo "❌ .env.production file not found!"
    exit 1
fi

# Check if this is an update (containers exist)
if docker ps -a | grep -q "prod-crabmailer"; then
    echo "🔄 Existing deployment detected. Creating backups..."
    
    # Create backups before update
    backup_database
    backup_volumes
    
    echo "🛑 Stopping existing containers..."
    docker compose -f $DOCKER_COMPOSE_FILE down
else
    echo "🆕 New deployment detected."
fi

# Pull latest images
echo "⬇️ Pulling latest images..."
docker compose -f $DOCKER_COMPOSE_FILE pull

# Build and start services
echo "🏗️ Building and starting services..."
docker compose -f $DOCKER_COMPOSE_FILE up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 30

# Check service health
echo "🔍 Checking service health..."
docker compose -f $DOCKER_COMPOSE_FILE ps

# Test if services are responding
echo "🧪 Testing services..."

# Test frontend
if curl -f -s -o /dev/null "http://localhost:80"; then
    echo "✅ Frontend is responding"
else
    echo "❌ Frontend is not responding"
fi

# Test API (if it has a health endpoint)
if curl -f -s -o /dev/null "http://localhost:9000/health" 2>/dev/null; then
    echo "✅ API is responding"
else
    echo "⚠️ API health check failed (this might be normal if no health endpoint exists)"
fi

# Show running containers
echo "🐳 Running containers:"
docker ps --filter "name=prod-crabmailer"

echo "✅ Deployment completed!"
echo "🌍 Your application should be available at:"
echo "   - Frontend: https://yourdomain.com"
echo "   - API: https://api.yourdomain.com"
echo ""
echo "📝 Next steps:"
echo "1. Update your domain names in the docker-compose.prod.yml file"
echo "2. Update your email in the Let's Encrypt configuration"
echo "3. Test your SSL certificates"
echo "4. Configure your DNS MX records for SMTP"
echo ""
echo "📊 To view logs: docker compose -f $DOCKER_COMPOSE_FILE logs -f [service_name]"
echo "🛑 To stop: docker compose -f $DOCKER_COMPOSE_FILE down"