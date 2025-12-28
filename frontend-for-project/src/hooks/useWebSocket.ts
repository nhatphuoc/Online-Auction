import { useEffect, useRef, useCallback, useState } from 'react';
import WebSocketService, { WebSocketMessage, CommentMessage, OrderMessage } from '../services/websocket.service';

export interface UseWebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void;
  onError?: (error: Event) => void;
  onOpen?: () => void;
  onClose?: () => void;
  autoConnect?: boolean;
}

export interface UseCommentWebSocketReturn {
  isConnected: boolean;
  sendComment: (content: string) => void;
  sendTyping: () => void;
  connect: () => Promise<void>;
  disconnect: () => void;
  messages: CommentMessage[];
}

export interface UseOrderWebSocketReturn {
  isConnected: boolean;
  sendMessage: (content: string) => void;
  sendTyping: () => void;
  connect: () => Promise<void>;
  disconnect: () => void;
  messages: OrderMessage[];
}

// Hook for comment WebSocket
export function useCommentWebSocket(
  productId: number,
  options: UseWebSocketOptions = {}
): UseCommentWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<CommentMessage[]>([]);
  const wsRef = useRef<WebSocketService | null>(null);
  const { onMessage, onError, onOpen, onClose, autoConnect = false } = options;

  const connect = useCallback(async () => {
    if (!wsRef.current) {
      wsRef.current = new WebSocketService();
    }

    try {
      await wsRef.current.connectToComments(productId);
    } catch (error) {
      console.error('Failed to connect to comment WebSocket:', error);
      throw error;
    }
  }, [productId]);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.disconnect();
      setIsConnected(false);
    }
  }, []);

  const sendComment = useCallback(
    (content: string) => {
      if (wsRef.current) {
        wsRef.current.sendComment(productId, content);
      }
    },
    [productId]
  );

  const sendTyping = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.sendTyping(productId);
    }
  }, [productId]);

  useEffect(() => {
    if (!wsRef.current) {
      wsRef.current = new WebSocketService();
    }

    const ws = wsRef.current;

    // Register handlers
    const unsubscribeMessage = ws.onMessage((message: WebSocketMessage) => {
      if (message.type === 'comment' && message.data) {
        setMessages((prev) => [...prev, message.data as CommentMessage]);
      }
      onMessage?.(message);
    });

    const unsubscribeError = ws.onError((error) => {
      setIsConnected(false);
      onError?.(error);
    });

    const unsubscribeOpen = ws.onOpen(() => {
      setIsConnected(true);
      onOpen?.();
    });

    const unsubscribeClose = ws.onClose(() => {
      setIsConnected(false);
      onClose?.();
    });

    // Auto connect if enabled
    if (autoConnect) {
      connect();
    }

    // Cleanup
    return () => {
      unsubscribeMessage();
      unsubscribeError();
      unsubscribeOpen();
      unsubscribeClose();
      disconnect();
    };
  }, [productId, onMessage, onError, onOpen, onClose, autoConnect, connect, disconnect]);

  return {
    isConnected,
    sendComment,
    sendTyping,
    connect,
    disconnect,
    messages,
  };
}

// Hook for order WebSocket
export function useOrderWebSocket(
  orderId: number,
  options: UseWebSocketOptions = {}
): UseOrderWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<OrderMessage[]>([]);
  const wsRef = useRef<WebSocketService | null>(null);
  const { onMessage, onError, onOpen, onClose, autoConnect = false } = options;

  const connect = useCallback(async () => {
    if (!wsRef.current) {
      wsRef.current = new WebSocketService();
    }

    try {
      await wsRef.current.connectToOrder(orderId);
    } catch (error) {
      console.error('Failed to connect to order WebSocket:', error);
      throw error;
    }
  }, [orderId]);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.disconnect();
      setIsConnected(false);
    }
  }, []);

  const sendMessage = useCallback(
    (content: string) => {
      if (wsRef.current) {
        wsRef.current.sendOrderMessage(content);
      }
    },
    []
  );

  const sendTyping = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.sendTyping();
    }
  }, []);

  useEffect(() => {
    if (!wsRef.current) {
      wsRef.current = new WebSocketService();
    }

    const ws = wsRef.current;

    // Register handlers
    const unsubscribeMessage = ws.onMessage((message: WebSocketMessage) => {
      if (message.type === 'message' && message.data) {
        setMessages((prev) => [...prev, message.data as OrderMessage]);
      }
      onMessage?.(message);
    });

    const unsubscribeError = ws.onError((error) => {
      setIsConnected(false);
      onError?.(error);
    });

    const unsubscribeOpen = ws.onOpen(() => {
      setIsConnected(true);
      onOpen?.();
    });

    const unsubscribeClose = ws.onClose(() => {
      setIsConnected(false);
      onClose?.();
    });

    // Auto connect if enabled
    if (autoConnect) {
      connect();
    }

    // Cleanup
    return () => {
      unsubscribeMessage();
      unsubscribeError();
      unsubscribeOpen();
      unsubscribeClose();
      disconnect();
    };
  }, [orderId, onMessage, onError, onOpen, onClose, autoConnect, connect, disconnect]);

  return {
    isConnected,
    sendMessage,
    sendTyping,
    connect,
    disconnect,
    messages,
  };
}
