import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../../hooks/useAuth';
import { apiClient } from '../../../services/api/client';
import { endpoints } from '../../../services/api/endpoints';
import { Product } from '../../../types';
import { formatCurrency, formatRelativeTime } from '../../../utils/formatters';
import { Plus, Package, Clock, Gavel, Edit, Eye } from 'lucide-react';
import { ProductSkeleton } from '../../../components/Common/Loading';

const MyProductsPage = () => {
  const { user } = useAuth();
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState<'all' | 'active' | 'ended'>('all');

  const fetchMyProducts = useCallback(async () => {
    if (!user?.id) return;
    
    try {
      setIsLoading(true);
      const response = await apiClient.get<{ success: boolean; data: Product[] }>(
        endpoints.products.bySeller(user.id)
      );
      setProducts(response.data.data || []);
    } catch (error) {
      console.error('Error fetching products:', error);
    } finally {
      setIsLoading(false);
    }
  }, [user]);

  useEffect(() => {
    fetchMyProducts();
  }, [fetchMyProducts]);

  const getFilteredProducts = () => {
    const now = Date.now();
    return products.filter((product) => {
      const endTime = new Date(product.endAt).getTime();
      if (filter === 'active') return endTime > now;
      if (filter === 'ended') return endTime <= now;
      return true;
    });
  };

  const filteredProducts = getFilteredProducts();

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Sản phẩm của tôi</h1>
          <p className="text-gray-600 mt-2">Quản lý các sản phẩm đấu giá của bạn</p>
        </div>
        <Link
          to="/seller/products/create"
          className="flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-sm"
        >
          <Plus className="w-5 h-5" />
          Đăng sản phẩm mới
        </Link>
      </div>

      {/* Filter Tabs */}
      <div className="flex gap-4 mb-6 border-b border-gray-200">
        <button
          onClick={() => setFilter('all')}
          className={`px-4 py-2 font-medium transition-colors ${
            filter === 'all'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          Tất cả ({products.length})
        </button>
        <button
          onClick={() => setFilter('active')}
          className={`px-4 py-2 font-medium transition-colors ${
            filter === 'active'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          Đang đấu giá
        </button>
        <button
          onClick={() => setFilter('ended')}
          className={`px-4 py-2 font-medium transition-colors ${
            filter === 'ended'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900'
          }`}
        >
          Đã kết thúc
        </button>
      </div>

      {/* Products Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: 6 }).map((_, i) => (
            <ProductSkeleton key={i} />
          ))}
        </div>
      ) : filteredProducts.length === 0 ? (
        <div className="text-center py-16">
          <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 text-lg mb-4">
            {filter === 'all'
              ? 'Bạn chưa có sản phẩm nào'
              : filter === 'active'
              ? 'Không có sản phẩm đang đấu giá'
              : 'Không có sản phẩm đã kết thúc'}
          </p>
          {filter === 'all' && (
            <Link
              to="/seller/products/create"
              className="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              <Plus className="w-5 h-5" />
              Đăng sản phẩm đầu tiên
            </Link>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredProducts.map((product) => {
            const isActive = new Date(product.endAt).getTime() > Date.now();
            return (
              <div
                key={product.id}
                className="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden"
              >
                <div className="relative h-48 bg-gray-200">
                  <img
                    src={product.thumbnailUrl}
                    alt={product.name}
                    className="w-full h-full object-cover"
                  />
                  <div
                    className={`absolute top-2 right-2 px-3 py-1 rounded-full text-xs font-semibold ${
                      isActive
                        ? 'bg-green-500 text-white'
                        : 'bg-gray-500 text-white'
                    }`}
                  >
                    {isActive ? 'Đang đấu giá' : 'Đã kết thúc'}
                  </div>
                </div>

                <div className="p-4">
                  <h3 className="font-semibold text-gray-900 line-clamp-2 mb-3">
                    {product.name}
                  </h3>

                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">Giá hiện tại:</span>
                      <span className="font-bold text-blue-600">
                        {formatCurrency(product.currentPrice)}
                      </span>
                    </div>

                    {product.buyNowPrice && (
                      <div className="flex justify-between items-center">
                        <span className="text-gray-600">Mua ngay:</span>
                        <span className="font-semibold text-green-600">
                          {formatCurrency(product.buyNowPrice)}
                        </span>
                      </div>
                    )}

                    <div className="flex items-center justify-between pt-2 border-t text-xs text-gray-500">
                      <div className="flex items-center gap-1">
                        <Gavel className="w-4 h-4" />
                        <span>
                          {product.highestBidder?.username
                            ? `Cao nhất: ${product.highestBidder.username}`
                            : 'Chưa có lượt đấu giá'}
                        </span>
                      </div>
                    </div>

                    <div
                      className={`flex items-center gap-1 text-xs ${
                        isActive ? 'text-gray-600' : 'text-red-600'
                      }`}
                    >
                      <Clock className="w-4 h-4" />
                      <span>
                        {isActive
                          ? `Còn ${formatRelativeTime(product.endAt)}`
                          : `Kết thúc ${formatRelativeTime(product.endAt)}`}
                      </span>
                    </div>
                  </div>

                  <div className="flex gap-2 mt-4">
                    <Link
                      to={`/products/${product.id}`}
                      className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors text-sm font-medium"
                    >
                      <Eye className="w-4 h-4" />
                      Xem
                    </Link>
                    <Link
                      to={`/seller/products/${product.id}/edit`}
                      className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
                    >
                      <Edit className="w-4 h-4" />
                      Sửa
                    </Link>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default MyProductsPage;
