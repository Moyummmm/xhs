import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';

import Layout from '@/components/Layout';
import ProtectedRoute from '@/components/ProtectedRoute';
import { LoginModal } from '@/components/Modal';

import HomePage from '@/pages/Home';
import NoteDetailPage from '@/pages/NoteDetail';
import PublishPage from '@/pages/Publish';
import ProfilePage from '@/pages/Profile';
import SearchPage from '@/pages/Search';
import NotificationsPage from '@/pages/Notifications';
import LoginPage from '@/pages/Login';
import RegisterPage from '@/pages/Register';
import EditProfilePage from '@/pages/EditProfile';
import NotFoundPage from '@/pages/NotFound';

import { useUIStore, useAuthStore } from '@/stores';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: 1,
    },
  },
});

const AppContent: React.FC = () => {
  const { isLoginModalOpen, closeLoginModal } = useUIStore();
  const { fetchCurrentUser } = useAuthStore();

  useEffect(() => {
    // Try to restore user session on mount
    fetchCurrentUser();
  }, [fetchCurrentUser]);

  return (
    <>
      <Routes>
        {/* Auth routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* Main routes with layout */}
        <Route element={<Layout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/search" element={<SearchPage />} />

          {/* Protected routes */}
          <Route
            path="/publish"
            element={
              <ProtectedRoute>
                <PublishPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/note/:id/edit"
            element={
              <ProtectedRoute>
                <PublishPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/notifications"
            element={
              <ProtectedRoute>
                <NotificationsPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/edit-profile"
            element={
              <ProtectedRoute>
                <EditProfilePage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/collections"
            element={
              <ProtectedRoute>
                <div className="container-custom py-8">
                  <h1 className="text-2xl font-bold mb-4">收藏夹</h1>
                  <p className="text-gray-500">功能开发中...</p>
                </div>
              </ProtectedRoute>
            }
          />
        </Route>

        {/* Routes without layout */}
        <Route path="/note/:id" element={<NoteDetailPage />} />
        <Route path="/user/:id" element={<ProfilePage />} />

        {/* 404 */}
        <Route path="*" element={<NotFoundPage />} />
      </Routes>

      {/* Global Login Modal */}
      <LoginModal open={isLoginModalOpen} onClose={closeLoginModal} />
    </>
  );
};

const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <ConfigProvider locale={zhCN}>
        <BrowserRouter>
          <AppContent />
        </BrowserRouter>
      </ConfigProvider>
    </QueryClientProvider>
  );
};

export default App;