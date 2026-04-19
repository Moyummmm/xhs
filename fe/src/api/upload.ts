import { http } from './request';
import { API_PATHS } from '@/constants/api';

/**
 * 上传图片
 */
export const uploadImage = (file: File, onProgress?: (percent: number) => void) => {
  return http.upload<{ url: string; width: number; height: number }>(
    API_PATHS.UPLOAD_IMAGE,
    file,
    onProgress
  );
};

/**
 * 上传视频
 */
export const uploadVideo = (file: File, onProgress?: (percent: number) => void) => {
  return http.upload<{ url: string; duration: number; cover_url: string }>(
    API_PATHS.UPLOAD_VIDEO,
    file,
    onProgress
  );
};
