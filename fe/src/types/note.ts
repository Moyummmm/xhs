import { User } from './user';

// 笔记图片
export interface NoteImage {
  id: number;
  url: string;
  width: number;
  height: number;
  sort_order: number;
}

// 笔记标签
export interface NoteTag {
  id: number;
  name: string;
}

// 话题
export interface Topic {
  id: number;
  name: string;
  description?: string;
  cover_url?: string;
  note_count?: number;
}

// 笔记类型
export type NoteType = 'image' | 'video';

// 笔记状态
export type NoteStatus = 0 | 1; // 0: 草稿, 1: 已发布

// 笔记信息
export interface Note {
  id: number;
  user_id: number;
  title: string;
  content?: string;
  cover_url?: string;
  video_url?: string;
  location?: string;
  topic_id?: number;
  status: NoteStatus;
  like_count: number;
  collect_count: number;
  comment_count: number;
  is_liked?: boolean;
  is_collected?: boolean;
  images?: NoteImage[];
  tags?: NoteTag[];
  topic?: Topic;
  user?: User;
  created_at?: string;
  updated_at?: string;
}

// 创建笔记请求
export interface CreateNoteRequest {
  title: string;
  content?: string;
  image_ids?: number[];
  video_url?: string;
  location?: string;
  topic_id?: number;
  tag_ids?: number[];
}

// 更新笔记请求
export interface UpdateNoteRequest extends Partial<CreateNoteRequest> {}

// 笔记列表查询参数
export interface NoteListParams {
  page?: number;
  page_size?: number;
  type?: 'recommend' | 'follow' | 'latest';
  user_id?: number;
  keyword?: string;
  topic_id?: number;
}
