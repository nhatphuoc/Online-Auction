import { useState, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { orderService } from '../../services/order.service';
import { OrderDetail } from '../../types';
import { formatCurrency, formatDate } from '../../utils/formatters';
import { 
  ArrowLeft, Package, CreditCard, MapPin, Truck, CheckCircle, 
  XCircle, MessageSquare, Star, AlertCircle 
} from 'lucide-react';
import { useUIStore } from '../../stores/ui.store';
import { useAuthStore } from '../../stores/auth.store';
import { OrderChat } from '../../components/Order/OrderChat';
import OrderWizard from '../../components/Order/OrderWizard';
import OrderChatErrorBoundary from '../../components/Order/OrderChatErrorBoundary';

const OrderDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const addToast = useUIStore((state) => state.addToast);
  const currentUser = useAuthStore((state) => state.user);
  
  const [order, setOrder] = useState<OrderDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'details' | 'chat' | 'rating'>('details');
  
  // Payment form
  const [paymentMethod, setPaymentMethod] = useState<'MOMO' | 'ZALOPAY' | 'VNPAY' | 'STRIPE' | 'PAYPAL'>('MOMO');
  const [paymentProof, setPaymentProof] = useState('');
  const [isPaymentProcessing, setIsPaymentProcessing] = useState(false);
  
  // Shipping address form
  const [shippingAddress, setShippingAddress] = useState('');
  const [shippingPhone, setShippingPhone] = useState('');
  const [isAddressSubmitting, setIsAddressSubmitting] = useState(false);
  
  // Shipping invoice form
  const [trackingNumber, setTrackingNumber] = useState('');
  const [shippingInvoice, setShippingInvoice] = useState('');
  const [isInvoiceSubmitting, setIsInvoiceSubmitting] = useState(false);
  
  // Cancel order form
  const [cancelReason, setCancelReason] = useState('');
  const [isCancelling, setIsCancelling] = useState(false);
  
  // Rating form
  const [rating, setRating] = useState<1 | -1>(1);
  const [ratingComment, setRatingComment] = useState('');
  const [isRatingSubmitting, setIsRatingSubmitting] = useState(false);

  const isBuyer = order && currentUser && order.winner_id === currentUser.id;
  const isSeller = order && currentUser && order.seller_id === currentUser.id;

  useEffect(() => {
    const tab = searchParams.get('tab');
    if (tab === 'chat' || tab === 'rating') {
      setActiveTab(tab);
    }
  }, [searchParams]);

  useEffect(() => {
    fetchOrderDetail();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const fetchOrderDetail = async () => {
    if (!id) return;
    try {
      setIsLoading(true);
      const data = await orderService.getOrderById(parseInt(id));
      setOrder(data);
    } catch (error) {
      console.error('Error fetching order:', error);
      addToast('error', 'Không thể tải thông tin đơn hàng');
      navigate('/orders');
    } finally {
      setIsLoading(false);
    }
  };

  const handlePayOrder = async () => {
    if (!id || !paymentMethod || isPaymentProcessing) return;
    try {
      setIsPaymentProcessing(true);
      await orderService.payOrder(parseInt(id), {
        payment_method: paymentMethod,
        payment_proof: paymentProof || undefined,
      });
      addToast('success', 'Thanh toán thành công!');
      await fetchOrderDetail();
      setPaymentProof('');
    } catch (error) {
      console.error('Error paying order:', error);
      addToast('error', 'Thanh toán thất bại');
    } finally {
      setIsPaymentProcessing(false);
    }
  };

  const handleProvideAddress = async () => {
    if (!id || !shippingAddress || !shippingPhone || isAddressSubmitting) return;
    try {
      setIsAddressSubmitting(true);
      await orderService.provideShippingAddress(parseInt(id), {
        shipping_address: shippingAddress,
        shipping_phone: shippingPhone,
      });
      addToast('success', 'Cập nhật địa chỉ giao hàng thành công!');
      await fetchOrderDetail();
      setShippingAddress('');
      setShippingPhone('');
    } catch (error) {
      console.error('Error providing address:', error);
      addToast('error', 'Cập nhật địa chỉ thất bại');
    } finally {
      setIsAddressSubmitting(false);
    }
  };

  const handleSendInvoice = async () => {
    if (!id || !trackingNumber || isInvoiceSubmitting) return;
    try {
      setIsInvoiceSubmitting(true);
      await orderService.sendShippingInvoice(parseInt(id), {
        tracking_number: trackingNumber,
        shipping_invoice: shippingInvoice || undefined,
      });
      addToast('success', 'Gửi thông tin vận chuyển thành công!');
      await fetchOrderDetail();
      setTrackingNumber('');
      setShippingInvoice('');
    } catch (error) {
      console.error('Error sending invoice:', error);
      addToast('error', 'Gửi thông tin vận chuyển thất bại');
    } finally {
      setIsInvoiceSubmitting(false);
    }
  };

  const handleConfirmDelivery = async () => {
    if (!id || !window.confirm('Xác nhận đã nhận được hàng?')) return;
    try {
      await orderService.confirmDelivery(parseInt(id));
      addToast('success', 'Xác nhận giao hàng thành công!');
      await fetchOrderDetail();
    } catch (error) {
      console.error('Error confirming delivery:', error);
      addToast('error', 'Xác nhận giao hàng thất bại');
    }
  };

  const handleCancelOrder = async () => {
    if (!id || !cancelReason || isCancelling) return;
    if (!window.confirm('Bạn có chắc muốn hủy đơn hàng này?')) return;
    
    try {
      setIsCancelling(true);
      await orderService.cancelOrder(parseInt(id), { cancel_reason: cancelReason });
      addToast('success', 'Hủy đơn hàng thành công!');
      await fetchOrderDetail();
      setCancelReason('');
    } catch (error) {
      console.error('Error cancelling order:', error);
      addToast('error', 'Hủy đơn hàng thất bại');
    } finally {
      setIsCancelling(false);
    }
  };

  const handleRateOrder = async () => {
    if (!id || isRatingSubmitting) return;
    try {
      setIsRatingSubmitting(true);
      await orderService.rateOrder(parseInt(id), {
        rating,
        comment: ratingComment || undefined,
      });
      addToast('success', 'Đánh giá thành công!');
      await fetchOrderDetail();
      setRatingComment('');
    } catch (error) {
      console.error('Error rating order:', error);
      addToast('error', 'Đánh giá thất bại');
    } finally {
      setIsRatingSubmitting(false);
    }
  };

  const getStatusBadge = (status: OrderDetail['status']) => {
    const styles: Record<OrderDetail['status'], string> = {
      PENDING_PAYMENT: 'bg-yellow-100 text-yellow-800',
      PAID: 'bg-blue-100 text-blue-800',
      ADDRESS_PROVIDED: 'bg-indigo-100 text-indigo-800',
      SHIPPING: 'bg-purple-100 text-purple-800',
      DELIVERED: 'bg-teal-100 text-teal-800',
      COMPLETED: 'bg-green-100 text-green-800',
      CANCELLED: 'bg-red-100 text-red-800',
    };
    const labels: Record<OrderDetail['status'], string> = {
      PENDING_PAYMENT: 'Chờ thanh toán',
      PAID: 'Đã thanh toán',
      ADDRESS_PROVIDED: 'Đã cung cấp địa chỉ',
      SHIPPING: 'Đang vận chuyển',
      DELIVERED: 'Đã giao hàng',
      COMPLETED: 'Hoàn thành',
      CANCELLED: 'Đã hủy',
    };
    return (
      <span className={`px-4 py-2 rounded-full text-sm font-semibold ${styles[status]}`}>
        {labels[status]}
      </span>
    );
  };

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="h-64 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8 text-center">
        <AlertCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
        <p className="text-xl text-gray-600">Không tìm thấy đơn hàng</p>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-6">
        <button
          onClick={() => navigate('/orders')}
          className="flex items-center text-blue-600 hover:text-blue-700 mb-4"
        >
          <ArrowLeft className="w-5 h-5 mr-2" />
          Quay lại
        </button>
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              Đơn hàng #{order.id}
            </h1>
            <p className="text-gray-600">
              Tạo lúc: {formatDate(order.created_at)}
            </p>
          </div>
          {getStatusBadge(order.status)}
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <div className="flex gap-6">
          <button
            onClick={() => setActiveTab('details')}
            className={`pb-4 px-2 font-medium transition-colors ${
              activeTab === 'details'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <Package className="w-4 h-4 inline mr-2" />
            Chi tiết
          </button>
          <button
            onClick={() => setActiveTab('chat')}
            className={`pb-4 px-2 font-medium transition-colors ${
              activeTab === 'chat'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <MessageSquare className="w-4 h-4 inline mr-2" />
            Trò chuyện
          </button>
          <button
            onClick={() => setActiveTab('rating')}
            className={`pb-4 px-2 font-medium transition-colors ${
              activeTab === 'rating'
                ? 'text-blue-600 border-b-2 border-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            <Star className="w-4 h-4 inline mr-2" />
            Đánh giá
          </button>
        </div>
      </div>

      {/* Tab Content */}
      {activeTab === 'details' && (
        <div className="space-y-6">
          {/* Order Wizard - Progress Tracker */}
          <OrderWizard order={order} isBuyer={!!isBuyer} />

          {/* Order Info */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold mb-4">Thông tin đơn hàng</h2>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-gray-600 mb-1">Người mua</p>
                <p className="font-medium">{order.buyer_name || order.buyer_info?.username || 'N/A'}</p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">Người bán</p>
                <p className="font-medium">{order.seller_name || order.seller_info?.username || 'N/A'}</p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">Giá cuối</p>
                <p className="font-medium text-blue-600 text-xl">
                  {formatCurrency(order.final_price)}
                </p>
              </div>
              <div>
                <p className="text-gray-600 mb-1">Phương thức thanh toán</p>
                <p className="font-medium">{order.payment_method || 'Chưa thanh toán'}</p>
              </div>
            </div>
          </div>

          {/* Shipping Info */}
          {order.shipping_address && (
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">Thông tin giao hàng</h2>
              <div className="space-y-3">
                <div>
                  <p className="text-gray-600 mb-1">Địa chỉ</p>
                  <p className="font-medium">{order.shipping_address}</p>
                </div>
                <div>
                  <p className="text-gray-600 mb-1">Số điện thoại</p>
                  <p className="font-medium">{order.shipping_phone}</p>
                </div>
                {order.tracking_number && (
                  <div>
                    <p className="text-gray-600 mb-1">Mã vận đơn</p>
                    <p className="font-mono font-medium">{order.tracking_number}</p>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Buyer Actions */}
          {isBuyer && (
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">Hành động</h2>
              
              {/* Pay Order */}
              {order.status === 'PENDING_PAYMENT' && (
                <div className="mb-6 pb-6 border-b">
                  <h3 className="font-medium mb-4 flex items-center">
                    <CreditCard className="w-5 h-5 mr-2 text-blue-600" />
                    Thanh toán đơn hàng
                  </h3>
                  <div className="space-y-3">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Phương thức thanh toán
                      </label>
                      <select
                        value={paymentMethod}
                        onChange={(e) => setPaymentMethod(e.target.value as 'MOMO' | 'ZALOPAY' | 'VNPAY' | 'STRIPE' | 'PAYPAL')}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      >
                        <option value="MOMO">MoMo</option>
                        <option value="ZALOPAY">ZaloPay</option>
                        <option value="VNPAY">VNPay</option>
                        <option value="STRIPE">Stripe</option>
                        <option value="PAYPAL">PayPal</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Link ảnh chứng từ (tùy chọn)
                      </label>
                      <input
                        type="url"
                        value={paymentProof}
                        onChange={(e) => setPaymentProof(e.target.value)}
                        placeholder="https://example.com/proof.jpg"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <button
                      onClick={handlePayOrder}
                      disabled={isPaymentProcessing}
                      className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 font-medium"
                    >
                      {isPaymentProcessing ? 'Đang xử lý...' : 'Thanh toán'}
                    </button>
                  </div>
                </div>
              )}

              {/* Provide Address */}
              {order.status === 'PAID' && (
                <div className="mb-6 pb-6 border-b">
                  <h3 className="font-medium mb-4 flex items-center">
                    <MapPin className="w-5 h-5 mr-2 text-blue-600" />
                    Cung cấp địa chỉ giao hàng
                  </h3>
                  <div className="space-y-3">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Địa chỉ nhận hàng
                      </label>
                      <textarea
                        value={shippingAddress}
                        onChange={(e) => setShippingAddress(e.target.value)}
                        placeholder="Số nhà, tên đường, phường/xã, quận/huyện, tỉnh/thành phố"
                        rows={3}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Số điện thoại
                      </label>
                      <input
                        type="tel"
                        value={shippingPhone}
                        onChange={(e) => setShippingPhone(e.target.value)}
                        placeholder="0901234567"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <button
                      onClick={handleProvideAddress}
                      disabled={isAddressSubmitting || !shippingAddress || !shippingPhone}
                      className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 font-medium"
                    >
                      {isAddressSubmitting ? 'Đang gửi...' : 'Gửi địa chỉ'}
                    </button>
                  </div>
                </div>
              )}

              {/* Confirm Delivery */}
              {order.status === 'SHIPPING' && (
                <button
                  onClick={handleConfirmDelivery}
                  className="w-full px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 font-medium flex items-center justify-center"
                >
                  <CheckCircle className="w-5 h-5 mr-2" />
                  Xác nhận đã nhận hàng
                </button>
              )}
            </div>
          )}

          {/* Seller Actions */}
          {isSeller && (
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">Hành động</h2>
              
              {/* Send Invoice */}
              {order.status === 'ADDRESS_PROVIDED' && (
                <div className="mb-6 pb-6 border-b">
                  <h3 className="font-medium mb-4 flex items-center">
                    <Truck className="w-5 h-5 mr-2 text-blue-600" />
                    Gửi thông tin vận chuyển
                  </h3>
                  <div className="space-y-3">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Mã vận đơn *
                      </label>
                      <input
                        type="text"
                        value={trackingNumber}
                        onChange={(e) => setTrackingNumber(e.target.value)}
                        placeholder="VN123456789"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Link hóa đơn vận chuyển (tùy chọn)
                      </label>
                      <input
                        type="url"
                        value={shippingInvoice}
                        onChange={(e) => setShippingInvoice(e.target.value)}
                        placeholder="https://example.com/invoice.jpg"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <button
                      onClick={handleSendInvoice}
                      disabled={isInvoiceSubmitting || !trackingNumber}
                      className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 font-medium"
                    >
                      {isInvoiceSubmitting ? 'Đang gửi...' : 'Gửi thông tin'}
                    </button>
                  </div>
                </div>
              )}

              {/* Cancel Order */}
              {order.status !== 'COMPLETED' && order.status !== 'CANCELLED' && (
                <div>
                  <h3 className="font-medium mb-4 flex items-center text-red-600">
                    <XCircle className="w-5 h-5 mr-2" />
                    Hủy đơn hàng
                  </h3>
                  <div className="space-y-3">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Lý do hủy
                      </label>
                      <textarea
                        value={cancelReason}
                        onChange={(e) => setCancelReason(e.target.value)}
                        placeholder="Vui lòng nhập lý do hủy đơn hàng..."
                        rows={3}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-red-500"
                      />
                    </div>
                    <button
                      onClick={handleCancelOrder}
                      disabled={isCancelling || !cancelReason}
                      className="w-full px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-gray-300 font-medium"
                    >
                      {isCancelling ? 'Đang hủy...' : 'Hủy đơn hàng'}
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {activeTab === 'chat' && order && (
        <div className="bg-white rounded-lg shadow overflow-hidden" style={{ height: '600px' }}>
          <OrderChatErrorBoundary>
            <OrderChat
              orderId={order.id}
              buyerId={order.winner_id}
              sellerId={order.seller_id}
            />
          </OrderChatErrorBoundary>
        </div>
      )}

      {activeTab === 'rating' && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-6">Đánh giá giao dịch</h2>
          
          {order.rating && (
            <div className="mb-6 pb-6 border-b">
              <h3 className="font-medium mb-4">Đánh giá hiện tại</h3>
              <div className="grid grid-cols-2 gap-6">
                <div>
                  <p className="text-sm text-gray-600 mb-2">Đánh giá của Buyer</p>
                  {order.rating.buyer_rating ? (
                    <>
                      <p className={`font-semibold ${order.rating.buyer_rating === 1 ? 'text-green-600' : 'text-red-600'}`}>
                        {order.rating.buyer_rating === 1 ? '👍 Tích cực' : '👎 Tiêu cực'}
                      </p>
                      {order.rating.buyer_comment && (
                        <p className="text-sm text-gray-700 mt-2">{order.rating.buyer_comment}</p>
                      )}
                    </>
                  ) : (
                    <p className="text-gray-400">Chưa đánh giá</p>
                  )}
                </div>
                <div>
                  <p className="text-sm text-gray-600 mb-2">Đánh giá của Seller</p>
                  {order.rating.seller_rating ? (
                    <>
                      <p className={`font-semibold ${order.rating.seller_rating === 1 ? 'text-green-600' : 'text-red-600'}`}>
                        {order.rating.seller_rating === 1 ? '👍 Tích cực' : '👎 Tiêu cực'}
                      </p>
                      {order.rating.seller_comment && (
                        <p className="text-sm text-gray-700 mt-2">{order.rating.seller_comment}</p>
                      )}
                    </>
                  ) : (
                    <p className="text-gray-400">Chưa đánh giá</p>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Rating Form */}
          {((isBuyer && !order.rating?.buyer_rating) || (isSeller && !order.rating?.seller_rating)) && 
           order.status === 'DELIVERED' && (
            <div>
              <h3 className="font-medium mb-4">
                Đánh giá {isBuyer ? 'người bán' : 'người mua'}
              </h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-3">
                    Đánh giá của bạn
                  </label>
                  <div className="flex gap-4">
                    <button
                      onClick={() => setRating(1)}
                      className={`flex-1 px-6 py-3 rounded-lg border-2 transition-colors ${
                        rating === 1
                          ? 'border-green-600 bg-green-50 text-green-700'
                          : 'border-gray-300 hover:border-green-300'
                      }`}
                    >
                      <span className="text-2xl">👍</span>
                      <p className="font-medium mt-1">Tích cực</p>
                    </button>
                    <button
                      onClick={() => setRating(-1)}
                      className={`flex-1 px-6 py-3 rounded-lg border-2 transition-colors ${
                        rating === -1
                          ? 'border-red-600 bg-red-50 text-red-700'
                          : 'border-gray-300 hover:border-red-300'
                      }`}
                    >
                      <span className="text-2xl">👎</span>
                      <p className="font-medium mt-1">Tiêu cực</p>
                    </button>
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Nhận xét (tùy chọn)
                  </label>
                  <textarea
                    value={ratingComment}
                    onChange={(e) => setRatingComment(e.target.value)}
                    placeholder="Chia sẻ trải nghiệm của bạn..."
                    rows={4}
                    maxLength={500}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <p className="text-xs text-gray-500 mt-1">{ratingComment.length}/500</p>
                </div>
                <button
                  onClick={handleRateOrder}
                  disabled={isRatingSubmitting}
                  className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 font-medium"
                >
                  {isRatingSubmitting ? 'Đang gửi...' : 'Gửi đánh giá'}
                </button>
              </div>
            </div>
          )}

          {order.status !== 'DELIVERED' && order.status !== 'COMPLETED' && (
            <p className="text-gray-500 text-center py-8">
              Bạn chỉ có thể đánh giá sau khi đơn hàng đã được giao
            </p>
          )}
        </div>
      )}
    </div>
  );
};

export default OrderDetailPage;
