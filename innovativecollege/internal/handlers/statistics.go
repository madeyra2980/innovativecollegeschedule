package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetLessonStatistics получает статистику уроков
func (h *Handlers) GetLessonStatistics(c *gin.Context) {
	// Получаем параметры
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	groupID := c.Query("group_id")
	teacherID := c.Query("teacher_id")

	// Устанавливаем период по умолчанию (последние 30 дней)
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Парсим даты
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты окончания"})
		return
	}

	// Создаем фильтр
	filter := bson.M{
		"date": bson.M{
			"$gte": time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location()),
			"$lte": time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location()),
			"$ne":  nil, // Исключаем уроки без даты
		},
	}

	if groupID != "" {
		if id, err := primitive.ObjectIDFromHex(groupID); err == nil {
			filter["group_id"] = id
		}
	}

	if teacherID != "" {
		if id, err := primitive.ObjectIDFromHex(teacherID); err == nil {
			filter["teacher_id"] = id
		}
	}

	collection := h.db.Collection("lessons")

	// Общее количество уроков
	totalLessons, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подсчета уроков"})
		return
	}

	// Уроки по сменам
	firstShiftFilter := bson.M{}
	for k, v := range filter {
		firstShiftFilter[k] = v
	}
	firstShiftFilter["shift"] = 1
	firstShiftCount, _ := collection.CountDocuments(context.Background(), firstShiftFilter)

	secondShiftFilter := bson.M{}
	for k, v := range filter {
		secondShiftFilter[k] = v
	}
	secondShiftFilter["shift"] = 2
	secondShiftCount, _ := collection.CountDocuments(context.Background(), secondShiftFilter)

	// Уроки по дням недели
	dayStats := make(map[string]int)
	for i := 1; i <= 7; i++ {
		dayFilter := bson.M{}
		for k, v := range filter {
			dayFilter[k] = v
		}
		dayFilter["day_of_week"] = i
		count, _ := collection.CountDocuments(context.Background(), dayFilter)
		dayNames := []string{"", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
		dayStats[dayNames[i]] = int(count)
	}

	// Топ преподавателей по количеству уроков
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$teacher_id",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка агрегации"})
		return
	}
	defer cursor.Close(context.Background())

	var topTeachers []bson.M
	cursor.All(context.Background(), &topTeachers)

	// Топ групп по количеству уроков
	groupPipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$group_id",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	groupCursor, err := collection.Aggregate(context.Background(), groupPipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка агрегации групп"})
		return
	}
	defer groupCursor.Close(context.Background())

	var topGroups []bson.M
	groupCursor.All(context.Background(), &topGroups)

	// Заполняем имена преподавателей и групп
	for i, teacher := range topTeachers {
		teacherID := teacher["_id"].(primitive.ObjectID)
		teacherCollection := h.db.Collection("teachers")
		var teacherDoc bson.M
		if err := teacherCollection.FindOne(context.Background(), bson.M{"_id": teacherID}).Decode(&teacherDoc); err == nil {
			topTeachers[i]["name"] = teacherDoc["first_name"].(string) + " " + teacherDoc["last_name"].(string)
		}
	}

	for i, group := range topGroups {
		groupID := group["_id"].(primitive.ObjectID)
		groupCollection := h.db.Collection("groups")
		var groupDoc bson.M
		if err := groupCollection.FindOne(context.Background(), bson.M{"_id": groupID}).Decode(&groupDoc); err == nil {
			topGroups[i]["name"] = groupDoc["name"].(string)
		}
	}

	statistics := gin.H{
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"total_lessons": totalLessons,
		"by_shift": gin.H{
			"first_shift":  firstShiftCount,
			"second_shift": secondShiftCount,
		},
		"by_day_of_week": dayStats,
		"top_teachers":   topTeachers,
		"top_groups":     topGroups,
	}

	c.JSON(http.StatusOK, statistics)
}
