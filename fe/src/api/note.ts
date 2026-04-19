import { http } from './request';
import { API_PATHS } from '@/constants/api';
import type {
  Note,
  CreateNoteRequest,
  UpdateNoteRequest,
  NoteListParams,
  PaginatedList,
} from '@/types';

/**
 * 获取笔记列表 (Feed)
 */
export const getNotesFeed = (params: NoteListParams = {}) => {
  return http.get<PaginatedList<Note>>(API_PATHS.NOTES_FEED, { params });
};

/**
 * 获取笔记详情
 */
export const getNoteDetail = (id: number) => {
  return http.get<Note>(API_PATHS.NOTE_DETAIL(id));
};

/**
 * 创建笔记
 */
export const createNote = (data: CreateNoteRequest) => {
  return http.post<Note>(API_PATHS.NOTES, data);
};

/**
 * 更新笔记
 */
export const updateNote = (id: number, data: UpdateNoteRequest) => {
  return http.put<Note>(API_PATHS.NOTE_DETAIL(id), data);
};

/**
 * 删除笔记
 */
export const deleteNote = (id: number) => {
  return http.delete(API_PATHS.NOTE_DETAIL(id));
};

/**
 * 搜索笔记
 */
export const searchNotes = (keyword: string, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Note>>(API_PATHS.NOTES_SEARCH, {
    params: { keyword, page, page_size: pageSize },
  });
};

/**
 * 点赞笔记
 */
export const likeNote = (id: number) => {
  return http.post(API_PATHS.NOTE_LIKE(id));
};

/**
 * 取消点赞
 */
export const unlikeNote = (id: number) => {
  return http.delete(API_PATHS.NOTE_LIKE(id));
};

/**
 * 收藏笔记
 */
export const collectNote = (id: number) => {
  return http.post(API_PATHS.NOTE_COLLECT(id));
};

/**
 * 取消收藏
 */
export const uncollectNote = (id: number) => {
  return http.delete(API_PATHS.NOTE_COLLECT(id));
};
