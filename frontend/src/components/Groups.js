import React, { useState, useEffect } from 'react';
import { groupsApi } from '../services/api';

const Groups = () => {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingGroup, setEditingGroup] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: ''
  });

  useEffect(() => {
    fetchGroups();
  }, []);

  const fetchGroups = async () => {
    try {
      setLoading(true);
      const response = await groupsApi.getAll();
      setGroups(response.data || []);
      setError(null);
    } catch (err) {
      setError('Ошибка загрузки групп: ' + (err.response?.data?.error || err.message));
      setGroups([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingGroup) {
        await groupsApi.update(editingGroup.id, formData);
      } else {
        await groupsApi.create(formData);
      }
      setShowModal(false);
      setEditingGroup(null);
      setFormData({ name: '', description: '' });
      fetchGroups();
    } catch (err) {
      setError('Ошибка сохранения: ' + (err.response?.data?.error || err.message));
    }
  };

  const handleEdit = (group) => {
    setEditingGroup(group);
    setFormData({
      name: group.name,
      description: group.description || ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить эту группу?')) {
      try {
        await groupsApi.delete(id);
        fetchGroups();
      } catch (err) {
        setError('Ошибка удаления: ' + (err.response?.data?.error || err.message));
      }
    }
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingGroup(null);
    setFormData({ name: '', description: '' });
  };

  if (loading) {
    return <div className="loading">Загрузка групп...</div>;
  }

  return (
    <div>
      <div className="card">
        <h2>Группы</h2>
        {error && <div className="error">{error}</div>}
        
        <button 
          className="btn"
          onClick={() => setShowModal(true)}
        >
          Добавить группу
        </button>

        {groups && groups.length > 0 ? (
          <table className="table">
            <thead>
              <tr>
                <th>Название</th>
                <th>Описание</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {groups.map((group) => (
                <tr key={group.id}>
                  <td>{group.name}</td>
                  <td>{group.description || '-'}</td>
                  <td>
                    <div className="btn-group">
                      <button 
                        className="btn btn-secondary btn-small"
                        onClick={() => handleEdit(group)}
                      >
                        Редактировать
                      </button>
                      <button 
                        className="btn btn-danger btn-small"
                        onClick={() => handleDelete(group.id)}
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
            <p>Список групп пуст</p>
            <p>Добавьте первую группу, нажав кнопку "Добавить группу"</p>
          </div>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h3>{editingGroup ? 'Редактировать группу' : 'Добавить группу'}</h3>
              <button className="close-btn" onClick={handleCloseModal}>×</button>
            </div>
            <div className="modal-content">
              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <label>Название группы:</label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({...formData, name: e.target.value})}
                    required
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
                    {editingGroup ? 'Сохранить' : 'Создать'}
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

export default Groups;
