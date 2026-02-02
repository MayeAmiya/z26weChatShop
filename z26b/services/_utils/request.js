/**
 * 统一的HTTP请求工具
 * 用于与后端API通信
 */

import { getApiBaseUrl, IS_DEV_MODE } from '../../config/env';

// API基础地址
const API_BASE_URL = getApiBaseUrl();

// 请求超时时间（毫秒）
const REQUEST_TIMEOUT = 30000;

// 开发环境测试 OpenID（仅开发环境使用）
const DEV_TEST_OPENID = 'oTest_dev_openid_001';

/**
 * 获取当前用户的 OpenID
 * 优先使用登录后存储的 openid，开发环境使用测试 openid
 */
function getCurrentOpenID() {
  try {
    const userInfo = wx.getStorageSync('userInfo');
    if (userInfo?.openid) {
      return userInfo.openid;
    }
  } catch (e) {
    console.warn('获取用户信息失败:', e);
  }
  
  // 仅开发环境使用测试 OpenID
  if (IS_DEV_MODE) {
    return DEV_TEST_OPENID;
  }
  
  return null;
}

/**
 * 获取当前用户的 Token（如果已登录）
 */
function getCurrentToken() {
  try {
    return wx.getStorageSync('token') || null;
  } catch (e) {
    return null;
  }
}

/**
 * 发起HTTP请求
 * @param {string} url - 请求路径
 * @param {object} options - 请求配置
 * @returns {Promise<any>}
 */
export function request(url, options = {}) {
  const {
    method = 'GET',
    data = null,
    header = {},
    requireAuth = false, // 是否需要认证
  } = options;

  const openid = getCurrentOpenID();
  const token = getCurrentToken();
  
  // 开发环境下跳过认证检查，生产环境需要认证的接口必须有 openid
  if (!IS_DEV_MODE && requireAuth && !openid) {
    return Promise.reject(new Error('请先登录'));
  }
  
  // 开发环境下打印调试信息
  if (IS_DEV_MODE) {
    console.log(`[DEV] ${method} ${url}`, openid ? `(OpenID: ${openid.substring(0, 10)}...)` : '(无认证)');
  }

  // 构建请求头
  const requestHeader = {
    'Content-Type': 'application/json',
    ...header,
  };

  // 添加认证信息
  if (openid) {
    requestHeader['X-OpenID'] = openid;
  }
  if (token) {
    requestHeader['Authorization'] = `Bearer ${token}`;
  }

  return new Promise((resolve, reject) => {
    wx.request({
      url: `${API_BASE_URL}${url}`,
      method,
      data,
      timeout: REQUEST_TIMEOUT,
      header: requestHeader,
      success(res) {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data);
        } else if (res.statusCode === 401) {
          // 未授权，清除登录状态
          try {
            wx.removeStorageSync('userInfo');
            wx.removeStorageSync('token');
          } catch (e) {
            // 忽略清除错误
          }
          reject(new Error('未授权，请重新登录'));
        } else if (res.statusCode === 404) {
          reject(new Error('请求的资源不存在'));
        } else if (res.statusCode === 429) {
          reject(new Error('请求过于频繁，请稍后重试'));
        } else if (res.statusCode >= 500) {
          reject(new Error('服务器错误，请稍后重试'));
        } else {
          reject(new Error(res.data?.error || `请求失败: ${res.statusCode}`));
        }
      },
      fail(err) {
        if (err.errMsg?.includes('timeout')) {
          reject(new Error('请求超时，请检查网络连接'));
        } else if (err.errMsg?.includes('fail')) {
          reject(new Error('网络连接失败，请检查网络'));
        } else {
          reject(new Error(err.errMsg || '网络请求失败'));
        }
      },
    });
  });
}

/**
 * GET请求
 * @param {string} url - 请求路径
 * @param {object} params - 查询参数
 * @param {object} options - 额外选项
 */
export function get(url, params = {}, options = {}) {
  const query = Object.keys(params)
    .filter(key => params[key] !== undefined && params[key] !== null)
    .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
    .join('&');
  
  const fullUrl = query ? `${url}?${query}` : url;
  return request(fullUrl, { method: 'GET', ...options });
}

/**
 * POST请求
 * @param {string} url - 请求路径
 * @param {object} data - 请求数据
 * @param {object} options - 额外选项
 */
export function post(url, data = {}, options = {}) {
  return request(url, {
    method: 'POST',
    data,
    ...options,
  });
}

/**
 * PUT请求
 * @param {string} url - 请求路径
 * @param {object} data - 请求数据
 * @param {object} options - 额外选项
 */
export function put(url, data = {}, options = {}) {
  return request(url, {
    method: 'PUT',
    data,
    ...options,
  });
}

/**
 * DELETE请求
 * @param {string} url - 请求路径
 * @param {object} options - 额外选项
 */
export function del(url, options = {}) {
  return request(url, {
    method: 'DELETE',
    ...options,
  });
}

/**
 * 带认证的请求方法（用于需要登录的接口）
 */
export const authRequest = {
  get: (url, params = {}) => get(url, params, { requireAuth: true }),
  post: (url, data = {}) => post(url, data, { requireAuth: true }),
  put: (url, data = {}) => put(url, data, { requireAuth: true }),
  del: (url) => del(url, { requireAuth: true }),
};
