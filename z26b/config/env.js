/**
 * 环境配置文件
 * 
 * ====================================
 * 当前处于【开发环境】
 * - 使用本地后端 localhost:8080
 * - 微信登录失败时自动使用测试用户
 * - 无需真实微信认证
 * 
 * 切换到生产环境时：
 * 1. 将 IS_DEV_MODE 改为 false
 * 2. 填写正确的 PRODUCTION_API_URL (HTTPS)
 * 3. 配置微信后台域名白名单
 * ====================================
 */

// ============ 主开关 ============
// true  = 开发模式：本地后端，允许测试用户，宽松验证
// false = 生产模式：远程API，必须真实微信登录
export const IS_DEV_MODE = true;

// ============ API地址配置 ============
export const API_CONFIG = {
  // 开发环境API地址
  development: 'http://localhost:8080/api',
  // 生产环境API地址（发布时需要修改为真实域名）
  // 注意：必须是 https 协议，且域名需要在小程序后台配置
  production: 'https://your-production-api.com/api',
};

// 获取当前环境的API地址
export function getApiBaseUrl() {
  if (IS_DEV_MODE) {
    return API_CONFIG.development;
  }
  return API_CONFIG.production;
}

// ============ 安全配置 ============
export const SECURITY_CONFIG = {
  // 是否启用调试模式（生产环境应设为 false）
  enableDebugLog: IS_DEV_MODE,
  // 请求超时时间（毫秒）
  requestTimeout: 30000,
  // 最大重试次数
  maxRetry: 3,
  // Token 过期提前刷新时间（秒）
  tokenRefreshBefore: 300,
};

// ============ 默认值配置 ============
export const DEFAULTS = {
  // 默认头像
  avatar: 'https://cdn-we-retail.ym.tencent.com/miniapp/usercenter/icon-user-center-avatar@2x.png',
  // 默认商品图片
  productImage: 'https://cdn-we-retail.ym.tencent.com/miniapp/goods/images/placeholder.png',
  // 默认店铺名称
  storeName: '默认店铺',
};

// ============ 功能开关 ============
export const FEATURES = {
  // 是否启用微信支付（需要后端配置 WECHAT_ENABLED=true）
  enableWxPay: false,
  // 是否启用余额支付
  enableBalancePay: true,
  // 是否显示调试信息
  showDebugInfo: IS_DEV_MODE,
};

// 打印当前环境信息（仅开发模式）
export function logEnvInfo() {
  if (!IS_DEV_MODE) return;
  
  console.log('========== 环境配置 ==========');
  console.log('开发模式:', IS_DEV_MODE);
  console.log('API地址:', getApiBaseUrl());
  console.log('调试日志:', SECURITY_CONFIG.enableDebugLog);
  console.log('==============================');
}

// ============ 生产环境检查 ============
/**
 * 检查生产环境配置是否完整
 * 在 app.js 启动时调用
 */
export function validateProductionConfig() {
  if (IS_DEV_MODE) return true;
  
  const errors = [];
  
  // 检查 API 地址
  if (API_CONFIG.production.includes('your-production-api.com')) {
    errors.push('请配置正确的生产环境 API 地址');
  }
  
  if (!API_CONFIG.production.startsWith('https://')) {
    errors.push('生产环境 API 必须使用 HTTPS 协议');
  }
  
  if (errors.length > 0) {
    console.error('⚠️ 生产环境配置错误:');
    errors.forEach(err => console.error('  -', err));
    return false;
  }
  
  return true;
}
