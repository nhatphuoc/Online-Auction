import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../../../hooks/useAuth';
import { productService } from '../../../services/product.service';
import { RichTextEditor } from '../../../components/Product/RichTextEditor';
import { ArrowLeft, Loader2, Save } from 'lucide-react';
import { useUIStore } from '../../../stores/ui.store';
import { Product } from '../../../types';

export const EditProductPage = () => {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();
  const addToast = useUIStore((state) => state.addToast);

  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [additionalDescription, setAdditionalDescription] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      loadProduct(parseInt(id));
    }
  }, [id]);

  const loadProduct = async (productId: number) => {
    try {
      setIsLoading(true);
      const data = await productService.getProductDetail(productId);
      console.log(data);
      setProduct(data);

      // Check if user is the seller
      if (user && data.sellerInfo.id !== user.id) {
        addToast('error', 'Bạn không có quyền chỉnh sửa sản phẩm này');
        navigate('/seller/products');
      }
    } catch (error) {
      console.error('Error loading product:', error);
      addToast('error', 'Không thể tải thông tin sản phẩm');
      navigate('/seller/products');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!additionalDescription.trim()) {
      setError('Vui lòng nhập nội dung mô tả bổ sung');
      return;
    }

    if (!product) return;

    try {
      setIsSubmitting(true);
      const updated = await productService.updateDescription(
        product.id,
        additionalDescription
      );

      setProduct(updated);
      setAdditionalDescription('');
      setError('');
      addToast('success', 'Đã cập nhật mô tả sản phẩm');
    } catch (error) {
      console.error('Error updating product:', error);
      const errorMessage = error instanceof Error ? error.message : 'Lỗi khi cập nhật mô tả';
      addToast('error', errorMessage);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!user || user.userRole !== 'ROLE_SELLER') {
    return (
      <div className="max-w-7xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold text-gray-900">Không có quyền truy cập</h1>
        <p className="text-gray-600 mt-2">Chỉ người bán mới có thể chỉnh sửa sản phẩm</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-16 text-center">
        <Loader2 className="w-12 h-12 text-blue-600 mx-auto mb-4 animate-spin" />
        <p className="text-gray-600">Đang tải thông tin sản phẩm...</p>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold text-gray-900">Không tìm thấy sản phẩm</h1>
        <button
          onClick={() => navigate('/seller/products')}
          className="mt-4 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          Quay lại danh sách sản phẩm
        </button>
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
        <h1 className="text-3xl font-bold text-gray-900">Chỉnh sửa mô tả sản phẩm</h1>
        <p className="text-gray-600 mt-2">{product.name}</p>
      </div>

      {/* Product Info */}
      <div className="bg-gray-50 rounded-lg p-6 mb-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <img
              src={product.thumbnailUrl}
              alt={product.name}
              className="w-full h-48 object-cover rounded-lg"
            />
          </div>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-gray-600">Danh mục</p>
              <p className="font-medium">{product.parentCategoryName} → {product.categoryName}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Giá hiện tại</p>
              <p className="font-bold text-blue-600">
                {product.currentPrice.toLocaleString('vi-VN')} đ
              </p>
            </div>
            {product.buyNowPrice && (
              <div>
                <p className="text-sm text-gray-600">Giá mua ngay</p>
                <p className="font-semibold text-green-600">
                  {product.buyNowPrice.toLocaleString('vi-VN')} đ
                </p>
              </div>
            )}
            <div>
              <p className="text-sm text-gray-600">Thời gian kết thúc</p>
              <p className="font-medium">
                {new Date(product.endAt).toLocaleString('vi-VN')}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Current Description */}
      <div className="mb-8">
        <h2 className="text-xl font-bold text-gray-900 mb-4">Mô tả hiện tại</h2>
        <div
          className="prose max-w-none bg-white border border-gray-200 rounded-lg p-6"
          dangerouslySetInnerHTML={{
            __html: product.description.replace(/\n/g, '<br />')
          }}
        />
      </div>

      {/* Add Description Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <h2 className="text-xl font-bold text-gray-900 mb-4">
            Thêm mô tả bổ sung
          </h2>
          <p className="text-sm text-gray-600 mb-4">
            Nội dung mới sẽ được thêm vào cuối mô tả hiện tại. Bạn không thể
            chỉnh sửa hoặc xóa mô tả cũ.
          </p>
          <RichTextEditor
            value={additionalDescription}
            onChange={(value) => {
              setAdditionalDescription(value);
              setError('');
            }}
            placeholder="Nhập nội dung bổ sung..."
            error={error}
          />
        </div>

        {/* Preview */}
        {additionalDescription && (
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-3">
              Xem trước nội dung bổ sung
            </h3>
            <div
              className="prose max-w-none bg-gray-50 border border-gray-200 rounded-lg p-6"
              dangerouslySetInnerHTML={{
                __html: additionalDescription
                  .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
                  .replace(/\*(.*?)\*/g, '<em>$1</em>')
                  .replace(/^### (.*$)/gim, '<h3 class="text-lg font-bold mt-4 mb-2">$1</h3>')
                  .replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-4 mb-2">$1</h2>')
                  .replace(/^# (.*$)/gim, '<h2 class="text-2xl font-bold mt-4 mb-2">$1</h2>')
                  .replace(/\n/g, '<br />')
              }}
            />
          </div>
        )}

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
            disabled={isSubmitting || !additionalDescription.trim()}
            className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {isSubmitting ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                Đang lưu...
              </>
            ) : (
              <>
                <Save className="w-5 h-5" />
                Lưu mô tả bổ sung
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
};
