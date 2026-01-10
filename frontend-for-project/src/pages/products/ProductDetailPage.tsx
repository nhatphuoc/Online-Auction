import { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productService } from '../../services/product.service';
import { bidService } from '../../services/bid.service';
import { commentService } from '../../services/comment.service';
import { watchlistService } from '../../services/watchlist.service';
import { useAuthStore } from '../../stores/auth.store';
import { useUIStore } from '../../stores/ui.store';
import { Product, BidHistory, Comment, ProductListItem } from '../../types';
import { formatCurrency, formatDate, formatBidderName } from '../../utils/formatters';
import { CountdownTimer } from '../../components/UI/CountdownTimer';
import { Modal, ConfirmDialog } from '../../components/UI/Modal';
import { LoadingSpinner } from '../../components/Common/Loading';
import {
  Heart,
  Share2,
  AlertCircle,
  User,
  Gavel,
  Clock,
  DollarSign,
  ChevronLeft,
  ChevronRight,
  Send,
  MessageCircle,
} from 'lucide-react';

type TabType = 'description' | 'bidHistory' | 'questions';

const ProductDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuthStore();
  const addToast = useUIStore((state) => state.addToast);

  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [activeTab, setActiveTab] = useState<TabType>('description');

  // Bidding
  const [bidAmount, setBidAmount] = useState(0);
  const [showBidModal, setShowBidModal] = useState(false);
  const [showConfirmBid, setShowConfirmBid] = useState(false);
  const [bidHistory, setBidHistory] = useState<BidHistory[]>([]);
  const [isLoadingBids, setIsLoadingBids] = useState(false);

  // Buy Now
  const [showConfirmBuyNow, setShowConfirmBuyNow] = useState(false);
  const [isBuyingNow, setIsBuyingNow] = useState(false);

  // Comments/Q&A
  const [comments, setComments] = useState<Comment[]>([]);
  const [newQuestion, setNewQuestion] = useState('');
  const [isLoadingComments, setIsLoadingComments] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  // Related products
  const [relatedProducts, setRelatedProducts] = useState<ProductListItem[]>([]);

  // Watchlist
  const [isInWatchlist, setIsInWatchlist] = useState(false);
  const [isTogglingWatchlist, setIsTogglingWatchlist] = useState(false);

  // Auto Bid
  const [showAutoBidModal, setShowAutoBidModal] = useState(false);
  const [autoBidMax, setAutoBidMax] = useState(0);
  const [isRegisteringAutoBid, setIsRegisteringAutoBid] = useState(false);

  useEffect(() => {
    if (product) {
      setAutoBidMax(product.currentPrice + product.stepPrice * 5);
    }
  }, [product]);


  // Check if product is in watchlist
  const checkWatchlistStatus = useCallback(async () => {
    if (!id || !isAuthenticated) return;

    try {
      const inWatchlist = await watchlistService.isInWatchlist(parseInt(id));
      setIsInWatchlist(inWatchlist);
    } catch (error) {
      console.error('Failed to check watchlist status:', error);
    }
  }, [id, isAuthenticated]);

  // Toggle watchlist
  const handleToggleWatchlist = async () => {
    if (!isAuthenticated) {
      addToast('error', 'Vui lòng đăng nhập để sử dụng tính năng này');
      return;
    }

    setIsTogglingWatchlist(true);
    try {
      if (isInWatchlist) {
        await watchlistService.removeFromWatchlist(parseInt(id!));
        setIsInWatchlist(false);
        addToast('success', 'Đã xóa khỏi danh sách yêu thích');
      } else {
        await watchlistService.addToWatchlist(parseInt(id!));
        setIsInWatchlist(true);
        addToast('success', 'Đã thêm vào danh sách yêu thích');
      }
    } catch (error) {
      console.error('Failed to toggle watchlist:', error);
      addToast('error', 'Không thể cập nhật danh sách yêu thích');
    } finally {
      setIsTogglingWatchlist(false);
    }
  };

  const loadProductDetail = useCallback(async () => {
    if (!id) return;

    setIsLoading(true);
    try {
      const data = await productService.getProductDetail(parseInt(id));
      setProduct(data);
      setBidAmount(data.currentPrice + data.stepPrice);

      // Load related products
      if (data.categoryId) {
        loadRelatedProducts(data.categoryId, parseInt(id));
      }

      // Check watchlist status
      if (isAuthenticated) {
        checkWatchlistStatus();
      }
    } catch (error) {
      console.error('Failed to load product:', error);
      addToast('error', 'Không thể tải thông tin sản phẩm');
    } finally {
      setIsLoading(false);
    }
  }, [id, addToast, isAuthenticated, checkWatchlistStatus]);

  const loadRelatedProducts = async (categoryId: number, currentProductId: number) => {
    try {
      const response = await productService.searchProducts({
        categoryId,
        pageSize: 5,
      });
      if (response.success && response.data) {
        setRelatedProducts(
          response.data.content.filter(p => p.id !== currentProductId).slice(0, 5)
        );
      }
    } catch (error) {
      console.error('Failed to load related products:', error);
    }
  };

  const loadBidHistory = useCallback(async () => {
    if (!id) return;

    setIsLoadingBids(true);
    try {
      const pageData = await bidService.getBidsByProduct(
        parseInt(id),
        {
          status: 'SUCCESS',
          size: 50,
          page: 0,
        }
      );

      // ✅ SET ĐÚNG ARRAY
      setBidHistory(pageData.content);

    } catch (error) {
      console.error('Failed to load bid history:', error);
    } finally {
      setIsLoadingBids(false);
    }
  }, [id]);

  const loadComments = useCallback(async () => {
    if (!id) return;

    setIsLoadingComments(true);
    try {
      const data = await commentService.getProductComments(parseInt(id), {
        limit: 5,
      });
      setComments(data);
    } catch (error) {
      console.error('Failed to load comments:', error);
    } finally {
      setIsLoadingComments(false);
    }
  }, [id]);

  const setupWebSocket = useCallback(async () => {
    if (!id || !isAuthenticated) return;

    try {
      const wsInfo = await commentService.getWebSocketInfo();
      const websocket = commentService.createWebSocketConnection(
        parseInt(id),
        wsInfo.internal_jwt,
        wsInfo.comment_service_websocket_url
      );

      websocket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log('WebSocket message received:', message);
        if (message.type === 'comment' && message.data) {
          // Backend sends: { type: 'comment', data: CommentResponse }
          setComments(prev => [...prev, message.data]);
        }
      };

      websocket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      setWs(websocket);
    } catch (error) {
      console.error('Failed to setup WebSocket:', error);
    }
  }, [id, isAuthenticated]);

  useEffect(() => {
    loadProductDetail();
  }, [loadProductDetail]);

  // Load bid history only once when tab is activated
  useEffect(() => {
    if (activeTab === 'bidHistory' && bidHistory.length === 0) {
      loadBidHistory();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, loadBidHistory]);

  // Load comments only once when tab is activated
  useEffect(() => {
    if (activeTab === 'questions' && comments.length === 0) {
      loadComments();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab, loadComments]);

  useEffect(() => {
    if (isAuthenticated && activeTab === 'questions') {
      setupWebSocket();
    }

    return () => {
      if (ws) {
        ws.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, activeTab]);

  const handlePlaceBid = async () => {
    if (!product || !isAuthenticated) {
      addToast('error', 'Vui lòng đăng nhập để đấu giá');
      navigate('/login');
      return;
    }

    if (bidAmount < product.currentPrice + product.stepPrice) {
      addToast('error', `Giá đấu tối thiểu là ${formatCurrency(product.currentPrice + product.stepPrice)}`);
      return;
    }

    try {
      const requestId = `bid-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      const response = await bidService.placeBid({
        productId: product.id,
        amount: bidAmount,
        requestId,
      });

      if (response.success) {
        addToast('success', 'Đấu giá thành công!');
        setShowBidModal(false);
        setShowConfirmBid(false);
        await loadProductDetail();
        await loadBidHistory();
      } else {
        addToast('error', response.message || 'Đấu giá thất bại');
      }
    } catch {
      addToast('error', 'Lỗi khi đấu giá. Vui lòng thử lại');
    }
  };

  const handleBuyNow = async () => {
    if (!product || !isAuthenticated) {
      addToast('error', 'Vui lòng đăng nhập để mua ngay');
      navigate('/login');
      return;
    }

    setShowConfirmBuyNow(false);
    setIsBuyingNow(true);

    try {
      const response = await productService.buyNow(product.id);

      if (response.success) {
        addToast('success', 'Mua ngay thành công! Đơn hàng đang được xử lý');
        await loadProductDetail();
        // Redirect to orders page after a short delay
        setTimeout(() => {
          navigate('/orders');
        }, 2000);
      } else {
        addToast('error', response.message || 'Mua ngay thất bại');
      }
    } catch (error) {
      const err = error as { response?: { data?: { message?: string } } };
      const errorMessage = err?.response?.data?.message || 'Lỗi khi mua ngay. Vui lòng thử lại';
      addToast('error', errorMessage);
    } finally {
      setIsBuyingNow(false);
    }
  };

  const handleSendQuestion = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newQuestion.trim() || !ws || !isAuthenticated) {
      if (!isAuthenticated) {
        addToast('error', 'Vui lòng đăng nhập để đặt câu hỏi');
        navigate('/login');
      }
      return;
    }

    try {
      const message = {
        type: 'comment',
        content: newQuestion.trim(),
      };

      console.log('Sending WebSocket message:', message);
      ws.send(JSON.stringify(message));
      setNewQuestion('');
      addToast('success', 'Câu hỏi của bạn đã được gửi!');
    } catch (error) {
      console.error('Failed to send question:', error);
      addToast('error', 'Không thể gửi câu hỏi');
    }
  };

  const nextImage = () => {
    if (product && product.images && product.images.length > 0) {
      const total = [product.thumbnailUrl, ...product.images].length;
      setCurrentImageIndex((prev) => (prev + 1) % total);
    }
  };

  const prevImage = () => {
    if (product && product.images && product.images.length > 0) {
      const total = [product.thumbnailUrl, ...product.images].length;
      setCurrentImageIndex((prev) => (prev - 1 + total) % total);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  if (!product) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <AlertCircle className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Không tìm thấy sản phẩm</h2>
          <button
            onClick={() => navigate('/')}
            className="mt-4 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Về trang chủ
          </button>
        </div>
      </div>
    );
  }

  const allImages = [product.thumbnailUrl, ...(product.images || [])];
  const isAuctionEnded = new Date(product.endAt) < new Date();
  const isSeller = user?.id === product.sellerId;
  const suggestedBid = product.currentPrice + product.stepPrice;

  function maskName(name?: string, visibleChars = 2): string {
    if (!name) return '';
    if (name.length <= visibleChars) return '*'.repeat(name.length);

    const maskedLength = name.length - visibleChars;
    return '*'.repeat(maskedLength) + name.slice(-visibleChars);
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Main Content */}
        <div className="bg-white rounded-lg shadow-sm overflow-hidden mb-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 p-6">
            {/* Image Gallery */}
            <div>
              <div className="relative aspect-square bg-gray-200 rounded-lg overflow-hidden mb-4">
                {allImages[currentImageIndex] ? (
                  <img
                    src={allImages[currentImageIndex]}
                    alt={product.name}
                    className="w-full h-full object-cover"
                    onError={(e) => {
                      e.currentTarget.style.display = 'none';
                      const parent = e.currentTarget.parentElement;
                      if (parent) {
                        parent.innerHTML = '<div class="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-200 to-gray-300"><div class="text-center text-gray-500"><svg class="w-24 h-24 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg><p class="text-sm">Hình ảnh không khả dụng</p></div></div>';
                      }
                    }}
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-200 to-gray-300">
                    <div className="text-center text-gray-500">
                      <svg className="w-24 h-24 mx-auto mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                      </svg>
                      <p className="text-sm">Chưa có hình ảnh</p>
                    </div>
                  </div>
                )}

                {allImages.length > 1 && (
                  <>
                    <button
                      onClick={prevImage}
                      className="absolute left-2 top-1/2 -translate-y-1/2 p-2 bg-black/50 hover:bg-black/70 text-white rounded-full transition-colors"
                    >
                      <ChevronLeft className="w-6 h-6" />
                    </button>
                    <button
                      onClick={nextImage}
                      className="absolute right-2 top-1/2 -translate-y-1/2 p-2 bg-black/50 hover:bg-black/70 text-white rounded-full transition-colors"
                    >
                      <ChevronRight className="w-6 h-6" />
                    </button>
                  </>
                )}
              </div>

              {/* Thumbnails */}
              {allImages.length > 1 && (
                <div className="grid grid-cols-5 gap-2">
                  {allImages.map((img, idx) => (
                    <button
                      key={idx}
                      onClick={() => setCurrentImageIndex(idx)}
                      className={`aspect-square rounded-lg overflow-hidden border-2 transition-colors ${idx === currentImageIndex
                        ? 'border-blue-600'
                        : 'border-gray-200 hover:border-gray-300'
                        }`}
                    >
                      <img
                        src={img}
                        alt={`${product.name} ${idx + 1}`}
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          const target = e.currentTarget;
                          target.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%23e5e7eb" width="100" height="100"/%3E%3C/svg%3E';
                        }}
                      />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Product Info & Bidding */}
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-4">
                {product.name}
              </h1>

              <div className="flex items-center gap-4 mb-6 flex-wrap">
                <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
                  {product.parentCategoryName} › {product.categoryName}
                </span>
                {isAuctionEnded && (
                  <span className="px-3 py-1 bg-gray-500 text-white rounded-full text-sm font-medium">
                    Đã kết thúc
                  </span>
                )}
                {isSeller && (
                  <span className="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
                    Sản phẩm của bạn
                  </span>
                )}
              </div>

              {/* Seller Info */}
              <div className="bg-gray-50 rounded-lg p-4 mb-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                      {product.sellerInfo.avatarUrl ? (
                        <img src={product.sellerInfo.avatarUrl} alt={product.sellerInfo.fullName} className="w-full h-full rounded-full object-cover" />
                      ) : (
                        <User className="w-6 h-6 text-blue-600" />
                      )}
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">
                        {maskName(product.sellerInfo.fullName)}
                      </p>
                      <p className="text-sm text-gray-500">Người bán</p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Highest Bidder Info */}
              <div className="bg-gray-50 rounded-lg p-4 mb-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                      {product.highestBidder?.avatarUrl ? (
                        <img
                          src={product.highestBidder.avatarUrl}
                          alt={product.highestBidder?.fullName ?? "Người dùng ẩn danh"}
                          className="w-full h-full rounded-full object-cover"
                        />
                      ) : (
                        <User className="w-6 h-6 text-blue-600" />
                      )}
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">
                        {product.highestBidder
                          ? maskName(product.highestBidder.fullName)
                          : "Chưa có người đấu giá"}
                      </p>
                      <p className="text-sm text-gray-500">Người mua cao nhất</p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Pricing */}
              <div className="space-y-4 mb-6">
                <div className="flex justify-between items-baseline border-b pb-4">
                  <span className="text-gray-600">Giá khởi điểm:</span>
                  <span className="text-lg font-semibold">
                    {formatCurrency(product.startingPrice)}
                  </span>
                </div>

                <div className="flex justify-between items-baseline border-b pb-4">
                  <span className="text-gray-600">Giá hiện tại:</span>
                  <span className="text-3xl font-bold text-blue-600">
                    {formatCurrency(product.currentPrice)}
                  </span>
                </div>

                {product.buyNowPrice && (
                  <div className="flex justify-between items-baseline border-b pb-4">
                    <span className="text-gray-600">Giá mua ngay:</span>
                    <span className="text-2xl font-bold text-green-600">
                      {formatCurrency(product.buyNowPrice)}
                    </span>
                  </div>
                )}

                <div className="flex justify-between items-baseline">
                  <span className="text-gray-600">Bước giá:</span>
                  <span className="font-semibold">
                    {formatCurrency(product.stepPrice)}
                  </span>
                </div>
              </div>

              {/* Countdown */}
              {!isAuctionEnded && (
                <div className="bg-orange-50 border border-orange-200 rounded-lg p-4 mb-6">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-700">
                      Thời gian còn lại:
                    </span>
                    <CountdownTimer endTime={product.endAt} />
                  </div>
                </div>
              )}

              {/* Action Buttons */}
              <div className="space-y-3">
                {!isAuctionEnded && !isSeller && isAuthenticated && (
                  <>
                    <button
                      onClick={() => setShowBidModal(true)}
                      className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold transition-colors"
                    >
                      <Gavel className="w-5 h-5" />
                      Đấu giá
                    </button>

                    {/* Auto Bid Button */}
                    <button
                      onClick={() => setShowAutoBidModal(true)}
                      className="w-full flex items-center justify-center gap-2 px-6 py-3 border border-blue-600 text-blue-600 hover:bg-blue-50 rounded-lg font-semibold transition-colors"
                    >
                      <Clock className="w-5 h-5" />
                      Đặt giá tự động (Auto Bid)
                    </button>

                    {product.buyNowPrice && (
                      <button
                        onClick={() => setShowConfirmBuyNow(true)}
                        disabled={isBuyingNow}
                        className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-green-600 hover:bg-green-700 text-white rounded-lg font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <DollarSign className="w-5 h-5" />
                        {isBuyingNow ? 'Đang xử lý...' : `Mua ngay ${formatCurrency(product.buyNowPrice)}`}
                      </button>
                    )}
                  </>
                )}

                {!isAuctionEnded && isSeller && (
                  <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
                    <p className="text-sm text-purple-800 text-center">
                      Bạn không thể đấu giá sản phẩm của chính mình
                    </p>
                  </div>
                )}

                {!isAuthenticated && !isAuctionEnded && (
                  <button
                    onClick={() => navigate('/login')}
                    className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold transition-colors"
                  >
                    Đăng nhập để đấu giá
                  </button>
                )}

                <div className="flex gap-3">
                  {isAuthenticated && !isSeller && (
                    <button
                      onClick={handleToggleWatchlist}
                      disabled={isTogglingWatchlist}
                      className={`flex-1 flex items-center justify-center gap-2 px-4 py-2 border rounded-lg font-medium transition-colors ${isInWatchlist
                        ? 'border-red-300 bg-red-50 text-red-700'
                        : 'border-gray-300 hover:bg-gray-50'
                        } ${isTogglingWatchlist ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                      <Heart className={`w-5 h-5 ${isInWatchlist ? 'fill-current' : ''}`} />
                      {isInWatchlist ? 'Đã lưu' : 'Lưu'}
                    </button>
                  )}
                  <button className="flex-1 flex items-center justify-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium transition-colors">
                    <Share2 className="w-5 h-5" />
                    Chia sẻ
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Tabs Section */}
        <div className="bg-white rounded-lg shadow-sm overflow-hidden mb-6">
          <div className="border-b border-gray-200">
            <nav className="flex">
              <button
                onClick={() => setActiveTab('description')}
                className={`px-6 py-4 border-b-2 font-semibold transition-colors ${activeTab === 'description'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
                  }`}
              >
                Mô tả
              </button>
              <button
                onClick={() => setActiveTab('bidHistory')}
                className={`px-6 py-4 border-b-2 font-semibold transition-colors ${activeTab === 'bidHistory'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
                  }`}
              >
                Lịch sử đấu giá
              </button>
              <button
                onClick={() => setActiveTab('questions')}
                className={`px-6 py-4 border-b-2 font-semibold transition-colors ${activeTab === 'questions'
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
                  }`}
              >
                Câu hỏi
              </button>
            </nav>
          </div>

          <div className="p-6">
            {/* Description Tab */}
            {activeTab === 'description' && (
              <div>
                <div className="prose max-w-none mb-8">
                  <div dangerouslySetInnerHTML={{ __html: product.description }} />
                </div>

                {/* Product Details */}
                <div className="grid grid-cols-2 gap-4 p-4 bg-gray-50 rounded-lg">
                  <div>
                    <p className="text-sm text-gray-600">Ngày đăng</p>
                    <p className="font-semibold">{formatDate(product.createdAt)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Kết thúc</p>
                    <p className="font-semibold">{formatDate(product.endAt)}</p>
                  </div>
                  {product.autoExtend && (
                    <div className="col-span-2">
                      <p className="text-sm text-blue-600 flex items-center gap-2">
                        <Clock className="w-4 h-4" />
                        Tự động gia hạn khi có đấu giá mới trong {product.extendThresholdMinutes || 5} phút cuối
                      </p>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Bid History Tab */}
            {activeTab === 'bidHistory' && (
              <div>
                {isLoadingBids ? (
                  <div className="text-center py-8">
                    <LoadingSpinner />
                  </div>
                ) : bidHistory.length === 0 ? (
                  <div className="text-center py-12 text-gray-500">
                    <Gavel className="w-16 h-16 mx-auto mb-4 opacity-30" />
                    <p>Chưa có lượt đấu giá nào</p>
                  </div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Thời điểm</th>
                          <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase">Người đấu giá</th>
                          <th className="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase">Giá đấu</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-200">
                        {bidHistory.map((bid) => (
                          <tr key={bid.id} className="hover:bg-gray-50">
                            <td className="px-4 py-3 text-sm text-gray-600">
                              {formatDate(bid.createdAt)}
                            </td>
                            <td className="px-4 py-3 text-sm font-medium">
                              {formatBidderName(`Người dùng ${bid.bidderName}`)}
                            </td>
                            <td className="px-4 py-3 text-sm font-bold text-blue-600 text-right">
                              {formatCurrency(bid.amount)}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </div>
            )}

            {/* Questions Tab */}
            {activeTab === 'questions' && (
              <div>
                {isLoadingComments ? (
                  <div className="text-center py-8">
                    <LoadingSpinner />
                  </div>
                ) : (
                  <>
                    {/* Question Form */}
                    {isAuthenticated && (
                      <form onSubmit={handleSendQuestion} className="mb-6 pb-6 border-b">
                        <div className="flex gap-3">
                          <input
                            type="text"
                            value={newQuestion}
                            onChange={(e) => setNewQuestion(e.target.value)}
                            placeholder="Đặt câu hỏi về sản phẩm..."
                            className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                          <button
                            type="submit"
                            disabled={!newQuestion.trim()}
                            className="px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white rounded-lg font-medium transition-colors flex items-center gap-2"
                          >
                            <Send className="w-5 h-5" />
                            Gửi
                          </button>
                        </div>
                      </form>
                    )}

                    {/* Comments List */}
                    {!comments || comments.length === 0 ? (
                      <div className="text-center py-12 text-gray-500">
                        <MessageCircle className="w-16 h-16 mx-auto mb-4 opacity-30" />
                        <p>Chưa có câu hỏi nào</p>
                        {!isAuthenticated && (
                          <button
                            onClick={() => navigate('/login')}
                            className="mt-4 text-blue-600 hover:text-blue-700 font-medium"
                          >
                            Đăng nhập để đặt câu hỏi
                          </button>
                        )}
                      </div>
                    ) : (
                      <div className="space-y-4">
                        {comments.map((comment) => (
                          <div key={comment.id} className="bg-gray-50 rounded-lg p-4">
                            <div className="flex items-start gap-3">
                              <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0">
                                <User className="w-5 h-5 text-blue-600" />
                              </div>
                              <div className="flex-1">
                                <div className="flex items-center justify-between mb-2">
                                  <span className="font-semibold text-gray-900">
                                    {comment.sender_name}
                                  </span>
                                  <span className="text-sm text-gray-500">
                                    {formatDate(comment.created_at)}
                                  </span>
                                </div>
                                <p className="text-gray-700">{comment.content}</p>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Related Products */}
        {relatedProducts.length > 0 && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">
              Sản phẩm liên quan
            </h2>
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4">
              {relatedProducts.map((relProduct) => (
                <a
                  key={relProduct.id}
                  href={`/products/${relProduct.id}`}
                  className="group"
                >
                  <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden mb-2">
                    <img
                      src={relProduct.thumbnailUrl}
                      alt={relProduct.name}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                      onError={(e) => {
                        e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="200" height="200"%3E%3Crect fill="%23e5e7eb" width="200" height="200"/%3E%3C/svg%3E';
                      }}
                    />
                  </div>
                  <h3 className="text-sm font-medium text-gray-900 line-clamp-2 group-hover:text-blue-600">
                    {relProduct.name}
                  </h3>
                  <p className="text-sm font-bold text-blue-600 mt-1">
                    {formatCurrency(relProduct.currentPrice)}
                  </p>
                </a>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Bid Modal */}
      <Modal
        isOpen={showBidModal}
        onClose={() => setShowBidModal(false)}
        title="Đặt giá thầu"
      >
        <div className="space-y-4">
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-sm text-blue-800">
              Giá đề nghị: <span className="font-bold">{formatCurrency(suggestedBid)}</span>
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Số tiền đấu giá
            </label>
            <input
              type="number"
              value={bidAmount}
              onChange={(e) => setBidAmount(parseInt(e.target.value) || 0)}
              min={suggestedBid}
              step={product.stepPrice}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <p className="text-sm text-gray-500 mt-1">
              Tối thiểu: {formatCurrency(suggestedBid)}
            </p>
          </div>

          <button
            onClick={() => {
              setShowBidModal(false);
              setShowConfirmBid(true);
            }}
            disabled={bidAmount < suggestedBid}
            className="w-full px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white rounded-lg font-semibold transition-colors"
          >
            Tiếp tục
          </button>
        </div>
      </Modal>

      {/* Confirm Bid Dialog */}
      <ConfirmDialog
        isOpen={showConfirmBid}
        onClose={() => setShowConfirmBid(false)}
        onConfirm={handlePlaceBid}
        title="Xác nhận đấu giá"
        message={`Bạn có chắc chắn muốn đấu giá ${formatCurrency(bidAmount)} cho sản phẩm này?`}
        confirmText="Xác nhận"
        variant="info"
      />

      {/* Confirm Buy Now Dialog */}
      <ConfirmDialog
        isOpen={showConfirmBuyNow}
        onClose={() => setShowConfirmBuyNow(false)}
        onConfirm={handleBuyNow}
        title="Xác nhận mua ngay"
        message={`Bạn có chắc chắn muốn mua ngay sản phẩm này với giá ${formatCurrency(product?.buyNowPrice || 0)}? Phiên đấu giá sẽ kết thúc ngay lập tức.`}
        confirmText="Mua ngay"
        variant="warning"
      />

      <Modal
        isOpen={showAutoBidModal}
        onClose={() => setShowAutoBidModal(false)}
        title="Đặt giá tự động"
      >
        <div className="space-y-4">
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 text-sm text-yellow-800">
            Hệ thống sẽ tự động đấu giá cho bạn với bước giá tối thiểu,
            cho đến khi đạt mức tối đa bạn đặt.
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Giá tối đa bạn sẵn sàng trả
            </label>
            <input
              type="number"
              value={autoBidMax}
              min={product.currentPrice + product.stepPrice}
              step={product.stepPrice}
              onChange={(e) => setAutoBidMax(Number(e.target.value) || 0)}
              className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
            />
            <p className="text-sm text-gray-500 mt-1">
              Tối thiểu: {formatCurrency(product.currentPrice + product.stepPrice)}
            </p>
          </div>

          <button
            disabled={isRegisteringAutoBid}
            onClick={async () => {
              try {
                setIsRegisteringAutoBid(true);
                await bidService.registerAutoBid(product.id, autoBidMax);
                addToast('success', 'Đã đăng ký đấu giá tự động');
                setShowAutoBidModal(false);
              } catch {
                addToast('error', 'Không thể đăng ký auto bid');
              } finally {
                setIsRegisteringAutoBid(false);
              }
            }}
            className="w-full px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold disabled:opacity-50"
          >
            Xác nhận Auto Bid
          </button>
        </div>
      </Modal>
    </div>
  );
};

export default ProductDetailPage;
