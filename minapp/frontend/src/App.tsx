import React, { useState, useEffect } from 'react';
import { WebApp } from './telegram-webapp';
import { TransactionForm } from './components/TransactionForm';
import { CategoryManager } from './components/CategoryManager';
import { CategoryLimits } from './components/CategoryLimits';
import { PlannedExpenses } from './components/PlannedExpenses';
import { PlannedIncome } from './components/PlannedIncome';
import { Notifications } from './components/Notifications';
import { Navigation } from './components/Navigation';
import { apiService } from './services/api';
import './App.css';

export type TabType = 'transactions' | 'categories' | 'limits' | 'planned-expenses' | 'planned-income' | 'notifications';

function App() {
  const [activeTab, setActiveTab] = useState<TabType>('transactions');
  const [categories, setCategories] = useState<any[]>([]);
  const [transactions, setTransactions] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Initialize Telegram WebApp
    WebApp.ready();
    WebApp.expand();
    
    // Set up Telegram WebApp theme
    WebApp.setHeaderColor('#2481cc');
    WebApp.setBackgroundColor('#ffffff');

    // Load initial data
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [categoriesData, transactionsData] = await Promise.all([
        apiService.getCategories(),
        apiService.getTransactions()
      ]);
      setCategories(categoriesData);
      setTransactions(transactionsData);
    } catch (error) {
      console.error('Error loading data:', error);
      WebApp.showAlert('Ошибка загрузки данных');
    } finally {
      setLoading(false);
    }
  };

  const handleTransactionAdded = (transaction: any) => {
    setTransactions(prev => [transaction, ...prev]);
    
    // Show success feedback
    WebApp.showAlert('Транзакция добавлена!');
    
    // Check for limit warnings after adding transaction
    apiService.checkLimitWarnings().catch(console.error);
  };

  const handleCategoryAdded = (category: any) => {
    setCategories(prev => [...prev, category]);
    WebApp.showAlert('Категория добавлена!');
  };

  const renderActiveTab = () => {
    switch (activeTab) {
      case 'transactions':
        return (
          <TransactionForm
            categories={categories}
            onTransactionAdded={handleTransactionAdded}
          />
        );
      case 'categories':
        return (
          <CategoryManager
            categories={categories}
            onCategoryAdded={handleCategoryAdded}
          />
        );
      case 'limits':
        return <CategoryLimits categories={categories} />;
      case 'planned-expenses':
        return <PlannedExpenses categories={categories} />;
      case 'planned-income':
        return <PlannedIncome />;
      case 'notifications':
        return <Notifications />;
      default:
        return null;
    }
  };

  if (loading) {
    return (
      <div className="app">
        <div className="loading-container">
          <div className="loading">Загрузка...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <div className="app-header">
        <h1>💰 FMP</h1>
        <p>Финансовый менеджер</p>
      </div>
      
      <Navigation activeTab={activeTab} onTabChange={setActiveTab} />
      
      <div className="app-content">
        {renderActiveTab()}
      </div>
    </div>
  );
}

export default App;