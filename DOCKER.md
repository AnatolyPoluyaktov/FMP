# 🐳 FMP Docker Setup

Полная локальная среда разработки для Financial Manager Platform.

## 🚀 Быстрый старт

```bash
# Запустить всю платформу
./docker.sh up

# Или через docker-compose
docker-compose up -d
```

## 📋 Доступные сервисы

| Сервис | URL | Описание |
|--------|-----|----------|
| **Mini App Frontend** | http://localhost:3000 | Telegram Mini App |
| **Analytics Dashboard** | http://localhost:3001 | Аналитический дашборд |
| **FMP Core API** | http://localhost:8080 | Основной API |
| **Mini App Backend** | http://localhost:8081 | Backend для Mini App |
| **PostgreSQL** | localhost:5432 | База данных |
| **Redis** | localhost:6379 | Кэш и сессии |

## 🛠️ Управление

### Основные команды

```bash
# Запуск
./docker.sh up

# Остановка
./docker.sh down

# Перезапуск
./docker.sh restart

# Просмотр логов
./docker.sh logs
./docker.sh logs fmp-core  # для конкретного сервиса

# Сборка образов
./docker.sh build

# Очистка
./docker.sh clean

# Статус сервисов
./docker.sh status
```

### Отладка

```bash
# Подключение к контейнеру
./docker.sh shell fmp-core
./docker.sh shell minapp-backend

# Подключение к базе данных
./docker.sh db

# Подключение к Redis
./docker.sh redis
```

## 🔧 Конфигурация

### Переменные окружения

Скопируйте `env.example` в `.env` и настройте под свои нужды:

```bash
cp env.example .env
```

### Порты

- **3000** - Mini App Frontend
- **3001** - Analytics Dashboard  
- **8080** - FMP Core API
- **8081** - Mini App Backend
- **5432** - PostgreSQL
- **6379** - Redis

## 📊 Мониторинг

### Логи сервисов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f fmp-core
```

### Статус контейнеров

```bash
docker-compose ps
```

### Использование ресурсов

```bash
docker stats
```

## 🗄️ База данных

### Подключение

```bash
# Через скрипт
./docker.sh db

# Напрямую
docker-compose exec postgres psql -U fmp_user -d fmp_db
```

### Миграции

Миграции автоматически применяются при запуске контейнера PostgreSQL.

### Резервное копирование

```bash
# Создать бэкап
docker-compose exec postgres pg_dump -U fmp_user fmp_db > backup.sql

# Восстановить из бэкапа
docker-compose exec -T postgres psql -U fmp_user fmp_db < backup.sql
```

## 🔄 Разработка

### Hot Reload

Для разработки с hot reload используйте volumes в docker-compose.yml:

```yaml
volumes:
  - ./fmp-core:/app
  - ./minapp/frontend:/app
```

### Отладка

```bash
# Запуск в debug режиме
docker-compose -f docker-compose.yml -f docker-compose.debug.yml up

# Просмотр логов в реальном времени
docker-compose logs -f --tail=100
```

## 🚨 Troubleshooting

### Проблемы с портами

```bash
# Проверить занятые порты
lsof -i :3000
lsof -i :8080

# Остановить конфликтующие процессы
sudo kill -9 <PID>
```

### Проблемы с Docker

```bash
# Очистить все
docker system prune -a

# Пересобрать образы
docker-compose build --no-cache
```

### Проблемы с базой данных

```bash
# Пересоздать volume
docker-compose down -v
docker-compose up -d postgres
```

## 📝 Полезные команды

```bash
# Просмотр всех контейнеров
docker ps -a

# Просмотр образов
docker images

# Очистка неиспользуемых ресурсов
docker system prune

# Просмотр использования диска
docker system df
```
