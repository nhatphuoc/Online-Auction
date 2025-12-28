import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productService } from '../../services/product.service';
import { bidService } from '../../services/bid.service';
import { commentService } from '../../services/comment.service';
import { useAuthStore } from '../../stores/auth.store';
import { useUIStore } from '../../stores/ui.store';
import { Product, BidHistory, Comment, ProductListItem } from '../../types';
import { formatCurrency, formatDate } from '../../utils/formatters';
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
} from 'lucide-react';

const ProductDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user, isAuthenticated } = useAuthStore();
  const addToast = useUIStore((state) => state.addToast);

  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  
  // Bidding
  const [bidAmount, setBidAmount] = useState(0);
  const [showBidModal, setShowBidModal] = useState(false);
  const [showConfirmBid, setShowConfirmBid] = useState(false);
  const [bidHistory, setBidHistory] = useState<BidHistory[]>([]);
  
  // Comments/Q&A
  const [comments, setComments] = useState<Comment[]>([]);
  
  // Related products
  const [relatedProducts, setRelatedProducts] = useState<ProductListItem[]>([]);
  
  // Watchlist
  const [isInWatchlist, setIsInWatchlist] = useState(false);

  useEffect(() => {
    if (id) {
      const loadData = async () => {
        await loadProductDetail();
        await loadBidHistory();
        await loadComments();
      };
      loadData();
    }
  }, [id]);

  const loadProductDetail = async () => {
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
    } catch (error) {
      console.error('Failed to load product:', error);
      addToast('error', 'Không thể tải thông tin sản phẩm');
    } finally {
      setIsLoading(false);
    }
  };

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

  const loadBidHistory = async () => {
    if (!id) return;
    
    try {
      const response = await bidService.searchBidHistory({
        productId: parseInt(id),
        status: 'SUCCESS',
        size: 20,
      });
      setBidHistory(response.content);
    } catch (error) {
      console.error('Failed to load bid history:', error);
    }
  };

  const loadComments = async () => {
    if (!id) return;
    
    try {
      const data = await commentService.getProductComments(parseInt(id), {
        limit: 50,
      });
      setComments(data);
    } catch (error) {
      console.error('Failed to load comments:', error);
    }
  };

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
        loadProductDetail();
        loadBidHistory();
      } else {
        addToast('error', response.message || 'Đấu giá thất bại');
      }
    } catch {
      addToast('error', 'Lỗi khi đấu giá. Vui lòng thử lại');
    }
  };

  const nextImage = () => {
    if (product) {
      setCurrentImageIndex((prev) => 
        prev === product.images.length - 1 ? 0 : prev + 1
      );
    }
  };

  const prevImage = () => {
    if (product) {
      setCurrentImageIndex((prev) => 
        prev === 0 ? product.images.length - 1 : prev - 1
      );
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

  const allImages = [product.thumbnailUrl, ...product.images];
  const isAuctionEnded = new Date(product.endAt) < new Date();
  const isSeller = user?.id === product.sellerId;
  const suggestedBid = product.currentPrice + product.stepPrice;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Main Content */}
        <div className="bg-white rounded-lg shadow-sm overflow-hidden mb-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 p-6">
            {/* Image Gallery */}
            <div>
              <div className="relative aspect-square bg-gray-100 rounded-lg overflow-hidden mb-4">
                <img
                  src={allImages[currentImageIndex]}
                  alt={product.name}
                  className="w-full h-full object-cover"
                  onError={(e) => {
                    e.currentTarget.src = '/placeholder-image.jpg';
                  }}
                />
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
                      className={`aspect-square rounded-lg overflow-hidden border-2 transition-colors ${
                        idx === currentImageIndex
                          ? 'border-blue-600'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                    >
                      <img
                        src={img}
                        alt={`${product.name} ${idx + 1}`}
                        className="w-full h-full object-cover"
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

              <div className="flex items-center gap-4 mb-6">
                <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
                  {product.categoryName}
                </span>
                {isAuctionEnded && (
                  <span className="px-3 py-1 bg-gray-500 text-white rounded-full text-sm font-medium">
                    Đã kết thúc
                  </span>
                )}
              </div>

              {/* Seller Info */}
              <div className="bg-gray-50 rounded-lg p-4 mb-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                      <User className="w-6 h-6 text-blue-600" />
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">
                        {product.sellerInfo.username}
                      </p>
                      <p className="text-sm text-gray-500">Người bán</p>
                    </div>
                  </div>
                  {/* Seller rating would go here */}
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
                    {product.buyNowPrice && (
                      <button className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-green-600 hover:bg-green-700 text-white rounded-lg font-semibold transition-colors">
                        <DollarSign className="w-5 h-5" />
                        Mua ngay {formatCurrency(product.buyNowPrice)}
                      </button>
                    )}
                  </>
                )}

                <div className="flex gap-3">
                  <button
                    onClick={() => setIsInWatchlist(!isInWatchlist)}
                    className={`flex-1 flex items-center justify-center gap-2 px-4 py-2 border rounded-lg font-medium transition-colors ${
                      isInWatchlist
                        ? 'border-red-300 bg-red-50 text-red-700'
                        : 'border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    <Heart className={`w-5 h-5 ${isInWatchlist ? 'fill-current' : ''}`} />
                    {isInWatchlist ? 'Đã lưu' : 'Lưu'}
                  </button>
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
              <button className="px-6 py-4 border-b-2 border-blue-600 font-semibold text-blue-600">
                Mô tả
              </button>
              <button className="px-6 py-4 border-b-2 border-transparent hover:border-gray-300 font-medium text-gray-600">
                Lịch sử đấu giá ({bidHistory.length})
              </button>
              <button className="px-6 py-4 border-b-2 border-transparent hover:border-gray-300 font-medium text-gray-600">
                Câu hỏi ({comments.length})
              </button>
            </nav>
          </div>

          <div className="p-6">
            {/* Description */}
            <div className="prose max-w-none">
              <div dangerouslySetInnerHTML={{ __html: product.description }} />
            </div>

            {/* Product Details */}
            <div className="mt-8 grid grid-cols-2 gap-4 p-4 bg-gray-50 rounded-lg">
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
                    Tự động gia hạn khi có đấu giá mới
                  </p>
                </div>
              )}
            </div>
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
              onChange={(e) => setBidAmount(parseInt(e.target.value))}
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
    </div>
  );
};

export default ProductDetailPage;
