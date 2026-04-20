import { http } from './request';
import { API_PATHS } from '@/constants/api';
import type { User, UpdateUserRequest, Note, PaginatedList } from '@/types';

/**
 * 获取用户信息
 */
export const getUserInfo = (id: number) => {
  return http.get<User>(API_PATHS.USER_INFO(id));
};

/**
 * 更新用户信息
 */
export const updateUserInfo = (id: number, data: UpdateUserRequest) => {
  return http.put<User>(API_PATHS.USER_UPDATE(id), data);
};

/**
 * 获取用户笔记列表
 */
export const getUserNotes = (id: number, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Note>>(API_PATHS.USER_NOTES(id), {
    params: { page, page_size: pageSize },
  });
};

/**
 * 获取用户赞过的笔记列表
 */
export const getUserLikes = (id: number, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Note>>(API_PATHS.USER_LIKES(id), {
    params: { page, page_size: pageSize },
  });
};

/**
 * 获取用户收藏列表
 */
export const getUserCollections = (id: number, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<Note>>(API_PATHS.USER_COLLECTIONS, {
    params: { user_id: id, page, page_size: pageSize },
  });
};

/**
 * 关注用户
 */
export const followUser = (id: number) => {
  return http.post(API_PATHS.USER_FOLLOW(id));
};

/**
 * 取消关注
 */
export const unfollowUser = (id: number) => {
  return http.delete(API_PATHS.USER_FOLLOW(id));
};

/**
 * 获取粉丝列表
 */
export const getFollowers = (id: number, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<User>>(API_PATHS.USER_FOLLOWERS(id), {
    params: { page, page_size: pageSize },
  });
};

/**
 * 获取关注列表
 */
export const getFollowings = (id: number, page = 1, pageSize = 20) => {
  return http.get<PaginatedList<User>>(API_PATHS.USER_FOLLOWINGS(id), {
    params: { page, page_size: pageSize },
  });
};
