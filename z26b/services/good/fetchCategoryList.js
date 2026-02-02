import { config } from '../../config/index';
import { get } from '../_utils/request';

/** 获取商品列表 */
function mockFetchGoodCategory() {
  const { delay } = require('../_utils/delay');
  const { getCategoryList } = require('../../model/category');
  return delay().then(() => getCategoryList());
}

/** 获取商品分类列表 */
export async function getCategoryList() {
  if (config.useMock) {
    return mockFetchGoodCategory();
  }
  
  try {
    const response = await get('/goods/category/list');
    return response.data;
  } catch (error) {
    console.error('获取分类列表失败:', error);
    throw error;
  }
}
