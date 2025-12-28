import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import {
  UserProfile,
  UpgradeRequest,
  PaginationResponse,
  ApiResponse,
} from '../types';

export const userService = {
  async getUserProfile() {
    const response = await apiClient.get<ApiResponse<UserProfile>>(
      endpoints.users.profile
    );
    return response.data.data;
  },

  async getUserByEmail(email: string) {
    const response = await apiClient.get<ApiResponse<UserProfile>>(
      endpoints.users.simple(email)
    );
    return response.data.data;
  },

  async getUserById(id: number) {
    const response = await apiClient.get<ApiResponse<UserProfile>>(
      endpoints.users.simpleById(id)
    );
    return response.data.data;
  },

  async searchUsers(params: {
    keyword?: string;
    role?: 'ROLE_BIDDER' | 'ROLE_SELLER' | 'ROLE_ADMIN';
    page?: number;
    size?: number;
  }) {
    const queryParams = new URLSearchParams();
    if (params.keyword) queryParams.append('keyword', params.keyword);
    if (params.role) queryParams.append('role', params.role);
    queryParams.append('page', (params.page || 0).toString());
    queryParams.append('size', (params.size || 10).toString());

    const response = await apiClient.get<PaginationResponse<UserProfile>>(
      `${endpoints.users.search}?${queryParams.toString()}`
    );
    return response.data;
  },

  async requestUpgradeToSeller(reason: string) {
    const response = await apiClient.post<string>(
      endpoints.users.upgradeToSeller(reason)
    );
    return response.data;
  },

  async approveUpgradeRequest(requestId: number) {
    const response = await apiClient.post<string>(
      endpoints.users.approveUpgrade(requestId)
    );
    return response.data;
  },

  async getUpgradeRequests(params: {
    page?: number;
    size?: number;
    sort?: string;
    direction?: 'asc' | 'desc';
  }) {
    const queryParams = new URLSearchParams();
    queryParams.append('page', (params.page || 0).toString());
    queryParams.append('size', (params.size || 10).toString());
    queryParams.append('sort', params.sort || 'createdAt');
    queryParams.append('direction', params.direction || 'desc');

    const response = await apiClient.get<PaginationResponse<UpgradeRequest>>(
      `${endpoints.users.upgradeRequests}?${queryParams.toString()}`
    );
    return response.data;
  },

  async updateProfile(data: {
    fullName?: string;
    phoneNumber?: string;
    address?: string;
    dateOfBirth?: string;
  }) {
    const response = await apiClient.put<ApiResponse<UserProfile>>(
      endpoints.users.profile,
      data
    );
    return response.data.data;
  },
};
