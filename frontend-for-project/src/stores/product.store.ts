import { create } from 'zustand';
import { Product, ProductListItem } from '../types';

interface ProductState {
  currentProduct: Product | null;
  products: ProductListItem[];
  selectedCategoryId: number | null;
  searchQuery: string;
  currentPage: number;
  pageSize: number;
  totalPages: number;
  sortBy: 'endAt' | 'price' | 'bids';

  setCurrentProduct: (product: Product | null) => void;
  setProducts: (products: ProductListItem[]) => void;
  setSelectedCategoryId: (categoryId: number | null) => void;
  setSearchQuery: (query: string) => void;
  setCurrentPage: (page: number) => void;
  setPageSize: (size: number) => void;
  setTotalPages: (pages: number) => void;
  setSortBy: (sort: 'endAt' | 'price' | 'bids') => void;
  reset: () => void;
}

export const useProductStore = create<ProductState>((set) => ({
  currentProduct: null,
  products: [],
  selectedCategoryId: null,
  searchQuery: '',
  currentPage: 0,
  pageSize: 12,
  totalPages: 0,
  sortBy: 'endAt',

  setCurrentProduct: (product) => set({ currentProduct: product }),
  setProducts: (products) => set({ products }),
  setSelectedCategoryId: (categoryId) => set({ selectedCategoryId: categoryId }),
  setSearchQuery: (query) => set({ searchQuery: query, currentPage: 0 }),
  setCurrentPage: (page) => set({ currentPage: page }),
  setPageSize: (size) => set({ pageSize: size }),
  setTotalPages: (pages) => set({ totalPages: pages }),
  setSortBy: (sort) => set({ sortBy: sort, currentPage: 0 }),
  reset: () =>
    set({
      currentProduct: null,
      products: [],
      selectedCategoryId: null,
      searchQuery: '',
      currentPage: 0,
      pageSize: 12,
      totalPages: 0,
      sortBy: 'endAt',
    }),
}));
