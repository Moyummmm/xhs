import { create } from 'zustand';

interface UIState {
  // 登录弹窗
  isLoginModalOpen: boolean;

  // 侧边栏 (移动端)
  isSidebarOpen: boolean;

  // Actions
  openLoginModal: () => void;
  closeLoginModal: () => void;
  toggleLoginModal: () => void;
  openSidebar: () => void;
  closeSidebar: () => void;
  toggleSidebar: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  isLoginModalOpen: false,
  isSidebarOpen: false,

  openLoginModal: () => set({ isLoginModalOpen: true }),
  closeLoginModal: () => set({ isLoginModalOpen: false }),
  toggleLoginModal: () => set((state) => ({ isLoginModalOpen: !state.isLoginModalOpen })),

  openSidebar: () => set({ isSidebarOpen: true }),
  closeSidebar: () => set({ isSidebarOpen: false }),
  toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),
}));
