import React, { useState, useEffect } from 'react';
import { Plus, DollarSign, Edit, Trash2 } from 'lucide-react';
import { apiService, PlannedIncome as PlannedIncomeType } from '../services/api';

export const PlannedIncome: React.FC = () => {
  const [incomes, setIncomes] = useState<PlannedIncomeType[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    amount: '',
    description: '',
    month: new Date().getMonth() + 1,
    year: new Date().getFullYear(),
  });
  const [loading, setLoading] = useState(false);
  const [selectedYear, setSelectedYear] = useState(new Date().getFullYear());

  const months = [
    'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
    'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'
  ];

  const years = Array.from({ length: 5 }, (_, i) => new Date().getFullYear() - 2 + i);

  useEffect(() => {
    loadIncomes();
  }, [selectedYear]);

  const loadIncomes = async () => {
    try {
      const data = await apiService.getPlannedIncome({
        year: selectedYear,
      });
      setIncomes(data);
    } catch (error) {
      console.error('Error loading planned income:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.amount) return;

    try {
      setLoading(true);
      const income = await apiService.createPlannedIncome({
        amount: parseFloat(formData.amount),
        description: formData.description || undefined,
        month: formData.month,
        year: formData.year,
      });
      
      setIncomes(prev => [...prev, income]);
      
      // Reset form
      setFormData({
        amount: '',
        description: '',
        month: new Date().getMonth() + 1,
        year: new Date().getFullYear(),
      });
      setShowForm(false);
    } catch (error) {
      console.error('Error creating planned income:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiService.deletePlannedIncome(id);
      setIncomes(prev => prev.filter(i => i.id !== id));
    } catch (error) {
      console.error('Error deleting planned income:', error);
    }
  };

  const getIncomeForMonth = (month: number) => {
    return incomes.find(income => income.month === month);
  };

  const getTotalIncome = () => {
    return incomes.reduce((total, income) => total + income.amount, 0);
  };

  const copyFromPreviousMonth = async (month: number) => {
    const previousMonth = month === 1 ? 12 : month - 1;
    const previousYear = month === 1 ? selectedYear - 1 : selectedYear;
    
    const previousIncome = await apiService.getPlannedIncome({
      month: previousMonth,
      year: previousYear,
    });

    if (previousIncome.length > 0) {
      const income = previousIncome[0];
      try {
        await apiService.createPlannedIncome({
          amount: income.amount,
          description: income.description || `Скопировано из ${months[previousMonth - 1]} ${previousYear}`,
          month: month,
          year: selectedYear,
        });
        loadIncomes();
      } catch (error) {
        console.error('Error copying income:', error);
      }
    }
  };

  return (
    <div className="planned-income">
      <div className="planned-income-header">
        <h2>Планируемые доходы</h2>
        <div className="year-selector">
          <label htmlFor="year-select">Год:</label>
          <select
            id="year-select"
            value={selectedYear}
            onChange={(e) => setSelectedYear(parseInt(e.target.value))}
          >
            {years.map(year => (
              <option key={year} value={year}>{year}</option>
            ))}
          </select>
        </div>
        <button
          className="add-button"
          onClick={() => setShowForm(!showForm)}
        >
          <Plus size={20} />
          Добавить доход
        </button>
      </div>

      <div className="income-summary">
        <div className="summary-card">
          <h3>Общий планируемый доход за {selectedYear} год</h3>
          <div className="summary-amount">{getTotalIncome().toFixed(2)} ₽</div>
        </div>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="planned-income-form">
          <div className="form-group">
            <label htmlFor="amount">Сумма</label>
            <div className="amount-input">
              <DollarSign className="amount-icon" />
              <input
                id="amount"
                type="number"
                step="0.01"
                placeholder="0.00"
                value={formData.amount}
                onChange={(e) => setFormData(prev => ({ ...prev, amount: e.target.value }))}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="description">Описание (необязательно)</label>
            <input
              id="description"
              type="text"
              placeholder="Описание планируемого дохода..."
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
            />
          </div>

          <div className="form-group">
            <label htmlFor="month">Месяц</label>
            <select
              id="month"
              value={formData.month}
              onChange={(e) => setFormData(prev => ({ ...prev, month: parseInt(e.target.value) }))}
              required
            >
              {months.map((month, index) => (
                <option key={index + 1} value={index + 1}>{month}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="year">Год</label>
            <select
              id="year"
              value={formData.year}
              onChange={(e) => setFormData(prev => ({ ...prev, year: parseInt(e.target.value) }))}
              required
            >
              {years.map(year => (
                <option key={year} value={year}>{year}</option>
              ))}
            </select>
          </div>

          <div className="form-actions">
            <button
              type="button"
              className="cancel-button"
              onClick={() => {
                setShowForm(false);
                setFormData({
                  amount: '',
                  description: '',
                  month: new Date().getMonth() + 1,
                  year: new Date().getFullYear(),
                });
              }}
            >
              Отмена
            </button>
            <button
              type="submit"
              className="submit-button"
              disabled={loading || !formData.amount}
            >
              {loading ? 'Добавление...' : 'Добавить доход'}
            </button>
          </div>
        </form>
      )}

      <div className="income-calendar">
        <h3>Доходы по месяцам</h3>
        <div className="months-grid">
          {months.map((month, index) => {
            const monthNumber = index + 1;
            const income = getIncomeForMonth(monthNumber);
            
            return (
              <div key={monthNumber} className="month-card">
                <div className="month-header">
                  <h4>{month}</h4>
                  {!income && (
                    <button
                      className="copy-button"
                      onClick={() => copyFromPreviousMonth(monthNumber)}
                      title="Скопировать из предыдущего месяца"
                    >
                      📋
                    </button>
                  )}
                </div>
                <div className="month-content">
                  {income ? (
                    <div className="income-details">
                      <div className="income-amount">{income.amount.toFixed(2)} ₽</div>
                      {income.description && (
                        <div className="income-description">{income.description}</div>
                      )}
                      <div className="income-actions">
                        <button className="edit-button" title="Редактировать">
                          <Edit size={14} />
                        </button>
                        <button
                          className="delete-button"
                          onClick={() => handleDelete(income.id)}
                          title="Удалить"
                        >
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </div>
                  ) : (
                    <div className="no-income">
                      <span>Нет данных</span>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
