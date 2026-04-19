# 小红书类应用 - 前端开发详细计划

## 1. 技术栈详述

### 1.1 核心技术
| 组件 | 版本 | 说明 |
|------|------|------|
| React | 18.x | UI 框架 |
| TypeScript | 5.x | 类型系统 |
| Vite | 5.x | 构建工具 |
| React Router | 6.x | 路由管理 |
| Zustand | 4.x | 状态管理 |
| Axios | 1.x | HTTP 客户端 |
| Tailwind CSS | 3.x | 原子化 CSS |

### 1.2 UI 组件库
| 组件 | 说明 |
|------|------|
| Ant Design | 企业级组件库 (仅使用复杂组件如 Modal, Form) |
| @ant-design/icons | 图标库 |
| react-masonry-css | 瀑布流布局 (替代手写实现) |
| react-lazy-load-image-component | 图片懒加载 |

> **样式方案说明**: 
> - Tailwind CSS: 用于布局、间距、颜色等原子化样式
> - Ant Design: 仅用于复杂交互组件 (Modal, DatePicker 等)
> - 自定义 CSS: 仅用于动画和特殊效果
> - **优先级**: Tailwind > 自定义 CSS > Ant Design 默认样式

### 1.3 辅助库
| 库名 | 用途 |
|------|------|
| dayjs | 日期处理 |
| lodash-es | 工具库 (按需引入) |
| ahooks | React Hooks 库 |
| @tanstack/react-query | 数据请求与缓存 (替代手写 hooks) |
| react-player | 视频播放器 |
| clsx + tailwind-merge | 条件类名合并 |

> **说明**: 
> - 使用 Ant Design Form 替代 formily，简化表单开发
> - 使用 textarea + Markdown 预览替代 wangEditor，降低复杂度
> - 使用 @tanstack/react-query 管理 API 请求状态，减少自定义 hooks

---

## 2. 项目结构 (React)

```
fe/
├── public/
│   ├── favicon.ico
│   └── index.html
│
├── src/
│   │
│   ├── api/                        # API 请求层
│   │   ├── request.ts             # Axios 实例配置
│   │   ├── interceptors.ts        # 拦截器
│   │   ├── auth.ts                # 认证相关 API
│   │   ├── user.ts                # 用户相关 API
│   │   ├── note.ts                # 笔记相关 API
│   │   ├── comment.ts             # 评论相关 API
│   │   ├── topic.ts               # 话题相关 API
│   │   ├── search.ts              # 搜索相关 API
│   │   └── notification.ts       # 通知相关 API
│   │
│   ├── components/                 # 公共组件
│   │   ├── Layout/
│   │   │   ├── index.tsx          # 主布局
│   │   │   ├── Header.tsx         # 顶部导航
│   │   │   ├── Sidebar.tsx        # 侧边栏
│   │   │   └── Footer.tsx         # 底部
│   │   │
│   │   ├── NoteCard/              # 笔记卡片
│   │   │   ├── index.tsx
│   │   │   ├── NoteCard.tsx       # 图片笔记卡片
│   │   │   ├── VideoCard.tsx      # 视频笔记卡片
│   │   │   └── styles.ts
│   │   │
│   │   ├── UserAvatar/            # 用户头像
│   │   │   ├── index.tsx
│   │   │   └── UserAvatar.tsx
│   │   │
│   │   ├── ImageUploader/         # 图片上传
│   │   │   ├── index.tsx
│   │   │   ├── ImageUploader.tsx
│   │   │   └── ImagePreview.tsx
│   │   │
│   │   ├── VideoUploader/         # 视频上传
│   │   │   ├── index.tsx
│   │   │   └── VideoUploader.tsx
│   │   │
│   │   ├── Comment/              # 评论组件
│   │   │   ├── index.tsx
│   │   │   ├── CommentList.tsx
│   │   │   ├── CommentItem.tsx
│   │   │   ├── CommentInput.tsx
│   │   │   └── CommentReply.tsx
│   │   │
│   │   ├── InfiniteScroll/        # 无限滚动
│   │   │   ├── index.tsx
│   │   │   └── MasonryGrid.tsx    # 瀑布流
│   │   │
│   │   ├── Modal/                # 弹窗
│   │   │   ├── LoginModal.tsx
│   │   │   ├── ShareModal.tsx
│   │   │   └── ConfirmModal.tsx
│   │   │
│   │   ├── Skeleton/              # 骨架屏
│   │   │   ├── NoteCardSkeleton.tsx
│   │   │   └── UserSkeleton.tsx
│   │   │
│   │   ├── Empty/                 # 空状态
│   │   │   └── Empty.tsx
│   │   │
│   │   └── Common/                # 通用组件
│   │       ├── Button.tsx
│   │       ├── Input.tsx
│   │       ├── Avatar.tsx
│   │       ├── Image.tsx
│   │       ├── Badge.tsx
│   │       ├── Dropdown.tsx
│   │       ├── Tabs.tsx
│   │       ├── Tag.tsx
│   │       ├── Toast.tsx
│   │       └── Loading.tsx
│   │
│   ├── pages/                      # 页面
│   │   ├── Home/                   # 首页
│   │   │   ├── index.tsx
│   │   │   ├── HomePage.tsx
│   │   │   ├── components/
│   │   │   │   ├── FeedTabs.tsx
│   │   │   │   ├── NoteList.tsx
│   │   │   │   └── TopicBanner.tsx
│   │   │   │   └── styles.ts
│   │   │   └── hooks/
│   │   │       └── useFeed.ts
│   │   │
│   │   ├── NoteDetail/             # 笔记详情
│   │   │   ├── index.tsx
│   │   │   ├── NoteDetailPage.tsx
│   │   │   ├── components/
│   │   │   │   ├── NoteHeader.tsx
│   │   │   │   ├── NoteContent.tsx
│   │   │   │   ├── NoteImages.tsx
│   │   │   │   ├── NoteVideo.tsx
│   │   │   │   ├── NoteActions.tsx
│   │   │   │   ├── NoteAuthor.tsx
│   │   │   │   ├── RelatedNotes.tsx
│   │   │   │   └── CommentSection.tsx
│   │   │   └── hooks/
│   │   │       └── useNoteDetail.ts
│   │   │
│   │   ├── Publish/                # 发布笔记
│   │   │   ├── index.tsx
│   │   │   ├── PublishPage.tsx
│   │   │   ├── components/
│   │   │   │   ├── PublishHeader.tsx
│   │   │   │   ├── ImageGrid.tsx
│   │   │   │   ├── VideoUploader.tsx
│   │   │   │   ├── LocationPicker.tsx
│   │   │   │   ├── TopicPicker.tsx
│   │   │   │   └── TagInput.tsx
│   │   │   └── hooks/
│   │   │       └── usePublish.ts
│   │   │
│   │   ├── Profile/               # 用户主页
│   │   │   ├── index.tsx
│   │   │   ├── ProfilePage.tsx
│   │   │   ├── components/
│   │   │   │   ├── ProfileHeader.tsx
│   │   │   │   ├── ProfileTabs.tsx
│   │   │   │   ├── ProfileNotes.tsx
│   │   │   │   ├── ProfileLikes.tsx
│   │   │   │   └── ProfileCollections.tsx
│   │   │   └── hooks/
│   │   │       └── useProfile.ts
│   │   │
│   │   ├── EditProfile/           # 编辑资料
│   │   │   ├── index.tsx
│   │   │   └── EditProfilePage.tsx
│   │   │
│   │   ├── Search/                # 搜索页
│   │   │   ├── index.tsx
│   │   │   ├── SearchPage.tsx
│   │   │   ├── components/
│   │   │   │   ├── SearchHeader.tsx
│   │   │   │   ├── SearchTabs.tsx
│   │   │   │   ├── UserResults.tsx
│   │   │   │   ├── NoteResults.tsx
│   │   │   │   └── TopicResults.tsx
│   │   │   └── hooks/
│   │   │       └── useSearch.ts
│   │   │
│   │   ├── Topic/                 # 话题页
│   │   │   ├── index.tsx
│   │   │   ├── TopicPage.tsx
│   │   │   └── components/
│   │   │       ├── TopicHeader.tsx
│   │   │       └── TopicNotes.tsx
│   │   │
│   │   ├── Notifications/         # 通知页
│   │   │   ├── index.tsx
│   │   │   ├── NotificationsPage.tsx
│   │   │   ├── components/
│   │   │   │   ├── NotificationList.tsx
│   │   │   │   └── NotificationItem.tsx
│   │   │   └── hooks/
│   │   │       └── useNotifications.ts
│   │   │
│   │   ├── Login/                 # 登录页
│   │   │   ├── index.tsx
│   │   │   ├── LoginPage.tsx
│   │   │   └── components/
│   │   │       ├── PhoneLogin.tsx
│   │   │       ├── PasswordLogin.tsx
│   │   │       └── RegisterForm.tsx
│   │   │
│   │   ├── Follow/                # 关注列表
│   │   │   ├── index.tsx
│   │   │   ├── FollowersPage.tsx
│   │   │   └── FollowingsPage.tsx
│   │   │
│   │   ├── Collection/            # 收藏夹
│   │   │   ├── index.tsx
│   │   │   ├── CollectionsPage.tsx
│   │   │   ├── CollectionDetailPage.tsx
│   │   │   └── components/
│   │   │       ├── CollectionList.tsx
│   │   │       └── CollectionItem.tsx
│   │   │
│   │   └── NotFound/              # 404
│   │       └── NotFoundPage.tsx
│   │
│   ├── stores/                     # 状态管理
│   │   ├── useAuthStore.ts        # 认证状态
│   │   ├── useUserStore.ts        # 用户状态
│   │   ├── useNoteStore.ts        # 笔记状态
│   │   ├── useNotificationStore.ts # 通知状态
│   │   └── useUIStore.ts          # UI 状态
│   │
│   ├── hooks/                      # 自定义 Hooks
│   │   ├── useAsync.ts            # 异步请求
│   │   ├── useDebounce.ts         # 防抖
│   │   ├── useThrottle.ts         # 节流
│   │   ├── useInfiniteScroll.ts   # 无限滚动
│   │   ├── useImageUpload.ts      # 图片上传
│   │   ├── useVideoUpload.ts      # 视频上传
│   │   └── useCountUp.ts          # 数字动画
│   │
│   ├── utils/                      # 工具函数
│   │   ├── format.ts              # 格式化 (时间、数字)
│   │   ├── validation.ts          # 表单校验
│   │   ├── storage.ts             # 本地存储
│   │   ├── cookie.ts              # Cookie 操作
│   │   ├── clipboard.ts           # 剪贴板
│   │   ├── share.ts               # 分享功能
│   │   └── browser.ts             # 浏览器检测
│   │
│   ├── types/                      # TypeScript 类型
│   │   ├── api.ts                 # API 响应类型
│   │   ├── user.ts                # 用户类型
│   │   ├── note.ts                # 笔记类型
│   │   ├── comment.ts             # 评论类型
│   │   ├── topic.ts               # 话题类型
│   │   └── notification.ts        # 通知类型
│   │
│   ├── styles/                     # 样式
│   │   ├── global.css             # 全局样式
│   │   ├── variables.css          # CSS 变量
│   │   ├── reset.css              # CSS 重置
│   │   └── animation.css          # 动画
│   │
│   ├── constants/                  # 常量
│   │   ├── api.ts                 # API 路径常量
│   │   ├── config.ts              # 配置常量
│   │   └── reg.ts                 # 正则常量
│   │
│   ├── App.tsx                     # 根组件
│   ├── main.tsx                   # 入口文件
│   └── vite-env.d.ts
│
├── index.html
├── vite.config.ts                 # Vite 配置
├── tsconfig.json                  # TypeScript 配置
├── tsconfig.node.json
├── tailwind.config.js             # Tailwind 配置
├── postcss.config.js              # PostCSS 配置
├── package.json
├── .env.example                   # 环境变量示例
└── Dockerfile
```

---

## 3. 页面路由设计

### 3.1 路由结构
```tsx
// App.tsx
const routes = [
  // 首页
  {
    path: '/',
    component: HomePage,
  },
  // 笔记详情
  {
    path: '/note/:id',
    component: NoteDetailPage,
  },
  // 发布笔记
  {
    path: '/publish',
    component: PublishPage,
    auth: true,
  },
  // 编辑笔记
  {
    path: '/note/:id/edit',
    component: PublishPage,
    auth: true,
  },
  // 用户主页
  {
    path: '/user/:id',
    component: ProfilePage,
  },
  // 编辑资料
  {
    path: '/edit-profile',
    component: EditProfilePage,
    auth: true,
  },
  // 搜索
  {
    path: '/search',
    component: SearchPage,
  },
  // 话题
  {
    path: '/topic/:id',
    component: TopicPage,
  },
  // 通知
  {
    path: '/notifications',
    component: NotificationsPage,
    auth: true,
  },
  // 登录
  {
    path: '/login',
    component: LoginPage,
  },
  // 收藏夹
  {
    path: '/collections',
    component: CollectionsPage,
    auth: true,
  },
  // 收藏夹详情
  {
    path: '/collection/:id',
    component: CollectionDetailPage,
  },
  // 粉丝
  {
    path: '/user/:id/followers',
    component: FollowersPage,
  },
  // 关注
  {
    path: '/user/:id/followings',
    component: FollowingsPage,
  },
  // 404
  {
    path: '/404',
    component: NotFoundPage,
  },
];
```

---

## 4. 组件详细设计

### 4.1 布局组件

#### Header (顶部导航)
```
┌─────────────────────────────────────────────────────────────┐
│  [Logo]  [首页] [发现] [消息]    [搜索框]    [登录] [发布]  │
└─────────────────────────────────────────────────────────────┘
功能:
- Logo 点击跳转首页
- 导航项高亮
- 搜索框快捷搜索
- 未登录显示登录按钮
- 已登录显示头像/消息通知
```

#### 布局结构
```
PC端 (1920px+):
┌────────────────────────────────────────────────────┐
│                    Header (60px)                   │
├──────────┬─────────────────────────────┬────────────┤
│          │                             │            │
│  Sidebar │      Main Content          │  RightBar  │
│  (200px) │      (自适应)              │  (300px)   │
│          │                             │            │
│          │                             │            │
└──────────┴─────────────────────────────┴────────────┘

移动端:
┌─────────────────────┐
│      Header         │
├─────────────────────┤
│                     │
│     Main Content    │
│                     │
├─────────────────────┤
│  底部导航 (TabBar)  │
└─────────────────────┘
```

### 4.2 笔记卡片 (NoteCard)

#### 图文笔记卡片
```
┌───────────────────────┐
│                       │
│    ┌─────────────┐   │
│    │             │   │
│    │   Image 1   │   │
│    │             │   │
│    └─────────────┘   │
│                       │
│   [图片2] [图片3][+3] │  (多图时显示)
│                       │
│   笔记标题...         │
│   笔记内容描述...     │
│                       │
│  [头像] 用户名 · 1小时前 │
│        ❤️ 1.2万 💬 234 │
└───────────────────────┘
```

#### 视频笔记卡片
```
┌───────────────────────┐
│    ┌─────────────┐   │
│    │             │   │
│    │   Video     │   │
│    │   Cover     │   │
│    │    ▶        │   │
│    │   0:30      │   │
│    └─────────────┘   │
│                       │
│   视频笔记标题...     │
│                       │
│  [头像] 用户名 · 1小时前 │
│        ❤️ 1.2万 💬 234 │
└───────────────────────┘
```

### 4.3 笔记详情页

```
┌──────────────────────────────────────────────────────────────┐
│  [← 返回]                           [分享] [举报] [···]     │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  [用户头像] 用户名                            [关注] [私信]  │
│  2024-01-01 12:00  IP属地: 北京                                │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │                                                        │  │
│  │   [图片1]                                              │  │
│  │                                                        │  │
│  │   [图片2]                                              │  │
│  │                                                        │  │
│  │   [图片3]                                              │  │
│  │                                                        │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  笔记标题                                                    │
│  笔记正文内容...                                             │
│                                                              │
│  #标签1 #标签2 #标签3                                        │
│                                                              │
│  📍 北京·朝阳区                                              │
│                                                              │
│  ❤️ 1.2万   💬 234   ⭐ 收藏   📤 分享                       │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│                        评论 (234)                            │
├──────────────────────────────────────────────────────────────┤
│  [评论输入框]                         [表情] [图片] [发送]  │
│                                                              │
│  热门评论                                                    │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ [头像] 用户名 · 1小时前              [回复] ❤️ 123     │  │
│  │ 评论内容...                                             │  │
│  │   └─ [头像] 子评论内容...                             │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  最新评论                                                    │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ ...                                                    │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

### 4.4 发布页

```
┌──────────────────────────────────────────────────────────────┐
│  [取消]              发布笔记                    [发布]      │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐    │
│  │        │ │        │ │   +   │ │        │ │        │    │
│  │ 图片1  │ │ 图片2  │ │ 添加   │ │ 图片3  │ │ 图片4  │    │
│  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘    │
│                                                              │
│  (最多9张图片，支持拖拽排序)                                │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  添加笔记正文...                                        │  │
│  │                                                        │  │
│  │  支持 @ 用户、# 话题标签                                │  │
│  │                                                        │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  ──────────────────────────────────────────────────────────  │
│                                                              │
│  话题标签                                                    │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  # 请选择话题                                          │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  地点                                                    │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  📍 显示/添加位置                                       │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  ──────────────────────────────────────────────────────────  │
│                                                              │
│  @好友                                                    │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  分享给好友或让他们帮你修改                            │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 4.5 用户主页

```
┌──────────────────────────────────────────────────────────────┐
│  [← 返回]                                                    │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│                   [大头像]                                   │
│                                                              │
│                   用户昵称                                   │
│                   @user_id                                  │
│                                                              │
│                   用户简介...                                │
│                                                              │
│     [编辑资料]    [关注 100]    [粉丝 50]    [获赞 1.2万]    │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │                  [关注]                               │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│  [笔记 50]  [收藏 10]  [赞过 100]                           │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                        │
│  │         │ │         │ │         │                        │
│  │  图片1  │ │  图片2  │ │  图片3  │                        │
│  │         │ │         │ │         │                        │
│  └─────────┘ └─────────┘ └─────────┘                        │
│                                                              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐                        │
│  │         │ │         │ │         │                        │
│  │  图片4  │ │  图片5  │ │  图片6  │                        │
│  │         │ │         │ │         │                        │
│  └─────────┘ └─────────┘ └─────────┘                        │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

---

## 5. 核心功能实现

### 5.1 瀑布流布局 (使用 react-masonry-css)

```tsx
// components/InfiniteScroll/MasonryGrid.tsx
import React from 'react';
import Masonry from 'react-masonry-css';
import { useInfiniteQuery } from '@tanstack/react-query';
import NoteCard from '@/components/NoteCard';
import { fetchNotes } from '@/api/note';

const breakpointColumns = {
  default: 4,    // 大屏 4 列
  1440: 3,       // 桌面 3 列
  1024: 2,       // 平板 2 列
  768: 2,        // 小平板 2 列
  640: 1,        // 手机 1 列
};

export const MasonryGrid: React.FC = () => {
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery({
    queryKey: ['notes', 'feed'],
    queryFn: ({ pageParam = 1 }) => fetchNotes({ page: pageParam, pageSize: 20 }),
    getNextPageParam: (lastPage) => {
      return lastPage.has_more ? lastPage.pagination.page + 1 : undefined;
    },
  });

  // 展平所有页的数据
  const notes = data?.pages.flatMap(page => page.list) ?? [];

  return (
    <Masonry
      breakpointCols={breakpointColumns}
      className="flex gap-4"
      columnClassName="flex flex-col gap-4"
    >
      {notes.map((note) => (
        <NoteCard key={note.id} note={note} />
      ))}
    </Masonry>
  );
};
```

> **优势**: 
> - 自动平衡列高度，避免简单取模导致的 uneven columns
> - 响应式断点配置，无需手动计算
> - 更好的性能，内部已优化

### 5.2 图片上传

```tsx
// hooks/useImageUpload.ts
import { useState, useCallback } from 'react';
import { message } from 'antd';
import { uploadImage } from '@/api/upload';

interface UploadedImage {
  url: string;
  width: number;
  height: number;
}

export const useImageUpload = (maxCount = 9) => {
  const [images, setImages] = useState<UploadedImage[]>([]);
  const [uploading, setUploading] = useState(false);

  const upload = useCallback(async (files: FileList) => {
    if (images.length + files.length > maxCount) {
      message.error(`最多上传 ${maxCount} 张图片`);
      return;
    }

    setUploading(true);
    try {
      const uploadPromises = Array.from(files).map(file => uploadImage(file));
      const results = await Promise.all(uploadPromises);
      setImages(prev => [...prev, ...results]);
    } catch (error) {
      message.error('上传失败，请重试');
    } finally {
      setUploading(false);
    }
  }, [images.length, maxCount]);

  const remove = useCallback((index: number) => {
    setImages(prev => prev.filter((_, i) => i !== index));
  }, []);

  const reorder = useCallback((fromIndex: number, toIndex: number) => {
    setImages(prev => {
      const newImages = [...prev];
      const [removed] = newImages.splice(fromIndex, 1);
      newImages.splice(toIndex, 0, removed);
      return newImages;
    });
  }, []);

  return { images, uploading, upload, remove, reorder };
};
```

### 5.3 API 请求管理 (@tanstack/react-query)

```tsx
// api/note.ts - API 定义
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from '@/api/request';

// 获取笔记列表
export const fetchNotes = async ({ page = 1, pageSize = 20 }) => {
  const { data } = await axios.get('/notes/feed', {
    params: { page, page_size: pageSize },
  });
  return data;
};

// 获取笔记详情
export const fetchNoteDetail = async (id: number) => {
  const { data } = await axios.get(`/notes/${id}`);
  return data;
};

// React Query Hooks
export const useNotesFeed = (page: number) => {
  return useQuery({
    queryKey: ['notes', 'feed', page],
    queryFn: () => fetchNotes({ page }),
    staleTime: 5 * 60 * 1000, // 5 分钟内不重新请求
  });
};

export const useNoteDetail = (id: number) => {
  return useQuery({
    queryKey: ['note', id],
    queryFn: () => fetchNoteDetail(id),
    enabled: !!id,
  });
};

// 点赞 mutation
export const useLikeNote = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (noteId: number) => 
      axios.post(`/notes/${noteId}/like`),
    onSuccess: (_, noteId) => {
      // 更新缓存中的点赞状态
      queryClient.invalidateQueries({ queryKey: ['note', noteId] });
      queryClient.invalidateQueries({ queryKey: ['notes', 'feed'] });
    },
  });
};
```

> **优势**:
> - 自动缓存和去重
> - 背景自动刷新
> - 乐观更新支持
> - 减少大量自定义 hooks 代码

### 5.4 认证状态管理 (Zustand)

```tsx
// stores/useAuthStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { getCurrentUser, login, logout } from '@/api/auth';

interface User {
  id: number;
  nickname: string;
  avatar: string;
  // ... other fields
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  
  // Actions
  login: (data: LoginData) => Promise<void>;
  logout: () => Promise<void>;
  fetchCurrentUser: () => Promise<void>;
  setUser: (user: User) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (data) => {
        set({ isLoading: true });
        try {
          const { token, user } = await login(data);
          set({ 
            token, 
            user, 
            isAuthenticated: true,
            isLoading: false 
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: async () => {
        await logout();
        set({ 
          user: null, 
          token: null, 
          isAuthenticated: false 
        });
      },

      fetchCurrentUser: async () => {
        const token = get().token;
        if (!token) return;
        
        try {
          const user = await getCurrentUser();
          set({ user, isAuthenticated: true });
        } catch {
          set({ user: null, token: null, isAuthenticated: false });
        }
      },

      setUser: (user) => set({ user }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        token: state.token 
      }),
    }
  )
);
```

### 5.4 登录弹窗

```tsx
// components/Modal/LoginModal.tsx
import React, { useState } from 'react';
import { Modal, Tabs } from 'antd';
import { useAuthStore } from '@/stores/useAuthStore';
import PhoneLogin from '@/pages/Login/components/PhoneLogin';
import PasswordLogin from '@/pages/Login/components/PasswordLogin';

interface LoginModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export const LoginModal: React.FC<LoginModalProps> = ({
  open,
  onClose,
  onSuccess,
}) => {
  const [activeTab, setActiveTab] = useState('phone');
  const login = useAuthStore(state => state.login);

  const handleLoginSuccess = async (data: LoginData) => {
    await login(data);
    onSuccess?.();
    onClose();
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      width={400}
      centered
      className="login-modal"
    >
      <div className="p-6">
        <h2 className="text-2xl font-bold text-center mb-6">登录</h2>
        
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          centered
          items={[
            {
              key: 'phone',
              label: '手机登录',
              children: <PhoneLogin onSuccess={handleLoginSuccess} />,
            },
            {
              key: 'password',
              label: '密码登录',
              children: <PasswordLogin onSuccess={handleLoginSuccess} />,
            },
          ]}
        />
      </div>
    </Modal>
  );
};
```

---

## 6. 样式设计

### 6.1 Tailwind 配置

```js
// tailwind.config.js
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#fff1f2',
          100: '#ffe4e6',
          500: '#f43f5e',  // 小红书红
          600: '#e11d48',
          700: '#be123c',
        },
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
      },
      borderRadius: {
        'xl': '12px',
        '2xl': '16px',
      },
      boxShadow: {
        'card': '0 2px 8px rgba(0, 0, 0, 0.08)',
        'modal': '0 8px 32px rgba(0, 0, 0, 0.12)',
      },
    },
  },
  plugins: [],
};
```

### 6.2 CSS 变量

```css
/* styles/variables.css */
:root {
  /* 主题色 */
  --color-primary: #f43f5e;
  --color-primary-light: #fff1f2;
  --color-primary-dark: #e11d48;
  
  /* 文字颜色 */
  --color-text-primary: #333333;
  --color-text-secondary: #666666;
  --color-text-tertiary: #999999;
  --color-text-placeholder: #cccccc;
  
  /* 背景色 */
  --color-bg-page: #f5f5f5;
  --color-bg-card: #ffffff;
  --color-bg-mask: rgba(0, 0, 0, 0.5);
  
  /* 边框 */
  --color-border: #eeeeee;
  --color-border-light: #f5f5f5;
  
  /* 间距 */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 12px;
  --spacing-lg: 16px;
  --spacing-xl: 24px;
  --spacing-2xl: 32px;
  
  /* 圆角 */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;
  
  /* 阴影 */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 2px 8px rgba(0, 0, 0, 0.08);
  --shadow-lg: 0 4px 16px rgba(0, 0, 0, 0.1);
  
  /* 动画 */
  --transition-fast: 150ms ease;
  --transition-normal: 300ms ease;
  --transition-slow: 500ms ease;
}
```

---

## 7. 开发任务分解

### Phase 1: 项目搭建 (Week 1-2)

#### Week 1: 环境配置
- [ ] 初始化 Vite + React + TypeScript 项目
- [ ] 配置 Tailwind CSS
- [ ] 配置 ESLint + Prettier
- [ ] 配置路径别名 (@/ → src/)
- [ ] 安装基础依赖
- [ ] 配置 Ant Design
- [ ] 创建目录结构
- [ ] 基础样式文件

#### Week 2: 布局与路由
- [ ] 配置 React Router
- [ ] 创建 Layout 组件
- [ ] Header 组件开发
- [ ] 响应式布局处理
- [ ] 登录注册页面框架
- [ ] 路由守卫

### Phase 2: 认证功能 (Week 3)

#### Week 3: 登录注册
- [ ] API 请求封装 (Axios)
- [ ] 请求/响应拦截器
- [ ] Token 存储与刷新
- [ ] 手机号登录
- [ ] 密码登录
- [ ] 用户注册
- [ ] 登录弹窗组件
- [ ] 登录状态管理
- [ ] 用户信息存储

### Phase 3: 首页与发现 (Week 4-5)

#### Week 4: 首页基础
- [ ] 首页路由与布局
- [ ] Feed Tabs 切换
- [ ] 瀑布流布局组件
- [ ] 笔记卡片组件
- [ ] 无限滚动加载
- [ ] 骨架屏加载
- [ ] 笔记类型区分 (图文/视频)

#### Week 5: 发现功能
- [ ] 推荐算法对接
- [ ] 关注 Tab 开发
- [ ] 最新 Tab 开发
- [ ] 热门话题 Banner
- [ ] 话题跳转
- [ ] 下拉刷新

### Phase 4: 笔记详情 (Week 6)

#### Week 6: 笔记详情页
- [ ] 笔记详情页面
- [ ] 图片预览组件 (支持放大/滑动)
- [ ] 视频播放器组件
- [ ] 笔记内容展示
- [ ] 标签/话题展示
- [ ] 地点展示
- [ ] 作者信息卡片
- [ ] 点赞/收藏/分享
- [ ] 相关笔记推荐

### Phase 5: 社交互动 (Week 7-8)

#### Week 7: 点赞收藏评论
- [ ] 点赞功能 (笔记/评论)
- [ ] 收藏功能
- [ ] 评论列表组件
- [ ] 评论输入组件
- [ ] 回复功能
- [ ] 评论点赞
- [ ] 空状态处理

#### Week 8: 关注与用户
- [ ] 关注/取消关注
- [ ] 用户主页开发
- [ ] 粉丝列表
- [ ] 关注列表
- [ ] 编辑资料页面
- [ ] 头像上传
- [ ] 用户数据展示

### Phase 6: 发布与搜索 (Week 9-10)

#### Week 9: 发布笔记
- [ ] 发布页面开发
- [ ] 图片上传组件 (支持拖拽排序)
- [ ] 图片预览与删除
- [ ] 视频上传组件
- [ ] 富文本编辑器
- [ ] @好友功能
- [ ] 话题选择器
- [ ] 地点选择器
- [ ] 标签输入
- [ ] 发布进度
- [ ] 发布成功跳转

#### Week 10: 搜索功能
- [ ] 搜索页面
- [ ] 搜索历史
- [ ] 热门搜索
- [ ] 搜索建议
- [ ] Tab 切换 (用户/笔记/话题)
- [ ] 搜索结果展示
- [ ] 话题详情页
- [ ] 话题笔记列表

### Phase 7: 通知与收藏 (Week 11)

#### Week 11: 通知收藏
- [ ] 通知列表页面
- [ ] 通知类型区分
- [ ] 未读标记
- [ ] 收藏夹列表
- [ ] 收藏夹详情
- [ ] 创建收藏夹
- [ ] 收藏夹编辑
- [ ] 收藏笔记展示

### Phase 8: 优化与完善 (Week 12)

#### Week 12: 性能优化
- [ ] 图片懒加载
- [ ] 列表虚拟化
- [ ] 路由懒加载
- [ ] 组件缓存
- [ ] 错误边界
- [ ] Loading 状态优化
- [ ] 空状态优化
- [ ] 移动端适配
- [ ] 浏览器兼容
- [ ] 生产构建优化

---

## 8. 环境变量配置

### 8.1 .env.example
```env
# 开发环境
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=XHS
VITE_APP_ENV=development

# 生产环境
# VITE_API_BASE_URL=https://api.example.com/api/v1
# VITE_APP_TITLE=XHS
# VITE_APP_ENV=production

# 功能开关
VITE_ENABLE_MOCK=false
VITE_ENABLE_ANALYTICS=false
```

---

## 9. 常用命令

```bash
# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建生产
pnpm build

# 预览生产
pnpm preview

# 代码检查
pnpm lint

# 代码格式化
pnpm format

# 类型检查
pnpm typecheck

# 单元测试
pnpm test

# Docker 构建
docker build -t xhs-fe:latest .
```

---

## 10. 响应式断点

| 设备 | 宽度 | 布局 |
|------|------|------|
| 移动端 | < 768px | 单列，底部导航 |
| 平板 | 768px - 1024px | 双列，侧边折叠 |
| 桌面 | 1024px - 1440px | 双列 + 右侧栏 |
| 大屏 | > 1440px | 三列布局 |

---

## 11. Git 分支规范

```
main          # 主分支 (生产)
├── develop    # 开发分支
├── feature/  # 功能分支
│   ├── feature/login
│   ├── feature/home
│   └── feature/note-detail
├── bugfix/   # Bug 修复分支
├── hotfix/   # 紧急修复分支
└── release/ # 发布分支
```

### 提交规范
```
feat: 新功能
fix: Bug 修复
docs: 文档更新
style: 代码格式
refactor: 重构
perf: 性能优化
test: 测试
chore: 构建/工具
```

---

## 12. 异常处理与边界情况

### 12.1 网络异常处理

```tsx
// api/request.ts - Axios 拦截器
import axios from 'axios';
import { message } from 'antd';

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 15000,
});

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器
instance.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.code === 'ECONNABORTED') {
      message.error('请求超时，请检查网络连接');
    } else if (!error.response) {
      message.error('网络连接失败，请检查网络');
    } else {
      const { status, data } = error.response;
      switch (status) {
        case 401:
          // Token 过期，跳转登录
          localStorage.removeItem('token');
          window.location.href = '/login';
          break;
        case 403:
          message.error('无权限访问');
          break;
        case 404:
          message.error('资源不存在');
          break;
        case 500:
          message.error('服务器错误，请稍后重试');
          break;
        default:
          message.error(data?.msg || '请求失败');
      }
    }
    return Promise.reject(error);
  }
);

export default instance;
```

### 12.2 图片加载失败处理

```tsx
// components/Common/Image.tsx
import React, { useState } from 'react';

interface ImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  fallback?: string;
}

export const Image: React.FC<ImageProps> = ({ 
  src, 
  alt, 
  fallback = '/images/placeholder.png',
  className = '',
  ...props 
}) => {
  const [error, setError] = useState(false);
  const [loading, setLoading] = useState(true);

  return (
    <div className={`relative ${className}`}>
      {loading && (
        <div className="absolute inset-0 bg-gray-100 animate-pulse" />
      )}
      <img
        src={error ? fallback : src}
        alt={alt}
        onLoad={() => setLoading(false)}
        onError={() => {
          setError(true);
          setLoading(false);
        }}
        className={`${loading ? 'opacity-0' : 'opacity-100'} transition-opacity`}
        {...props}
      />
    </div>
  );
};
```

### 12.3 空状态处理

```tsx
// components/Empty/Empty.tsx
import React from 'react';

interface EmptyProps {
  description?: string;
  action?: React.ReactNode;
}

export const Empty: React.FC<EmptyProps> = ({ 
  description = '暂无数据',
  action 
}) => {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <svg className="w-24 h-24 text-gray-300" fill="none" viewBox="0 0 24 24">
        {/* 空状态图标 */}
      </svg>
      <p className="mt-4 text-gray-500">{description}</p>
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
};
```

### 12.4 表单校验

```tsx
// utils/validation.ts
export const validators = {
  phone: (value: string) => {
    if (!value) return '请输入手机号';
    if (!/^1[3-9]\d{9}$/.test(value)) return '手机号格式不正确';
    return '';
  },
  
  password: (value: string) => {
    if (!value) return '请输入密码';
    if (value.length < 8) return '密码至少 8 位';
    if (!/[a-zA-Z]/.test(value) || !/\d/.test(value)) {
      return '密码必须包含字母和数字';
    }
    return '';
  },
  
  code: (value: string) => {
    if (!value) return '请输入验证码';
    if (!/^\d{6}$/.test(value)) return '验证码为 6 位数字';
    return '';
  },
};
```

### 12.5 大文件上传处理

```tsx
// hooks/useVideoUpload.ts
import { useState } from 'react';
import { message, Progress } from 'antd';
import axios from '@/api/request';

export const useVideoUpload = () => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const upload = async (file: File) => {
    // 文件大小限制 100MB
    if (file.size > 100 * 1024 * 1024) {
      message.error('视频大小不能超过 100MB');
      return null;
    }

    // 文件类型校验
    const allowedTypes = ['video/mp4', 'video/quicktime'];
    if (!allowedTypes.includes(file.type)) {
      message.error('仅支持 MP4 和 MOV 格式');
      return null;
    }

    setUploading(true);
    setProgress(0);

    try {
      const formData = new FormData();
      formData.append('file', file);

      const { data } = await axios.post('/upload/video', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (progressEvent) => {
          const percent = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total!
          );
          setProgress(percent);
        },
      });

      message.success('上传成功');
      return data;
    } catch (error) {
      message.error('上传失败，请重试');
      return null;
    } finally {
      setUploading(false);
      setProgress(0);
    }
  };

  return { upload, uploading, progress };
};
```

---

## 13. 测试策略

### 13.1 测试框架选型
| 测试类型 | 工具 | 说明 |
|----------|------|------|
| 单元测试 | Vitest + Testing Library | 组件和工具函数测试 |
| E2E 测试 | Playwright | 核心用户流程 |
| 视觉回归 | Percy (可选) | UI 一致性检查 |

### 13.2 关键测试用例

```tsx
// __tests__/NoteCard.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { NoteCard } from '@/components/NoteCard';

describe('NoteCard', () => {
  const mockNote = {
    id: 1,
    title: '测试笔记',
    cover_url: 'https://example.com/cover.jpg',
    user: { nickname: '测试用户', avatar: 'https://example.com/avatar.jpg' },
    like_count: 100,
  };

  it('渲染笔记卡片', () => {
    render(<NoteCard note={mockNote} />);
    expect(screen.getByText('测试笔记')).toBeInTheDocument();
    expect(screen.getByText('测试用户')).toBeInTheDocument();
  });

  it('点击卡片跳转到详情页', () => {
    const mockNavigate = vi.fn();
    render(<NoteCard note={mockNote} />, { 
      wrapper: ({ children }) => (
        <MemoryRouter>{children}</MemoryRouter>
      )
    });
    fireEvent.click(screen.getByRole('article'));
    // 验证路由跳转
  });
});
```

### 13.3 E2E 测试示例

```typescript
// e2e/login.spec.ts
import { test, expect } from '@playwright/test';

test('用户登录流程', async ({ page }) => {
  await page.goto('/login');
  
  // 输入手机号
  await page.fill('input[name="phone"]', '13800138000');
  
  // 获取验证码
  await page.click('button:has-text("获取验证码")');
  
  // 输入验证码
  await page.fill('input[name="code"]', '123456');
  
  // 点击登录
  await page.click('button:has-text("登录")');
  
  // 验证跳转首页
  await expect(page).toHaveURL('/');
  await expect(page.locator('[data-testid="user-avatar"]')).toBeVisible();
});
```

### 13.4 CI 集成
```yaml
# .github/workflows/test.yml
name: Frontend Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: pnpm/action-setup@v2
      - run: pnpm install
      - run: pnpm typecheck
      - run: pnpm lint
      - run: pnpm test -- --coverage
      - run: pnpm build
```
