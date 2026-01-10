import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../../hooks/useAuth';
import { productService } from '../../../services/product.service';
import { mediaService } from '../../../services/media.service';
import { CategorySelector } from '../../../components/Product/CategorySelector';
import { RichTextEditor } from '../../../components/Product/RichTextEditor';
import { ArrowLeft, Loader2, DollarSign, Calendar, Image as ImageIcon } from 'lucide-react';
import { useUIStore } from '../../../stores/ui.store';

interface ProductFormData {
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
}

export const CreateProductPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const addToast = useUIStore((state) => state.addToast);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [uploadingImages, setUploadingImages] = useState(false);

  const [formData, setFormData] = useState<Partial<ProductFormData>>({
    name: '',
    description: '',
    images: [],
    thumbnailUrl: '',
    startingPrice: 0,
    stepPrice: 0,
    buyNowPrice: undefined,
    autoExtend: true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name?.trim()) {
      newErrors.name = 'Vui lòng nhập tên sản phẩm';
    }

    if (!formData.thumbnailUrl) {
      newErrors.thumbnailUrl = 'Vui lòng chọn ảnh đại diện';
    }

    if (!formData.images || formData.images.length < 3) {
      newErrors.images = 'Vui lòng tải lên ít nhất 3 ảnh sản phẩm';
    }

    if (!formData.description?.trim()) {
      newErrors.description = 'Vui lòng nhập mô tả sản phẩm';
    }

    if (!formData.categoryId) {
      newErrors.category = 'Vui lòng chọn danh mục sản phẩm';
    }

    if (!formData.startingPrice || formData.startingPrice <= 0) {
      newErrors.startingPrice = 'Giá khởi điểm phải lớn hơn 0';
    }

    if (!formData.stepPrice || formData.stepPrice <= 0) {
      newErrors.stepPrice = 'Bước giá phải lớn hơn 0';
    }

    if (formData.buyNowPrice && formData.buyNowPrice <= (formData.startingPrice || 0)) {
      newErrors.buyNowPrice = 'Giá mua ngay phải lớn hơn giá khởi điểm';
    }

    if (!formData.endAt) {
      newErrors.endAt = 'Vui lòng chọn thời gian kết thúc';
    } else {
      const endTime = new Date(formData.endAt).getTime();
      const now = Date.now();
      if (endTime <= now) {
        newErrors.endAt = 'Thời gian kết thúc phải sau thời điểm hiện tại';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleImageUpload = async (files: File[]) => {
    try {
      setUploadingImages(true);

      // Upload using presigned URL flow
      const results = await mediaService.uploadMultipleFiles(files, 'products');

      // Extract successful uploads
      const uploadedUrls = results
        .filter((r) => r.imageUrl !== null)
        .map((r) => r.imageUrl!);

      // Check for failures
      const failedCount = results.filter((r) => r.imageUrl === null).length;
      if (failedCount > 0) {
        addToast('warning', `${failedCount} ảnh tải lên thất bại`);
      }

      if (uploadedUrls.length === 0) {
        addToast('error', 'Không có ảnh nào được tải lên thành công');
        return;
      }

      const allImages = [...(formData.images || []), ...uploadedUrls];
      setFormData({
        ...formData,
        images: allImages,
        thumbnailUrl: formData.thumbnailUrl || allImages[0],
      });

      addToast('success', `Đã tải lên ${uploadedUrls.length} ảnh`);
    } catch (error) {
      console.error('Error uploading images:', error);
      addToast('error', 'Lỗi khi tải ảnh lên');
    } finally {
      setUploadingImages(false);
    }
  };

  const handleRemoveImage = (index: number) => {
    const newImages = formData.images?.filter((_, i) => i !== index) || [];
    const updatedFormData = {
      ...formData,
      images: newImages,
    };

    // If removed image was thumbnail, set new thumbnail
    if (formData.images?.[index] === formData.thumbnailUrl) {
      updatedFormData.thumbnailUrl = newImages[0] || '';
    }

    setFormData(updatedFormData);
  };

  const handleSetThumbnail = (url: string) => {
    setFormData({ ...formData, thumbnailUrl: url });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      addToast('error', 'Vui lòng kiểm tra lại thông tin');
      return;
    }

    try {
      setIsSubmitting(true);

      const productData: ProductFormData = {
        name: formData.name!,
        thumbnailUrl: formData.thumbnailUrl!,
        images: formData.images!,
        description: formData.description!,
        categoryId: formData.categoryId!,
        categoryName: formData.categoryName!,
        parentCategoryId: formData.parentCategoryId!,
        parentCategoryName: formData.parentCategoryName!,
        startingPrice: formData.startingPrice!,
        buyNowPrice: formData.buyNowPrice,
        stepPrice: formData.stepPrice!,
        endAt: formData.endAt!,
        autoExtend: formData.autoExtend!,
      };

      await productService.createProduct(productData);
      addToast('success', 'Tạo sản phẩm thành công!');
      navigate('/seller/products');
    } catch (error) {
      console.error('Error creating product:', error);
      const errorMessage = error instanceof Error ? error.message : 'Lỗi khi tạo sản phẩm';
      addToast('error', errorMessage);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!user || user.userRole !== 'ROLE_SELLER') {
    return (
      <div className="max-w-7xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold text-gray-900">Không có quyền truy cập</h1>
        <p className="text-gray-600 mt-2">Chỉ người bán mới có thể tạo sản phẩm</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <button
          onClick={() => navigate('/seller/products')}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-5 h-5" />
          Quay lại danh sách sản phẩm
        </button>
        <h1 className="text-3xl font-bold text-gray-900">Đăng sản phẩm mới</h1>
        <p className="text-gray-600 mt-2">Tạo phiên đấu giá mới cho sản phẩm của bạn</p>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-8">
        {/* Product Name */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Tên sản phẩm <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            value={formData.name}
            onChange={(e) =>
              setFormData({ ...formData, name: e.target.value })
            }
            className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.name ? 'border-red-500' : 'border-gray-300'
              }`}
            placeholder="Nhập tên sản phẩm"
          />
          {errors.name && (
            <p className="mt-1 text-sm text-red-600">{errors.name}</p>
          )}
        </div>

        {/* Images */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Hình ảnh sản phẩm <span className="text-red-500">*</span>
            <span className="text-sm text-gray-500 ml-2">(Tối thiểu 3 ảnh)</span>
          </label>

          {/* Image Grid */}
          {formData.images && formData.images.length > 0 && (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
              {formData.images.map((url, index) => (
                <div key={index} className="relative group">
                  <img
                    src={url}
                    alt={`Product ${index + 1}`}
                    className={`w-full h-32 object-cover rounded-lg ${url === formData.thumbnailUrl
                      ? 'ring-4 ring-blue-500'
                      : ''
                      }`}
                  />
                  {url === formData.thumbnailUrl && (
                    <div className="absolute top-2 left-2 bg-blue-500 text-white text-xs px-2 py-1 rounded">
                      Ảnh đại diện
                    </div>
                  )}
                  <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 transition-all rounded-lg flex items-center justify-center gap-2">
                    {url !== formData.thumbnailUrl && (
                      <button
                        type="button"
                        onClick={() => handleSetThumbnail(url)}
                        className="opacity-0 group-hover:opacity-100 px-3 py-1 bg-blue-500 text-white text-sm rounded hover:bg-blue-600"
                      >
                        Đặt làm ảnh đại diện
                      </button>
                    )}
                    <button
                      type="button"
                      onClick={() => handleRemoveImage(index)}
                      className="opacity-0 group-hover:opacity-100 px-3 py-1 bg-red-500 text-white text-sm rounded hover:bg-red-600"
                    >
                      Xóa
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Upload Button */}
          <div>
            <input
              type="file"
              multiple
              accept="image/*"
              onChange={(e) => {
                const files = Array.from(e.target.files || []);
                handleImageUpload(files);
              }}
              disabled={uploadingImages}
              className="hidden"
              id="image-upload"
            />
            <label
              htmlFor="image-upload"
              className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-500 transition-colors cursor-pointer block"
            >
              {uploadingImages ? (
                <Loader2 className="w-12 h-12 text-gray-400 mx-auto mb-4 animate-spin" />
              ) : (
                <ImageIcon className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              )}
              <p className="text-gray-600">
                {uploadingImages
                  ? 'Đang tải ảnh lên...'
                  : 'Click để chọn ảnh'}
              </p>
              <p className="text-sm text-gray-500 mt-2">
                PNG, JPG, JPEG - Tối đa 5MB
              </p>
            </label>
          </div>
          {errors.images && (
            <p className="mt-1 text-sm text-red-600">{errors.images}</p>
          )}
        </div>

        {/* Category */}
        <CategorySelector
          value={
            formData.categoryId
              ? {
                categoryId: formData.categoryId,
                categoryName: formData.categoryName!,
                parentCategoryId: formData.parentCategoryId!,
                parentCategoryName: formData.parentCategoryName!,
              }
              : undefined
          }
          onChange={(category) => {
            setFormData({
              ...formData,
              categoryId: category.categoryId,
              categoryName: category.categoryName,
              parentCategoryId: category.parentCategoryId,
              parentCategoryName: category.parentCategoryName,
            });
          }}
          error={errors.category}
        />

        {/* Description */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Mô tả sản phẩm <span className="text-red-500">*</span>
          </label>
          <RichTextEditor
            value={formData.description || ''}
            onChange={(value) =>
              setFormData({ ...formData, description: value })
            }
            error={errors.description}
          />
        </div>

        {/* Pricing */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Giá khởi điểm <span className="text-red-500">*</span>
            </label>
            <div className="relative">
              <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={formData.startingPrice
                  ? new Intl.NumberFormat('vi-VN').format(formData.startingPrice)
                  : ''
                }
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    startingPrice: Number(e.target.value.replace(/[^\d]/g, '')) || 0,
                  })
                }
                className={`w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.startingPrice ? 'border-red-500' : 'border-gray-300'
                  }`}
                placeholder="0"
              />
            </div>
            {errors.startingPrice && (
              <p className="mt-1 text-sm text-red-600">{errors.startingPrice}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Bước giá <span className="text-red-500">*</span>
            </label>
            <div className="relative">
              <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={
                  formData.stepPrice
                    ? new Intl.NumberFormat('vi-VN').format(formData.stepPrice)
                    : ''
                }
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    stepPrice: Number(e.target.value.replace(/[^\d]/g, '')) || 0,
                  })
                }
                className={`w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.stepPrice ? 'border-red-500' : 'border-gray-300'
                  }`}
                placeholder="0"
              />
            </div>
            {errors.stepPrice && (
              <p className="mt-1 text-sm text-red-600">{errors.stepPrice}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Giá mua ngay (tùy chọn)
            </label>
            <div className="relative">
              <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                value={
                  formData.buyNowPrice !== undefined
                    ? new Intl.NumberFormat('vi-VN').format(formData.buyNowPrice)
                    : ''
                }
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    buyNowPrice: e.target.value
                      ? Number(e.target.value.replace(/[^\d]/g, ''))
                      : undefined,
                  })
                }
                className={`w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.buyNowPrice ? 'border-red-500' : 'border-gray-300'
                  }`}
                placeholder="0"
              />
            </div>
            {errors.buyNowPrice && (
              <p className="mt-1 text-sm text-red-600">{errors.buyNowPrice}</p>
            )}
          </div>
        </div>

        {/* End Date */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Thời gian kết thúc <span className="text-red-500">*</span>
          </label>
          <div className="relative">
            <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="datetime-local"
              value={formData.endAt || ''}
              onChange={(e) =>
                setFormData({ ...formData, endAt: e.target.value })
              }
              className={`w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent ${errors.endAt ? 'border-red-500' : 'border-gray-300'
                }`}
            />
          </div>
          {errors.endAt && (
            <p className="mt-1 text-sm text-red-600">{errors.endAt}</p>
          )}
        </div>

        {/* Auto Extend */}
        <div className="flex items-start">
          <input
            type="checkbox"
            id="autoExtend"
            checked={formData.autoExtend}
            onChange={(e) =>
              setFormData({ ...formData, autoExtend: e.target.checked })
            }
            className="mt-1 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label htmlFor="autoExtend" className="ml-3">
            <span className="block text-sm font-medium text-gray-700">
              Tự động gia hạn
            </span>
            <span className="block text-sm text-gray-500">
              Tự động gia hạn thêm 10 phút nếu có lượt đấu giá mới trong 5 phút
              cuối
            </span>
          </label>
        </div>

        {/* Submit Buttons */}
        <div className="flex gap-4 pt-6 border-t">
          <button
            type="button"
            onClick={() => navigate('/seller/products')}
            className="flex-1 px-6 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Hủy
          </button>
          <button
            type="submit"
            disabled={isSubmitting || uploadingImages}
            className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {isSubmitting ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Đang tạo...
              </>
            ) : (
              'Đăng sản phẩm'
            )}
          </button>
        </div>
      </form>
    </div>
  );
};
