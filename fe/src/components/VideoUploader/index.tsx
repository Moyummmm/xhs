import React, { useState } from 'react';
import { Upload, message, Progress } from 'antd';
import type { UploadProps } from 'antd';
import { VideoCameraOutlined, LoadingOutlined } from '@ant-design/icons';
import { uploadVideo } from '@/api/upload';
import { CONFIG } from '@/constants/config';
import { formatFileSize } from '@/utils';

interface VideoUploaderProps {
  value?: string;
  onChange?: (url: string) => void;
  disabled?: boolean;
}

export const VideoUploader: React.FC<VideoUploaderProps> = ({
  value,
  onChange,
  disabled = false,
}) => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const handleUpload: UploadProps['customRequest'] = async (options: any) => {
    const { file, onSuccess, onError } = options;
    const uploadFile = file as File;

    // 文件大小校验
    if (uploadFile.size > CONFIG.MAX_VIDEO_SIZE) {
      message.error(`视频大小不能超过 ${formatFileSize(CONFIG.MAX_VIDEO_SIZE)}`);
      onError?.(new Error('File too large'));
      return;
    }

    // 文件类型校验
    if (!CONFIG.ALLOWED_VIDEO_TYPES.includes(uploadFile.type as any)) {
      message.error('仅支持 MP4 和 MOV 格式');
      onError?.(new Error('Invalid file type'));
      return;
    }

    setUploading(true);
    setProgress(0);
    try {
      const result = await uploadVideo(uploadFile, setProgress);
      onSuccess?.(result);
      onChange?.(result.url);
      message.success('上传成功');
    } catch (error) {
      onError?.(error as Error);
      message.error('上传失败，请重试');
    } finally {
      setUploading(false);
      setProgress(0);
    }
  };

  if (value) {
    return (
      <div className="relative w-full aspect-video bg-black rounded-lg overflow-hidden group">
        <video src={value} className="w-full h-full object-contain" controls />
        {!disabled && (
          <button
            onClick={() => onChange?.('')}
            className="absolute top-2 right-2 bg-black/50 text-white px-3 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity"
          >
            删除
          </button>
        )}
      </div>
    );
  }

  return (
    <div className="w-full">
      <Upload
        accept="video/mp4,video/quicktime"
        showUploadList={false}
        customRequest={handleUpload}
        disabled={disabled || uploading}
      >
        <div className="w-full aspect-video flex flex-col items-center justify-center border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-500 transition-colors cursor-pointer bg-gray-50">
          {uploading ? (
            <div className="flex flex-col items-center gap-4">
              <LoadingOutlined className="text-3xl text-primary-500" />
              <Progress percent={progress} className="w-48" />
              <span className="text-sm text-gray-500">上传中...</span>
            </div>
          ) : (
            <>
              <VideoCameraOutlined className="text-4xl text-gray-400" />
              <span className="text-sm text-gray-500 mt-2">点击上传视频</span>
              <span className="text-xs text-gray-400 mt-1">支持 MP4、MOV，最大 100MB</span>
            </>
          )}
        </div>
      </Upload>
    </div>
  );
};

export default VideoUploader;
