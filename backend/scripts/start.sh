#!/bin/bash
# All-in-one deployment script for petricia.store
# This script builds, starts services, and obtains SSL certificates

set -e

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "ERROR: .env file not found!"
    echo "Please create .env file with DATABASE_URL variable."
    echo "Example: DATABASE_URL=postgres://user:pass@host:5432/dbname"
    exit 1
fi

echo "=== Starting petricia.store deployment ==="

# Create necessary directories
mkdir -p ssl-certificates html

# Check if SSL certificates already exist
if [ -f "ssl-certificates/live/petricia.store/fullchain.pem" ]; then
    echo "SSL certificates already exist. Starting all services..."
    docker compose up -d
    echo "=== Deployment complete ==="
    echo "Your site is available at https://petricia.store"
else
    echo "No SSL certificates found. Starting services for initial setup..."
    
    # Start nginx first (HTTP only) - we'll use a temp config for certbot
    cat > nginx/temp-http.conf << 'EOF'
server {
    listen 80;
    server_name petricia.store www.petricia.store;
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }
    location / {
        proxy_pass http://petricia-api:9090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

    # Build and start api first
    echo "Building and starting API..."
    docker compose build --no-cache api
    docker compose up -d api
    
    # Temporarily use HTTP config for nginx
    cp nginx/nginx-ssl.conf nginx/nginx-ssl.conf.bak
    cp nginx/temp-http.conf nginx/nginx-ssl.conf
    
    # Start nginx
    echo "Starting nginx..."
    docker compose up -d nginx
    
    # Wait for services
    sleep 5
    
    # Get SSL certificates
    echo "Obtaining SSL certificates from Let's Encrypt..."
    docker compose run --rm certbot certonly \
        --webroot \
        --webroot-path=/var/www/html \
        --email admin@petricia.store \
        --agree-tos \
        --no-eff-email \
        -d petricia.store \
        -d www.petricia.store
    
    # Restore SSL config
    mv nginx/nginx-ssl.conf.bak nginx/nginx-ssl.conf
    
    # Restart nginx with SSL
    echo "Restarting nginx with SSL..."
    docker compose restart nginx
    
    echo "=== Deployment complete ==="
    echo "Your site is available at https://petricia.store"
    
    # Clean up temp file
    rm -f nginx/temp-http.conf
fi
