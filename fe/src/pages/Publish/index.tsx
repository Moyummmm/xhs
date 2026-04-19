import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input, message, Select } from 'antd';
import { LeftOutlined } from '@ant-design/icons';
import { useMutation } from '@tanstack/react-query';
import { createNote } from '@/api/note';
import { ImageUploader } from '@/components/ImageUploader';
import { VideoUploader } from '@/components/VideoUploader';
import type { UploadFile } from 'antd';
import type { Note, CreateNoteRequest } from '@/types';

const PublishPage: React.FC = () => {
  const navigate = useNavigate();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [images, setImages] = useState<UploadFile[]>([]);
  const [videoUrl, setVideoUrl] = useState('');
  const [location, setLocation] = useState('');
  const [topicId, setTopicId] = useState<number | undefined>();

  const publishMutation = useMutation<Note, Error, CreateNoteRequest>({
    mutationFn: createNote,
    onSuccess: (data) => {
      message.success('发布成功');
      navigate(`/note/${data.id}`);
    },
    onError: () => {
      message.error('发布失败，请重试');
    },
  });

  const handleSubmit = () => {
    if (!title.trim()) {
      message.error('请输入标题');
      return;
    }

    const imageIds = images
      .filter((img) => img.response?.url)
      .map((img) => img.response.url);

    if (imageIds.length === 0 && !videoUrl) {
      message.error('请上传图片或视频');
      return;
    }

    publishMutation.mutate({
      title,
      content,
      image_ids: imageIds.length > 0 ? imageIds : undefined,
      video_url: videoUrl || undefined,
      location: location || undefined,
      topic_id: topicId,
    });
  };

  return (
    <div className="min-h-screen bg-[#f5f5f5]">
      {/* Header */}
      <div className="sticky top-0 z-40 bg-white border-b border-gray-100 px-4 py-3 flex items-center justify-between">
        <button onClick={() => navigate(-1)} className="flex items-center gap-1 text-gray-600 hover:text-gray-800">
          <LeftOutlined />
          <span>取消</span>
        </button>
        <h1 className="text-lg font-bold">发布笔记</h1>
        <Button
          type="primary"
          onClick={handleSubmit}
          loading={publishMutation.isPending}
        >
          发布
        </Button>
      </div>

      <div className="container-custom py-4 max-w-3xl">
        <div className="bg-white rounded-xl p-6 space-y-6">
          {/* Media Upload */}
          <div>
            <h3 className="text-sm font-medium text-gray-700 mb-3">图片/视频</h3>
            {!videoUrl ? (
              <ImageUploader value={images} onChange={setImages} maxCount={9} />
            ) : (
              <VideoUploader value={videoUrl} onChange={setVideoUrl} />
            )}
            {images.length === 0 && !videoUrl && (
              <div className="mt-2 text-center">
                <button
                  onClick={() => setVideoUrl('temp')}
                  className="text-sm text-primary-500 hover:text-primary-600"
                >
                  或上传视频
                </button>
              </div>
            )}
          </div>

          {/* Title */}
          <div>
            <h3 className="text-sm font-medium text-gray-700 mb-2">标题</h3>
            <Input
              placeholder="填写标题会有更多赞哦~"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              maxLength={100}
              size="large"
            />
          </div>

          {/* Content */}
          <div>
            <h3 className="text-sm font-medium text-gray-700 mb-2">正文</h3>
            <Input.TextArea
              placeholder="添加正文"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={8}
              className="resize-none"
            />
          </div>

          {/* Location */}
          <div>
            <h3 className="text-sm font-medium text-gray-700 mb-2">地点</h3>
            <Input
              placeholder="添加位置，让更多人看到你的笔记"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              prefix="\ud83d\udccd"
            />
          </div>

          {/* Topic */}
          <div>
            <h3 className="text-sm font-medium text-gray-700 mb-2">话题</h3>
            <Select
              placeholder="选择话题"
              value={topicId}
              onChange={setTopicId}
              className="w-full"
              allowClear
              options={[
                { label: '#日常', value: 1 },
                { label: '#美食', value: 2 },
                { label: '#旅行', value: 3 },
                { label: '#穿搭', value: 4 },
              ]}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default PublishPage;
