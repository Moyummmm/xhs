import { http } from './request';
import { API_PATHS } from '@/constants/api';
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '@/types';

/**
 * 用户登录
 */
export const login = (data: LoginRequest) => {
  return http.post<AuthResponse>(API_PATHS.AUTH_LOGIN, data);
};

/**
 * 用户注册
 */
export const register = (data: RegisterRequest) => {
  return http.post<AuthResponse>(API_PATHS.AUTH_REGISTER, data);
};

/**
 * 用户登出
 */
export const logout = () => {
  return http.post(API_PATHS.AUTH_LOGOUT);
};

/**
 * 刷新 Token
 */
export const refreshToken = () => {
  return http.post<{ token: string }>(API_PATHS.AUTH_REFRESH);
};

/**
 * 获取当前用户信息
 */
export const getCurrentUser = () => {
  return http.get<User>(API_PATHS.USER_ME).catch(() => {
    // 如果获取失败，尝试从本地存储获取
    return null;
  });
};
