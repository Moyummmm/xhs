import React, { useState } from 'react';
import { Button, Input, message } from 'antd';
import { useAuthStore } from '@/stores';
import { validators } from '@/utils/validation';

interface PasswordLoginProps {
  onSuccess?: () => void;
  onClose?: () => void;
  onSwitchToRegister?: () => void;
}

const PasswordLogin: React.FC<PasswordLoginProps> = ({ onSuccess, onClose, onSwitchToRegister }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ username?: string; password?: string }>({});
  const { login, isLoading } = useAuthStore();

  const handleSubmit = async () => {
    const usernameError = validators.required(username, '用户名');
    const passwordError = validators.password(password);

    if (usernameError || passwordError) {
      setErrors({ username: usernameError, password: passwordError });
      return;
    }

    try {
      await login({ username, password });
      message.success('登录成功');
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
          placeholder="请输入用户名"
          value={username}
          onChange={(e) => {
            setUsername(e.target.value);
            setErrors((prev) => ({ ...prev, username: undefined }));
          }}
          status={errors.username ? 'error' : undefined}
          size="large"
        />
        {errors.username && <p className="text-red-500 text-xs mt-1">{errors.username}</p>}
      </div>

      <div>
        <Input.Password
          placeholder="请输入密码"
          value={password}
          onChange={(e) => {
            setPassword(e.target.value);
            setErrors((prev) => ({ ...prev, password: undefined }));
          }}
          status={errors.password ? 'error' : undefined}
          size="large"
        />
        {errors.password && <p className="text-red-500 text-xs mt-1">{errors.password}</p>}
      </div>

      <Button type="primary" block size="large" loading={isLoading} onClick={handleSubmit}>
        登录
      </Button>

      <div className="text-center">
        <span className="text-gray-500 text-sm">没有账号？</span>
        <button
          type="button"
          className="text-primary-500 text-sm ml-1 hover:underline"
          onClick={onSwitchToRegister}
        >
          去注册
        </button>
      </div>
    </div>
  );
};

export default PasswordLogin;
