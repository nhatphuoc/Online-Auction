# Frontend Implementation Plan - Online Auction System

## Project Overview
Building a comprehensive React + TypeScript frontend for an Online Auction System with 4 user roles: Guest, Bidder, Seller, and Admin. The system includes real-time bidding, product management, order processing, and chat functionality.

---

## Phase 1: Infrastructure Setup

### 1.1 Add Required Dependencies
**Packages to install:**
- `react-router-dom@latest` - Routing & navigation
- `zustand@latest` - State management (lightweight)
- `react-hook-form@latest` + `zod@latest` - Form handling & validation
- `axios@latest` - HTTP client for API calls
- `date-fns@latest` - Date formatting utilities
- `recharts@latest` - Dashboard charts (for admin)

### 1.2 API Client Setup
**Create: `src/services/api/client.ts`**
- Initialize Axios with base URL
- Add token header interceptor
- Handle token refresh logic

**Create: `src/services/api/endpoints.ts`**
- Centralized API endpoint constants
- Organized by service

### 1.3 Authentication Service
**Create: `src/services/auth.ts`**
- Implement auth service using API Gateway
- Methods: register, verifyOtp, signIn, validateJWT, signInWithGoogle
- Token storage/retrieval from localStorage

---

## Phase 2: State Management & Core Services

### 2.1 Zustand Stores
**Create: `src/stores/auth.store.ts`** - Current user & auth status
**Create: `src/stores/ui.store.ts`** - Global UI state
**Create: `src/stores/product.store.ts`** - Product search & filters

### 2.2 WebSocket Services
**Create: `src/services/websocket/comment.ts`** - Product Q&A real-time
**Create: `src/services/websocket/order.ts`** - Order chat real-time

---

## Phase 3: Page Structure & Routing

### 3.1 Main Routing & Layout
- `src/App.tsx` - Update with Router setup
- `src/routes/index.tsx` - Route definitions
- `src/layouts/MainLayout.tsx` - Main layout wrapper
- `src/components/Navigation/Navbar.tsx` - Top navigation

### 3.2 Authentication Pages
- `src/pages/auth/LoginPage.tsx` - Email/password + Google OAuth
- `src/pages/auth/RegisterPage.tsx` - Multi-step form with OTP
- `src/pages/auth/VerifyOtpPage.tsx` - OTP verification
- `src/pages/auth/ForgotPasswordPage.tsx` - Password reset

### 3.3 Guest/Public Pages
- `src/pages/home/HomePage.tsx` - Hero + top 5 products sections
- `src/pages/products/ProductListPage.tsx` - Category products with pagination
- `src/pages/products/ProductDetailPage.tsx` - Full product details + Q&A
- `src/pages/products/SearchPage.tsx` - Full-text search + filters

### 3.4 Bidder Pages
- `src/pages/bidding/BidPage.tsx` - Place bid with auto-bid option
- `src/pages/bidding/MyBidsPage.tsx` - Active bids list
- `src/pages/bidding/WatchListPage.tsx` - Saved products
- `src/pages/profile/BidderProfilePage.tsx` - Profile + rating history

### 3.5 Seller Pages
- `src/pages/selling/CreateProductPage.tsx` - Multi-step product creation
- `src/pages/selling/EditProductPage.tsx` - Append description only
- `src/pages/selling/MyProductsPage.tsx` - Products dashboard
- `src/pages/selling/ProductQuestionsPage.tsx` - Answer bidder questions
- `src/pages/profile/SellerProfilePage.tsx` - Seller-specific profile

### 3.6 Order Management
- `src/pages/orders/OrderCheckoutPage.tsx` - 5-step checkout + chat
- `src/pages/orders/MyOrdersPage.tsx` - Orders list with filters

### 3.7 Admin Pages
- `src/pages/admin/DashboardPage.tsx` - Statistics & charts
- `src/pages/admin/CategoriesPage.tsx` - Category CRUD
- `src/pages/admin/ProductsPage.tsx` - Product management
- `src/pages/admin/UsersPage.tsx` - User management
- `src/pages/admin/SellerRequestsPage.tsx` - Seller approval requests

---

## Phase 4: Shared Components

### 4.1 Product Components
- `src/components/Products/ProductCard.tsx` - List item display
- `src/components/Products/ProductGallery.tsx` - Image carousel
- `src/components/Products/BidHistoryTable.tsx` - Bid history with masked names
- `src/components/Products/QuestionsSection.tsx` - Q&A display

### 4.2 Form Components
- `src/components/Forms/TextInput.tsx` - React Hook Form integration
- `src/components/Forms/FileUpload.tsx` - Single & multiple upload
- `src/components/Forms/DatePicker.tsx` - Date/time selection

### 4.3 Layout Components
- `src/components/Common/Loading.tsx` - Skeletons & spinners
- `src/components/Common/Modal.tsx` - Generic modal
- `src/components/Common/Toast.tsx` - Notifications
- `src/components/Common/Pagination.tsx` - Page navigation

### 4.4 Chat Components
- `src/components/Chat/ChatWidget.tsx` - Order chat (WebSocket)
- `src/components/Chat/CommentSection.tsx` - Product Q&A (WebSocket)

---

## Phase 5: Utilities & Helpers

### 5.1 Formatters & Validators
- `src/utils/formatters.ts` - formatCurrency, formatDate, formatTimeRemaining, etc.
- `src/utils/validators.ts` - Email, password, phone, price validation

### 5.2 Custom Hooks
- `src/hooks/useAuth.ts` - Auth store access
- `src/hooks/useFetch.ts` - Data fetching
- `src/hooks/useWebSocket.ts` - WebSocket connection
- `src/hooks/useAsync.ts` - Loading/error state
- `src/hooks/useLocalStorage.ts` - Persistent state

---

## Phase 6: Styling & Design

### 6.1 Design System (Tailwind)
- Color palette: primary, secondary, accent, success, warning, error, neutral
- 8px spacing system
- Typography system with 3 font weights
- Responsive breakpoints: 640px, 768px, 1024px, 1280px

### 6.2 Features
- Mobile-first responsive design
- Dark mode support (optional future)
- Consistent animations & transitions

---

## MVP Scope & Priorities (Based on Discussion)

### Phase 1: MVP Foundation (Week 1-2)
**Priority 1 - MUST HAVE:**
- Authentication (Register → OTP → Login)
- Home page with top 5 products
- Product listing by category (with pagination)
- Product detail page
- Search functionality

**Priority 2 - IMPORTANT (Week 2-3):**
- Bidding system (place bid, bid history)
- User profile (bidder + seller)
- Order checkout process
- Real-time chat (WebSocket)

**Priority 3 - NICE-TO-HAVE (Week 4+):**
- Admin dashboard
- Automated payment integration
- Advanced seller features
- Dark mode

### Phase 2: Post-MVP Enhancements
- Payment gateway integration (Stripe/VNPay for production)
- Advanced analytics
- Performance optimizations

---

## Design System Decision: Modern Minimal

**Color Palette:**
- Primary: #3B82F6 (Blue 500)
- Secondary: #6B7280 (Gray 500)
- Accent: #10B981 (Emerald 500)
- Success: #10B981 (Green)
- Warning: #F59E0B (Amber)
- Error: #EF4444 (Red)
- Neutral: #F3F4F6 to #1F2937 (Gray scale)

**Typography:**
- Headings: 3 weights max (600, 700, 800)
- Body: 400, 500
- Line height: 150% body, 120% headings

**Spacing:** 8px base unit system

---

## Environment Variables
```
VITE_API_BASE_URL=http://localhost:8080/api
VITE_GOOGLE_CLIENT_ID=your_google_client_id
```

---

## Key Decisions

1. **Routing:** React Router v6 - industry standard
2. **State:** Zustand - lightweight, minimal boilerplate
3. **Forms:** React Hook Form + Zod - type-safe, performant
4. **HTTP:** Axios - better error handling
5. **Styling:** Tailwind CSS - already configured
6. **Icons:** Lucide React - available
7. **Dates:** date-fns - lightweight
8. **WebSocket:** Native API + custom service wrapper
9. **Code Splitting:** React.lazy() per route

---

## Critical Success Factors

1. **Authentication Flow:** Properly implement token management & refresh
2. **WebSocket Integration:** Real-time chat & comments must be reliable
3. **Form Validation:** Type-safe validation with Zod
4. **API Error Handling:** Centralized error management
5. **State Synchronization:** Zustand stores properly initialized
6. **Mobile Responsive:** Test all pages on mobile devices
7. **Role-Based Access:** Proper authorization checks for all pages
8. **Loading States:** Show skeleton screens & spinners

---

## File Structure (Final)
```
src/
├── types/                      # Shared TypeScript types
├── services/                   # API & business logic
│   ├── api/
│   ├── auth.ts
│   ├── product.ts
│   ├── bidding.ts
│   ├── order.ts
│   ├── user.ts
│   ├── category.ts
│   └── websocket/
├── stores/                     # Zustand state management
├── routes/                     # Route definitions
├── layouts/                    # Layout components
├── pages/                      # Page components
│   ├── auth/
│   ├── home/
│   ├── products/
│   ├── bidding/
│   ├── selling/
│   ├── orders/
│   ├── profile/
│   └── admin/
├── components/                 # Shared components
│   ├── Navigation/
│   ├── Products/
│   ├── Forms/
│   ├── Chat/
│   └── Common/
├── hooks/                      # Custom React hooks
├── utils/                      # Utilities & helpers
├── styles/                     # Design tokens
├── App.tsx
├── main.tsx
└── index.css
```

**Total Scope:** ~80-100 files, ~10,000 LOC
**Estimated Timeline:** 5-6 weeks for full implementation
