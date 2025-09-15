package models

import (
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group представляет группу студентов
type Group struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Subject представляет предмет
type Subject struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Code        string             `bson:"code" json:"code"` // Например: ОН 3.1, ОН 2.2
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Student представляет студента
type Student struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IIN       string             `bson:"iin" json:"iin"` // ИИН для входа
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	GroupID   primitive.ObjectID `bson:"group_id" json:"group_id"`
	Group     *Group             `bson:"group,omitempty" json:"group,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Teacher представляет преподавателя
type Teacher struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IIN       string             `bson:"iin" json:"iin"` // ИИН для входа
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Subjects  []string           `bson:"subjects" json:"subjects"` // Предметы которые ведет
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Schedule представляет расписание
type Schedule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID     primitive.ObjectID `bson:"group_id" json:"group_id"`
	Group       *Group             `bson:"group,omitempty" json:"group,omitempty"`
	TeacherID   primitive.ObjectID `bson:"teacher_id" json:"teacher_id"`
	Teacher     *Teacher           `bson:"teacher,omitempty" json:"teacher,omitempty"`
	SubjectID   primitive.ObjectID `bson:"subject_id" json:"subject_id"`
	Subject     *Subject           `bson:"subject,omitempty" json:"subject,omitempty"`
	Room        string             `bson:"room" json:"room"`
	DayOfWeek   int                `bson:"day_of_week" json:"day_of_week"` // 1-7 (понедельник-воскресенье)
	StartTime   string             `bson:"start_time" json:"start_time"`   // "12:40"
	EndTime     string             `bson:"end_time" json:"end_time"`       // "14:00"
	Shift       int                `bson:"shift" json:"shift"`             // 1 или 2 смена
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Lesson представляет урок в календаре
type Lesson struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID     primitive.ObjectID `bson:"group_id" json:"group_id"`
	Group       *Group             `bson:"group,omitempty" json:"group,omitempty"`
	TeacherID   primitive.ObjectID `bson:"teacher_id" json:"teacher_id"`
	Teacher     *Teacher           `bson:"teacher,omitempty" json:"teacher,omitempty"`
	SubjectID   primitive.ObjectID `bson:"subject_id" json:"subject_id"`
	Subject     *Subject           `bson:"subject,omitempty" json:"subject,omitempty"`
	Room        string             `bson:"room" json:"room"`
	Date        *time.Time         `bson:"date,omitempty" json:"date,omitempty"`             // Конкретная дата урока
	StartTime   string             `bson:"start_time,omitempty" json:"start_time,omitempty"` // "12:40"
	EndTime     string             `bson:"end_time,omitempty" json:"end_time,omitempty"`     // "14:00"
	Shift       int                `bson:"shift,omitempty" json:"shift,omitempty"`           // 1 или 2 смена (автоматически определяется)
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateGroupRequest запрос на создание группы
type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateSubjectRequest запрос на создание предмета
type CreateSubjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreateStudentRequest запрос на создание студента
type CreateStudentRequest struct {
	IIN       string `json:"iin" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	GroupID   string `json:"group_id" binding:"required"`
}

// CreateTeacherRequest запрос на создание преподавателя
type CreateTeacherRequest struct {
	IIN       string   `json:"iin" binding:"required"`
	FirstName string   `json:"first_name" binding:"required"`
	LastName  string   `json:"last_name" binding:"required"`
	Subjects  []string `json:"subjects"`
}

// CreateScheduleRequest запрос на создание расписания
type CreateScheduleRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	TeacherID   string `json:"teacher_id" binding:"required"`
	SubjectID   string `json:"subject_id" binding:"required"`
	Room        string `json:"room" binding:"required"`
	DayOfWeek   int    `json:"day_of_week" binding:"required,min=1,max=7"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
	Shift       int    `json:"shift" binding:"required,min=1,max=2"`
	Description string `json:"description,omitempty"`
}

// UpdateGroupRequest запрос на обновление группы
type UpdateGroupRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateSubjectRequest запрос на обновление предмета
type UpdateSubjectRequest struct {
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateStudentRequest запрос на обновление студента
type UpdateStudentRequest struct {
	IIN       string `json:"iin,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
}

// UpdateTeacherRequest запрос на обновление преподавателя
type UpdateTeacherRequest struct {
	IIN       string   `json:"iin,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Subjects  []string `json:"subjects,omitempty"`
}

// UpdateScheduleRequest запрос на обновление расписания
type UpdateScheduleRequest struct {
	GroupID     string `json:"group_id,omitempty"`
	TeacherID   string `json:"teacher_id,omitempty"`
	SubjectID   string `json:"subject_id,omitempty"`
	Room        string `json:"room,omitempty"`
	DayOfWeek   *int   `json:"day_of_week,omitempty"` // Указатель для проверки на 0
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
	Shift       *int   `json:"shift,omitempty"` // Указатель для проверки на 0
	Description string `json:"description,omitempty"`
}

// CreateLessonRequest запрос на создание урока
type CreateLessonRequest struct {
	GroupID     string `json:"group_id" binding:"required"`
	TeacherID   string `json:"teacher_id" binding:"required"`
	SubjectID   string `json:"subject_id" binding:"required"`
	Room        string `json:"room" binding:"required"`
	Date        string `json:"date,omitempty"`       // "2024-01-15" - необязательное поле
	StartTime   string `json:"start_time,omitempty"` // "12:40" - необязательное поле
	EndTime     string `json:"end_time,omitempty"`   // "14:00" - необязательное поле
	Shift       int    `json:"shift,omitempty"`      // 1 или 2 смена - необязательное поле (автоматически определяется)
	Description string `json:"description,omitempty"`
}

// UpdateLessonRequest запрос на обновление урока
type UpdateLessonRequest struct {
	GroupID     string `json:"group_id,omitempty"`
	TeacherID   string `json:"teacher_id,omitempty"`
	SubjectID   string `json:"subject_id,omitempty"`
	Room        string `json:"room,omitempty"`
	Date        string `json:"date,omitempty"`       // "2024-01-15"
	StartTime   string `json:"start_time,omitempty"` // "12:40"
	EndTime     string `json:"end_time,omitempty"`   // "14:00"
	Shift       *int   `json:"shift,omitempty"`      // 1 или 2 смена (указатель для проверки на 0)
	Description string `json:"description,omitempty"`
}

// DetermineShift определяет смену по времени начала урока
// Первая смена: 8:00-12:30, Вторая смена: 12:40-17:00
func DetermineShift(startTime string) int {
	if startTime == "" {
		return 0
	}

	// Парсим время в формате "HH:MM"
	parts := strings.Split(startTime, ":")
	if len(parts) != 2 {
		return 0
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}

	// Конвертируем в минуты для удобства сравнения
	totalMinutes := hour*60 + minute

	// Первая смена: 8:00 (480 мин) - 12:30 (750 мин)
	// Вторая смена: 12:40 (760 мин) - 17:00 (1020 мин)
	if totalMinutes >= 480 && totalMinutes <= 750 {
		return 1 // Первая смена
	} else if totalMinutes >= 760 && totalMinutes <= 1020 {
		return 2 // Вторая смена
	}

	return 0 // Неопределенная смена
}

// TimeSlot представляет временной слот в расписании
type TimeSlot struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StartTime string             `bson:"start_time" json:"start_time"` // "12:40"
	EndTime   string             `bson:"end_time" json:"end_time"`     // "14:00"
	Shift     int                `bson:"shift" json:"shift"`           // 1 или 2
	Label     string             `bson:"label" json:"label"`           // "12:40-14:00"
	IsActive  bool               `bson:"is_active" json:"is_active"`   // Активен ли слот
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreateTimeSlotRequest запрос на создание временного слота
type CreateTimeSlotRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Shift     int    `json:"shift" binding:"required,min=1,max=2"`
	Label     string `json:"label,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
}

// UpdateTimeSlotRequest запрос на обновление временного слота
type UpdateTimeSlotRequest struct {
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
	Shift     *int    `json:"shift,omitempty"`
	Label     *string `json:"label,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
