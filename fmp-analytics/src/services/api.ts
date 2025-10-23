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

export interface CategoryLimit {
  id: string;
  category_id: string;
  limit: number;
  month: number;
  year: number;
  created_at: string;
  updated_at: string;
}

export interface LimitExceeded {
  id: string;
  category_id: string;
  limit: number;
  actual: number;
  month: number;
  year: number;
  created_at: string;
}

export const apiService = {
  // Categories
  getCategories: async (): Promise<Category[]> => {
    const response = await api.get('/categories');
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

  // Analytics
  getMonthlySummary: async (month: number, year: number): Promise<MonthlySummary> => {
    const response = await api.get(`/analytics/monthly-summary?month=${month}&year=${year}`);
    return response.data;
  },

  getCategorySummary: async (filters?: {
    category_id?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<CategorySummary[]> => {
    const params = new URLSearchParams();
    if (filters?.category_id) params.append('category_id', filters.category_id);
    if (filters?.start_date) params.append('start_date', filters.start_date);
    if (filters?.end_date) params.append('end_date', filters.end_date);
    
    const response = await api.get(`/analytics/category-summary?${params.toString()}`);
    return response.data;
  },

  getLimitExceeded: async (): Promise<LimitExceeded[]> => {
    const response = await api.get('/analytics/limit-exceeded');
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
};
