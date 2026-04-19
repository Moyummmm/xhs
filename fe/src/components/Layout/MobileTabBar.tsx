import React from 'react';
import { HomeOutlined, CompassOutlined, PlusOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons';
import { Link, useLocation } from 'react-router-dom';
import { useAuthStore, useUIStore } from '@/stores';

export const MobileTabBar: React.FC = () => {
  const location = useLocation();
  const { user, isAuthenticated } = useAuthStore();
  const { openLoginModal } = useUIStore();

  const isActive = (path: string) => location.pathname === path;

  const handlePublishClick = (e: React.MouseEvent) => {
    if (!isAuthenticated) {
      e.preventDefault();
      openLoginModal();
    }
  };

  return (
    <nav className="lg:hidden fixed bottom-0 left-0 right-0 z-50 bg-white border-t border-gray-100 h-[60px]">
      <div className="flex items-center justify-around h-full">
        <Link
          to="/"
          className={`flex flex-col items-center gap-1 ${
            isActive('/') ? 'text-primary-500' : 'text-gray-500'
          }`}
        >
          <HomeOutlined className="text-xl" />
          <span className="text-xs">首页</span>
        </Link>

        <Link
          to="/search"
          className={`flex flex-col items-center gap-1 ${
            isActive('/search') ? 'text-primary-500' : 'text-gray-500'
          }`}
        >
          <CompassOutlined className="text-xl" />
          <span className="text-xs">发现</span>
        </Link>

        <Link
          to="/publish"
          onClick={handlePublishClick}
          className="flex flex-col items-center justify-center -mt-2"
        >
          <div className="w-12 h-12 bg-primary-500 rounded-xl flex items-center justify-center shadow-lg">
            <PlusOutlined className="text-2xl text-white" />
          </div>
        </Link>

        <Link
          to="/notifications"
          className={`flex flex-col items-center gap-1 ${
            isActive('/notifications') ? 'text-primary-500' : 'text-gray-500'
          }`}
        >
          <MessageOutlined className="text-xl" />
          <span className="text-xs">消息</span>
        </Link>

        <Link
          to={isAuthenticated ? `/user/${user?.id}` : '#'}
          onClick={(e) => {
            if (!isAuthenticated) {
              e.preventDefault();
              openLoginModal();
            }
          }}
          className={`flex flex-col items-center gap-1 ${
            isActive(`/user/${user?.id}`) ? 'text-primary-500' : 'text-gray-500'
          }`}
        >
          <UserOutlined className="text-xl" />
          <span className="text-xs">我</span>
        </Link>
      </div>
    </nav>
  );
};

export default MobileTabBar;