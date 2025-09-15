import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Группы
export const groupsApi = {
  getAll: () => api.get('/groups'),
};

// Преподаватели
export const teachersApi = {
  getAll: () => api.get('/teachers'),
  getSchedule: (iin: string) => api.get(`/teachers/${iin}/schedule`),
};

// Предметы
export const subjectsApi = {
  getAll: () => api.get('/subjects'),
};

// Студенты
export const studentsApi = {
  getSchedule: (iin: string) => api.get(`/students/${iin}/schedule`),
};

// Расписание
export const schedulesApi = {
  getAll: () => api.get('/schedules'),
  getByDay: (day: number) => api.get(`/schedules/day/${day}`),
};

// Уроки
export const lessonsApi = {
  getAll: (params = {}) => api.get('/lessons', { params }),
  getByDate: (date: string) => api.get(`/lessons/date/${date}`),
  getByGroup: (groupId: string) => api.get('/lessons', { params: { group_id: groupId } }),
  getByTeacher: (teacherId: string) => api.get('/lessons', { params: { teacher_id: teacherId } }),
  getByShift: (shift: number) => api.get('/lessons', { params: { shift } }),
  getByDateAndShift: (date: string, shift: number) => api.get('/lessons', { params: { date, shift } }),
};

export default api;
