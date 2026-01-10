import { useState, useRef, useEffect, useMemo } from 'react';
import { useOrderChat } from '../../hooks/useOrderChat';
import { OrderChatMessage } from '../../services/orderChat.service';
import { useAuthStore } from '../../stores/auth.store';
import { useUIStore } from '../../stores/ui.store';
import { formatDate } from '../../utils/formatters';
import { Send, Loader, WifiOff, Wifi } from 'lucide-react';

interface OrderChatProps {
  orderId: number;
  buyerId: number;
  sellerId: number;
}

/**
 * OrderChat Component
 * Real-time chat component for order communication between buyer and seller
 * Uses WebSocket for instant messaging
 */
export const OrderChat: React.FC<OrderChatProps> = ({ orderId, buyerId, sellerId }) => {
  const { user } = useAuthStore();
  const addToast = useUIStore((state) => state.addToast);
  const [newMessage, setNewMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const messagesContainerRef = useRef<HTMLDivElement>(null); // Container ref
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const typingTimeoutRef = useRef<NodeJS.Timeout>();

  const {
    messages,
    isConnected,
    isLoading,
    sendMessage,
    sendTyping,
  } = useOrderChat({
    orderId,
    enabled: true,
    onMessage: (message: OrderChatMessage) => {
      console.log('New message received:', message);
      // Scroll to bottom when new message arrives
      scrollToBottom();
    },
    onError: (error) => {
      console.error('WebSocket error:', error);
      addToast('error', 'Kết nối chat bị gián đoạn');
    },
    onConnect: () => {
      console.log('Chat connected');
      addToast('success', 'Đã kết nối chat');
    },
    onDisconnect: () => {
      console.log('Chat disconnected');
      addToast('warning', 'Mất kết nối chat');
    },
  });

  const scrollToBottom = () => {
    // Scroll only the messages container, not the whole page
    if (messagesContainerRef.current) {
      const container = messagesContainerRef.current;
      // Always scroll to bottom when new message arrives
      container.scrollTo({
        top: container.scrollHeight,
        behavior: 'smooth'
      });
    }
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newMessage.trim() || !isConnected) {
      return;
    }

    sendMessage(newMessage.trim());
    setNewMessage('');
    setIsTyping(false);
  };

  const handleTyping = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewMessage(e.target.value);

    // Send typing indicator
    if (!isTyping) {
      setIsTyping(true);
      sendTyping();
    }

    // Clear previous timeout
    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    // Set timeout to stop typing indicator
    typingTimeoutRef.current = setTimeout(() => {
      setIsTyping(false);
    }, 2000);
  };

  const isSender = (message: OrderChatMessage) => {
    return message.sender_id === user?.id;
  };

  const getSenderName = (message: OrderChatMessage) => {
    if (message.sender_name) return message.sender_name;
    if (message.sender_id === buyerId) return 'Người mua';
    if (message.sender_id === sellerId) return 'Người bán';
    return 'Unknown';
  };

  // Deduplicate messages by ID (extra safety layer)
  const uniqueMessages = useMemo(() => {
    if (!messages || !Array.isArray(messages)) {
      return [];
    }
    const seen = new Set<number>();
    return messages.filter(msg => {
      if (!msg.id) return true; // Keep messages without ID (shouldn't happen)
      if (seen.has(msg.id)) {
        console.warn('Duplicate message detected in render:', msg.id);
        return false;
      }
      seen.add(msg.id);
      return true;
    });
  }, [messages]);

  return (
    <div className="flex flex-col h-full bg-white rounded-lg shadow-sm">
      {/* Header */}
      <div className="px-4 py-3 border-b border-gray-200 flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Trò chuyện</h3>
        <div className="flex items-center gap-2">
          {isLoading && (
            <div className="flex items-center gap-2 text-gray-500">
              <Loader className="w-4 h-4 animate-spin" />
              <span className="text-sm">Đang kết nối...</span>
            </div>
          )}
          {!isLoading && isConnected && (
            <div className="flex items-center gap-2 text-green-600">
              <Wifi className="w-4 h-4" />
              <span className="text-sm">Đã kết nối</span>
            </div>
          )}
          {!isLoading && !isConnected && (
            <div className="flex items-center gap-2 text-red-600">
              <WifiOff className="w-4 h-4" />
              <span className="text-sm">Mất kết nối</span>
            </div>
          )}
        </div>
      </div>

      {/* Messages */}
      <div 
        ref={messagesContainerRef}
        className="flex-1 overflow-y-auto p-4 space-y-4"
      >
        {uniqueMessages.length === 0 ? (
          <div className="text-center text-gray-500 py-8">
            <p>Chưa có tin nhắn nào</p>
            <p className="text-sm mt-2">Hãy bắt đầu cuộc trò chuyện!</p>
          </div>
        ) : (
          uniqueMessages.map((message) => (
            <div
              key={message.id || `temp-${message.created_at}`}
              className={`flex ${isSender(message) ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                  isSender(message)
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-900'
                }`}
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className={`text-xs font-semibold ${
                    isSender(message) ? 'text-blue-100' : 'text-gray-600'
                  }`}>
                    {getSenderName(message)}
                  </span>
                  <span className={`text-xs ${
                    isSender(message) ? 'text-blue-100' : 'text-gray-500'
                  }`}>
                    {formatDate(message.created_at)}
                  </span>
                </div>
                <p className="text-sm break-words">{message.message}</p>
              </div>
            </div>
          ))
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSendMessage} className="px-4 py-3 border-t border-gray-200">
        <div className="flex gap-2">
          <input
            type="text"
            value={newMessage}
            onChange={handleTyping}
            placeholder={isConnected ? 'Nhập tin nhắn...' : 'Đang kết nối...'}
            disabled={!isConnected || isLoading}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
          />
          <button
            type="submit"
            disabled={!isConnected || !newMessage.trim() || isLoading}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
          >
            <Send className="w-5 h-5" />
            <span>Gửi</span>
          </button>
        </div>
      </form>
    </div>
  );
};

export default OrderChat;
