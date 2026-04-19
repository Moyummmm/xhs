import { apiClient } from '../helpers/api-client';
import { authHelpers, API_PATHS } from '../helpers/auth-helpers';
import { generateNote } from '../helpers/test-data';
import { validators } from '../helpers/validators';

describe('Collect API', () => {
  let userA: { userId: number; token: string };

  beforeAll(async () => {
    const users = await authHelpers.initTestUsers();
    userA = users.userA;
    authHelpers.setToken(userA.token);
  });

  let testNoteId: number;

  beforeEach(async () => {
    // 创建一个测试笔记用于收藏测试
    const noteData = generateNote();
    const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
    testNoteId = createResponse.id;
  });

  describe('POST /notes/:id/collect', () => {
    it('should collect a note', async () => {
      const response = await apiClient.post(API_PATHS.NOTE_COLLECT(testNoteId));
      validators.validateSuccessResponse(response);
    });
  });

  describe('DELETE /notes/:id/collect', () => {
    it('should uncollect a note', async () => {
      // 先收藏
      await apiClient.post(API_PATHS.NOTE_COLLECT(testNoteId));

      // 取消收藏
      const response = await apiClient.delete(API_PATHS.NOTE_COLLECT(testNoteId));
      validators.validateSuccessResponse(response);
    });
  });

  describe('GET /collects', () => {
    it('should get user collect list', async () => {
      // 先收藏笔记
      await apiClient.post(API_PATHS.NOTE_COLLECT(testNoteId));

      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.COLLECTS);

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
    });
  });

  describe('GET /collects/count', () => {
    it('should get user collect count', async () => {
      const response = await apiClient.get<{ count: number }>(API_PATHS.COLLECTS_COUNT);
      validators.validateSuccessResponse(response);
      expect(typeof response.count).toBe('number');
    });
  });
});
