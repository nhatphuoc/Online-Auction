import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import { Order, OrderDetail } from '../types';

export const orderService = {
  async createOrder(orderData: {
    auctionId: number;
    winnerId: number;
    sellerId: number;
    finalPrice: number;
  }) {
    const response = await apiClient.post<Order>(
      endpoints.orders.create,
      orderData
    );
    return response.data;
  },

  async getOrderById(id: number) {
    const response = await apiClient.get<OrderDetail>(
      endpoints.orders.detail(id)
    );
    return response.data;
  },

  async getUserOrders(params?: {
    role?: 'buyer' | 'seller';
    status?: Order['status'];
  }) {
    const queryParams = new URLSearchParams();
    if (params?.role) queryParams.append('role', params.role);
    if (params?.status) queryParams.append('status', params.status);

    const url = queryParams.toString()
      ? `${endpoints.orders.list}?${queryParams.toString()}`
      : endpoints.orders.list;

    const response = await apiClient.get<Order[]>(url);
    return response.data;
  },

  async getWebSocketInfo() {
    const response = await apiClient.get<{
      order_service_websocket_url: string;
      internal_jwt: string;
    }>(endpoints.orders.websocket);
    return response.data;
  },

  createWebSocketConnection(
    orderId: number,
    internalJwt: string,
    wsUrl: string
  ) {
    const token = localStorage.getItem('accessToken');
    const url = `${wsUrl}?orderId=${orderId}&X-User-Token=${token}&X-Internal-JWT=${internalJwt}`;
    return new WebSocket(url);
  },
};
