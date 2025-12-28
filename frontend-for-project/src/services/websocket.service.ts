import { apiClient } from './api/client';
import { endpoints } from './api/endpoints';

export interface CommentMessage {
  id?: number;
  product_id: number;
  sender_id: number;
  content: string;
  created_at: string;
}

export interface OrderMessage {
  id?: number;
  order_id: number;
  sender_id: number;
  message: string;
  created_at: string;
}

export interface WebSocketMessage {
  type: 'comment' | 'message' | 'typing';
  data?: CommentMessage | OrderMessage;
  product_id?: number;
  content?: string;
}

type MessageHandler = (message: WebSocketMessage) => void;
type ErrorHandler = (error: Event) => void;
type ConnectionHandler = () => void;

class WebSocketService {
  private ws: WebSocket | null = null;
  private messageHandlers: Set<MessageHandler> = new Set();
  private errorHandlers: Set<ErrorHandler> = new Set();
  private openHandlers: Set<ConnectionHandler> = new Set();
  private closeHandlers: Set<ConnectionHandler> = new Set();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectTimeout: NodeJS.Timeout | null = null;

  // Connect to comment service WebSocket
  async connectToComments(productId: number): Promise<void> {
    try {
      const response = await apiClient.get(endpoints.comments.websocket);
      const wsInfo = response.data;

      const wsUrl = wsInfo.comment_service_websocket_url;
      const internalJWT = wsInfo.internal_jwt;
      const userToken = localStorage.getItem('accessToken');

      if (!wsUrl || !internalJWT || !userToken) {
        throw new Error('Missing WebSocket connection information');
      }

      const wsFullUrl = `${wsUrl}?productId=${productId}&X-User-Token=${encodeURIComponent(
        userToken
      )}&X-Internal-JWT=${encodeURIComponent(internalJWT)}`;

      await this.connect(wsFullUrl);
    } catch (error) {
      console.error('Failed to connect to comment WebSocket:', error);
      throw error;
    }
  }

  // Connect to order service WebSocket
  async connectToOrder(orderId: number): Promise<void> {
    try {
      const response = await apiClient.get(endpoints.orders.websocket);
      const wsInfo = response.data;

      const wsUrl = wsInfo.order_service_websocket_url;
      const internalJWT = wsInfo.internal_jwt;
      const userToken = localStorage.getItem('accessToken');

      if (!wsUrl || !internalJWT || !userToken) {
        throw new Error('Missing WebSocket connection information');
      }

      const wsFullUrl = `${wsUrl}?orderId=${orderId}&X-User-Token=${encodeURIComponent(
        userToken
      )}&X-Internal-JWT=${encodeURIComponent(internalJWT)}`;

      await this.connect(wsFullUrl);
    } catch (error) {
      console.error('Failed to connect to order WebSocket:', error);
      throw error;
    }
  }

  // Generic connect method
  private connect(url: string): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          this.openHandlers.forEach((handler) => handler());
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            this.messageHandlers.forEach((handler) => handler(message));
          } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.errorHandlers.forEach((handler) => handler(error));
          reject(error);
        };

        this.ws.onclose = () => {
          console.log('WebSocket disconnected');
          this.closeHandlers.forEach((handler) => handler());
          this.handleReconnect(url);
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  // Handle reconnection logic
  private handleReconnect(url: string): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
      
      console.log(`Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})`);
      
      this.reconnectTimeout = setTimeout(() => {
        console.log('Reconnecting...');
        this.connect(url).catch((error) => {
          console.error('Reconnection failed:', error);
        });
      }, delay);
    } else {
      console.error('Max reconnection attempts reached');
    }
  }

  // Send a message
  send(message: WebSocketMessage): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error('WebSocket is not connected');
    }
  }

  // Send a comment
  sendComment(productId: number, content: string): void {
    this.send({
      type: 'comment',
      product_id: productId,
      content: content,
    });
  }

  // Send an order message
  sendOrderMessage(content: string): void {
    this.send({
      type: 'message',
      content: content,
    });
  }

  // Send typing indicator
  sendTyping(productId?: number): void {
    const message: WebSocketMessage = {
      type: 'typing',
    };
    if (productId) {
      message.product_id = productId;
    }
    this.send(message);
  }

  // Register message handler
  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  // Register error handler
  onError(handler: ErrorHandler): () => void {
    this.errorHandlers.add(handler);
    return () => this.errorHandlers.delete(handler);
  }

  // Register connection open handler
  onOpen(handler: ConnectionHandler): () => void {
    this.openHandlers.add(handler);
    return () => this.openHandlers.delete(handler);
  }

  // Register connection close handler
  onClose(handler: ConnectionHandler): () => void {
    this.closeHandlers.add(handler);
    return () => this.closeHandlers.delete(handler);
  }

  // Disconnect
  disconnect(): void {
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.messageHandlers.clear();
    this.errorHandlers.clear();
    this.openHandlers.clear();
    this.closeHandlers.clear();
    this.reconnectAttempts = 0;
  }

  // Check if connected
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}

// Export singleton instance
export const commentWebSocket = new WebSocketService();
export const orderWebSocket = new WebSocketService();

// Export class for creating new instances if needed
export default WebSocketService;
