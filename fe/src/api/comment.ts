import { http } from './request';
import { API_PATHS } from '@/constants/api';
import type { Comment, CreateCommentRequest, CommentListParams, PaginatedList } from '@/types';

/**
 * 获取评论列表
 */
export const getComments = (params: CommentListParams) => {
  return http.get<PaginatedList<Comment>>(API_PATHS.NOTE_COMMENTS(params.note_id), {
    params: { page: params.page, page_size: params.page_size, type: params.type },
  });
};

/**
 * 创建评论
 */
export const createComment = (data: CreateCommentRequest) => {
  return http.post<Comment>(API_PATHS.NOTE_COMMENTS(data.note_id), data);
};

/**
 * 删除评论
 */
export const deleteComment = (id: number) => {
  return http.delete(`${API_PATHS.COMMENTS}/${id}`);
};

/**
 * 点赞评论
 */
export const likeComment = (id: number) => {
  return http.post(API_PATHS.COMMENT_LIKE(id));
};

/**
 * 取消点赞评论
 */
export const unlikeComment = (id: number) => {
  return http.delete(API_PATHS.COMMENT_LIKE(id));
};
