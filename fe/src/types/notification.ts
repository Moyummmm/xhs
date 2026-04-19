import { User } from './user';

// 通知类型
export type NotificationType = 'like' | 'comment' | 'follow' | 'system';

// 通知信息
export interface Notification {
  id: number;
  user_id: number;
  type: NotificationType;
  title: string;
  content: string;
  target_id?: number;
  target_type?: string;
  is_read: boolean;
  from_user?: User;
  created_at?: string;
}

// 通知列表查询参数
export interface NotificationListParams {
  page?: number;
  page_size?: number;
  type?: NotificationType;
}
