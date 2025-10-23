import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { apiService, CategorySummary, Category } from '../services/api';

export const CategoryAnalytics: React.FC = () => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [categorySummary, setCategorySummary] = useState<CategorySummary[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [dateRange, setDateRange] = useState({
    start: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
    end: new Date().toISOString().split('T')[0]
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, [selectedCategory, dateRange]);

  const loadData = async () => {
    try {
      setLoading(true);
      
      const [categoriesData, summaryData] = await Promise.all([
        apiService.getCategories(),
        apiService.getCategorySummary({
          category_id: selectedCategory || undefined,
          start_date: dateRange.start,
          end_date: dateRange.end
        })
      ]);

      setCategories(categoriesData);
      setCategorySummary(summaryData);
    } catch (error) {
      console.error('Error loading category data:', error);
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

  const totalAmount = categorySummary.reduce((sum, item) => sum + item.amount, 0);

  if (loading) {
    return (
      <div className="category-analytics">
        <div className="loading">Загрузка данных...</div>
      </div>
    );
  }

  return (
    <div className="category-analytics">
      <div className="analytics-header">
        <h2>Аналитика по категориям</h2>
        
        <div className="filters">
          <div className="filter-group">
            <label htmlFor="category-select">Категория:</label>
            <select
              id="category-select"
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
            >
              <option value="">Все категории</option>
              {categories.map(category => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </select>
          </div>

          <div className="filter-group">
            <label htmlFor="start-date">От:</label>
            <input
              id="start-date"
              type="date"
              value={dateRange.start}
              onChange={(e) => setDateRange(prev => ({ ...prev, start: e.target.value }))}
            />
          </div>

          <div className="filter-group">
            <label htmlFor="end-date">До:</label>
            <input
              id="end-date"
              type="date"
              value={dateRange.end}
              onChange={(e) => setDateRange(prev => ({ ...prev, end: e.target.value }))}
            />
          </div>
        </div>
      </div>

      {/* Summary */}
      <div className="summary-section">
        <div className="summary-card">
          <h3>Общая сумма</h3>
          <p className="summary-value">{totalAmount.toLocaleString('ru-RU')} ₽</p>
        </div>
        <div className="summary-card">
          <h3>Количество категорий</h3>
          <p className="summary-value">{categorySummary.length}</p>
        </div>
        <div className="summary-card">
          <h3>Превышения лимитов</h3>
          <p className="summary-value">{categorySummary.filter(item => item.is_exceeded).length}</p>
        </div>
      </div>

      {/* Charts */}
      <div className="charts-grid">
        <div className="chart-container">
          <h3>Распределение по категориям</h3>
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

      {/* Category Details */}
      <div className="category-details">
        <h3>Детализация по категориям</h3>
        <div className="details-table">
          <div className="table-header">
            <span>Категория</span>
            <span>Сумма</span>
            <span>Лимит</span>
            <span>Статус</span>
          </div>
          {categorySummary.map(item => (
            <div key={item.category_id} className="table-row">
              <span className="category-name">{item.category_name}</span>
              <span className="amount">{item.amount.toLocaleString('ru-RU')} ₽</span>
              <span className="limit">
                {item.limit ? `${item.limit.toLocaleString('ru-RU')} ₽` : 'Не установлен'}
              </span>
              <span className={`status ${item.is_exceeded ? 'exceeded' : 'normal'}`}>
                {item.is_exceeded ? 'Превышен' : 'Норма'}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
