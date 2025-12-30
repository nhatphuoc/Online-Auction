import { useState, useEffect } from 'react';
import { Category } from '../../types';
import { categoryService } from '../../services/category.service';
import { ChevronDown } from 'lucide-react';

interface CategorySelectorProps {
  value?: {
    categoryId: number;
    categoryName: string;
    parentCategoryId: number;
    parentCategoryName: string;
  };
  onChange: (category: {
    categoryId: number;
    categoryName: string;
    parentCategoryId: number;
    parentCategoryName: string;
  }) => void;
  error?: string;
}

export const CategorySelector = ({
  value,
  onChange,
  error,
}: CategorySelectorProps) => {
  const [parentCategories, setParentCategories] = useState<Category[]>([]);
  const [childCategories, setChildCategories] = useState<Category[]>([]);
  const [selectedParentId, setSelectedParentId] = useState<number | undefined>(
    value?.parentCategoryId
  );
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadParentCategories();
  }, []);

  useEffect(() => {
    if (selectedParentId) {
      loadChildCategories(selectedParentId);
    } else {
      setChildCategories([]);
    }
  }, [selectedParentId]);

  const loadParentCategories = async () => {
    try {
      setIsLoading(true);
      const categories = await categoryService.getAllCategories({ level: 1 });
      setParentCategories(categories);
    } catch (error) {
      console.error('Error loading parent categories:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadChildCategories = async (parentId: number) => {
    try {
      const children = await categoryService.getCategoriesByParent(parentId);
      setChildCategories(children);
    } catch (error) {
      console.error('Error loading child categories:', error);
    }
  };

  const handleParentChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const parentId = parseInt(e.target.value);
    const parent = parentCategories.find((c) => c.id === parentId);
    
    if (parent) {
      setSelectedParentId(parentId);
      // Reset child selection when parent changes
      setChildCategories([]);
    } else {
      setSelectedParentId(undefined);
      setChildCategories([]);
    }
  };

  const handleChildChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const childId = parseInt(e.target.value);
    const child = childCategories.find((c) => c.id === childId);
    const parent = parentCategories.find((c) => c.id === selectedParentId);

    if (child && parent) {
      onChange({
        categoryId: child.id,
        categoryName: child.name,
        parentCategoryId: parent.id,
        parentCategoryName: parent.name,
      });
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="animate-pulse">
          <div className="h-4 w-32 bg-gray-200 rounded mb-2"></div>
          <div className="h-10 bg-gray-200 rounded"></div>
        </div>
        <div className="animate-pulse">
          <div className="h-4 w-32 bg-gray-200 rounded mb-2"></div>
          <div className="h-10 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Parent Category */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Danh mục chính <span className="text-red-500">*</span>
        </label>
        <div className="relative">
          <select
            value={selectedParentId || ''}
            onChange={handleParentChange}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent appearance-none bg-white"
          >
            <option value="">-- Chọn danh mục chính --</option>
            {parentCategories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400 pointer-events-none" />
        </div>
      </div>

      {/* Child Category */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Danh mục con <span className="text-red-500">*</span>
        </label>
        <div className="relative">
          <select
            value={value?.categoryId || ''}
            onChange={handleChildChange}
            disabled={!selectedParentId || childCategories.length === 0}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent appearance-none bg-white disabled:bg-gray-100 disabled:cursor-not-allowed"
          >
            <option value="">
              {selectedParentId
                ? childCategories.length === 0
                  ? '-- Không có danh mục con --'
                  : '-- Chọn danh mục con --'
                : '-- Vui lòng chọn danh mục chính trước --'}
            </option>
            {childCategories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400 pointer-events-none" />
        </div>
        {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
      </div>
    </div>
  );
};
