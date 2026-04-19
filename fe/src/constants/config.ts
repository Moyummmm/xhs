// 配置常量
export const CONFIG = {
  // 分页
  DEFAULT_PAGE_SIZE: 20,
  MAX_PAGE_SIZE: 100,

  // 图片
  MAX_IMAGE_COUNT: 9,
  MAX_IMAGE_SIZE: 10 * 1024 * 1024, // 10MB

  // 视频
  MAX_VIDEO_SIZE: 100 * 1024 * 1024, // 100MB
  ALLOWED_VIDEO_TYPES: ['video/mp4', 'video/quicktime'],

  // Token
  TOKEN_KEY: 'xhs_token',
  USER_KEY: 'xhs_user',

  // 防抖/节流
  DEBOUNCE_DELAY: 300,
  THROTTLE_DELAY: 500,
} as const;
