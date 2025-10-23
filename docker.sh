#!/bin/bash

# FMP Docker Management Script

set -e

case "$1" in
    "up")
        echo "🚀 Starting FMP platform..."
        if [ ! -f .env ]; then
            echo "⚠️  .env file not found. Creating from example..."
            cp env.example .env
            echo "📝 Please edit .env file with your configuration"
        fi
        docker-compose --env-file .env up -d
        echo "✅ Platform started!"
        echo "📱 Mini App: http://localhost:3000"
        echo "📊 Analytics: http://localhost:3001"
        echo "🔌 API: http://localhost:8080"
        ;;
    "down")
        echo "🛑 Stopping FMP platform..."
        docker-compose --env-file .env down
        echo "✅ Platform stopped!"
        ;;
    "restart")
        echo "🔄 Restarting FMP platform..."
        docker-compose --env-file .env restart
        echo "✅ Platform restarted!"
        ;;
    "logs")
        service=${2:-""}
        if [ -n "$service" ]; then
            echo "📋 Showing logs for $service..."
            docker-compose --env-file .env logs -f "$service"
        else
            echo "📋 Showing logs for all services..."
            docker-compose --env-file .env logs -f
        fi
        ;;
    "build")
        echo "🔨 Building FMP platform..."
        docker-compose --env-file .env build --no-cache
        echo "✅ Build completed!"
        ;;
    "clean")
        echo "🧹 Cleaning up..."
        docker-compose --env-file .env down -v
        docker system prune -f
        echo "✅ Cleanup completed!"
        ;;
    "status")
        echo "📊 Platform status:"
        docker-compose --env-file .env ps
        ;;
    "shell")
        service=${2:-"fmp-core"}
        echo "🐚 Opening shell for $service..."
        docker-compose --env-file .env exec "$service" sh
        ;;
    "db")
        echo "🗄️ Opening database shell..."
        docker-compose --env-file .env exec postgres psql -U ${POSTGRES_USER:-fmp_user} -d ${POSTGRES_DB:-fmp_db}
        ;;
    "redis")
        echo "🔴 Opening Redis shell..."
        docker-compose --env-file .env exec redis redis-cli
        ;;
    *)
        echo "FMP Docker Management Script"
        echo ""
        echo "Usage: $0 {up|down|restart|logs|build|clean|status|shell|db|redis}"
        echo ""
        echo "Commands:"
        echo "  up        - Start all services"
        echo "  down      - Stop all services"
        echo "  restart   - Restart all services"
        echo "  logs      - Show logs (optionally for specific service)"
        echo "  build     - Build all images"
        echo "  clean     - Stop services and clean up volumes"
        echo "  status    - Show status of all services"
        echo "  shell     - Open shell in service (default: fmp-core)"
        echo "  db        - Open PostgreSQL shell"
        echo "  redis     - Open Redis shell"
        echo ""
        echo "Examples:"
        echo "  $0 up"
        echo "  $0 logs fmp-core"
        echo "  $0 shell minapp-backend"
        ;;
esac
