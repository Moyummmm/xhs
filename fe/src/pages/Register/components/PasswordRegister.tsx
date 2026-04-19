import React, { useState } from 'react';
import { Button, Input, message } from 'antd';
import { useAuthStore } from '@/stores';
import { validators } from '@/utils/validation';

interface PasswordRegisterProps {
  onSuccess?: () => void;
  onClose?: () => void;
  onSwitchToLogin?: () => void;
}

const PasswordRegister: React.FC<PasswordRegisterProps> = ({
  onSuccess,
  onClose,
  onSwitchToLogin,
}) => {
  const [nickname, setNickname] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<{
    nickname?: string;
    password?: string;
    confirmPassword?: string;
  }>({});
  const { register, isLoading } = useAuthStore();

  const handleSubmit = async () => {
    const nicknameError = validators.nickname(nickname);
    const passwordError = validators.password(password);
    const confirmError = password !== confirmPassword ? '两次密码输入不一致' : '';

    if (nicknameError || passwordError || confirmError) {
      setErrors({
        nickname: nicknameError,
        password: passwordError,
        confirmPassword: confirmError,
      });
      return;
    }

    try {
      await register({ nickname, password });
      message.success('注册成功');
      onSuccess?.();
      onClose?.();
    } catch (error) {
      // Error handled in store
    }
  };

  return (
    <div className="space-y-4">
      <div>
        <Input
          placeholder="请输入昵称 (2-20字符)"
          value={nickname}
          onChange={(e) => {
            setNickname(e.target.value);
            setErrors((prev) => ({ ...prev, nickname: undefined }));
          }}
          status={errors.nickname ? 'error' : undefined}
          size="large"
        />
        {errors.nickname && (
          <p className="text-red-500 text-xs mt-1">{errors.nickname}</p>
        )}
      </div>

      <div>
        <Input.Password
          placeholder="请输入密码 (至少6位)"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
            setErrors((prev) => ({ ...prev, password: undefined }));
          }}
          status={errors.password ? 'error' : undefined}
          size="large"
        />
        {errors.password && (
          <p className="text-red-500 text-xs mt-1">{errors.password}</p>
        )}
      </div>

      <div>
        <Input.Password
          placeholder="请再次输入密码"
          value={confirmPassword}
          onChange={(e) => {
            setConfirmPassword(e.target.value);
            setErrors((prev) => ({ ...prev, confirmPassword: undefined }));
          }}
          status={errors.confirmPassword ? 'error' : undefined}
          size="large"
        />
        {errors.confirmPassword && (
          <p className="text-red-500 text-xs mt-1">{errors.confirmPassword}</p>
        )}
      </div>

      <Button
        type="primary"
        block
        size="large"
        loading={isLoading}
        onClick={handleSubmit}
      >
        注册
      </Button>

      <div className="text-center">
        <span className="text-gray-500 text-sm">已有账号？</span>
        <button
          type="button"
          className="text-primary-500 text-sm ml-1 hover:underline"
          onClick={onSwitchToLogin}
        >
          去登录
        </button>
      </div>
    </div>
  );
};

export default PasswordRegister;
