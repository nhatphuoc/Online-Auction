import { formatDistanceToNow, format, parseISO } from 'date-fns';
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
    const parsedDate = typeof date === 'string' ? parseISO(date) : date;
    return format(parsedDate, dateFormat, { locale: vi });
  } catch {
    return '';
  }
};

export const formatTimeRemaining = (endDate: string | Date): string => {
  try {
    const parsedDate = typeof endDate === 'string' ? parseISO(endDate) : endDate;
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
    const parsedDate = typeof date === 'string' ? parseISO(date) : date;
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
