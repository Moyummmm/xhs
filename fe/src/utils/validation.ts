import { REGEX } from '@/constants/reg';

/**
 * 表单校验器
 */
export const validators = {
  password: (value: string): string => {
    if (!value) return '请输入密码';
    if (value.length < 8) return '密码至少 8 位';
    if (!/(?=.*[a-zA-Z])(?=.*\d)/.test(value)) {
      return '密码必须包含字母和数字';
    }
    return '';
  },

  code: (value: string): string => {
    if (!value) return '请输入验证码';
    if (!REGEX.CODE.test(value)) return '验证码为 6 位数字';
    return '';
  },

  nickname: (value: string): string => {
    if (!value) return '请输入昵称';
    if (value.length < 2) return '昵称至少 2 个字符';
    if (value.length > 20) return '昵称不能超过 20 个字符';
    return '';
  },

  required: (value: unknown, fieldName = '此字段'): string => {
    if (!value && value !== 0) return `请输入${fieldName}`;
    return '';
  },
};

/**
 * 解析路由参数中的 ID
 * @param id 路由参数
 * @returns 正整数 ID，非法时返回 null
 */
export const parseRouteId = (id: string | undefined): number | null => {
  if (!id) return null;
  const num = Number(id);
  return Number.isNaN(num) || num <= 0 ? null : num;
};
