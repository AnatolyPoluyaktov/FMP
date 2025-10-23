# 🚀 GitHub Deployment Guide

Полное руководство по деплою FMP платформы на GitHub с использованием секретов.

## 📋 Предварительные требования

1. **GitHub репозиторий** - создан и настроен
2. **Сервер для деплоя** - VPS или облачный сервер
3. **SSH доступ** к серверу
4. **Docker** установлен на сервере

## 🔐 Настройка секретов

### 1. Перейдите в настройки репозитория
```
Settings → Secrets and variables → Actions
```

### 2. Добавьте необходимые секреты

#### 🗄️ Database Secrets
```
DB_PASSWORD=your_secure_database_password
POSTGRES_PASSWORD=your_secure_database_password
```

#### 🔑 API Secrets
```
API_SECRET_KEY=your-32-character-secret-key-for-api
JWT_SECRET=your-32-character-secret-key-for-jwt
```

#### 🤖 Telegram Bot Secrets
```
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook
```

#### 🚀 Deployment Secrets
```
PROD_HOST=your-server-ip-or-domain
PROD_USER=your-ssh-username
PROD_SSH_KEY=your-private-ssh-key
```

## 🛠️ Настройка сервера

### 1. Установка Docker
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# CentOS/RHEL
yum install -y docker
systemctl start docker
systemctl enable docker
```

### 2. Создание пользователя для деплоя
```bash
# Создать пользователя
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy

# Настроить SSH ключи
sudo mkdir -p /home/deploy/.ssh
sudo chown deploy:deploy /home/deploy/.ssh
sudo chmod 700 /home/deploy/.ssh
```

### 3. Настройка директорий
```bash
sudo mkdir -p /opt/fmp
sudo chown deploy:deploy /opt/fmp
```

## 📁 Структура проекта на сервере

```
/opt/fmp/
├── docker-compose.yml
├── .env
├── fmp-core/
│   ├── Dockerfile
│   └── migrations/
├── minapp/
│   ├── backend/
│   │   └── Dockerfile
│   └── frontend/
│       └── Dockerfile
├── fmp-analytics/
│   └── Dockerfile
└── deploy/
    └── nginx/
        └── nginx.conf
```

## 🔄 Процесс деплоя

### 1. GitHub Actions автоматически:
- Запускает тесты
- Собирает приложения
- Деплоит на сервер при push в main

### 2. Ручной деплой
```bash
# На сервере
cd /opt/fmp
git pull origin main
docker-compose --env-file .env up -d --build
```

## 🌐 Настройка домена и SSL

### 1. Настройка Nginx
```bash
# Создать конфигурацию сайта
sudo nano /etc/nginx/sites-available/fmp

# Содержимое:
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /api/ {
        proxy_pass http://localhost:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. Активация сайта
```bash
sudo ln -s /etc/nginx/sites-available/fmp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. SSL сертификат (Let's Encrypt)
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## 📊 Мониторинг

### 1. Логи приложений
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f fmp-core
```

### 2. Статус сервисов
```bash
docker-compose ps
docker stats
```

### 3. Мониторинг ресурсов
```bash
# Использование диска
df -h

# Использование памяти
free -h

# Загрузка системы
htop
```

## 🔧 Troubleshooting

### Проблема: Сервисы не запускаются
```bash
# Проверить логи
docker-compose logs

# Проверить конфигурацию
docker-compose config

# Пересобрать образы
docker-compose build --no-cache
```

### Проблема: База данных недоступна
```bash
# Проверить подключение
docker-compose exec postgres psql -U fmp_user -d fmp_db

# Проверить миграции
docker-compose exec fmp-core ./migrate -path migrations -database "postgres://fmp_user:fmp_password@postgres:5432/fmp_db?sslmode=disable" up
```

### Проблема: SSL сертификат
```bash
# Обновить сертификат
sudo certbot renew

# Проверить статус
sudo certbot certificates
```

## 🔄 Обновление

### Автоматическое обновление
- При push в main ветку автоматически запускается деплой

### Ручное обновление
```bash
# На сервере
cd /opt/fmp
git pull origin main
docker-compose --env-file .env up -d --build
```

## 📝 Полезные команды

```bash
# Перезапуск сервисов
docker-compose restart

# Остановка всех сервисов
docker-compose down

# Очистка неиспользуемых ресурсов
docker system prune -a

# Бэкап базы данных
docker-compose exec postgres pg_dump -U fmp_user fmp_db > backup.sql

# Восстановление из бэкапа
docker-compose exec -T postgres psql -U fmp_user fmp_db < backup.sql
```

## 🚨 Безопасность

### Рекомендации:
- Используйте сложные пароли
- Регулярно обновляйте зависимости
- Настройте файрвол
- Используйте HTTPS
- Регулярно делайте бэкапы
- Мониторьте логи на предмет подозрительной активности

### Файрвол (UFW)
```bash
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw enable
```
