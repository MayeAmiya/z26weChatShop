import { config } from '../../config/index';
import { get } from '../_utils/request';

/** 获取商品列表 */
function mockFetchPromotion(ID = 0) {
  const { delay } = require('../_utils/delay');
  const { getPromotion } = require('../../model/promotion');
  return delay().then(() => getPromotion(ID));
}

/** 获取促销详情 */
export async function fetchPromotion(ID = 0) {
  if (config.useMock) {
    return mockFetchPromotion(ID);
  }
  
  try {
    const response = await get(`/home/promotions?id=${ID}`);
    return response.data;
  } catch (error) {
    console.error('获取促销详情失败:', error);
    throw error;
  }
}
