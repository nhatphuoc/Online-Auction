import { useParams } from 'react-router-dom';
import { LoadingSpinner } from '../../components/Common/Loading';

const ProductListPage = () => {
  const { categoryId } = useParams();

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Danh sách sản phẩm - Category {categoryId}</h1>
        <LoadingSpinner />
        <p className="text-center text-gray-600 mt-4">Tính năng này sắp sửa có...</p>
      </div>
    </div>
  );
};

export default ProductListPage;
