import { createBrowserRouter } from 'react-router-dom';
import MainLayout from '../layouts/MainLayout';
import HomePage from '../pages/home/HomePage';
import LoginPage from '../pages/auth/LoginPage';
import RegisterPage from '../pages/auth/RegisterPage';
import VerifyOtpPage from '../pages/auth/VerifyOtpPage';
import ProductListPage from '../pages/products/ProductListPage';
import ProductDetailPage from '../pages/products/ProductDetailPage';
import SearchPage from '../pages/products/SearchPage';
import ProfilePage from '../pages/profile/ProfilePage';
import WatchlistPage from '../pages/bidder/WatchlistPage';
import MyBidsPage from '../pages/bidder/MyBidsPage';
import MyProductsPage from '../pages/seller/products/MyProductsPage';
import { CreateProductPage } from '../pages/seller/products/CreateProductPage';
import { EditProductPage } from '../pages/seller/products/EditProductPage';
import OrdersPage from '../pages/orders/OrdersPage';
import OrderDetailPage from '../pages/orders/OrderDetailPage';
import NotFoundPage from '../pages/common/NotFoundPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'category/:categoryId',
        element: <ProductListPage />,
      },
      {
        path: 'products/:id',
        element: <ProductDetailPage />,
      },
      {
        path: 'search',
        element: <SearchPage />,
      },
      {
        path: 'profile',
        element: <ProfilePage />,
      },
      {
        path: 'watchlist',
        element: <WatchlistPage />,
      },
      {
        path: 'my-bids',
        element: <MyBidsPage />,
      },
      {
        path: 'orders',
        element: <OrdersPage />,
      },
      {
        path: 'orders/:id',
        element: <OrderDetailPage />,
      },
      {
        path: 'seller/products',
        element: <MyProductsPage />,
      },
      {
        path: 'seller/products/create',
        element: <CreateProductPage />,
      },
      {
        path: 'seller/products/:id/edit',
        element: <EditProductPage />,
      },
    ],
  },
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/verify-otp',
    element: <VerifyOtpPage />,
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
]);
