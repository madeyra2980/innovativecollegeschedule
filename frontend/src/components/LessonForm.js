import React, { useState } from 'react';
import CardSelector from './CardSelector';

const LessonForm = ({ onSubmit, onCancel, groups, teachers, subjects, initialData = null }) => {
  const [formData, setFormData] = useState({
    group_id: initialData?.group_id || '',
    teacher_id: initialData?.teacher_id || '',
    subject_id: initialData?.subject_id || '',
    room: initialData?.room || '',
    description: initialData?.description || ''
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  return (
    <div className="lesson-form">
      <h3>{initialData ? 'Редактировать урок' : 'Создать урок'}</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Группа:</label>
          <CardSelector
            options={groups}
            value={formData.group_id}
            onChange={(value) => setFormData({...formData, group_id: value})}
            titleKey="name"
            subtitleKey="code"
          />
        </div>

        <div className="form-group">
          <label>Преподаватель:</label>
          <CardSelector
            options={teachers.map(teacher => ({
              ...teacher,
              name: `${teacher.first_name} ${teacher.last_name}`
            }))}
            value={formData.teacher_id}
            onChange={(value) => setFormData({...formData, teacher_id: value})}
            titleKey="name"
            subtitleKey="specialization"
          />
        </div>

        <div className="form-group">
          <label>Предмет:</label>
          <CardSelector
            options={subjects}
            value={formData.subject_id}
            onChange={(value) => setFormData({...formData, subject_id: value})}
            titleKey="name"
            subtitleKey="code"
          />
        </div>

        <div className="form-group">
          <label>Аудитория:</label>
          <input
            type="text"
            name="room"
            value={formData.room}
            onChange={handleChange}
            required
            placeholder="Номер аудитории"
          />
        </div>

        <div className="form-info">
          <p><strong>Информация:</strong> Дата и время урока будут установлены автоматически при перетаскивании в расписание.</p>
        </div>

        <div className="form-group">
          <label>Описание урока:</label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows="3"
            placeholder="Описание урока или темы"
          />
        </div>

        <div className="form-actions">
          <button type="submit" className="btn btn-primary">
            {initialData ? 'Сохранить' : 'Создать урок'}
          </button>
          <button type="button" className="btn btn-secondary" onClick={onCancel}>
            Отмена
          </button>
        </div>
      </form>
    </div>
  );
};

export default LessonForm;
