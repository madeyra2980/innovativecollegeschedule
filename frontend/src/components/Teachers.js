import React, { useState, useEffect } from 'react';
import { teachersApi } from '../services/api';

const Teachers = () => {
  const [teachers, setTeachers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingTeacher, setEditingTeacher] = useState(null);
  const [formData, setFormData] = useState({
    iin: '',
    first_name: '',
    last_name: '',
    subjects: ''
  });

  useEffect(() => {
    fetchTeachers();
  }, []);

  const fetchTeachers = async () => {
    try {
      setLoading(true);
      const response = await teachersApi.getAll();
      setTeachers(response.data || []);
      setError(null);
    } catch (err) {
      setError('Ошибка загрузки преподавателей: ' + (err.response?.data?.error || err.message));
      setTeachers([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const subjectsArray = formData.subjects
        .split(',')
        .map(s => s.trim())
        .filter(s => s.length > 0);

      const dataToSubmit = {
        ...formData,
        subjects: subjectsArray
      };

      if (editingTeacher) {
        await teachersApi.update(editingTeacher.id, dataToSubmit);
      } else {
        await teachersApi.create(dataToSubmit);
      }
      setShowModal(false);
      setEditingTeacher(null);
      setFormData({ iin: '', first_name: '', last_name: '', subjects: '' });
      fetchTeachers();
    } catch (err) {
      setError('Ошибка сохранения: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleEdit = (teacher) => {
    setEditingTeacher(teacher);
    setFormData({
      iin: teacher.iin,
      first_name: teacher.first_name,
      last_name: teacher.last_name,
      subjects: teacher.subjects.join(', ')
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить этого преподавателя?')) {
      try {
        await teachersApi.delete(id);
        fetchTeachers();
      } catch (err) {
        setError('Ошибка удаления: ' + (err.response?.data?.error || err.message));
      }
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingTeacher(null);
    setFormData({ iin: '', first_name: '', last_name: '', subjects: '' });
  };

  if (loading) {
    return <div className="loading">Загрузка преподавателей...</div>;
  }

  return (
    <div>
      <div className="card">
        <h2>Преподаватели</h2>
        {error && <div className="error">{error}</div>}
        
        <button 
          className="btn"
          onClick={() => setShowModal(true)}
        >
          Добавить преподавателя
        </button>

        {teachers && teachers.length > 0 ? (
          <table className="table">
            <thead>
              <tr>
                <th>ИИН</th>
                <th>Имя</th>
                <th>Фамилия</th>
                <th>Предметы</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {teachers.map((teacher) => (
                <tr key={teacher.id}>
                  <td>{teacher.iin}</td>
                  <td>{teacher.first_name}</td>
                  <td>{teacher.last_name}</td>
                  <td>
                    {teacher.subjects.length > 0 ? (
                      <div>
                        {teacher.subjects.map((subject, index) => (
                          <span key={index} className="schedule-item">
                            {subject}
                          </span>
                        ))}
                      </div>
                    ) : '-'}
                  </td>
                  <td>
                    <div className="btn-group">
                      <button 
                        className="btn btn-secondary btn-small"
                        onClick={() => handleEdit(teacher)}
                      >
                        Редактировать
                      </button>
                      <button 
                        className="btn btn-danger btn-small"
                        onClick={() => handleDelete(teacher.id)}
                      >
                        Удалить
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div className="empty-state">
            <p>Список преподавателей пуст</p>
            <p>Добавьте первого преподавателя, нажав кнопку "Добавить преподавателя"</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h3>{editingTeacher ? 'Редактировать преподавателя' : 'Добавить преподавателя'}</h3>
              <button className="close-btn" onClick={handleCloseModal}>×</button>
            </div>
            <div className="modal-content">
              <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>ИИН:</label>
                <input
                  type="text"
                  value={formData.iin}
                  onChange={(e) => setFormData({...formData, iin: e.target.value})}
                  required
                  maxLength="12"
                />
              </div>
              <div className="form-group">
                <label>Имя:</label>
                <input
                  type="text"
                  value={formData.first_name}
                  onChange={(e) => setFormData({...formData, first_name: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Фамилия:</label>
                <input
                  type="text"
                  value={formData.last_name}
                  onChange={(e) => setFormData({...formData, last_name: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Предметы (через запятую):</label>
                <textarea
                  value={formData.subjects}
                  onChange={(e) => setFormData({...formData, subjects: e.target.value})}
                  rows="3"
                  placeholder="Например: Программирование, Базы данных, Веб-разработка"
                />
              </div>
              <div className="btn-group">
                <button type="submit" className="btn">
                  {editingTeacher ? 'Сохранить' : 'Создать'}
                </button>
                <button type="button" className="btn btn-secondary" onClick={handleCloseModal}>
                  Отмена
                </button>
              </div>
            </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Teachers;
