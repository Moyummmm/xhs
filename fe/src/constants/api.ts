// API 路径常量
export const API_PATHS = {
  // 认证
  AUTH_LOGIN: '/auth/login',
  AUTH_REGISTER: '/auth/register',
  AUTH_LOGOUT: '/auth/logout',
  AUTH_REFRESH: '/auth/refresh',

  // 用户
  USER_ME: '/users/me',
  USER_INFO: (id: number) => `/users/${id}`,
  USER_UPDATE: (id: number) => `/users/${id}`,
  USER_NOTES: (id: number) => `/users/${id}/notes`,
  USER_LIKES: (id: number) => `/users/${id}/likes`,
  USER_FOLLOW: (id: number) => `/users/${id}/follow`,
  USER_FOLLOWERS: (id: number) => `/users/${id}/followers`,
  USER_FOLLOWINGS: (id: number) => `/users/${id}/followings`,

  // 笔记
  NOTES: '/notes',
  NOTE_DETAIL: (id: number) => `/notes/${id}`,
  NOTES_FEED: '/notes/feed',
  NOTES_SEARCH: '/notes/search',

  // 互动
  NOTE_LIKE: (id: number) => `/notes/${id}/like`,
  NOTE_COLLECT: (id: number) => `/notes/${id}/collect`,
  NOTE_COMMENTS: (id: number) => `/notes/${id}/comments`,
  USER_COLLECTIONS: '/collects',

  // 评论
  COMMENTS: '/comments',
  COMMENT_LIKE: (id: number) => `/comments/${id}/like`,

  // 话题
  TOPICS: '/topics',
  TOPIC_DETAIL: (id: number) => `/topics/${id}`,

  // 搜索
  SEARCH: '/search',

  // 通知
  NOTIFICATIONS: '/notifications',
  NOTIFICATION_READ: (id: number) => `/notifications/${id}/read`,
  NOTIFICATIONS_READ_ALL: '/notifications/read-all',

  // 上传
  UPLOAD_IMAGE: '/upload/image',
  UPLOAD_VIDEO: '/upload/video',
} as const;
