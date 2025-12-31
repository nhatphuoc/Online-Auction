import { useState } from 'react';
import { UpgradeRequest } from '../../types';
import { CheckCircle, XCircle, User, Clock, MessageSquare } from 'lucide-react';

interface UpgradeRequestListProps {
  requests: UpgradeRequest[];
  onApprove: (requestId: number) => void;
  onReject: (requestId: number, reason?: string) => void;
}

export default function UpgradeRequestList({ requests, onApprove, onReject }: UpgradeRequestListProps) {
  const [rejectingId, setRejectingId] = useState<number | null>(null);
  const [rejectReason, setRejectReason] = useState('');

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('vi-VN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const handleRejectClick = (requestId: number) => {
    setRejectingId(requestId);
    setRejectReason('');
  };

  const handleRejectSubmit = (requestId: number) => {
    onReject(requestId, rejectReason);
    setRejectingId(null);
    setRejectReason('');
  };

  const handleRejectCancel = () => {
    setRejectingId(null);
    setRejectReason('');
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'PENDING':
        return <span className="px-3 py-1 text-xs font-medium rounded-full bg-yellow-100 text-yellow-800">Chờ Duyệt</span>;
      case 'APPROVED':
        return <span className="px-3 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800">Đã Duyệt</span>;
      case 'REJECTED':
        return <span className="px-3 py-1 text-xs font-medium rounded-full bg-red-100 text-red-800">Đã Từ Chối</span>;
      default:
        return <span className="px-3 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-800">{status}</span>;
    }
  };

  if (requests.length === 0) {
    return (
      <div className="bg-white rounded-lg shadow-md p-12 text-center text-gray-500">
        <User className="w-16 h-16 mx-auto mb-4 text-gray-300" />
        <p className="text-lg">Không có yêu cầu nào</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 border-b">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Người Dùng
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Lý Do
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Trạng Thái
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Ngày Gửi
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Ngày Xử Lý
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Thao Tác
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {requests.map((request) => (
              <>
                <tr key={request.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div className="flex items-center">
                      <div className="flex-shrink-0 h-10 w-10 bg-gray-200 rounded-full flex items-center justify-center">
                        <User className="w-6 h-6 text-gray-600" />
                      </div>
                      <div className="ml-4">
                        <div className="text-sm font-medium text-gray-900">
                          {request.user.fullName}
                        </div>
                        <div className="text-sm text-gray-500">
                          {request.user.email}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <div className="text-sm text-gray-900 max-w-xs truncate" title={request.reason}>
                      {request.reason}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {getStatusBadge(request.status)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center text-sm text-gray-500">
                      <Clock className="w-4 h-4 mr-1" />
                      {formatDate(request.createdAt)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {request.reviewedAt ? formatDate(request.reviewedAt) : '-'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    {request.status === 'PENDING' && (
                      <div className="flex justify-end gap-2">
                        <button
                          onClick={() => onApprove(request.id)}
                          className="inline-flex items-center px-3 py-1.5 bg-green-600 text-white rounded-md hover:bg-green-700 transition"
                          title="Duyệt yêu cầu"
                        >
                          <CheckCircle className="w-4 h-4 mr-1" />
                          Duyệt
                        </button>
                        <button
                          onClick={() => handleRejectClick(request.id)}
                          className="inline-flex items-center px-3 py-1.5 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
                          title="Từ chối yêu cầu"
                        >
                          <XCircle className="w-4 h-4 mr-1" />
                          Từ Chối
                        </button>
                      </div>
                    )}
                  </td>
                </tr>

                {/* Rejection Reason Row */}
                {request.status === 'REJECTED' && request.rejectionReason && (
                  <tr className="bg-red-50">
                    <td colSpan={6} className="px-6 py-3">
                      <div className="flex items-start">
                        <MessageSquare className="w-4 h-4 text-red-600 mr-2 mt-0.5" />
                        <div>
                          <span className="text-sm font-medium text-red-900">Lý do từ chối:</span>
                          <p className="text-sm text-red-800 mt-1">{request.rejectionReason}</p>
                        </div>
                      </div>
                    </td>
                  </tr>
                )}

                {/* Reject Form Row */}
                {rejectingId === request.id && (
                  <tr className="bg-gray-50">
                    <td colSpan={6} className="px-6 py-4">
                      <div className="max-w-2xl">
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Lý do từ chối (tùy chọn)
                        </label>
                        <textarea
                          value={rejectReason}
                          onChange={(e) => setRejectReason(e.target.value)}
                          rows={3}
                          className="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent"
                          placeholder="Nhập lý do từ chối yêu cầu..."
                        />
                        <div className="flex gap-2 mt-3">
                          <button
                            onClick={() => handleRejectSubmit(request.id)}
                            className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
                          >
                            Xác Nhận Từ Chối
                          </button>
                          <button
                            onClick={handleRejectCancel}
                            className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition"
                          >
                            Hủy
                          </button>
                        </div>
                      </div>
                    </td>
                  </tr>
                )}
              </>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
