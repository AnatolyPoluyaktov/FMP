import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, AlertTriangle } from 'lucide-react';
import { apiService, CategoryLimit, Category } from '../services/api';

interface CategoryLimitsProps {
  categories: Category[];
}

export const CategoryLimits: React.FC<CategoryLimitsProps> = ({ categories }) => {
  const [limits, setLimits] = useState<CategoryLimit[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editingLimit, setEditingLimit] = useState<CategoryLimit | null>(null);
  const [formData, setFormData] = useState({
    category_id: '',
    limit: '',
    month: new Date().getMonth() + 1,
    year: new Date().getFullYear(),
  });

  useEffect(() => {
    loadLimits();
  }, []);

  const loadLimits = async () => {
    try {
      setLoading(true);
      const data = await apiService.getCategoryLimits();
      setLimits(data);
    } catch (error) {
      console.error('Error loading limits:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.category_id || !formData.limit) return;

    try {
      if (editingLimit) {
        // Update existing limit
        await apiService.updateCategoryLimit(editingLimit.id, {
          category_id: formData.category_id,
          limit: parseFloat(formData.limit),
          month: formData.month,
          year: formData.year,
        });
      } else {
        // Create new limit
        await apiService.createCategoryLimit({
          category_id: formData.category_id,
          limit: parseFloat(formData.limit),
          month: formData.month,
          year: formData.year,
        });
      }

      await loadLimits();
      resetForm();
    } catch (error) {
      console.error('Error saving limit:', error);
    }
  };

  const handleEdit = (limit: CategoryLimit) => {
    setEditingLimit(limit);
    setFormData({
      category_id: limit.category_id,
      limit: limit.limit.toString(),
      month: limit.month,
      year: limit.year,
    });
    setShowForm(true);
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот лимит?')) return;

    try {
      await apiService.deleteCategoryLimit(id);
      await loadLimits();
    } catch (error) {
      console.error('Error deleting limit:', error);
    }
  };

  const resetForm = () => {
    setFormData({
      category_id: '',
      limit: '',
      month: new Date().getMonth() + 1,
      year: new Date().getFullYear(),
    });
    setEditingLimit(null);
    setShowForm(false);
  };

  const getCategoryName = (categoryId: string) => {
    const category = categories.find(c => c.id === categoryId);
    return category ? category.name : 'Неизвестная категория';
  };

  const getMonthName = (month: number) => {
    const months = [
      'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
      'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'
    ];
    return months[month - 1];
  };

  if (loading) {
    return (
      <div className="category-limits">
        <div className="loading">Загрузка лимитов...</div>
      </div>
    );
  }

  return (
    <div className="category-limits">
      <div className="limits-header">
        <h2>Лимиты по категориям</h2>
        <button
          className="add-button"
          onClick={() => setShowForm(true)}
        >
          <Plus size={20} />
          Добавить лимит
        </button>
      </div>

      {showForm && (
        <div className="limit-form">
          <h3>{editingLimit ? 'Редактировать лимит' : 'Добавить лимит'}</h3>
          
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="category-select">Категория</label>
              <select
                id="category-select"
                value={formData.category_id}
                onChange={(e) => setFormData(prev => ({ ...prev, category_id: e.target.value }))}
                required
              >
                <option value="">Выберите категорию</option>
                {categories.map(category => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label htmlFor="limit">Лимит (₽)</label>
              <input
                id="limit"
                type="number"
                step="0.01"
                min="0"
                placeholder="0.00"
                value={formData.limit}
                onChange={(e) => setFormData(prev => ({ ...prev, limit: e.target.value }))}
                required
              />
            </div>

            <div className="form-row">
              <div className="form-group">
                <label htmlFor="month">Месяц</label>
                <select
                  id="month"
                  value={formData.month}
                  onChange={(e) => setFormData(prev => ({ ...prev, month: parseInt(e.target.value) }))}
                  required
                >
                  {Array.from({ length: 12 }, (_, i) => (
                    <option key={i + 1} value={i + 1}>
                      {getMonthName(i + 1)}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label htmlFor="year">Год</label>
                <input
                  id="year"
                  type="number"
                  min="2020"
                  max="2030"
                  value={formData.year}
                  onChange={(e) => setFormData(prev => ({ ...prev, year: parseInt(e.target.value) }))}
                  required
                />
              </div>
            </div>

            <div className="form-actions">
              <button type="submit" className="save-button">
                {editingLimit ? 'Сохранить' : 'Добавить'}
              </button>
              <button type="button" className="cancel-button" onClick={resetForm}>
                Отмена
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="limits-list">
        {limits.length === 0 ? (
          <div className="no-limits">
            <AlertTriangle className="no-limits-icon" />
            <p>Лимиты не установлены</p>
            <p>Добавьте лимиты для контроля расходов</p>
          </div>
        ) : (
          limits.map(limit => (
            <div key={limit.id} className="limit-item">
              <div className="limit-info">
                <h3>{getCategoryName(limit.category_id)}</h3>
                <p className="limit-period">
                  {getMonthName(limit.month)} {limit.year}
                </p>
                <p className="limit-amount">{limit.limit.toLocaleString('ru-RU')} ₽</p>
              </div>
              
              <div className="limit-actions">
                <button
                  className="edit-button"
                  onClick={() => handleEdit(limit)}
                  title="Редактировать"
                >
                  <Edit size={16} />
                </button>
                <button
                  className="delete-button"
                  onClick={() => handleDelete(limit.id)}
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
