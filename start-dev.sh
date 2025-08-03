#!/bin/bash

echo "🚀 Starting Kanban API Server..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the kanban-api directory"
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Check if port 8080 is already in use
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  Port 8080 is already in use. Stopping existing process..."
    lsof -ti:8080 | xargs kill -9
fi

# Start the server
echo "🔧 Starting server on port 8080..."
go run cmd/api/main.go 