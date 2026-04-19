import React from 'react';
import clsx from 'clsx';

interface LoadingProps {
  size?: 'sm' | 'md' | 'lg';
  text?: string;
  className?: string;
}

const sizeClasses = {
  sm: 'w-4 h-4 border-2',
  md: 'w-8 h-8 border-3',
  lg: 'w-12 h-12 border-4',
};

export const Loading: React.FC<LoadingProps> = ({ size = 'md', text, className = '' }) => {
  return (
    <div className={clsx('flex flex-col items-center justify-center', className)}>
      <div
        className={clsx(
          'rounded-full border-primary-500 border-t-transparent animate-spin',
          sizeClasses[size]
        )}
      />
      {text && <p className="mt-2 text-gray-500 text-sm">{text}</p>}
    </div>
  );
};

export default Loading;
