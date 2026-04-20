import React, { useState } from 'react';
import clsx from 'clsx';

interface ImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  fallback?: string;
  lazy?: boolean;
}

const DEFAULT_FALLBACK = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1 1"%3E%3C/svg%3E';

export const Image: React.FC<ImageProps> = ({
  src,
  alt,
  fallback = DEFAULT_FALLBACK,
  className = '',
  style,
  lazy = true,
  ...props
}) => {
  const [error, setError] = useState(false);
  const [loading, setLoading] = useState(true);

  return (
    <div className={clsx('relative overflow-hidden', className)} style={style}>
      {loading && !error && (
        <div className="absolute inset-0 bg-gray-100 animate-pulse" />
      )}
      <img
        src={error ? fallback : src}
        alt={alt}
        loading={lazy ? 'lazy' : undefined}
        onLoad={() => setLoading(false)}
        onError={() => {
          setError(true);
          setLoading(false);
        }}
        className={clsx(
          'w-full transition-opacity duration-300',
          loading ? 'opacity-0' : 'opacity-100'
        )}
        {...props}
      />
    </div>
  );
};

export default Image;
