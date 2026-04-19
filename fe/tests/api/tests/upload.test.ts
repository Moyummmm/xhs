import { apiClient } from '../helpers/api-client';
import { authHelpers, API_PATHS } from '../helpers/auth-helpers';
import { generateImageFile, generateVideoFile } from '../helpers/test-data';

describe('Upload API', () => {
  let userA: { userId: number; token: string };

  beforeAll(async () => {
    const users = await authHelpers.initTestUsers();
    userA = users.userA;
    authHelpers.setToken(userA.token);
  });

  describe('POST /upload/image', () => {
    it('should upload an image', async () => {
      const file = generateImageFile();

      const response = await apiClient.upload<{ url: string }>(
        API_PATHS.UPLOAD_IMAGE,
        file
      );

      expect(response).toBeDefined();
      expect((response as { url?: string }).url).toBeDefined();
    });

    it('should reject upload without file', async () => {
      await expect(
        apiClient.post(API_PATHS.UPLOAD_IMAGE, {})
      ).rejects.toThrow();
    });

    it('should reject unauthenticated upload', async () => {
      authHelpers.setToken(null);
      const file = generateImageFile();

      await expect(
        apiClient.upload(API_PATHS.UPLOAD_IMAGE, file)
      ).rejects.toThrow();

      // 恢复 token
      authHelpers.setToken(userA.token);
    });
  });

  describe('POST /upload/video', () => {
    it('should upload a video', async () => {
      const file = generateVideoFile();

      const response = await apiClient.upload<{ url: string }>(
        API_PATHS.UPLOAD_VIDEO,
        file
      );

      expect(response).toBeDefined();
      expect((response as { url?: string }).url).toBeDefined();
    });

    it('should reject unauthenticated video upload', async () => {
      authHelpers.setToken(null);
      const file = generateVideoFile();

      await expect(
        apiClient.upload(API_PATHS.UPLOAD_VIDEO, file)
      ).rejects.toThrow();

      // 恢复 token
      authHelpers.setToken(userA.token);
    });
  });
});
