import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { TrendingUp, Eye, Clock, Trophy, AlertCircle, Gavel } from 'lucide-react';
import { bidService } from '../../services/bid.service';
import { useAuthStore } from '../../stores/auth.store';
import { useUIStore } from '../../stores/ui.store';
import { UserBidResponse } from '../../types';
import { formatCurrency, formatDate } from '../../utils/formatters';
import { CountdownTimer } from '../../components/UI/CountdownTimer';
import { LoadingSpinner } from '../../components/Common/Loading';

interface BidWithStatus extends UserBidResponse {
  isWinning: boolean;
  isEnded: boolean;
  isWon: boolean;
}

const MyBidsPage = () => {
  const [bids, setBids] = useState<BidWithStatus[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState<'all' | 'winning' | 'outbid'>('all');
  const { user } = useAuthStore();
  const addToast = useUIStore((state) => state.addToast);

  const loadMyBids = async () => {
    if (!user) return;

    setIsLoading(true);
    try {
      const response = await bidService.searchBidHistory({
        bidderId: user.id,
        status: 'SUCCESS',
        size: 100,
      });

      console.log(response);

      // Map backend response sang BidWithStatus
      const bidsWithStatus: BidWithStatus[] = response.content.map((bid) => {
        const isEnded = bid.endAt ? new Date(bid.endAt) <= new Date() : false;
        const isWinning = !isEnded && bid.currentBidder === user.id;
        const isWon = isEnded && bid.currentBidder === user.id;

        return {
          ...bid,
          isEnded,
          isWinning,
          isWon
        };
      });

      // Group theo productId -> lấy bid mới nhất cho mỗi sản phẩm
      const bidsByProduct = new Map<number, BidWithStatus>();
      bidsWithStatus.forEach((bid) => {
        const existing = bidsByProduct.get(bid.productId);
        if (!existing || new Date(bid.bidCreatedAt) > new Date(existing.bidCreatedAt)) {
          bidsByProduct.set(bid.productId, bid);
        }
      });

      setBids(Array.from(bidsByProduct.values()));
    } catch (error) {
      console.error('Failed to load bids:', error);
      addToast('error', 'Không thể tải danh sách đấu giá');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadMyBids();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user]);

  const filteredBids = bids.filter((bid) => {
    if (filter === 'winning') return bid.isWinning && !bid.isEnded;
    if (filter === 'outbid') return !bid.isWinning && !bid.isEnded;
    return true;
  });

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
          <TrendingUp className="w-8 h-8 text-blue-500" />
          <h1 className="text-3xl font-bold text-gray-900">Đấu giá của tôi</h1>
        </div>

        {bids.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <Gavel className="w-24 h-24 text-gray-300 mx-auto mb-4" />
            <h2 className="text-2xl font-semibold text-gray-900 mb-2">
              Chưa tham gia đấu giá nào
            </h2>
            <p className="text-gray-600 mb-6">
              Bắt đầu đấu giá để có cơ hội sở hữu sản phẩm yêu thích
            </p>
            <Link
              to="/"
              className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Khám phá sản phẩm
            </Link>
          </div>
        ) : (
          <>
            {/* Filter Tabs */}
            <div className="bg-white rounded-lg shadow-sm mb-6">
              <div className="flex border-b border-gray-200">
                <button
                  onClick={() => setFilter('all')}
                  className={`flex-1 px-6 py-4 text-sm font-medium transition-colors ${filter === 'all'
                    ? 'text-blue-600 border-b-2 border-blue-600'
                    : 'text-gray-600 hover:text-gray-900'
                    }`}
                >
                  Tất cả ({bids.length})
                </button>
                <button
                  onClick={() => setFilter('winning')}
                  className={`flex-1 px-6 py-4 text-sm font-medium transition-colors ${filter === 'winning'
                    ? 'text-blue-600 border-b-2 border-blue-600'
                    : 'text-gray-600 hover:text-gray-900'
                    }`}
                >
                  Đang thắng ({bids.filter((b) => b.isWinning && !b.isEnded).length})
                </button>
                <button
                  onClick={() => setFilter('outbid')}
                  className={`flex-1 px-6 py-4 text-sm font-medium transition-colors ${filter === 'outbid'
                    ? 'text-blue-600 border-b-2 border-blue-600'
                    : 'text-gray-600 hover:text-gray-900'
                    }`}
                >
                  Bị vượt giá ({bids.filter((b) => !b.isWinning && !b.isEnded).length})
                </button>
              </div>
            </div>

            {/* Bids List */}
            <div className="bg-white rounded-lg shadow-sm overflow-hidden">
              <div className="divide-y divide-gray-200">
                {filteredBids.length === 0 ? (
                  <div className="p-12 text-center">
                    <AlertCircle className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                    <p className="text-gray-600">Không có đấu giá nào trong danh mục này</p>
                  </div>
                ) : (
                  filteredBids.map((bid) => (
                    <div key={bid.id} className="p-6 hover:bg-gray-50 transition-colors">
                      <div className="flex gap-6">
                        {/* Product Image */}
                        <Link
                          to={`/products/${bid.productId}`}
                          className="flex-shrink-0 w-32 h-32 bg-gray-200 rounded-lg overflow-hidden group"
                        >
                          {bid.thumbnailUrl ? (
                            <img
                              src={bid.thumbnailUrl}
                              alt={bid.productName || 'Product'}
                              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                              onError={(e) => {
                                e.currentTarget.src =
                                  'https://via.placeholder.com/400?text=No+Image';
                              }}
                            />
                          ) : (
                            <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-200 to-gray-300">
                              <Gavel className="w-12 h-12 text-gray-400" />
                            </div>
                          )}
                        </Link>

                        {/* Bid Info */}
                        <div className="flex-1">
                          <div className="flex items-start justify-between gap-4 mb-3">
                            <div>
                              <Link
                                to={`/products/${bid.productId}`}
                                className="text-lg font-semibold text-gray-900 hover:text-blue-600 transition-colors line-clamp-1"
                              >
                                {bid.productName || `Sản phẩm #${bid.productId}`}
                              </Link>
                              <p className="text-sm text-gray-500 mt-1">
                                Đấu giá lúc: {formatDate(bid.bidCreatedAt)}
                              </p>
                            </div>

                            {/* Status Badge */}
                            {bid.isWon ? (
                              <span className="px-3 py-1 bg-yellow-100 text-green-800 rounded-full text-sm font-medium flex items-center gap-1">
                                <Trophy className="w-4 h-4" />
                                Đã thắng
                              </span>
                            ) : bid.isEnded ? (
                              <span className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm font-medium">
                                Đã kết thúc
                              </span>
                            ) : bid.isWinning ? (
                              <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium flex items-center gap-1">
                                <Trophy className="w-4 h-4" />
                                Đang thắng
                              </span>
                            ) : (
                              <span className="px-3 py-1 bg-orange-100 text-orange-700 rounded-full text-sm font-medium flex items-center gap-1">
                                <AlertCircle className="w-4 h-4" />
                                Bị vượt giá
                              </span>
                            )}
                          </div>

                          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                            {/* My Bid */}
                            <div>
                              <p className="text-sm text-gray-600 mb-1">Giá đấu của bạn</p>
                              <p className="text-lg font-bold text-blue-600">
                                {formatCurrency(bid.bidAmount)}
                              </p>
                            </div>

                            {/* Current Price */}
                            {bid.currentPrice && (
                              <div>
                                <p className="text-sm text-gray-600 mb-1">Giá hiện tại</p>
                                <p className="text-lg font-bold text-gray-900">
                                  {formatCurrency(bid.currentPrice)}
                                </p>
                              </div>
                            )}

                            {/* Time Remaining */}
                            {bid.endAt && (
                              <div>
                                <p className="text-sm text-gray-600 mb-1">Thời gian còn lại</p>
                                {bid.isEnded ? (
                                  <p className="text-sm font-semibold text-red-600">
                                    Đã kết thúc
                                  </p>
                                ) : (
                                  <div className="flex items-center gap-2">
                                    <Clock className="w-4 h-4 text-orange-500" />
                                    <CountdownTimer endTime={bid.endAt} showIcon={false} />
                                  </div>
                                )}
                              </div>
                            )}
                          </div>

                          <Link
                            to={`/products/${bid.productId}`}
                            className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                          >
                            <Eye className="w-4 h-4" />
                            Xem chi tiết
                          </Link>
                        </div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default MyBidsPage;
