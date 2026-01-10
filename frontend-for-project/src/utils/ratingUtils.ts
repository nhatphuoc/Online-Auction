import { UserProfile } from '../types';

/**
 * Calculate rating percentage safely
 * Handles division by zero and ensures valid percentage
 */
export const calculateRatingPercentage = (
  totalNumberGoodReviews: number = 0,
  totalNumberReviews: number = 0
): number => {
  if (totalNumberReviews === 0) {
    return 0; // No reviews = 0%
  }
  
  const percentage = (totalNumberGoodReviews / totalNumberReviews) * 100;
  return Math.round(percentage * 10) / 10; // Round to 1 decimal place
};

/**
 * Calculate negative reviews count
 */
export const calculateNegativeReviews = (
  totalNumberReviews: number = 0,
  totalNumberGoodReviews: number = 0
): number => {
  return Math.max(0, totalNumberReviews - totalNumberGoodReviews);
};

/**
 * Get rating display data from user profile
 * Safely handles undefined/null values
 */
export const getRatingDisplay = (profile: UserProfile | null) => {
  if (!profile) {
    return {
      totalRatings: 0,
      positiveRatings: 0,
      negativeRatings: 0,
      ratingPercentage: 0,
    };
  }

  const totalRatings = profile.totalNumberReviews || 0;
  const positiveRatings = profile.totalNumberGoodReviews || 0;
  const negativeRatings = calculateNegativeReviews(totalRatings, positiveRatings);
  const ratingPercentage = calculateRatingPercentage(positiveRatings, totalRatings);

  return {
    totalRatings,
    positiveRatings,
    negativeRatings,
    ratingPercentage,
  };
};

/**
 * Get rating status badge color and text
 */
export const getRatingStatus = (ratingPercentage: number) => {
  if (ratingPercentage >= 90) {
    return {
      color: 'bg-green-100 text-green-700 border-green-300',
      text: 'Xuất sắc',
      emoji: '🌟',
    };
  } else if (ratingPercentage >= 75) {
    return {
      color: 'bg-blue-100 text-blue-700 border-blue-300',
      text: 'Tốt',
      emoji: '👍',
    };
  } else if (ratingPercentage >= 50) {
    return {
      color: 'bg-yellow-100 text-yellow-700 border-yellow-300',
      text: 'Trung bình',
      emoji: '👌',
    };
  } else {
    return {
      color: 'bg-red-100 text-red-700 border-red-300',
      text: 'Cần cải thiện',
      emoji: '⚠️',
    };
  }
};
