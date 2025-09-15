import React, { useState, useEffect } from 'react';
import { subjectsApi } from '../services/api';

const Subjects = () => {
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingSubject, setEditingSubject] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    code: '',
    description: ''
  });

  useEffect(() => {
    fetchSubjects();
  }, []);

  const fetchSubjects = async () => {
    try {
      setLoading(true);
      const response = await subjectsApi.getAll();
      setSubjects(response.data || []);
      setError(null);
    } catch (err) {
      setError('Ошибка загрузки предметов: ' + (err.response?.data?.error || err.message));
      setSubjects([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingSubject) {
        await subjectsApi.update(editingSubject.id, formData);
      } else {
        await subjectsApi.create(formData);
      }
      setShowModal(false);
      setEditingSubject(null);
      setFormData({ name: '', code: '', description: '' });
      fetchSubjects();
    } catch (err) {
      setError('Ошибка сохранения: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleEdit = (subject) => {
    setEditingSubject(subject);
    setFormData({
      name: subject.name,
      code: subject.code,
      description: subject.description || ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить этот предмет?')) {
      try {
        await subjectsApi.delete(id);
        fetchSubjects();
      } catch (err) {
        setError('Ошибка удаления: ' + (err.response?.data?.error || err.message));
      }
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingSubject(null);
    setFormData({ name: '', code: '', description: '' });
  };

  if (loading) {
    return <div className="loading">Загрузка предметов...</div>;
  }

  return (
    <div>
      <div className="card">
        <h2>Предметы</h2>
        {error && <div className="error">{error}</div>}
        
        <button 
          className="btn"
          onClick={() => setShowModal(true)}
        >
          Добавить предмет
        </button>

        {subjects && subjects.length > 0 ? (
          <table className="table subjects-table">
            <thead>
              <tr>
                <th>Название</th>
                <th>Код</th>
                <th>Описание</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {subjects.map((subject) => (
                <tr key={subject.id}>
                  <td>{subject.name}</td>
                  <td>{subject.code}</td>
                  <td>{subject.description || '-'}</td>
                  <td>
                    <div className="btn-group">
                      <button 
                        className="btn btn-secondary btn-small"
                        onClick={() => handleEdit(subject)}
                      >
                        Редактировать
                      </button>
                      <button 
                        className="btn btn-danger btn-small"
                        onClick={() => handleDelete(subject.id)}
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
            <p>Список предметов пуст</p>
            <p>Добавьте первый предмет, нажав кнопку "Добавить предмет"</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h3>{editingSubject ? 'Редактировать предмет' : 'Добавить предмет'}</h3>
              <button className="close-btn" onClick={handleCloseModal}>×</button>
            </div>
            <div className="modal-content">
              <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Название предмета:</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Код предмета:</label>
                <input
                  type="text"
                  value={formData.code}
                  onChange={(e) => setFormData({...formData, code: e.target.value})}
                  required
                  placeholder="Например: ОН 3.1"
                />
              </div>
              <div className="form-group">
                <label>Описание:</label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({...formData, description: e.target.value})}
                  rows="3"
                />
              </div>
              <div className="btn-group">
                <button type="submit" className="btn">
                  {editingSubject ? 'Сохранить' : 'Создать'}
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

export default Subjects;
