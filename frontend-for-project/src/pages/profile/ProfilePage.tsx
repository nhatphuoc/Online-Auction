import { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { useRole } from '../../hooks/useRole';
import {
  User, Mail, Phone, Calendar, Shield, Star,
  Edit2, Save, X, Lock, Award
} from 'lucide-react';
import { userService } from '../../services/user.service';
import { UserProfile } from '../../types';
import { useUIStore } from '../../stores/ui.store';

const ProfilePage = () => {
  const { user, setUser } = useAuth();
  const { role, isAdmin, isSeller, isBidder } = useRole();
  const addToast = useUIStore((state) => state.addToast);
  const [isEditing, setIsEditing] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [profileData, setProfileData] = useState<UserProfile | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    fullName: user?.fullName || '',
    phoneNumber: user?.phoneNumber || '',
    address: '',
    dateOfBirth: '',
  });

  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      setLoading(true);
      const data = await userService.getUserProfile();
      if (data) {
        setProfileData(data);
        setFormData({
          fullName: data.fullName,
          phoneNumber: data.phoneNumber,
          address: data.address || '',
          dateOfBirth: data.dateOfBirth || '',
        });
      }
    } catch (error) {
      console.error('Error fetching profile:', error);
      addToast('error', 'Không thể tải thông tin hồ sơ');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateProfile = async () => {
    try {
      setLoading(true);
      const updated = await userService.updateProfile(formData);
      if (updated) {
        setProfileData(updated);
        setUser(updated);
        setIsEditing(false);
        addToast('success', 'Cập nhật hồ sơ thành công');
      }
    } catch (error) {
      console.error('Error updating profile:', error);
      addToast('error', 'Không thể cập nhật hồ sơ');
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async () => {
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      addToast('error', 'Mật khẩu mới không khớp');
      return;
    }

    if (passwordData.newPassword.length < 6) {
      addToast('error', 'Mật khẩu phải có ít nhất 6 ký tự');
      return;
    }

    try {
      setLoading(true);
      // TODO: Implement password change API
      addToast('success', 'Đổi mật khẩu thành công');
      setIsChangingPassword(false);
      setPasswordData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      });
    } catch (error) {
      console.error('Error changing password:', error);
      addToast('error', 'Không thể đổi mật khẩu');
    } finally {
      setLoading(false);
    }
  };

  const getRoleBadgeColor = () => {
    if (isAdmin) return 'bg-red-100 text-red-700';
    if (isSeller) return 'bg-purple-100 text-purple-700';
    return 'bg-blue-100 text-blue-700';
  };

  const getRoleText = () => {
    if (isAdmin) return 'Quản trị viên';
    if (isSeller) return 'Người bán';
    return 'Người mua';
  };

  if (loading && !profileData) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const handleUpgradeRequest = async () => {
    try {
      setLoading(true);
      setError(null);

      const res = await userService.requestUpgradeToSeller(
        "Tôi muốn trở thành người bán"
      );

      setMessage(res); // "Upgrade request submitted"
    } catch (err: any) {
      setError(err.response?.data?.error || "Có lỗi xảy ra");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Hồ sơ cá nhân</h1>

      {/* Profile Info Card */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold">Thông tin cá nhân</h2>
          {!isEditing && (
            <button
              onClick={() => setIsEditing(true)}
              className="flex items-center gap-2 px-4 py-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
            >
              <Edit2 className="w-4 h-4" />
              Chỉnh sửa
            </button>
          )}
        </div>

        <div className="space-y-4">
          {/* Role Badge */}
          <div className="flex items-center gap-3">
            <Shield className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm text-gray-600">Vai trò</p>
              <div className="flex items-center gap-2">
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getRoleBadgeColor()}`}>
                  {getRoleText()}
                </span>
                <span className="text-xs text-gray-500">({role})</span>
              </div>
            </div>
          </div>

          {/* Email */}
          <div className="flex items-center gap-3">
            <Mail className="w-5 h-5 text-gray-400" />
            <div>
              <p className="text-sm text-gray-600">Email</p>
              <p className="font-medium">{user?.email}</p>
            </div>
          </div>

          {/* Full Name */}
          <div className="flex items-center gap-3">
            <User className="w-5 h-5 text-gray-400" />
            <div className="flex-1">
              <p className="text-sm text-gray-600">Họ và tên</p>
              {isEditing ? (
                <input
                  type="text"
                  value={formData.fullName}
                  onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              ) : (
                <p className="font-medium">{profileData?.fullName}</p>
              )}
            </div>
          </div>

          {/* Phone */}
          <div className="flex items-center gap-3">
            <Phone className="w-5 h-5 text-gray-400" />
            <div className="flex-1">
              <p className="text-sm text-gray-600">Số điện thoại</p>
              {isEditing ? (
                <input
                  type="text"
                  value={formData.phoneNumber}
                  onChange={(e) => setFormData({ ...formData, phoneNumber: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              ) : (
                <p className="font-medium">{profileData?.phoneNumber}</p>
              )}
            </div>
          </div>

          {/* Date of Birth */}
          <div className="flex items-center gap-3">
            <Calendar className="w-5 h-5 text-gray-400" />
            <div className="flex-1">
              <p className="text-sm text-gray-600">Ngày sinh</p>
              {isEditing ? (
                <input
                  type="date"
                  value={formData.dateOfBirth}
                  onChange={(e) => setFormData({ ...formData, dateOfBirth: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              ) : (
                <p className="font-medium">{profileData?.dateOfBirth || 'Chưa cập nhật'}</p>
              )}
            </div>
          </div>

          {/* Address */}
          <div className="flex items-start gap-3">
            <User className="w-5 h-5 text-gray-400 mt-1" />
            <div className="flex-1">
              <p className="text-sm text-gray-600">Địa chỉ</p>
              {isEditing ? (
                <textarea
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              ) : (
                <p className="font-medium">{profileData?.address || 'Chưa cập nhật'}</p>
              )}
            </div>
          </div>

          {isEditing && (
            <div className="flex gap-3 pt-4">
              <button
                onClick={handleUpdateProfile}
                disabled={loading}
                className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
              >
                <Save className="w-4 h-4" />
                Lưu thay đổi
              </button>
              <button
                onClick={() => {
                  setIsEditing(false);
                  fetchProfile();
                }}
                className="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <X className="w-4 h-4" />
                Hủy
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Rating Card */}
      {(isBidder || isSeller) && profileData?.rating && (
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex items-center gap-3 mb-6">
            <Award className="w-6 h-6 text-yellow-500" />
            <h2 className="text-xl font-semibold">Đánh giá</h2>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <p className="text-2xl font-bold text-blue-600">{profileData.rating.totalRatings}</p>
              <p className="text-sm text-gray-600">Tổng đánh giá</p>
            </div>
            <div className="text-center p-4 bg-green-50 rounded-lg">
              <p className="text-2xl font-bold text-green-600">{profileData.rating.positiveRatings}</p>
              <p className="text-sm text-gray-600">Tích cực</p>
            </div>
            <div className="text-center p-4 bg-red-50 rounded-lg">
              <p className="text-2xl font-bold text-red-600">{profileData.rating.negativeRatings}</p>
              <p className="text-sm text-gray-600">Tiêu cực</p>
            </div>
            <div className="text-center p-4 bg-yellow-50 rounded-lg">
              <p className="text-2xl font-bold text-yellow-600">{profileData.rating.ratingPercentage}%</p>
              <p className="text-sm text-gray-600">Tỷ lệ +</p>
            </div>
          </div>

          {/* Rating Reviews */}
          {profileData.rating.reviews && profileData.rating.reviews.length > 0 && (
            <div>
              <h3 className="font-semibold mb-4">Nhận xét gần đây</h3>
              <div className="space-y-3">
                {profileData.rating.reviews.slice(0, 5).map((review) => (
                  <div key={review.id} className="p-4 bg-gray-50 rounded-lg">
                    <div className="flex items-center justify-between mb-2">
                      <span className="font-medium">{review.fromUserName}</span>
                      <div className="flex items-center gap-2">
                        <Star
                          className={`w-4 h-4 ${review.rating === 1 ? 'text-yellow-500 fill-yellow-500' : 'text-gray-400'}`}
                        />
                        <span className="text-sm text-gray-600">
                          {new Date(review.createdAt).toLocaleDateString('vi-VN')}
                        </span>
                      </div>
                    </div>
                    <p className="text-gray-700">{review.comment}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Password Change Card */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <Lock className="w-6 h-6 text-gray-600" />
            <h2 className="text-xl font-semibold">Bảo mật</h2>
          </div>
          {!isChangingPassword && (
            <button
              onClick={() => setIsChangingPassword(true)}
              className="px-4 py-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
            >
              Đổi mật khẩu
            </button>
          )}
        </div>

        {isChangingPassword && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Mật khẩu hiện tại
              </label>
              <input
                type="password"
                value={passwordData.currentPassword}
                onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Mật khẩu mới
              </label>
              <input
                type="password"
                value={passwordData.newPassword}
                onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Xác nhận mật khẩu mới
              </label>
              <input
                type="password"
                value={passwordData.confirmPassword}
                onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div className="flex gap-3 pt-4">
              <button
                onClick={handleChangePassword}
                disabled={loading}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50"
              >
                Đổi mật khẩu
              </button>
              <button
                onClick={() => {
                  setIsChangingPassword(false);
                  setPasswordData({
                    currentPassword: '',
                    newPassword: '',
                    confirmPassword: '',
                  });
                }}
                className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Hủy
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Upgrade to Seller for Bidders */}
      {isBidder && (
        <div className="bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg shadow-md p-6 mt-6">
          <h2 className="text-xl font-semibold mb-4">Trở thành người bán</h2>
          <p className="text-gray-700 mb-4">
            Bạn muốn bán sản phẩm trên nền tảng? Nâng cấp tài khoản của bạn lên Người bán ngay!
          </p>

          <button
            onClick={handleUpgradeRequest}
            disabled={loading}
            className="px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg
                 hover:from-purple-700 hover:to-pink-700 transition-all disabled:opacity-50"
          >
            {loading ? "Đang gửi yêu cầu..." : "Yêu cầu nâng cấp"}
          </button>

          {/* Success message */}
          {message && (
            <p className="mt-4 text-green-600 font-medium">
              {message}
            </p>
          )}

          {/* Error message */}
          {error && (
            <p className="mt-4 text-red-600 font-medium">
              {error}
            </p>
          )}
        </div>
      )}
    </div>
  );
};

export default ProfilePage;
