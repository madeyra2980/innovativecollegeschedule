import React from 'react';
import { Group, Teacher, FilterType } from '../types';

interface FilterSelectorProps {
  groups: Group[];
  teachers: Teacher[];
  selectedGroup: Group | null;
  selectedTeacher: Teacher | null;
  onGroupSelect: (group: Group | null) => void;
  onTeacherSelect: (teacher: Teacher | null) => void;
  filterType: FilterType;
  onFilterTypeChange: (type: FilterType) => void;
  loading: boolean;
}

const FilterSelector: React.FC<FilterSelectorProps> = ({
  groups,
  teachers,
  selectedGroup,
  selectedTeacher,
  onGroupSelect,
  onTeacherSelect,
  filterType,
  onFilterTypeChange,
  loading
}) => {
  return (
    <div className="filter-selector">
      <div className="filter-header">
        <h2>Выберите фильтр</h2>
        
        <div className="filter-type-buttons">
          <button
            className={`filter-type-btn ${filterType === 'group' ? 'active' : ''}`}
            onClick={() => onFilterTypeChange('group')}
          >
            По группам
          </button>
          <button
            className={`filter-type-btn ${filterType === 'teacher' ? 'active' : ''}`}
            onClick={() => onFilterTypeChange('teacher')}
          >
            По преподавателям
          </button>
        </div>
      </div>

      {loading ? (
        <div className="loading">
          <div className="loading-spinner"></div>
          <p>Загрузка данных...</p>
        </div>
      ) : (
        <div className="filter-content">
          {filterType === 'group' ? (
            <div className="groups-list">
              <h3>Группы</h3>
              {groups.length === 0 ? (
                <div className="empty-state">
                  <div className="empty-icon"></div>
                  <p>Группы не найдены</p>
                  <small>Попробуйте обновить страницу или обратитесь к администратору</small>
                </div>
              ) : (
                <div className="items-grid">
                  {groups.map((group) => (
                    <div
                      key={group.id}
                      className={`item-card ${selectedGroup?.id === group.id ? 'selected' : ''}`}
                      onClick={() => onGroupSelect(selectedGroup?.id === group.id ? null : group)}
                    >
                      <div className="item-header">
                        <h4>{group.name}</h4>
                        {group.code && <span className="item-code">{group.code}</span>}
                      </div>
                      {group.description && (
                        <p className="item-description">{group.description}</p>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          ) : (
            <div className="teachers-list">
              <h3>Преподаватели</h3>
              {teachers.length === 0 ? (
                <div className="empty-state">
                  <div className="empty-icon">👨‍🏫</div>
                  <p>Преподаватели не найдены</p>
                  <small>Попробуйте обновить страницу или обратитесь к администратору</small>
                </div>
              ) : (
                <div className="items-grid">
                  {teachers.map((teacher) => (
                    <div
                      key={teacher.id}
                      className={`item-card ${selectedTeacher?.id === teacher.id ? 'selected' : ''}`}
                      onClick={() => onTeacherSelect(selectedTeacher?.id === teacher.id ? null : teacher)}
                    >
                      <div className="item-header">
                        <h4>{teacher.first_name} {teacher.last_name}</h4>
                      </div>
                      {teacher.subjects.length > 0 && (
                        <div className="subjects">
                          <strong>Предметы:</strong>
                          <div className="subjects-list">
                            {teacher.subjects.map((subject, index) => (
                              <span key={index} className="subject-tag">
                                {subject}
                              </span>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default FilterSelector;
