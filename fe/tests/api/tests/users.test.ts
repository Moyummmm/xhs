import { apiClient } from '../helpers/api-client';
import { authHelpers, API_PATHS } from '../helpers/auth-helpers';
import { generateUser } from '../helpers/test-data';
import { validators } from '../helpers/validators';

describe('Users API', () => {
  let userA: { userId: number; token: string };
  let userB: { userId: number; token: string };

  beforeAll(async () => {
    const users = await authHelpers.initTestUsers();
    userA = users.userA;
    userB = users.userB;
    authHelpers.setToken(userA.token);
  });

  describe('GET /users/:id', () => {
    it('should get user info by id', async () => {
      const response = await apiClient.get<{
        id: number;
        username: string;
        nickname: string;
      }>(API_PATHS.USER_INFO(userB.userId));

      validators.validateSuccessResponse(response);
      validators.validateUser(response);
      expect(response.id).toBe(userB.userId);
    });

    it('should return 404 for non-existent user', async () => {
      await expect(
        apiClient.get(API_PATHS.USER_INFO(999999999))
      ).rejects.toThrow();
    });
  });

  describe('GET /users/:id/notes', () => {
    it('should get user notes', async () => {
      // 先创建一篇笔记
      const noteId = await authHelpers.createTestNote();

      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.USER_NOTES(userA.userId));

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
      expect(Array.isArray(response.list)).toBe(true);
    });
  });

  describe('GET /users/:id/followers', () => {
    it('should get user followers', async () => {
      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.USER_FOLLOWERS(userB.userId));

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
    });
  });

  describe('GET /users/:id/followings', () => {
    it('should get user followings', async () => {
      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.USER_FOLLOWINGS(userB.userId));

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
    });
  });

  describe('PUT /users/:id', () => {
    it('should update user info', async () => {
      const response = await apiClient.put<{
        id: number;
        nickname: string;
        bio: string;
      }>(API_PATHS.USER_UPDATE(userA.userId), {
        nickname: 'updated_nickname_' + Date.now(),
        bio: 'Updated bio',
      });

      validators.validateSuccessResponse(response);
      expect(response.id).toBe(userA.userId);
    });

    it('should return 403 when updating another user', async () => {
      await expect(
        apiClient.put(API_PATHS.USER_UPDATE(userB.userId), {
          nickname: 'hacker',
        })
      ).rejects.toThrow();
    });
  });

  describe('DELETE /users/:id', () => {
    it('should delete own user', async () => {
      // 创建一个临时用户来删除
      const tempUserData = generateUser();
      const tempAuth = await apiClient.post<{ token: string; user: { id: number } }>(
        API_PATHS.AUTH_REGISTER,
        {
          nickname: tempUserData.nickname,
          password: tempUserData.password,
        }
      );

      authHelpers.setToken(tempAuth.token);

      await expect(
        apiClient.delete(API_PATHS.USER_INFO(tempAuth.user.id))
      ).resolves.toBeDefined();

      // 恢复 token
      authHelpers.setToken(userA.token);
    });
  });

  describe('POST /users/:id/follow', () => {
    it('should follow another user', async () => {
      authHelpers.setToken(userA.token);

      const response = await apiClient.post(API_PATHS.USER_FOLLOW(userB.userId));
      validators.validateSuccessResponse(response);
    });

    it('should return error when following self', async () => {
      authHelpers.setToken(userA.token);

      await expect(
        apiClient.post(API_PATHS.USER_FOLLOW(userA.userId))
      ).rejects.toThrow();
    });
  });

  describe('DELETE /users/:id/follow', () => {
    it('should unfollow a user', async () => {
      authHelpers.setToken(userA.token);

      // 先关注
      await apiClient.post(API_PATHS.USER_FOLLOW(userB.userId));

      // 再取消关注
      const response = await apiClient.delete(API_PATHS.USER_FOLLOW(userB.userId));
      validators.validateSuccessResponse(response);
    });
  });
});
