import { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Shield, Loader, RotateCcw } from 'lucide-react';
import { authService } from '../../services/auth';
import { useUIStore } from '../../stores/ui.store';

const VerifyOtpPage = () => {
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email') || '';
  const [otpCode, setOtpCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [resendTimer, setResendTimer] = useState(0);
  const navigate = useNavigate();
  const addToast = useUIStore((state) => state.addToast);

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (resendTimer > 0) {
      interval = setInterval(() => {
        setResendTimer((prev) => prev - 1);
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [resendTimer]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email) {
      addToast('error', 'Email không hợp lệ');
      return;
    }

    setIsLoading(true);
    try {
      const result = await authService.verifyOtp(email, otpCode);
      if (result.success) {
        addToast('success', 'Xác thực OTP thành công! Vui lòng đăng nhập');
        navigate('/login');
      } else {
        addToast('error', result.message || 'Xác thực OTP thất bại');
      }
    } catch (error) {
      addToast('error', 'Lỗi kết nối đến server');
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendOtp = async () => {
    setResendTimer(60);
    addToast('info', 'Mã OTP mới đã được gửi');
  };

  if (!email) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4">
        <div className="text-center">
          <p className="text-gray-600 mb-4">Email không hợp lệ</p>
          <button
            onClick={() => navigate('/register')}
            className="text-blue-600 hover:text-blue-700 font-medium"
          >
            Quay lại đăng ký
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-xl shadow-lg p-8">
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-12 h-12 bg-blue-100 rounded-lg mb-4">
              <Shield className="w-6 h-6 text-blue-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900">Xác Thực OTP</h1>
            <p className="text-gray-600 text-sm mt-2">Nhập mã OTP được gửi đến email</p>
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
            <p className="text-sm text-blue-800">
              Mã OTP đã được gửi đến <strong>{email}</strong>
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Mã OTP</label>
              <input
                type="text"
                value={otpCode}
                onChange={(e) => setOtpCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                placeholder="000000"
                maxLength={6}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent text-center text-3xl tracking-widest font-semibold"
                required
              />
              <p className="text-xs text-gray-500 mt-2">Nhập 6 chữ số</p>
            </div>

            <button
              type="submit"
              disabled={isLoading || otpCode.length < 6}
              className="w-full flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-3 rounded-lg transition-colors mt-6"
            >
              {isLoading ? <Loader className="w-5 h-5 animate-spin" /> : <Shield className="w-5 h-5" />}
              {isLoading ? 'Đang xác thực...' : 'Xác Thực'}
            </button>
          </form>

          <button
            onClick={handleResendOtp}
            disabled={resendTimer > 0}
            className="w-full flex items-center justify-center gap-2 mt-4 text-blue-600 hover:text-blue-700 disabled:text-gray-400 font-medium text-sm"
          >
            <RotateCcw className="w-4 h-4" />
            {resendTimer > 0 ? `Gửi lại trong ${resendTimer}s` : 'Gửi lại mã OTP'}
          </button>

          <div className="text-center mt-6 pt-6 border-t border-gray-200">
            <p className="text-gray-600 text-sm">
              <button
                onClick={() => navigate('/register')}
                className="text-blue-600 hover:text-blue-700 font-medium"
              >
                Quay lại đăng ký
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default VerifyOtpPage;
