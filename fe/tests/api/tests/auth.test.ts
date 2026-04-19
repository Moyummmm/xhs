import { apiClient } from '../helpers/api-client';
import { authHelpers, API_PATHS } from '../helpers/auth-helpers';
import { generateUser } from '../helpers/test-data';
import { validators } from '../helpers/validators';

describe('Auth API', () => {
  describe('POST /auth/register', () => {
    it('should register a new user successfully', async () => {
      const userData = generateUser();
      const response = await apiClient.post<{ token: string; user: { id: number; username: string } }>(
        API_PATHS.AUTH_REGISTER,
        {
          nickname: userData.nickname,
          password: userData.password,
        }
      );

      validators.validateSuccessResponse(response);
      expect(response.token).toBeDefined();
      expect(response.user.id).toBeDefined();
      expect(response.user.username).toBeDefined();
    });

    it('should reject duplicate username (if applicable)', async () => {
      // 先注册一个用户
      const userData = generateUser();
      await apiClient.post(API_PATHS.AUTH_REGISTER, {
        nickname: userData.nickname,
        password: userData.password,
      });

      // 尝试用相同昵称注册（如果后端检查昵称重复）
      await expect(
        apiClient.post(API_PATHS.AUTH_REGISTER, {
          nickname: userData.nickname,
          password: userData.password,
        })
      ).rejects.toThrow();
    });

    it('should reject registration with missing fields', async () => {
      await expect(
        apiClient.post(API_PATHS.AUTH_REGISTER, {})
      ).rejects.toThrow();
    });
  });

  describe('POST /auth/login', () => {
    it('should login successfully with valid credentials', async () => {
      const userA = authHelpers.getUserA();
      if (!userA) throw new Error('Test user not initialized');

      const response = await apiClient.post<{ token: string; user: { id: number } }>(
        API_PATHS.AUTH_LOGIN,
        {
          username: userA.username,
          password: userA.password,
        }
      );

      validators.validateSuccessResponse(response);
      expect(response.token).toBeDefined();
      expect(response.user.id).toBe(userA.userId);
    });

    it('should reject login with invalid credentials', async () => {
      await expect(
        apiClient.post(API_PATHS.AUTH_LOGIN, {
          username: 'nonexistent_user',
          password: 'wrong_password',
        })
      ).rejects.toThrow();
    });

    it('should reject login with missing fields', async () => {
      await expect(
        apiClient.post(API_PATHS.AUTH_LOGIN, {})
      ).rejects.toThrow();
    });
  });

  describe('POST /auth/logout', () => {
    it('should logout successfully', async () => {
      const userA = authHelpers.getUserA();
      if (!userA) throw new Error('Test user not initialized');

      // 确保已登录
      authHelpers.setToken(userA.token);

      // 登出应该成功
      await expect(apiClient.post(API_PATHS.AUTH_LOGOUT)).resolves.toBeDefined();
    });
  });

  describe('POST /auth/refresh', () => {
    it('should refresh token successfully', async () => {
      const userA = authHelpers.getUserA();
      if (!userA) throw new Error('Test user not initialized');

      authHelpers.setToken(userA.token);

      const response = await apiClient.post<{ token: string }>(API_PATHS.AUTH_REFRESH);
      validators.validateSuccessResponse(response);
      expect(response.token).toBeDefined();
    });
  });

  describe('GET /users/me', () => {
    it('should get current user info', async () => {
      const userA = authHelpers.getUserA();
      if (!userA) throw new Error('Test user not initialized');

      authHelpers.setToken(userA.token);

      const response = await apiClient.get<{
        id: number;
        username: string;
        nickname: string;
      }>(API_PATHS.USER_ME);

      validators.validateSuccessResponse(response);
      expect(response.id).toBe(userA.userId);
    });

    it('should return 401 when not authenticated', async () => {
      authHelpers.setToken(null);

      await expect(apiClient.get(API_PATHS.USER_ME)).rejects.toThrow();
    });
  });
});
