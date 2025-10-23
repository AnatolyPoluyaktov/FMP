import React from 'react';
import { TabType } from '../App';
import { Home, Tag, Settings, Calendar, DollarSign, Bell } from 'lucide-react';

interface NavigationProps {
  activeTab: TabType;
  onTabChange: (tab: TabType) => void;
}

const tabs = [
  { id: 'transactions' as TabType, label: 'Транзакции', icon: Home },
  { id: 'categories' as TabType, label: 'Категории', icon: Tag },
  { id: 'planned-expenses' as TabType, label: 'Планы расходов', icon: Calendar },
  { id: 'planned-income' as TabType, label: 'Планы доходов', icon: DollarSign },
  { id: 'limits' as TabType, label: 'Лимиты', icon: Settings },
  { id: 'notifications' as TabType, label: 'Уведомления', icon: Bell },
];

export const Navigation: React.FC<NavigationProps> = ({ activeTab, onTabChange }) => {
  return (
    <nav className="navigation">
      {tabs.map(({ id, label, icon: Icon }) => (
        <button
          key={id}
          className={`nav-button ${activeTab === id ? 'active' : ''}`}
          onClick={() => onTabChange(id)}
        >
          <Icon size={20} />
          <span>{label}</span>
        </button>
      ))}
    </nav>
  );
};
