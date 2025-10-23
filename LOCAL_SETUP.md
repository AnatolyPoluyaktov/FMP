# Инструкция по локальному запуску FMP

## Быстрый старт

### 1. Подготовка окружения

```bash
# Клонируйте репозиторий и перейдите в директорию
cd /Users/andromeda/personal_projects/FMP

# Скопируйте файл с переменными окружения
cp env.example .env

# Отредактируйте .env файл с вашими настройками
# Особенно важно указать правильный TELEGRAM_BOT_TOKEN
```

### 2. Запуск базы данных и Redis

```bash
# Запустите только базу данных и Redis
docker-compose --env-file .env up -d postgres redis

# Проверьте, что сервисы запущены
docker ps
```

### 3. Запуск fmp-core локально

```bash
# Перейдите в директорию fmp-core
cd fmp-core

# Установите зависимости
go mod tidy

# Запустите сервер с переменными окружения
DB_HOST=localhost DB_PORT=5432 DB_NAME=fmp_db DB_USER=fmp_user DB_PASSWORD=fmp_password_secure_123 go run main.go
```

### 4. Запуск minapp-backend локально

```bash
# В новом терминале, перейдите в директорию minapp/backend
cd minapp/backend

# Установите зависимости
go mod tidy

# Запустите сервер
FMP_CORE_URL=http://localhost:8080/api/v1 PORT=8081 go run main.go
```

### 5. Запуск minapp-frontend локально

```bash
# В новом терминале, перейдите в директорию minapp/frontend
cd minapp/frontend

# Установите зависимости
npm install

# Запустите фронтенд
npm run dev
```

## Проверка работы

После запуска всех сервисов:

1. **fmp-core API**: http://localhost:8080/health
2. **minapp-backend**: http://localhost:8081/health  
3. **minapp-frontend**: http://localhost:3000
4. **База данных**: localhost:5432 (fmp_db)
5. **Redis**: localhost:6379

## Проблемы и решения

### fmp-core не подключается к базе данных
- Убедитесь, что PostgreSQL запущен: `docker ps`
- Проверьте переменные окружения в .env файле
- Проверьте логи: `docker logs fmp-postgres`

### minapp-frontend ошибки сборки
- Убедитесь, что установлены все зависимости: `npm install`
- Проверьте версии Node.js и npm
- Очистите кэш: `npm cache clean --force`

### fmp-analytics проблемы сборки
- Проблема с ajv решена в Dockerfile
- Если проблемы продолжаются, запустите локально: `npm install && npm run build`

## Docker Compose (полный запуск)

```bash
# Запуск всех сервисов
./docker.sh up

# Проверка статуса
./docker.sh status

# Просмотр логов
./docker.sh logs

# Остановка
./docker.sh down
```

## Структура проекта

```
FMP/
├── fmp-core/           # Основной API сервер (Go)
├── minapp/             # Telegram Mini App
│   ├── backend/        # Backend для мини-приложения (Go)
│   └── frontend/       # Frontend мини-приложения (React)
├── fmp-analytics/      # Аналитический дашборд (React)
├── deploy/             # Ansible playbooks для деплоя
├── docker-compose.yml  # Docker Compose конфигурация
└── .env               # Переменные окружения
```

## API Endpoints

### fmp-core (порт 8080)
- `GET /health` - проверка здоровья
- `GET /api/v1/categories` - список категорий
- `POST /api/v1/categories` - создание категории
- `GET /api/v1/transactions` - список транзакций
- `POST /api/v1/transactions` - создание транзакции
- `GET /swagger/*` - Swagger документация

### minapp-backend (порт 8081)
- `GET /health` - проверка здоровья
- `POST /api/webhook` - Telegram webhook
- `GET /api/categories` - прокси к fmp-core
- `POST /api/transactions` - прокси к fmp-core

## Переменные окружения

Основные переменные в .env файле:

```bash
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fmp_db
DB_USER=fmp_user
DB_PASSWORD=fmp_password_secure_123

# API
API_PORT=8080
FMP_CORE_URL=http://localhost:8080/api/v1

# Telegram
TELEGRAM_BOT_TOKEN=test-token
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook

# Frontend
REACT_APP_API_URL=http://localhost:8080/api/v1
REACT_APP_MINAPP_BACKEND_URL=http://localhost:8081
```
