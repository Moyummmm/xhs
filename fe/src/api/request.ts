import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { message } from 'antd';
import { CONFIG } from '@/constants/config';
import { storage } from '@/utils/storage';

// 创建 Axios 实例
const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    // 兜底：拦截包含 NaN 的非法请求
    if (config.url && config.url.includes('NaN')) {
      return Promise.reject(new Error('非法请求参数'));
    }
    const token = storage.get<string>(CONFIG.TOKEN_KEY);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器
instance.interceptors.response.use(
  (response) => {
    const { data } = response;
    // 如果后端返回统一格式 { code, msg, data }
    if (data?.code !== undefined) {
      if (data.code === 0 || data.code === 200) {
        return data.data ?? data;
      }
      message.error(data.msg || '请求失败');
      return Promise.reject(new Error(data.msg));
    }
    return data;
  },
  (error: AxiosError) => {
    if (error.code === 'ECONNABORTED') {
      message.error('请求超时，请检查网络连接');
    } else if (!error.response) {
      message.error('网络连接失败，请检查网络');
    } else {
      const { status, data } = error.response as { status: number; data?: { msg?: string } };
      switch (status) {
        case 401:
          // Token 过期，清除并跳转登录
          storage.remove(CONFIG.TOKEN_KEY);
          storage.remove(CONFIG.USER_KEY);
          if (!window.location.pathname.includes('/login')) {
            window.location.href = '/login';
          }
          break;
        case 403:
          message.error('无权限访问');
          break;
        case 404:
          message.error('资源不存在');
          break;
        case 500:
          message.error('服务器错误，请稍后重试');
          break;
        default:
          message.error((data as { msg?: string })?.msg || '请求失败');
      }
    }
    return Promise.reject(error);
  }
);

export default instance;

// 封装常用请求方法 - 拦截器已解包 data，直接返回 T
export const http = {
  get: async <T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    const response = await instance.get<T>(url, config);
    return response as unknown as T;
  },

  post: async <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await instance.post<T>(url, data, config);
    return response as unknown as T;
  },

  put: async <T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> => {
    const response = await instance.put<T>(url, data, config);
    return response as unknown as T;
  },

  delete: async <T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    const response = await instance.delete<T>(url, config);
    return response as unknown as T;
  },

  upload: async <T = unknown>(
    url: string,
    file: File,
    onProgress?: (percent: number) => void
  ): Promise<T> => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await instance.post<T>(url, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(percent);
        }
      },
    });
    return response as unknown as T;
  },
};
