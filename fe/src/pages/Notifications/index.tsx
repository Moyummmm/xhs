import React from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Button, message } from 'antd';
import { getNotifications, markAllNotificationsRead } from '@/api/notification';
import { formatRelativeTime } from '@/utils';
import { Empty } from '@/components/Common';
import UserAvatar from '@/components/UserAvatar';
import type { Notification, PaginatedList } from '@/types';

const NotificationsPage: React.FC = () => {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery<PaginatedList<Notification>>({
    queryKey: ['notifications'],
    queryFn: () => getNotifications({ page: 1, page_size: 50 }),
  });

  const markAllReadMutation = useMutation({
    mutationFn: markAllNotificationsRead,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
      message.success('已全部标记为已读');
    },
  });

  const notifications = data?.list || [];
  const unreadCount = notifications.filter((n: Notification) => !n.is_read).length;

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'like':
        return '\u2764\ufe0f';
      case 'comment':
        return '\ud83d\udcac';
      case 'follow':
        return '\ud83d\udc64';
      case 'system':
        return '\ud83d\udce2';
      default:
        return '\ud83d\udd14';
    }
  };

  return (
    <div className="container-custom py-4">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-bold">通知</h1>
        {unreadCount > 0 && (
          <Button
            size="small"
            onClick={() => markAllReadMutation.mutate()}
            loading={markAllReadMutation.isPending}
          >
            全部已读
          </Button>
        )}
      </div>

      <div className="bg-white rounded-xl divide-y divide-gray-50">
        {isLoading ? (
          <div className="p-4 space-y-4">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3 animate-pulse">
                <div className="w-10 h-10 bg-gray-200 rounded-full" />
                <div className="flex-1 space-y-2">
                  <div className="h-4 bg-gray-200 rounded w-3/4" />
                  <div className="h-3 bg-gray-200 rounded w-1/2" />
                </div>
              </div>
            ))}
          </div>
        ) : notifications.length > 0 ? (
          notifications.map((notification: Notification) => (
            <div
              key={notification.id}
              className={`flex items-start gap-3 p-4 ${
                !notification.is_read ? 'bg-primary-50/50' : ''
              }`}
            >
              <span className="text-2xl">{getTypeIcon(notification.type)}</span>
              <div className="flex-1 min-w-0">
                {notification.from_user && (
                  <UserAvatar user={notification.from_user} size="sm" className="mb-1" />
                )}
                <p className="text-sm text-gray-800">{notification.content}</p>
                <p className="text-xs text-gray-400 mt-1">
                  {formatRelativeTime(notification.created_at!)}
                </p>
              </div>
              {!notification.is_read && (
                <span className="w-2 h-2 bg-primary-500 rounded-full flex-shrink-0 mt-2" />
              )}
            </div>
          ))
        ) : (
          <Empty description="暂无通知" />
        )}
      </div>
    </div>
  );
};

export default NotificationsPage;
