import React, { useState } from 'react';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Plus, Calendar, DollarSign, Tag, FileText } from 'lucide-react';
import { apiService, Category, Transaction } from '../services/api';

interface TransactionFormProps {
  categories: Category[];
  onTransactionAdded: (transaction: Transaction) => void;
}

export const TransactionForm: React.FC<TransactionFormProps> = ({
  categories,
  onTransactionAdded,
}) => {
  const [formData, setFormData] = useState({
    category_id: '',
    amount: '',
    description: '',
    date: format(new Date(), 'yyyy-MM-dd'),
  });
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const filteredCategories = categories.filter(category =>
    category.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.category_id || !formData.amount) return;

    try {
      setLoading(true);
      const transaction = await apiService.createTransaction({
        category_id: formData.category_id,
        amount: parseFloat(formData.amount),
        description: formData.description || undefined,
        date: formData.date,
      });
      
      onTransactionAdded(transaction);
      
      // Reset form
      setFormData({
        category_id: '',
        amount: '',
        description: '',
        date: format(new Date(), 'yyyy-MM-dd'),
      });
      setSearchTerm('');
    } catch (error) {
      console.error('Error creating transaction:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="transaction-form">
      <div className="form-header">
        <h2>💰 Добавить транзакцию</h2>
        <p>Зафиксируйте ваши доходы и расходы</p>
      </div>
      
      <form onSubmit={handleSubmit} className="modern-form">
        <div className="form-group">
          <label htmlFor="category-search" className="form-label">
            <Tag className="label-icon" />
            Категория
          </label>
          <div className="search-container">
            <input
              id="category-search"
              type="text"
              placeholder="🔍 Поиск категории..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input modern-input"
            />
            {searchTerm && (
              <div className="category-dropdown modern-dropdown">
                {filteredCategories.map(category => (
                  <button
                    key={category.id}
                    type="button"
                    className="category-option modern-option"
                    onClick={() => {
                      setFormData(prev => ({ ...prev, category_id: category.id }));
                      setSearchTerm(category.name);
                    }}
                  >
                    <span className="category-name">{category.name}</span>
                  </button>
                ))}
                {filteredCategories.length === 0 && (
                  <div className="no-results">❌ Категории не найдены</div>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="amount" className="form-label">
            <DollarSign className="label-icon" />
            Сумма
          </label>
          <div className="amount-input modern-input-group">
            <DollarSign className="input-icon" />
            <input
              id="amount"
              type="number"
              step="0.01"
              placeholder="0.00"
              value={formData.amount}
              onChange={(e) => setFormData(prev => ({ ...prev, amount: e.target.value }))}
              className="modern-input amount-field"
              required
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="description" className="form-label">
            <FileText className="label-icon" />
            Описание (необязательно)
          </label>
          <input
            id="description"
            type="text"
            placeholder="📝 Описание транзакции..."
            value={formData.description}
            onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
            className="modern-input"
          />
        </div>

        <div className="form-group">
          <label htmlFor="date" className="form-label">
            <Calendar className="label-icon" />
            Дата
          </label>
          <div className="date-input modern-input-group">
            <Calendar className="input-icon" />
            <input
              id="date"
              type="date"
              value={formData.date}
              onChange={(e) => setFormData(prev => ({ ...prev, date: e.target.value }))}
              className="modern-input date-field"
              required
            />
          </div>
        </div>

        <button
          type="submit"
          className="submit-button modern-button primary"
          disabled={loading || !formData.category_id || !formData.amount}
        >
          <Plus size={20} />
          {loading ? '⏳ Добавление...' : '✅ Добавить транзакцию'}
        </button>
      </form>
    </div>
  );
};
