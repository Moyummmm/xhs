// 正则表达式常量
export const REGEX = {
  PASSWORD: /^(?=.*[a-zA-Z])(?=.*\d).{8,}$/,
  CODE: /^\d{6}$/,
  URL: /^https?:\/\/.+/,
} as const;
