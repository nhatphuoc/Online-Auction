import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import {
  Order,
  OrderDetail,
  CreateOrderRequest,
  PayOrderRequest,
  ShippingAddressRequest,
  ShippingInvoiceRequest,
  CancelOrderRequest,
  // SendMessageRequest, // DEPRECATED: Use WebSocket
  RateOrderRequest,
  // OrderMessage, // DEPRECATED: Use WebSocket
  OrderRating,
  UserRatingStats,
} from '../types';

export const orderService = {
  // Create order (Internal - called by auction service)
  async createOrder(orderData: CreateOrderRequest): Promise<Order> {
    const response = await apiClient.post<Order>(
      endpoints.orders.create,
      orderData
    );
    return response.data;
  },

  // Get order by ID
  async getOrderById(id: number): Promise<OrderDetail> {
    const response = await apiClient.get<OrderDetail>(
      endpoints.orders.detail(id)
    );
    return response.data;
  },

  // Get user orders
  async getUserOrders(params?: {
    role?: 'buyer' | 'seller' | 'ROLE_BIDDER' | 'ROLE_SELLER';
    status?: Order['status'];
    page?: number;
    limit?: number;
  }): Promise<{ data: Order[]; pagination: { page: number; limit: number; total: number; total_pages: number } }> {
    const queryParams = new URLSearchParams();
    if (params?.role) queryParams.append('role', params.role);
    if (params?.status) queryParams.append('status', params.status);
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const url = queryParams.toString()
      ? `${endpoints.orders.list}?${queryParams.toString()}`
      : endpoints.orders.list;

    const response = await apiClient.get<{ data: Order[]; pagination: { page: number; limit: number; total: number; total_pages: number } }>(url);
    return response.data;
  },

  // Pay for order (Buyer)
  async payOrder(id: number, paymentData: PayOrderRequest): Promise<Order> {
    const response = await apiClient.post<Order>(
      endpoints.orders.pay(id),
      paymentData
    );
    return response.data;
  },

  // Provide shipping address (Buyer)
  async provideShippingAddress(
    id: number,
    addressData: ShippingAddressRequest
  ): Promise<Order> {
    const response = await apiClient.post<Order>(
      endpoints.orders.shippingAddress(id),
      addressData
    );
    return response.data;
  },

  // Send shipping invoice (Seller)
  async sendShippingInvoice(
    id: number,
    invoiceData: ShippingInvoiceRequest
  ): Promise<Order> {
    const response = await apiClient.post<Order>(
      endpoints.orders.shippingInvoice(id),
      invoiceData
    );
    return response.data;
  },

  // Confirm delivery (Buyer)
  async confirmDelivery(id: number): Promise<Order> {
    const response = await apiClient.post<Order>(
      endpoints.orders.confirmDelivery(id)
    );
    return response.data;
  },

  // Cancel order (Seller)
  async cancelOrder(id: number, cancelData: CancelOrderRequest): Promise<Order> {
    const response = await apiClient.post<Order>(
      endpoints.orders.cancel(id),
      cancelData
    );
    return response.data;
  },

  // DEPRECATED: Use WebSocket for real-time messaging
  // Send message
  // async sendMessage(id: number, messageData: SendMessageRequest): Promise<OrderMessage> {
  //   const response = await apiClient.post<OrderMessage>(
  //     endpoints.orders.sendMessage(id),
  //     messageData
  //   );
  //   return response.data;
  // },

  // DEPRECATED: Use WebSocket for real-time messaging
  // Get messages
  // async getMessages(id: number, params?: { limit?: number; offset?: number }): Promise<OrderMessage[]> {
  //   const queryParams = new URLSearchParams();
  //   if (params?.limit) queryParams.append('limit', params.limit.toString());
  //   if (params?.offset) queryParams.append('offset', params.offset.toString());

  //   const url = queryParams.toString()
  //     ? `${endpoints.orders.getMessages(id)}?${queryParams.toString()}`
  //     : endpoints.orders.getMessages(id);

  //   const response = await apiClient.get<OrderMessage[]>(url);
  //   return response.data;
  // },

  // Rate order
  async rateOrder(id: number, ratingData: RateOrderRequest): Promise<OrderRating> {
    const response = await apiClient.post<OrderRating>(
      endpoints.orders.rate(id),
      ratingData
    );
    return response.data;
  },

  // Get order rating
  async getOrderRating(id: number): Promise<OrderRating> {
    const response = await apiClient.get<OrderRating>(
      endpoints.orders.getRating(id)
    );
    return response.data;
  },

  // Get user rating statistics
  async getUserRating(userId: number): Promise<UserRatingStats> {
    const response = await apiClient.get<UserRatingStats>(
      endpoints.orders.getUserRating(userId)
    );
    return response.data;
  },

  // Get all orders (Admin only)
  async getAllOrders(params?: {
    status?: Order['status'];
    limit?: number;
    offset?: number;
  }): Promise<Order[]> {
    const queryParams = new URLSearchParams();
    if (params?.status) queryParams.append('status', params.status);
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.offset) queryParams.append('offset', params.offset.toString());

    const url = queryParams.toString()
      ? `${endpoints.orders.adminOrders}?${queryParams.toString()}`
      : endpoints.orders.adminOrders;

    const response = await apiClient.get<Order[]>(url);
    return response.data;
  },

  // Get WebSocket info for order chat
  async getWebSocketInfo(): Promise<{
    order_service_websocket_url: string;
    internal_jwt: string;
  }> {
    const response = await apiClient.get<{
      order_service_websocket_url: string;
      internal_jwt: string;
    }>(endpoints.orders.websocket);
    return response.data;
  },

  // Create WebSocket connection for order chat
  createWebSocketConnection(
    orderId: number,
    internalJwt: string,
    wsUrl: string
  ): WebSocket {
    const token = localStorage.getItem('accessToken');
    const url = `${wsUrl}?orderId=${orderId}&X-User-Token=${token}&X-Internal-JWT=${internalJwt}`;
    return new WebSocket(url);
  },
};

export type { Order, OrderDetail };
