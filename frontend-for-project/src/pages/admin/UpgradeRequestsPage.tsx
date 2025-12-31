import { useState, useEffect } from 'react';
import { userService } from '../../services/user.service';
import { UpgradeRequest } from '../../types';
import { useAuth } from '../../hooks/useAuth';
import { useUIStore } from '../../stores/ui.store';
import UpgradeRequestList from '../../components/Admin/UpgradeRequestList';
import Loading from '../../components/Common/Loading';
import { Pagination } from '../../components/UI/Pagination';

export default function UpgradeRequestsPage() {
  const { user } = useAuth();
  const addToast = useUIStore((state) => state.addToast);
  const [requests, setRequests] = useState<UpgradeRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<'PENDING' | 'APPROVED' | 'REJECTED' | 'ALL'>('PENDING');
  const [currentPage, setCurrentPage] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [totalElements, setTotalElements] = useState(0);
  const pageSize = 10;

  useEffect(() => {
    if (user?.userRole !== 'ROLE_ADMIN') {
      addToast('error', 'Bạn không có quyền truy cập trang này');
      window.location.href = '/';
      return;
    }
    fetchRequests();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user, filter, currentPage]);

  const fetchRequests = async () => {
    try {
      setLoading(true);
      const response = await userService.getUpgradeRequests({
        status: filter === 'ALL' ? undefined : filter,
        page: currentPage,
        size: pageSize,
        sort: 'createdAt',
        direction: 'desc'
      });
      setRequests(response.content);
      setTotalPages(response.totalPages);
      setTotalElements(response.totalElements);
    } catch (error) {
      addToast('error', 'Không thể tải danh sách yêu cầu');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (requestId: number) => {
    if (!window.confirm('Bạn có chắc chắn muốn duyệt yêu cầu này?')) {
      return;
    }

    try {
      await userService.approveUpgradeRequest(requestId);
      addToast('success', 'Đã duyệt yêu cầu thành công');
      fetchRequests();
    } catch (error) {
      const errorMessage = error instanceof Error && 'response' in error 
        ? (error as unknown as { response?: { data?: { message?: string } } }).response?.data?.message 
        : 'Không thể duyệt yêu cầu';
      addToast('error', errorMessage || 'Không thể duyệt yêu cầu');
      console.error(error);
    }
  };

  const handleReject = async (requestId: number, reason?: string) => {
    try {
      await userService.rejectUpgradeRequest(requestId, reason);
      addToast('success', 'Đã từ chối yêu cầu');
      fetchRequests();
    } catch (error) {
      const errorMessage = error instanceof Error && 'response' in error 
        ? (error as unknown as { response?: { data?: { message?: string } } }).response?.data?.message 
        : 'Không thể từ chối yêu cầu';
      addToast('error', errorMessage || 'Không thể từ chối yêu cầu');
      console.error(error);
    }
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  if (loading && requests.length === 0) {
    return <Loading />;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-4">Quản Lý Yêu Cầu Nâng Cấp</h1>
        
        {/* Filter Tabs */}
        <div className="flex gap-2 border-b">
          <button
            onClick={() => {
              setFilter('PENDING');
              setCurrentPage(0);
            }}
            className={`px-4 py-2 font-medium transition ${
              filter === 'PENDING'
                ? 'border-b-2 border-blue-600 text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Chờ Duyệt
          </button>
          <button
            onClick={() => {
              setFilter('APPROVED');
              setCurrentPage(0);
            }}
            className={`px-4 py-2 font-medium transition ${
              filter === 'APPROVED'
                ? 'border-b-2 border-blue-600 text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Đã Duyệt
          </button>
          <button
            onClick={() => {
              setFilter('REJECTED');
              setCurrentPage(0);
            }}
            className={`px-4 py-2 font-medium transition ${
              filter === 'REJECTED'
                ? 'border-b-2 border-blue-600 text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Đã Từ Chối
          </button>
          <button
            onClick={() => {
              setFilter('ALL');
              setCurrentPage(0);
            }}
            className={`px-4 py-2 font-medium transition ${
              filter === 'ALL'
                ? 'border-b-2 border-blue-600 text-blue-600'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Tất Cả
          </button>
        </div>

        {/* Stats */}
        <div className="mt-4 text-sm text-gray-600">
          Tổng số: <span className="font-semibold">{totalElements}</span> yêu cầu
        </div>
      </div>

      {loading ? (
        <Loading />
      ) : (
        <>
          <UpgradeRequestList
            requests={requests}
            onApprove={handleApprove}
            onReject={handleReject}
          />

          {totalPages > 1 && (
            <div className="mt-6">
              <Pagination
                currentPage={currentPage}
                totalPages={totalPages}
                onPageChange={handlePageChange}
              />
            </div>
          )}
        </>
      )}
    </div>
  );
}
