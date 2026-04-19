import { User } from './user';

// 评论信息
export interface Comment {
  id: number;
  user_id: number;
  note_id: number;
  parent_id?: number;
  content: string;
  like_count: number;
  is_liked?: boolean;
  user?: User;
  replies?: Comment[];
  reply_count?: number;
  created_at?: string;
}

// 创建评论请求
export interface CreateCommentRequest {
  note_id: number;
  content: string;
  parent_id?: number;
  reply_to_user_id?: number;
}

// 评论列表查询参数
export interface CommentListParams {
  note_id: number;
  page?: number;
  page_size?: number;
  type?: 'hot' | 'latest';
}
