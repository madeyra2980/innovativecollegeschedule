package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"innovativecollege/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateTimeSlot создает новый временной слот
func (h *Handlers) CreateTimeSlot(c *gin.Context) {
	var req models.CreateTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерируем label если не указан
	if req.Label == "" {
		req.Label = req.StartTime + "-" + req.EndTime
	}

	timeSlot := models.TimeSlot{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Shift:     req.Shift,
		Label:     req.Label,
		IsActive:  req.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := h.db.Collection("time_slots")
	result, err := collection.InsertOne(context.Background(), timeSlot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания временного слота"})
		return
	}

	timeSlot.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, timeSlot)
}

// GetTimeSlots получает все временные слоты
func (h *Handlers) GetTimeSlots(c *gin.Context) {
	collection := h.db.Collection("time_slots")

	// Получаем параметры фильтрации
	shift := c.Query("shift")
	isActive := c.Query("is_active")

	filter := bson.M{}
	if shift != "" {
		if shiftInt, err := strconv.Atoi(shift); err == nil {
			filter["shift"] = shiftInt
		}
	}
	if isActive != "" {
		if isActiveBool, err := strconv.ParseBool(isActive); err == nil {
			filter["is_active"] = isActiveBool
		}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения временных слотов"})
		return
	}
	defer cursor.Close(context.Background())

	var timeSlots []models.TimeSlot
	if err = cursor.All(context.Background(), &timeSlots); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки временных слотов"})
		return
	}

	c.JSON(http.StatusOK, timeSlots)
}

// GetTimeSlot получает временной слот по ID
func (h *Handlers) GetTimeSlot(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	collection := h.db.Collection("time_slots")
	var timeSlot models.TimeSlot
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&timeSlot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Временной слот не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска временного слота"})
		}
		return
	}

	c.JSON(http.StatusOK, timeSlot)
}

// UpdateTimeSlot обновляет временной слот
func (h *Handlers) UpdateTimeSlot(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	var req models.UpdateTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.Collection("time_slots")

	// Проверяем существование временного слота
	var existingTimeSlot models.TimeSlot
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&existingTimeSlot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Временной слот не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска временного слота"})
		}
		return
	}

	// Строим обновление
	update := bson.M{"updated_at": time.Now()}
	if req.StartTime != nil {
		update["start_time"] = *req.StartTime
	}
	if req.EndTime != nil {
		update["end_time"] = *req.EndTime
	}
	if req.Shift != nil {
		update["shift"] = *req.Shift
	}
	if req.Label != nil {
		update["label"] = *req.Label
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}

	// Если изменилось время, обновляем label
	if req.StartTime != nil || req.EndTime != nil {
		startTime := existingTimeSlot.StartTime
		endTime := existingTimeSlot.EndTime
		if req.StartTime != nil {
			startTime = *req.StartTime
		}
		if req.EndTime != nil {
			endTime = *req.EndTime
		}
		update["label"] = startTime + "-" + endTime
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления временного слота"})
		return
	}

	// Получаем обновленный временной слот
	var updatedTimeSlot models.TimeSlot
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&updatedTimeSlot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного временного слота"})
		return
	}

	c.JSON(http.StatusOK, updatedTimeSlot)
}

// DeleteTimeSlot удаляет временной слот
func (h *Handlers) DeleteTimeSlot(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	collection := h.db.Collection("time_slots")

	// Проверяем существование временного слота
	var existingTimeSlot models.TimeSlot
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&existingTimeSlot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Временной слот не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска временного слота"})
		}
		return
	}

	// Проверяем, есть ли уроки в этом временном слоте
	lessonCollection := h.db.Collection("lessons")
	lessonCount, err := lessonCollection.CountDocuments(context.Background(), bson.M{
		"start_time": existingTimeSlot.StartTime,
		"end_time":   existingTimeSlot.EndTime,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки связанных уроков"})
		return
	}

	if lessonCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя удалить временной слот, в котором есть уроки"})
		return
	}

	// Удаляем временной слот
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления временного слота"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Временной слот успешно удален"})
}
