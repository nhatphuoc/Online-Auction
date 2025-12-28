import { useEffect } from 'react';
import { useAsync } from '../../hooks/useAsync';
import { apiClient } from '../../services/api/client';
import { endpoints } from '../../services/api/endpoints';
import { ProductListItem, ApiResponse } from '../../types';
import { LoadingSpinner, ProductSkeleton } from '../../components/Common/Loading';
import { formatCurrency, formatTimeRemaining, formatRelativeTime } from '../../utils/formatters';
import { TrendingUp, Clock, DollarSign } from 'lucide-react';

const HomePage = () => {
  const topEnding = useAsync<ApiResponse<ProductListItem[]>>(
    () => apiClient.get(endpoints.products.topEnding),
    true,
    []
  );

  const topBids = useAsync<ApiResponse<ProductListItem[]>>(
    () => apiClient.get(endpoints.products.topMostBids),
    true,
    []
  );

  const topPrice = useAsync<ApiResponse<ProductListItem[]>>(
    () => apiClient.get(endpoints.products.topHighestPrice),
    true,
    []
  );

  const ProductSection = ({
    title,
    icon: Icon,
    data,
    isLoading,
  }: {
    title: string;
    icon: any;
    // API responses may vary — accept array or objects that wrap an array (e.g. { content: [...] } or ApiResponse)
    data: any;
    isLoading: boolean;
  }) => {
    // Normalize different possible payload shapes into an array or null
    const list: ProductListItem[] | null = (() => {
      if (!data) return null;
      if (Array.isArray(data)) return data as ProductListItem[];
      if (Array.isArray(data.content)) return data.content as ProductListItem[];
      if (Array.isArray(data.data)) return data.data as ProductListItem[]; // nested data
      if (Array.isArray(data.items)) return data.items as ProductListItem[];
      return null;
    })();

    return (
      <section className="py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3 mb-8">
            <Icon className="w-6 h-6 text-blue-600" />
            <h2 className="text-2xl font-bold text-gray-900">{title}</h2>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-6">
            {isLoading ? (
              Array.from({ length: 5 }).map((_, i) => <ProductSkeleton key={i} />)
            ) : list && list.length > 0 ? (
              list.map((product) => (
                <a
                  key={product.id}
                  href={`/products/${product.id}`}
                  className="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden group"
                >
                  <div className="relative h-48 bg-gray-200 overflow-hidden">
                    <img
                      src={product.thumbnailUrl}
                      alt={product.name}
                      className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                    />
                  </div>
                  <div className="p-4">
                    <h3 className="font-semibold text-gray-900 line-clamp-2 mb-2 group-hover:text-blue-600">
                      {product.name}
                    </h3>
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-600">Giá hiện tại:</span>
                        <span className="font-bold text-blue-600">{formatCurrency(product.currentPrice)}</span>
                      </div>
                      <div className="flex justify-between items-center text-xs text-gray-500">
                        <span>Lượt đấu: {product.bidCount}</span>
                      </div>
                    </div>
                  </div>
                </a>
              ))
            ) : (
              <div className="col-span-full text-center text-gray-500 py-8">Không có sản phẩm nào.</div>
            )}
          </div>
        </div>
      </section>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-blue-600 to-blue-800 text-white py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-4xl md:text-5xl font-bold mb-4">Sàn Đấu Giá Trực Tuyến</h1>
          <p className="text-lg md:text-xl text-blue-100 mb-8">
            Tìm kiếm những sản phẩm tuyệt vời với giá tốt nhất
          </p>
          <form onSubmit={(e) => {
            e.preventDefault();
            const query = (e.target as any).search.value;
            window.location.href = `/search?q=${encodeURIComponent(query)}`;
          }} className="max-w-2xl mx-auto">
            <div className="relative">
              <input
                type="text"
                name="search"
                placeholder="Tìm kiếm sản phẩm..."
                className="w-full px-6 py-3 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900"
              />
              <button
                type="submit"
                className="absolute right-2 top-1/2 -translate-y-1/2 px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium"
              >
                Tìm
              </button>
            </div>
          </form>
        </div>
      </section>

      {/* Top Products Sections */}
      <ProductSection
        title="Sắp Kết Thúc"
        icon={Clock}
        data={topEnding.data?.data || null}
        isLoading={topEnding.isLoading}
      />

      <ProductSection
        title="Nhiều Lượt Đấu Nhất"
        icon={TrendingUp}
        data={topBids.data?.data || null}
        isLoading={topBids.isLoading}
      />

      <ProductSection
        title="Giá Cao Nhất"
        icon={DollarSign}
        data={topPrice.data?.data || null}
        isLoading={topPrice.isLoading}
      />
    </div>
  );
};

export default HomePage;
