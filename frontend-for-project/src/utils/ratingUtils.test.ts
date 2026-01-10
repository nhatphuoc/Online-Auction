/**
 * Rating Utilities Tests
 * Manual verification for rating calculations
 */

import { calculateRatingPercentage, calculateNegativeReviews, getRatingDisplay } from './ratingUtils';

console.log('=== Rating Utils Tests ===\n');

// Test 1: Division by zero
console.log('Test 1: No reviews (division by zero)');
const result1 = calculateRatingPercentage(0, 0);
console.log(`calculateRatingPercentage(0, 0) = ${result1}%`);
console.log(`Expected: 0%, Result: ${result1 === 0 ? 'PASS ✅' : 'FAIL ❌'}\n`);

// Test 2: All positive reviews
console.log('Test 2: All positive reviews');
const result2 = calculateRatingPercentage(10, 10);
console.log(`calculateRatingPercentage(10, 10) = ${result2}%`);
console.log(`Expected: 100%, Result: ${result2 === 100 ? 'PASS ✅' : 'FAIL ❌'}\n`);

// Test 3: Mixed reviews
console.log('Test 3: Mixed reviews (1 good out of 1 total)');
const result3 = calculateRatingPercentage(1, 1);
console.log(`calculateRatingPercentage(1, 1) = ${result3}%`);
console.log(`Expected: 100%, Result: ${result3 === 100 ? 'PASS ✅' : 'FAIL ❌'}\n`);

// Test 4: Calculate negative reviews
console.log('Test 4: Calculate negative reviews');
const result4 = calculateNegativeReviews(10, 7);
console.log(`calculateNegativeReviews(10, 7) = ${result4}`);
console.log(`Expected: 3, Result: ${result4 === 3 ? 'PASS ✅' : 'FAIL ❌'}\n`);

// Test 5: Edge case - more good than total (shouldn't happen but handle it)
console.log('Test 5: Edge case - safety check');
const result5 = calculateNegativeReviews(5, 10);
console.log(`calculateNegativeReviews(5, 10) = ${result5}`);
console.log(`Expected: 0 (clamped), Result: ${result5 === 0 ? 'PASS ✅' : 'FAIL ❌'}\n`);

// Test 6: getRatingDisplay with null profile
console.log('Test 6: getRatingDisplay with null profile');
const result6 = getRatingDisplay(null);
console.log(`getRatingDisplay(null) =`, result6);
console.log(`Expected: all zeros, Result: ${
  result6.totalRatings === 0 && 
  result6.positiveRatings === 0 && 
  result6.negativeRatings === 0 && 
  result6.ratingPercentage === 0 
    ? 'PASS ✅' : 'FAIL ❌'
}\n`);

// Test 7: Real API response simulation
console.log('Test 7: Real API response from user');
const mockProfile = {
  id: 17,
  email: 'vonhatphuoc32@gmail.com',
  fullName: 'nhatphuoc32',
  phoneNumber: '',
  userRole: 'ROLE_SELLER' as const,
  isEmailVerified: false,
  createdAt: '2026-01-10',
  updatedAt: '2026-01-10',
  totalNumberReviews: 1,
  totalNumberGoodReviews: 1,
};

const result7 = getRatingDisplay(mockProfile);
console.log(`Profile: totalNumberReviews=${mockProfile.totalNumberReviews}, totalNumberGoodReviews=${mockProfile.totalNumberGoodReviews}`);
console.log(`Result:`, result7);
console.log(`Expected: 1 total, 1 positive, 0 negative, 100%`);
console.log(`Result: ${
  result7.totalRatings === 1 && 
  result7.positiveRatings === 1 && 
  result7.negativeRatings === 0 && 
  result7.ratingPercentage === 100 
    ? 'PASS ✅' : 'FAIL ❌'
}\n`);

console.log('=== All Tests Complete ===');
