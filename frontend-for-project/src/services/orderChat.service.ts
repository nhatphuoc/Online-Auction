import apiClient from './api/client';
import { endpoints } from './api/endpoints';

export interface OrderChatMessage {
  id?: number;
  order_id: number;
  sender_id: number;
  message: string;
  created_at: string;
  sender_name?: string;
}

export interface OrderWebSocketInfo {
  order_service_websocket_url: string;
  internal_jwt: string;
}

export interface OrderWebSocketMessage {
  type: 'message' | 'history' | 'typing' | 'error';
  data?: OrderChatMessage | OrderChatMessage[];
  message?: string;
}

export interface OrderChatHistoryResponse {
  data: OrderChatMessage[];
  pagination: {
    total: number;
    limit: number;
    offset: number;
  };
}

/**
 * Order Chat Service using WebSocket and REST API
 * - REST API for loading initial chat history
 * - WebSocket for real-time messaging
 */
export const orderChatService = {
  /**
   * Get chat history using REST API
   * @param orderId - Order ID
   * @param limit - Number of messages to retrieve (default: 50, max: 100)
   * @param offset - Offset for pagination (default: 0)
   * @returns Chat history with pagination
   */
  async getChatHistory(
    orderId: number,
    limit: number = 50,
    offset: number = 0
  ): Promise<OrderChatHistoryResponse> {
    console.log('Fetching chat history:', { orderId, limit, offset });
    console.log('Endpoint:', endpoints.orders.getMessages(orderId));
    
    const response = await apiClient.get<OrderChatHistoryResponse>(
      endpoints.orders.getMessages(orderId),
      {
        params: { limit, offset },
      }
    );
    
    console.log('Chat history response:', response.data);
    return response.data;
  },

  /**
   * Get WebSocket connection information from API Gateway
   */
  async getWebSocketInfo(): Promise<OrderWebSocketInfo> {
    const response = await apiClient.get<OrderWebSocketInfo>(
      endpoints.orders.websocket
    );
    return response.data;
  },

  /**
   * Create WebSocket connection for order chat
   * @param orderId - Order ID to connect to
   * @param internalJwt - Internal JWT token from API Gateway
   * @param wsUrl - WebSocket URL from API Gateway
   * @returns WebSocket connection
   */
  createWebSocketConnection(
    orderId: number,
    internalJwt: string,
    wsUrl: string
  ): WebSocket {
    const token = localStorage.getItem('accessToken');
    
    if (!token) {
      throw new Error('No access token found');
    }

    const url = `${wsUrl}?orderId=${orderId}&X-User-Token=${encodeURIComponent(
      token
    )}&X-Internal-JWT=${encodeURIComponent(internalJwt)}`;

    const ws = new WebSocket(url);

    ws.onerror = (error) => {
      console.error('Order WebSocket error:', error);
    };

    return ws;
  },

  /**
   * Send a message through WebSocket
   * @param ws - WebSocket connection
   * @param message - Message content
   */
  sendMessage(ws: WebSocket, message: string): void {
    if (ws.readyState === WebSocket.OPEN) {
      const payload = {
        type: 'message',
        content: message,
      };
      ws.send(JSON.stringify(payload));
    } else {
      console.error('WebSocket is not open. ReadyState:', ws.readyState);
    }
  },

  /**
   * Send typing indicator through WebSocket
   * @param ws - WebSocket connection
   */
  sendTyping(ws: WebSocket): void {
    if (ws.readyState === WebSocket.OPEN) {
      const payload = {
        type: 'typing',
      };
      ws.send(JSON.stringify(payload));
    }
  },

  /**
   * Close WebSocket connection
   * @param ws - WebSocket connection
   */
  closeConnection(ws: WebSocket): void {
    if (ws && ws.readyState !== WebSocket.CLOSED) {
      ws.close();
    }
  },
};

export default orderChatService;
