#!/bin/bash
set -e

echo "🚀 Hypercommit Local Build & Deploy Script"
echo "=========================================="

# Configuration
SSH_HOST="root@hypercode.ovh"
SSH_PORT="2222"
REMOTE_PATH="/opt/hypercode"
SERVICE_NAME="hypercode"

echo "📦 Building CSS assets..."
bunx @tailwindcss/cli -i ./views/styles/main.css -o ./public/assets/styles.css

echo "🔨 Building binary..."
mkdir -p bin
CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/hypercommit ./cmd/server

echo "📤 Uploading binary to server..."
scp -P $SSH_PORT bin/hypercommit $SSH_HOST:$REMOTE_PATH/hypercommit.new

echo "🔄 Deploying on server..."
ssh -p $SSH_PORT $SSH_HOST << 'ENDSSH'
set -e
cd /opt/hypercode

# Stop service
echo "⏸️  Stopping service..."
systemctl stop hypercode

# Backup current binary
if [ -f hypercode ]; then
    BACKUP_NAME="hypercode.backup.$(date +%Y%m%d_%H%M%S)"
    mv hypercode "$BACKUP_NAME"
    echo "💾 Created backup: $BACKUP_NAME"
fi

# Deploy new binary
mv hypercommit.new hypercommit
chmod +x hypercommit

# Start service
echo "▶️  Starting service..."
systemctl start hypercode

# Check status
sleep 2
if systemctl is-active --quiet hypercode; then
    echo "✅ Service is running"
    systemctl status hypercode --no-pager | head -10
else
    echo "❌ Service failed to start"
    journalctl -u hypercode --since '30 seconds ago' --no-pager
    exit 1
fi
ENDSSH

echo ""
echo "✨ Deployment complete!"
echo ""
echo "View logs: ssh -p $SSH_PORT $SSH_HOST 'journalctl -u $SERVICE_NAME -f'"
