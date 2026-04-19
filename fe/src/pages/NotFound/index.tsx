import React from 'react';

const NotFoundPage: React.FC = () => {
  return (
    <div className="min-h-[60vh] flex flex-col items-center justify-center">
      <h1 className="text-6xl font-bold text-gray-300">404</h1>
      <p className="text-xl text-gray-500 mt-4">页面不存在</p>
      <a href="/" className="mt-6 text-primary-500 hover:text-primary-600">
        返回首页
      </a>
    </div>
  );
};

export default NotFoundPage;