import { config } from '../../config/index';
import { get } from '../_utils/request';

/** 获取商品列表 */
function mockFetchGood(ID = 0) {

  const { delay } = require('../_utils/delay');
  const { genGood } = require('../../model/good');

  return delay().then(() => {
    const res = genGood(ID);
    return res;
  });
}

/** 获取商品详情 */
export async function fetchGood(ID = 0) {
  if (config.useMock) {
    return mockFetchGood(ID);
  }
  
  try {
    const response = await get(`/goods/${ID}`);
    const spu = response.data;
    // 映射后端字段到前端格式
    return {
      ...spu,
      spuId: spu._id || spu.id,
      primaryImage: spu.cover_image || spu.primary_image,
      minSalePrice: spu.skus?.[0]?.price || spu.min_sale_price || 0,
      maxLinePrice: spu.skus?.[spu.skus.length - 1]?.price || spu.max_line_price || 0,
      isPutOnSale: spu.status === 'ENABLED' ? 1 : 0,
      skuList: spu.skus || [],
    };
  } catch (error) {
    console.error('获取商品详情失败:', error);
    throw error;
  }
}
