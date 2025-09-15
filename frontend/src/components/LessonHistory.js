import React, { useState, useEffect } from 'react';
import { lessonsApi, groupsApi, teachersApi, subjectsApi } from '../services/api';
import CardSelector from './CardSelector';

const LessonHistory = () => {
  const [lessons, setLessons] = useState([]);
  const [groups, setGroups] = useState([]);
  const [teachers, setTeachers] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    startDate: '',
    endDate: '',
    groupId: '',
    teacherId: '',
    shift: ''
  });
  const [filtersExpanded, setFiltersExpanded] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    fetchLessons();
  }, [filters]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [groupsResponse, teachersResponse, subjectsResponse] = await Promise.all([
        groupsApi.getAll(),
        teachersApi.getAll(),
        subjectsApi.getAll()
      ]);
      setGroups(groupsResponse.data || []);
      setTeachers(teachersResponse.data || []);
      setSubjects(subjectsResponse.data || []);
      setError(null);
    } catch (err) {
      setError('Ошибка загрузки данных: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  const fetchLessons = async () => {
    try {
      const params = {};
      
      if (filters.startDate && filters.endDate) {
        params.start_date = filters.startDate;
        params.end_date = filters.endDate;
      } else if (filters.startDate) {
        params.start_date = filters.startDate;
        params.end_date = filters.startDate;
      }
      
      if (filters.groupId) params.group_id = filters.groupId;
      if (filters.teacherId) params.teacher_id = filters.teacherId;
      if (filters.shift) params.shift = parseInt(filters.shift);

      const response = await lessonsApi.getAll(params);
      setLessons(response.data || []);
    } catch (err) {
      setError('Ошибка загрузки уроков: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleFilterChange = (field, value) => {
    setFilters(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const clearFilters = () => {
    setFilters({
      startDate: '',
      endDate: '',
      groupId: '',
      teacherId: '',
      shift: ''
    });
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
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
    return subject ? subject.code : 'Н/Д';
  };

  const getShiftName = (shift) => {
    return shift === 1 ? 'Первая смена' : shift === 2 ? 'Вторая смена' : 'Не указана';
  };

  if (loading) {
    return <div className="loading">Загрузка истории уроков...</div>;
  }

  return (
    <div className="lesson-history">
      <div className="card">
        <div className="card-header">
          <h2>История уроков</h2>
          <p>Просмотр проведенных уроков по датам и фильтрам</p>
        </div>

        {/* Фильтры */}
      <div className="filters-section">
        <div className="filters-header">
          <h3>Фильтры</h3>
          <button 
            className="btn btn-secondary filters-toggle"
            onClick={() => setFiltersExpanded(!filtersExpanded)}
          >
            {filtersExpanded ? '▼ Скрыть фильтры' : '▶ Показать фильтры'}
          </button>
        </div>
        <div className={`filters-grid ${filtersExpanded ? 'expanded' : 'collapsed'}`}>
            <div className="filter-group">
              <label>Период с:</label>
              <input
                type="date"
                value={filters.startDate}
                onChange={(e) => handleFilterChange('startDate', e.target.value)}
                className="form-control"
              />
            </div>
            
            <div className="filter-group">
              <label>Период по:</label>
              <input
                type="date"
                value={filters.endDate}
                onChange={(e) => handleFilterChange('endDate', e.target.value)}
                className="form-control"
              />
            </div>

            <div className="filter-group">
              <label>Группа:</label>
              <CardSelector
                options={[{id: '', name: 'Все группы'}, ...groups]}
                value={filters.groupId}
                onChange={(value) => handleFilterChange('groupId', value)}
                titleKey="name"
                subtitleKey="code"
              />
            </div>

            <div className="filter-group">
              <label>Преподаватель:</label>
              <CardSelector
                options={[
                  {id: '', name: 'Все преподаватели'}, 
                  ...teachers.map(teacher => ({
                    ...teacher,
                    name: `${teacher.first_name} ${teacher.last_name}`
                  }))
                ]}
                value={filters.teacherId}
                onChange={(value) => handleFilterChange('teacherId', value)}
                titleKey="name"
                subtitleKey="specialization"
              />
            </div>

            <div className="filter-group">
              <label>Смена:</label>
              <CardSelector
                options={[
                  {id: '', name: 'Все смены'},
                  {id: '1', name: 'Первая смена'},
                  {id: '2', name: 'Вторая смена'}
                ]}
                value={filters.shift}
                onChange={(value) => handleFilterChange('shift', value)}
                titleKey="name"
              />
            </div>

            <div className="filter-group">
              <button onClick={clearFilters} className="btn btn-secondary">
                Очистить фильтры
              </button>
            </div>
          </div>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* Результаты */}
        <div className="results-section">
          <div className="results-header">
            <h3>Результаты ({lessons.length} уроков)</h3>
          </div>

          {lessons.length === 0 ? (
            <div className="empty-state">
              <p>Уроки не найдены</p>
              <p>Попробуйте изменить фильтры поиска</p>
            </div>
          ) : (
            <div className="lessons-list">
              {lessons.map(lesson => (
                <div key={lesson.id} className="lesson-card">
                  <div className="lesson-header">
                    <div className="lesson-date">
                      <strong>{formatDate(lesson.date)}</strong>
                      <span className="lesson-time">
                        {lesson.start_time} - {lesson.end_time}
                      </span>
                    </div>
                    <div className="lesson-shift">
                      {getShiftName(lesson.shift)}
                    </div>
                  </div>
                  
                  <div className="lesson-content">
                    <div className="lesson-info">
                      <div className="info-item">
                        <strong>Группа:</strong> {getGroupName(lesson.group_id)}
                      </div>
                      <div className="info-item">
                        <strong>Преподаватель:</strong> {getTeacherName(lesson.teacher_id)}
                      </div>
                      <div className="info-item">
                        <strong>Предмет:</strong> {getSubjectName(lesson.subject_id)}
                        <span className="subject-code">({getSubjectCode(lesson.subject_id)})</span>
                      </div>
                      <div className="info-item">
                        <strong>Аудитория:</strong> {lesson.room}
                      </div>
                      {lesson.description && (
                        <div className="info-item">
                          <strong>Описание:</strong> {lesson.description}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default LessonHistory;
