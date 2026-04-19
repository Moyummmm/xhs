import { faker } from '@faker-js/faker';

// 生成随机用户数据
export function generateUser(overrides?: Partial<UserData>): UserData {
  const id = faker.number.int({ min: 100000, max: 999999 });
  return {
    username: `test_${id}_${faker.string.alphanumeric(8)}`,
    password: 'Test123456',
    nickname: faker.person.fullName(),
    avatar: faker.image.avatar(),
    bio: faker.lorem.sentence(),
    ...overrides,
  };
}

// 生成随机笔记数据
export function generateNote(overrides?: Partial<NoteData>): NoteData {
  return {
    title: faker.lorem.sentence({ min: 3, max: 8 }),
    content: faker.lorem.paragraphs({ min: 1, max: 3 }),
    images: [
      faker.image.url({ width: 640, height: 480 }),
      faker.image.url({ width: 640, height: 480 }),
    ],
    topics: [faker.word.noun(), faker.word.noun()],
    ...overrides,
  };
}

// 生成随机评论数据
export function generateComment(overrides?: Partial<{ content: string }>): { content: string } {
  return {
    content: faker.lorem.sentence(),
    ...overrides,
  };
}

// 创建测试用图片文件 (1x1 PNG)
export function generateImageFile(): File {
  // 1x1 透明 PNG 的 base64
  const pngBase64 = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==';
  const blob = base64ToBlob(pngBase64, 'image/png');
  return new File([blob], 'test.png', { type: 'image/png' });
}

// 创建测试用视频文件 (很小的占位文件)
export function generateVideoFile(): File {
  // 只是一个很小的占位文件，实际是无效的
  const blob = new Blob(['video placeholder'], { type: 'video/mp4' });
  return new File([blob], 'test.mp4', { type: 'video/mp4' });
}

// base64 转 Blob
function base64ToBlob(base64: string, mimeType: string): Blob {
  const byteCharacters = atob(base64);
  const byteNumbers = new Array(byteCharacters.length);
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);
  return new Blob([byteArray], { type: mimeType });
}

// 类型定义
export interface UserData {
  username: string;
  password: string;
  nickname: string;
  avatar?: string;
  bio?: string;
}

export interface NoteData {
  title: string;
  content: string;
  images?: string[];
  topics?: string[];
}

export interface TestUser {
  userId: number;
  username: string;
  password: string;
  token: string;
}
