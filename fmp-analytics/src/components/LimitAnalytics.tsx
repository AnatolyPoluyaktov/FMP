import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { AlertTriangle, TrendingUp, DollarSign } from 'lucide-react';
import { apiService, LimitExceeded, CategoryLimit, Category } from '../services/api';

export const LimitAnalytics: React.FC = () => {
  const [limitExceeded, setLimitExceeded] = useState<LimitExceeded[]>([]);
  const [categoryLimits, setCategoryLimits] = useState<CategoryLimit[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedYear, setSelectedYear] = useState(new Date().getFullYear());
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, [selectedYear]);

  const loadData = async () => {
    try {
      setLoading(true);
      
      const [exceededData, limitsData, categoriesData] = await Promise.all([
        apiService.getLimitExceeded(),
        apiService.getCategoryLimits({ year: selectedYear }),
        apiService.getCategories()
      ]);

      setLimitExceeded(exceededData);
      setCategoryLimits(limitsData);
      setCategories(categoriesData);
    } catch (error) {
      console.error('Error loading limit data:', error);
    } finally {
      setLoading(false);
    }
  };

  const getCategoryName = (categoryId: string) => {
    const category = categories.find(c => c.id === categoryId);
    return category ? category.name : `Категория ${categoryId}`;
  };

  const exceededByMonth = limitExceeded.reduce((acc, item) => {
    const month = item.month;
    if (!acc[month]) {
      acc[month] = [];
    }
    acc[month].push(item);
    return acc;
  }, {} as Record<number, LimitExceeded[]>);

  const chartData = Object.entries(exceededByMonth).map(([month, items]) => ({
    month: new Date(selectedYear, parseInt(month) - 1).toLocaleDateString('ru-RU', { month: 'short' }),
    count: items.length,
    totalExceeded: items.reduce((sum, item) => sum + (item.actual - item.limit), 0)
  }));

  const totalExceeded = limitExceeded.length;
  const totalOverLimit = limitExceeded.reduce((sum, item) => sum + (item.actual - item.limit), 0);
  const averageExceeded = totalExceeded > 0 ? totalOverLimit / totalExceeded : 0;

  if (loading) {
    return (
      <div className="limit-analytics">
        <div className="loading">Загрузка данных...</div>
      </div>
    );
  }

  return (
    <div className="limit-analytics">
      <div className="analytics-header">
        <h2>Аналитика лимитов</h2>
        
        <div className="year-selector">
          <button 
            onClick={() => setSelectedYear(selectedYear - 1)}
            className="year-button"
          >
            ←
          </button>
          <span className="current-year">{selectedYear}</span>
          <button 
            onClick={() => setSelectedYear(selectedYear + 1)}
            className="year-button"
          >
            →
          </button>
        </div>
      </div>

      {/* Summary Stats */}
      <div className="summary-stats">
        <div className="stat-card">
          <div className="stat-icon">
            <AlertTriangle size={24} />
          </div>
          <div className="stat-content">
            <h3>Всего превышений</h3>
            <p className="stat-value">{totalExceeded}</p>
          </div>
        </div>
        
        <div className="stat-card">
          <div className="stat-icon">
            <DollarSign size={24} />
          </div>
          <div className="stat-content">
            <h3>Сумма превышений</h3>
            <p className="stat-value">{totalOverLimit.toLocaleString('ru-RU')} ₽</p>
          </div>
        </div>
        
        <div className="stat-card">
          <div className="stat-icon">
            <TrendingUp size={24} />
          </div>
          <div className="stat-content">
            <h3>Среднее превышение</h3>
            <p className="stat-value">{averageExceeded.toLocaleString('ru-RU')} ₽</p>
          </div>
        </div>
      </div>

      {/* Charts */}
      <div className="charts-grid">
        <div className="chart-container">
          <h3>Превышения по месяцам</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip 
                formatter={(value, name) => [
                  name === 'count' ? `${value} превышений` : `${value.toLocaleString('ru-RU')} ₽`,
                  name === 'count' ? 'Количество' : 'Сумма превышений'
                ]}
              />
              <Bar dataKey="count" fill="#ff4444" name="Количество превышений" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="chart-container">
          <h3>Сумма превышений по месяцам</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip formatter={(value) => [`${value.toLocaleString('ru-RU')} ₽`, 'Сумма превышений']} />
              <Bar dataKey="totalExceeded" fill="#ff8042" name="Сумма превышений" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Current Limits */}
      <div className="current-limits">
        <h3>Текущие лимиты на {selectedYear} год</h3>
        <div className="limits-table">
          <div className="table-header">
            <span>Категория</span>
            <span>Месяц</span>
            <span>Лимит</span>
          </div>
          {categoryLimits.map(limit => (
            <div key={limit.id} className="table-row">
              <span className="category-name">{getCategoryName(limit.category_id)}</span>
              <span className="month">
                {new Date(selectedYear, limit.month - 1).toLocaleDateString('ru-RU', { month: 'long' })}
              </span>
              <span className="limit">{limit.limit.toLocaleString('ru-RU')} ₽</span>
            </div>
          ))}
        </div>
      </div>

      {/* Recent Exceeded */}
      <div className="recent-exceeded">
        <h3>Недавние превышения</h3>
        <div className="exceeded-list">
          {limitExceeded.slice(0, 10).map(item => (
            <div key={item.id} className="exceeded-item">
              <div className="exceeded-header">
                <span className="category-name">{getCategoryName(item.category_id)}</span>
                <span className="exceeded-date">
                  {new Date(item.created_at).toLocaleDateString('ru-RU')}
                </span>
              </div>
              <div className="exceeded-details">
                <div className="amount-detail">
                  <span className="label">Лимит:</span>
                  <span className="limit">{item.limit.toLocaleString('ru-RU')} ₽</span>
                </div>
                <div className="amount-detail">
                  <span className="label">Потрачено:</span>
                  <span className="actual">{item.actual.toLocaleString('ru-RU')} ₽</span>
                </div>
                <div className="amount-detail">
                  <span className="label">Превышение:</span>
                  <span className="exceeded">
                    +{(item.actual - item.limit).toLocaleString('ru-RU')} ₽
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
