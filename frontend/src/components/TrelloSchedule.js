  import React, { useState, useEffect } from 'react';
import { groupsApi, teachersApi, subjectsApi, lessonsApi, timeSlotsApi } from '../services/api';
import LessonForm from './LessonForm';

const TrelloSchedule = () => {
  const [groups, setGroups] = useState([]);
  const [teachers, setTeachers] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [allLessons, setAllLessons] = useState([]);
  const [scheduledLessons, setScheduledLessons] = useState([]);
  const [timeSlots, setTimeSlots] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showLessonForm, setShowLessonForm] = useState(false);
  const [editingLesson, setEditingLesson] = useState(null);
  const [selectedDay, setSelectedDay] = useState(1); // Понедельник по умолчанию
  const [selectedShift, setSelectedShift] = useState(1);
  const [selectedLesson, setSelectedLesson] = useState(null);
  const [removingLessons, setRemovingLessons] = useState(new Set());
  const [isCreatingLesson, setIsCreatingLesson] = useState(false);
  const [notification, setNotification] = useState(null);

  const daysOfWeek = [
    { value: 1, label: 'Понедельник' },
    { value: 2, label: 'Вторник' },
    { value: 3, label: 'Среда' },
    { value: 4, label: 'Четверг' },
    { value: 5, label: 'Пятница' }
  ];

  const shifts = [
    { value: 1, label: 'Первая смена' },
    { value: 2, label: 'Вторая смена' }
  ];

  // Функция для показа уведомлений
  const showNotification = (message, type = 'error') => {
    setNotification({ message, type });
    setTimeout(() => {
      setNotification(null);
    }, 4000); // Автоматически скрываем через 4 секунды
  };

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    fetchLessons();
  }, [selectedDay, selectedShift]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [groupsResponse, teachersResponse, subjectsResponse, availableLessonsResponse, timeSlotsResponse] = await Promise.all([
        groupsApi.getAll(),
        teachersApi.getAll(),
        subjectsApi.getAll(),
        lessonsApi.getAvailable(),
        timeSlotsApi.getActive()
      ]);
      setGroups(groupsResponse.data || []);
      setTeachers(teachersResponse.data || []);
      setSubjects(subjectsResponse.data || []);
      setAllLessons(availableLessonsResponse.data || []);
      setTimeSlots(timeSlotsResponse.data || []);
      setError(null);
      
      // Загружаем уроки для текущей смены
      await fetchLessons();
    } catch (err) {
      setError('Ошибка загрузки данных: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  const fetchLessons = async () => {
    try {
      // Загружаем все уроки
      const lessonsResponse = await lessonsApi.getAll();
      const allLessonsData = lessonsResponse.data || [];
      
      // Обновляем общий список уроков
      setAllLessons(allLessonsData);
      
      // Фильтруем уроки по дню недели и смене
      const filteredLessons = allLessonsData.filter(lesson => {
        // Показываем уроки без даты в доступных уроках
        if (!lesson.date) return false;
        
        // Проверяем день недели
        const lessonDate = new Date(lesson.date);
        const dayOfWeek = lessonDate.getDay() === 0 ? 7 : lessonDate.getDay(); // Воскресенье = 7
        const isCorrectDay = dayOfWeek === selectedDay;
        
        // Проверяем смену (если указана)
        const isCorrectShift = !lesson.shift || lesson.shift === selectedShift;
        
        
        return isCorrectDay && isCorrectShift;
      });
      setScheduledLessons(filteredLessons);
    } catch (err) {
      setError('Ошибка загрузки уроков: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleLessonClick = (lesson) => {
    setSelectedLesson(lesson);
  };

  const handleCancelSelection = () => {
    setSelectedLesson(null);
  };

  const handleSlotSelect = async (timeSlotId) => {
    if (!selectedLesson) return;

    const currentSlots = getTimeSlotsByShift(selectedShift);
    const slot = currentSlots.find(s => s.id === timeSlotId);

    if (!slot) return;

    // Проверяем, не добавлен ли уже такой урок в этот слот
    if (isLessonAlreadyInSlot(timeSlotId, selectedLesson)) {
      showNotification('Этот урок уже добавлен в данный слот! Один урок не может быть назначен дважды в одно время.', 'warning');
      return;
    }

    try {
      const startTime = slot.start_time;
      const endTime = slot.end_time;
      
      // Вычисляем дату для выбранного дня недели относительно текущей недели
      const today = new Date();
      const currentDayOfWeek = today.getDay() === 0 ? 7 : today.getDay(); // Воскресенье = 7
      
      // Вычисляем понедельник текущей недели
      const mondayOfCurrentWeek = new Date(today);
      mondayOfCurrentWeek.setDate(today.getDate() - currentDayOfWeek + 1);
      
      // Вычисляем дату для выбранного дня недели
      const targetDate = new Date(mondayOfCurrentWeek);
      targetDate.setDate(mondayOfCurrentWeek.getDate() + selectedDay - 1);
      
      // Создаем НОВЫЙ урок на основе выбранного (копию)
      const newScheduledLesson = {
        group_id: selectedLesson.group_id,
        teacher_id: selectedLesson.teacher_id,
        subject_id: selectedLesson.subject_id,
        room: selectedLesson.room,
        description: selectedLesson.description,
        date: targetDate.toISOString().split('T')[0],
        start_time: startTime,
        end_time: endTime,
        shift: selectedShift
      };
      
      
      // Создаем новый урок в расписании
      await lessonsApi.create(newScheduledLesson);
      
      // Перезагружаем доступные уроки (оригинал должен остаться)
      const availableLessonsResponse = await lessonsApi.getAvailable();
      setAllLessons(availableLessonsResponse.data || []);
      
      // Перезагружаем уроки для текущей смены
      await fetchLessons();

      setSelectedLesson(null);
      setError(''); // Очищаем ошибки при успешном добавлении
      showNotification('Урок успешно добавлен в расписание!', 'success');
    } catch (err) {
      setError('Ошибка обновления урока: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleRemoveFromSchedule = async (lessonId) => {
    if (removingLessons.has(lessonId)) return; // Предотвращаем множественные клики для конкретного урока
    
    setRemovingLessons(prev => new Set(prev).add(lessonId));
    try {
      // Ищем урок в scheduledLessons (уроки в расписании)
      const lesson = scheduledLessons.find(l => l.id === lessonId);
      if (lesson) {
        // Используем DELETE запрос для удаления урока из расписания
        await lessonsApi.delete(lessonId);
        
        // Перезагружаем доступные уроки с сервера
        const availableLessonsResponse = await lessonsApi.getAvailable();
        setAllLessons(availableLessonsResponse.data || []);
        
        // Перезагружаем данные с сервера для подтверждения
        await fetchLessons();
        showNotification('Урок удален из расписания!', 'success');
      }
    } catch (err) {
      setError('Ошибка удаления урока: ' + (err.response?.data?.error || err.message));
    } finally {
      setRemovingLessons(prev => {
        const newSet = new Set(prev);
        newSet.delete(lessonId);
        return newSet;
      });
    }
  };

  const handleCreateLesson = async (lessonData) => {
    try {
      setIsCreatingLesson(true);
      setError(null);
      
      // Создаем урок только с основными данными, без даты и времени
      const newLesson = {
        group_id: lessonData.group_id,
        teacher_id: lessonData.teacher_id,
        subject_id: lessonData.subject_id,
        room: lessonData.room,
        description: lessonData.description,
        // Дата и время будут установлены при перетаскивании
        date: null,
        start_time: null,
        end_time: null
      };
      
      const response = await lessonsApi.create(newLesson);
      
      // Перезагружаем доступные уроки
      const availableLessonsResponse = await lessonsApi.getAvailable();
      setAllLessons(availableLessonsResponse.data || []);
      
      // Перезагружаем уроки для текущей смены
      await fetchLessons();
      
      setShowLessonForm(false);
      showNotification('Урок успешно создан!', 'success');
    } catch (err) {
      setError('Ошибка создания урока: ' + (err.response?.data?.error || err.message));
    } finally {
      setIsCreatingLesson(false);
    }
  };

  const handleEditLesson = async (lessonData) => {
    try {
      // Убеждаемся, что дата в правильном формате YYYY-MM-DD
      let formattedDate = editingLesson.date;
      if (formattedDate) {
        // Если дата в формате ISO string, извлекаем только дату
        if (formattedDate.includes('T')) {
          formattedDate = formattedDate.split('T')[0];
        }
        // Если дата в формате Date object, конвертируем в YYYY-MM-DD
        else if (formattedDate instanceof Date) {
          formattedDate = formattedDate.toISOString().split('T')[0];
        }
        // Если дата уже в правильном формате, оставляем как есть
        // Дополнительная проверка формата YYYY-MM-DD
        if (formattedDate && !/^\d{4}-\d{2}-\d{2}$/.test(formattedDate)) {
          console.error('Invalid date format:', formattedDate);
          setError('Неверный формат даты. Ожидается YYYY-MM-DD');
          return;
        }
      }
      
      const updatedLesson = {
        ...editingLesson,
        ...lessonData,
        date: formattedDate
      };
      
      // Отладочная информация
      console.log('Original lesson date:', editingLesson.date);
      console.log('Formatted date:', formattedDate);
      console.log('Updating lesson with data:', updatedLesson);
      
      await lessonsApi.update(editingLesson.id, updatedLesson);
      
      // Обновляем урок в локальном состоянии allLessons
      setAllLessons(prevLessons => 
        prevLessons.map(l => 
          l.id === editingLesson.id 
            ? updatedLesson
            : l
        )
      );
      
      // Перезагружаем уроки для текущей смены
      await fetchLessons();
      
      setShowLessonForm(false);
      setEditingLesson(null);
    } catch (err) {
      setError('Ошибка обновления урока: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleDeleteLesson = async (lessonId) => {
    if (window.confirm('Вы уверены, что хотите удалить этот урок?')) {
      try {
        await lessonsApi.delete(lessonId);
        
        // Перезагружаем доступные уроки с сервера
        const availableLessonsResponse = await lessonsApi.getAvailable();
        setAllLessons(availableLessonsResponse.data || []);
        
        // Перезагружаем уроки для текущей смены
        await fetchLessons();
      } catch (err) {
        setError('Ошибка удаления урока: ' + (err.response?.data?.error || err.message));
      }
    }
  };

  const getGroupName = (groupId) => {
    const group = groups.find(g => g.id === groupId);
    return group ? group.name : 'Неизвестная группа';
  };

  const getTeacherName = (teacherId) => {
    const teacher = teachers.find(t => t.id === teacherId);
    return teacher ? `${teacher.first_name} ${teacher.last_name}` : 'Неизвестный преподаватель';
  };

  const getSubjectName = (subjectId) => {
    const subject = subjects.find(s => s.id === subjectId);
    return subject ? subject.name : 'Неизвестный предмет';
  };

  const getSubjectCode = (subjectId) => {
    const subject = subjects.find(s => s.id === subjectId);
    return subject ? subject.code : 'ОН';
  };

  // Получаем слоты по смене
  const getTimeSlotsByShift = (shift) => {
    return timeSlots.filter(slot => slot.shift === shift);
  };

  const getLessonsForTimeSlot = (timeSlotId) => {
    const currentSlots = getTimeSlotsByShift(selectedShift);
    const slot = currentSlots.find(s => s.id === timeSlotId);
    if (!slot) return [];

    const slotStartTime = slot.start_time;
    const slotEndTime = slot.end_time;
    
    const filteredLessons = scheduledLessons.filter(lesson => {
      // Нормализуем время (убираем лишние пробелы и приводим к единому формату)
      const lessonStartTime = lesson.start_time ? lesson.start_time.trim() : '';
      const lessonEndTime = lesson.end_time ? lesson.end_time.trim() : '';
      
      return lessonStartTime === slotStartTime && lessonEndTime === slotEndTime;
    });
    return filteredLessons;
  };

  // Проверяем, есть ли уже такой урок в слоте
  const isLessonAlreadyInSlot = (timeSlotId, lessonToCheck) => {
    const existingLessons = getLessonsForTimeSlot(timeSlotId);
    
    return existingLessons.some(lesson => 
      lesson.group_id === lessonToCheck.group_id &&
      lesson.teacher_id === lessonToCheck.teacher_id &&
      lesson.subject_id === lessonToCheck.subject_id
    );
  };

  const currentTimeSlots = getTimeSlotsByShift(selectedShift);

  if (loading) {
    return <div className="loading">Загрузка расписания...</div>;
  }

  return (
    <div 
      className={`trello-schedule ${selectedLesson ? 'slot-selection-mode' : ''}`}
      onClick={(e) => {
        // Если клик по затемненной области (не по дочерним элементам)
        if (e.target === e.currentTarget || e.target.classList.contains('trello-schedule')) {
          handleCancelSelection();
        }
      }}
    >
      {/* Уведомления */}
      {notification && (
        <div className={`notification notification-${notification.type}`}>
          <div className="notification-content">
            <span className="notification-icon">
              {notification.type === 'warning' ? '⚠️' : notification.type === 'success' ? '✅' : '❌'}
            </span>
            <span className="notification-message">{notification.message}</span>
            <button 
              className="notification-close"
              onClick={() => setNotification(null)}
            >
              ×
            </button>
          </div>
        </div>
      )}

      {/* Затемненный фон при выборе слота */}
      {selectedLesson && (
        <div className="background-overlay" onClick={handleCancelSelection}></div>
      )}
      
      <div className="schedule-header">
        {error && <div className="error">{error}</div>}
        
        
        <div className="controls">
          <div className="control-group">
            <label>День недели:</label>
            <div className="day-buttons">
              {daysOfWeek.map((day) => (
                <button
                  key={day.value}
                  className={`day-button ${selectedDay === day.value ? 'active' : ''}`}
                  onClick={() => setSelectedDay(day.value)}
                >
                  {day.label}
                </button>
              ))}
            </div>
          </div>

          <button 
            className="btn btn-primary"
            onClick={() => setShowLessonForm(true)}
            disabled={isCreatingLesson}
          >
            {isCreatingLesson ? '⏳ Создание...' : '+ Добавить урок'}
          </button>

          <div className="control-group">
            <label>Смена:</label>
            <div className="shift-buttons">
              {shifts.map((shift) => (
                <button
                  key={shift.value}
                  className={`shift-button ${selectedShift === shift.value ? 'active' : ''}`}
                  onClick={() => setSelectedShift(shift.value)}
                >
                  {shift.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Секция с доступными уроками */}
      <div className="available-lessons">
        <h3>Доступные уроки</h3>
        {selectedLesson && (
          <div className="selected-lesson-info">
            <h4>Выбранный урок: {getSubjectName(selectedLesson.subject_id)}</h4>
            <p>Группа: {getGroupName(selectedLesson.group_id)} | Преподаватель: {getTeacherName(selectedLesson.teacher_id)}</p>
            <button className="btn btn-sm btn-secondary" onClick={handleCancelSelection}>
              Отмена
            </button>
          </div>
        )}
        <div className="lessons-pool">
          {allLessons
            .filter(lesson => !lesson.date) // Показываем только уроки без даты (незапланированные)
            .map(lesson => (
              <div 
                key={lesson.id} 
                className={`available-lesson-card clickable ${selectedLesson && selectedLesson.id === lesson.id ? 'selected' : ''}`}
                onClick={() => handleLessonClick(lesson)}
              >
                <div className="lesson-info">
                  <div className="lesson-subject">{getSubjectName(lesson.subject_id)}</div>
                  <div className="lesson-group">{getGroupName(lesson.group_id)}</div>
                  <div className="lesson-teacher">{getTeacherName(lesson.teacher_id)}</div>
                  <div className="lesson-room">Аудитория: {lesson.room}</div>
                </div>
                <div className="lesson-actions">
                  <button 
                    className="btn btn-sm btn-secondary"
                    onClick={(e) => {
                      e.stopPropagation();
                      setEditingLesson(lesson);
                      setShowLessonForm(true);
                    }}
                  >
                    Редактировать
                  </button>
                  <button 
                    className="btn btn-sm btn-danger"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeleteLesson(lesson.id);
                    }}
                  >
                    Удалить
                  </button>
                </div>
              </div>
            ))}
        </div>
      </div>

      <div className="schedule-board">
        {currentTimeSlots.map((timeSlot) => {
          const isDuplicate = selectedLesson && isLessonAlreadyInSlot(timeSlot.id, selectedLesson);
          return (
            <div key={timeSlot.id} className={`time-slot-column ${isDuplicate ? 'duplicate-slot' : ''}`}>
              <div className="time-slot-header">
                <h3>{timeSlot.label}</h3>
                <small style={{color: '#666'}}>
                  Уроков: {getLessonsForTimeSlot(timeSlot.id).length}
                </small>
                {isDuplicate && (
                  <div style={{color: '#ff6b6b', fontSize: '12px', fontWeight: 'bold'}}>
                    Уже добавлен!
                  </div>
                )}
              </div>
            
            <div className="lesson-list">
              {getLessonsForTimeSlot(timeSlot.id).map((lesson) => (
                <div key={lesson.id} className="lesson-card">
                  <div className="lesson-header">
                    <span className="course-code">{getSubjectCode(lesson.subject_id)}</span>
                    <div className="lesson-actions">
                      <button
                        className="btn-edit"
                        onClick={(e) => {
                          e.stopPropagation();
                          e.preventDefault();
                          setEditingLesson(lesson);
                          setShowLessonForm(true);
                        }}
                      >
                        Редактировать
                      </button>
                      <button
                        className="btn-delete"
                        onClick={(e) => {
                          e.stopPropagation();
                          e.preventDefault();
                          handleRemoveFromSchedule(lesson.id);
                        }}
                        title="Убрать из расписания"
                        disabled={removingLessons.has(lesson.id)}
                      >
                        {removingLessons.has(lesson.id) ? '⏳' : '📤'}
                      </button>
                    </div>
                  </div>
                  
                  <div className="lesson-content">
                    <div className="lesson-subject">{getSubjectName(lesson.subject_id)}</div>
                    <div className="lesson-group">{getGroupName(lesson.group_id)}</div>
                    <div className="lesson-teacher">{getTeacherName(lesson.teacher_id)}</div>
                    <div className="lesson-room">Аудитория: {lesson.room}</div>
                    {lesson.description && (
                      <div className="lesson-description">{lesson.description}</div>
                    )}
                  </div>
                </div>
              ))}
              
              {/* Кнопка + для назначения выбранного урока в слот */}
              {selectedLesson && (
                <button 
                  className="add-lesson-btn slot-selection"
                  onClick={() => handleSlotSelect(timeSlot.id)}
                >
                  +
                </button>
              )}
            </div>
          </div>
          );
        })}
      </div>


      {showLessonForm && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h3>{editingLesson ? 'Редактировать урок' : 'Создать урок'}</h3>
              <button 
                className="close-btn"
                onClick={() => {
                  setShowLessonForm(false);
                  setEditingLesson(null);
                }}
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <LessonForm
                onSubmit={editingLesson ? handleEditLesson : handleCreateLesson}
                onCancel={() => {
                  setShowLessonForm(false);
                  setEditingLesson(null);
                }}
                groups={groups}
                teachers={teachers}
                subjects={subjects}
                initialData={editingLesson}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TrelloSchedule;
