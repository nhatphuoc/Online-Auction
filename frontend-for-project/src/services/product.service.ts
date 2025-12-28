import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import {
  Product,
  ProductListItem,
  ApiResponse,
  SearchResponse,
} from '../types';

export const productService = {
  async getTopEnding() {
    const response = await apiClient.get<ApiResponse<ProductListItem[]>>(
      endpoints.products.topEnding
    );
    return response.data;
  },

  async getTopMostBids() {
    const response = await apiClient.get<ApiResponse<ProductListItem[]>>(
      endpoints.products.topMostBids
    );
    return response.data;
  },

  async getTopHighestPrice() {
    const response = await apiClient.get<ApiResponse<ProductListItem[]>>(
      endpoints.products.topHighestPrice
    );
    return response.data;
  },

  async getProductDetail(id: number) {
    const response = await apiClient.get<Product>(
      endpoints.products.detail(id)
    );
    return response.data;
  },

  async getProductsBySeller(sellerId: number) {
    const response = await apiClient.get<ProductListItem[]>(
      endpoints.products.bySeller(sellerId)
    );
    return response.data;
  },

  async searchProducts(params: {
    query?: string;
    parentCategoryId?: number;
    categoryId?: number;
    page?: number;
    pageSize?: number;
  }) {
    const queryParams = new URLSearchParams();
    if (params.query) queryParams.append('query', params.query);
    if (params.parentCategoryId)
      queryParams.append('parentCategoryId', params.parentCategoryId.toString());
    if (params.categoryId)
      queryParams.append('categoryId', params.categoryId.toString());
    queryParams.append('page', (params.page || 0).toString());
    queryParams.append('pageSize', (params.pageSize || 12).toString());

    const response = await apiClient.get<SearchResponse<ProductListItem>>(
      `${endpoints.products.search}?${queryParams.toString()}`
    );
    return response.data;
  },

  async createProduct(productData: {
    name: string;
    thumbnailUrl: string;
    images: string[];
    description: string;
    categoryId: number;
    categoryName: string;
    parentCategoryId: number;
    parentCategoryName: string;
    startingPrice: number;
    buyNowPrice?: number;
    stepPrice: number;
    endAt: string;
    autoExtend: boolean;
  }) {
    const response = await apiClient.post<Product>(
      endpoints.products.create,
      productData
    );
    return response.data;
  },

  async updateDescription(productId: number, additionalDescription: string) {
    const response = await apiClient.patch<Product>(
      endpoints.products.updateDescription(productId),
      { additionalDescription }
    );
    return response.data;
  },
};
