import React from 'react';
import { Link } from 'react-router-dom';
import clsx from 'clsx';
import { HeartOutlined } from '@ant-design/icons';
import { Image } from '@/components/Common';
import type { Note } from '@/types';
import { formatNumber } from '@/utils';

interface NoteCardProps {
  note: Note;
  className?: string;
}

export const NoteCard: React.FC<NoteCardProps> = ({ note, className }) => {
  const { id, title, cover_url, images, user, like_count } = note;

  const firstImage = images && images.length > 0 ? images[0] : null;

  const cardContent = (
    <>
      {/* Cover Image */}
      <div className="relative w-full">
        <Image
          src={cover_url || firstImage?.url || ''}
          alt={title}
          className="w-full h-auto block"
        />
        {images && images.length > 1 && (
          <div className="absolute bottom-2 right-2 bg-black/50 text-white text-xs px-2 py-1 rounded-full">
            {images.length}
          </div>
        )}
      </div>

      {/* Content */}
      <div className="p-3">
        <h3 className="text-sm font-medium text-gray-800 line-clamp-2 mb-2">{title}</h3>

        {/* Author & Stats */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {user?.avatar ? (
              <img src={user.avatar} alt={user.nickname} className="w-5 h-5 rounded-full object-cover" />
            ) : (
              <div className="w-5 h-5 rounded-full bg-gray-200" />
            )}
            <span className="text-xs text-gray-500 truncate max-w-[80px]">{user?.nickname}</span>
          </div>
          <div className="flex items-center gap-3 text-gray-400 text-xs">
            <span className="flex items-center gap-1">
              <HeartOutlined />
              {formatNumber(like_count || 0)}
            </span>
          </div>
        </div>
      </div>
    </>
  );

  if (!id) {
    return (
      <div className={clsx('block bg-white rounded-xl overflow-hidden shadow-card', className)}>
        {cardContent}
      </div>
    );
  }

  return (
    <Link to={`/note/${id}`} className={clsx('block bg-white rounded-xl overflow-hidden shadow-card hover:shadow-lg transition-shadow', className)}>
      {cardContent}
    </Link>
  );
};

export default NoteCard;
