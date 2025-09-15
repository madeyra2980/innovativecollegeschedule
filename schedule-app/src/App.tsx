import React, { useState, useEffect } from 'react';
import { groupsApi, teachersApi, studentsApi, schedulesApi, lessonsApi, subjectsApi } from './services/api';
import { Group, Teacher, Schedule, Lesson, FilterType, StudentScheduleResponse, TeacherScheduleResponse, Subject } from './types';
import FilterSelector from './components/FilterSelector';
import ScheduleView from './components/ScheduleView';
import './App.css';

const App: React.FC = () => {
  const [groups, setGroups] = useState<Group[]>([]);
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [schedules, setSchedules] = useState<Schedule[]>([]);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [selectedTeacher, setSelectedTeacher] = useState<Teacher | null>(null);
  const [filterType, setFilterType] = useState<FilterType>('group');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    if (selectedGroup || selectedTeacher) {
      fetchScheduleData();
    }
  }, [selectedGroup, selectedTeacher, filterType]);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [groupsResponse, teachersResponse, subjectsResponse] = await Promise.all([
        groupsApi.getAll(),
        teachersApi.getAll(),
        subjectsApi.getAll()
      ]);
      
      setGroups(groupsResponse.data || []);
      setTeachers(teachersResponse.data || []);
      setSubjects(subjectsResponse.data || []);
      
      console.log('Loaded data:', {
        groups: groupsResponse.data?.length || 0,
        teachers: teachersResponse.data?.length || 0,
        subjects: subjectsResponse.data?.length || 0
      });
    } catch (err: any) {
      setError('Ошибка загрузки данных: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  const fetchScheduleData = async () => {
    try {
      setError(null);
      
      if (filterType === 'group' && selectedGroup) {
        // Получаем расписание группы
        const schedulesResponse = await schedulesApi.getAll();
        const allSchedules = schedulesResponse.data || [];
        const groupSchedules = allSchedules.filter((schedule: Schedule) => schedule.group_id === selectedGroup.id);
        
        // Получаем уроки группы
        const lessonsResponse = await lessonsApi.getByGroup(selectedGroup.id);
        const groupLessons = lessonsResponse.data || [];
        
        setSchedules(groupSchedules);
        setLessons(groupLessons);
      } else if (filterType === 'teacher' && selectedTeacher) {
        // Получаем расписание преподавателя
        const schedulesResponse = await schedulesApi.getAll();
        const allSchedules = schedulesResponse.data || [];
        const teacherSchedules = allSchedules.filter((schedule: Schedule) => schedule.teacher_id === selectedTeacher.id);
        
        // Получаем уроки преподавателя
        const lessonsResponse = await lessonsApi.getByTeacher(selectedTeacher.id);
        const teacherLessons = lessonsResponse.data || [];
        
        setSchedules(teacherSchedules);
        setLessons(teacherLessons);
      }
    } catch (err: any) {
      setError('Ошибка загрузки расписания: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleGroupSelect = (group: Group | null) => {
    setSelectedGroup(group);
    setSelectedTeacher(null);
  };

  const handleTeacherSelect = (teacher: Teacher | null) => {
    setSelectedTeacher(teacher);
    setSelectedGroup(null);
  };

  const handleFilterTypeChange = (type: FilterType) => {
    setFilterType(type);
    setSelectedGroup(null);
    setSelectedTeacher(null);
    setSchedules([]);
    setLessons([]);
  };

  return (
    <div className="app">
      <header className="app-header">
        <div className="container">
          <div className="header-content">
            <div className="logo">
              <img 
                src="https://innovativecollege.kz/wp-content/uploads/2025/03/cropped-img-3928png-300x142.png" 
                alt="Инновационный колледж" 
                className="logo-image"
              />
              <span className="logo-subtitle">Система управления учебным процессом</span>
            </div>
            <div className="header-description">
              <h1>Расписание колледжа</h1>
              <p>Выберите группу или преподавателя для просмотра расписания</p>
            </div>
          </div>
        </div>
      </header>

      <main className="app-main">
        <div className="container">
          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <FilterSelector
            groups={groups}
            teachers={teachers}
            selectedGroup={selectedGroup}
            selectedTeacher={selectedTeacher}
            onGroupSelect={handleGroupSelect}
            onTeacherSelect={handleTeacherSelect}
            filterType={filterType}
            onFilterTypeChange={handleFilterTypeChange}
            loading={loading}
          />

          {(selectedGroup || selectedTeacher) && (
            <ScheduleView
              schedules={schedules}
              lessons={lessons}
              groups={groups}
              teachers={teachers}
              subjects={subjects}
              filterType={filterType}
              selectedItem={selectedGroup || selectedTeacher}
            />
          )}
        </div>
      </main>

      <footer className="footer">
        <div className="container">
          <div className="footer-content">
            <div className="footer-section">
              <h3>Innovative College</h3>
              <p>Миссия колледжа: "Подготовка конкурентоспособных специалистов, востребованных на рынке труда, способных внести вклад в развитие экономики нового Казахстана в современном мире"</p>
              <div className="footer-info">
                <span>Email: semey@innovativecollege.kz</span>
                <span>Телефон: +7 775 399 2283</span>
                <span>Адрес: ул. Жамбыла 20, г. Семей</span>
                <span>Instagram: @innovativecollege</span>
                <span>WhatsApp: +7(747)1849036</span>
              </div>
            </div>
            
            <div className="footer-section">
              <h4>Быстрые ссылки</h4>
              <div className="footer-links">
                <span>Группы</span>
                <span>Студенты</span>
                <span>Преподаватели</span>
                <span>Предметы</span>
              </div>
            </div>

            <div className="footer-section">
              <h4>Специальности</h4>
              <div className="footer-links">
                <span>Программное обеспечение</span>
                <span>Пожарная безопасность</span>
                <span>Организация дорожного движения</span>
                <span>Социальная работа</span>
                <span>Правоведение</span>
                <span>Учет и аудит</span>
              </div>
            </div>

            <div className="footer-section">
              <h4>Рабочий график</h4>
              <div className="footer-links">
                <span>Понедельник-Пятница</span>
                <span>8:00 — 17:00</span>
                <span>Прием на обучение</span>
                <span>Виртуальная приемная комиссия</span>
                <span>График работы приемной комиссии</span>
              </div>
            </div>
          </div>
          
          <div className="footer-bottom">
            <p>&copy; 2025 Innovative College | Powered by Innovative College.</p>
            <div className="footer-tech">
              <span>React & Go</span>
              <span>Безопасность данных</span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;