package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"innovativecollege/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Handlers {
	return &Handlers{db: db}
}

// ========== ГРУППЫ ==========

// CreateGroup создает новую группу
func (h *Handlers) CreateGroup(c *gin.Context) {
	var req models.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group := models.Group{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := h.db.Collection("groups")
	result, err := collection.InsertOne(context.Background(), group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания группы"})
		return
	}

	group.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, group)
}

// GetGroups получает все группы
func (h *Handlers) GetGroups(c *gin.Context) {
	collection := h.db.Collection("groups")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения групп"})
		return
	}
	defer cursor.Close(context.Background())

	var groups []models.Group
	if err = cursor.All(context.Background(), &groups); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки групп"})
		return
	}

	// Если нет групп, возвращаем пустой массив вместо null
	if groups == nil {
		groups = []models.Group{}
	}

	c.JSON(http.StatusOK, groups)
}

// UpdateGroup обновляет группу
func (h *Handlers) UpdateGroup(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
		return
	}

	var req models.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование группы
	collection := h.db.Collection("groups")
	var existingGroup models.Group
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingGroup)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Группа не найдена"})
		return
	}

	// Создаем объект для обновления
	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}

	// Обновляем группу
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления группы"})
		return
	}

	// Получаем обновленную группу
	var updatedGroup models.Group
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленной группы"})
		return
	}

	c.JSON(http.StatusOK, updatedGroup)
}

// DeleteGroup удаляет группу
func (h *Handlers) DeleteGroup(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
		return
	}

	// Проверяем существование группы
	collection := h.db.Collection("groups")
	var group models.Group
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Группа не найдена"})
		return
	}

	// Проверяем, есть ли студенты в этой группе
	studentCollection := h.db.Collection("students")
	studentCount, err := studentCollection.CountDocuments(context.Background(), bson.M{"group_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки студентов"})
		return
	}

	// Проверяем, есть ли уроки с этой группой
	lessonCollection := h.db.Collection("lessons")
	lessonCount, err := lessonCollection.CountDocuments(context.Background(), bson.M{"group_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки уроков"})
		return
	}

	if studentCount > 0 || lessonCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя удалить группу, в которой есть студенты или уроки"})
		return
	}

	// Удаляем группу
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления группы"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Группа успешно удалена"})
}

// ========== ПРЕДМЕТЫ ==========

// CreateSubject создает новый предмет
func (h *Handlers) CreateSubject(c *gin.Context) {
	var req models.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject := models.Subject{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := h.db.Collection("subjects")
	result, err := collection.InsertOne(context.Background(), subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания предмета"})
		return
	}

	subject.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, subject)
}

// GetSubjects получает все предметы
func (h *Handlers) GetSubjects(c *gin.Context) {
	collection := h.db.Collection("subjects")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения предметов"})
		return
	}
	defer cursor.Close(context.Background())

	var subjects []models.Subject
	if err = cursor.All(context.Background(), &subjects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки предметов"})
		return
	}

	// Если нет предметов, возвращаем пустой массив вместо null
	if subjects == nil {
		subjects = []models.Subject{}
	}

	c.JSON(http.StatusOK, subjects)
}

// UpdateSubject обновляет предмет
func (h *Handlers) UpdateSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID предмета"})
		return
	}

	var req models.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование предмета
	collection := h.db.Collection("subjects")
	var existingSubject models.Subject
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingSubject)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предмет не найден"})
		return
	}

	// Создаем объект для обновления
	update := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Code != "" {
		update["code"] = req.Code
	}
	if req.Description != "" {
		update["description"] = req.Description
	}

	// Обновляем предмет
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления предмета"})
		return
	}

	// Получаем обновленный предмет
	var updatedSubject models.Subject
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedSubject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного предмета"})
		return
	}

	c.JSON(http.StatusOK, updatedSubject)
}

// DeleteSubject удаляет предмет
func (h *Handlers) DeleteSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID предмета"})
		return
	}

	// Проверяем существование предмета
	collection := h.db.Collection("subjects")
	var existingSubject models.Subject
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingSubject)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предмет не найден"})
		return
	}

	// Проверяем, есть ли уроки с этим предметом
	lessonCollection := h.db.Collection("lessons")
	lessonCount, err := lessonCollection.CountDocuments(context.Background(), bson.M{"subject_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки уроков"})
		return
	}

	if lessonCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя удалить предмет, который используется в уроках"})
		return
	}

	// Удаляем предмет
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления предмета"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Предмет успешно удален"})
}

// ========== СТУДЕНТЫ ==========

// CreateStudent создает нового студента
func (h *Handlers) CreateStudent(c *gin.Context) {
	var req models.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование группы
	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
		return
	}

	groupCollection := h.db.Collection("groups")
	var group models.Group
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": groupID}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Группа не найдена"})
		return
	}

	student := models.Student{
		IIN:       req.IIN,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		GroupID:   groupID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := h.db.Collection("students")
	result, err := collection.InsertOne(context.Background(), student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания студента"})
		return
	}

	student.ID = result.InsertedID.(primitive.ObjectID)
	student.Group = &group
	c.JSON(http.StatusCreated, student)
}

// GetStudents получает всех студентов
func (h *Handlers) GetStudents(c *gin.Context) {
	collection := h.db.Collection("students")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения студентов"})
		return
	}
	defer cursor.Close(context.Background())

	var students []models.Student
	if err = cursor.All(context.Background(), &students); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки студентов"})
		return
	}

	// Загружаем информацию о группах
	for i := range students {
		groupCollection := h.db.Collection("groups")
		var group models.Group
		err = groupCollection.FindOne(context.Background(), bson.M{"_id": students[i].GroupID}).Decode(&group)
		if err == nil {
			students[i].Group = &group
		}
	}

	c.JSON(http.StatusOK, students)
}

// GetStudentSchedule получает расписание студента по ИИН
func (h *Handlers) GetStudentSchedule(c *gin.Context) {
	iin := c.Param("iin")

	// Находим студента по ИИН
	studentCollection := h.db.Collection("students")
	var student models.Student
	err := studentCollection.FindOne(context.Background(), bson.M{"iin": iin}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Студент не найден"})
		return
	}

	// Получаем расписание группы студента
	scheduleCollection := h.db.Collection("schedules")
	cursor, err := scheduleCollection.Find(context.Background(), bson.M{"group_id": student.GroupID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения расписания"})
		return
	}
	defer cursor.Close(context.Background())

	var schedules []models.Schedule
	if err = cursor.All(context.Background(), &schedules); err != nil {
		// Если есть ошибка декодирования, возвращаем пустой массив
		schedules = []models.Schedule{}
	}

	// Загружаем информацию о преподавателях
	for i := range schedules {
		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		err = teacherCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].TeacherID}).Decode(&teacher)
		if err == nil {
			schedules[i].Teacher = &teacher
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"student":   student,
		"schedules": schedules,
	})
}

// UpdateStudent обновляет студента
func (h *Handlers) UpdateStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID студента"})
		return
	}

	var req models.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование студента
	collection := h.db.Collection("students")
	var existingStudent models.Student
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingStudent)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Студент не найден"})
		return
	}

	// Если обновляется группа, проверяем её существование
	if req.GroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(req.GroupID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
			return
		}

		groupCollection := h.db.Collection("groups")
		var group models.Group
		err = groupCollection.FindOne(context.Background(), bson.M{"_id": groupID}).Decode(&group)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Группа не найдена"})
			return
		}
	}

	// Создаем объект для обновления
	update := bson.M{"updated_at": time.Now()}
	if req.IIN != "" {
		update["iin"] = req.IIN
	}
	if req.FirstName != "" {
		update["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		update["last_name"] = req.LastName
	}
	if req.GroupID != "" {
		groupID, _ := primitive.ObjectIDFromHex(req.GroupID)
		update["group_id"] = groupID
	}

	// Обновляем студента
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления студента"})
		return
	}

	// Получаем обновленного студента с информацией о группе
	var updatedStudent models.Student
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedStudent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного студента"})
		return
	}

	// Загружаем информацию о группе
	groupCollection := h.db.Collection("groups")
	var group models.Group
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": updatedStudent.GroupID}).Decode(&group)
	if err == nil {
		updatedStudent.Group = &group
	}

	c.JSON(http.StatusOK, updatedStudent)
}

// DeleteStudent удаляет студента
func (h *Handlers) DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID студента"})
		return
	}

	// Проверяем существование студента
	collection := h.db.Collection("students")
	var student models.Student
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Студент не найден"})
		return
	}

	// Удаляем студента
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления студента"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Студент успешно удален"})
}

// ========== ПРЕПОДАВАТЕЛИ ==========

// CreateTeacher создает нового преподавателя
func (h *Handlers) CreateTeacher(c *gin.Context) {
	var req models.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacher := models.Teacher{
		IIN:       req.IIN,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Subjects:  req.Subjects,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := h.db.Collection("teachers")
	result, err := collection.InsertOne(context.Background(), teacher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания преподавателя"})
		return
	}

	teacher.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, teacher)
}

// GetTeachers получает всех преподавателей
func (h *Handlers) GetTeachers(c *gin.Context) {
	collection := h.db.Collection("teachers")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения преподавателей"})
		return
	}
	defer cursor.Close(context.Background())

	var teachers []models.Teacher
	if err = cursor.All(context.Background(), &teachers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки преподавателей"})
		return
	}

	c.JSON(http.StatusOK, teachers)
}

// GetTeacherSchedule получает расписание преподавателя по ИИН
func (h *Handlers) GetTeacherSchedule(c *gin.Context) {
	iin := c.Param("iin")

	// Находим преподавателя по ИИН
	teacherCollection := h.db.Collection("teachers")
	var teacher models.Teacher
	err := teacherCollection.FindOne(context.Background(), bson.M{"iin": iin}).Decode(&teacher)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Преподаватель не найден"})
		return
	}

	// Получаем расписание преподавателя
	scheduleCollection := h.db.Collection("schedules")
	cursor, err := scheduleCollection.Find(context.Background(), bson.M{"teacher_id": teacher.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения расписания"})
		return
	}
	defer cursor.Close(context.Background())

	var schedules []models.Schedule
	if err = cursor.All(context.Background(), &schedules); err != nil {
		// Если есть ошибка декодирования, возвращаем пустой массив
		schedules = []models.Schedule{}
	}

	// Загружаем информацию о группах
	for i := range schedules {
		groupCollection := h.db.Collection("groups")
		var group models.Group
		err = groupCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].GroupID}).Decode(&group)
		if err == nil {
			schedules[i].Group = &group
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"teacher":   teacher,
		"schedules": schedules,
	})
}

// UpdateTeacher обновляет преподавателя
func (h *Handlers) UpdateTeacher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID преподавателя"})
		return
	}

	var req models.UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование преподавателя
	collection := h.db.Collection("teachers")
	var existingTeacher models.Teacher
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingTeacher)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Преподаватель не найден"})
		return
	}

	// Создаем объект для обновления
	update := bson.M{"updated_at": time.Now()}
	if req.IIN != "" {
		update["iin"] = req.IIN
	}
	if req.FirstName != "" {
		update["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		update["last_name"] = req.LastName
	}
	if req.Subjects != nil {
		update["subjects"] = req.Subjects
	}

	// Обновляем преподавателя
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления преподавателя"})
		return
	}

	// Получаем обновленного преподавателя
	var updatedTeacher models.Teacher
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedTeacher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного преподавателя"})
		return
	}

	c.JSON(http.StatusOK, updatedTeacher)
}

// DeleteTeacher удаляет преподавателя
func (h *Handlers) DeleteTeacher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID преподавателя"})
		return
	}

	// Проверяем существование преподавателя
	collection := h.db.Collection("teachers")
	var teacher models.Teacher
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&teacher)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Преподаватель не найден"})
		return
	}

	// Проверяем, есть ли расписания с этим преподавателем
	scheduleCollection := h.db.Collection("schedules")
	scheduleCount, err := scheduleCollection.CountDocuments(context.Background(), bson.M{"teacher_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки расписаний"})
		return
	}

	// Проверяем, есть ли уроки с этим преподавателем
	lessonCollection := h.db.Collection("lessons")
	lessonCount, err := lessonCollection.CountDocuments(context.Background(), bson.M{"teacher_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки уроков"})
		return
	}

	if scheduleCount > 0 || lessonCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя удалить преподавателя, у которого есть расписания или уроки"})
		return
	}

	// Удаляем преподавателя
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления преподавателя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Преподаватель успешно удален"})
}

// ========== РАСПИСАНИЕ ==========

// CreateSchedule создает новое расписание
func (h *Handlers) CreateSchedule(c *gin.Context) {
	var req models.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование группы
	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
		return
	}

	groupCollection := h.db.Collection("groups")
	var group models.Group
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": groupID}).Decode(&group)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Группа не найдена"})
		return
	}

	// Проверяем существование преподавателя
	teacherID, err := primitive.ObjectIDFromHex(req.TeacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID преподавателя"})
		return
	}

	teacherCollection := h.db.Collection("teachers")
	var teacher models.Teacher
	err = teacherCollection.FindOne(context.Background(), bson.M{"_id": teacherID}).Decode(&teacher)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Преподаватель не найден"})
		return
	}

	// Проверяем существование предмета
	subjectID, err := primitive.ObjectIDFromHex(req.SubjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID предмета"})
		return
	}

	subjectCollection := h.db.Collection("subjects")
	var subject models.Subject
	err = subjectCollection.FindOne(context.Background(), bson.M{"_id": subjectID}).Decode(&subject)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Предмет не найден"})
		return
	}

	schedule := models.Schedule{
		GroupID:     groupID,
		TeacherID:   teacherID,
		SubjectID:   subjectID,
		Room:        req.Room,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Shift:       req.Shift,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := h.db.Collection("schedules")
	result, err := collection.InsertOne(context.Background(), schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания расписания"})
		return
	}

	schedule.ID = result.InsertedID.(primitive.ObjectID)
	schedule.Group = &group
	schedule.Teacher = &teacher
	c.JSON(http.StatusCreated, schedule)
}

// GetSchedules получает все расписания
func (h *Handlers) GetSchedules(c *gin.Context) {
	collection := h.db.Collection("schedules")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения расписания"})
		return
	}
	defer cursor.Close(context.Background())

	var schedules []models.Schedule
	if err = cursor.All(context.Background(), &schedules); err != nil {
		// Если есть ошибка декодирования, возвращаем пустой массив
		schedules = []models.Schedule{}
	}

	// Загружаем информацию о группах и преподавателях
	for i := range schedules {
		// Группа
		groupCollection := h.db.Collection("groups")
		var group models.Group
		err = groupCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].GroupID}).Decode(&group)
		if err == nil {
			schedules[i].Group = &group
		}

		// Преподаватель
		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		err = teacherCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].TeacherID}).Decode(&teacher)
		if err == nil {
			schedules[i].Teacher = &teacher
		}

		// Предмет - проверяем, что SubjectID не пустой
		if !schedules[i].SubjectID.IsZero() {
			subjectCollection := h.db.Collection("subjects")
			var subject models.Subject
			err = subjectCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].SubjectID}).Decode(&subject)
			if err == nil {
				schedules[i].Subject = &subject
			} else {
				// Если предмет не найден, создаем пустой объект
				schedules[i].Subject = &models.Subject{
					ID:   schedules[i].SubjectID,
					Name: "Предмет не найден",
					Code: "Н/Д",
				}
			}
		} else {
			// Если SubjectID пустой, создаем пустой объект
			schedules[i].Subject = &models.Subject{
				Name: "Предмет не указан",
				Code: "Н/Д",
			}
		}
	}

	c.JSON(http.StatusOK, schedules)
}

// GetSchedulesByDay получает расписание по дню недели
func (h *Handlers) GetSchedulesByDay(c *gin.Context) {
	dayStr := c.Param("day")
	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный день недели (1-7)"})
		return
	}

	collection := h.db.Collection("schedules")
	cursor, err := collection.Find(context.Background(), bson.M{"day_of_week": day})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения расписания"})
		return
	}
	defer cursor.Close(context.Background())

	var schedules []models.Schedule
	if err = cursor.All(context.Background(), &schedules); err != nil {
		// Если есть ошибка декодирования, возвращаем пустой массив
		schedules = []models.Schedule{}
	}

	// Загружаем информацию о группах и преподавателях
	for i := range schedules {
		// Группа
		groupCollection := h.db.Collection("groups")
		var group models.Group
		err = groupCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].GroupID}).Decode(&group)
		if err == nil {
			schedules[i].Group = &group
		}

		// Преподаватель
		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		err = teacherCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].TeacherID}).Decode(&teacher)
		if err == nil {
			schedules[i].Teacher = &teacher
		}

		// Предмет - проверяем, что SubjectID не пустой
		if !schedules[i].SubjectID.IsZero() {
			subjectCollection := h.db.Collection("subjects")
			var subject models.Subject
			err = subjectCollection.FindOne(context.Background(), bson.M{"_id": schedules[i].SubjectID}).Decode(&subject)
			if err == nil {
				schedules[i].Subject = &subject
			} else {
				// Если предмет не найден, создаем пустой объект
				schedules[i].Subject = &models.Subject{
					ID:   schedules[i].SubjectID,
					Name: "Предмет не найден",
					Code: "Н/Д",
				}
			}
		} else {
			// Если SubjectID пустой, создаем пустой объект
			schedules[i].Subject = &models.Subject{
				Name: "Предмет не указан",
				Code: "Н/Д",
			}
		}
	}

	c.JSON(http.StatusOK, schedules)
}

// UpdateSchedule обновляет расписание
func (h *Handlers) UpdateSchedule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID расписания"})
		return
	}

	var req models.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование расписания
	collection := h.db.Collection("schedules")
	var existingSchedule models.Schedule
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingSchedule)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Расписание не найдено"})
		return
	}

	// Если обновляется группа, проверяем её существование
	if req.GroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(req.GroupID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
			return
		}

		groupCollection := h.db.Collection("groups")
		var group models.Group
		err = groupCollection.FindOne(context.Background(), bson.M{"_id": groupID}).Decode(&group)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Группа не найдена"})
			return
		}
	}

	// Если обновляется преподаватель, проверяем его существование
	if req.TeacherID != "" {
		teacherID, err := primitive.ObjectIDFromHex(req.TeacherID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID преподавателя"})
			return
		}

		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		err = teacherCollection.FindOne(context.Background(), bson.M{"_id": teacherID}).Decode(&teacher)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Преподаватель не найден"})
			return
		}
	}

	// Если обновляется предмет, проверяем его существование
	if req.SubjectID != "" {
		subjectID, err := primitive.ObjectIDFromHex(req.SubjectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID предмета"})
			return
		}

		subjectCollection := h.db.Collection("subjects")
		var subject models.Subject
		err = subjectCollection.FindOne(context.Background(), bson.M{"_id": subjectID}).Decode(&subject)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Предмет не найден"})
			return
		}
	}

	// Создаем объект для обновления
	update := bson.M{"updated_at": time.Now()}
	if req.GroupID != "" {
		groupID, _ := primitive.ObjectIDFromHex(req.GroupID)
		update["group_id"] = groupID
	}
	if req.TeacherID != "" {
		teacherID, _ := primitive.ObjectIDFromHex(req.TeacherID)
		update["teacher_id"] = teacherID
	}
	if req.SubjectID != "" {
		subjectID, _ := primitive.ObjectIDFromHex(req.SubjectID)
		update["subject_id"] = subjectID
	}
	if req.Room != "" {
		update["room"] = req.Room
	}
	if req.DayOfWeek != nil {
		if *req.DayOfWeek < 1 || *req.DayOfWeek > 7 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "День недели должен быть от 1 до 7"})
			return
		}
		update["day_of_week"] = *req.DayOfWeek
	}
	if req.StartTime != "" {
		update["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		update["end_time"] = req.EndTime
	}
	if req.Shift != nil {
		if *req.Shift < 1 || *req.Shift > 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Смена должна быть 1 или 2"})
			return
		}
		update["shift"] = *req.Shift
	}
	if req.Description != "" {
		update["description"] = req.Description
	}

	// Обновляем расписание
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления расписания"})
		return
	}

	// Получаем обновленное расписание с информацией о группе и преподавателе
	var updatedSchedule models.Schedule
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedSchedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного расписания"})
		return
	}

	// Загружаем информацию о группе
	groupCollection := h.db.Collection("groups")
	var group models.Group
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": updatedSchedule.GroupID}).Decode(&group)
	if err == nil {
		updatedSchedule.Group = &group
	}

	// Загружаем информацию о преподавателе
	teacherCollection := h.db.Collection("teachers")
	var teacher models.Teacher
	err = teacherCollection.FindOne(context.Background(), bson.M{"_id": updatedSchedule.TeacherID}).Decode(&teacher)
	if err == nil {
		updatedSchedule.Teacher = &teacher
	}

	c.JSON(http.StatusOK, updatedSchedule)
}

// DeleteSchedule удаляет расписание
func (h *Handlers) DeleteSchedule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID расписания"})
		return
	}

	// Проверяем существование расписания
	collection := h.db.Collection("schedules")
	var schedule models.Schedule
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&schedule)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Расписание не найдено"})
		return
	}

	// Удаляем расписание
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления расписания"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Расписание успешно удалено"})
}

// ========== КАЛЕНДАРЬ (УРОКИ) ==========

// CreateLesson создает новый урок
func (h *Handlers) CreateLesson(c *gin.Context) {
	var req models.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем ID
	groupID, err := primitive.ObjectIDFromHex(req.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
		return
	}

	teacherID, err := primitive.ObjectIDFromHex(req.TeacherID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID преподавателя"})
		return
	}

	subjectID, err := primitive.ObjectIDFromHex(req.SubjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID предмета"})
		return
	}

	// Создаем урок с базовыми данными
	lesson := models.Lesson{
		GroupID:     groupID,
		TeacherID:   teacherID,
		SubjectID:   subjectID,
		Room:        req.Room,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Если дата и время предоставлены, парсим их
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Используйте YYYY-MM-DD"})
			return
		}
		lesson.Date = &date
	}

	if req.StartTime != "" {
		lesson.StartTime = req.StartTime
		// Автоматически определяем смену по времени начала
		lesson.Shift = models.DetermineShift(req.StartTime)
	}

	if req.EndTime != "" {
		lesson.EndTime = req.EndTime
	}

	// Если смена указана вручную, используем её
	if req.Shift > 0 {
		lesson.Shift = req.Shift
	}

	collection := h.db.Collection("lessons")
	result, err := collection.InsertOne(context.Background(), lesson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания урока"})
		return
	}

	lesson.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, lesson)
}

// GetLessons получает все уроки
func (h *Handlers) GetLessons(c *gin.Context) {
	collection := h.db.Collection("lessons")

	// Получаем параметры запроса
	date := c.Query("date")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	groupID := c.Query("group_id")
	teacherID := c.Query("teacher_id")
	shift := c.Query("shift")

	filter := bson.M{}

	// Приоритет: start_date/end_date > date
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			end, err := time.Parse("2006-01-02", endDate)
			if err == nil {
				// Фильтр по диапазону дат
				startOfDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
				endOfDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
				filter["date"] = bson.M{
					"$gte": startOfDay,
					"$lte": endOfDay,
				}
			}
		}
	} else if date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err == nil {
			// Фильтр по дате (весь день)
			startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
			endOfDay := startOfDay.Add(24 * time.Hour)
			filter["date"] = bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			}
		}
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

	if shift != "" {
		if shiftNum, err := strconv.Atoi(shift); err == nil && (shiftNum == 1 || shiftNum == 2) {
			filter["shift"] = shiftNum
		}
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения уроков"})
		return
	}
	defer cursor.Close(context.Background())

	var lessons []models.Lesson
	if err = cursor.All(context.Background(), &lessons); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки уроков"})
		return
	}

	// Заполняем связанные данные
	for i := range lessons {
		// Получаем группу
		groupCollection := h.db.Collection("groups")
		var group models.Group
		if err := groupCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].GroupID}).Decode(&group); err == nil {
			lessons[i].Group = &group
		}

		// Получаем преподавателя
		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		if err := teacherCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].TeacherID}).Decode(&teacher); err == nil {
			lessons[i].Teacher = &teacher
		}

		// Получаем предмет
		subjectCollection := h.db.Collection("subjects")
		var subject models.Subject
		if err := subjectCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].SubjectID}).Decode(&subject); err == nil {
			lessons[i].Subject = &subject
		}
	}

	c.JSON(http.StatusOK, lessons)
}

// GetLessonsByDate получает уроки по конкретной дате
func (h *Handlers) GetLessonsByDate(c *gin.Context) {
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Используйте YYYY-MM-DD"})
		return
	}

	// Фильтр по дате (весь день)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	collection := h.db.Collection("lessons")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения уроков"})
		return
	}
	defer cursor.Close(context.Background())

	var lessons []models.Lesson
	if err = cursor.All(context.Background(), &lessons); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки уроков"})
		return
	}

	// Заполняем связанные данные
	for i := range lessons {
		// Получаем группу
		groupCollection := h.db.Collection("groups")
		var group models.Group
		if err := groupCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].GroupID}).Decode(&group); err == nil {
			lessons[i].Group = &group
		}

		// Получаем преподавателя
		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		if err := teacherCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].TeacherID}).Decode(&teacher); err == nil {
			lessons[i].Teacher = &teacher
		}

		// Получаем предмет
		subjectCollection := h.db.Collection("subjects")
		var subject models.Subject
		if err := subjectCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].SubjectID}).Decode(&subject); err == nil {
			lessons[i].Subject = &subject
		}
	}

	c.JSON(http.StatusOK, lessons)
}

// UpdateLesson обновляет урок
func (h *Handlers) UpdateLesson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID урока"})
		return
	}

	var req models.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := h.db.Collection("lessons")

	// Проверяем существование урока
	var existingLesson models.Lesson
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingLesson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Урок не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска урока"})
		}
		return
	}

	// Подготавливаем обновления
	update := bson.M{
		"updated_at": time.Now(),
	}

	if req.GroupID != "" {
		groupID, err := primitive.ObjectIDFromHex(req.GroupID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID группы"})
			return
		}
		update["group_id"] = groupID
	}

	if req.TeacherID != "" {
		teacherID, err := primitive.ObjectIDFromHex(req.TeacherID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID преподавателя"})
			return
		}
		update["teacher_id"] = teacherID
	}

	if req.SubjectID != "" {
		subjectID, err := primitive.ObjectIDFromHex(req.SubjectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID предмета"})
			return
		}
		update["subject_id"] = subjectID
	}

	if req.Room != "" {
		update["room"] = req.Room
	}

	if req.Date != "" {
		// Логируем входящую дату для отладки
		log.Printf("Parsing date: '%s'", req.Date)
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			log.Printf("Date parsing error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Используйте YYYY-MM-DD"})
			return
		}
		update["date"] = &date
	}

	if req.StartTime != "" {
		update["start_time"] = req.StartTime
		// Автоматически определяем смену по времени начала
		update["shift"] = models.DetermineShift(req.StartTime)
	}

	if req.EndTime != "" {
		update["end_time"] = req.EndTime
	}

	// Если смена указана вручную, используем её
	if req.Shift != nil && *req.Shift > 0 {
		update["shift"] = *req.Shift
	}

	if req.Description != "" {
		update["description"] = req.Description
	}

	// Обновляем урок
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления урока"})
		return
	}

	// Получаем обновленный урок
	var updatedLesson models.Lesson
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedLesson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного урока"})
		return
	}

	c.JSON(http.StatusOK, updatedLesson)
}

// DeleteLesson удаляет урок
func (h *Handlers) DeleteLesson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID урока"})
		return
	}

	collection := h.db.Collection("lessons")

	// Проверяем существование урока
	var existingLesson models.Lesson
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingLesson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Урок не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка поиска урока"})
		}
		return
	}

	// Удаляем урок
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления урока"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Урок успешно удален"})
}

// GetAvailableLessons получает уроки без даты и времени (доступные для назначения)
func (h *Handlers) GetAvailableLessons(c *gin.Context) {
	collection := h.db.Collection("lessons")

	// Фильтр для уроков без даты (доступные уроки)
	filter := bson.M{
		"date": bson.M{"$exists": false},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доступных уроков"})
		return
	}
	defer cursor.Close(context.Background())

	var lessons []models.Lesson
	if err = cursor.All(context.Background(), &lessons); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки доступных уроков"})
		return
	}

	// Заполняем связанные данные
	for i := range lessons {
		// Получаем группу
		groupCollection := h.db.Collection("groups")
		var group models.Group
		if err := groupCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].GroupID}).Decode(&group); err == nil {
			lessons[i].Group = &group
		}

		// Получаем преподавателя
		teacherCollection := h.db.Collection("teachers")
		var teacher models.Teacher
		if err := teacherCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].TeacherID}).Decode(&teacher); err == nil {
			lessons[i].Teacher = &teacher
		}

		// Получаем предмет
		subjectCollection := h.db.Collection("subjects")
		var subject models.Subject
		if err := subjectCollection.FindOne(context.Background(), bson.M{"_id": lessons[i].SubjectID}).Decode(&subject); err == nil {
			lessons[i].Subject = &subject
		}
	}

	c.JSON(http.StatusOK, lessons)
}
