import { formatDistanceToNow, format } from 'date-fns';
import { vi } from 'date-fns/locale';

export const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
};

export const formatDate = (date: string | Date, dateFormat = 'dd/MM/yyyy HH:mm'): string => {
  try {
    // Backend sends time with 'Z' suffix but it's actually Vietnam time (not UTC)
    // Remove 'Z' to parse as local time
    let parsedDate: Date;
    if (typeof date === 'string') {
      const dateStr = date.replace('Z', '');
      parsedDate = new Date(dateStr);
    } else {
      parsedDate = date;
    }
    return format(parsedDate, dateFormat, { locale: vi });
  } catch {
    return '';
  }
};

export const formatTimeRemaining = (endDate: string | Date): string => {
  try {
    // Backend sends time like "2024-01-17T10:00:00Z" but it's actually Vietnam time (not UTC)
    // We need to parse it as local time by removing the 'Z' suffix
    let parsedDate: Date;
    if (typeof endDate === 'string') {
      // Remove 'Z' suffix if present and parse as local time
      const dateStr = endDate.replace('Z', '');
      parsedDate = new Date(dateStr);
    } else {
      parsedDate = endDate;
    }
    
    const now = new Date();

    if (parsedDate < now) {
      return 'Đã kết thúc';
    }

    const diff = parsedDate.getTime() - now.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

    if (days > 0) {
      return `${days} ngày ${hours} giờ`;
    } else if (hours > 0) {
      return `${hours} giờ ${minutes} phút`;
    } else {
      return `${minutes} phút`;
    }
  } catch {
    return '';
  }
};

export const formatRelativeTime = (date: string | Date): string => {
  try {
    // Backend sends time with 'Z' suffix but it's actually Vietnam time (not UTC)
    // Remove 'Z' to parse as local time
    let parsedDate: Date;
    if (typeof date === 'string') {
      const dateStr = date.replace('Z', '');
      parsedDate = new Date(dateStr);
    } else {
      parsedDate = date;
    }
    return formatDistanceToNow(parsedDate, { addSuffix: true, locale: vi });
  } catch {
    return '';
  }
};

export const formatBidderName = (name: string): string => {
  if (name.length <= 4) return '*'.repeat(name.length);
  return '*'.repeat(name.length - 2) + name.slice(-2);
};

export const formatNumber = (num: number): string => {
  return new Intl.NumberFormat('vi-VN').format(num);
};

export const truncateText = (text: string, length: number): string => {
  if (text.length <= length) return text;
  return text.substring(0, length) + '...';
};
