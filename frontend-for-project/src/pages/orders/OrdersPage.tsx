import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { orderService } from '../../services/order.service';
import { Order } from '../../types';
import { formatCurrency, formatRelativeTime } from '../../utils/formatters';
import { Package, ShoppingCart, Clock, User, Star, MessageSquare } from 'lucide-react';
import { useUIStore } from '../../stores/ui.store';

const OrdersPage = () => {
  const navigate = useNavigate();
  const addToast = useUIStore((state) => state.addToast);
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'buyer' | 'seller'>('buyer');
  const [statusFilter, setStatusFilter] = useState<Order['status'] | 'all'>('all');

  const fetchOrders = useCallback(async () => {
    try {
      setIsLoading(true);
      const params: { role: 'buyer' | 'seller'; status?: Order['status'] } = { role: activeTab };
      if (statusFilter !== 'all') {
        params.status = statusFilter;
      }
      const data = await orderService.getUserOrders(params);
      // Ensure data is an array
      if (Array.isArray(data)) {
        setOrders(data);
      } else {
        console.error('Invalid orders data:', data);
        setOrders([]);
        addToast('error', 'Dữ liệu đơn hàng không hợp lệ');
      }
    } catch (error) {
      console.error('Error fetching orders:', error);
      setOrders([]);
      addToast('error', 'Không thể tải đơn hàng');
    } finally {
      setIsLoading(false);
    }
  }, [activeTab, statusFilter, addToast]);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  const getStatusBadge = (status: Order['status']) => {
    const styles: Record<Order['status'], string> = {
      PENDING_PAYMENT: 'bg-yellow-100 text-yellow-800',
      PAID: 'bg-blue-100 text-blue-800',
      ADDRESS_PROVIDED: 'bg-indigo-100 text-indigo-800',
      SHIPPING: 'bg-purple-100 text-purple-800',
      DELIVERED: 'bg-teal-100 text-teal-800',
      COMPLETED: 'bg-green-100 text-green-800',
      CANCELLED: 'bg-red-100 text-red-800',
    };
    const labels: Record<Order['status'], string> = {
      PENDING_PAYMENT: 'Chờ thanh toán',
      PAID: 'Đã thanh toán',
      ADDRESS_PROVIDED: 'Đã cung cấp địa chỉ',
      SHIPPING: 'Đang vận chuyển',
      DELIVERED: 'Đã giao hàng',
      COMPLETED: 'Hoàn thành',
      CANCELLED: 'Đã hủy',
    };
    return (
      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${styles[status]}`}>
        {labels[status]}
      </span>
    );
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Đơn hàng</h1>

      {/* Tabs */}
      <div className="flex gap-4 mb-6 border-b border-gray-200">
        <button
          onClick={() => setActiveTab('buyer')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'buyer'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          <ShoppingCart className="w-4 h-4 inline mr-2" />
          Đơn mua
        </button>
        <button
          onClick={() => setActiveTab('seller')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'seller'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          <Package className="w-4 h-4 inline mr-2" />
          Đơn bán
        </button>
      </div>

      {/* Status Filter */}
      <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
        {['all', 'PENDING_PAYMENT', 'PAID', 'ADDRESS_PROVIDED', 'SHIPPING', 'DELIVERED', 'COMPLETED', 'CANCELLED'].map((status) => (
          <button
            key={status}
            onClick={() => setStatusFilter(status as Order['status'] | 'all')}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              statusFilter === status
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            {status === 'all' ? 'Tất cả' : 
             status === 'PENDING_PAYMENT' ? 'Chờ thanh toán' :
             status === 'PAID' ? 'Đã thanh toán' :
             status === 'ADDRESS_PROVIDED' ? 'Đã cung cấp địa chỉ' :
             status === 'SHIPPING' ? 'Đang giao' :
             status === 'DELIVERED' ? 'Đã giao hàng' :
             status === 'COMPLETED' ? 'Hoàn thành' : 'Đã hủy'}
          </button>
        ))}
      </div>

      {/* Orders List */}
      {isLoading ? (
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="bg-white rounded-lg p-6 shadow-sm animate-pulse">
              <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
              <div className="h-6 bg-gray-200 rounded w-3/4"></div>
            </div>
          ))}
        </div>
      ) : orders.length === 0 ? (
        <div className="text-center py-16">
          <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 text-lg">Không có đơn hàng nào</p>
        </div>
      ) : (
        <div className="space-y-4">
          {orders.map((order) => (
            <div
              key={order.id}
              className="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden cursor-pointer"
              onClick={() => navigate(`/orders/${order.id}`)}
            >
              <div className="p-6">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <p className="text-sm text-gray-500 mb-1">
                      Mã đơn hàng: #{order.id}
                    </p>
                    <p className="text-xs text-gray-400">
                      <Clock className="w-3 h-3 inline mr-1" />
                      {formatRelativeTime(order.created_at)}
                    </p>
                  </div>
                  {getStatusBadge(order.status)}
                </div>

                <div className="flex items-center gap-4 mb-4">
                  {order.product_image && (
                    <img
                      src={order.product_image}
                      alt={order.product_name}
                      className="w-20 h-20 object-cover rounded-lg"
                    />
                  )}
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900 mb-2">
                      {order.product_name || 'Sản phẩm đấu giá'}
                    </h3>
                    <p className="text-sm text-gray-600 flex items-center gap-2">
                      <User className="w-4 h-4" />
                      {activeTab === 'buyer' 
                        ? `Người bán: ${order.seller_info?.username || 'N/A'}`
                        : `Người mua: ${order.buyer_info?.username || 'N/A'}`
                      }
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-gray-600 mb-1">Tổng tiền:</p>
                    <p className="text-xl font-bold text-blue-600">
                      {formatCurrency(order.final_price)}
                    </p>
                  </div>
                </div>

                {order.shipping_address && (
                  <p className="text-sm text-gray-600 mb-2">
                    📍 Địa chỉ: {order.shipping_address}
                  </p>
                )}

                {order.tracking_number && (
                  <p className="text-sm text-gray-600 mb-2">
                    📦 Mã vận đơn: <span className="font-mono">{order.tracking_number}</span>
                  </p>
                )}

                <div className="flex gap-2 mt-4 pt-4 border-t">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/orders/${order.id}`);
                    }}
                    className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
                  >
                    Xem chi tiết
                  </button>
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/orders/${order.id}?tab=chat`);
                    }}
                    className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors text-sm font-medium"
                  >
                    <MessageSquare className="w-4 h-4 inline mr-1" />
                    Chat
                  </button>
                  {order.status === 'DELIVERED' && !order.rating?.[activeTab === 'buyer' ? 'buyer_rating' : 'seller_rating'] && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        navigate(`/orders/${order.id}?tab=rating`);
                      }}
                      className="px-4 py-2 bg-yellow-100 text-yellow-700 rounded-lg hover:bg-yellow-200 transition-colors text-sm font-medium"
                    >
                      <Star className="w-4 h-4 inline mr-1" />
                      Đánh giá
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default OrdersPage;
