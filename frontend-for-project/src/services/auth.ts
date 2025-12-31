import apiClient from './api/client';
import { endpoints } from './api/endpoints';
import { AuthResponse, User } from '../types';

export const authService = {
  async register(email: string, password: string, fullName: string, phoneNumber: string) {
    try {
      const response = await apiClient.post<AuthResponse>(endpoints.auth.register, {
        email,
        password,
        fullName,
        phoneNumber,
      });
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        message: error.response?.data?.message || 'Registration failed',
      };
    }
  },

  async verifyOtp(email: string, otpCode: string) {
    try {
      const response = await apiClient.post<AuthResponse>(endpoints.auth.verifyOtp, {
        email,
        otpCode,
      });
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        message: error.response?.data?.message || 'OTP verification failed',
      };
    }
  },

  async signIn(email: string, password: string) {
    try {
      const response = await apiClient.post<AuthResponse>(endpoints.auth.signIn, {
        email,
        password,
      });
      if (response.data.success && response.data.accessToken && response.data.refreshToken) {
        localStorage.setItem('accessToken', response.data.accessToken);
        localStorage.setItem('refreshToken', response.data.refreshToken);
      }
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        message: error.response?.data?.message || 'Sign in failed',
      };
    }
  },

  async signInWithGoogle(idToken: string) {
    try {
      const response = await apiClient.post<AuthResponse>(endpoints.auth.signInGoogle, {
        idToken,
      });
      if (response.data.success && response.data.accessToken && response.data.refreshToken) {
        localStorage.setItem('accessToken', response.data.accessToken);
        localStorage.setItem('refreshToken', response.data.refreshToken);
      }
      return response.data;
    } catch (error: any) {
      return {
        success: false,
        message: error.response?.data?.message || 'Google sign in failed',
      };
    }
  },

  async validateJwt(token: string) {
    try {
      const response = await apiClient.post<{ valid: boolean }>(endpoints.auth.validateJwt, {
        token,
      });
      return response.data.valid;
    } catch {
      return false;
    }
  },

  async getCurrentUser() {
    try {
      const response = await apiClient.get<{ success: boolean; data: any }>(endpoints.users.profile);
      if (response.data.success && response.data.data) {
        const backendUser = response.data.data;
        console.log('Backend user data:', backendUser);

        // Map backend response to User interface
        const user: User = {
          id: backendUser.id,
          email: backendUser.email,
          fullName: backendUser.fullName,
          phoneNumber: backendUser.phoneNumber || '',
          userRole: backendUser.role, // Backend uses 'role', we use 'userRole'
          isEmailVerified: backendUser.emailVerified,
          createdAt: backendUser.createdAt || new Date().toISOString(),
          updatedAt: backendUser.updatedAt || new Date().toISOString(),
        };

        console.log('Mapped user with role:', user.userRole);
        return user;
      }
      return null;
    } catch (error) {
      console.error('Error fetching current user:', error);
      return null;
    }
  },

  logout() {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  },

  getToken() {
    return localStorage.getItem('accessToken');
  },

  isAuthenticated() {
    return !!localStorage.getItem('accessToken');
  },
};
