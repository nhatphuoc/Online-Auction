import { ThumbsUp, ThumbsDown } from 'lucide-react';

interface RatingProps {
  positiveCount: number;
  negativeCount: number;
  className?: string;
  showPercentage?: boolean;
}

export const Rating = ({
  positiveCount,
  negativeCount,
  className = '',
  showPercentage = true,
}: RatingProps) => {
  const total = positiveCount + negativeCount;
  const percentage = total > 0 ? Math.round((positiveCount / total) * 100) : 0;

  const getColorClass = () => {
    if (percentage >= 80) return 'text-green-600';
    if (percentage >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className={`flex items-center gap-3 ${className}`}>
      <div className="flex items-center gap-2">
        <div className="flex items-center gap-1 text-green-600">
          <ThumbsUp className="w-4 h-4" />
          <span className="font-semibold">{positiveCount}</span>
        </div>
        <span className="text-gray-400">/</span>
        <div className="flex items-center gap-1 text-red-600">
          <ThumbsDown className="w-4 h-4" />
          <span className="font-semibold">{negativeCount}</span>
        </div>
      </div>
      {showPercentage && total > 0 && (
        <div className={`font-bold ${getColorClass()}`}>
          ({percentage}%)
        </div>
      )}
    </div>
  );
};

interface RatingButtonProps {
  type: 'positive' | 'negative';
  onClick: () => void;
  disabled?: boolean;
  className?: string;
}

export const RatingButton = ({
  type,
  onClick,
  disabled = false,
  className = '',
}: RatingButtonProps) => {
  const isPositive = type === 'positive';

  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${
        isPositive
          ? 'bg-green-100 text-green-700 hover:bg-green-200'
          : 'bg-red-100 text-red-700 hover:bg-red-200'
      } disabled:opacity-50 disabled:cursor-not-allowed ${className}`}
    >
      {isPositive ? (
        <ThumbsUp className="w-5 h-5" />
      ) : (
        <ThumbsDown className="w-5 h-5" />
      )}
      <span>{isPositive ? 'Tích cực' : 'Tiêu cực'}</span>
    </button>
  );
};
