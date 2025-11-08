import React, { createContext, useContext, useState, useEffect } from 'react';
import { accountApi } from '../api/client';

const AuthContext = createContext(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Check if user is authenticated on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const [result, error] = await accountApi.me();
      if (error) {
        throw new Error('Auth check failed');
      }
      setUser(result);
      setIsAuthenticated(true);
    } catch (error) {
      console.error('Auth check failed:', error);
      setUser(null);
      setIsAuthenticated(false);
    } finally {
      setLoading(false);
    }
  };

  const login = async (email, password) => {
    try {
      const result = await accountApi.login({ email, password });

      if (result.success) {
        // After successful login, fetch user info
        await checkAuth();
        return { success: true };
      } else {
        return {
          success: false,
          error: result.error === 'invalid_credentials'
            ? '이메일 또는 비밀번호가 올바르지 않습니다'
            : '로그인에 실패했습니다',
        };
      }
    } catch (error) {
      console.error('Login failed:', error);
      return {
        success: false,
        error: '로그인에 실패했습니다. 다시 시도해주세요.',
      };
    }
  };

  const signup = async (email, password, nickname) => {
    try {
      const result = await accountApi.signup({ email, password, nickname });

      if (result.success) {
        return { success: true };
      } else {
        return {
          success: false,
          error: '회원가입에 실패했습니다',
        };
      }
    } catch (error) {
      console.error('Signup failed:', error);
      return {
        success: false,
        error: '회원가입에 실패했습니다. 다시 시도해주세요.',
      };
    }
  };

  const logout = async () => {
    try {
      await accountApi.logout();
    } catch (error) {
      console.error('Logout failed:', error);
    } finally {
      setUser(null);
      setIsAuthenticated(false);
    }
  };

  const value = {
    user,
    loading,
    isAuthenticated,
    login,
    signup,
    logout,
    checkAuth,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
