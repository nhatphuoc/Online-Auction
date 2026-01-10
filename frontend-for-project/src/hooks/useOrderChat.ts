import { useState, useEffect, useCallback, useRef } from 'react';
import orderChatService, { OrderChatMessage, OrderWebSocketMessage } from '../services/orderChat.service';

interface UseOrderChatOptions {
  orderId: number;
  enabled?: boolean;
  onMessage?: (message: OrderChatMessage) => void;
  onError?: (error: Event) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
}

export const useOrderChat = ({
  orderId,
  enabled = true,
  onMessage,
  onError,
  onConnect,
  onDisconnect,
}: UseOrderChatOptions) => {
  const [messages, setMessages] = useState<OrderChatMessage[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const hasLoadedHistory = useRef(false);

  // Reset history flag when orderId changes
  useEffect(() => {
    hasLoadedHistory.current = false;
    setMessages([]);
  }, [orderId]);

  // Load chat history using REST API
  const loadChatHistory = useCallback(async () => {
    if (hasLoadedHistory.current) {
      console.log('Chat history already loaded, skipping...');
      return;
    }
    
    console.log('Loading chat history for order:', orderId);
    try {
      setIsLoading(true);
      const response = await orderChatService.getChatHistory(orderId, 50, 0);
      console.log('Chat history loaded:', response);
      
      // Ensure messages is always an array
      if (response && response.data && Array.isArray(response.data)) {
        setMessages(response.data);
      } else {
        console.warn('Invalid chat history response, using empty array');
        setMessages([]);
      }
      
      hasLoadedHistory.current = true;
    } catch (error) {
      console.error('Failed to load chat history:', error);
      setMessages([]); // Set to empty array on error
    } finally {
      setIsLoading(false);
    }
  }, [orderId]);

  // Connect to WebSocket for real-time messaging
  const connect = useCallback(async () => {
    if (!enabled || isConnected) return;

    try {
      setIsLoading(true);
      const wsInfo = await orderChatService.getWebSocketInfo();
      
      const ws = orderChatService.createWebSocketConnection(
        orderId,
        wsInfo.internal_jwt,
        wsInfo.order_service_websocket_url
      );

      ws.onopen = () => {
        console.log('Order chat WebSocket connected');
        setIsConnected(true);
        setIsLoading(false);
        onConnect?.();
      };

      ws.onmessage = (event) => {
        try {
          const data: OrderWebSocketMessage = JSON.parse(event.data);
          
          if (data.type === 'message' && data.data) {
            // New message received via WebSocket
            const message = data.data as OrderChatMessage;
            
            // Deduplicate messages by ID
            setMessages(prev => {
              // Ensure prev is an array
              const prevMessages = Array.isArray(prev) ? prev : [];
              
              // Skip if message already exists (prevent duplicate)
              if (message.id && prevMessages.some(m => m.id === message.id)) {
                console.log('Duplicate message detected, skipping:', message.id);
                return prevMessages;
              }
              console.log('Adding new message from WebSocket:', message);
              return [...prevMessages, message];
            });
            
            onMessage?.(message);
          } else if (data.type === 'error') {
            console.error('WebSocket error:', data.message);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      ws.onerror = (error) => {
        console.error('Order chat WebSocket error:', error);
        setIsLoading(false);
        onError?.(error);
      };

      ws.onclose = () => {
        console.log('Order chat WebSocket disconnected');
        setIsConnected(false);
        setIsLoading(false);
        onDisconnect?.();
      };

      wsRef.current = ws;
    } catch (error) {
      console.error('Failed to connect to order chat:', error);
      setIsLoading(false);
    }
  }, [orderId, enabled, isConnected, onMessage, onError, onConnect, onDisconnect]);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    if (wsRef.current) {
      orderChatService.closeConnection(wsRef.current);
      wsRef.current = null;
    }
  }, []);

  // Send message
  const sendMessage = useCallback((message: string) => {
    if (wsRef.current && isConnected) {
      orderChatService.sendMessage(wsRef.current, message);
    } else {
      console.error('Cannot send message: WebSocket not connected');
    }
  }, [isConnected]);

  // Send typing indicator
  const sendTyping = useCallback(() => {
    if (wsRef.current && isConnected) {
      orderChatService.sendTyping(wsRef.current);
    }
  }, [isConnected]);

  // Load more messages (pagination)
  const loadMoreMessages = useCallback(async (offset: number) => {
    try {
      const response = await orderChatService.getChatHistory(orderId, 50, offset);
      
      // Ensure response.data is an array before prepending
      if (response && response.data && Array.isArray(response.data)) {
        // Prepend older messages to the beginning of the list
        setMessages(prev => [...response.data, ...prev]);
        return response.pagination;
      } else {
        console.warn('Invalid load more messages response');
        return null;
      }
    } catch (error) {
      console.error('Failed to load more messages:', error);
      return null;
    }
  }, [orderId]);

  // Load chat history and connect to WebSocket
  useEffect(() => {
    if (!enabled) return;

    console.log('useOrderChat effect triggered for orderId:', orderId);
    
    // First load history via REST API
    loadChatHistory();
    
    // Then connect to WebSocket for real-time updates
    connect();

    return () => {
      console.log('useOrderChat cleanup for orderId:', orderId);
      disconnect();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, orderId]); // Only re-run when enabled or orderId changes

  return {
    messages,
    isConnected,
    isLoading,
    sendMessage,
    sendTyping,
    loadMoreMessages,
    connect,
    disconnect,
  };
};

export default useOrderChat;
