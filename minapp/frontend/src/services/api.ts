import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Category {
  id: string;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: string;
  category_id: string;
  amount: number;
  description?: string;
  date: string;
  created_at: string;
  updated_at: string;
}

export interface CategoryLimit {
  id: string;
  category_id: string;
  limit: number;
  month: number;
  year: number;
  created_at: string;
  updated_at: string;
}

export interface MonthlySummary {
  month: number;
  year: number;
  categories: CategorySummary[];
  total: number;
}

export interface CategorySummary {
  category_id: string;
  category_name: string;
  amount: number;
  limit?: number;
  is_exceeded: boolean;
}

export interface PlannedExpense {
  id: string;
  category_id: string;
  amount: number;
  description?: string;
  planned_date: string;
  is_completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface PlannedIncome {
  id: string;
  amount: number;
  description?: string;
  month: number;
  year: number;
  created_at: string;
  updated_at: string;
}

export interface Notification {
  id: string;
  type: 'daily_reminder' | 'limit_warning' | 'limit_exceeded' | 'income_reminder';
  title: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationStats {
  unread_count: number;
  total_count: number;
}

export const apiService = {
  // Categories
  getCategories: async (): Promise<Category[]> => {
    const response = await api.get('/categories');
    return response.data;
  },

  createCategory: async (data: { name: string; description?: string }): Promise<Category> => {
    const response = await api.post('/categories', data);
    return response.data;
  },

  // Transactions
  getTransactions: async (filters?: {
    category_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<Transaction[]> => {
    const params = new URLSearchParams();
    if (filters?.category_id) params.append('category_id', filters.category_id);
    if (filters?.start_date) params.append('start_date', filters.start_date);
    if (filters?.end_date) params.append('end_date', filters.end_date);
    
    const response = await api.get(`/transactions?${params.toString()}`);
    return response.data;
  },

  createTransaction: async (data: {
    category_id: string;
    amount: number;
    description?: string;
    date?: string;
  }): Promise<Transaction> => {
    const response = await api.post('/transactions', data);
    return response.data;
  },

  // Category Limits
  getCategoryLimits: async (filters?: {
    category_id?: string;
    month?: number;
    year?: number;
  }): Promise<CategoryLimit[]> => {
    const params = new URLSearchParams();
    if (filters?.category_id) params.append('category_id', filters.category_id);
    if (filters?.month) params.append('month', filters.month.toString());
    if (filters?.year) params.append('year', filters.year.toString());
    
    const response = await api.get(`/category-limits?${params.toString()}`);
    return response.data;
  },

  createCategoryLimit: async (data: {
    category_id: string;
    limit: number;
    month: number;
    year: number;
  }): Promise<CategoryLimit> => {
    const response = await api.post('/category-limits', data);
    return response.data;
  },

  updateCategoryLimit: async (id: string, data: {
    category_id: string;
    limit: number;
    month: number;
    year: number;
  }): Promise<CategoryLimit> => {
    const response = await api.put(`/category-limits/${id}`, data);
    return response.data;
  },

  deleteCategoryLimit: async (id: string): Promise<void> => {
    await api.delete(`/category-limits/${id}`);
  },

  // Analytics
  getMonthlySummary: async (month: number, year: number): Promise<MonthlySummary> => {
    const response = await api.get(`/monthly-summary?month=${month}&year=${year}`);
    return response.data;
  },

  // Planned Expenses
  getPlannedExpenses: async (filters?: {
    category_id?: string;
    start_date?: string;
    end_date?: string;
    is_completed?: boolean;
  }): Promise<PlannedExpense[]> => {
    const params = new URLSearchParams();
    if (filters?.category_id) params.append('category_id', filters.category_id);
    if (filters?.start_date) params.append('start_date', filters.start_date);
    if (filters?.end_date) params.append('end_date', filters.end_date);
    if (filters?.is_completed !== undefined) params.append('is_completed', filters.is_completed.toString());
    
    const response = await api.get(`/planned-expenses?${params.toString()}`);
    return response.data;
  },

  createPlannedExpense: async (data: {
    category_id: string;
    amount: number;
    description?: string;
    planned_date: string;
  }): Promise<PlannedExpense> => {
    const response = await api.post('/planned-expenses', data);
    return response.data;
  },

  updatePlannedExpense: async (id: string, data: {
    category_id: string;
    amount: number;
    description?: string;
    planned_date: string;
  }): Promise<PlannedExpense> => {
    const response = await api.put(`/planned-expenses/${id}`, data);
    return response.data;
  },

  deletePlannedExpense: async (id: string): Promise<void> => {
    await api.delete(`/planned-expenses/${id}`);
  },

  // Planned Income
  getPlannedIncome: async (filters?: {
    month?: number;
    year?: number;
  }): Promise<PlannedIncome[]> => {
    const params = new URLSearchParams();
    if (filters?.month) params.append('month', filters.month.toString());
    if (filters?.year) params.append('year', filters.year.toString());
    
    const response = await api.get(`/planned-income?${params.toString()}`);
    return response.data;
  },

  createPlannedIncome: async (data: {
    amount: number;
    description?: string;
    month: number;
    year: number;
  }): Promise<PlannedIncome> => {
    const response = await api.post('/planned-income', data);
    return response.data;
  },

  updatePlannedIncome: async (id: string, data: {
    amount: number;
    description?: string;
    month: number;
    year: number;
  }): Promise<PlannedIncome> => {
    const response = await api.put(`/planned-income/${id}`, data);
    return response.data;
  },

  deletePlannedIncome: async (id: string): Promise<void> => {
    await api.delete(`/planned-income/${id}`);
  },

  // Notifications
  getNotifications: async (): Promise<Notification[]> => {
    const response = await api.get('/notifications');
    return response.data;
  },

  createNotification: async (data: {
    type: 'daily_reminder' | 'limit_warning' | 'limit_exceeded' | 'income_reminder';
    title: string;
    message: string;
  }): Promise<Notification> => {
    const response = await api.post('/notifications', data);
    return response.data;
  },

  markNotificationAsRead: async (id: string): Promise<void> => {
    await api.put(`/notifications/${id}/read`);
  },

  getNotificationStats: async (): Promise<NotificationStats> => {
    const response = await api.get('/notifications/stats');
    return response.data;
  },

  checkDailyReminder: async (): Promise<void> => {
    await api.post('/notifications/check-daily');
  },

  checkLimitWarnings: async (): Promise<void> => {
    await api.post('/notifications/check-limits');
  },
};
