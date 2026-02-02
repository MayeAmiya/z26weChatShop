/**
 * 安全工具函数
 * 用于输入验证和防止 XSS 攻击
 * 
 * 注意：开发环境下验证较宽松，生产环境会更严格
 */

import { IS_DEV_MODE } from '../config/env';

/**
 * 转义 HTML 特殊字符，防止 XSS 攻击
 * @param {string} str - 需要转义的字符串
 * @returns {string} 转义后的字符串
 */
export function escapeHtml(str) {
  if (typeof str !== 'string') {
    return str;
  }
  
  const htmlEntities = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#x27;',
    '/': '&#x2F;',
  };
  
  return str.replace(/[&<>"'/]/g, (char) => htmlEntities[char]);
}

/**
 * 清理用户输入，移除潜在的危险内容
 * @param {string} str - 用户输入
 * @returns {string} 清理后的字符串
 */
export function sanitizeInput(str) {
  if (typeof str !== 'string') {
    return str;
  }
  
  return str
    .trim()
    // 移除 script 标签
    .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
    // 移除 on* 事件处理器
    .replace(/\bon\w+\s*=/gi, '')
    // 移除 javascript: 协议
    .replace(/javascript:/gi, '')
    // 移除 data: 协议（可能用于注入）
    .replace(/data:/gi, '');
}

/**
 * 验证手机号格式
 * @param {string} phone - 手机号
 * @returns {boolean} 是否有效
 */
export function isValidPhone(phone) {
  if (!phone || typeof phone !== 'string') {
    return false;
  }
  // 中国大陆手机号：1开头，第二位3-9，共11位
  const phoneRegex = /^1[3-9]\d{9}$/;
  return phoneRegex.test(phone);
}

/**
 * 验证用户名格式
 * @param {string} name - 用户名
 * @returns {{valid: boolean, message?: string}} 验证结果
 */
export function validateName(name) {
  if (!name || typeof name !== 'string') {
    return { valid: false, message: '姓名不能为空' };
  }
  
  const trimmed = name.trim();
  
  if (trimmed.length < 2) {
    return { valid: false, message: '姓名至少2个字符' };
  }
  
  if (trimmed.length > 20) {
    return { valid: false, message: '姓名不能超过20个字符' };
  }
  
  // 只允许中文、英文、数字
  if (!/^[\u4e00-\u9fa5a-zA-Z0-9]+$/.test(trimmed)) {
    return { valid: false, message: '姓名只能包含中文、英文和数字' };
  }
  
  return { valid: true };
}

/**
 * 验证地址格式
 * @param {string} address - 地址
 * @returns {{valid: boolean, message?: string}} 验证结果
 */
export function validateAddress(address) {
  if (!address || typeof address !== 'string') {
    return { valid: false, message: '地址不能为空' };
  }
  
  const trimmed = address.trim();
  
  if (trimmed.length < 5) {
    return { valid: false, message: '请输入详细地址' };
  }
  
  if (trimmed.length > 200) {
    return { valid: false, message: '地址不能超过200个字符' };
  }
  
  return { valid: true };
}

/**
 * 验证评论内容
 * @param {string} comment - 评论内容
 * @returns {{valid: boolean, message?: string}} 验证结果
 */
export function validateComment(comment) {
  if (!comment || typeof comment !== 'string') {
    return { valid: false, message: '评论内容不能为空' };
  }
  
  const trimmed = comment.trim();
  
  if (trimmed.length < 5) {
    return { valid: false, message: '评论内容至少5个字符' };
  }
  
  if (trimmed.length > 500) {
    return { valid: false, message: '评论内容不能超过500个字符' };
  }
  
  return { valid: true };
}

/**
 * 验证数量
 * @param {number} quantity - 数量
 * @param {number} max - 最大值
 * @returns {{valid: boolean, message?: string}} 验证结果
 */
export function validateQuantity(quantity, max = 999) {
  const num = parseInt(quantity, 10);
  
  if (isNaN(num)) {
    return { valid: false, message: '请输入有效数量' };
  }
  
  if (num < 1) {
    return { valid: false, message: '数量至少为1' };
  }
  
  if (num > max) {
    return { valid: false, message: `数量不能超过${max}` };
  }
  
  return { valid: true };
}

/**
 * 验证价格
 * @param {number} price - 价格
 * @returns {{valid: boolean, message?: string}} 验证结果
 */
export function validatePrice(price) {
  const num = parseFloat(price);
  
  if (isNaN(num)) {
    return { valid: false, message: '请输入有效金额' };
  }
  
  if (num < 0) {
    return { valid: false, message: '金额不能为负数' };
  }
  
  if (num > 9999999) {
    return { valid: false, message: '金额超出范围' };
  }
  
  return { valid: true };
}

/**
 * 防止重复提交
 * 开发环境下间隔时间较短，便于测试
 * 
 * @param {string} key - 操作标识
 * @param {number} delay - 间隔时间（毫秒）
 * @returns {function} 检查函数
 */
const submitTimestamps = {};
export function debounceSubmit(key, delay = 1000) {
  // 开发环境下缩短间隔时间，便于快速测试
  const actualDelay = IS_DEV_MODE ? Math.min(delay, 300) : delay;
  return () => {
    const now = Date.now();
    const lastTime = submitTimestamps[key] || 0;
    
    if (now - lastTime < actualDelay) {
      return false;
    }
    
    submitTimestamps[key] = now;
    return true;
  };
}

/**
 * 敏感信息脱敏
 */
export const mask = {
  /**
   * 手机号脱敏 (保留前3后4)
   * @param {string} phone 
   * @returns {string}
   */
  phone(phone) {
    if (!phone || phone.length !== 11) return phone;
    return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2');
  },
  
  /**
   * 姓名脱敏 (只显示姓)
   * @param {string} name 
   * @returns {string}
   */
  name(name) {
    if (!name || name.length < 2) return name;
    return name[0] + '*'.repeat(name.length - 1);
  },
  
  /**
   * 地址脱敏 (隐藏详细门牌号)
   * @param {string} address 
   * @returns {string}
   */
  address(address) {
    if (!address || address.length < 10) return address;
    // 保留前20个字符
    return address.substring(0, 20) + '****';
  },
};
