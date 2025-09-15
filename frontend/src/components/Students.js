import React, { useState, useEffect } from 'react';
import { studentsApi, groupsApi } from '../services/api';
import CardSelector from './CardSelector';

const Students = () => {
  const [students, setStudents] = useState([]);
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingStudent, setEditingStudent] = useState(null);
  const [formData, setFormData] = useState({
    iin: '',
    first_name: '',
    last_name: '',
    group_id: ''
  });

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [studentsResponse, groupsResponse] = await Promise.all([
        studentsApi.getAll(),
        groupsApi.getAll()
      ]);
      setStudents(studentsResponse.data || []);
      setGroups(groupsResponse.data || []);
      setError(null);
    } catch (err) {
      setError('Ошибка загрузки данных: ' + (err.response?.data?.error || err.message));
      setStudents([]);
      setGroups([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingStudent) {
        await studentsApi.update(editingStudent.id, formData);
      } else {
        await studentsApi.create(formData);
      }
      setShowModal(false);
      setEditingStudent(null);
      setFormData({ iin: '', first_name: '', last_name: '', group_id: '' });
      fetchData();
    } catch (err) {
      setError('Ошибка сохранения: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleEdit = (student) => {
    setEditingStudent(student);
    setFormData({
      iin: student.iin,
      first_name: student.first_name,
      last_name: student.last_name,
      group_id: student.group_id
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить этого студента?')) {
      try {
        await studentsApi.delete(id);
        fetchData();
      } catch (err) {
        setError('Ошибка удаления: ' + (err.response?.data?.error || err.message));
      }
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingStudent(null);
    setFormData({ iin: '', first_name: '', last_name: '', group_id: '' });
  };

  const getGroupName = (groupId) => {
    const group = groups.find(g => g.id === groupId);
    return group ? group.name : 'Неизвестная группа';
  };

  if (loading) {
    return <div className="loading">Загрузка студентов...</div>;
  }

  return (
    <div>
      <div className="card">
        <h2>Студенты</h2>
        {error && <div className="error">{error}</div>}
        
        <button 
          className="btn"
          onClick={() => setShowModal(true)}
        >
          Добавить студента
        </button>

        {students && students.length > 0 ? (
          <table className="table">
            <thead>
              <tr>
                <th>ИИН</th>
                <th>Имя</th>
                <th>Фамилия</th>
                <th>Группа</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {students.map((student) => (
                <tr key={student.id}>
                  <td>{student.iin}</td>
                  <td>{student.first_name}</td>
                  <td>{student.last_name}</td>
                  <td>{getGroupName(student.group_id)}</td>
                  <td>
                    <div className="btn-group">
                      <button 
                        className="btn btn-secondary btn-small"
                        onClick={() => handleEdit(student)}
                      >
                        Редактировать
                      </button>
                      <button 
                        className="btn btn-danger btn-small"
                        onClick={() => handleDelete(student.id)}
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
            <p>Список студентов пуст</p>
            <p>Добавьте первого студента, нажав кнопку "Добавить студента"</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h3>{editingStudent ? 'Редактировать студента' : 'Добавить студента'}</h3>
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
                <label>Группа:</label>
                <CardSelector
                  options={groups}
                  value={formData.group_id}
                  onChange={(value) => setFormData({...formData, group_id: value})}
                  titleKey="name"
                  subtitleKey="code"
                />
              </div>
              <div className="btn-group">
                <button type="submit" className="btn">
                  {editingStudent ? 'Сохранить' : 'Создать'}
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

export default Students;
