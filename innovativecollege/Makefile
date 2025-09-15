.PHONY: run build init-data clean

# Запуск приложения
run:
	go run main.go

# Сборка приложения
build:
	go build -o bin/innovativecollege main.go

# Инициализация тестовых данных
init-data:
	go run scripts/init_data.go

# Очистка собранных файлов
clean:
	rm -rf bin/

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Запуск с инициализацией данных
start: run init-data

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  run        - Запустить приложение"
	@echo "  build      - Собрать приложение"
	@echo "  init-data  - Инициализировать тестовые данные"
	@echo "  clean      - Очистить собранные файлы"
	@echo "  deps       - Установить зависимости"
	@echo "  start      - Запустить с инициализацией данных"
	@echo "  help       - Показать эту справку"

