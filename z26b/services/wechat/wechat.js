/**
 * 微信服务 - 登录和支付
 * 
 * 注意：当前使用本地余额系统进行测试
 * 微信登录和支付功能需要在后端启用 WECHAT_ENABLED=true 后才能使用
 */

import { post, get } from '../_utils/request';

// ============================================
// 微信登录相关
// ============================================

/**
 * 微信登录
 * 使用 wx.login() 获取的 code 换取用户信息和 token
 * 
 * @returns {Promise<{data: {token: string, userInfo: Object, isNewUser: boolean}}>}
 */
export async function wxLogin() {
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

    // 3. 保存 token 到本地
    if (response.data?.token) {
      wx.setStorageSync('token', response.data.token);
      wx.setStorageSync('userInfo', response.data.userInfo);
    }

    return response;
  } catch (error) {
    console.error('微信登录失败:', error);
    
    // 检查是否是微信功能未启用的错误
    if (error.message?.includes('WECHAT_DISABLED') || error.message?.includes('未启用')) {
      return {
        error: true,
        code: 'WECHAT_DISABLED',
        message: '微信登录功能未启用，请使用测试账号',
      };
    }
    
    throw error;
  }
}

/**
 * 更新用户信息（头像、昵称）
 * 
 * @param {Object} userInfo - 用户信息
 * @param {string} userInfo.nickName - 昵称
 * @param {string} userInfo.avatarUrl - 头像URL
 * @returns {Promise<Object>}
 */
export async function updateWxUserInfo({ nickName, avatarUrl }) {
  try {
    const response = await post('/wechat/updateUserInfo', {
      nickName,
      avatarUrl,
    });
    
    // 更新本地存储
    const userInfo = wx.getStorageSync('userInfo') || {};
    userInfo.nickName = nickName;
    userInfo.avatar = avatarUrl;
    wx.setStorageSync('userInfo', userInfo);
    
    return response;
  } catch (error) {
    console.error('更新用户信息失败:', error);
    throw error;
  }
}

/**
 * 检查微信登录状态
 * 
 * @returns {Promise<boolean>}
 */
export async function checkWxLoginStatus() {
  return new Promise((resolve) => {
    wx.checkSession({
      success: () => resolve(true),
      fail: () => resolve(false),
    });
  });
}

/**
 * 获取本地存储的用户信息
 * 
 * @returns {Object|null}
 */
export function getLocalUserInfo() {
  return wx.getStorageSync('userInfo') || null;
}

/**
 * 获取本地存储的 token
 * 
 * @returns {string|null}
 */
export function getLocalToken() {
  return wx.getStorageSync('token') || null;
}

/**
 * 清除登录状态
 */
export function clearLoginStatus() {
  wx.removeStorageSync('token');
  wx.removeStorageSync('userInfo');
}

// ============================================
// 微信支付相关
// ============================================

/**
 * 创建微信支付订单并调起支付
 * 
 * @param {string} orderId - 订单ID
 * @returns {Promise<{success: boolean, message: string}>}
 */
export async function wxPay(orderId) {
  try {
    // 1. 调用后端创建支付订单
    const response = await post('/wechat/pay/create', {
      orderId,
    });

    if (!response.data) {
      throw new Error('创建支付订单失败');
    }

    const { timeStamp, nonceStr, package: packageStr, signType, paySign } = response.data;

    // 2. 调起微信支付
    const payResult = await new Promise((resolve, reject) => {
      wx.requestPayment({
        timeStamp,
        nonceStr,
        package: packageStr,
        signType,
        paySign,
        success: resolve,
        fail: reject,
      });
    });

    return {
      success: true,
      message: '支付成功',
      data: payResult,
    };
  } catch (error) {
    console.error('微信支付失败:', error);
    
    // 检查是否是微信支付未启用的错误
    if (error.message?.includes('WECHAT_PAY_DISABLED') || error.message?.includes('未启用')) {
      return {
        success: false,
        code: 'WECHAT_PAY_DISABLED',
        message: '微信支付功能未启用，请使用余额支付',
      };
    }
    
    // 用户取消支付
    if (error.errMsg?.includes('cancel')) {
      return {
        success: false,
        code: 'USER_CANCEL',
        message: '用户取消支付',
      };
    }
    
    return {
      success: false,
      message: error.message || '支付失败',
    };
  }
}

/**
 * 查询订单支付状态
 * 
 * @param {string} orderId - 订单ID
 * @returns {Promise<{isPaid: boolean, status: string}>}
 */
export async function queryPayStatus(orderId) {
  try {
    const response = await get(`/wechat/pay/query/${orderId}`);
    return response.data;
  } catch (error) {
    console.error('查询支付状态失败:', error);
    throw error;
  }
}

/**
 * 申请微信退款
 * 
 * @param {string} orderId - 订单ID
 * @param {string} reason - 退款原因
 * @returns {Promise<Object>}
 */
export async function wxRefund(orderId, reason = '') {
  try {
    const response = await post('/wechat/pay/refund', {
      orderId,
      reason,
    });
    return response;
  } catch (error) {
    console.error('申请退款失败:', error);
    
    // 检查是否是功能未实现的错误
    if (error.message?.includes('NOT_IMPLEMENTED')) {
      return {
        success: false,
        code: 'NOT_IMPLEMENTED',
        message: '微信退款功能暂未开放，请联系客服处理',
      };
    }
    
    throw error;
  }
}

// ============================================
// 支付方式选择
// ============================================

/**
 * 检查微信支付是否可用
 * 
 * @returns {Promise<boolean>}
 */
export async function isWxPayAvailable() {
  try {
    // 尝试调用微信支付查询接口，如果返回 WECHAT_PAY_DISABLED 则不可用
    const response = await get('/wechat/pay/query/test');
    return true;
  } catch (error) {
    if (error.message?.includes('WECHAT_PAY_DISABLED')) {
      return false;
    }
    // 其他错误（如订单不存在）不代表支付不可用
    return true;
  }
}

/**
 * 获取可用的支付方式
 * 
 * @returns {Promise<Array<{type: string, name: string, enabled: boolean}>>}
 */
export async function getAvailablePayMethods() {
  const methods = [
    {
      type: 'balance',
      name: '余额支付',
      enabled: true, // 余额支付始终可用
    },
    {
      type: 'wechat',
      name: '微信支付',
      enabled: await isWxPayAvailable(),
    },
  ];
  
  return methods.filter(m => m.enabled);
}
