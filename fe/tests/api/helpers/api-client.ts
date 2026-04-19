import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';

// API 响应类型
interface ApiResponse<T = unknown> {
  code: number;
  msg: string;
  data: T;
}

// 创建 Axios 实例
const instance = axios.create({
  baseURL: process.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Token 管理
let currentToken: string | null = null;

export const apiClient = {
  setToken(token: string | null) {
    currentToken = token;
  },

  getToken(): string | null {
    return currentToken;
  },

  // 通用请求方法
  async get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await instance.get<ApiResponse<T>>(url, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
    return this.handleResponse(response);
  },

  async post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await instance.post<ApiResponse<T>>(url, data, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
    return this.handleResponse(response);
  },

  async put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await instance.put<ApiResponse<T>>(url, data, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
    return this.handleResponse(response);
  },

  async delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await instance.delete<ApiResponse<T>>(url, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
    return this.handleResponse(response);
  },

  async upload<T = unknown>(
    url: string,
    file: File,
    onProgress?: (percent: number) => void
  ): Promise<T> {
    const formData = new FormData();
    formData.append('file', file);
    const response = await instance.post<ApiResponse<T>>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(percent);
        }
      },
    });
    return this.handleResponse(response);
  },

  // 处理响应，解包 { code, msg, data }
  handleResponse<T>(response: AxiosResponse<ApiResponse<T>>): T {
    const { data } = response;
    if (data?.code !== undefined) {
      if (data.code === 0 || data.code === 200) {
        return data.data as T;
      }
      throw new Error(data.msg || `API Error: code=${data.code}`);
    }
    return data as T;
  },

  // 获取原始响应（用于测试负向场景）
  async getRaw<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<ApiResponse<T>>> {
    return instance.get<ApiResponse<T>>(url, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
  },

  async postRaw<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<AxiosResponse<ApiResponse<T>>> {
    return instance.post<ApiResponse<T>>(url, data, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
  },

  async deleteRaw<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<ApiResponse<T>>> {
    return instance.delete<ApiResponse<T>>(url, {
      ...config,
      headers: {
        ...config?.headers,
        ...(currentToken ? { Authorization: `Bearer ${currentToken}` } : {}),
      },
    });
  },
};

// 导出 http 别名，方便直接调用
export const http = apiClient;
