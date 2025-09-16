import axios from 'axios';

const API_BASE_URL = process.env.NODE_ENV === 'production'
  ? '/api/v1'  // В продакшене используем относительный путь (nginx прокси)
  : 'http://localhost:8080/api/v1';  // В разработке используем localhost

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
