import React, { useState } from 'react';
import { useInfiniteQuery } from '@tanstack/react-query';
import { getNotesFeed } from '@/api/note';
import NoteCard from '@/components/NoteCard';
import InfiniteScroll from '@/components/InfiniteScroll';
import type { PaginatedList, Note } from '@/types';

const FeedTabs: React.FC<{ activeTab: string; onChange: (tab: string) => void }> = ({
  activeTab,
  onChange,
}) => {
  const tabs = [
    { key: 'recommend', label: '推荐' },
    { key: 'follow', label: '关注' },
    { key: 'latest', label: '最新' },
  ];

  return (
    <div className="flex items-center gap-6 border-b border-gray-100 bg-white">
      {tabs.map((tab) => (
        <button
          key={tab.key}
          onClick={() => onChange(tab.key)}
          className={`py-4 px-2 text-base font-medium transition-colors relative ${
            activeTab === tab.key ? 'text-primary-500' : 'text-gray-600 hover:text-gray-800'
          }`}
        >
          {tab.label}
          {activeTab === tab.key && (
            <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-6 h-0.5 bg-primary-500 rounded-full" />
          )}
        </button>
      ))}
    </div>
  );
};

const HomePage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('recommend');

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery<PaginatedList<Note>>({
    queryKey: ['notes', 'feed', activeTab],
    queryFn: ({ pageParam }) =>
      getNotesFeed({ type: activeTab as any, page: pageParam as number, page_size: 20 }),
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      return lastPage.pagination.has_more ? lastPage.pagination.page + 1 : undefined;
    },
  });

  const notes = data?.pages.flatMap((page) => page.list) ?? [];
  const isEmpty = !isLoading && notes.length === 0;

  return (
    <div className="container-custom py-4">
      <FeedTabs activeTab={activeTab} onChange={setActiveTab} />
      <div className="mt-4">
        <InfiniteScroll
          hasNextPage={hasNextPage || false}
          isFetchingNextPage={isFetchingNextPage}
          fetchNextPage={fetchNextPage}
          isEmpty={isEmpty}
          emptyText="暂无内容"
          isLoading={isLoading}
        >
          {notes.map((note, index) => (
            <NoteCard key={note.id ?? `fallback-${index}`} note={note} />
          ))}
        </InfiniteScroll>
      </div>
    </div>
  );
};

export default HomePage;
