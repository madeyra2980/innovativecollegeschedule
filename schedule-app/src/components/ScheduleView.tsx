import React from 'react';
import { Schedule, Lesson, Group, Teacher, Subject } from '../types';

type ScheduleItem = Schedule | Lesson;

interface ScheduleViewProps {
  schedules: Schedule[];
  lessons: Lesson[];
  groups: Group[];
  teachers: Teacher[];
  subjects: Subject[];
  filterType: 'group' | 'teacher';
  selectedItem: any;
}

const ScheduleView: React.FC<ScheduleViewProps> = ({ 
  schedules, 
  lessons, 
  groups,
  teachers,
  subjects,
  filterType, 
  selectedItem 
}) => {
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

  const getScheduleForDay = (day: number) => {
    return schedules.filter(schedule => schedule.day_of_week === day);
  };

  const getLessonsForDay = (day: number) => {
    return lessons.filter(lesson => {
      if (!lesson.date) return false;
      const lessonDate = new Date(lesson.date);
      const dayOfWeek = lessonDate.getDay() === 0 ? 7 : lessonDate.getDay();
      return dayOfWeek === day;
    });
  };

  const getScheduleForDayAndShift = (day: number, shift: number) => {
    return schedules.filter(schedule => 
      schedule.day_of_week === day && schedule.shift === shift
    );
  };

  const getLessonsForDayAndShift = (day: number, shift: number) => {
    return lessons.filter(lesson => {
      if (!lesson.date) return false;
      const lessonDate = new Date(lesson.date);
      const dayOfWeek = lessonDate.getDay() === 0 ? 7 : lessonDate.getDay();
      return dayOfWeek === day && lesson.shift === shift;
    });
  };

  const getGroupName = (groupId: string) => {
    if (!groupId) return 'Неизвестная группа';
    const group = groups.find(g => g.id === groupId);
    return group ? group.name : 'Неизвестная группа';
  };

  const getTeacherName = (teacherId: string) => {
    if (!teacherId) return 'Неизвестный преподаватель';
    const teacher = teachers.find(t => t.id === teacherId);
    return teacher ? `${teacher.first_name} ${teacher.last_name}` : 'Неизвестный преподаватель';
  };

  const getSubjectName = (subjectId: string) => {
    if (!subjectId) return 'Неизвестный предмет';
    const subject = subjects.find(s => s.id === subjectId);
    if (!subject) {
      console.warn('Subject not found for ID:', subjectId, 'Available subjects:', subjects.map(s => ({ id: s.id, name: s.name })));
    }
    return subject ? subject.name : 'Неизвестный предмет';
  };

  const getSubjectCode = (subjectId: string) => {
    if (!subjectId) return 'Н/Д';
    const subject = subjects.find(s => s.id === subjectId);
    return subject ? subject.code : 'Н/Д';
  };

  const hasAnySchedule = schedules.length > 0 || lessons.length > 0;

  return (
    <div className="schedule-view">
      <div className="schedule-header">
        <h2>
          Расписание {filterType === 'group' ? 'группы' : 'преподавателя'}: 
          {filterType === 'group' ? 
            (selectedItem?.name || 'Не выбрано') : 
            (selectedItem ? `${selectedItem.first_name} ${selectedItem.last_name}` : 'Не выбрано')
          }
        </h2>
      </div>

      {!hasAnySchedule ? (
        <div className="empty-schedule">
          <div className="empty-icon"></div>
          <h3>Расписание не найдено</h3>
          <p>
            {filterType === 'group' 
              ? 'Для выбранной группы пока нет расписания' 
              : 'Для выбранного преподавателя пока нет расписания'
            }
          </p>
          <small>Обратитесь к администратору для добавления расписания</small>
        </div>
      ) : (
        <div className="schedule-grid">
          {daysOfWeek.map((day) => {
            const daySchedules = getScheduleForDay(day.value);
            const dayLessons = getLessonsForDay(day.value);
            const hasDaySchedule = daySchedules.length > 0 || dayLessons.length > 0;
            
            return (
              <div key={day.value} className="day-column">
                <div className="day-header">
                  <h3>{day.label}</h3>
                  {hasDaySchedule && (
                    <span className="day-count">
                      {daySchedules.length + dayLessons.length} занятий
                    </span>
                  )}
                </div>
                
                {!hasDaySchedule ? (
                  <div className="empty-day">
                    <div className="empty-day-icon"></div>
                    <p>Выходной</p>
                  </div>
                ) : (
                  shifts.map((shift) => {
                    const daySchedules = getScheduleForDayAndShift(day.value, shift.value);
                    const dayLessons = getLessonsForDayAndShift(day.value, shift.value);
                    const allItems = [...daySchedules, ...dayLessons];
                    
                    return (
                      <div key={shift.value} className="shift-section">
                        <div className="shift-header">
                          <h4>{shift.label}</h4>
                          <span className="item-count">({allItems.length} занятий)</span>
                        </div>
                        
                        <div className="items-list">
                          {allItems.length === 0 ? (
                            <div className="empty-slot">
                              <span>Нет занятий</span>
                            </div>
                          ) : (
                            allItems.map((item, index) => (
                              <div key={`${item.id}-${index}`} className="schedule-item">
                                <div className="item-header">
                                  <span className="subject-code">
                                    {item.subject_id ? getSubjectCode(item.subject_id) : 'Урок'}
                                  </span>
                                  <span className="time">
                                    {item.start_time} - {item.end_time}
                                  </span>
                                </div>
                                
                                <div className="item-content">
                                  <div className="subject-name">
                                    {item.subject_id ? getSubjectName(item.subject_id) : 'Урок'}
                                  </div>
                                  
                                  {filterType === 'group' ? (
                                    <div className="teacher-name">
                                      {getTeacherName(item.teacher_id)}
                                    </div>
                                  ) : (
                                    <div className="group-name">
                                      {getGroupName(item.group_id)}
                                    </div>
                                  )}
                                  
                                  <div className="room">
                                    Аудитория: {item.room}
                                  </div>
                                  
                                  {item.description && (
                                    <div className="description">
                                      {item.description}
                                    </div>
                                  )}
                                </div>
                              </div>
                            ))
                          )}
                        </div>
                      </div>
                    );
                  })
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default ScheduleView;
