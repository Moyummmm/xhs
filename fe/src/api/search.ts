import { http } from './request';
import type { Note, User, Topic, PaginatedList } from '@/types';

/**
 * 搜索笔记
 */
export const searchNotes = (keyword: string, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Note>>('/notes/search', {
    params: { keyword, page, page_size: pageSize },
  });
};

/**
 * 搜索用户
 */
export const searchUsers = (keyword: string, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<User>>('/users/search', {
    params: { keyword, page, page_size: pageSize },
  });
};

/**
 * 搜索话题
 */
export const searchTopics = (keyword: string, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Topic>>('/topics/search', {
    params: { keyword, page, page_size: pageSize },
  });
};
