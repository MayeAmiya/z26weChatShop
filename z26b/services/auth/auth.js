import { post } from '../_utils/request';
import { IS_DEV_MODE } from '../../config/env';

// 开发环境测试用户数据
const DEV_TEST_USER = {
  _id: 'dev_user_001',
  openid: 'oTest_dev_openid_001',
  nickName: '开发测试用户',
  avatar: '',
};

/**
 * 微信登录（主要登录方式）
 * 使用 wx.login 获取 code 进行登录
 * 开发环境下如果微信登录失败，会自动使用测试用户
 * @returns {Promise<{_id: string, openid: string, nickName: string, avatar: string, token?: string}>}
 */
export async function login() {
  // 开发环境：优先尝试真实登录，失败则使用测试用户
  if (IS_DEV_MODE) {
    try {
      return await _realLogin();
    } catch (error) {
      console.warn('[DEV] 微信登录失败，使用测试用户:', error.message);
      return devLogin();
    }
  }
  return _realLogin();
}

/**
 * 开发环境模拟登录
 * 直接使用测试用户数据，无需联网
 */
export function devLogin() {
  if (!IS_DEV_MODE) {
    throw new Error('devLogin 仅限开发环境使用');
  }
  
  const userInfo = { ...DEV_TEST_USER };
  wx.setStorageSync('userInfo', userInfo);
  wx.setStorageSync('userId', userInfo._id);
  console.log('[DEV] 使用测试用户登录:', userInfo.nickName);
  return userInfo;
}

/**
 * 真实微信登录逻辑
 */
async function _realLogin() {
  try {
    // 1. 调用微信登录获取 code
    const loginResult = await new Promise((resolve, reject) => {
      wx.login({
        success: resolve,
        fail: reject,
      });
    });

    if (!loginResult.code) {
      throw new Error('微信登录失败：未获取到code');
    }

    // 2. 将 code 发送到后端换取用户信息
    const response = await post('/wechat/login', {
      code: loginResult.code,
    });

    // 3. 保存用户信息和 token 到本地
    if (response.data) {
      const userData = response.data;
      const userInfo = {
        _id: userData.user?._id || userData._id || userData.id,
        openid: userData.user?.open_id || userData.openid || '',
        nickName: userData.user?.nick_name || userData.nickName || '微信用户',
        avatar: userData.user?.avatar || userData.avatar || '',
      };
      
      // 保存用户信息
      wx.setStorageSync('userInfo', userInfo);
      wx.setStorageSync('userId', userInfo._id);
      
      // 保存 token（如果后端返回）
      if (userData.token) {
        wx.setStorageSync('token', userData.token);
      }
      
      return userInfo;
    }

    throw new Error('登录失败：未获取到用户信息');
  } catch (error) {
    console.error('微信登录失败:', error);
    throw error;
  }
}

/**
 * 静默登录（仅获取 openid，不获取用户信息）
 * 用于初次进入小程序时建立会话
 * 开发环境下如果失败会使用测试用户
 */
export async function silentLogin() {
  // 开发环境：尝试真实登录，失败则使用测试用户
  if (IS_DEV_MODE) {
    try {
      return await _silentLoginReal();
    } catch (error) {
      console.warn('[DEV] 静默登录失败，使用测试用户');
      return devLogin();
    }
  }
  return _silentLoginReal();
}

async function _silentLoginReal() {
  const loginResult = await new Promise((resolve, reject) => {
    wx.login({
      success: resolve,
      fail: reject,
    });
  });

  if (!loginResult.code) {
    throw new Error('微信登录失败：未获取到code');
  }

  const response = await post('/wechat/login', {
    code: loginResult.code,
  });

  if (response.data) {
    const userData = response.data;
    const userInfo = {
      _id: userData.user?._id || userData._id || '',
      openid: userData.user?.open_id || userData.openid || '',
      nickName: userData.user?.nick_name || '微信用户',
      avatar: userData.user?.avatar || '',
    };
    
    wx.setStorageSync('userInfo', userInfo);
    
    if (userData.token) {
      wx.setStorageSync('token', userData.token);
    }
    
    return userInfo;
  }
  
  throw new Error('未获取到用户信息');
}

/**
 * 获取本地存储的用户信息
 */
export function getUserInfo() {
  try {
    return wx.getStorageSync('userInfo') || null;
  } catch (error) {
    return null;
  }
}

/**
 * 检查用户是否已登录
 */
export function isLoggedIn() {
  const userInfo = getUserInfo();
  // 开发环境下，只要有 userInfo 就算已登录
  // 生产环境下，需要有有效的 openid
  if (IS_DEV_MODE) {
    return !!userInfo;
  }
  return !!(userInfo && userInfo.openid);
}

/**
 * 检查登录状态，未登录则跳转登录页
 * @param {string} redirectUrl - 登录后跳转的页面
 * @returns {boolean} 是否已登录
 */
export function checkLoginStatus(redirectUrl = '') {
  if (isLoggedIn()) {
    return true;
  }
  
  // 未登录，跳转到登录页
  const url = redirectUrl 
    ? `/pages/usercenter/login/index?redirect=${encodeURIComponent(redirectUrl)}`
    : '/pages/usercenter/login/index';
  
  wx.navigateTo({ url });
  return false;
}

/**
 * 清除用户信息（登出）
 */
export function logout() {
  try {
    wx.removeStorageSync('userInfo');
    wx.removeStorageSync('userId');
    wx.removeStorageSync('token');
  } catch (error) {
    console.error('登出失败:', error);
  }
}

/**
 * 获取本地存储的 token
 */
export function getToken() {
  try {
    return wx.getStorageSync('token') || null;
  } catch (error) {
    return null;
  }
}
