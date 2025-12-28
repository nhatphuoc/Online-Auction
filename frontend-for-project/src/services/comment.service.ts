import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import { Comment } from '../types';

export const commentService = {
  async getProductComments(
    productId: number,
    params: { limit?: number; offset?: number } = {}
  ) {
    const queryParams = new URLSearchParams();
    queryParams.append('limit', (params.limit || 50).toString());
    queryParams.append('offset', (params.offset || 0).toString());

    const response = await apiClient.get<Comment[]>(
      `${endpoints.comments.history(productId)}?${queryParams.toString()}`
    );
    return response.data;
  },

  async getWebSocketInfo() {
    const response = await apiClient.get<{
      comment_service_websocket_url: string;
      internal_jwt: string;
    }>(endpoints.comments.websocket);
    return response.data;
  },

  createWebSocketConnection(
    productId: number,
    internalJwt: string,
    wsUrl: string
  ) {
    const token = localStorage.getItem('accessToken');
    const url = `${wsUrl}?productId=${productId}&X-User-Token=${token}&X-Internal-JWT=${internalJwt}`;
    return new WebSocket(url);
  },
};
