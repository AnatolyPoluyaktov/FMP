import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { BarChart3, PieChart, TrendingUp, AlertTriangle } from 'lucide-react';

export const Navigation: React.FC = () => {
  const location = useLocation();

  const navItems = [
    { path: '/', label: 'Обзор', icon: BarChart3 },
    { path: '/monthly', label: 'Месячная аналитика', icon: TrendingUp },
    { path: '/categories', label: 'Категории', icon: PieChart },
    { path: '/limits', label: 'Лимиты', icon: AlertTriangle },
  ];

  return (
    <nav className="navigation">
      <div className="nav-brand">
        <h1>📊 FMP Analytics</h1>
      </div>
      <div className="nav-links">
        {navItems.map(({ path, label, icon: Icon }) => (
          <Link
            key={path}
            to={path}
            className={`nav-link ${location.pathname === path ? 'active' : ''}`}
          >
            <Icon size={20} />
            <span>{label}</span>
          </Link>
        ))}
      </div>
    </nav>
  );
};
