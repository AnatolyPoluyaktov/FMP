import React, { useState } from 'react';
import { Plus, Tag, Edit, Trash2, FileText } from 'lucide-react';
import { apiService, Category } from '../services/api';

interface CategoryManagerProps {
  categories: Category[];
  onCategoryAdded: (category: Category) => void;
}

export const CategoryManager: React.FC<CategoryManagerProps> = ({
  categories,
  onCategoryAdded,
}) => {
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
  });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) return;

    try {
      setLoading(true);
      const category = await apiService.createCategory({
        name: formData.name.trim(),
        description: formData.description.trim() || undefined,
      });
      
      onCategoryAdded(category);
      
      // Reset form
      setFormData({ name: '', description: '' });
      setShowForm(false);
    } catch (error) {
      console.error('Error creating category:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="category-manager">
      <div className="category-header">
        <div className="header-content">
          <h2>🏷️ Управление категориями</h2>
          <p>Создавайте и организуйте категории для ваших транзакций</p>
        </div>
        <button
          className="add-button modern-button primary"
          onClick={() => setShowForm(!showForm)}
        >
          <Plus size={20} />
          {showForm ? 'Скрыть форму' : 'Добавить категорию'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="category-form modern-form">
          <div className="form-group">
            <label htmlFor="category-name" className="form-label">
              <Tag className="label-icon" />
              Название категории
            </label>
            <input
              id="category-name"
              type="text"
              placeholder="🏪 Например: Еда, Транспорт, Развлечения..."
              value={formData.name}
              onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
              className="modern-input"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="category-description" className="form-label">
              <FileText className="label-icon" />
              Описание (необязательно)
            </label>
            <input
              id="category-description"
              type="text"
              placeholder="📝 Описание категории..."
              value={formData.description}
              onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              className="modern-input"
            />
          </div>

          <div className="form-actions">
            <button
              type="button"
              className="cancel-button modern-button secondary"
              onClick={() => {
                setShowForm(false);
                setFormData({ name: '', description: '' });
              }}
            >
              ❌ Отмена
            </button>
            <button
              type="submit"
              className="submit-button modern-button primary"
              disabled={loading || !formData.name.trim()}
            >
              {loading ? '⏳ Создание...' : '✅ Создать категорию'}
            </button>
          </div>
        </form>
      )}

      <div className="categories-list">
        {categories.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">
              <Tag size={48} />
            </div>
            <h3>📂 Категории не найдены</h3>
            <p>Создайте первую категорию для начала работы с транзакциями</p>
          </div>
        ) : (
          <div className="categories-grid">
            {categories.map(category => (
              <div key={category.id} className="category-item modern-card">
                <div className="category-info">
                  <div className="category-header-item">
                    <Tag className="category-icon" />
                    <h3>{category.name}</h3>
                  </div>
                  {category.description && (
                    <p className="category-description">{category.description}</p>
                  )}
                  <p className="category-date">
                    📅 Создано: {new Date(category.created_at).toLocaleDateString('ru-RU')}
                  </p>
                </div>
                <div className="category-actions">
                  <button className="edit-button modern-button-icon" title="Редактировать">
                    <Edit size={16} />
                  </button>
                  <button className="delete-button modern-button-icon danger" title="Удалить">
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
