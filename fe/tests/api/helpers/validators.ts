// 响应校验工具

import { AxiosResponse } from 'axios';

interface ApiResponse<T = unknown> {
  code: number;
  msg: string;
  data: T;
}

// 校验成功响应
export function validateSuccessResponse(response: unknown): void {
  expect(response).toBeDefined();
  expect(response).not.toBeNull();
}

// 校验分页列表
export function validatePaginatedList(response: unknown): void {
  expect(response).toBeDefined();
  const data = response as { list?: unknown[]; pagination?: unknown };
  expect(Array.isArray(data.list)).toBe(true);
  expect(data.pagination).toBeDefined();
  const pagination = data.pagination as { total?: number; page?: number; page_size?: number; has_more?: boolean };
  expect(typeof pagination.total).toBe('number');
  expect(typeof pagination.page).toBe('number');
  expect(typeof pagination.page_size).toBe('number');
  expect(typeof pagination.has_more).toBe('boolean');
}

// 校验用户对象
export function validateUser(user: unknown): void {
  expect(user).toBeDefined();
  const u = user as { id?: number; username?: string; nickname?: string };
  expect(typeof u.id).toBe('number');
  expect(typeof u.username).toBe('string');
  expect(typeof u.nickname).toBe('string');
}

// 校验笔记对象
export function validateNote(note: unknown): void {
  expect(note).toBeDefined();
  const n = note as { id?: number; title?: string; content?: string };
  expect(typeof n.id).toBe('number');
  expect(typeof n.title).toBe('string');
  expect(typeof n.content).toBe('string');
}

// 校验评论对象
export function validateComment(comment: unknown): void {
  expect(comment).toBeDefined();
  const c = comment as { id?: number; content?: string };
  expect(typeof c.id).toBe('number');
  expect(typeof c.content).toBe('string');
}

// 校验错误响应
export function validateErrorResponse(
  response: AxiosResponse<ApiResponse> | unknown,
  expectedCode?: number
): void {
  if (expectedCode !== undefined) {
    const r = response as AxiosResponse<ApiResponse>;
    expect(r.data).toBeDefined();
    expect(r.data.code).toBe(expectedCode);
  }
}

// 校验收藏对象
export function validateCollect(collect: unknown): void {
  expect(collect).toBeDefined();
  const c = collect as { id?: number; note_id?: number };
  expect(typeof c.id).toBe('number');
  expect(typeof c.note_id).toBe('number');
}

// 校验通知对象
export function validateNotification(notification: unknown): void {
  expect(notification).toBeDefined();
  const n = notification as { id?: number; type?: string; content?: string };
  expect(typeof n.id).toBe('number');
  expect(typeof n.type).toBe('string');
}

// 导出所有校验器
export const validators = {
  validateSuccessResponse,
  validatePaginatedList,
  validateUser,
  validateNote,
  validateComment,
  validateErrorResponse,
  validateCollect,
  validateNotification,
};
