import React, { useEffect, useRef, useCallback } from 'react';
import Masonry from 'react-masonry-css';
import { NoteCardSkeleton } from '@/components/Common';
import Empty from '@/components/Common/Empty';

const breakpointColumns = {
  default: 4,
  1440: 3,
  1024: 2,
  768: 2,
  640: 1,
};

interface InfiniteScrollProps {
  children: React.ReactNode;
  hasNextPage: boolean;
  isFetchingNextPage: boolean;
  fetchNextPage: () => void;
  isEmpty?: boolean;
  emptyText?: string;
  isLoading?: boolean;
}

export const InfiniteScroll: React.FC<InfiniteScrollProps> = ({
  children,
  hasNextPage,
  isFetchingNextPage,
  fetchNextPage,
  isEmpty = false,
  emptyText = '暂无内容',
  isLoading = false,
}) => {
  const observerRef = useRef<IntersectionObserver | null>(null);
  const loadMoreRef = useRef<HTMLDivElement>(null);

  const handleObserver = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const [entry] = entries;
      if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage();
      }
    },
    [hasNextPage, isFetchingNextPage, fetchNextPage]
  );

  useEffect(() => {
    if (observerRef.current) {
      observerRef.current.disconnect();
    }

    observerRef.current = new IntersectionObserver(handleObserver, {
      rootMargin: '200px',
    });

    if (loadMoreRef.current) {
      observerRef.current.observe(loadMoreRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, []);

  if (isLoading) {
    return (
      <Masonry
        breakpointCols={breakpointColumns}
        className="flex gap-4"
        columnClassName="flex flex-col gap-4"
      >
        {Array.from({ length: 8 }).map((_, i) => (
          <NoteCardSkeleton key={i} />
        ))}
      </Masonry>
    );
  }

  if (isEmpty) {
    return <Empty description={emptyText} />;
  }

  return (
    <>
      <Masonry
        breakpointCols={breakpointColumns}
        className="flex gap-4"
        columnClassName="flex flex-col gap-4"
      >
        {children}
      </Masonry>
      <div ref={loadMoreRef} className="h-10 flex items-center justify-center">
        {isFetchingNextPage && <NoteCardSkeleton />}
      </div>
    </>
  );
};

export default InfiniteScroll;