import apiClient from './api/client';
import { ProductListItem, ApiResponse } from '../types';

export interface WatchlistItem {
  id: number;
  userId: number;
  productId: number;
  product: ProductListItem;
  createdAt: string;
}

export const watchlistService = {
  async getWatchlist() {
    const response = await apiClient.get<ApiResponse<WatchlistItem[]>>('/watchlist');
    return response.data.data || [];
  },

  async addToWatchlist(productId: number) {
    const response = await apiClient.post<ApiResponse<WatchlistItem>>(
      `/watchlist/${productId}`
    );
    return response.data;
  },

  async removeFromWatchlist(productId: number) {
    const response = await apiClient.delete<ApiResponse<void>>(
      `/watchlist/${productId}`
    );
    return response.data;
  },

  async isInWatchlist(productId: number) {
    try {
      const watchlist = await this.getWatchlist();
      return watchlist.some(item => item.productId === productId);
    } catch {
      return false;
    }
  },
};
