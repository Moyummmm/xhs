import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input, message, Radio, DatePicker } from 'antd';
import { LeftOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores';
import { updateUserInfo } from '@/api/user';
import dayjs from 'dayjs';

const EditProfilePage: React.FC = () => {
  const navigate = useNavigate();
  const { user, setUser } = useAuthStore();
  const [loading, setLoading] = useState(false);

  const [nickname, setNickname] = useState('');
  const [avatar, setAvatar] = useState('');
  const [bio, setBio] = useState('');
  const [gender, setGender] = useState(0);
  const [birthday, setBirthday] = useState<string | undefined>(undefined);

  useEffect(() => {
    if (user) {
      setNickname(user.nickname || '');
      setAvatar(user.avatar || '');
      setBio(user.bio || '');
      setGender(user.gender || 0);
      setBirthday(user.birthday);
    }
  }, [user]);

  const handleSubmit = async () => {
    if (!user?.id) return;
    if (!nickname.trim()) {
      message.error('请输入昵称');
      return;
    }

    setLoading(true);
    try {
      const updated = await updateUserInfo(user.id, {
        nickname: nickname.trim(),
        avatar: avatar.trim() || undefined,
        bio: bio.trim() || undefined,
        gender: gender || undefined,
        birthday: birthday || undefined,
      });
      setUser(updated);
      message.success('资料更新成功');
      navigate(`/user/${user.id}`);
    } catch {
      message.error('更新失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#f5f5f5]">
      {/* Header */}
      <div className="bg-white border-b border-gray-100 px-4 py-3 flex items-center sticky top-0 z-40">
        <button onClick={() => navigate(-1)} className="flex items-center gap-1 text-gray-600 hover:text-gray-800">
          <LeftOutlined />
          <span>返回</span>
        </button>
        <h1 className="text-lg font-bold ml-4">编辑资料</h1>
      </div>

      <div className="container-custom py-6 max-w-xl">
        <div className="bg-white rounded-xl p-6 space-y-6">
          {/* Avatar */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">头像</label>
            <div className="flex items-center gap-4">
              <div className="w-20 h-20 rounded-full overflow-hidden bg-gray-100">
                {avatar ? (
                  <img src={avatar} alt="avatar" className="w-full h-full object-cover" />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-400 text-2xl font-bold">
                    {nickname?.charAt(0)?.toUpperCase() || '?'}
                  </div>
                )}
              </div>
              <Input
                placeholder="头像 URL"
                value={avatar}
                onChange={(e) => setAvatar(e.target.value)}
                className="flex-1"
              />
            </div>
          </div>

          {/* Nickname */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">昵称 *</label>
            <Input
              value={nickname}
              onChange={(e) => setNickname(e.target.value)}
              maxLength={20}
              showCount
              placeholder="填写昵称"
            />
          </div>

          {/* Bio */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">简介</label>
            <Input.TextArea
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              maxLength={200}
              showCount
              rows={4}
              placeholder="介绍一下自己..."
              className="resize-none"
            />
          </div>

          {/* Gender */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">性别</label>
            <Radio.Group value={gender} onChange={(e) => setGender(e.target.value)}>
              <Radio value={0}>未知</Radio>
              <Radio value={1}>男</Radio>
              <Radio value={2}>女</Radio>
            </Radio.Group>
          </div>

          {/* Birthday */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">生日</label>
            <DatePicker
              value={birthday ? dayjs(birthday) : null}
              onChange={(date) => setBirthday(date ? date.format('YYYY-MM-DD') : undefined)}
              placeholder="选择生日"
              className="w-full"
            />
          </div>

          {/* Submit */}
          <Button
            type="primary"
            size="large"
            block
            loading={loading}
            onClick={handleSubmit}
          >
            保存
          </Button>
        </div>
      </div>
    </div>
  );
};

export default EditProfilePage;
