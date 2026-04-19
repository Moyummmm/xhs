import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Input, Tabs } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { searchNotes as searchNotesApi } from '@/api/search';
import NoteCard from '@/components/NoteCard';
import Empty from '@/components/Common/Empty';
import type { Note, PaginatedList } from '@/types';

const SearchPage: React.FC = () => {
  const [keyword, setKeyword] = useState('');
  const [activeTab, setActiveTab] = useState('notes');

  const { data: searchData, isLoading } = useQuery<PaginatedList<Note>>({
    queryKey: ['search', activeTab, keyword],
    queryFn: () => searchNotesApi(keyword, 1, 20),
    enabled: !!keyword && activeTab === 'notes',
  });

  const handleSearch = (value: string) => {
    if (value.trim()) {
      setKeyword(value.trim());
    }
  };

  const hotSearches = ['穿搭分享', '美食探店', '旅行攻略', '护肤心得', '健身打卡'];

  return (
    <div className="container-custom py-4">
      {/* Search Input */}
      <div className="max-w-2xl mx-auto mb-6">
        <Input.Search
          placeholder="搜索笔记、用户"
          size="large"
          allowClear
          onSearch={handleSearch}
          prefix={<SearchOutlined className="text-gray-400" />}
        />
      </div>

      {!keyword ? (
        /* Hot Searches */
        <div className="max-w-2xl mx-auto">
          <h3 className="text-sm font-medium text-gray-500 mb-3">热门搜索</h3>
          <div className="flex flex-wrap gap-2">
            {hotSearches.map((term) => (
              <button
                key={term}
                onClick={() => setKeyword(term)}
                className="px-4 py-2 bg-white rounded-full text-sm text-gray-700 hover:bg-primary-50 hover:text-primary-500 transition-colors"
              >
                {term}
              </button>
            ))}
          </div>
        </div>
      ) : (
        /* Search Results */
        <div>
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            items={[
              { key: 'notes', label: '笔记' },
              { key: 'users', label: '用户' },
              { key: 'topics', label: '话题' },
            ]}
          />

          <div className="mt-4">
            {isLoading ? (
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {Array.from({ length: 8 }).map((_, i) => (
                  <div key={i} className="bg-white rounded-xl overflow-hidden animate-pulse">
                    <div className="aspect-[3/4] bg-gray-200" />
                    <div className="p-3 space-y-2">
                      <div className="h-4 bg-gray-200 rounded w-3/4" />
                      <div className="h-3 bg-gray-200 rounded w-1/2" />
                    </div>
                  </div>
                ))}
              </div>
            ) : activeTab === 'notes' ? (
              searchData?.list && searchData.list.length > 0 ? (
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                  {searchData.list.map((note: Note) => (
                    <NoteCard key={note.id} note={note} />
                  ))}
                </div>
              ) : (
                <Empty description="没有找到相关笔记" />
              )
            ) : (
              <Empty description="暂无数据" />
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default SearchPage;
