import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../../hooks/useAuth';
import { apiClient } from '../../../services/api/client';
import { endpoints } from '../../../services/api/endpoints';
import { productService } from '../../../services/product.service';
import { Product } from '../../../types';
import { formatCurrency, formatRelativeTime } from '../../../utils/formatters';
import { Plus, Package, Clock, Gavel, Edit, Eye, Trash2, Loader2 } from 'lucide-react';
import { ProductSkeleton } from '../../../components/Common/Loading';
import { useUIStore } from '../../../stores/ui.store';
import { Modal } from '../../../components/UI/Modal';

const MyProductsPage = () => {
  const { user } = useAuth();
  const addToast = useUIStore((state) => state.addToast);
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState<'all' | 'active' | 'ended'>('all');
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [productToDelete, setProductToDelete] = useState<Product | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const fetchMyProducts = useCallback(async () => {
    if (!user?.id) return;
    
    try {
      setIsLoading(true);
      // API trả về array trực tiếp, không có wrapper
      const response = await apiClient.get<Product[]>(
        endpoints.products.bySeller(user.id)
      );
      setProducts(response.data || []);
    } catch (error) {
      console.error('Error fetching products:', error);
      addToast('error', 'Không thể tải danh sách sản phẩm');
    } finally {
      setIsLoading(false);
    }
  }, [user, addToast]);

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

  const handleDeleteClick = (product: Product) => {
    setProductToDelete(product);
    setDeleteModalOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!productToDelete) return;

    try {
      setIsDeleting(true);
      await productService.deleteProduct(productToDelete.id);
      addToast('success', 'Đã xóa sản phẩm thành công');
      setDeleteModalOpen(false);
      setProductToDelete(null);
      // Refresh product list
      fetchMyProducts();
    } catch (error) {
      console.error('Error deleting product:', error);
      addToast('error', 'Lỗi khi xóa sản phẩm');
    } finally {
      setIsDeleting(false);
    }
  };

  const filteredProducts = getFilteredProducts();

  // Calculate stats
  const stats = {
    total: products.length,
    active: products.filter(p => new Date(p.endAt).getTime() > Date.now()).length,
    ended: products.filter(p => new Date(p.endAt).getTime() <= Date.now()).length,
    totalValue: products.reduce((sum, p) => sum + p.currentPrice, 0),
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Sản phẩm của tôi</h1>
          <p className="text-gray-600 mt-2">Quản lý các sản phẩm đấu giá của bạn</p>
        </div>
        <Link
          to="/seller/products/create"
          className="flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-sm font-medium"
        >
          <Plus className="w-5 h-5" />
          Đăng sản phẩm mới
        </Link>
      </div>

      {/* Stats Cards */}
      {!isLoading && products.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow-sm p-6 border-l-4 border-blue-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 font-medium">Tổng sản phẩm</p>
                <p className="text-3xl font-bold text-gray-900 mt-1">{stats.total}</p>
              </div>
              <Package className="w-12 h-12 text-blue-500 opacity-20" />
            </div>
          </div>
          
          <div className="bg-white rounded-lg shadow-sm p-6 border-l-4 border-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 font-medium">Đang đấu giá</p>
                <p className="text-3xl font-bold text-gray-900 mt-1">{stats.active}</p>
              </div>
              <Clock className="w-12 h-12 text-green-500 opacity-20" />
            </div>
          </div>
          
          <div className="bg-white rounded-lg shadow-sm p-6 border-l-4 border-gray-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 font-medium">Đã kết thúc</p>
                <p className="text-3xl font-bold text-gray-900 mt-1">{stats.ended}</p>
              </div>
              <Gavel className="w-12 h-12 text-gray-500 opacity-20" />
            </div>
          </div>
          
          <div className="bg-white rounded-lg shadow-sm p-6 border-l-4 border-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 font-medium">Tổng giá trị</p>
                <p className="text-2xl font-bold text-gray-900 mt-1">{formatCurrency(stats.totalValue)}</p>
              </div>
              <div className="text-purple-500 opacity-20 text-3xl font-bold">₫</div>
            </div>
          </div>
        </div>
      )}

      {/* Filter Tabs */}
      <div className="bg-white rounded-lg shadow-sm mb-6">
        <div className="flex gap-4 px-6 border-b border-gray-200">
          <button
            onClick={() => setFilter('all')}
            className={`px-4 py-4 font-medium transition-colors relative ${
              filter === 'all'
                ? 'text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Tất cả
            <span className={`ml-2 px-2 py-0.5 text-xs rounded-full ${
              filter === 'all' 
                ? 'bg-blue-100 text-blue-700' 
                : 'bg-gray-100 text-gray-600'
            }`}>
              {products.length}
            </span>
            {filter === 'all' && (
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600"></div>
            )}
          </button>
          <button
            onClick={() => setFilter('active')}
            className={`px-4 py-4 font-medium transition-colors relative ${
              filter === 'active'
                ? 'text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Đang đấu giá
            <span className={`ml-2 px-2 py-0.5 text-xs rounded-full ${
              filter === 'active' 
                ? 'bg-blue-100 text-blue-700' 
                : 'bg-gray-100 text-gray-600'
            }`}>
              {stats.active}
            </span>
            {filter === 'active' && (
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600"></div>
            )}
          </button>
          <button
            onClick={() => setFilter('ended')}
            className={`px-4 py-4 font-medium transition-colors relative ${
              filter === 'ended'
                ? 'text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Đã kết thúc
            <span className={`ml-2 px-2 py-0.5 text-xs rounded-full ${
              filter === 'ended' 
                ? 'bg-blue-100 text-blue-700' 
                : 'bg-gray-100 text-gray-600'
            }`}>
              {stats.ended}
            </span>
            {filter === 'ended' && (
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-600"></div>
            )}
          </button>
        </div>
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
            const hasBids = product.currentPrice > product.startingPrice;
            
            return (
              <div
                key={product.id}
                className="bg-white rounded-lg shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden border border-gray-100"
              >
                {/* Thumbnail */}
                <div className="relative h-56 bg-gray-200 overflow-hidden group">
                  <img
                    src={product.thumbnailUrl}
                    alt={product.name}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    onError={(e) => {
                      e.currentTarget.src = 'https://via.placeholder.com/400x300?text=No+Image';
                    }}
                  />
                  {/* Status Badge */}
                  <div
                    className={`absolute top-3 right-3 px-3 py-1 rounded-full text-xs font-semibold shadow-lg ${
                      isActive
                        ? 'bg-green-500 text-white'
                        : 'bg-gray-700 text-white'
                    }`}
                  >
                    {isActive ? '🟢 Đang đấu giá' : '⚫ Đã kết thúc'}
                  </div>
                  
                  {/* Quick Actions Overlay */}
                  <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-40 transition-all duration-300 flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100">
                    <Link
                      to={`/products/${product.id}`}
                      className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
                      title="Xem chi tiết"
                    >
                      <Eye className="w-5 h-5 text-gray-700" />
                    </Link>
                    <Link
                      to={`/seller/products/${product.id}/edit`}
                      className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
                      title="Chỉnh sửa"
                    >
                      <Edit className="w-5 h-5 text-blue-600" />
                    </Link>
                    <button
                      onClick={() => handleDeleteClick(product)}
                      className="p-2 bg-white rounded-full hover:bg-gray-100 transition-colors"
                      title="Xóa"
                    >
                      <Trash2 className="w-5 h-5 text-red-600" />
                    </button>
                  </div>
                </div>

                {/* Content */}
                <div className="p-5">
                  {/* Product Name */}
                  <h3 className="font-semibold text-gray-900 line-clamp-2 mb-3 text-lg hover:text-blue-600 transition-colors">
                    <Link to={`/products/${product.id}`}>
                      {product.name}
                    </Link>
                  </h3>

                  {/* Price Info */}
                  <div className="space-y-2.5 mb-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-600">Giá khởi điểm:</span>
                      <span className="font-medium text-gray-700">
                        {formatCurrency(product.startingPrice)}
                      </span>
                    </div>
                    
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-600">Giá hiện tại:</span>
                      <span className={`font-bold text-lg ${hasBids ? 'text-blue-600' : 'text-gray-900'}`}>
                        {formatCurrency(product.currentPrice)}
                      </span>
                    </div>

                    {product.buyNowPrice && (
                      <div className="flex justify-between items-center pt-2 border-t border-gray-100">
                        <span className="text-sm text-gray-600">Mua ngay:</span>
                        <span className="font-semibold text-green-600">
                          {formatCurrency(product.buyNowPrice)}
                        </span>
                      </div>
                    )}
                  </div>

                  {/* Bidding Info */}
                  <div className="flex items-center justify-between py-2.5 px-3 bg-gray-50 rounded-lg mb-4">
                    <div className="flex items-center gap-2">
                      <Gavel className="w-4 h-4 text-gray-500" />
                      <span className="text-sm text-gray-600">
                        {product.highestBidder?.username ? (
                          <>
                            <span className="font-medium text-gray-900">{product.highestBidder.username}</span>
                            <span className="text-gray-500"> đang dẫn đầu</span>
                          </>
                        ) : (
                          <span className="text-gray-500">Chưa có lượt đấu giá</span>
                        )}
                      </span>
                    </div>
                  </div>

                  {/* Time Info */}
                  <div
                    className={`flex items-center gap-2 py-2 px-3 rounded-lg mb-4 ${
                      isActive 
                        ? 'bg-green-50 text-green-700' 
                        : 'bg-red-50 text-red-700'
                    }`}
                  >
                    <Clock className="w-4 h-4" />
                    <span className="text-sm font-medium">
                      {isActive
                        ? `Còn ${formatRelativeTime(product.endAt)}`
                        : `Kết thúc ${formatRelativeTime(product.endAt)}`}
                    </span>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex gap-2">
                    <Link
                      to={`/products/${product.id}`}
                      className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors text-sm font-medium"
                    >
                      <Eye className="w-4 h-4" />
                      Xem
                    </Link>
                    <Link
                      to={`/seller/products/${product.id}/edit`}
                      className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
                    >
                      <Edit className="w-4 h-4" />
                      Sửa
                    </Link>
                    <button
                      onClick={() => handleDeleteClick(product)}
                      className="flex items-center justify-center gap-2 px-4 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors text-sm font-medium"
                      title="Xóa sản phẩm"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {deleteModalOpen && productToDelete && (
        <Modal
          isOpen={deleteModalOpen}
          onClose={() => setDeleteModalOpen(false)}
          title="Xác nhận xóa sản phẩm"
        >
          <div className="space-y-4">
            <p className="text-gray-600">
              Bạn có chắc chắn muốn xóa sản phẩm <strong>{productToDelete.name}</strong>?
            </p>
            <p className="text-sm text-red-600">
              Hành động này không thể hoàn tác.
            </p>
            <div className="flex gap-3 pt-4">
              <button
                onClick={() => setDeleteModalOpen(false)}
                disabled={isDeleting}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 disabled:opacity-50"
              >
                Hủy
              </button>
              <button
                onClick={handleDeleteConfirm}
                disabled={isDeleting}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {isDeleting ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    Đang xóa...
                  </>
                ) : (
                  'Xóa sản phẩm'
                )}
              </button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
};

export default MyProductsPage;
