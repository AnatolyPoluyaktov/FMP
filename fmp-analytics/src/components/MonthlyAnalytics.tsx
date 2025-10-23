import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { apiService, MonthlySummary } from '../services/api';

export const MonthlyAnalytics: React.FC = () => {
  const [year, setYear] = useState(new Date().getFullYear());
  const [monthlyData, setMonthlyData] = useState<MonthlySummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadMonthlyData();
  }, [year]);

  const loadMonthlyData = async () => {
    try {
      setLoading(true);
      const data = [];
      
      for (let month = 1; month <= 12; month++) {
        try {
          const summary = await apiService.getMonthlySummary(month, year);
          data.push(summary);
        } catch (error) {
          // If no data for this month, add empty data
          data.push({
            month,
            year,
            categories: [],
            total: 0
          });
        }
      }
      
      setMonthlyData(data);
    } catch (error) {
      console.error('Error loading monthly data:', error);
    } finally {
      setLoading(false);
    }
  };

  const chartData = monthlyData.map(item => ({
    month: new Date(year, item.month - 1).toLocaleDateString('ru-RU', { month: 'short' }),
    amount: item.total,
    fullMonth: new Date(year, item.month - 1).toLocaleDateString('ru-RU', { month: 'long' })
  }));

  const totalYear = monthlyData.reduce((sum, item) => sum + item.total, 0);
  const averageMonth = totalYear / 12;
  const maxMonth = Math.max(...monthlyData.map(item => item.total));
  const minMonth = Math.min(...monthlyData.map(item => item.total));

  if (loading) {
    return (
      <div className="monthly-analytics">
        <div className="loading">Загрузка данных...</div>
      </div>
    );
  }

  return (
    <div className="monthly-analytics">
      <div className="analytics-header">
        <h2>Месячная аналитика</h2>
        <div className="year-selector">
          <button 
            onClick={() => setYear(year - 1)}
            className="year-button"
          >
            ←
          </button>
          <span className="current-year">{year}</span>
          <button 
            onClick={() => setYear(year + 1)}
            className="year-button"
          >
            →
          </button>
        </div>
      </div>

      {/* Summary Stats */}
      <div className="summary-stats">
        <div className="stat-card">
          <h3>Общие расходы за год</h3>
          <p className="stat-value">{totalYear.toLocaleString('ru-RU')} ₽</p>
        </div>
        <div className="stat-card">
          <h3>Средние расходы в месяц</h3>
          <p className="stat-value">{averageMonth.toLocaleString('ru-RU')} ₽</p>
        </div>
        <div className="stat-card">
          <h3>Максимальный месяц</h3>
          <p className="stat-value">{maxMonth.toLocaleString('ru-RU')} ₽</p>
        </div>
        <div className="stat-card">
          <h3>Минимальный месяц</h3>
          <p className="stat-value">{minMonth.toLocaleString('ru-RU')} ₽</p>
        </div>
      </div>

      {/* Monthly Chart */}
      <div className="chart-container">
        <h3>Динамика расходов по месяцам</h3>
        <ResponsiveContainer width="100%" height={400}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="month" />
            <YAxis />
            <Tooltip 
              formatter={(value) => [`${value.toLocaleString('ru-RU')} ₽`, 'Сумма']}
              labelFormatter={(label, payload) => {
                if (payload && payload[0]) {
                  return payload[0].payload.fullMonth;
                }
                return label;
              }}
            />
            <Line 
              type="monotone" 
              dataKey="amount" 
              stroke="#8884d8" 
              strokeWidth={2}
              dot={{ fill: '#8884d8', strokeWidth: 2, r: 4 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Monthly Breakdown */}
      <div className="monthly-breakdown">
        <h3>Детализация по месяцам</h3>
        <div className="breakdown-list">
          {monthlyData.map(item => (
            <div key={`${item.year}-${item.month}`} className="breakdown-item">
              <div className="month-info">
                <span className="month-name">
                  {new Date(item.year, item.month - 1).toLocaleDateString('ru-RU', { month: 'long' })}
                </span>
                <span className="month-total">{item.total.toLocaleString('ru-RU')} ₽</span>
              </div>
              {item.categories.length > 0 && (
                <div className="categories-list">
                  {item.categories.slice(0, 3).map(category => (
                    <div key={category.category_id} className="category-item">
                      <span className="category-name">{category.category_name}</span>
                      <span className="category-amount">{category.amount.toLocaleString('ru-RU')} ₽</span>
                    </div>
                  ))}
                  {item.categories.length > 3 && (
                    <div className="more-categories">
                      и еще {item.categories.length - 3} категорий
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
