import { apiClient } from '../helpers/api-client';
import { authHelpers, API_PATHS } from '../helpers/auth-helpers';
import { validators } from '../helpers/validators';

describe('Notifications API', () => {
  let userA: { userId: number; token: string };

  beforeAll(async () => {
    const users = await authHelpers.initTestUsers();
    userA = users.userA;
    authHelpers.setToken(userA.token);
  });

  describe('GET /notifications', () => {
    it('should get notification list', async () => {
      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.NOTIFICATIONS);

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
    });

    it('should require authentication', async () => {
      authHelpers.setToken(null);

      await expect(apiClient.get(API_PATHS.NOTIFICATIONS)).rejects.toThrow();

      // 恢复 token
      authHelpers.setToken(userA.token);
    });
  });

  describe('PUT /notifications/:id/read', () => {
    it('should mark notification as read', async () => {
      // 先获取通知列表
      const response = await apiClient.get<{
        list: Array<{ id: number }>;
      }>(API_PATHS.NOTIFICATIONS);

      if (response.list && response.list.length > 0) {
        const notificationId = response.list[0].id;
        const markReadResponse = await apiClient.put(
          API_PATHS.NOTIFICATION_READ(notificationId)
        );
        validators.validateSuccessResponse(markReadResponse);
      } else {
        // 如果没有通知，测试应该返回空列表而不是报错
        expect(response.list).toEqual([]);
      }
    });
  });

  describe('PUT /notifications/read-all', () => {
    it('should mark all notifications as read', async () => {
      const response = await apiClient.put(API_PATHS.NOTIFICATIONS_READ_ALL);
      validators.validateSuccessResponse(response);
    });
  });
});
