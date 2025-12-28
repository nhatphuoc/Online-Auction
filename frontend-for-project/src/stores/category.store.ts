import { create } from 'zustand';
import { Category } from '../types';

interface CategoryState {
  categories: Category[];
  selectedParentId: number | null;
  selectedCategoryId: number | null;
  isLoading: boolean;
  
  setCategories: (categories: Category[]) => void;
  setSelectedParentId: (id: number | null) => void;
  setSelectedCategoryId: (id: number | null) => void;
  setLoading: (loading: boolean) => void;
  
  getParentCategories: () => Category[];
  getChildCategories: (parentId: number) => Category[];
  getCategoryById: (id: number) => Category | undefined;
}

export const useCategoryStore = create<CategoryState>((set, get) => ({
  categories: [],
  selectedParentId: null,
  selectedCategoryId: null,
  isLoading: false,

  setCategories: (categories) => set({ categories }),
  setSelectedParentId: (id) => set({ selectedParentId: id }),
  setSelectedCategoryId: (id) => set({ selectedCategoryId: id }),
  setLoading: (loading) => set({ isLoading: loading }),

  getParentCategories: () => {
    const { categories } = get();
    // API returns categories with children already nested
    return categories.filter((cat) => cat.level === 1 && cat.is_active);
  },

  getChildCategories: (parentId) => {
    const { categories } = get();
    // First try to find in the parent's children array
    const parent = categories.find((cat) => cat.id === parentId);
    if (parent && parent.children) {
      return parent.children.filter((child) => child.is_active);
    }
    // Fallback to flat structure
    return categories.filter(
      (cat) => cat.parent_id === parentId && cat.level === 2 && cat.is_active
    );
  },

  getCategoryById: (id) => {
    const { categories } = get();
    return categories.find((cat) => cat.id === id);
  },
}));
