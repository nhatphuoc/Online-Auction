import { useAuthStore } from '../stores/auth.store';

// Hooks for role checking
export const useRole = () => {
  const user = useAuthStore((state) => state.user);

  return {
    role: user?.userRole || null,
    isAdmin: user?.userRole === 'ROLE_ADMIN',
    isSeller: user?.userRole === 'ROLE_SELLER',
    isBidder: user?.userRole === 'ROLE_BIDDER',
    isGuest: !user,
    isAuthenticated: !!user,
    user,
  };
};
