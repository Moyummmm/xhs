import React from 'react';
import { HomeOutlined, CompassOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons';
import { Link, useLocation } from 'react-router-dom';
import { useAuthStore, useUIStore } from '@/stores';

export const Sidebar: React.FC = () => {
  const location = useLocation();
  const { user } = useAuthStore();
  const { closeSidebar } = useUIStore();

  const isActive = (path: string) => location.pathname === path;

  const menuItems = [
    { path: '/', icon: <HomeOutlined />, label: '首页' },
    { path: '/search', icon: <CompassOutlined />, label: '发现' },
    { path: '/notifications', icon: <MessageOutlined />, label: '消息' },
    { path: `/user/${user?.id}`, icon: <UserOutlined />, label: '我' },
  ];

  return (
    <aside className="hidden lg:block fixed left-0 top-[60px] bottom-0 w-[200px] bg-white border-r border-gray-100 overflow-y-auto">
      <nav className="p-4 space-y-2">
        {menuItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            onClick={closeSidebar}
            className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
              isActive(item.path)
                ? 'bg-primary-50 text-primary-500'
                : 'text-gray-700 hover:bg-gray-50'
            }`}
          >
            <span className="text-xl">{item.icon}</span>
            <span className="font-medium">{item.label}</span>
          </Link>
        ))}
      </nav>
    </aside>
  );
};

export default Sidebar;