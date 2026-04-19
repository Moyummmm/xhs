import { apiClient } from '../helpers/api-client';
import { authHelpers, API_PATHS } from '../helpers/auth-helpers';
import { generateNote } from '../helpers/test-data';
import { validators } from '../helpers/validators';

describe('Notes API', () => {
  let userA: { userId: number; token: string };

  beforeAll(async () => {
    const users = await authHelpers.initTestUsers();
    userA = users.userA;
    authHelpers.setToken(userA.token);
  });

  describe('GET /notes/feed', () => {
    it('should get notes feed', async () => {
      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.NOTES_FEED);

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
      expect(Array.isArray(response.list)).toBe(true);
    });
  });

  describe('POST /notes', () => {
    it('should create a new note', async () => {
      const noteData = generateNote();
      const response = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);

      validators.validateSuccessResponse(response);
      expect(response.id).toBeDefined();
    });

    it('should reject note without title', async () => {
      await expect(
        apiClient.post(API_PATHS.NOTES, { content: 'test content' })
      ).rejects.toThrow();
    });
  });

  describe('GET /notes/:id', () => {
    it('should get note by id', async () => {
      // 先创建一个笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      const response = await apiClient.get<{
        id: number;
        title: string;
        content: string;
      }>(API_PATHS.NOTE_DETAIL(noteId));

      validators.validateSuccessResponse(response);
      validators.validateNote(response);
      expect(response.id).toBe(noteId);
    });

    it('should return 404 for non-existent note', async () => {
      await expect(
        apiClient.get(API_PATHS.NOTE_DETAIL(999999999))
      ).rejects.toThrow();
    });
  });

  describe('PUT /notes/:id', () => {
    it('should update own note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      // 更新笔记
      const response = await apiClient.put<{ id: number; title: string }>(
        API_PATHS.NOTE_DETAIL(noteId),
        { title: 'Updated Title' }
      );

      validators.validateSuccessResponse(response);
      expect(response.id).toBe(noteId);
    });
  });

  describe('DELETE /notes/:id', () => {
    it('should delete own note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      // 删除笔记
      await expect(apiClient.delete(API_PATHS.NOTE_DETAIL(noteId))).resolves.toBeDefined();
    });
  });

  describe('GET /notes/search', () => {
    it('should search notes by keyword', async () => {
      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.NOTES_SEARCH, {
        params: { keyword: 'test' },
      });

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
    });
  });

  describe('POST /notes/:id/like', () => {
    it('should like a note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      const response = await apiClient.post(API_PATHS.NOTE_LIKE(noteId));
      validators.validateSuccessResponse(response);
    });

    it('should not throw when liking already liked note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      // 点赞两次
      await apiClient.post(API_PATHS.NOTE_LIKE(noteId));
      await expect(apiClient.post(API_PATHS.NOTE_LIKE(noteId))).resolves.toBeDefined();
    });
  });

  describe('DELETE /notes/:id/like', () => {
    it('should unlike a note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      // 点赞
      await apiClient.post(API_PATHS.NOTE_LIKE(noteId));

      // 取消点赞
      const response = await apiClient.delete(API_PATHS.NOTE_LIKE(noteId));
      validators.validateSuccessResponse(response);
    });
  });

  describe('GET /notes/:id/comments', () => {
    it('should get note comments', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      const response = await apiClient.get<{
        list: unknown[];
        pagination: { total: number; page: number; page_size: number; has_more: boolean };
      }>(API_PATHS.NOTE_COMMENTS(noteId));

      validators.validateSuccessResponse(response);
      validators.validatePaginatedList(response);
    });
  });

  describe('POST /notes/:id/comments', () => {
    it('should create a comment', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      const response = await apiClient.post<{ id: number }>(API_PATHS.NOTE_COMMENTS(noteId), {
        content: 'Test comment',
      });

      validators.validateSuccessResponse(response);
      expect(response.id).toBeDefined();
    });
  });

  describe('DELETE /notes/:id/comments/:comment_id', () => {
    it('should delete own comment', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      // 创建评论
      const commentResponse = await apiClient.post<{ id: number }>(
        API_PATHS.NOTE_COMMENTS(noteId),
        { content: 'Test comment to delete' }
      );
      const commentId = commentResponse.id;

      // 删除评论
      await expect(
        apiClient.delete(`${API_PATHS.NOTE_COMMENTS(noteId)}/${commentId}`)
      ).resolves.toBeDefined();
    });
  });

  describe('POST /notes/:id/collect', () => {
    it('should collect a note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      const response = await apiClient.post(API_PATHS.NOTE_COLLECT(noteId));
      validators.validateSuccessResponse(response);
    });
  });

  describe('DELETE /notes/:id/collect', () => {
    it('should uncollect a note', async () => {
      // 先创建笔记
      const noteData = generateNote();
      const createResponse = await apiClient.post<{ id: number }>(API_PATHS.NOTES, noteData);
      const noteId = createResponse.id;

      // 收藏
      await apiClient.post(API_PATHS.NOTE_COLLECT(noteId));

      // 取消收藏
      const response = await apiClient.delete(API_PATHS.NOTE_COLLECT(noteId));
      validators.validateSuccessResponse(response);
    });
  });
});
