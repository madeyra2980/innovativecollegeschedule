import axios from 'axios';

const API_BASE_URL = 'http://localhost:9090/api/v1';
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000, // 10 секунд таймаут
});

// Интерцептор для обработки ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.code === 'ECONNABORTED') {
      error.message = 'Превышено время ожидания запроса';
    } else if (error.code === 'ERR_NETWORK') {
      error.message = 'Ошибка сети. Проверьте подключение к серверу';
    } else if (error.response?.status === 500) {
      error.message = 'Внутренняя ошибка сервера';
    } else if (error.response?.status === 404) {
      error.message = 'Ресурс не найден';
    }
    return Promise.reject(error);
  }
);

// Группы
export const groupsApi = {
  getAll: () => api.get('/groups'),
  create: (data) => api.post('/groups', data),
  update: (id, data) => api.put(`/groups/${id}`, data),
  delete: (id) => api.delete(`/groups/${id}`),
};

// Предметы
export const subjectsApi = {
  getAll: () => api.get('/subjects'),
  create: (data) => api.post('/subjects', data),
  update: (id, data) => api.put(`/subjects/${id}`, data),
  delete: (id) => api.delete(`/subjects/${id}`),
};

// Студенты
export const studentsApi = {
  getAll: () => api.get('/students'),
  create: (data) => api.post('/students', data),
  update: (id, data) => api.put(`/students/${id}`, data),
  delete: (id) => api.delete(`/students/${id}`),
  getSchedule: (iin) => api.get(`/students/${iin}/schedule`),
};

// Преподаватели
export const teachersApi = {
  getAll: () => api.get('/teachers'),
  create: (data) => api.post('/teachers', data),
  update: (id, data) => api.put(`/teachers/${id}`, data),
  delete: (id) => api.delete(`/teachers/${id}`),
  getSchedule: (iin) => api.get(`/teachers/${iin}/schedule`),
};

// Расписание
export const schedulesApi = {
  getAll: () => api.get('/schedules'),
  create: (data) => api.post('/schedules', data),
  update: (id, data) => api.put(`/schedules/${id}`, data),
  delete: (id) => api.delete(`/schedules/${id}`),
  getByDay: (day) => api.get(`/schedules/day/${day}`),
};

// Календарь (Уроки)
export const lessonsApi = {
  getAll: (params = {}) => api.get('/lessons', { params }),
  getAvailable: () => api.get('/lessons/available'),
  create: (data) => api.post('/lessons', data),
  update: (id, data) => api.put(`/lessons/${id}`, data),
  delete: (id) => api.delete(`/lessons/${id}`),
  getByDate: (date) => api.get(`/lessons/date/${date}`),
  getByGroup: (groupId) => api.get('/lessons', { params: { group_id: groupId } }),
  getByTeacher: (teacherId) => api.get('/lessons', { params: { teacher_id: teacherId } }),
  getByShift: (shift) => api.get('/lessons', { params: { shift } }),
  getByDateAndShift: (date, shift) => api.get('/lessons', { params: { date, shift } }),
};

// Временные слоты
export const timeSlotsApi = {
  getAll: (params = {}) => api.get('/time-slots', { params }),
  getById: (id) => api.get(`/time-slots/${id}`),
  create: (data) => api.post('/time-slots', data),
  update: (id, data) => api.put(`/time-slots/${id}`, data),
  delete: (id) => api.delete(`/time-slots/${id}`),
  getByShift: (shift) => api.get('/time-slots', { params: { shift } }),
  getActive: () => api.get('/time-slots', { params: { is_active: true } }),
};

// Статистика
export const statisticsApi = {
  getLessonStatistics: (params = {}) => api.get('/statistics/lessons', { params }),
};

export default api;
