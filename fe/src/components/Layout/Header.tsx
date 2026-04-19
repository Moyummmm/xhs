import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { SearchOutlined, BellOutlined, PlusOutlined, UserOutlined } from '@ant-design/icons';
import { useAuthStore, useUIStore } from '@/stores';
import { useNotificationStore } from '@/stores/useNotificationStore';

export const Header: React.FC = () => {
  const location = useLocation();
  const { user, isAuthenticated } = useAuthStore();
  const { openLoginModal } = useUIStore();
  const { unreadCount } = useNotificationStore();

  const isActive = (path: string) => location.pathname === path;

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-white shadow-sm h-[60px]">
      <div className="container-custom h-full flex items-center justify-between">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2">
          <span className="text-2xl font-bold text-primary-500">小红书</span>
        </Link>

        {/* Navigation - Desktop */}
        <nav className="hidden md:flex items-center gap-8">
          <Link
            to="/"
            className={`text-base font-medium transition-colors ${
              isActive('/') ? 'text-primary-500' : 'text-gray-700 hover:text-primary-500'
            }`}
          >
            首页
          </Link>
          <Link
            to="/search"
            className={`text-base font-medium transition-colors ${
              isActive('/search') ? 'text-primary-500' : 'text-gray-700 hover:text-primary-500'
            }`}
          >
            发现
          </Link>
        </nav>

        {/* Search Bar */}
        <div className="hidden md:flex flex-1 max-w-md mx-8">
          <Link
            to="/search"
            className="w-full flex items-center bg-gray-100 rounded-full px-4 py-2 hover:bg-gray-200 transition-colors"
          >
            <SearchOutlined className="text-gray-400" />
            <span className="ml-2 text-gray-400 text-sm">搜索笔记、用户</span>
          </Link>
        </div>

        {/* Right Actions */}
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <>
              {/* Notifications */}
              <Link to="/notifications" className="relative p-2 hover:bg-gray-100 rounded-full transition-colors">
                <BellOutlined className="text-xl text-gray-600" />
                {unreadCount > 0 && (
                  <span className="absolute top-0 right-0 w-4 h-4 bg-primary-500 text-white text-xs rounded-full flex items-center justify-center">
                    {unreadCount > 9 ? '9+' : unreadCount}
                  </span>
                )}
              </Link>

              {/* Publish Button */}
              <Link
                to="/publish"
                className="btn-primary flex items-center gap-1"
              >
                <PlusOutlined />
                <span className="hidden sm:inline">发布</span>
              </Link>

              {/* User Avatar */}
              <Link to={`/user/${user?.id}`} className="w-9 h-9 rounded-full overflow-hidden border-2 border-gray-200 hover:border-primary-500 transition-colors">
                {user?.avatar ? (
                  <img src={user.avatar} alt={user.nickname} className="w-full h-full object-cover" />
                ) : (
                  <div className="w-full h-full bg-gray-200 flex items-center justify-center">
                    <UserOutlined className="text-gray-400" />
                  </div>
                )}
              </Link>
            </>
          ) : (
            <button onClick={openLoginModal} className="btn-primary">
              登录
            </button>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;