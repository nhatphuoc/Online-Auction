import { Outlet } from 'react-router-dom';
import Navbar from '../components/Navigation/Navbar';
import Toast from '../components/Common/Toast';

const MainLayout = () => {
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      <Navbar />
      <main className="flex-1">
        <Outlet />
      </main>
      <footer className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div>
              <h3 className="text-lg font-bold mb-4">Sàn Đấu Giá Online</h3>
              <p className="text-gray-400">Platform đấu giá trực tuyến uy tín</p>
            </div>
            <div>
              <h4 className="text-sm font-semibold mb-4 text-gray-300">Hỗ Trợ</h4>
              <ul className="text-gray-400 space-y-2">
                <li><a href="#" className="hover:text-white">Trung tâm trợ giúp</a></li>
                <li><a href="#" className="hover:text-white">Liên hệ</a></li>
              </ul>
            </div>
            <div>
              <h4 className="text-sm font-semibold mb-4 text-gray-300">Về Chúng Tôi</h4>
              <ul className="text-gray-400 space-y-2">
                <li><a href="#" className="hover:text-white">Giới thiệu</a></li>
                <li><a href="#" className="hover:text-white">Điều khoản</a></li>
              </ul>
            </div>
            <div>
              <h4 className="text-sm font-semibold mb-4 text-gray-300">Pháp Lý</h4>
              <ul className="text-gray-400 space-y-2">
                <li><a href="#" className="hover:text-white">Chính sách bảo mật</a></li>
                <li><a href="#" className="hover:text-white">Quy tắc</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
            <p>&copy; 2025 Sàn Đấu Giá Online. All rights reserved.</p>
          </div>
        </div>
      </footer>
      <Toast />
    </div>
  );
};

export default MainLayout;
