export interface User {
  id: number;
  email: string;
  fullName: string;
  phoneNumber: string;
  userRole: 'ROLE_BIDDER' | 'ROLE_SELLER' | 'ROLE_ADMIN';
  isEmailVerified: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface AuthResponse {
  success: boolean;
  accessToken?: string;
  refreshToken?: string;
  message?: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  description?: string;
  parent_id?: number;
  level: 1 | 2;
  is_active: boolean;
  display_order: number;
  children?: Category[];
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  name: string;
  thumbnailUrl: string;
  images: string[];
  description: string;
  parentCategoryId: number;
  parentCategoryName: string;
  categoryId: number;
  categoryName: string;
  startingPrice: number;
  currentPrice: number;
  buyNowPrice?: number;
  stepPrice: number;
  createdAt: string;
  endAt: string;
  autoExtend: boolean;
  extendThresholdMinutes?: number;
  extendDurationMinutes?: number;
  sellerId: number;
  sellerInfo: {
    userId: number;
    username: string;
    avatarUrl?: string;
  };
  highestBidder?: {
    userId: number;
    username: string;
    avatarUrl?: string;
  };
}

export interface ProductListItem {
  id: number;
  thumbnailUrl: string;
  name: string;
  currentPrice: number;
  buyNowPrice?: number;
  createdAt: string;
  endAt: string;
  bidCount: number;
  categoryParentId: number;
  categoryParentName: string;
  categoryId: number;
  categoryName: string;
}

export interface BidRequest {
  productId: number;
  amount: number;
  requestId: string;
}

export interface BidResponse {
  success: boolean;
  message: string;
  data?: {
    newHighest: number;
    previousHighestBidder?: number;
  };
}

export interface BidHistory {
  id: number;
  productId: number;
  bidderId: number;
  amount: number;
  status: 'SUCCESS' | 'FAILED';
  requestId: string;
  createdAt: string;
}

export interface Order {
  id: number;
  auctionId: number;
  winnerId: number;
  sellerId: number;
  finalPrice: number;
  status: 'PENDING_PAYMENT' | 'PAID' | 'SHIPPING' | 'COMPLETED' | 'CANCELLED';
  rating?: {
    id: number;
    orderId: number;
    sellerRating?: number;
    sellerComment?: string;
    buyerRating?: number;
    buyerComment?: string;
    createdAt: string;
    updatedAt: string;
  };
  createdAt: string;
  updatedAt: string;
}

export interface ChatMessage {
  id: number;
  sender_id: number;
  message: string;
  created_at: string;
}

export interface Comment {
  id: number;
  product_id: number;
  sender_id: number;
  content: string;
  created_at: string;
}

export interface PaginationResponse<T> {
  content: T[];
  pageable: {
    pageNumber: number;
    pageSize: number;
  };
  totalElements: number;
  totalPages: number;
}

export interface SearchResponse<T> {
  success: boolean;
  data: {
    content: T[];
    totalElements: number;
    totalPages: number;
    size: number;
    number: number;
    numberOfElements: number;
    first: boolean;
    last: boolean;
    empty: boolean;
  };
  message: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
}
