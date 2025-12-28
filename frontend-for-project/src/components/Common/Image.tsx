import { useState } from 'react';
import { ImageOff } from 'lucide-react';

interface ImageWithFallbackProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  src?: string;
  alt: string;
  fallbackSrc?: string;
  showPlaceholder?: boolean;
}

export const ImageWithFallback = ({
  src,
  alt,
  fallbackSrc,
  showPlaceholder = true,
  className = '',
  ...props
}: ImageWithFallbackProps) => {
  const [error, setError] = useState(false);
  const [loading, setLoading] = useState(true);

  const handleError = () => {
    setError(true);
    setLoading(false);
  };

  const handleLoad = () => {
    setLoading(false);
  };

  // If no src provided or error occurred and no fallback
  if (!src || (error && !fallbackSrc)) {
    if (!showPlaceholder) return null;
    
    return (
      <div className={`flex items-center justify-center bg-gray-100 ${className}`}>
        <ImageOff className="w-12 h-12 text-gray-300" />
      </div>
    );
  }

  return (
    <>
      {loading && showPlaceholder && (
        <div className={`absolute inset-0 flex items-center justify-center bg-gray-100 animate-pulse ${className}`}>
          <div className="w-12 h-12 rounded-full bg-gray-200 animate-pulse" />
        </div>
      )}
      <img
        {...props}
        src={error && fallbackSrc ? fallbackSrc : src}
        alt={alt}
        className={`${className} ${loading ? 'opacity-0' : 'opacity-100'} transition-opacity duration-300`}
        onError={handleError}
        onLoad={handleLoad}
      />
    </>
  );
};

interface ProductImageProps {
  src?: string;
  alt: string;
  className?: string;
}

export const ProductImage = ({ src, alt, className = '' }: ProductImageProps) => {
  return (
    <ImageWithFallback
      src={src}
      alt={alt}
      className={className}
      showPlaceholder={true}
    />
  );
};

interface AvatarImageProps {
  src?: string;
  alt: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
}

export const AvatarImage = ({ src, alt, size = 'md' }: AvatarImageProps) => {
  const [error, setError] = useState(false);

  const sizeClasses = {
    sm: 'w-8 h-8',
    md: 'w-10 h-10',
    lg: 'w-12 h-12',
    xl: 'w-16 h-16',
  };

  const textSizes = {
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-base',
    xl: 'text-xl',
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((word) => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  if (!src || error) {
    return (
      <div
        className={`${sizeClasses[size]} rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold ${textSizes[size]}`}
      >
        {getInitials(alt)}
      </div>
    );
  }

  return (
    <img
      src={src}
      alt={alt}
      className={`${sizeClasses[size]} rounded-full object-cover border-2 border-gray-200`}
      onError={() => setError(true)}
    />
  );
};
