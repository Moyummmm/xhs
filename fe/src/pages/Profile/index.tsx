import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import { Button, Tabs, message } from 'antd';
import { LeftOutlined, EditOutlined } from '@ant-design/icons';
import { getUserInfo, followUser, unfollowUser, getUserNotes } from '@/api/user';
import NoteCard from '@/components/NoteCard';
import InfiniteScroll from '@/components/InfiniteScroll';
import { formatNumber, parseRouteId } from '@/utils';
import { useAuthStore } from '@/stores';
import type { User, Note, PaginatedList } from '@/types';

const ProfilePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user: currentUser, isAuthenticated } = useAuthStore();
  const userId = parseRouteId(id);
  const isOwnProfile = currentUser?.id === userId;

  // 获取用户信息
  const { data: user } = useQuery<User>({
    queryKey: ['user', userId],
    queryFn: () => getUserInfo(userId!),
    enabled: userId !== null,
  });

  // 获取用户笔记
  const {
    data: notesData,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery<PaginatedList<Note>>({
    queryKey: ['user-notes', userId],
    queryFn: ({ pageParam }) => getUserNotes(userId!, pageParam as number, 20),
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      return lastPage.pagination.has_more ? lastPage.pagination.page + 1 : undefined;
    },
    enabled: userId !== null,
  });

  // 关注/取消关注
  const followMutation = useMutation({
    mutationFn: () => followUser(userId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', userId] });
      message.success('关注成功');
    },
  });

  const unfollowMutation = useMutation({
    mutationFn: () => unfollowUser(userId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', userId] });
      message.success('已取消关注');
    },
  });

  const handleFollow = () => {
    if (!isAuthenticated) {
      message.info('请先登录');
      return;
    }
    if (user?.is_following) {
      unfollowMutation.mutate();
    } else {
      followMutation.mutate();
    }
  };

  if (userId === null) {
    return (
      <div className="min-h-screen bg-[#f5f5f5]">
        <div className="container-custom py-16 text-center">
          <p className="text-gray-500">用户ID无效</p>
          <Button onClick={() => navigate(-1)} className="mt-4">
            返回
          </Button>
        </div>
      </div>
    );
  }

  const notes = notesData?.pages.flatMap((page) => page.list) ?? [];
  const isEmpty = !isLoading && notes.length === 0;

  return (
    <div className="min-h-screen bg-[#f5f5f5]">
      {/* Header */}
      <div className="bg-white border-b border-gray-100 px-4 py-3 flex items-center justify-between sticky top-[60px] z-40">
        <button onClick={() => navigate(-1)} className="flex items-center gap-1 text-gray-600 hover:text-gray-800">
          <LeftOutlined />
        </button>
        {isOwnProfile && (
          <Button icon={<EditOutlined />} onClick={() => navigate('/edit-profile')}>
            编辑资料
          </Button>
        )}
      </div>

      {/* Profile Header */}
      {user && (
        <div className="bg-white">
          <div className="container-custom py-8">
            <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
              <div className="w-24 h-24 rounded-full overflow-hidden border-4 border-white shadow-lg">
                {user.avatar ? (
                  <img src={user.avatar} alt={user.nickname} className="w-full h-full object-cover" />
                ) : (
                  <div className="w-full h-full bg-gradient-to-br from-primary-100 to-primary-500 flex items-center justify-center text-white text-3xl font-bold">
                    {user.nickname?.charAt(0)?.toUpperCase()}
                  </div>
                )}
              </div>

              <div className="flex-1 text-center md:text-left">
                <h1 className="text-2xl font-bold text-gray-800">{user.nickname}</h1>
                <p className="text-gray-500 mt-1">{user.bio || '这个人很懒，什么都没写~'}</p>

                <div className="flex items-center justify-center md:justify-start gap-6 mt-4">
                  <div className="text-center">
                    <p className="text-lg font-bold">{formatNumber(user.note_count || 0)}</p>
                    <p className="text-xs text-gray-500">笔记</p>
                  </div>
                  <div className="text-center">
                    <p className="text-lg font-bold">{formatNumber(user.follower_count || 0)}</p>
                    <p className="text-xs text-gray-500">粉丝</p>
                  </div>
                  <div className="text-center">
                    <p className="text-lg font-bold">{formatNumber(user.follow_count || 0)}</p>
                    <p className="text-xs text-gray-500">关注</p>
                  </div>
                  <div className="text-center">
                    <p className="text-lg font-bold">{formatNumber(user.like_count || 0)}</p>
                    <p className="text-xs text-gray-500">获赞与收藏</p>
                  </div>
                </div>

                {!isOwnProfile && (
                  <div className="mt-4">
                    <Button
                      type={user.is_following ? 'default' : 'primary'}
                      onClick={handleFollow}
                      loading={followMutation.isPending || unfollowMutation.isPending}
                    >
                      {user.is_following ? '已关注' : '关注'}
                    </Button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Notes Grid */}
      <div className="container-custom py-4">
        <Tabs
          defaultActiveKey="notes"
          items={[
            { key: 'notes', label: '笔记' },
            { key: 'likes', label: '赞过' },
            { key: 'collections', label: '收藏' },
          ]}
        />
        <div className="mt-4">
          <InfiniteScroll
            hasNextPage={hasNextPage || false}
            isFetchingNextPage={isFetchingNextPage}
            fetchNextPage={fetchNextPage}
            isEmpty={isEmpty}
            emptyText="暂无笔记"
            isLoading={isLoading}
          >
            {notes.map((note: Note) => (
              <NoteCard key={note.id} note={note} />
            ))}
          </InfiniteScroll>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;
