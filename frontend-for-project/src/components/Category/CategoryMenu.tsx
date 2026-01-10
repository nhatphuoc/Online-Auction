import { useEffect, useState, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { ChevronDown, ChevronRight, Layers } from 'lucide-react';
import { useCategoryStore } from '../../stores/category.store';
import { categoryService } from '../../services/category.service';

export const CategoryMenu = () => {
  const [expandedParents, setExpandedParents] = useState<Set<number>>(new Set());
  const [isLoading, setIsLoading] = useState(true);
  
  const {
    setCategories,
    getParentCategories,
    getChildCategories,
  } = useCategoryStore();

  const loadCategories = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await categoryService.getAllCategories();
      setCategories(data);
    } catch (error) {
      console.error('Failed to load categories:', error);
    } finally {
      setIsLoading(false);
    }
  }, [setCategories]);

  useEffect(() => {
    loadCategories();
  }, [loadCategories]);

  const toggleParent = (parentId: number) => {
    const newExpanded = new Set(expandedParents);
    if (newExpanded.has(parentId)) {
      newExpanded.delete(parentId);
    } else {
      newExpanded.add(parentId);
    }
    setExpandedParents(newExpanded);
  };

  const parentCategories = getParentCategories();

  if (isLoading) {
    return (
      <div className="bg-white rounded-lg shadow-sm p-4">
        <div className="flex items-center gap-2 mb-4">
          <Layers className="w-5 h-5 text-gray-400 animate-pulse" />
          <div className="h-5 bg-gray-200 rounded w-32 animate-pulse" />
        </div>
        <div className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-8 bg-gray-100 rounded animate-pulse" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-sm overflow-hidden">
      <div className="bg-gradient-to-r from-blue-600 to-blue-700 text-white p-4">
        <div className="flex items-center gap-2">
          <Layers className="w-5 h-5" />
          <h2 className="font-bold text-lg">Danh Mục</h2>
        </div>
      </div>

      <nav className="p-2">
        {parentCategories.length === 0 ? (
          <div className="text-center text-gray-500 py-8">
            Không có danh mục nào
          </div>
        ) : (
          <ul className="space-y-1">
            {/* Mục "Tất cả" */}
            <li>
              <Link
                to="/category"
                className="block px-3 py-2 rounded-lg hover:bg-blue-50 hover:text-blue-600 transition-colors font-medium"
              >
                Tất cả sản phẩm
              </Link>
            </li>
            
            {/* Divider */}
            <li className="my-2">
              <hr className="border-gray-200" />
            </li>

            {parentCategories.map((parent) => {
              const children = getChildCategories(parent.id);
              const isExpanded = expandedParents.has(parent.id);

              return (
                <li key={parent.id}>
                  <div className="flex items-center">
                    {children.length > 0 && (
                      <button
                        onClick={() => toggleParent(parent.id)}
                        className="p-1 hover:bg-gray-100 rounded transition-colors"
                      >
                        {isExpanded ? (
                          <ChevronDown className="w-4 h-4 text-gray-600" />
                        ) : (
                          <ChevronRight className="w-4 h-4 text-gray-600" />
                        )}
                      </button>
                    )}
                    <Link
                      to={`/category/${parent.id}`}
                      className="flex-1 px-3 py-2 rounded-lg hover:bg-blue-50 hover:text-blue-600 transition-colors font-medium"
                    >
                      {parent.name}
                    </Link>
                  </div>

                  {isExpanded && children.length > 0 && (
                    <ul className="ml-6 mt-1 space-y-1 border-l-2 border-gray-200 pl-2">
                      {children.map((child) => (
                        <li key={child.id}>
                          <Link
                            to={`/category/${child.id}`}
                            className="block px-3 py-2 rounded-lg hover:bg-blue-50 hover:text-blue-600 transition-colors text-sm"
                          >
                            {child.name}
                          </Link>
                        </li>
                      ))}
                    </ul>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </nav>
    </div>
  );
};

export const CategoryBreadcrumb = ({ categoryId }: { categoryId: number }) => {
  const { getCategoryById } = useCategoryStore();
  const category = getCategoryById(categoryId);

  if (!category) return null;

  const parent = category.parent_id ? getCategoryById(category.parent_id) : null;

  return (
    <nav className="flex items-center gap-2 text-sm text-gray-600 mb-4">
      <Link to="/" className="hover:text-blue-600 transition-colors">
        Trang chủ
      </Link>
      <ChevronRight className="w-4 h-4" />
      {parent && (
        <>
          <Link
            to={`/category/${parent.id}`}
            className="hover:text-blue-600 transition-colors"
          >
            {parent.name}
          </Link>
          <ChevronRight className="w-4 h-4" />
        </>
      )}
      <span className="text-gray-900 font-medium">{category.name}</span>
    </nav>
  );
};
