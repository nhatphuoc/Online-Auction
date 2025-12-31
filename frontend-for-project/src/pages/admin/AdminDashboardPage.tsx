import { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { useUIStore } from '../../stores/ui.store';
import { userService } from '../../services/user.service';
import Loading from '../../components/Common/Loading';
import { Users, Package, Clock, CheckCircle, XCircle, TrendingUp } from 'lucide-react';

interface DashboardStats {
  totalUsers: number;
  totalProducts: number;
  pendingUpgrades: number;
  approvedUpgrades: number;
  rejectedUpgrades: number;
}

export default function AdminDashboardPage() {
  const { user } = useAuth();
  const addToast = useUIStore((state) => state.addToast);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    totalProducts: 0,
    pendingUpgrades: 0,
    approvedUpgrades: 0,
    rejectedUpgrades: 0,
  });

  useEffect(() => {
    if (user?.userRole !== 'ROLE_ADMIN') {
      addToast('error', 'Bạn không có quyền truy cập trang này');
      window.location.href = '/';
      return;
    }
    fetchStats();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user]);

  const fetchStats = async () => {
    try {
      setLoading(true);
      
      // Fetch upgrade requests stats
      const [pendingRes, approvedRes, rejectedRes] = await Promise.all([
        userService.getUpgradeRequests({ status: 'PENDING', size: 1 }),
        userService.getUpgradeRequests({ status: 'APPROVED', size: 1 }),
        userService.getUpgradeRequests({ status: 'REJECTED', size: 1 }),
      ]);

      setStats({
        totalUsers: 0, // We don't have a specific endpoint for this
        totalProducts: 0, // We don't have a specific endpoint for this
        pendingUpgrades: pendingRes.totalElements,
        approvedUpgrades: approvedRes.totalElements,
        rejectedUpgrades: rejectedRes.totalElements,
      });
    } catch (error) {
      addToast('error', 'Không thể tải thống kê');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <Loading />;
  }

  const statCards = [
    {
      title: 'Yêu Cầu Chờ Duyệt',
      value: stats.pendingUpgrades,
      icon: Clock,
      color: 'bg-yellow-500',
      textColor: 'text-yellow-600',
      bgColor: 'bg-yellow-50',
      link: '/admin/upgrade-requests',
    },
    {
      title: 'Đã Duyệt',
      value: stats.approvedUpgrades,
      icon: CheckCircle,
      color: 'bg-green-500',
      textColor: 'text-green-600',
      bgColor: 'bg-green-50',
    },
    {
      title: 'Đã Từ Chối',
      value: stats.rejectedUpgrades,
      icon: XCircle,
      color: 'bg-red-500',
      textColor: 'text-red-600',
      bgColor: 'bg-red-50',
    },
    {
      title: 'Tổng Yêu Cầu',
      value: stats.pendingUpgrades + stats.approvedUpgrades + stats.rejectedUpgrades,
      icon: TrendingUp,
      color: 'bg-blue-500',
      textColor: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Trang Quản Trị</h1>
        <p className="text-gray-600 mt-2">Tổng quan hệ thống đấu giá trực tuyến</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {statCards.map((stat, index) => (
          <div
            key={index}
            className={`${stat.bgColor} rounded-lg shadow-md p-6 transition-transform hover:scale-105 ${
              stat.link ? 'cursor-pointer' : ''
            }`}
            onClick={() => stat.link && (window.location.href = stat.link)}
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600 mb-1">{stat.title}</p>
                <p className={`text-3xl font-bold ${stat.textColor}`}>{stat.value}</p>
              </div>
              <div className={`${stat.color} p-3 rounded-full`}>
                <stat.icon className="w-6 h-6 text-white" />
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4 flex items-center">
            <Users className="w-5 h-5 mr-2 text-blue-600" />
            Quản Lý Người Dùng
          </h2>
          <p className="text-gray-600 mb-4">
            Xem và quản lý yêu cầu nâng cấp tài khoản từ Bidder lên Seller
          </p>
          <a
            href="/admin/upgrade-requests"
            className="inline-block px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            Xem Yêu Cầu
          </a>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4 flex items-center">
            <Package className="w-5 h-5 mr-2 text-purple-600" />
            Quản Lý Danh Mục
          </h2>
          <p className="text-gray-600 mb-4">
            Thêm, sửa, xóa danh mục sản phẩm trong hệ thống
          </p>
          <a
            href="/admin/categories"
            className="inline-block px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition"
          >
            Quản Lý Danh Mục
          </a>
        </div>
      </div>

      {/* Recent Activity - Placeholder */}
      <div className="mt-8 bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold mb-4">Hoạt Động Gần Đây</h2>
        <div className="text-center py-8 text-gray-500">
          <p>Chức năng đang phát triển...</p>
        </div>
      </div>
    </div>
  );
}
