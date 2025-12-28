import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import { BidRequest, BidResponse, BidHistory, PaginationResponse } from '../types';

export const bidService = {
  async placeBid(bidData: BidRequest) {
    const response = await apiClient.post<BidResponse>(
      endpoints.bids.place,
      bidData
    );
    return response.data;
  },

  async searchBidHistory(params: {
    productId?: number;
    bidderId?: number;
    status?: 'SUCCESS' | 'FAILED';
    requestId?: string;
    from?: string;
    to?: string;
    page?: number;
    size?: number;
  }) {
    const queryParams = new URLSearchParams();
    if (params.productId)
      queryParams.append('productId', params.productId.toString());
    if (params.bidderId)
      queryParams.append('bidderId', params.bidderId.toString());
    if (params.status) queryParams.append('status', params.status);
    if (params.requestId) queryParams.append('requestId', params.requestId);
    if (params.from) queryParams.append('from', params.from);
    if (params.to) queryParams.append('to', params.to);
    queryParams.append('page', (params.page || 0).toString());
    queryParams.append('size', (params.size || 10).toString());

    const response = await apiClient.get<PaginationResponse<BidHistory>>(
      `${endpoints.bids.search}?${queryParams.toString()}`
    );
    return response.data;
  },
};
