# Примеры использования API

## Быстрый старт

1. Запустите MongoDB
2. Запустите приложение: `go run main.go`
3. Инициализируйте тестовые данные: `go run scripts/init_data.go`
4. Откройте `test.html` в браузере для тестирования

## Основные API endpoints

### Группы

#### Создать группу
```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Программное обеспечение",
    "code": "ПО-31",
    "description": "Группа по программированию"
  }'
```

#### Получить все группы
```bash
curl http://localhost:8080/api/v1/groups
```

### Студенты

#### Создать студента
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{
    "iin": "111111111111",
    "first_name": "Алма",
    "last_name": "Нурланова",
    "group_id": "GROUP_ID_HERE"
  }'
```

#### Получить всех студентов
```bash
curl http://localhost:8080/api/v1/students
```

#### Получить расписание студента по ИИН
```bash
curl http://localhost:8080/api/v1/students/111111111111/schedule
```

### Преподаватели

#### Создать преподавателя
```bash
curl -X POST http://localhost:8080/api/v1/teachers \
  -H "Content-Type: application/json" \
  -d '{
    "iin": "123456789012",
    "first_name": "Айдар",
    "last_name": "Караев",
    "subjects": ["Программирование", "Базы данных"]
  }'
```

#### Получить всех преподавателей
```bash
curl http://localhost:8080/api/v1/teachers
```

#### Получить расписание преподавателя по ИИН
```bash
curl http://localhost:8080/api/v1/teachers/123456789012/schedule
```

### Расписание

#### Создать расписание
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
    "description": "Основы программирования"
  }'
```

#### Получить все расписания
```bash
curl http://localhost:8080/api/v1/schedules
```

#### Получить расписание по дню недели
```bash
curl http://localhost:8080/api/v1/schedules/day/1
```

## Обновление записей

### Обновить группу
```bash
curl -X PUT http://localhost:8080/api/v1/groups/GROUP_ID_HERE \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Новое название группы",
    "description": "Обновленное описание"
  }'
```

### Обновить студента
```bash
curl -X PUT http://localhost:8080/api/v1/students/STUDENT_ID_HERE \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Новое имя",
    "last_name": "Новая фамилия"
  }'
```

### Обновить преподавателя
```bash
curl -X PUT http://localhost:8080/api/v1/teachers/TEACHER_ID_HERE \
  -H "Content-Type: application/json" \
  -d '{
    "subjects": ["Новый предмет", "Другой предмет"]
  }'
```

### Обновить расписание
```bash
curl -X PUT http://localhost:8080/api/v1/schedules/SCHEDULE_ID_HERE \
  -H "Content-Type: application/json" \
  -d '{
    "room": "Новая аудитория",
    "start_time": "14:00",
    "end_time": "15:20"
  }'
```

## Удаление записей

### Удалить группу
```bash
curl -X DELETE http://localhost:8080/api/v1/groups/GROUP_ID_HERE
```

### Удалить студента
```bash
curl -X DELETE http://localhost:8080/api/v1/students/STUDENT_ID_HERE
```

### Удалить преподавателя
```bash
curl -X DELETE http://localhost:8080/api/v1/teachers/TEACHER_ID_HERE
```

### Удалить расписание
```bash
curl -X DELETE http://localhost:8080/api/v1/schedules/SCHEDULE_ID_HERE
```

## Правила удаления

- **Группу** можно удалить только если в ней нет студентов
- **Преподавателя** можно удалить только если у него нет расписаний
- **Студента** можно удалить в любое время
- **Расписание** можно удалить в любое время

## Частичное обновление

Все UPDATE операции поддерживают частичное обновление - можно передать только те поля, которые нужно изменить. Остальные поля останутся без изменений.

Например, чтобы изменить только название группы:
```bash
curl -X PUT http://localhost:8080/api/v1/groups/GROUP_ID_HERE \
  -H "Content-Type: application/json" \
  -d '{"name": "Только новое название"}'
```

## Тестовые данные

После запуска `go run scripts/init_data.go` будут созданы:

### Группы:
- ПО-31 (Программное обеспечение)
- ӨҚ-31 (Экономика)
- Құқық-31 (Право)
- ЖҚҰ-31 (Дорожная полиция)

### Студенты:
- ИИН: 111111111111 (Алма Нурланова)
- ИИН: 222222222222 (Данияр Токтаров)
- ИИН: 333333333333 (Айжан Калиева)
- ИИН: 444444444444 (Ерлан Жумабаев)

### Преподаватели:
- ИИН: 123456789012 (Айдар Караев)
- ИИН: 234567890123 (Мария Иванова)
- ИИН: 345678901234 (Сергей Петров)
- ИИН: 456789012345 (Анна Сидорова)

## Структура ответов

### Расписание студента
```json
{
  "student": {
    "id": "...",
    "iin": "111111111111",
    "first_name": "Алма",
    "last_name": "Нурланова",
    "group_id": "...",
    "created_at": "...",
    "updated_at": "..."
  },
  "schedules": [
    {
      "id": "...",
      "group_id": "...",
      "teacher_id": "...",
      "teacher": {
        "id": "...",
        "iin": "123456789012",
        "first_name": "Айдар",
        "last_name": "Караев",
        "subjects": ["Программирование", "Базы данных"]
      },
      "subject": "Анализ и оценка экономических процессов",
      "room": "508",
      "day_of_week": 1,
      "start_time": "12:40",
      "end_time": "14:00",
      "shift": 2,
      "description": "...",
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

## Коды дней недели

- 1 = Понедельник
- 2 = Вторник
- 3 = Среда
- 4 = Четверг
- 5 = Пятница
- 6 = Суббота
- 7 = Воскресенье

## Смены

- 1 = Первая смена
- 2 = Вторая смена
