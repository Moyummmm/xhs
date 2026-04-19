import React from 'react';
import { LoginModal } from '@/components/Modal';
import type { AuthMode } from '@/components/Modal';

interface LoginPageProps {
  mode?: AuthMode;
}

const LoginPage: React.FC<LoginPageProps> = ({ mode = 'login' }) => {
  return (
    <div className="min-h-[calc(100vh-120px)] flex items-center justify-center">
      <div className="w-full max-w-md bg-white rounded-2xl shadow-card p-8">
        <h1 className="text-3xl font-bold text-center mb-8">
          {mode === 'login' ? '欢迎使用小红书' : '欢迎加入小红书'}
        </h1>
        <LoginModal open={true} mode={mode} onClose={() => window.history.back()} />
      </div>
    </div>
  );
};

export default LoginPage;