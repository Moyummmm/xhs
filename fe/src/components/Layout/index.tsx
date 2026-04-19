import React from 'react';
import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { MobileTabBar } from './MobileTabBar';

export const Layout: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#f5f5f5]">
      <Header />
      <div className="pt-[60px] pb-[60px] lg:pb-0">
        <Sidebar />
        <main className="lg:ml-[200px] min-h-[calc(100vh-60px)]">
          <Outlet />
        </main>
      </div>
      <MobileTabBar />
    </div>
  );
};

export default Layout;