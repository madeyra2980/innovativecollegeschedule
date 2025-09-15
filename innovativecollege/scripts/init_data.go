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
	fmt.Println("Инициализация тестовых данных...")

	// Ждем запуска сервера
	time.Sleep(2 * time.Second)

	// Создаем группы
	groups := []map[string]string{
		{"name": "Программное обеспечение", "code": "ПО-31", "description": "Группа по программированию"},
		{"name": "Экономика", "code": "ӨҚ-31", "description": "Экономическая группа"},
		{"name": "Право", "code": "Құқық-31", "description": "Юридическая группа"},
		{"name": "Дорожная полиция", "code": "ЖҚҰ-31", "description": "Группа дорожной полиции"},
	}

	var groupIDs []string
	for _, group := range groups {
		groupID := createGroup(group)
		if groupID != "" {
			groupIDs = append(groupIDs, groupID)
			fmt.Printf("Создана группа: %s (ID: %s)\n", group["name"], groupID)
		}
	}

	// Создаем преподавателей
	teachers := []map[string]interface{}{
		{"iin": "123456789012", "first_name": "Айдар", "last_name": "Караев", "subjects": []string{"Программирование", "Базы данных"}},
		{"iin": "234567890123", "first_name": "Мария", "last_name": "Иванова", "subjects": []string{"Экономика", "Маркетинг"}},
		{"iin": "345678901234", "first_name": "Сергей", "last_name": "Петров", "subjects": []string{"Право", "Гражданское право"}},
		{"iin": "456789012345", "first_name": "Анна", "last_name": "Сидорова", "subjects": []string{"Дорожная безопасность", "ПДД"}},
	}

	var teacherIDs []string
	for _, teacher := range teachers {
		teacherID := createTeacher(teacher)
		if teacherID != "" {
			teacherIDs = append(teacherIDs, teacherID)
			fmt.Printf("Создан преподаватель: %s %s (ID: %s)\n", teacher["first_name"], teacher["last_name"], teacherID)
		}
	}

	// Создаем студентов
	students := []map[string]string{
		{"iin": "111111111111", "first_name": "Алма", "last_name": "Нурланова", "group_id": groupIDs[0]},
		{"iin": "222222222222", "first_name": "Данияр", "last_name": "Токтаров", "group_id": groupIDs[0]},
		{"iin": "333333333333", "first_name": "Айжан", "last_name": "Калиева", "group_id": groupIDs[1]},
		{"iin": "444444444444", "first_name": "Ерлан", "last_name": "Жумабаев", "group_id": groupIDs[2]},
	}

	for _, student := range students {
		createStudent(student)
		fmt.Printf("Создан студент: %s %s\n", student["first_name"], student["last_name"])
	}

	// Создаем расписание (понедельник, 2 смена)
	schedules := []map[string]interface{}{
		{
			"group_id":    groupIDs[0],
			"teacher_id":  teacherIDs[0],
			"subject":     "Анализ и оценка экономических процессов в предприятии",
			"room":        "508",
			"day_of_week": 1,
			"start_time":  "12:40",
			"end_time":    "14:00",
			"shift":       2,
			"description": "Использование информационно-справочных и интерактивных веб-порталов",
		},
		{
			"group_id":    groupIDs[1],
			"teacher_id":  teacherIDs[1],
			"subject":     "Определение уровня пожарной опасности зданий и сооружений",
			"room":        "503",
			"day_of_week": 1,
			"start_time":  "12:40",
			"end_time":    "14:00",
			"shift":       2,
			"description": "Понятия о пожаре или чрезвычайной ситуации",
		},
		{
			"group_id":    groupIDs[2],
			"teacher_id":  teacherIDs[2],
			"subject":     "Правовые основы правоохранительных органов",
			"room":        "210",
			"day_of_week": 1,
			"start_time":  "8:00",
			"end_time":    "9:20",
			"shift":       1,
			"description": "Предоставление информации об их деятельности",
		},
	}

	for _, schedule := range schedules {
		createSchedule(schedule)
		fmt.Printf("Создано расписание: %s в %s\n", schedule["subject"], schedule["room"])
	}

	fmt.Println("\nТестовые данные успешно созданы!")
	fmt.Println("\nПримеры запросов:")
	fmt.Println("1. Расписание студента: GET http://localhost:8080/api/v1/students/111111111111/schedule")
	fmt.Println("2. Расписание преподавателя: GET http://localhost:8080/api/v1/teachers/123456789012/schedule")
	fmt.Println("3. Все группы: GET http://localhost:8080/api/v1/groups")
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

func createGroup(group map[string]string) string {
	url := baseURL + "/groups"
	resp, err := makeRequest("POST", url, group)
	if err != nil {
		fmt.Printf("Ошибка создания группы: %v\n", err)
		return ""
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)
	if id, ok := result["id"].(string); ok {
		return id
	}
	return ""
}

func createTeacher(teacher map[string]interface{}) string {
	url := baseURL + "/teachers"
	resp, err := makeRequest("POST", url, teacher)
	if err != nil {
		fmt.Printf("Ошибка создания преподавателя: %v\n", err)
		return ""
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)
	if id, ok := result["id"].(string); ok {
		return id
	}
	return ""
}

func createStudent(student map[string]string) string {
	url := baseURL + "/students"
	resp, err := makeRequest("POST", url, student)
	if err != nil {
		fmt.Printf("Ошибка создания студента: %v\n", err)
		return ""
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)
	if id, ok := result["id"].(string); ok {
		return id
	}
	return ""
}

func createSchedule(schedule map[string]interface{}) string {
	url := baseURL + "/schedules"
	resp, err := makeRequest("POST", url, schedule)
	if err != nil {
		fmt.Printf("Ошибка создания расписания: %v\n", err)
		return ""
	}

	var result map[string]interface{}
	json.Unmarshal(resp, &result)
	if id, ok := result["id"].(string); ok {
		return id
	}
	return ""
}
