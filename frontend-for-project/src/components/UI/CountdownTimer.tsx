import { useEffect, useState } from 'react';
import { formatTimeRemaining } from '../../utils/formatters';
import { Clock } from 'lucide-react';

interface CountdownTimerProps {
  endTime: string;
  onExpire?: () => void;
  className?: string;
  showIcon?: boolean;
}

export const CountdownTimer = ({
  endTime,
  onExpire,
  className = '',
  showIcon = true,
}: CountdownTimerProps) => {
  const [timeLeft, setTimeLeft] = useState<string>('');
  const [isExpired, setIsExpired] = useState(false);

  useEffect(() => {
    const updateTimer = () => {
      const now = Date.now();
      const end = new Date(endTime).getTime();
      const diff = end - now;

      if (diff <= 0) {
        setTimeLeft('Đã kết thúc');
        setIsExpired(true);
        if (onExpire) onExpire();
        return;
      }

      setTimeLeft(formatTimeRemaining(endTime));
    };

    updateTimer();
    const interval = setInterval(updateTimer, 1000);

    return () => clearInterval(interval);
  }, [endTime, onExpire]);

  const getColorClass = () => {
    if (isExpired) return 'text-gray-500';
    const now = Date.now();
    const end = new Date(endTime).getTime();
    const diff = end - now;
    const hours = diff / (1000 * 60 * 60);

    if (hours < 1) return 'text-red-600 font-bold animate-pulse';
    if (hours < 24) return 'text-orange-600 font-semibold';
    return 'text-gray-700';
  };

  return (
    <div className={`flex items-center gap-2 ${getColorClass()} ${className}`}>
      {showIcon && <Clock className="w-5 h-5" />}
      <span>{timeLeft}</span>
    </div>
  );
};
