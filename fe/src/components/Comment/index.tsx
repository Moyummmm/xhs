import React from 'react';
import { MessageOutlined } from '@ant-design/icons';
import type { Comment as CommentType } from '@/types';
import { formatRelativeTime, formatNumber } from '@/utils';
import UserAvatar from '@/components/UserAvatar';

interface CommentItemProps {
  comment: CommentType;
  onReply?: (comment: CommentType) => void;
  onLike?: (id: number) => void;
}

export const CommentItem: React.FC<CommentItemProps> = ({ comment, onReply, onLike }) => {
  return (
    <div className="flex gap-3 py-3 border-b border-gray-50 last:border-0">
      <UserAvatar user={comment.user!} size="sm" showName={false} />
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <span className="text-sm font-medium text-gray-700">{comment.user?.nickname}</span>
          <span className="text-xs text-gray-400">{formatRelativeTime(comment.created_at!)}</span>
        </div>
        <p className="text-sm text-gray-800 break-words">{comment.content}</p>
        <div className="flex items-center gap-4 mt-2">
          <button
            onClick={() => onLike?.(comment.id)}
            className={`flex items-center gap-1 text-xs ${
              comment.is_liked ? 'text-primary-500' : 'text-gray-400 hover:text-gray-600'
            }`}
          >
            <span>{comment.is_liked ? '❤️' : '🤍'}</span>
            {comment.like_count > 0 && formatNumber(comment.like_count)}
          </button>
          <button
            onClick={() => onReply?.(comment)}
            className="text-xs text-gray-400 hover:text-gray-600"
          >
            回复
          </button>
        </div>
      </div>
    </div>
  );
};

interface CommentListProps {
  comments: CommentType[];
  loading?: boolean;
  onReply?: (comment: CommentType) => void;
  onLike?: (id: number) => void;
}

export const CommentList: React.FC<CommentListProps> = ({
  comments,
  loading = false,
  onReply,
  onLike,
}) => {
  if (!loading && comments.length === 0) {
    return (
      <div className="py-8 text-center text-gray-400">
        <MessageOutlined className="text-4xl mb-2" />
        <p>暂无评论，快来抢沙发吧</p>
      </div>
    );
  }

  return (
    <div className="space-y-0">
      {comments.map((comment) => (
        <CommentItem
          key={comment.id}
          comment={comment}
          onReply={onReply}
          onLike={onLike}
        />
      ))}
    </div>
  );
};

interface CommentInputProps {
  value: string;
  onChange: (value: string) => void;
  onSubmit: () => void;
  placeholder?: string;
  disabled?: boolean;
}

export const CommentInput: React.FC<CommentInputProps> = ({
  value,
  onChange,
  onSubmit,
  placeholder = '说点什么...',
  disabled = false,
}) => {
  const handleSubmit = () => {
    if (!value.trim()) return;
    onSubmit();
  };

  return (
    <div className="flex items-center gap-2 p-3 bg-white border-t border-gray-100">
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
        placeholder={placeholder}
        disabled={disabled}
        className="flex-1 px-4 py-2 bg-gray-100 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 disabled:opacity-50"
      />
      <button
        onClick={handleSubmit}
        disabled={!value.trim() || disabled}
        className="px-4 py-2 bg-primary-500 text-white text-sm rounded-full disabled:opacity-50 disabled:cursor-not-allowed hover:bg-primary-600 transition-colors"
      >
        发送
      </button>
    </div>
  );
};

export default CommentList;
