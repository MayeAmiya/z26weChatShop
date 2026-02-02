/**
 * 用户服务 - 用户信息管理
 */

import { get, put } from '../_utils/request';

/**
 * 获取用户信息
 * @returns {Promise<{data: {_id: string, nickName: string, avatar: string}}>}
 */
export async function getUserInfo() {
  try {
    const response = await get('/user/info');
    return response;
  } catch (error) {
    console.error('获取用户信息失败:', error);
    throw error;
  }
}

/**
 * 更新用户信息
 * @param {{nickName?: string, avatar?: string}} data
 */
export async function updateUserInfo(data) {
  try {
    const response = await put('/user/info', data);
    return response;
  } catch (error) {
    console.error('更新用户信息失败:', error);
    throw error;
  }
}
