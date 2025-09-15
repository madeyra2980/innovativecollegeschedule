export interface Group {
  id: string;
  name: string;
  code?: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Teacher {
  id: string;
  iin: string;
  first_name: string;
  last_name: string;
  subjects: string[];
  created_at: string;
  updated_at: string;
}

export interface Student {
  id: string;
  iin: string;
  first_name: string;
  last_name: string;
  group_id: string;
  group?: Group;
  created_at: string;
  updated_at: string;
}

export interface Subject {
  id: string;
  name: string;
  code: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Schedule {
  id: string;
  group_id: string;
  group?: Group;
  teacher_id: string;
  teacher?: Teacher;
  subject_id: string;
  subject?: Subject;
  room: string;
  day_of_week: number;
  start_time: string;
  end_time: string;
  shift: number;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Lesson {
  id: string;
  group_id: string;
  group?: Group;
  teacher_id: string;
  teacher?: Teacher;
  subject_id: string;
  subject?: Subject;
  room: string;
  date?: string;
  start_time?: string;
  end_time?: string;
  shift?: number;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface StudentScheduleResponse {
  student: Student;
  schedules: Schedule[];
}

export interface TeacherScheduleResponse {
  teacher: Teacher;
  schedules: Schedule[];
}

export type FilterType = 'group' | 'teacher';
