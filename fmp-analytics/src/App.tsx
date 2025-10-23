import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Navigation } from './components/Navigation';
import { Dashboard } from './components/Dashboard';
import { MonthlyAnalytics } from './components/MonthlyAnalytics';
import { CategoryAnalytics } from './components/CategoryAnalytics';
import { LimitAnalytics } from './components/LimitAnalytics';
import './App.css';

const App: React.FC = () => {
  return (
    <Router>
      <div className="app">
        <Navigation />
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/monthly" element={<MonthlyAnalytics />} />
            <Route path="/categories" element={<CategoryAnalytics />} />
            <Route path="/limits" element={<LimitAnalytics />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
};

export default App;
