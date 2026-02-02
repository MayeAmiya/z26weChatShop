import updateManager from './common/updateManager';
import { login, isLoggedIn } from './services/auth/auth';
import { logEnvInfo } from './config/env';

App({
  onLaunch: function () {
    // 打印环境配置
    logEnvInfo();
    this.initUser();
  },
  
  onShow: function () {
    updateManager();
  },

  async initUser() {
    try {
      // 检查是否已登录
      if (!isLoggedIn()) {
        // 使用微信登录
        await login();
        console.log('✓ 微信登录成功');
      } else {
        console.log('✓ 用户已登录');
      }
    } catch (error) {
      console.error('✗ 登录失败:', error);
      // 显示登录失败提示
      wx.showToast({
        title: '登录失败，请重试',
        icon: 'none',
        duration: 2000,
      });
    }
  },
});
