import { useEffect, useState, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { productService } from '../../services/product.service';
import { categoryService } from '../../services/category.service';
import { ProductListItem, Category } from '../../types';
import { ProductGrid } from '../../components/UI/ProductCard';
import { Pagination } from '../../components/UI/Pagination';
import { Search, Filter, X } from 'lucide-react';

const SearchPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  
  const [products, setProducts] = useState<ProductListItem[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [totalElements, setTotalElements] = useState(0);
  
  const [searchQuery, setSearchQuery] = useState(query);
  const [selectedParentCategory, setSelectedParentCategory] = useState<number | null>(null);
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null);
  const [sortBy, setSortBy] = useState('endAt');
  const [showFilters, setShowFilters] = useState(false);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      const data = await categoryService.getAllCategories();
      setCategories(data);
    } catch (error) {
      console.error('Failed to load categories:', error);
    }
  };

  const searchProducts = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await productService.searchProducts({
        query: query || undefined,
        parentCategoryId: selectedParentCategory || undefined,
        categoryId: selectedCategory || undefined,
        page: currentPage,
        pageSize: 12,
      });

      if (response.success && response.data) {
        setProducts(response.data.content);
        setTotalPages(response.data.totalPages);
        setTotalElements(response.data.totalElements);
      }
    } catch (error) {
      console.error('Search failed:', error);
      setProducts([]);
    } finally {
      setIsLoading(false);
    }
  }, [query, selectedParentCategory, selectedCategory, currentPage]);

  useEffect(() => {
    searchProducts();
  }, [searchProducts]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      setSearchParams({ q: searchQuery.trim() });
      setCurrentPage(0);
    }
  };

  const handleClearFilters = () => {
    setSelectedParentCategory(null);
    setSelectedCategory(null);
    setCurrentPage(0);
  };

  const parentCategories = categories.filter(cat => cat.level === 1);
  const childCategories = selectedParentCategory
    ? categories.filter(cat => cat.level === 2 && cat.parent_id === selectedParentCategory)
    : [];

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Search Header */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
          <form onSubmit={handleSearch} className="flex gap-3">
            <div className="relative flex-1">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Tìm kiếm sản phẩm..."
                className="w-full pl-12 pr-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <button
              type="submit"
              className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
            >
              Tìm kiếm
            </button>
            <button
              type="button"
              onClick={() => setShowFilters(!showFilters)}
              className="lg:hidden px-4 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <Filter className="w-5 h-5" />
            </button>
          </form>
        </div>

        <div className="flex gap-6">
          {/* Filters Sidebar */}
          <aside className={`w-64 space-y-4 ${showFilters ? 'block' : 'hidden lg:block'}`}>
            <div className="bg-white rounded-lg shadow-sm p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="font-bold text-gray-900">Bộ Lọc</h3>
                {(selectedParentCategory || selectedCategory) && (
                  <button
                    onClick={handleClearFilters}
                    className="text-sm text-blue-600 hover:text-blue-700 flex items-center gap-1"
                  >
                    <X className="w-4 h-4" />
                    Xóa
                  </button>
                )}
              </div>

              {/* Parent Category Filter */}
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Danh Mục Chính
                </label>
                <select
                  value={selectedParentCategory || ''}
                  onChange={(e) => {
                    setSelectedParentCategory(e.target.value ? parseInt(e.target.value) : null);
                    setSelectedCategory(null);
                    setCurrentPage(0);
                  }}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">Tất cả</option>
                  {parentCategories.map(cat => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>

              {/* Child Category Filter */}
              {selectedParentCategory && childCategories.length > 0 && (
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Danh Mục Con
                  </label>
                  <select
                    value={selectedCategory || ''}
                    onChange={(e) => {
                      setSelectedCategory(e.target.value ? parseInt(e.target.value) : null);
                      setCurrentPage(0);
                    }}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">Tất cả</option>
                    {childCategories.map(cat => (
                      <option key={cat.id} value={cat.id}>{cat.name}</option>
                    ))}
                  </select>
                </div>
              )}

              {/* Sort */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Sắp Xếp
                </label>
                <select
                  value={sortBy}
                  onChange={(e) => {
                    setSortBy(e.target.value);
                    setCurrentPage(0);
                  }}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="endAt">Sắp kết thúc</option>
                  <option value="price_asc">Giá tăng dần</option>
                  <option value="price_desc">Giá giảm dần</option>
                  <option value="bids">Nhiều lượt đấu</option>
                  <option value="newest">Mới nhất</option>
                </select>
              </div>
            </div>
          </aside>

          {/* Results */}
          <div className="flex-1">
            <div className="mb-6">
              <h2 className="text-2xl font-bold text-gray-900">
                Kết Quả Tìm Kiếm
                {query && <span className="text-blue-600"> "{query}"</span>}
              </h2>
              {!isLoading && (
                <p className="text-gray-600 mt-1">
                  Tìm thấy {totalElements} sản phẩm
                </p>
              )}
            </div>

            <ProductGrid
              products={products}
              isLoading={isLoading}
              emptyMessage={
                query
                  ? `Không tìm thấy sản phẩm nào cho "${query}"`
                  : 'Nhập từ khóa để tìm kiếm sản phẩm'
              }
            />

            {!isLoading && totalPages > 1 && (
              <Pagination
                currentPage={currentPage}
                totalPages={totalPages}
                onPageChange={(page) => {
                  setCurrentPage(page);
                  window.scrollTo({ top: 0, behavior: 'smooth' });
                }}
                className="mt-8"
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SearchPage;
