import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import { Category } from '../types';

export const categoryService = {
  async getAllCategories(params?: { parent_id?: number; level?: 1 | 2 }) {
    const queryParams = new URLSearchParams();
    if (params?.parent_id)
      queryParams.append('parent_id', params.parent_id.toString());
    if (params?.level) queryParams.append('level', params.level.toString());

    const url = queryParams.toString()
      ? `${endpoints.categories.list}?${queryParams.toString()}`
      : endpoints.categories.list;

    const response = await apiClient.get<{ categories: Category[] }>(url);
    return response.data.categories;
  },

  async getCategoryById(id: number) {
    const response = await apiClient.get<Category>(
      endpoints.categories.detail(id)
    );
    return response.data;
  },

  async getCategoriesByParent(parentId: number) {
    const response = await apiClient.get<Category[]>(
      endpoints.categories.byParent(parentId)
    );
    return response.data;
  },

  async createCategory(data: {
    name: string;
    slug: string;
    description?: string;
    parent_id?: number;
    level: 1 | 2;
    display_order: number;
  }) {
    const response = await apiClient.post<Category>(
      endpoints.categories.create,
      data
    );
    return response.data;
  },

  async updateCategory(
    id: number,
    data: {
      name?: string;
      slug?: string;
      description?: string;
      is_active?: boolean;
      display_order?: number;
      parent_id?: number | null;
    }
  ) {
    const response = await apiClient.put<Category>(
      endpoints.categories.update(id),
      data
    );
    return response.data;
  },

  async deleteCategory(id: number) {
    const response = await apiClient.delete<{ message: string }>(
      endpoints.categories.delete(id)
    );
    return response.data;
  },
};
