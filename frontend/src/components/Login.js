import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const Login = ({ onLogin }) => {
  const [credentials, setCredentials] = useState({
    username: '',
    password: ''
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setCredentials(prev => ({
      ...prev,
      [name]: value
    }));
    // Очищаем ошибку при изменении полей
    if (error) setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    // Имитация задержки для реалистичности
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Проверка логина и пароля
    if (credentials.username === 'user' && credentials.password === '123') {
      // Успешная авторизация
      localStorage.setItem('isAuthenticated', 'true');
      localStorage.setItem('username', credentials.username);
      onLogin(true);
      navigate('/');
    } else {
      setError('Неверный логин или пароль');
    }
    
    setIsLoading(false);
  };

  return (
    <div className="login-container">
      <div className="login-background">
        <div className="login-shapes">
          <div className="shape shape-1"></div>
          <div className="shape shape-2"></div>
          <div className="shape shape-3"></div>
        </div>
      </div>
      
      <div className="login-card">
        <div className="login-header">
          <div className="login-logo">
            <img 
              src="https://innovativecollege.kz/wp-content/uploads/2025/03/cropped-img-3928png-300x142.png" 
              alt="Innovative College"
              className="logo-image"
            />
          </div>
          <h1 className="login-title">Добро пожаловать</h1>
          <p className="login-subtitle">Система управления учебным процессом</p>
        </div>

        <form className="login-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username" className="form-label">
              Логин
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={credentials.username}
              onChange={handleChange}
              className="form-input"
              placeholder="Введите логин"
              required
              disabled={isLoading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="password" className="form-label">
              Пароль
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={credentials.password}
              onChange={handleChange}
              className="form-input"
              placeholder="Введите пароль"
              required
              disabled={isLoading}
            />
          </div>

          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <button 
            type="submit" 
            className="login-button"
            disabled={isLoading}
          >
            {isLoading ? (
              <span className="loading-spinner"></span>
            ) : (
              'Войти в систему'
            )}
          </button>
        </form>

      </div>
    </div>
  );
};

export default Login;
