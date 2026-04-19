import { create } from 'zustand';
import { getNotifications, markAllNotificationsRead as markAllReadApi } from '@/api/notification';
import type { Notification, PaginatedList } from '@/types';

interface NotificationState {
  notifications: Notification[];
  unreadCount: number;
  isLoading: boolean;

  // Actions
  fetchNotifications: (page?: number) => Promise<void>;
  markAllRead: () => Promise<void>;
  addNotification: (notification: Notification) => void;
  resetUnreadCount: () => void;
}

export const useNotificationStore = create<NotificationState>((set, get) => ({
  notifications: [],
  unreadCount: 0,
  isLoading: false,

  fetchNotifications: async (page = 1) => {
    set({ isLoading: true });
    try {
      const data: PaginatedList<Notification> = await getNotifications({ page, page_size: 20 });
      const unreadCount = data.list.filter((n: Notification) => !n.is_read).length;
      set({
        notifications: page === 1 ? data.list : [...get().notifications, ...data.list],
        unreadCount,
        isLoading: false,
      });
    } catch {
      set({ isLoading: false });
    }
  },

  markAllRead: async () => {
    try {
      await markAllReadApi();
      set({ unreadCount: 0 });
      // 更新所有通知为已读
      set((state) => ({
        notifications: state.notifications.map((n) => ({ ...n, is_read: true })),
      }));
    } catch {
      // Ignore errors
    }
  },

  addNotification: (notification: Notification) => {
    set((state) => ({
      notifications: [notification, ...state.notifications],
      unreadCount: state.unreadCount + 1,
    }));
  },

  resetUnreadCount: () => set({ unreadCount: 0 }),
}));
