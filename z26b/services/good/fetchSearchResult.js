/* eslint-disable no-param-reassign */
import { config } from '../../config/index';
import { get } from '../_utils/request';

/** 获取搜索历史 */
function mockSearchResult(params) {
  const { delay } = require('../_utils/delay');
  const { getSearchResult } = require('../../model/search');

  const data = getSearchResult(params);

  if (data.spuList.length) {
    data.spuList.forEach((item) => {
      item.spuId = item.spuId;
      item.thumb = item.primaryImage;
      item.title = item.title;
      item.price = item.minSalePrice;
      item.originPrice = item.maxLinePrice;
      if (item.spuTagList) {
        item.tags = item.spuTagList.map((tag) => ({ title: tag.title }));
      } else {
        item.tags = [];
      }
    });
  }
  return delay().then(() => {
    return data;
  });
}

/** 获取搜索结果 */
export async function getSearchResult(params) {
  if (config.useMock) {
    return mockSearchResult(params);
  }
  
  try {
    const response = await get('/goods/search', params);
    // 后端返回 { data: { records: [...], total: number } }
    const records = response.data?.records || response.data || [];
    const total = response.data?.total || records.length;
    const spuList = records.map(spu => ({
      spuId: spu._id || spu.id,
      thumb: spu.cover_image || spu.primary_image,
      title: spu.name || spu.title,
      price: spu.min_sale_price,
      originPrice: spu.max_line_price,
      tags: spu.tags || [],
    }));
    return { spuList, totalCount: total };
  } catch (error) {
    console.error('搜索失败:', error);
    throw error;
  }
}
