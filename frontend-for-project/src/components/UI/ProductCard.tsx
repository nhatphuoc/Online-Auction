import { Link } from 'react-router-dom';
import { ProductListItem } from '../../types';
import { formatCurrency, formatRelativeTime } from '../../utils/formatters';
import { Clock, Gavel, Eye } from 'lucide-react';
import { ProductImage } from '../Common/Image';

interface ProductCardProps {
  product: ProductListItem;
  showCategory?: boolean;
  isNew?: boolean;
}

export const ProductCard = ({ product, showCategory = true, isNew = false }: ProductCardProps) => {
  const timeLeft = new Date(product.endAt).getTime() - Date.now();
  const isEndingSoon = timeLeft < 3 * 24 * 60 * 60 * 1000; // Less than 3 days

  return (
    <Link
      to={`/products/${product.id}`}
      className="bg-white rounded-lg shadow-sm hover:shadow-lg transition-all duration-200 overflow-hidden group relative"
    >
      {isNew && (
        <div className="absolute top-2 left-2 z-10 bg-red-500 text-white text-xs font-bold px-2 py-1 rounded-md">
          MỚI
        </div>
      )}
      
      <div className="relative h-48 bg-gray-200 overflow-hidden">
        <ProductImage
          src={product.thumbnailUrl}
          alt={product.name}
          className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
        />
        {product.buyNowPrice && (
          <div className="absolute top-2 right-2 bg-green-500 text-white text-xs font-semibold px-2 py-1 rounded-md">
            MUA NGAY
          </div>
        )}
      </div>

      <div className="p-4">
        <h3 className="font-semibold text-gray-900 line-clamp-2 mb-2 group-hover:text-blue-600 min-h-[3rem]">
          {product.name}
        </h3>

        {showCategory && (
          <div className="text-xs text-gray-500 mb-2">
            {product.categoryParentName} › {product.categoryName}
          </div>
        )}

        <div className="space-y-2">
          <div className="flex justify-between items-baseline">
            <span className="text-sm text-gray-600">Giá hiện tại:</span>
            <span className="font-bold text-lg text-blue-600">
              {formatCurrency(product.currentPrice)}
            </span>
          </div>

          {product.buyNowPrice && (
            <div className="flex justify-between items-baseline">
              <span className="text-xs text-gray-500">Mua ngay:</span>
              <span className="font-semibold text-sm text-green-600">
                {formatCurrency(product.buyNowPrice)}
              </span>
            </div>
          )}

          <div className="flex items-center justify-between text-xs text-gray-500 pt-2 border-t">
            <div className="flex items-center gap-1">
              <Gavel className="w-4 h-4" />
              <span>{product.bidCount} lượt</span>
            </div>
            <div className={`flex items-center gap-1 ${isEndingSoon ? 'text-red-600 font-semibold' : ''}`}>
              <Clock className="w-4 h-4" />
              <span>{formatRelativeTime(product.endAt)}</span>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
};

interface ProductGridProps {
  products: ProductListItem[];
  isLoading?: boolean;
  emptyMessage?: string;
}

export const ProductGrid = ({ products, isLoading, emptyMessage = 'Không có sản phẩm nào' }: ProductGridProps) => {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {Array.from({ length: 8 }).map((_, i) => (
          <ProductCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="text-center py-12">
        <Eye className="w-16 h-16 text-gray-300 mx-auto mb-4" />
        <p className="text-gray-500 text-lg">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      {products.map((product) => {
        const minutesSinceCreated =
          (Date.now() - new Date(product.createdAt).getTime()) / (1000 * 60);
        const isNew = minutesSinceCreated < 60; // Product created within last 60 minutes

        return <ProductCard key={product.id} product={product} isNew={isNew} />;
      })}
    </div>
  );
};

const ProductCardSkeleton = () => (
  <div className="bg-white rounded-lg shadow-sm overflow-hidden animate-pulse">
    <div className="h-48 bg-gray-200" />
    <div className="p-4">
      <div className="h-4 bg-gray-200 rounded mb-2" />
      <div className="h-4 bg-gray-200 rounded w-3/4 mb-4" />
      <div className="space-y-2">
        <div className="h-3 bg-gray-200 rounded" />
        <div className="h-3 bg-gray-200 rounded" />
        <div className="flex justify-between mt-4">
          <div className="h-3 bg-gray-200 rounded w-1/3" />
          <div className="h-3 bg-gray-200 rounded w-1/4" />
        </div>
      </div>
    </div>
  </div>
);
