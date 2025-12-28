import { ReactNode } from 'react';
import { useAuthStore } from '../../stores/auth.store';
import { UserRole } from '../../types';

interface RoleBasedProps {
  children: ReactNode;
  allowedRoles: UserRole[];
  fallback?: ReactNode;
}

export const RoleBased = ({ children, allowedRoles, fallback = null }: RoleBasedProps) => {
  const user = useAuthStore((state) => state.user);

  if (!user || !allowedRoles.includes(user.userRole as UserRole)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
};

interface RoleBasedNavigationProps {
  bidderContent?: ReactNode;
  sellerContent?: ReactNode;
  adminContent?: ReactNode;
  guestContent?: ReactNode;
}

export const RoleBasedNavigation = ({
  bidderContent,
  sellerContent,
  adminContent,
  guestContent,
}: RoleBasedNavigationProps) => {
  const user = useAuthStore((state) => state.user);

  if (!user) {
    return <>{guestContent}</>;
  }

  switch (user.userRole) {
    case 'ROLE_ADMIN':
      return <>{adminContent}</>;
    case 'ROLE_SELLER':
      return <>{sellerContent}</>;
    case 'ROLE_BIDDER':
      return <>{bidderContent}</>;
    default:
      return <>{guestContent}</>;
  }
};

