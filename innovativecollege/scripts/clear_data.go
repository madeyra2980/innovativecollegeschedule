package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

func main() {
	fmt.Println("Очистка тестовых данных...")

	// Ждем запуска сервера
	time.Sleep(2 * time.Second)

	// Очищаем все данные
	clearLessons()
	clearSchedules()
	clearStudents()
	clearTeachers()
	clearGroups()
	clearSubjects()

	fmt.Println("\nТестовые данные успешно очищены!")
}

func makeRequest(method, url string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func clearLessons() {
	fmt.Println("Очистка уроков...")
	// Получаем все уроки
	resp, err := makeRequest("GET", baseURL+"/lessons", nil)
	if err != nil {
		fmt.Printf("Ошибка получения уроков: %v\n", err)
		return
	}

	var lessons []map[string]interface{}
	json.Unmarshal(resp, &lessons)

	// Удаляем каждый урок
	for _, lesson := range lessons {
		if id, ok := lesson["id"].(string); ok {
			_, err := makeRequest("DELETE", baseURL+"/lessons/"+id, nil)
			if err != nil {
				fmt.Printf("Ошибка удаления урока %s: %v\n", id, err)
			} else {
				fmt.Printf("Удален урок: %s\n", id)
			}
		}
	}
}

func clearSchedules() {
	fmt.Println("Очистка расписания...")
	// Получаем все расписания
	resp, err := makeRequest("GET", baseURL+"/schedules", nil)
	if err != nil {
		fmt.Printf("Ошибка получения расписаний: %v\n", err)
		return
	}

	var schedules []map[string]interface{}
	json.Unmarshal(resp, &schedules)

	// Удаляем каждое расписание
	for _, schedule := range schedules {
		if id, ok := schedule["id"].(string); ok {
			_, err := makeRequest("DELETE", baseURL+"/schedules/"+id, nil)
			if err != nil {
				fmt.Printf("Ошибка удаления расписания %s: %v\n", id, err)
			} else {
				fmt.Printf("Удалено расписание: %s\n", id)
			}
		}
	}
}

func clearStudents() {
	fmt.Println("Очистка студентов...")
	// Получаем всех студентов
	resp, err := makeRequest("GET", baseURL+"/students", nil)
	if err != nil {
		fmt.Printf("Ошибка получения студентов: %v\n", err)
		return
	}

	var students []map[string]interface{}
	json.Unmarshal(resp, &students)

	// Удаляем каждого студента
	for _, student := range students {
		if iin, ok := student["iin"].(string); ok {
			_, err := makeRequest("DELETE", baseURL+"/students/"+iin, nil)
			if err != nil {
				fmt.Printf("Ошибка удаления студента %s: %v\n", iin, err)
			} else {
				fmt.Printf("Удален студент: %s\n", iin)
			}
		}
	}
}

func clearTeachers() {
	fmt.Println("Очистка преподавателей...")
	// Получаем всех преподавателей
	resp, err := makeRequest("GET", baseURL+"/teachers", nil)
	if err != nil {
		fmt.Printf("Ошибка получения преподавателей: %v\n", err)
		return
	}

	var teachers []map[string]interface{}
	json.Unmarshal(resp, &teachers)

	// Удаляем каждого преподавателя
	for _, teacher := range teachers {
		if iin, ok := teacher["iin"].(string); ok {
			_, err := makeRequest("DELETE", baseURL+"/teachers/"+iin, nil)
			if err != nil {
				fmt.Printf("Ошибка удаления преподавателя %s: %v\n", iin, err)
			} else {
				fmt.Printf("Удален преподаватель: %s\n", iin)
			}
		}
	}
}

func clearGroups() {
	fmt.Println("Очистка групп...")
	// Получаем все группы
	resp, err := makeRequest("GET", baseURL+"/groups", nil)
	if err != nil {
		fmt.Printf("Ошибка получения групп: %v\n", err)
		return
	}

	var groups []map[string]interface{}
	json.Unmarshal(resp, &groups)

	// Удаляем каждую группу
	for _, group := range groups {
		if id, ok := group["id"].(string); ok {
			_, err := makeRequest("DELETE", baseURL+"/groups/"+id, nil)
			if err != nil {
				fmt.Printf("Ошибка удаления группы %s: %v\n", id, err)
			} else {
				fmt.Printf("Удалена группа: %s\n", id)
			}
		}
	}
}

func clearSubjects() {
	fmt.Println("Очистка предметов...")
	// Получаем все предметы
	resp, err := makeRequest("GET", baseURL+"/subjects", nil)
	if err != nil {
		fmt.Printf("Ошибка получения предметов: %v\n", err)
		return
	}

	var subjects []map[string]interface{}
	json.Unmarshal(resp, &subjects)

	// Удаляем каждый предмет
	for _, subject := range subjects {
		if id, ok := subject["id"].(string); ok {
			_, err := makeRequest("DELETE", baseURL+"/subjects/"+id, nil)
			if err != nil {
				fmt.Printf("Ошибка удаления предмета %s: %v\n", id, err)
			} else {
				fmt.Printf("Удален предмет: %s\n", id)
			}
		}
	}
}
