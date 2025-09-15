package routes

import (
	"innovativecollege/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handlers) {
	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Группы
		api.POST("/groups", h.CreateGroup)
		api.GET("/groups", h.GetGroups)
		api.PUT("/groups/:id", h.UpdateGroup)
		api.DELETE("/groups/:id", h.DeleteGroup)

		// Предметы
		api.POST("/subjects", h.CreateSubject)
		api.GET("/subjects", h.GetSubjects)
		api.PUT("/subjects/:id", h.UpdateSubject)
		api.DELETE("/subjects/:id", h.DeleteSubject)

		// Студенты
		api.POST("/students", h.CreateStudent)
		api.GET("/students", h.GetStudents)
		api.GET("/students/:iin/schedule", h.GetStudentSchedule)
		api.PUT("/students/:id", h.UpdateStudent)
		api.DELETE("/students/:id", h.DeleteStudent)

		// Преподаватели
		api.POST("/teachers", h.CreateTeacher)
		api.GET("/teachers", h.GetTeachers)
		api.GET("/teachers/:iin/schedule", h.GetTeacherSchedule)
		api.PUT("/teachers/:id", h.UpdateTeacher)
		api.DELETE("/teachers/:id", h.DeleteTeacher)

		// Расписание
		api.POST("/schedules", h.CreateSchedule)
		api.GET("/schedules", h.GetSchedules)
		api.GET("/schedules/day/:day", h.GetSchedulesByDay)
		api.PUT("/schedules/:id", h.UpdateSchedule)
		api.DELETE("/schedules/:id", h.DeleteSchedule)

		// Календарь (Уроки)
		api.POST("/lessons", h.CreateLesson)
		api.GET("/lessons", h.GetLessons)
		api.GET("/lessons/available", h.GetAvailableLessons)
		api.GET("/lessons/date/:date", h.GetLessonsByDate)
		api.PUT("/lessons/:id", h.UpdateLesson)
		api.DELETE("/lessons/:id", h.DeleteLesson)

		// Временные слоты
		api.POST("/time-slots", h.CreateTimeSlot)
		api.GET("/time-slots", h.GetTimeSlots)
		api.GET("/time-slots/:id", h.GetTimeSlot)
		api.PUT("/time-slots/:id", h.UpdateTimeSlot)
		api.DELETE("/time-slots/:id", h.DeleteTimeSlot)

		// Статистика
		api.GET("/statistics/lessons", h.GetLessonStatistics)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
