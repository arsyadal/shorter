import axios, { AxiosResponse, AxiosError, InternalAxiosRequestConfig } from 'axios';
import type {
  CreateURLRequest,
  CreateURLResponse,
  URLListResponse,
  ClickStats,
  ApiError,
} from '@/types';

// Create axios instance
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  (error: AxiosError<ApiError>) => {
    if (error.response?.data?.error) {
      throw new Error(error.response.data.error);
    }
    throw new Error(error.message || 'An error occurred');
  }
);

export const urlApi = {
  // Create a short URL
  createURL: async (data: CreateURLRequest): Promise<CreateURLResponse> => {
    const response = await api.post('/shorten', data);
    return response.data;
  },

  // Get all URLs with pagination
  getURLs: async (page: number = 1, limit: number = 10): Promise<URLListResponse> => {
    const response = await api.get('/urls', {
      params: { page, limit },
    });
    return response.data;
  },

  // Get statistics for a specific URL
  getURLStats: async (shortCode: string): Promise<ClickStats> => {
    const response = await api.get(`/stats/${shortCode}`);
    return response.data;
  },

  // Health check
  healthCheck: async (): Promise<{ status: string; message: string }> => {
    const response = await axios.get(
      process.env.NEXT_PUBLIC_API_URL?.replace('/api', '') || 'http://localhost:8080' + '/health'
    );
    return response.data;
  },
};

export default api; 