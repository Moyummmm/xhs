import { http } from './request';
import type { Topic, Note, PaginatedList } from '@/types';

/**
 * 获取话题列表
 */
export const getTopics = (page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Topic>>('/topics', {
    params: { page, page_size: pageSize },
  });
};

/**
 * 获取话题详情
 */
export const getTopicDetail = (id: number) => {
  return http.get<Topic>(`/topics/${id}`);
};

/**
 * 获取话题下的笔记列表
 */
export const getTopicNotes = (id: number, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Note>>(`/topics/${id}/notes`, {
    params: { page, page_size: pageSize },
  });
};
