import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Heart, Trash2, Eye, Clock, Gavel } from 'lucide-react';
import { watchlistService, WatchlistItem } from '../../services/watchlist.service';
import { useUIStore } from '../../stores/ui.store';
import { formatCurrency, formatDate } from '../../utils/formatters';
import { CountdownTimer } from '../../components/UI/CountdownTimer';
import { LoadingSpinner } from '../../components/Common/Loading';

const WatchlistPage = () => {
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const addToast = useUIStore((state) => state.addToast);

  const loadWatchlist = async () => {
    setIsLoading(true);
    try {
      const data = await watchlistService.getWatchlist();
      setWatchlist(data);
    } catch (error) {
      console.error('Failed to load watchlist:', error);
      addToast('error', 'Không thể tải danh sách yêu thích');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadWatchlist();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleRemove = async (productId: number) => {
    try {
      await watchlistService.removeFromWatchlist(productId);
      setWatchlist(prev => prev.filter(item => item.productId !== productId));
      addToast('success', 'Đã xóa khỏi danh sách yêu thích');
    } catch (error) {
      console.error('Failed to remove from watchlist:', error);
      addToast('error', 'Không thể xóa khỏi danh sách yêu thích');
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-3 mb-6">
          <Heart className="w-8 h-8 text-red-500 fill-red-500" />
          <h1 className="text-3xl font-bold text-gray-900">Sản phẩm yêu thích</h1>
        </div>

        {watchlist.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <Heart className="w-24 h-24 text-gray-300 mx-auto mb-4" />
            <h2 className="text-2xl font-semibold text-gray-900 mb-2">
              Chưa có sản phẩm yêu thích
            </h2>
            <p className="text-gray-600 mb-6">
              Thêm sản phẩm vào danh sách yêu thích để theo dõi dễ dàng hơn
            </p>
            <Link
              to="/"
              className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Khám phá sản phẩm
            </Link>
          </div>
        ) : (
          <div className="bg-white rounded-lg shadow-sm overflow-hidden">
            <div className="p-4 bg-gray-50 border-b border-gray-200">
              <p className="text-sm text-gray-600">
                Bạn có <span className="font-semibold text-gray-900">{watchlist.length}</span> sản phẩm trong danh sách yêu thích
              </p>
            </div>

            <div className="divide-y divide-gray-200">
              {watchlist.map((item) => {
                const product = item.product;
                const isEnded = new Date(product.endAt) < new Date();

                return (
                  <div key={item.id} className="p-6 hover:bg-gray-50 transition-colors">
                    <div className="flex gap-6">
                      {/* Product Image */}
                      <Link
                        to={`/products/${product.id}`}
                        className="flex-shrink-0 w-48 h-48 bg-gray-200 rounded-lg overflow-hidden group"
                      >
                        <img
                          src={product.thumbnailUrl}
                          alt={product.name}
                          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                          onError={(e) => {
                            e.currentTarget.src = 'https://via.placeholder.com/400?text=No+Image';
                          }}
                        />
                      </Link>

                      {/* Product Info */}
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-4 mb-3">
                          <div className="flex-1">
                            <Link
                              to={`/products/${product.id}`}
                              className="text-xl font-semibold text-gray-900 hover:text-blue-600 transition-colors line-clamp-2"
                            >
                              {product.name}
                            </Link>
                            <p className="text-sm text-gray-600 mt-1">
                              {product.categoryParentName} › {product.categoryName}
                            </p>
                          </div>

                          <button
                            onClick={() => handleRemove(product.id)}
                            className="p-2 text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                            title="Xóa khỏi yêu thích"
                          >
                            <Trash2 className="w-5 h-5" />
                          </button>
                        </div>

                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                          {/* Current Price */}
                          <div>
                            <p className="text-sm text-gray-600 mb-1">Giá hiện tại</p>
                            <p className="text-lg font-bold text-blue-600">
                              {formatCurrency(product.currentPrice)}
                            </p>
                          </div>

                          {/* Buy Now Price */}
                          {product.buyNowPrice && (
                            <div>
                              <p className="text-sm text-gray-600 mb-1">Mua ngay</p>
                              <p className="text-lg font-bold text-green-600">
                                {formatCurrency(product.buyNowPrice)}
                              </p>
                            </div>
                          )}

                          {/* Bid Count */}
                          <div>
                            <p className="text-sm text-gray-600 mb-1">Lượt đấu giá</p>
                            <div className="flex items-center gap-2">
                              <Gavel className="w-4 h-4 text-gray-400" />
                              <p className="text-lg font-semibold text-gray-900">
                                {product.bidCount}
                              </p>
                            </div>
                          </div>

                          {/* Time Remaining */}
                          <div>
                            <p className="text-sm text-gray-600 mb-1">Thời gian còn lại</p>
                            {isEnded ? (
                              <p className="text-sm font-semibold text-red-600">
                                Đã kết thúc
                              </p>
                            ) : (
                              <div className="flex items-center gap-2">
                                <Clock className="w-4 h-4 text-orange-500" />
                                <CountdownTimer endTime={product.endAt} showIcon={false} />
                              </div>
                            )}
                          </div>
                        </div>

                        <div className="flex items-center gap-3">
                          <Link
                            to={`/products/${product.id}`}
                            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                          >
                            <Eye className="w-4 h-4" />
                            Xem chi tiết
                          </Link>

                          {isEnded && (
                            <span className="px-4 py-2 bg-red-100 text-red-700 rounded-lg text-sm font-medium">
                              Đấu giá đã kết thúc
                            </span>
                          )}

                          <span className="text-sm text-gray-500">
                            Thêm vào: {formatDate(item.createdAt)}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default WatchlistPage;
