import { authHelpers } from './helpers/auth-helpers';

// 全局 setup
beforeAll(async () => {
  // 初始化测试用户
  await authHelpers.initTestUsers();
});

// 全局 teardown
afterAll(async () => {
  // 清理所有测试数据
  await authHelpers.cleanup();
});

// 设置 Jest 超时时间
jest.setTimeout(30000);
