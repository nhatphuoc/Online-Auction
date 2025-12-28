import { useSearchParams } from 'react-router-dom';
import { LoadingSpinner } from '../../components/Common/Loading';

const SearchPage = () => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') || '';

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Kết quả tìm kiếm</h1>
        <p className="text-gray-600 mb-8">Tìm kiếm: "{query}"</p>
        <LoadingSpinner />
        <p className="text-center text-gray-600 mt-4">Tính năng này sắp sửa có...</p>
      </div>
    </div>
  );
};

export default SearchPage;
