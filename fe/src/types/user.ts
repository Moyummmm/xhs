// 用户信息
export interface User {
  id: number;
  nickname: string;
  avatar?: string;
  bio?: string;
  gender?: number; // 0: 未知, 1: 男, 2: 女
  birthday?: string;
  follow_count?: number;
  follower_count?: number;
  like_count?: number;
  note_count?: number;
  is_following?: boolean;
  created_at?: string;
  updated_at?: string;
}

// 登录请求
export interface LoginRequest {
  username: string;
  password: string;
}

// 注册请求
export interface RegisterRequest {
  nickname: string;
  password: string;
}

// 登录/注册响应
export interface AuthResponse {
  token: string;
  user: User;
}

// 更新用户信息
export interface UpdateUserRequest {
  nickname?: string;
  avatar?: string;
  bio?: string;
  gender?: number;
  birthday?: string;
}
