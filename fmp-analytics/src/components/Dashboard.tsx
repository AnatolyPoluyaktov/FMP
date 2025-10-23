import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { DollarSign, TrendingUp, AlertTriangle, Calendar } from 'lucide-react';
import { apiService, MonthlySummary, CategorySummary, LimitExceeded } from '../services/api';

export const Dashboard: React.FC = () => {
  const [currentMonth, setCurrentMonth] = useState(new Date());
  const [monthlySummary, setMonthlySummary] = useState<MonthlySummary | null>(null);
  const [categorySummary, setCategorySummary] = useState<CategorySummary[]>([]);
  const [limitExceeded, setLimitExceeded] = useState<LimitExceeded[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, [currentMonth]);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      const month = currentMonth.getMonth() + 1;
      const year = currentMonth.getFullYear();
      
      const [monthlyData, categoryData, limitData] = await Promise.all([
        apiService.getMonthlySummary(month, year),
        apiService.getCategorySummary({
          start_date: `${year}-${month.toString().padStart(2, '0')}-01`,
          end_date: `${year}-${month.toString().padStart(2, '0')}-31`
        }),
        apiService.getLimitExceeded()
      ]);

      setMonthlySummary(monthlyData);
      setCategorySummary(categoryData);
      setLimitExceeded(limitData);
    } catch (error) {
      console.error('Error loading dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82CA9D'];

  const pieData = categorySummary.map(item => ({
    name: item.category_name,
    value: item.amount,
    color: COLORS[categorySummary.indexOf(item) % COLORS.length]
  }));

  const barData = categorySummary.map(item => ({
    name: item.category_name,
    amount: item.amount,
    limit: item.limit || 0,
    exceeded: item.is_exceeded
  }));

  if (loading) {
    return (
      <div className="dashboard">
        <div className="loading">Загрузка данных...</div>
      </div>
    );
  }

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <h2>Обзор финансов</h2>
        <div className="month-selector">
          <button 
            onClick={() => setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1))}
            className="month-button"
          >
            ←
          </button>
          <span className="current-month">
            {currentMonth.toLocaleDateString('ru-RU', { month: 'long', year: 'numeric' })}
          </span>
          <button 
            onClick={() => setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1))}
            className="month-button"
          >
            →
          </button>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="summary-cards">
        <div className="summary-card">
          <div className="card-icon">
            <DollarSign size={24} />
          </div>
          <div className="card-content">
            <h3>Общие расходы</h3>
            <p className="card-value">{monthlySummary?.total.toLocaleString('ru-RU')} ₽</p>
          </div>
        </div>

        <div className="summary-card">
          <div className="card-icon">
            <TrendingUp size={24} />
          </div>
          <div className="card-content">
            <h3>Категорий</h3>
            <p className="card-value">{categorySummary.length}</p>
          </div>
        </div>

        <div className="summary-card">
          <div className="card-icon">
            <AlertTriangle size={24} />
          </div>
          <div className="card-content">
            <h3>Превышения лимитов</h3>
            <p className="card-value">{limitExceeded.length}</p>
          </div>
        </div>
      </div>

      {/* Charts */}
      <div className="charts-grid">
        <div className="chart-container">
          <h3>Расходы по категориям</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {pieData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => [`${value.toLocaleString('ru-RU')} ₽`, 'Сумма']} />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="chart-container">
          <h3>Сравнение с лимитами</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={barData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip formatter={(value) => [`${value.toLocaleString('ru-RU')} ₽`, 'Сумма']} />
              <Bar dataKey="amount" fill="#8884d8" name="Потрачено" />
              <Bar dataKey="limit" fill="#82ca9d" name="Лимит" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Recent Limit Exceeded */}
      {limitExceeded.length > 0 && (
        <div className="limit-exceeded-section">
          <h3>Недавние превышения лимитов</h3>
          <div className="limit-exceeded-list">
            {limitExceeded.slice(0, 5).map(item => (
              <div key={item.id} className="limit-exceeded-item">
                <div className="exceeded-info">
                  <span className="exceeded-category">Категория ID: {item.category_id}</span>
                  <span className="exceeded-date">
                    {new Date(item.created_at).toLocaleDateString('ru-RU')}
                  </span>
                </div>
                <div className="exceeded-amounts">
                  <span className="limit">Лимит: {item.limit.toLocaleString('ru-RU')} ₽</span>
                  <span className="actual">Потрачено: {item.actual.toLocaleString('ru-RU')} ₽</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
