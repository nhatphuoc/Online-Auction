import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import {
  MediaUploadResponse,
  PresignedUrlResponse,
} from '../types';

export const mediaService = {
  async uploadSingleFile(file: File, folder?: string) {
    const formData = new FormData();
    formData.append('file', file);

    const url = folder
      ? `${endpoints.media.upload}?folder=${encodeURIComponent(folder)}`
      : endpoints.media.upload;

    const response = await apiClient.post<MediaUploadResponse>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  async uploadMultipleFiles(files: File[], folder?: string) {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append('files', file);
    });

    const url = folder
      ? `${endpoints.media.uploadMultiple}?folder=${encodeURIComponent(folder)}`
      : endpoints.media.uploadMultiple;

    const response = await apiClient.post<{
      message: string;
      files: MediaUploadResponse[];
      success_count: number;
      failed_count: number;
    }>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  async getPresignedUrl(filename: string, folder?: string) {
    const queryParams = new URLSearchParams();
    queryParams.append('filename', filename);
    if (folder) queryParams.append('folder', folder);

    const response = await apiClient.get<PresignedUrlResponse>(
      `${endpoints.media.presign}?${queryParams.toString()}`
    );
    return response.data;
  },

  async uploadToPresignedUrl(presignedUrl: string, file: File) {
    const response = await fetch(presignedUrl, {
      method: 'PUT',
      body: file,
      headers: {
        'Content-Type': file.type,
      },
    });
    return response.ok;
  },
};
