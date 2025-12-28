import { useState, useRef } from 'react';
import { mediaService } from '../../services/media.service';
import { useAuthStore } from '../../stores/auth.store';
import { useUIStore } from '../../stores/ui.store';
import { Upload, X, Image as ImageIcon, Loader } from 'lucide-react';

interface ImageUploaderProps {
  maxFiles?: number;
  maxSizeInMB?: number;
  folder?: string;
  onUploadComplete?: (imageUrls: string[]) => void;
  initialImages?: string[];
}

interface UploadingFile {
  file: File;
  preview: string;
  status: 'pending' | 'uploading' | 'success' | 'error';
  progress: number;
  url?: string;
  error?: string;
}

export const ImageUploader = ({
  maxFiles = 5,
  maxSizeInMB = 5,
  folder = 'products',
  onUploadComplete,
  initialImages = [],
}: ImageUploaderProps) => {
  const { isAuthenticated } = useAuthStore();
  const addToast = useUIStore((state) => state.addToast);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [uploadingFiles, setUploadingFiles] = useState<UploadingFile[]>([]);
  const [uploadedUrls, setUploadedUrls] = useState<string[]>(initialImages);

  const validateFile = (file: File): string | null => {
    // Check file type
    if (!file.type.startsWith('image/')) {
      return 'Chỉ chấp nhận file ảnh (JPG, PNG, GIF, WebP)';
    }

    // Check file size
    const maxSize = maxSizeInMB * 1024 * 1024;
    if (file.size > maxSize) {
      return `Kích thước file không được vượt quá ${maxSizeInMB}MB`;
    }

    return null;
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (!isAuthenticated) {
      addToast('error', 'Vui lòng đăng nhập để upload ảnh');
      return;
    }

    const files = Array.from(event.target.files || []);
    
    // Check total files limit
    const currentTotal = uploadingFiles.length + uploadedUrls.length;
    const remainingSlots = maxFiles - currentTotal;
    
    if (files.length > remainingSlots) {
      addToast('error', `Chỉ có thể upload tối đa ${maxFiles} ảnh. Còn lại ${remainingSlots} ảnh.`);
      return;
    }

    // Validate each file
    const validFiles: File[] = [];
    for (const file of files) {
      const error = validateFile(file);
      if (error) {
        addToast('error', `${file.name}: ${error}`);
      } else {
        validFiles.push(file);
      }
    }

    if (validFiles.length === 0) return;

    // Create preview for valid files
    const newUploadingFiles: UploadingFile[] = validFiles.map((file) => ({
      file,
      preview: URL.createObjectURL(file),
      status: 'pending',
      progress: 0,
    }));

    setUploadingFiles((prev) => [...prev, ...newUploadingFiles]);

    // Start upload
    uploadFiles(newUploadingFiles);

    // Reset input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const uploadFiles = async (filesToUpload: UploadingFile[]) => {
    for (let i = 0; i < filesToUpload.length; i++) {
      const uploadingFile = filesToUpload[i];
      
      try {
        // Update status to uploading
        updateFileStatus(uploadingFile.preview, 'uploading', 10);

        // Step 1: Get presigned URL
        updateFileStatus(uploadingFile.preview, 'uploading', 30);
        const presignedData = await mediaService.getPresignedUrl(
          uploadingFile.file.name,
          folder
        );

        // Step 2: Upload to S3 using presigned URL
        updateFileStatus(uploadingFile.preview, 'uploading', 60);
        await mediaService.uploadToPresignedUrl(
          presignedData.presigned_url,
          uploadingFile.file
        );

        // Step 3: Success
        updateFileStatus(uploadingFile.preview, 'success', 100, presignedData.image_url);
        
        // Add to uploaded URLs
        setUploadedUrls((prev) => {
          const newUrls = [...prev, presignedData.image_url];
          if (onUploadComplete) {
            onUploadComplete(newUrls);
          }
          return newUrls;
        });

      } catch (error) {
        console.error('Upload failed:', error);
        const errorMessage = error instanceof Error ? error.message : 'Upload thất bại';
        updateFileStatus(uploadingFile.preview, 'error', 0, undefined, errorMessage);
        addToast('error', `${uploadingFile.file.name}: ${errorMessage}`);
      }
    }
  };

  const updateFileStatus = (
    preview: string,
    status: UploadingFile['status'],
    progress: number,
    url?: string,
    error?: string
  ) => {
    setUploadingFiles((prev) =>
      prev.map((file) =>
        file.preview === preview
          ? { ...file, status, progress, url, error }
          : file
      )
    );
  };

  const removeUploadingFile = (preview: string) => {
    setUploadingFiles((prev) => {
      const file = prev.find((f) => f.preview === preview);
      if (file) {
        URL.revokeObjectURL(file.preview);
      }
      return prev.filter((f) => f.preview !== preview);
    });
  };

  const removeUploadedImage = (url: string) => {
    setUploadedUrls((prev) => {
      const newUrls = prev.filter((u) => u !== url);
      if (onUploadComplete) {
        onUploadComplete(newUrls);
      }
      return newUrls;
    });
  };

  const triggerFileInput = () => {
    fileInputRef.current?.click();
  };

  const canUploadMore = uploadingFiles.length + uploadedUrls.length < maxFiles;

  return (
    <div className="space-y-4">
      {/* Upload Button */}
      {canUploadMore && (
        <div
          onClick={triggerFileInput}
          className="border-2 border-dashed border-gray-300 hover:border-blue-500 rounded-lg p-8 text-center cursor-pointer transition-colors bg-gray-50 hover:bg-blue-50"
        >
          <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <p className="text-sm font-medium text-gray-700 mb-1">
            Nhấn để chọn ảnh
          </p>
          <p className="text-xs text-gray-500">
            Tối đa {maxFiles} ảnh, mỗi ảnh không quá {maxSizeInMB}MB
          </p>
          <p className="text-xs text-gray-400 mt-1">
            ({uploadedUrls.length + uploadingFiles.length}/{maxFiles} đã chọn)
          </p>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            onChange={handleFileSelect}
            className="hidden"
          />
        </div>
      )}

      {/* Uploaded Images */}
      {uploadedUrls.length > 0 && (
        <div>
          <h4 className="text-sm font-semibold text-gray-700 mb-2">
            Ảnh đã upload ({uploadedUrls.length})
          </h4>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {uploadedUrls.map((url, index) => (
              <div
                key={url}
                className="relative aspect-square rounded-lg overflow-hidden border-2 border-green-200 bg-gray-100 group"
              >
                <img
                  src={url}
                  alt={`Uploaded ${index + 1}`}
                  className="w-full h-full object-cover"
                  onError={(e) => {
                    e.currentTarget.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%23e5e7eb" width="100" height="100"/%3E%3C/svg%3E';
                  }}
                />
                <button
                  onClick={() => removeUploadedImage(url)}
                  className="absolute top-2 right-2 p-1 bg-red-600 hover:bg-red-700 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                  title="Xóa ảnh"
                >
                  <X className="w-4 h-4" />
                </button>
                {index === 0 && (
                  <div className="absolute bottom-0 left-0 right-0 bg-green-600 text-white text-xs py-1 px-2 text-center">
                    Ảnh đại diện
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Uploading Files */}
      {uploadingFiles.length > 0 && (
        <div>
          <h4 className="text-sm font-semibold text-gray-700 mb-2">
            Đang upload ({uploadingFiles.filter(f => f.status === 'uploading' || f.status === 'pending').length})
          </h4>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {uploadingFiles.map((uploadFile) => (
              <div
                key={uploadFile.preview}
                className="relative aspect-square rounded-lg overflow-hidden border-2 border-gray-300 bg-gray-100"
              >
                <img
                  src={uploadFile.preview}
                  alt={uploadFile.file.name}
                  className="w-full h-full object-cover"
                />
                
                {/* Upload Overlay */}
                <div className="absolute inset-0 bg-black bg-opacity-50 flex flex-col items-center justify-center">
                  {uploadFile.status === 'pending' && (
                    <ImageIcon className="w-8 h-8 text-white opacity-70" />
                  )}
                  
                  {uploadFile.status === 'uploading' && (
                    <>
                      <Loader className="w-8 h-8 text-white animate-spin mb-2" />
                      <span className="text-white text-sm font-medium">
                        {uploadFile.progress}%
                      </span>
                      <div className="w-4/5 h-2 bg-gray-700 rounded-full mt-2 overflow-hidden">
                        <div
                          className="h-full bg-blue-500 transition-all duration-300"
                          style={{ width: `${uploadFile.progress}%` }}
                        />
                      </div>
                    </>
                  )}
                  
                  {uploadFile.status === 'success' && (
                    <div className="text-white text-center">
                      <div className="w-12 h-12 mx-auto mb-2 bg-green-500 rounded-full flex items-center justify-center">
                        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                      </div>
                      <span className="text-xs">Hoàn thành</span>
                    </div>
                  )}
                  
                  {uploadFile.status === 'error' && (
                    <div className="text-white text-center px-2">
                      <div className="w-12 h-12 mx-auto mb-2 bg-red-500 rounded-full flex items-center justify-center">
                        <X className="w-8 h-8" />
                      </div>
                      <span className="text-xs">
                        {uploadFile.error || 'Lỗi upload'}
                      </span>
                    </div>
                  )}
                </div>

                {/* Remove Button */}
                {(uploadFile.status === 'error' || uploadFile.status === 'success') && (
                  <button
                    onClick={() => removeUploadingFile(uploadFile.preview)}
                    className="absolute top-2 right-2 p-1 bg-red-600 hover:bg-red-700 text-white rounded-full"
                    title="Xóa"
                  >
                    <X className="w-4 h-4" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Info */}
      {uploadedUrls.length === 0 && uploadingFiles.length === 0 && (
        <div className="text-center py-4 text-gray-500 text-sm">
          <ImageIcon className="w-12 h-12 mx-auto mb-2 opacity-30" />
          <p>Chưa có ảnh nào được chọn</p>
        </div>
      )}
    </div>
  );
};

export default ImageUploader;
