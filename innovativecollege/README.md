# Система расписания колледжа

Простое backend приложение на Golang + MongoDB для управления расписанием, группами, студентами и преподавателями.

## Особенности

- **Без авторизации** в админ панели - все операции доступны напрямую
- **Вход по ИИН** - студенты и преподаватели заходят в свой профиль только по ИИН (без пароля)
- **Простая логика** - минимальный функционал без усложнений

## Требования

- Go 1.21+
- MongoDB 4.4+

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd innovativecollege
```

2. Установите зависимости:
```bash
go mod tidy
```

3. Убедитесь, что MongoDB запущен на `localhost:27017`

4. Запустите приложение:
```bash
go run main.go
```

Сервер запустится на порту 8080.

## API Endpoints

### Группы
- `POST /api/v1/groups` - Создать группу
- `GET /api/v1/groups` - Получить все группы
- `PUT /api/v1/groups/{id}` - Обновить группу
- `DELETE /api/v1/groups/{id}` - Удалить группу

### Студенты
- `POST /api/v1/students` - Создать студента
- `GET /api/v1/students` - Получить всех студентов
- `GET /api/v1/students/{iin}/schedule` - Получить расписание студента по ИИН
- `PUT /api/v1/students/{id}` - Обновить студента
- `DELETE /api/v1/students/{id}` - Удалить студента

### Преподаватели
- `POST /api/v1/teachers` - Создать преподавателя
- `GET /api/v1/teachers` - Получить всех преподавателей
- `GET /api/v1/teachers/{iin}/schedule` - Получить расписание преподавателя по ИИН
- `PUT /api/v1/teachers/{id}` - Обновить преподавателя
- `DELETE /api/v1/teachers/{id}` - Удалить преподавателя

### Расписание
- `POST /api/v1/schedules` - Создать расписание
- `GET /api/v1/schedules` - Получить все расписания
- `GET /api/v1/schedules/day/{day}` - Получить расписание по дню недели (1-7)
- `PUT /api/v1/schedules/{id}` - Обновить расписание
- `DELETE /api/v1/schedules/{id}` - Удалить расписание

### Health Check
- `GET /health` - Проверка состояния сервера

## Примеры использования

### Создание группы
```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Программное обеспечение",
    "code": "ПО-31",
    "description": "Группа по программированию"
  }'
```

### Создание студента
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{
    "iin": "123456789012",
    "first_name": "Айдар",
    "last_name": "Караев",
    "group_id": "GROUP_ID_HERE"
  }'
```

### Создание преподавателя
```bash
curl -X POST http://localhost:8080/api/v1/teachers \
  -H "Content-Type: application/json" \
  -d '{
    "iin": "987654321098",
    "first_name": "Мария",
    "last_name": "Иванова",
    "subjects": ["Программирование", "Базы данных"]
  }'
```

### Создание расписания
```bash
curl -X POST http://localhost:8080/api/v1/schedules \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": "GROUP_ID_HERE",
    "teacher_id": "TEACHER_ID_HERE",
    "subject": "Программирование",
    "room": "508",
    "day_of_week": 1,
    "start_time": "12:40",
    "end_time": "14:00",
    "shift": 2,
    "description": "Анализ и оценка экономических процессов"
  }'
```

### Получение расписания студента
```bash
curl http://localhost:8080/api/v1/students/123456789012/schedule
```

### Получение расписания преподавателя
```bash
curl http://localhost:8080/api/v1/teachers/987654321098/schedule
```

## Структура проекта

```
innovativecollege/
├── main.go                 # Точка входа
├── go.mod                  # Зависимости Go
├── config.env             # Конфигурация
├── internal/
│   ├── config/            # Конфигурация приложения
│   ├── database/          # Подключение к MongoDB
│   ├── handlers/          # HTTP обработчики
│   ├── models/            # Модели данных
│   └── routes/            # Маршруты API
└── README.md              # Документация
```

## Модели данных

### Group (Группа)
- ID, Name, Code, Description, CreatedAt, UpdatedAt

### Student (Студент)
- ID, IIN, FirstName, LastName, GroupID, CreatedAt, UpdatedAt

### Teacher (Преподаватель)
- ID, IIN, FirstName, LastName, Subjects[], CreatedAt, UpdatedAt

### Schedule (Расписание)
- ID, GroupID, TeacherID, Subject, Room, DayOfWeek, StartTime, EndTime, Shift, Description, CreatedAt, UpdatedAt

## Примечания

- День недели: 1 = понедельник, 7 = воскресенье
- Смена: 1 = первая смена, 2 = вторая смена
- ИИН используется как уникальный идентификатор для входа студентов и преподавателей
- Все времена хранятся в формате "HH:MM"

