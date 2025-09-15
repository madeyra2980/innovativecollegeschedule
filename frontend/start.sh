#!/bin/bash

# Скрипт для запуска фронтенда
echo "🚀 Запуск фронтенда календаря..."

# Проверяем, установлен ли Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js не установлен. Пожалуйста, установите Node.js"
    exit 1
fi

# Проверяем, установлен ли npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm не установлен. Пожалуйста, установите npm"
    exit 1
fi

# Устанавливаем зависимости если нужно
if [ ! -d "node_modules" ]; then
    echo "📦 Установка зависимостей..."
    npm install
fi

# Запускаем приложение
echo "🌟 Запуск React приложения..."
echo "📅 Календарь будет доступен по адресу: http://localhost:3000"
echo "🔗 Убедитесь, что бэкенд запущен на порту 8080"
echo ""

npm start