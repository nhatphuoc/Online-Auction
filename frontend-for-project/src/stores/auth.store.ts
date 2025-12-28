import { create } from 'zustand';
import { User } from '../types';
import { authService } from '../services/auth';

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  setLoading: (loading: boolean) => void;
  logout: () => void;
  initializeAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isLoading: false,
  isAuthenticated: false,

  setUser: (user) => {
    set({ user, isAuthenticated: !!user });
  },

  setLoading: (loading) => {
    set({ isLoading: loading });
  },

  logout: () => {
    authService.logout();
    set({ user: null, isAuthenticated: false });
  },

  initializeAuth: async () => {
    set({ isLoading: true });
    try {
      const token = authService.getToken();
      if (token) {
        const isValid = await authService.validateJwt(token);
        if (isValid) {
          const user = await authService.getCurrentUser();
          set({ user: user || null, isAuthenticated: !!user });
        } else {
          authService.logout();
          set({ user: null, isAuthenticated: false });
        }
      }
    } catch (error) {
      console.error('Failed to initialize auth:', error);
      set({ user: null, isAuthenticated: false });
    } finally {
      set({ isLoading: false });
    }
  },
}));
