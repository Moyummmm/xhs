import React, { useCallback, useState } from 'react';
import { Upload, message } from 'antd';
import type { UploadFile, UploadProps } from 'antd';
import { PlusOutlined, LoadingOutlined } from '@ant-design/icons';
import { uploadImage } from '@/api/upload';
import { CONFIG } from '@/constants/config';
import { formatFileSize } from '@/utils';

interface ImageUploaderProps {
  maxCount?: number;
  value?: UploadFile[];
  onChange?: (files: UploadFile[]) => void;
  disabled?: boolean;
}

export const ImageUploader: React.FC<ImageUploaderProps> = ({
  maxCount = CONFIG.MAX_IMAGE_COUNT,
  value = [],
  onChange,
  disabled = false,
}) => {
  const [uploading, setUploading] = useState(false);

  const handleUpload: UploadProps['customRequest'] = useCallback(async (options: any) => {
    const { file, onSuccess, onError } = options;
    const uploadFile = file as File;

    // 文件大小校验
    if (uploadFile.size > CONFIG.MAX_IMAGE_SIZE) {
      message.error(`图片大小不能超过 ${formatFileSize(CONFIG.MAX_IMAGE_SIZE)}`);
      onError?.(new Error('File too large'));
      return;
    }

    setUploading(true);
    try {
      const result = await uploadImage(uploadFile);
      onSuccess?.(result);
    } catch (error) {
      onError?.(error as Error);
      message.error('上传失败，请重试');
    } finally {
      setUploading(false);
    }
  }, []);

  const handleChange: UploadProps['onChange'] = useCallback((info: any) => {
    const { fileList } = info;
    if (fileList.length > maxCount) {
      message.error(`最多上传 ${maxCount} 张图片`);
      return;
    }
    onChange?.(fileList);
  }, [maxCount, onChange]);

  const uploadButton = (
    <div className="w-[100px] h-[100px] flex flex-col items-center justify-center border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-500 transition-colors cursor-pointer">
      {uploading ? <LoadingOutlined className="text-2xl text-gray-400" /> : <PlusOutlined className="text-2xl text-gray-400" />}
      <span className="text-xs text-gray-400 mt-1">上传图片</span>
    </div>
  );

  return (
    <Upload
      listType="picture-card"
      fileList={value}
      customRequest={handleUpload}
      onChange={handleChange}
      multiple
      maxCount={maxCount}
      accept="image/*"
      disabled={disabled || uploading}
      beforeUpload={() => false} // 阻止默认上传行为
    >
      {value.length < maxCount && uploadButton}
    </Upload>
  );
};

export default ImageUploader;
