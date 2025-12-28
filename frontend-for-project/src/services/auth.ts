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

  async signInWithGoogle(token: string) {
    try {
      const response = await apiClient.post<AuthResponse>(endpoints.auth.signInGoogle, {
        token,
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
      const response = await apiClient.get<{ success: boolean; data: User }>(endpoints.users.profile);
      if (response.data.success) {
        return response.data.data;
      }
      return null;
    } catch {
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
