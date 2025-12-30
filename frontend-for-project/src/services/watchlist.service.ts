import apiClient from './api/client';
import { ApiResponse } from '../types';

export interface WatchlistItem {
  id: number;
  product_id: number;
  thumbnailUrl: string;
  name: string;
  currentPrice: number;
  buyNowPrice: number;
  createdAt: string;
  endAt: string;
  bidCount: number;
  categoryName: string;
}

interface WatchlistResponse {
  message: string;
  data: WatchlistItem[];
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
}

interface CheckWatchlistResponse {
  is_in_watchlist: boolean;
  product_id: number;
}

export const watchlistService = {
  // Get user's watch list - GET /api/orders/data/watchlist
  async getWatchlist(page = 1, limit = 20) {
    const response = await apiClient.get<WatchlistResponse>(
      `/orders/data/watchlist?page=${page}&limit=${limit}`
    );
    return response.data.data || [];
  },

  // Add product to watch list - POST /api/orders/data/watchlist
  async addToWatchlist(productId: number) {
    const response = await apiClient.post<ApiResponse<WatchlistItem>>(
      `/orders/data/watchlist`,
      { product_id: productId }
    );
    return response.data;
  },

  // Remove from watch list - DELETE /api/orders/data/watchlist/{product_id}
  async removeFromWatchlist(productId: number) {
    const response = await apiClient.delete<ApiResponse<void>>(
      `/orders/data/watchlist/${productId}`
    );
    return response.data;
  },

  // Check if product is in watch list - GET /api/orders/data/watchlist/{product_id}/check
  async isInWatchlist(productId: number): Promise<boolean> {
    try {
      const response = await apiClient.get<CheckWatchlistResponse>(
        `/orders/data/watchlist/${productId}/check`
      );
      return response.data.is_in_watchlist;
    } catch {
      return false;
    }
  },
};
