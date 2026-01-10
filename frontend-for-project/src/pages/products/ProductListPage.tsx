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
    setIsLoading(true);
    try {
      const params: {
        page: number;
        pageSize: number;
        categoryId?: number;
      } = {
        page: currentPage,
        pageSize: 12,
      };

      // If categoryId exists, filter by category
      if (categoryId) {
        params.categoryId = parseInt(categoryId);
      }

      const response = await productService.searchProducts(params);

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
    loadProducts();
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
        {categoryId ? (
          <CategoryBreadcrumb categoryId={parseInt(categoryId)} />
        ) : (
          <nav className="flex mb-6" aria-label="Breadcrumb">
            <ol className="inline-flex items-center space-x-1 md:space-x-3">
              <li className="inline-flex items-center">
                <a href="/" className="text-gray-700 hover:text-blue-600 inline-flex items-center">
                  <svg className="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"></path>
                  </svg>
                  Trang chủ
                </a>
              </li>
              <li>
                <div className="flex items-center">
                  <svg className="w-6 h-6 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd"></path>
                  </svg>
                  <span className="ml-1 text-gray-500 md:ml-2">Tất cả sản phẩm</span>
                </div>
              </li>
            </ol>
          </nav>
        )}

        {/* Header with Sort */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {categoryId ? 'Danh Sách Sản Phẩm' : 'Tất Cả Sản Phẩm'}
            </h1>
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
