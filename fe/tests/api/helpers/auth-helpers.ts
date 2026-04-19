import { apiClient } from './api-client';
import { generateUser, generateNote, type TestUser, type UserData } from './test-data';

const API_PATHS = {
  AUTH_LOGIN: '/auth/login',
  AUTH_REGISTER: '/auth/register',
  AUTH_LOGOUT: '/auth/logout',
  AUTH_REFRESH: '/auth/refresh',
  USER_ME: '/users/me',
  USER_INFO: (id: number) => `/users/${id}`,
  USER_UPDATE: (id: number) => `/users/${id}`,
  USER_NOTES: (id: number) => `/users/${id}/notes`,
  USER_FOLLOW: (id: number) => `/users/${id}/follow`,
  USER_FOLLOWERS: (id: number) => `/users/${id}/followers`,
  USER_FOLLOWINGS: (id: number) => `/users/${id}/followings`,
  NOTES: '/notes',
  NOTE_DETAIL: (id: number) => `/notes/${id}`,
  NOTES_FEED: '/notes/feed',
  NOTES_SEARCH: '/notes/search',
  NOTE_LIKE: (id: number) => `/notes/${id}/like`,
  NOTE_COLLECT: (id: number) => `/notes/${id}/collect`,
  NOTE_COMMENTS: (id: number) => `/notes/${id}/comments`,
  COMMENTS: '/comments',
  COMMENT_LIKE: (id: number) => `/comments/${id}/like`,
  UPLOAD_IMAGE: '/upload/image',
  UPLOAD_VIDEO: '/upload/video',
  NOTIFICATIONS: '/notifications',
  NOTIFICATION_READ: (id: number) => `/notifications/${id}/read`,
  NOTIFICATIONS_READ_ALL: '/notifications/read-all',
  COLLECTS: '/collects',
  COLLECTS_COUNT: '/collects/count',
};

// Auth 响应类型
interface AuthResponse {
  token: string;
  user: {
    id: number;
    username: string;
    nickname: string;
    avatar?: string;
    bio?: string;
  };
}

// 全局测试用户存储
let testUserA: TestUser | null = null;
let testUserB: TestUser | null = null;
let createdNoteIds: number[] = [];

// 注册用户
async function registerUser(userData: UserData & { username?: string }): Promise<AuthResponse> {
  return apiClient.post<AuthResponse>(API_PATHS.AUTH_REGISTER, {
    nickname: userData.nickname,
    password: userData.password,
  });
}

// 登录用户
async function loginUser(userData: { username: string; password: string }): Promise<AuthResponse> {
  return apiClient.post<AuthResponse>(API_PATHS.AUTH_LOGIN, {
    username: userData.username,
    password: userData.password,
  });
}

// 创建测试用户（注册+登录）
async function createTestUser(): Promise<TestUser> {
  const userData = generateUser();
  // 先注册
  const registerResponse = await registerUser(userData);

  // 再登录获取完整用户信息
  const loginResponse = await loginUser({
    username: userData.username,
    password: userData.password,
  });

  return {
    userId: loginResponse.user.id,
    username: userData.username,
    password: userData.password,
    token: loginResponse.token,
  };
}

// 删除测试用户
async function deleteTestUser(userId: number): Promise<void> {
  const token = apiClient.getToken();
  try {
    await apiClient.delete(`${API_PATHS.USER_INFO(userId)}`);
  } catch (error) {
    // 忽略错误，可能是用户已经被删除
  } finally {
    // 恢复 token
    if (token) {
      apiClient.setToken(token);
    }
  }
}

// 清理测试笔记
async function cleanupTestNotes(): Promise<void> {
  for (const noteId of createdNoteIds) {
    try {
      await apiClient.delete(API_PATHS.NOTE_DETAIL(noteId));
    } catch {
      // 忽略错误
    }
  }
  createdNoteIds = [];
}

// 导出 auth helpers
export const authHelpers = {
  // 初始化两个测试用户
  async initTestUsers(): Promise<{ userA: TestUser; userB: TestUser }> {
    if (!testUserA) {
      testUserA = await createTestUser();
    }
    if (!testUserB) {
      testUserB = await createTestUser();
    }
    apiClient.setToken(testUserA.token);
    return { userA: testUserA, userB: testUserB };
  },

  // 获取测试用户 A
  getUserA(): TestUser | null {
    return testUserA;
  },

  // 获取测试用户 B
  getUserB(): TestUser | null {
    return testUserB;
  },

  // 设置当前 token
  setToken(token: string | null): void {
    apiClient.setToken(token);
  },

  // 获取当前 token
  getToken(): string | null {
    return apiClient.getToken();
  },

  // 创建笔记并记录 ID 用于清理
  async createTestNote(): Promise<number> {
    const noteData = generateNote();
    const note = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
    createdNoteIds.push(note.id);
    return note.id;
  },

  // 清理所有测试数据
  async cleanup(): Promise<void> {
    await cleanupTestNotes();
    if (testUserA) {
      await deleteTestUser(testUserA.userId);
      testUserA = null;
    }
    if (testUserB) {
      await deleteTestUser(testUserB.userId);
      testUserB = null;
    }
    apiClient.setToken(null);
  },
};

// 导出 API_PATHS 供测试用
export { API_PATHS };
