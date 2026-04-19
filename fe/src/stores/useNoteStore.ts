import { create } from 'zustand';
import type { Note } from '@/types';

interface NoteState {
  // 当前浏览的笔记 ID (用于详情页)
  currentNoteId: number | null;

  // 最近浏览的笔记
  recentNotes: Note[];

  // Actions
  setCurrentNoteId: (id: number | null) => void;
  addRecentNote: (note: Note) => void;
  clearRecentNotes: () => void;
}

const MAX_RECENT_NOTES = 20;

export const useNoteStore = create<NoteState>((set) => ({
  currentNoteId: null,
  recentNotes: [],

  setCurrentNoteId: (id: number | null) => set({ currentNoteId: id }),

  addRecentNote: (note: Note) =>
    set((state) => {
      // 移除重复的
      const filtered = state.recentNotes.filter((n) => n.id !== note.id);
      // 添加到开头，限制数量
      const newRecent = [note, ...filtered].slice(0, MAX_RECENT_NOTES);
      return { recentNotes: newRecent };
    }),

  clearRecentNotes: () => set({ recentNotes: [] }),
}));
