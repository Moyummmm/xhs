import React, { useState } from 'react';
import { Modal } from 'antd';
import PasswordLogin from '@/pages/Login/components/PasswordLogin';
import PasswordRegister from '@/pages/Register/components/PasswordRegister';

export type AuthMode = 'login' | 'register';

interface LoginModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess?: () => void;
  mode?: AuthMode;
}

export const LoginModal: React.FC<LoginModalProps> = ({
  open,
  onClose,
  onSuccess,
  mode: initialMode = 'login',
}) => {
  const [mode, setMode] = useState<AuthMode>(initialMode);

  const handleSwitchToLogin = () => setMode('login');
  const handleSwitchToRegister = () => setMode('register');

  const handleSuccess = () => {
    onSuccess?.();
    onClose?.();
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      footer={null}
      width={400}
      centered
      className="login-modal"
    >
      <div className="p-6">
        <h2 className="text-2xl font-bold text-center mb-6">
          {mode === 'login' ? '登录' : '注册'}
        </h2>
        {mode === 'login' ? (
          <PasswordLogin onSuccess={handleSuccess} onClose={onClose} onSwitchToRegister={handleSwitchToRegister} />
        ) : (
          <PasswordRegister
            onSuccess={handleSuccess}
            onClose={onClose}
            onSwitchToLogin={handleSwitchToLogin}
          />
        )}
      </div>
    </Modal>
  );
};

export default LoginModal;
