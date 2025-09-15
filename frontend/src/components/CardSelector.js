import React from 'react';

const CardSelector = ({ 
  options, 
  value, 
  onChange, 
  titleKey = 'name', 
  subtitleKey = 'code',
  placeholder = 'Выберите опцию'
}) => {
  return (
    <div className="card-selector">
      {options.map((option) => (
        <div
          key={option.id}
          className={`card-option ${value === option.id ? 'selected' : ''}`}
          onClick={() => onChange(option.id)}
        >
          <div className="card-option-check"></div>
          <div className="card-option-title">
            {option[titleKey]}
          </div>
          {option[subtitleKey] && (
            <div className="card-option-subtitle">
              {option[subtitleKey]}
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default CardSelector;
