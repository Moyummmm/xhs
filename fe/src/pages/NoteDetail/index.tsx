import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { message, Button } from 'antd';
import { HeartOutlined, HeartFilled, StarOutlined, StarFilled, ShareAltOutlined, LeftOutlined } from '@ant-design/icons';
import ReactPlayer from 'react-player';
import { getNoteDetail, likeNote, unlikeNote, collectNote, uncollectNote } from '@/api/note';
import { getComments, createComment } from '@/api/comment';
import { Image } from '@/components/Common';
import UserAvatar from '@/components/UserAvatar';
import { CommentList, CommentInput } from '@/components/Comment';
import { formatRelativeTime, formatNumber, parseRouteId } from '@/utils';
import { useAuthStore, useUIStore } from '@/stores';
import type { Note, Comment, PaginatedList } from '@/types';

const NoteDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuthStore();
  const { openLoginModal } = useUIStore();

  const noteId = parseRouteId(id);

  const [commentText, setCommentText] = useState('');
  const [activeImageIndex, setActiveImageIndex] = useState(0);

  // 进入详情页时滚动到顶部
  useEffect(() => {
    window.scrollTo(0, 0);
  }, [id]);

  // 获取笔记详情
  const { data: note, isLoading } = useQuery<Note>({
    queryKey: ['note', noteId],
    queryFn: () => getNoteDetail(noteId!),
    enabled: noteId !== null,
  });

  // 获取评论列表
  const { data: commentsData } = useQuery<PaginatedList<Comment>>({
    queryKey: ['comments', noteId],
    queryFn: () => getComments({ note_id: noteId!, page: 1, page_size: 20, type: 'hot' }),
    enabled: noteId !== null,
  });

  // 点赞
  const likeMutation = useMutation({
    mutationFn: (noteId: number) => likeNote(noteId),
    onMutate: async (noteId) => {
      await queryClient.cancelQueries({ queryKey: ['note', noteId] });
      const previous = queryClient.getQueryData(['note', noteId]);
      queryClient.setQueryData(['note', noteId], (old: Note | undefined) => {
        if (!old) return old;
        return {
          ...old,
          is_liked: true,
          like_count: (old.like_count || 0) + 1,
        };
      });
      return { previous };
    },
    onError: (_err, noteId, context) => {
      queryClient.setQueryData(['note', noteId], context?.previous);
      message.error('点赞失败');
    },
    onSettled: (noteId) => {
      queryClient.invalidateQueries({ queryKey: ['note', noteId] });
    },
  });

  const unlikeMutation = useMutation({
    mutationFn: (noteId: number) => unlikeNote(noteId),
    onMutate: async (noteId) => {
      await queryClient.cancelQueries({ queryKey: ['note', noteId] });
      const previous = queryClient.getQueryData(['note', noteId]);
      queryClient.setQueryData(['note', noteId], (old: Note | undefined) => {
        if (!old) return old;
        return {
          ...old,
          is_liked: false,
          like_count: Math.max((old.like_count || 1) - 1, 0),
        };
      });
      return { previous };
    },
    onError: (_err, noteId, context) => {
      queryClient.setQueryData(['note', noteId], context?.previous);
      message.error('取消点赞失败');
    },
    onSettled: (noteId) => {
      queryClient.invalidateQueries({ queryKey: ['note', noteId] });
    },
  });

  // 收藏
  const collectMutation = useMutation({
    mutationFn: (noteId: number) => collectNote(noteId),
    onMutate: async (noteId) => {
      await queryClient.cancelQueries({ queryKey: ['note', noteId] });
      const previous = queryClient.getQueryData(['note', noteId]);
      queryClient.setQueryData(['note', noteId], (old: Note | undefined) => {
        if (!old) return old;
        return {
          ...old,
          is_collected: true,
          collect_count: (old.collect_count || 0) + 1,
        };
      });
      return { previous };
    },
    onError: (_err, noteId, context) => {
      queryClient.setQueryData(['note', noteId], context?.previous);
      message.error('收藏失败');
    },
    onSettled: (noteId) => {
      queryClient.invalidateQueries({ queryKey: ['note', noteId] });
    },
  });

  const uncollectMutation = useMutation({
    mutationFn: (noteId: number) => uncollectNote(noteId),
    onMutate: async (noteId) => {
      await queryClient.cancelQueries({ queryKey: ['note', noteId] });
      const previous = queryClient.getQueryData(['note', noteId]);
      queryClient.setQueryData(['note', noteId], (old: Note | undefined) => {
        if (!old) return old;
        return {
          ...old,
          is_collected: false,
          collect_count: Math.max((old.collect_count || 1) - 1, 0),
        };
      });
      return { previous };
    },
    onError: (_err, noteId, context) => {
      queryClient.setQueryData(['note', noteId], context?.previous);
      message.error('取消收藏失败');
    },
    onSettled: (noteId) => {
      queryClient.invalidateQueries({ queryKey: ['note', noteId] });
    },
  });

  // 评论
  const commentMutation = useMutation({
    mutationFn: (content: string) =>
      createComment({ note_id: noteId!, content }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', noteId] });
      setCommentText('');
      message.success('评论成功');
    },
  });

  const handleLike = () => {
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }
    if (!noteId) return;
    if (note?.is_liked) {
      unlikeMutation.mutate(noteId);
    } else {
      likeMutation.mutate(noteId);
    }
  };

  const handleCollect = () => {
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }
    if (!noteId) return;
    if (note?.is_collected) {
      uncollectMutation.mutate(noteId);
    } else {
      collectMutation.mutate(noteId);
    }
  };

  const handleComment = () => {
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }
    commentMutation.mutate(commentText);
  };

  if (noteId === null) {
    return (
      <div className="container-custom py-4">
        <div className="text-center py-16">
          <p className="text-gray-500">笔记ID无效</p>
          <Button onClick={() => navigate(-1)} className="mt-4">
            返回
          </Button>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="container-custom py-4">
        <div className="bg-white rounded-xl p-6 animate-pulse">
          <div className="h-6 bg-gray-200 rounded w-1/3 mb-4" />
          <div className="aspect-video bg-gray-200 rounded mb-4" />
          <div className="space-y-2">
            <div className="h-4 bg-gray-200 rounded w-3/4" />
            <div className="h-4 bg-gray-200 rounded w-1/2" />
          </div>
        </div>
      </div>
    );
  }

  if (!note) {
    return (
      <div className="container-custom py-4">
        <div className="text-center py-16">
          <p className="text-gray-500">笔记不存在</p>
          <Button onClick={() => navigate(-1)} className="mt-4">
            返回
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <div className="sticky top-0 z-50 w-full bg-white/95 backdrop-blur border-b border-gray-100 px-4 py-3 flex items-center justify-between shadow-sm">
        <button onClick={() => navigate(-1)} className="flex items-center gap-1 text-gray-600 hover:text-gray-800">
          <LeftOutlined />
          <span>返回</span>
        </button>
        <div className="flex items-center gap-3">
          <button className="p-2 hover:bg-gray-100 rounded-full">
            <ShareAltOutlined className="text-gray-600" />
          </button>
        </div>
      </div>

      <div className="container-custom py-4">
        <div className="flex flex-col lg:flex-row gap-8">
          {/* Main Content - Images/Video */}
          <div className="flex-1">
            {note.video_url ? (
              <div className="rounded-xl overflow-hidden bg-black">
                <ReactPlayer
                  url={note.video_url}
                  controls
                  width="100%"
                  height="auto"
                  style={{ aspectRatio: '16/9' }}
                />
              </div>
            ) : note.images && note.images.length > 0 ? (
              <div>
                <div className="relative rounded-xl overflow-hidden bg-gray-100">
                  <Image
                    src={note.images[activeImageIndex]?.url || note.cover_url}
                    alt={note.title}
                    className="w-full"
                    style={{ maxHeight: '70vh' }}
                  />
                </div>
                {note.images.length > 1 && (
                  <div className="flex gap-2 mt-3 overflow-x-auto pb-2">
                    {note.images.map((img, idx) => (
                      <button
                        key={img.id}
                        onClick={() => setActiveImageIndex(idx)}
                        className={`flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden border-2 ${
                          idx === activeImageIndex ? 'border-primary-500' : 'border-transparent'
                        }`}
                      >
                        <img src={img.url} alt="" className="w-full h-full object-cover" />
                      </button>
                    ))}
                  </div>
                )}
              </div>
            ) : null}

            {/* Note Info */}
            <div className="mt-6">
              <h1 className="text-xl font-bold text-gray-800 mb-3">{note.title}</h1>
              <p className="text-gray-700 whitespace-pre-wrap mb-4">{note.content}</p>

              {note.topic && (
                <div className="mb-3">
                  <span className="text-primary-500 text-sm">#{note.topic.name}</span>
                </div>
              )}

              {note.location && (
                <div className="mb-3 text-sm text-gray-500">
                  &#128205; {note.location}
                </div>
              )}

              <p className="text-xs text-gray-400">
                {formatRelativeTime(note.created_at!)}
              </p>
            </div>
          </div>

          {/* Sidebar - Author & Actions */}
          <div className="lg:w-[320px] space-y-4">
            {/* Author Card */}
            {note.user && (
              <div className="bg-white rounded-xl border border-gray-100 p-4">
                <UserAvatar user={note.user} size="md" />
                <p className="text-sm text-gray-500 mt-2 line-clamp-2">{note.user.bio || '这个人很懒，什么都没写~'}</p>
              </div>
            )}

            {/* Action Buttons */}
            <div className="bg-white rounded-xl border border-gray-100 p-4 space-y-3">
              <button
                onClick={handleLike}
                className={`w-full flex items-center justify-center gap-2 py-3 rounded-lg transition-colors ${
                  note.is_liked
                    ? 'bg-primary-50 text-primary-500'
                    : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                }`}
              >
                {note.is_liked ? <HeartFilled /> : <HeartOutlined />}
                <span>{formatNumber(note.like_count)}</span>
              </button>

              <button
                onClick={handleCollect}
                className={`w-full flex items-center justify-center gap-2 py-3 rounded-lg transition-colors ${
                  note.is_collected
                    ? 'bg-yellow-50 text-yellow-500'
                    : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                }`}
              >
                {note.is_collected ? <StarFilled /> : <StarOutlined />}
                <span>收藏</span>
              </button>
            </div>
          </div>
        </div>

        {/* Comments Section */}
        <div className="mt-8 border-t border-gray-100 pt-6">
          <h3 className="text-lg font-bold mb-4">评论 ({note.comment_count})</h3>
          <CommentInput
            value={commentText}
            onChange={setCommentText}
            onSubmit={handleComment}
            disabled={commentMutation.isPending}
          />
          <div className="mt-4">
            <CommentList
              comments={commentsData?.list || []}
              onLike={(commentId) => console.log('Like comment', commentId)}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default NoteDetailPage;
