package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

func main() {
	fmt.Println("⚠️  ВНИМАНИЕ: Этот скрипт удалит ВСЕ расписания из базы данных!")
	fmt.Println("Нажмите Enter для продолжения или Ctrl+C для отмены...")
	fmt.Scanln()

	fmt.Println("Получение списка всех расписаний...")

	// Получаем все расписания
	schedules, err := getAllSchedules()
	if err != nil {
		fmt.Printf("❌ Ошибка получения расписаний: %v\n", err)
		return
	}

	if len(schedules) == 0 {
		fmt.Println("✅ Расписания не найдены. Нечего удалять.")
		return
	}

	fmt.Printf("📋 Найдено %d расписаний для удаления\n", len(schedules))

	// Удаляем каждое расписание
	successCount := 0
	errorCount := 0

	for i, schedule := range schedules {
		fmt.Printf("🗑️  Удаление расписания %d/%d (ID: %s)... ", i+1, len(schedules), schedule["id"])

		err := deleteSchedule(schedule["id"].(string))
		if err != nil {
			fmt.Printf("❌ Ошибка: %v\n", err)
			errorCount++
		} else {
			fmt.Println("✅ Успешно")
			successCount++
		}

		// Небольшая задержка между запросами
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n📊 Результат удаления:\n")
	fmt.Printf("✅ Успешно удалено: %d\n", successCount)
	fmt.Printf("❌ Ошибок: %d\n", errorCount)

	if errorCount == 0 {
		fmt.Println("🎉 Все расписания успешно удалены!")
	} else {
		fmt.Println("⚠️  Некоторые расписания не удалось удалить. Проверьте логи выше.")
	}
}

func getAllSchedules() ([]map[string]interface{}, error) {
	url := baseURL + "/schedules"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var schedules []map[string]interface{}
	err = json.Unmarshal(body, &schedules)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func deleteSchedule(id string) error {
	url := baseURL + "/schedules/" + id

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
