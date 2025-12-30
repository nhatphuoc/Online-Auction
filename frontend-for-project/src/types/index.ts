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
  auction_id: number;
  winner_id: number;
  seller_id: number;
  final_price: number;
  status: 'PENDING_PAYMENT' | 'PAID' | 'ADDRESS_PROVIDED' | 'SHIPPING' | 'DELIVERED' | 'COMPLETED' | 'CANCELLED';
  payment_method?: string;
  payment_proof?: string;
  paid_at?: string;
  shipping_address?: string;
  shipping_phone?: string;
  tracking_number?: string;
  shipping_invoice?: string;
  delivered_at?: string;
  completed_at?: string;
  cancelled_at?: string;
  cancel_reason?: string;
  // User names from backend JOIN
  buyer_name?: string;
  seller_name?: string;
  // Extended fields (may come from joined queries)
  product_name?: string;
  product_image?: string;
  seller_info?: {
    id: number;
    username: string;
    email: string;
  };
  buyer_info?: {
    id: number;
    username: string;
    email: string;
  };
  rating?: OrderRating;
  created_at: string;
  updated_at: string;
}

export interface OrderRating {
  id: number;
  order_id: number;
  buyer_rating?: number | null; // 1 or -1
  buyer_comment?: string;
  buyer_rated_at?: string | null;
  seller_rating?: number | null; // 1 or -1
  seller_comment?: string;
  seller_rated_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface OrderMessage {
  id: number;
  order_id: number;
  sender_id: number;
  message: string;
  created_at: string;
}

export interface CreateOrderRequest {
  auction_id: number;
  winner_id: number;
  seller_id: number;
  final_price: number;
}

export interface PayOrderRequest {
  payment_method: 'MOMO' | 'ZALOPAY' | 'VNPAY' | 'STRIPE' | 'PAYPAL';
  payment_proof?: string;
}

export interface ShippingAddressRequest {
  shipping_address: string;
  shipping_phone: string;
}

export interface ShippingInvoiceRequest {
  tracking_number: string;
  shipping_invoice?: string;
}

export interface CancelOrderRequest {
  cancel_reason: string;
}

export interface SendMessageRequest {
  message: string;
}

export interface RateOrderRequest {
  rating: 1 | -1; // +1 or -1
  comment?: string;
}

export interface UserRatingStats {
  user_id: number;
  total_number_reviews: number;
  total_number_good_reviews: number;
  rating_percentage: number;
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
  sender_name?: string;
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

export interface UserProfile extends User {
  address?: string;
  dateOfBirth?: string;
  rating?: UserRating;
}

export interface UserRating {
  totalRatings: number;
  positiveRatings: number;
  negativeRatings: number;
  ratingPercentage: number;
  reviews: RatingReview[];
}

export interface RatingReview {
  id: number;
  fromUserId: number;
  fromUserName: string;
  rating: 1 | -1; // +1 or -1
  comment: string;
  createdAt: string;
}

export interface UpgradeRequest {
  id: number;
  userId: number;
  userEmail: string;
  userName: string;
  reason: string;
  status: 'PENDING' | 'APPROVED' | 'REJECTED';
  createdAt: string;
  updatedAt: string;
}

export interface WatchlistItem {
  id: number;
  userId: number;
  productId: number;
  product: ProductListItem;
  createdAt: string;
}

export interface Notification {
  id: number;
  userId: number;
  type: 'BID_PLACED' | 'BID_OUTBID' | 'AUCTION_WON' | 'AUCTION_ENDED' | 'QUESTION_RECEIVED' | 'ANSWER_RECEIVED';
  title: string;
  message: string;
  read: boolean;
  link?: string;
  createdAt: string;
}

export interface ProductQuestion {
  id: number;
  productId: number;
  askerId: number;
  askerName: string;
  question: string;
  answer?: string;
  answeredAt?: string;
  createdAt: string;
}

export interface OrderDetail extends Order {
  messages?: OrderMessage[];
}

export interface MediaUploadResponse {
  message: string;
  file_name: string;
  file_size: number;
  mime_type: string;
  image_url: string;
  uploaded_at: string;
}

export interface PresignedUrlResponse {
  presigned_url: string;
  image_url: string;
  file_name: string;
  expires_in: number;
}

export interface WebSocketMessage {
  type: 'BID_UPDATE' | 'COMMENT' | 'ORDER_UPDATE' | 'NOTIFICATION';
  data: Record<string, unknown>;
  timestamp: string;
}
