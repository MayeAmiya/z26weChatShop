import { config } from '../../config/index';

const SEARCH_HISTORY_KEY = 'search_history';
const MAX_HISTORY_COUNT = 10;

/** 获取搜索历史 - Mock */
function mockSearchHistory() {
  const { delay } = require('../_utils/delay');
  const { getSearchHistory } = require('../../model/search');
  return delay().then(() => getSearchHistory());
}

/** 获取搜索历史 */
export function getSearchHistory() {
  if (config.useMock) {
    return mockSearchHistory();
  }
  // 使用本地存储
  return new Promise((resolve) => {
    try {
      const history = wx.getStorageSync(SEARCH_HISTORY_KEY) || [];
      resolve({ historyWords: history });
    } catch (e) {
      resolve({ historyWords: [] });
    }
  });
}

/** 添加搜索历史 */
export function addSearchHistory(keyword) {
  if (!keyword || !keyword.trim()) return;
  try {
    let history = wx.getStorageSync(SEARCH_HISTORY_KEY) || [];
    // 移除重复项
    history = history.filter(item => item !== keyword);
    // 添加到开头
    history.unshift(keyword);
    // 限制数量
    if (history.length > MAX_HISTORY_COUNT) {
      history = history.slice(0, MAX_HISTORY_COUNT);
    }
    wx.setStorageSync(SEARCH_HISTORY_KEY, history);
  } catch (e) {
    console.error('保存搜索历史失败:', e);
  }
}

/** 清除搜索历史 */
export function clearSearchHistory() {
  try {
    wx.removeStorageSync(SEARCH_HISTORY_KEY);
  } catch (e) {
    console.error('清除搜索历史失败:', e);
  }
}

/** 获取热门搜索 - Mock */
function mockSearchPopular() {
  const { delay } = require('../_utils/delay');
  const { getSearchPopular } = require('../../model/search');
  return delay().then(() => getSearchPopular());
}

/** 获取热门搜索 */
export function getSearchPopular() {
  if (config.useMock) {
    return mockSearchPopular();
  }
  // 返回静态热门词
  return new Promise((resolve) => {
    resolve({ 
      popularWords: ['热卖', '新品', '特价', '包邮', '限时优惠'] 
    });
  });
}
