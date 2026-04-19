import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { login as apiLogin, logout as apiLogout, register as apiRegister, getCurrentUser } from '@/api/auth';
import type { User, LoginRequest, RegisterRequest, AuthResponse } from '@/types';
import { CONFIG } from '@/constants/config';
import { storage } from '@/utils/storage';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;

  // Actions
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  fetchCurrentUser: () => Promise<void>;
  setUser: (user: User) => void;
  setToken: (token: string) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (data: LoginRequest) => {
        set({ isLoading: true });
        try {
          const response: AuthResponse = await apiLogin(data);
          const { token, user } = response;
          storage.set(CONFIG.TOKEN_KEY, token);
          storage.set(CONFIG.USER_KEY, user);
          set({
            token,
            user,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      register: async (data: RegisterRequest) => {
        set({ isLoading: true });
        try {
          const response: AuthResponse = await apiRegister(data);
          const { token, user } = response;
          storage.set(CONFIG.TOKEN_KEY, token);
          storage.set(CONFIG.USER_KEY, user);
          set({
            token,
            user,
            isAuthenticated: true,
            isLoading: false,
          });
        } catch (error) {
          set({ isLoading: false });
          throw error;
        }
      },

      logout: async () => {
        try {
          await apiLogout();
        } catch {
          // Ignore logout API errors
        }
        storage.remove(CONFIG.TOKEN_KEY);
        storage.remove(CONFIG.USER_KEY);
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        });
      },

      fetchCurrentUser: async () => {
        const token = get().token || storage.get<string>(CONFIG.TOKEN_KEY);
        if (!token) return;

        try {
          const user = await getCurrentUser();
          if (user && user.id) {
            storage.set(CONFIG.USER_KEY, user);
            set({ user, isAuthenticated: true });
          }
        } catch {
          set({ user: null, token: null, isAuthenticated: false });
          storage.remove(CONFIG.TOKEN_KEY);
          storage.remove(CONFIG.USER_KEY);
        }
      },

      setUser: (user: User) => {
        set({ user });
        storage.set(CONFIG.USER_KEY, user);
      },

      setToken: (token: string) => {
        set({ token });
        storage.set(CONFIG.TOKEN_KEY, token);
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
