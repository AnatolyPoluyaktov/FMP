# Financial Manager Platform (FMP)

A comprehensive financial management platform consisting of a Telegram mini app, analytics dashboard, and backend API.

## 🏗️ Architecture

The platform consists of three main components:

- **fmp-core**: Go-based API server providing CRUD operations for financial data
- **minapp**: Telegram mini app with React frontend and Go backend
- **fmp-analytics**: React-based analytics dashboard for data visualization

## 📋 Spec-First Development

This project follows a **spec-first approach**:

1. **OpenAPI Specification** (`api-spec.yaml`) - Single source of truth for API design
2. **Code Generation** - Go models, server interfaces, and clients generated from spec
3. **Swagger Documentation** - Auto-generated API documentation
4. **Type Safety** - Consistent types across frontend and backend

### Code Generation Workflow

```bash
# Install required tools
make install-tools

# Generate code from OpenAPI spec
make generate

# Build applications
make build
```

## 🚀 Quick Start

### Docker (Recommended)

```bash
# Copy environment variables template
cp env.example .env

# Edit .env file with your configuration
nano .env

# Start entire platform with Docker
make docker-up

# Or use the script directly
./docker.sh up
```

**Access URLs:**
- 📱 **Mini App**: http://localhost:3000
- 📊 **Analytics**: http://localhost:3001  
- 🔌 **API**: http://localhost:8080
- 🗄️ **Database**: localhost:5432

> ⚠️ **Important**: Before running, copy `env.example` to `.env` and configure your environment variables. See [GitHub Secrets Guide](GITHUB_SECRETS.md) for production setup.

### Local Development

```bash
# Install dependencies
make install-deps

# Generate code from OpenAPI spec
make generate

# Start only database services
docker-compose -f docker-compose.dev.yml up -d

# Start development servers locally
make dev-core    # Terminal 1
make dev-minapp  # Terminal 2
```

### Benefits of Spec-First Approach

- ✅ **Consistency** - Single source of truth for API contracts
- ✅ **Type Safety** - Generated types prevent runtime errors
- ✅ **Documentation** - Auto-generated Swagger docs
- ✅ **Client Generation** - Automatic client SDKs
- ✅ **Validation** - Request/response validation out of the box

## 📁 Project Structure

```
FMP/
├── fmp-core/                 # Main API server
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   ├── migrations/      # Database migrations
│   │   ├── models/         # Data models
│   │   └── services/       # Business logic
│   ├── migrations/         # SQL migration files
│   └── main.go
├── minapp/                  # Telegram mini app
│   ├── backend/            # Go backend for mini app
│   └── frontend/           # React frontend
├── fmp-analytics/          # Analytics dashboard
│   └── src/
│       ├── components/     # React components
│       └── services/       # API services
└── deploy/                 # Deployment configuration
    ├── ansible/           # Ansible playbooks
    └── deploy.sh          # Deployment script
```

## 🚀 Features

### Core API (fmp-core)
- ✅ Categories management
- ✅ Transactions CRUD operations
- ✅ Planned expenses and income
- ✅ Category limits and monitoring
- ✅ Analytics endpoints
- ✅ Notification system
- ✅ Swagger API documentation

### Telegram Mini App (minapp)
- ✅ Transaction entry with category search
- ✅ Category management
- ✅ Monthly summary
- ✅ Category limits
- ✅ Planned expenses/income
- ✅ Notifications with Telegram WebApp SDK
- ✅ Responsive mobile-first design

### Analytics Dashboard (fmp-analytics)
- ✅ Overview dashboard with key metrics
- ✅ Monthly analytics with trends
- ✅ Category-based analysis
- ✅ Limit monitoring and exceeded tracking
- ✅ Interactive charts (Recharts)
- ✅ Responsive design

## 🛠️ Technology Stack

- **Backend**: Go 1.21, Gin framework, PostgreSQL
- **Frontend**: React 18, TypeScript, Recharts
- **Telegram**: Telegram WebApp SDK
- **Deployment**: Ansible, GitHub Actions, Nginx
- **Database**: PostgreSQL with migrations

## 📋 Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 13+
- Ansible (for deployment)
- Telegram Bot Token

## 🔧 Local Development

### 1. Database Setup

```bash
# Start PostgreSQL
sudo systemctl start postgresql

# Create database
createdb fmp

# Run migrations
cd fmp-core
go run main.go migrate
```

### 2. FMP Core API

```bash
cd fmp-core

# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=fmp
export DB_USER=your_user
export DB_PASSWORD=your_password
export TELEGRAM_BOT_TOKEN=your_bot_token

# Run the server
go run main.go
```

The API will be available at `http://localhost:8080`

### 3. Mini App Frontend

```bash
cd minapp/frontend

# Install dependencies
npm install

# Set environment variables
export REACT_APP_API_URL=http://localhost:8080/api
export REACT_APP_TELEGRAM_BOT_TOKEN=your_bot_token

# Start development server
npm start
```

### 4. Analytics Dashboard

```bash
cd fmp-analytics

# Install dependencies
npm install

# Set environment variables
export REACT_APP_API_URL=http://localhost:8080/api

# Start development server
npm start
```

## 🚀 Deployment

### Using Ansible (Recommended)

1. **Configure inventory**:
   ```bash
   cp deploy/ansible/inventory.yml.example deploy/ansible/inventory.yml
   # Edit inventory.yml with your server details
   ```

2. **Deploy all components**:
   ```bash
   cd deploy
   ./deploy.sh all
   ```

3. **Deploy specific component**:
   ```bash
   ./deploy.sh core      # Deploy only API
   ./deploy.sh minapp    # Deploy only mini app
   ./deploy.sh analytics # Deploy only analytics
   ```

### Manual Deployment

1. **Build applications**:
   ```bash
   # Build fmp-core
   cd fmp-core && go build -o fmp-core .

   # Build minapp backend
   cd minapp/backend && go build -o fmp-minapp-backend .

   # Build minapp frontend
   cd minapp/frontend && npm run build

   # Build analytics
   cd fmp-analytics && npm run build
   ```

2. **Deploy to servers**:
   - Copy binaries to `/opt/fmp-*/`
   - Configure systemd services
   - Set up Nginx reverse proxy

## 📊 API Documentation

Once the API is running, visit:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- API endpoints: `http://localhost:8080/api/v1/`

### Key Endpoints

- `GET /api/v1/categories` - List categories
- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/analytics/monthly-summary` - Monthly analytics
- `GET /api/v1/notifications` - Get notifications

## 🔐 Environment Variables

### fmp-core
```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fmp
DB_USER=fmp
DB_PASSWORD=password
TELEGRAM_BOT_TOKEN=your_bot_token
PORT=8080
GIN_MODE=release
```

### minapp/frontend
```bash
REACT_APP_API_URL=http://localhost:8080/api
REACT_APP_TELEGRAM_BOT_TOKEN=your_bot_token
```

### fmp-analytics
```bash
REACT_APP_API_URL=http://localhost:8080/api
```

## 🧪 Testing

```bash
# Test fmp-core
cd fmp-core
go test ./...

# Test minapp backend
cd minapp/backend
go test ./...

# Test frontend applications
cd minapp/frontend
npm test

cd fmp-analytics
npm test
```

## 📈 Monitoring

The platform includes built-in monitoring:

- **Health checks**: API health endpoints
- **Logging**: Structured logging with different levels
- **Metrics**: Transaction counts, limit exceedances
- **Notifications**: Daily reminders and limit warnings

## 🔄 CI/CD

GitHub Actions workflows are configured for:

- **fmp-core**: Go tests, build, and deployment
- **minapp**: Frontend and backend tests, build, deployment
- **fmp-analytics**: React tests, build, deployment

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the API documentation
- Review the deployment logs

## 🎯 Roadmap

- [ ] Multi-user support
- [ ] Advanced reporting
- [ ] Mobile app (React Native)
- [ ] Integration with banking APIs
- [ ] Budget planning tools
- [ ] Export to Excel/PDF