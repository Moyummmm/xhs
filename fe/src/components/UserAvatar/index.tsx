import React from 'react';
import { Link } from 'react-router-dom';
import clsx from 'clsx';
import type { User } from '@/types';

interface UserAvatarProps {
  user: Pick<User, 'id' | 'nickname' | 'avatar'>;
  size?: 'sm' | 'md' | 'lg';
  showName?: boolean;
  className?: string;
}

const sizeClasses = {
  sm: 'w-6 h-6',
  md: 'w-9 h-9',
  lg: 'w-16 h-16',
};

const textSizeClasses = {
  sm: 'text-xs',
  md: 'text-sm',
  lg: 'text-base',
};

export const UserAvatar: React.FC<UserAvatarProps> = ({
  user,
  size = 'md',
  showName = true,
  className = '',
}) => {
  return (
    <Link to={`/user/${user.id}`} className={clsx('flex items-center gap-2', className)}>
      <div className={clsx('rounded-full overflow-hidden border border-gray-100 flex-shrink-0', sizeClasses[size])}>
        {user.avatar ? (
          <img src={user.avatar} alt={user.nickname} className="w-full h-full object-cover" />
        ) : (
          <div className="w-full h-full bg-gradient-to-br from-primary-100 to-primary-500 flex items-center justify-center text-white font-medium">
            {user.nickname?.charAt(0)?.toUpperCase() || 'U'}
          </div>
        )}
      </div>
      {showName && (
        <span className={clsx('text-gray-700 truncate', textSizeClasses[size])}>
          {user.nickname}
        </span>
      )}
    </Link>
  );
};

export default UserAvatar;