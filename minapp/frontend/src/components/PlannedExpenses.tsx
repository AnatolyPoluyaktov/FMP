import React, { useState, useEffect } from 'react';
import { format, addDays, addWeeks, addMonths } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Plus, Calendar, DollarSign, CheckCircle, Circle, Edit, Trash2 } from 'lucide-react';
import { apiService, Category, PlannedExpense } from '../services/api';

interface PlannedExpensesProps {
  categories: Category[];
}

export const PlannedExpenses: React.FC<PlannedExpensesProps> = ({ categories }) => {
  const [expenses, setExpenses] = useState<PlannedExpense[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    category_id: '',
    amount: '',
    description: '',
    planned_date: format(new Date(), 'yyyy-MM-dd'),
  });
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterPeriod, setFilterPeriod] = useState<'week' | 'month' | 'year'>('month');

  const filteredCategories = categories.filter(category =>
    category.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  useEffect(() => {
    loadExpenses();
  }, [filterPeriod]);

  const loadExpenses = async () => {
    try {
      const now = new Date();
      let startDate: string;
      let endDate: string;

      switch (filterPeriod) {
        case 'week':
          startDate = format(now, 'yyyy-MM-dd');
          endDate = format(addWeeks(now, 2), 'yyyy-MM-dd');
          break;
        case 'month':
          startDate = format(now, 'yyyy-MM-dd');
          endDate = format(addMonths(now, 1), 'yyyy-MM-dd');
          break;
        case 'year':
          startDate = format(now, 'yyyy-MM-dd');
          endDate = format(addMonths(now, 12), 'yyyy-MM-dd');
          break;
        default:
          startDate = format(now, 'yyyy-MM-dd');
          endDate = format(addMonths(now, 1), 'yyyy-MM-dd');
      }

      const data = await apiService.getPlannedExpenses({
        start_date: startDate,
        end_date: endDate,
      });
      setExpenses(data);
    } catch (error) {
      console.error('Error loading planned expenses:', error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.category_id || !formData.amount) return;

    try {
      setLoading(true);
      const expense = await apiService.createPlannedExpense({
        category_id: formData.category_id,
        amount: parseFloat(formData.amount),
        description: formData.description || undefined,
        planned_date: formData.planned_date,
      });
      
      setExpenses(prev => [...prev, expense]);
      
      // Reset form
      setFormData({
        category_id: '',
        amount: '',
        description: '',
        planned_date: format(new Date(), 'yyyy-MM-dd'),
      });
      setSearchTerm('');
      setShowForm(false);
    } catch (error) {
      console.error('Error creating planned expense:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleComplete = async (expense: PlannedExpense) => {
    try {
      await apiService.updatePlannedExpense(expense.id, {
        category_id: expense.category_id,
        amount: expense.amount,
        description: expense.description,
        planned_date: expense.planned_date,
      });
      
      setExpenses(prev => prev.map(e => 
        e.id === expense.id ? { ...e, is_completed: !e.is_completed } : e
      ));
    } catch (error) {
      console.error('Error updating planned expense:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiService.deletePlannedExpense(id);
      setExpenses(prev => prev.filter(e => e.id !== id));
    } catch (error) {
      console.error('Error deleting planned expense:', error);
    }
  };

  const getCategoryName = (categoryId: string) => {
    const category = categories.find(c => c.id === categoryId);
    return category?.name || 'Неизвестная категория';
  };

  const quickDateOptions = [
    { label: 'Сегодня', value: format(new Date(), 'yyyy-MM-dd') },
    { label: 'Завтра', value: format(addDays(new Date(), 1), 'yyyy-MM-dd') },
    { label: 'Через неделю', value: format(addWeeks(new Date(), 1), 'yyyy-MM-dd') },
    { label: 'Через месяц', value: format(addMonths(new Date(), 1), 'yyyy-MM-dd') },
  ];

  return (
    <div className="planned-expenses">
      <div className="planned-expenses-header">
        <h2>Планируемые расходы</h2>
        <div className="filter-buttons">
          <button
            className={filterPeriod === 'week' ? 'active' : ''}
            onClick={() => setFilterPeriod('week')}
          >
            2 недели
          </button>
          <button
            className={filterPeriod === 'month' ? 'active' : ''}
            onClick={() => setFilterPeriod('month')}
          >
            Месяц
          </button>
          <button
            className={filterPeriod === 'year' ? 'active' : ''}
            onClick={() => setFilterPeriod('year')}
          >
            Год
          </button>
        </div>
        <button
          className="add-button"
          onClick={() => setShowForm(!showForm)}
        >
          <Plus size={20} />
          Добавить расход
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="planned-expense-form">
          <div className="form-group">
            <label htmlFor="category-search">Категория</label>
            <div className="search-container">
              <input
                id="category-search"
                type="text"
                placeholder="Поиск категории..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="search-input"
              />
              {searchTerm && (
                <div className="category-dropdown">
                  {filteredCategories.map(category => (
                    <button
                      key={category.id}
                      type="button"
                      className="category-option"
                      onClick={() => {
                        setFormData(prev => ({ ...prev, category_id: category.id }));
                        setSearchTerm(category.name);
                      }}
                    >
                      {category.name}
                    </button>
                  ))}
                  {filteredCategories.length === 0 && (
                    <div className="no-results">Категории не найдены</div>
                  )}
                </div>
              )}
            </div>
          </div>

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
              placeholder="Описание планируемого расхода..."
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
            />
          </div>

          <div className="form-group">
            <label htmlFor="planned-date">Планируемая дата</label>
            <div className="date-input">
              <Calendar className="date-icon" />
              <input
                id="planned-date"
                type="date"
                value={formData.planned_date}
                onChange={(e) => setFormData(prev => ({ ...prev, planned_date: e.target.value }))}
                required
              />
            </div>
            <div className="quick-date-buttons">
              {quickDateOptions.map(option => (
                <button
                  key={option.value}
                  type="button"
                  className="quick-date-button"
                  onClick={() => setFormData(prev => ({ ...prev, planned_date: option.value }))}
                >
                  {option.label}
                </button>
              ))}
            </div>
          </div>

          <div className="form-actions">
            <button
              type="button"
              className="cancel-button"
              onClick={() => {
                setShowForm(false);
                setFormData({
                  category_id: '',
                  amount: '',
                  description: '',
                  planned_date: format(new Date(), 'yyyy-MM-dd'),
                });
                setSearchTerm('');
              }}
            >
              Отмена
            </button>
            <button
              type="submit"
              className="submit-button"
              disabled={loading || !formData.category_id || !formData.amount}
            >
              {loading ? 'Добавление...' : 'Добавить расход'}
            </button>
          </div>
        </form>
      )}

      <div className="expenses-list">
        {expenses.length === 0 ? (
          <div className="empty-state">
            <Calendar size={48} />
            <p>Планируемые расходы не найдены</p>
            <p>Добавьте первый планируемый расход</p>
          </div>
        ) : (
          expenses.map(expense => (
            <div key={expense.id} className={`expense-item ${expense.is_completed ? 'completed' : ''}`}>
              <div className="expense-info">
                <div className="expense-header">
                  <h3>{getCategoryName(expense.category_id)}</h3>
                  <span className="expense-amount">{expense.amount.toFixed(2)} ₽</span>
                </div>
                {expense.description && (
                  <p className="expense-description">{expense.description}</p>
                )}
                <p className="expense-date">
                  Планируется: {format(new Date(expense.planned_date), 'dd MMMM yyyy', { locale: ru })}
                </p>
              </div>
              <div className="expense-actions">
                <button
                  className="complete-button"
                  onClick={() => handleToggleComplete(expense)}
                  title={expense.is_completed ? 'Отметить как невыполненное' : 'Отметить как выполненное'}
                >
                  {expense.is_completed ? <CheckCircle size={20} /> : <Circle size={20} />}
                </button>
                <button className="edit-button" title="Редактировать">
                  <Edit size={16} />
                </button>
                <button
                  className="delete-button"
                  onClick={() => handleDelete(expense.id)}
                  title="Удалить"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
