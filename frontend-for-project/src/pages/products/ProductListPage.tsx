import { useEffect, useState, useCallback } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { productService } from '../../services/product.service';
import { ProductListItem } from '../../types';
import { ProductGrid } from '../../components/UI/ProductCard';
import { Pagination } from '../../components/UI/Pagination';
import { CategoryBreadcrumb } from '../../components/Category/CategoryMenu';
import { ArrowUpDown } from 'lucide-react';

const ProductListPage = () => {
  const { categoryId } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  
  const [products, setProducts] = useState<ProductListItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [sortBy, setSortBy] = useState<string>(
    searchParams.get('sort') || 'endAt'
  );

  const loadProducts = useCallback(async () => {
    if (!categoryId) return;

    setIsLoading(true);
    try {
      const response = await productService.searchProducts({
        categoryId: parseInt(categoryId),
        page: currentPage,
        pageSize: 12,
      });

      if (response.success && response.data) {
        setProducts(response.data.content);
        setTotalPages(response.data.totalPages);
      }
    } catch (error) {
      console.error('Failed to load products:', error);
      setProducts([]);
    } finally {
      setIsLoading(false);
    }
  }, [categoryId, currentPage]);

  useEffect(() => {
    if (categoryId) {
      loadProducts();
    }
  }, [categoryId, currentPage, sortBy, loadProducts]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleSortChange = (newSort: string) => {
    setSortBy(newSort);
    setCurrentPage(0);
    setSearchParams({ sort: newSort });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Breadcrumb */}
        {categoryId && <CategoryBreadcrumb categoryId={parseInt(categoryId)} />}

        {/* Header with Sort */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              Danh Sách Sản Phẩm
            </h1>
            {!isLoading && (
              <p className="text-gray-600 mt-1">
                Tìm thấy {products.length} sản phẩm
              </p>
            )}
          </div>

          {/* Sort Options */}
          <div className="flex items-center gap-2">
            <ArrowUpDown className="w-5 h-5 text-gray-400" />
            <select
              value={sortBy}
              onChange={(e) => handleSortChange(e.target.value)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="endAt">Sắp kết thúc</option>
              <option value="price_asc">Giá tăng dần</option>
              <option value="price_desc">Giá giảm dần</option>
              <option value="bids">Nhiều lượt đấu</option>
              <option value="newest">Mới nhất</option>
            </select>
          </div>
        </div>

        {/* Products Grid */}
        <ProductGrid
          products={products}
          isLoading={isLoading}
          emptyMessage="Không có sản phẩm nào trong danh mục này"
        />

        {/* Pagination */}
        {!isLoading && totalPages > 1 && (
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
            className="mt-8"
          />
        )}
      </div>
    </div>
  );
};

export default ProductListPage;
