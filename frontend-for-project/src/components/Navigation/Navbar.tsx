import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { 
  Menu, X, Gavel, LogOut, User, LogIn, 
  Heart, Package, ShoppingCart, Settings, TrendingUp, Search 
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { useRole } from '../../hooks/useRole';
import { RoleBased } from '../Common/RoleBased';

const Navbar = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [profileMenuOpen, setProfileMenuOpen] = useState(false);
  const navigate = useNavigate();
  const { isAuthenticated, user, logout } = useAuth();
  const { isAdmin, isSeller } = useRole();

  const handleLogout = () => {
    logout();
    setMobileMenuOpen(false);
    setProfileMenuOpen(false);
    navigate('/');
  };

  return (
    <nav className="bg-white border-b border-gray-200 sticky top-0 z-50 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo */}
          <div className="flex items-center gap-6">
            <Link to="/" className="flex items-center gap-2 font-bold text-2xl text-blue-600">
              <Gavel className="w-8 h-8" />
              <span className="hidden sm:inline">Đấu Giá</span>
            </Link>

            {/* Desktop Navigation Links */}
            {isAuthenticated && (
              <div className="hidden lg:flex items-center gap-4">
                <RoleBased allowedRoles={['ROLE_BIDDER', 'ROLE_SELLER']}>
                  <Link
                    to="/watchlist"
                    className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                  >
                    <Heart className="w-4 h-4" />
                    <span className="text-sm font-medium">Yêu thích</span>
                  </Link>
                </RoleBased>

                <RoleBased allowedRoles={['ROLE_SELLER']}>
                  <Link
                    to="/seller/products"
                    className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                  >
                    <Package className="w-4 h-4" />
                    <span className="text-sm font-medium">Sản phẩm của tôi</span>
                  </Link>
                </RoleBased>

                <RoleBased allowedRoles={['ROLE_BIDDER', 'ROLE_SELLER']}>
                  <Link
                    to="/orders"
                    className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                  >
                    <ShoppingCart className="w-4 h-4" />
                    <span className="text-sm font-medium">Đơn hàng</span>
                  </Link>
                </RoleBased>

                <RoleBased allowedRoles={['ROLE_ADMIN']}>
                  <Link
                    to="/admin"
                    className="flex items-center gap-2 px-3 py-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                  >
                    <Settings className="w-4 h-4" />
                    <span className="text-sm font-medium">Quản trị</span>
                  </Link>
                </RoleBased>
              </div>
            )}
          </div>

          {/* Quick Search Link - Desktop */}
          <div className="hidden md:flex items-center flex-1 mx-8 justify-end">
            <Link
              to="/search"
              className="flex items-center gap-2 px-4 py-2 text-gray-700 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
            >
              <Search className="w-5 h-5" />
              <span className="text-sm font-medium">Tìm kiếm</span>
            </Link>
          </div>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center gap-4">
            {!isAuthenticated ? (
              <>
                <Link
                  to="/login"
                  className="flex items-center gap-2 px-4 py-2 text-gray-700 hover:text-blue-600 font-medium transition-colors"
                >
                  <LogIn className="w-5 h-5" />
                  Đăng nhập
                </Link>
                <Link
                  to="/register"
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors"
                >
                  Đăng ký
                </Link>
              </>
            ) : (
              <div className="relative">
                <button
                  onClick={() => setProfileMenuOpen(!profileMenuOpen)}
                  className="flex items-center gap-2 px-4 py-2 text-gray-700 hover:text-blue-600 font-medium rounded-lg hover:bg-blue-50 transition-colors"
                >
                  <User className="w-5 h-5" />
                  <span className="max-w-32 truncate">{user?.fullName}</span>
                  <span className="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded">
                    {isAdmin ? 'Admin' : isSeller ? 'Seller' : 'Bidder'}
                  </span>
                </button>

                {profileMenuOpen && (
                  <div className="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-gray-200 py-2">
                    <Link
                      to="/profile"
                      className="flex items-center gap-3 px-4 py-2 text-gray-700 hover:bg-blue-50 transition-colors"
                      onClick={() => setProfileMenuOpen(false)}
                    >
                      <User className="w-4 h-4" />
                      Hồ sơ cá nhân
                    </Link>
                    
                    <RoleBased allowedRoles={['ROLE_BIDDER']}>
                      <Link
                        to="/my-bids"
                        className="flex items-center gap-3 px-4 py-2 text-gray-700 hover:bg-blue-50 transition-colors"
                        onClick={() => setProfileMenuOpen(false)}
                      >
                        <TrendingUp className="w-4 h-4" />
                        Đấu giá của tôi
                      </Link>
                    </RoleBased>
                    
                    <div className="border-t border-gray-200 my-2"></div>
                    
                    <button
                      onClick={handleLogout}
                      className="w-full flex items-center gap-3 px-4 py-2 text-red-600 hover:bg-red-50 transition-colors"
                    >
                      <LogOut className="w-4 h-4" />
                      Đăng xuất
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
              className="p-2 rounded-lg text-gray-700 hover:bg-gray-100"
            >
              {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {mobileMenuOpen && (
          <div className="md:hidden border-t border-gray-200 py-4 space-y-4">
            <div className="px-4 space-y-2">
              <Link
                to="/search"
                className="flex items-center gap-2 px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
                onClick={() => setMobileMenuOpen(false)}
              >
                <Search className="w-5 h-5" />
                Tìm kiếm
              </Link>
            </div>

            <div className="px-4 space-y-2">
              {!isAuthenticated ? (
                <>
                  <Link
                    to="/login"
                    className="block px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
                  >
                    Đăng nhập
                  </Link>
                  <Link
                    to="/register"
                    className="block px-4 py-2 bg-blue-600 text-white rounded-lg text-center font-medium hover:bg-blue-700"
                  >
                    Đăng ký
                  </Link>
                </>
              ) : (
                <>
                  <Link
                    to="/profile"
                    className="block px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
                  >
                    Hồ sơ cá nhân
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg"
                  >
                    Đăng xuất
                  </button>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
