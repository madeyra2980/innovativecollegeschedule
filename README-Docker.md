# 🐳 Docker Setup для Innovative College

Этот документ описывает, как запустить всю систему Innovative College с помощью Docker.

## 📋 Предварительные требования

- Docker 20.10+
- Docker Compose 2.0+

## 🚀 Быстрый запуск

### 1. Клонируйте репозиторий и перейдите в директорию проекта
```bash
cd /path/to/your/project
```

### 2. Запустите всю систему одной командой
```bash
docker-compose up -d
```

### 3. Инициализируйте тестовые данные (опционально)
```bash
# Подождите пока backend запустится, затем выполните:
docker-compose exec backend wget -O- http://localhost:8080/health

# Инициализируйте тестовые данные
docker-compose exec backend go run scripts/init_data.go
```

## 🌐 Доступ к приложениям

После запуска приложения будут доступны по следующим адресам:

- **Backend API**: http://localhost:8080
- **Admin Panel**: http://localhost:3000
- **Schedule Viewer**: http://localhost:3001
- **MongoDB**: localhost:27017

## 📊 Компоненты системы

### 1. MongoDB (Порт 27017)
- База данных для хранения всех данных системы
- Автоматически создается с пользователем admin/password123

### 2. Backend API (Порт 8080)
- Go приложение с REST API
- Автоматически подключается к MongoDB
- Health check на `/health`

### 3. Frontend Admin Panel (Порт 3000)
- React приложение для администраторов
- Nginx для раздачи статических файлов
- Автоматически проксирует API запросы на backend

### 4. Schedule Viewer (Порт 3001)
- React TypeScript приложение для просмотра расписания
- Современный UI для студентов и преподавателей
- Nginx для раздачи статических файлов

## 🛠️ Управление контейнерами

### Просмотр статуса
```bash
docker-compose ps
```

### Просмотр логов
```bash
# Все сервисы
docker-compose logs

# Конкретный сервис
docker-compose logs backend
docker-compose logs frontend
docker-compose logs schedule-app
docker-compose logs mongodb
```

### Остановка системы
```bash
docker-compose down
```

### Остановка с удалением данных
```bash
docker-compose down -v
```

### Пересборка контейнеров
```bash
docker-compose build --no-cache
docker-compose up -d
```

## 🔧 Разработка

### Запуск в режиме разработки
```bash
# Запустите только базу данных
docker-compose up -d mongodb

# Запустите backend локально
cd innovativecollege
go run main.go

# Запустите frontend локально
cd frontend
npm start

# Запустите schedule-app локально
cd schedule-app
npm start
```

### Обновление кода
```bash
# После изменения кода пересоберите и перезапустите
docker-compose build
docker-compose up -d
```

## 📁 Структура Docker файлов

```
app/
├── docker-compose.yml          # Главный файл оркестрации
├── .dockerignore              # Игнорируемые файлы
├── README-Docker.md           # Эта документация
├── innovativecollege/
│   └── Dockerfile             # Backend контейнер
├── frontend/
│   ├── Dockerfile             # Frontend контейнер
│   └── nginx.conf             # Nginx конфигурация
└── schedule-app/
    ├── Dockerfile             # Schedule app контейнер
    └── nginx.conf             # Nginx конфигурация
```

## 🔍 Отладка

### Проверка подключения к базе данных
```bash
docker-compose exec backend ping mongodb
```

### Проверка API
```bash
curl http://localhost:8080/health
```

### Доступ к MongoDB
```bash
docker-compose exec mongodb mongosh -u admin -p password123 --authenticationDatabase admin
```

### Просмотр переменных окружения
```bash
docker-compose exec backend env
```

## 🚨 Устранение проблем

### Порт уже используется
Если порты 3000, 3001, 8080 или 27017 уже используются:
```bash
# Измените порты в docker-compose.yml
ports:
  - "3002:80"  # Вместо 3000:80
```

### Проблемы с правами доступа
```bash
# На Linux/Mac
sudo docker-compose up -d
```

### Очистка Docker кэша
```bash
docker system prune -a
docker volume prune
```

## 📈 Мониторинг

### Использование ресурсов
```bash
docker stats
```

### Проверка здоровья сервисов
```bash
docker-compose ps
```

## 🔐 Безопасность

В production окружении обязательно измените:
- Пароли MongoDB
- Переменные окружения
- Настройки сети
- SSL сертификаты

## 📞 Поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose logs`
2. Убедитесь, что все порты свободны
3. Проверьте подключение к интернету для загрузки образов
4. Перезапустите систему: `docker-compose restart`


