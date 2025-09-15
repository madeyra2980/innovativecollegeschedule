import React, { useState, useEffect } from 'react';
import { Routes, Route, Link, useLocation } from 'react-router-dom';
import Groups from './components/Groups';
import Students from './components/Students';
import Teachers from './components/Teachers';
import Subjects from './components/Subjects';
import Schedule from './components/Schedule';
import LessonHistory from './components/LessonHistory';
import Login from './components/Login';

function App() {
  const location = useLocation();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isHeaderVisible, setIsHeaderVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return localStorage.getItem('isAuthenticated') === 'true';
  });

  const toggleMenu = () => {
    console.log('Toggle menu clicked, current state:', isMenuOpen);
    setIsMenuOpen(!isMenuOpen);
  };

  const handleLogin = (success) => {
    setIsAuthenticated(success);
  };

  const handleLogout = () => {
    localStorage.removeItem('isAuthenticated');
    localStorage.removeItem('username');
    setIsAuthenticated(false);
    setIsMenuOpen(false);
  };

  const closeMenu = () => {
    setIsMenuOpen(false);
  };

  // Закрытие меню при клике вне его области
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (isMenuOpen && !event.target.closest('.nav') && !event.target.closest('.burger-menu')) {
        closeMenu();
      }
    };

    if (isMenuOpen) {
      document.addEventListener('click', handleClickOutside);
      document.body.style.overflow = 'hidden'; // Блокируем скролл
    } else {
      document.body.style.overflow = 'unset'; // Разблокируем скролл
    }

    return () => {
      document.removeEventListener('click', handleClickOutside);
      document.body.style.overflow = 'unset';
    };
  }, [isMenuOpen]);

  // Скрытие/показ header при скролле
  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      
      if (currentScrollY < lastScrollY || currentScrollY < 10) {
        // Показываем header при скролле вверх или в самом верху
        setIsHeaderVisible(true);
      } else if (currentScrollY > lastScrollY && currentScrollY > 100) {
        // Скрываем header при скролле вниз (только если проскроллили больше 100px)
        setIsHeaderVisible(false);
      }
      
      setLastScrollY(currentScrollY);
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, [lastScrollY]);

  // Если пользователь не авторизован, показываем страницу входа
  if (!isAuthenticated) {
    return <Login onLogin={handleLogin} />;
  }

  return (
    <div className="App">
      {/* Фиксированное бургер-меню */}
      <button
        className={`burger-menu ${isMenuOpen ? 'active' : ''}`}
        onClick={toggleMenu}
        aria-label={isMenuOpen ? "Закрыть меню" : "Открыть меню"}
        title={isMenuOpen ? "Закрыть меню" : "Открыть меню"}
      >
        <span></span>
        <span></span>
        <span></span>
      </button>

      <header className={`header ${!isHeaderVisible ? 'header-hidden' : ''}`}>
        <div className="container">
          <div className="header-content">
            <div className="logo">
              <img 
                src="https://innovativecollege.kz/wp-content/uploads/2025/03/cropped-img-3928png-300x142.png" 
                alt="Инновационный колледж" 
                className="logo-image"
              />
              <span className="logo-subtitle">Система управления учебным процессом</span>
            </div>
            
          </div>
          
        </div>

        {/* Burger Menu Navigation */}
        <nav className={`nav ${isMenuOpen ? 'active' : ''}`}>
              <Link 
                to="/groups" 
                className={location.pathname === '/groups' ? 'active' : ''}
                onClick={closeMenu}
              >
                Группы
              </Link>
              <Link 
                to="/students" 
                className={location.pathname === '/students' ? 'active' : ''}
                onClick={closeMenu}
              >
                Студенты
              </Link>
              <Link 
                to="/teachers" 
                className={location.pathname === '/teachers' ? 'active' : ''}
                onClick={closeMenu}
              >
                Преподаватели
              </Link>
              <Link 
                to="/subjects" 
                className={location.pathname === '/subjects' ? 'active' : ''}
                onClick={closeMenu}
              >
                Предметы
              </Link>
              <Link 
                to="/schedule" 
                className={location.pathname === '/schedule' ? 'active' : ''}
                onClick={closeMenu}
              >
                Расписание
              </Link>
              <Link 
                to="/history" 
                className={location.pathname === '/history' ? 'active' : ''}
                onClick={closeMenu}
              >
                История уроков
              </Link>
              <button className="logout-button" onClick={handleLogout}>
                Выйти
              </button>
            </nav>
      </header>

      <main className="main-content">
        <div className="container">
          <Routes>
            <Route path="/" element={<Groups />} />
            <Route path="/groups" element={<Groups />} />
            <Route path="/students" element={<Students />} />
            <Route path="/teachers" element={<Teachers />} />
            <Route path="/subjects" element={<Subjects />} />
            <Route path="/schedule" element={<Schedule />} />
            <Route path="/history" element={<LessonHistory />} />
          </Routes>
        </div>
      </main>

      <footer className="footer">
        <div className="container">
          <div className="footer-content">
            <div className="footer-section">
              <h3>Innovative College</h3>
              <p>Миссия колледжа: "Подготовка конкурентоспособных специалистов, востребованных на рынке труда, способных внести вклад в развитие экономики нового Казахстана в современном мире"</p>
              <div className="footer-info">
                <span>Email: semey@innovativecollege.kz</span>
                <span>Телефон: +7 775 399 2283</span>
                <span>Адрес: ул. Жамбыла 20, г. Семей</span>
                <span>Instagram: @innovativecollege</span>
                <span>WhatsApp: +7(747)1849036</span>
              </div>
            </div>
            
            <div className="footer-section">
              <h4>Быстрые ссылки</h4>
              <div className="footer-links">
                <Link to="/groups">Группы</Link>
                <Link to="/students">Студенты</Link>
                <Link to="/teachers">Преподаватели</Link>
                <Link to="/subjects">Предметы</Link>
              </div>
            </div>

            <div className="footer-section">
              <h4>Специальности</h4>
              <div className="footer-links">
                <span>Программное обеспечение</span>
                <span>Пожарная безопасность</span>
                <span>Организация дорожного движения</span>
                <span>Социальная работа</span>
                <span>Правоведение</span>
                <span>Учет и аудит</span>
              </div>
            </div>

            <div className="footer-section">
              <h4>Рабочий график</h4>
              <div className="footer-links">
                <span>Понедельник-Пятница</span>
                <span>8:00 — 17:00</span>
                <span>Прием на обучение</span>
                <span>Виртуальная приемная комиссия</span>
                <span>График работы приемной комиссии</span>
              </div>
            </div>
          </div>
          
          <div className="footer-bottom">
            <p>&copy; 2025 Innovative College | Powered by Innovative College.</p>
            <div className="footer-tech">
              <span>React & Go</span>
              <span>Безопасность данных</span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

export default App;
