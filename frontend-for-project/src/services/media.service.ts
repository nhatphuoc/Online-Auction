import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import { PresignedUrlResponse } from '../types';

/**
 * Media Service - Sử dụng Presigned URL để upload trực tiếp lên S3
 * 
 * Flow:
 * 1. Gọi getPresignedUrl hoặc getPresignedUrls để lấy presigned URL từ backend
 * 2. Backend tạo presigned URL từ AWS S3 (có thời hạn 15 phút)
 * 3. Sử dụng uploadToPresignedUrl để upload file trực tiếp lên S3
 * 4. Sử dụng image_url (public URL) để lưu vào database
 */

export const mediaService = {
  /**
   * Lấy presigned URL cho 1 file
   * @param filename - Tên file (VD: product.jpg)
   * @param folder - Thư mục đích trong S3 (VD: products/, avatars/)
   * @returns Presigned URL, image URL, key, expires_in
   */
  async getPresignedUrl(filename: string, folder?: string) {
    const queryParams = new URLSearchParams();
    queryParams.append('filename', filename);
    if (folder) queryParams.append('folder', folder);

    const response = await apiClient.get<PresignedUrlResponse>(
      `${endpoints.media.presign}?${queryParams.toString()}`
    );
    return response.data;
  },

  /**
   * Lấy presigned URLs cho nhiều files
   * @param filenames - Mảng tên files
   * @param folder - Thư mục đích trong S3
   * @returns Array of presigned URL objects
   */
  async getPresignedUrls(filenames: string[], folder?: string) {
    const queryParams = folder ? `?folder=${encodeURIComponent(folder)}` : '';
    
    const response = await apiClient.post<{
      presigned: PresignedUrlResponse[];
    }>(`${endpoints.media.presignMultiple}${queryParams}`, filenames);
    
    return response.data.presigned;
  },

  /**
   * Upload file trực tiếp lên S3 sử dụng presigned URL
   * @param presignedUrl - Presigned URL từ getPresignedUrl
   * @param file - File object cần upload
   * @returns true nếu upload thành công
   */
  async uploadToPresignedUrl(presignedUrl: string, file: File): Promise<boolean> {
    try {
      const response = await fetch(presignedUrl, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': file.type,
        },
      });
      return response.ok;
    } catch (error) {
      console.error('Error uploading to S3:', error);
      return false;
    }
  },

  /**
   * Upload 1 file (Helper function - tổng hợp 2 bước)
   * @param file - File object cần upload
   * @param folder - Thư mục đích trong S3
   * @returns image_url (public URL) nếu thành công, null nếu failed
   */
  async uploadSingleFile(file: File, folder?: string): Promise<string | null> {
    try {
      // Bước 1: Lấy presigned URL
      const presignedData = await this.getPresignedUrl(file.name, folder);
      
      // Bước 2: Upload lên S3
      const success = await this.uploadToPresignedUrl(presignedData.presigned_url, file);
      
      if (success) {
        return presignedData.image_url;
      }
      return null;
    } catch (error) {
      console.error('Error in uploadSingleFile:', error);
      return null;
    }
  },

  /**
   * Upload nhiều files song song (Helper function)
   * @param files - Array of File objects
   * @param folder - Thư mục đích trong S3
   * @returns Array of image URLs (null for failed uploads)
   */
  async uploadMultipleFiles(
    files: File[],
    folder?: string
  ): Promise<{ imageUrl: string | null; filename: string }[]> {
    try {
      // Bước 1: Lấy presigned URLs cho tất cả files
      const filenames = files.map((f) => f.name);
      const presignedUrls = await this.getPresignedUrls(filenames, folder);

      // Bước 2: Upload tất cả files song song
      const uploadPromises = presignedUrls.map((presignedData, index) =>
        this.uploadToPresignedUrl(presignedData.presigned_url, files[index])
          .then((success) => ({
            imageUrl: success ? presignedData.image_url : null,
            filename: files[index].name,
          }))
          .catch(() => ({
            imageUrl: null,
            filename: files[index].name,
          }))
      );

      return await Promise.all(uploadPromises);
    } catch (error) {
      console.error('Error in uploadMultipleFiles:', error);
      return files.map((f) => ({ imageUrl: null, filename: f.name }));
    }
  },

  /**
   * Upload nhiều files với progress tracking
   * @param files - Array of File objects
   * @param folder - Thư mục đích
   * @param onProgress - Callback function nhận (uploaded, total)
   * @returns Array of image URLs
   */
  async uploadMultipleFilesWithProgress(
    files: File[],
    folder: string | undefined,
    onProgress?: (uploaded: number, total: number) => void
  ): Promise<{ imageUrl: string | null; filename: string }[]> {
    try {
      // Bước 1: Lấy presigned URLs
      const filenames = files.map((f) => f.name);
      const presignedUrls = await this.getPresignedUrls(filenames, folder);

      // Bước 2: Upload từng file và track progress
      const results: { imageUrl: string | null; filename: string }[] = [];
      let uploaded = 0;

      for (let i = 0; i < files.length; i++) {
        const success = await this.uploadToPresignedUrl(
          presignedUrls[i].presigned_url,
          files[i]
        );

        results.push({
          imageUrl: success ? presignedUrls[i].image_url : null,
          filename: files[i].name,
        });

        uploaded++;
        if (onProgress) {
          onProgress(uploaded, files.length);
        }
      }

      return results;
    } catch (error) {
      console.error('Error in uploadMultipleFilesWithProgress:', error);
      return files.map((f) => ({ imageUrl: null, filename: f.name }));
    }
  },
};
