import { http } from './request';
import { API_PATHS } from '@/constants/api';
import type { Notification, NotificationListParams, PaginatedList } from '@/types';

/**
 * 获取通知列表
 */
export const getNotifications = (params: NotificationListParams = {}) => {
  return http.get<PaginatedList<Notification>>(API_PATHS.NOTIFICATIONS, { params });
};

/**
 * 标记通知已读
 */
export const markNotificationRead = (id: number) => {
  return http.put(API_PATHS.NOTIFICATION_READ(id));
};

/**
 * 标记所有通知已读
 */
export const markAllNotificationsRead = () => {
  return http.put(API_PATHS.NOTIFICATIONS_READ_ALL);
};
