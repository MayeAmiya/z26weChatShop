/* eslint-disable no-param-reassign */
import { config } from '../../config/index';
import { get } from '../_utils/request';
import { getPrice } from './spu';

/** ��ȡ��Ʒ�б� */
function mockFetchGoodsList(params) {
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
      item.desc = '';
      if (item.spuTagList) {
        item.tags = item.spuTagList.map((tag) => tag.title);
      } else {
        item.tags = [];
      }
    });
  }
  return delay().then(() => {
    return data;
  });
}

/** 获取商品列表 */
export async function fetchGoodsList(params) {
  if (config.useMock) {
    return mockFetchGoodsList(params);
  }
  
  try {
    const response = await get('/goods/list', params);
    const data = response.data || {};
    const records = data.records || [];
    const spuList = await Promise.all(records.map(async (spu) => {
      let price = 0;
      try {
        price = await getPrice(spu._id || spu.id);
      } catch (e) {
        console.warn('获取商品价格失败:', spu._id || spu.id, e);
      }
      return {
        spuId: spu._id || spu.id,
        thumb: spu.cover_image || '',
        title: spu.name,
        price: price || 0,
        originPrice: 0,
        desc: spu.detail || '',
        tags: [],
        isPutOnSale: spu.status === 'ENABLED' ? 1 : 0,
      };
    }));
    return { spuList, totalCount: data.total || 0 };
  } catch (error) {
    console.error('获取商品列表失败:', error);
    throw error;
  }
}
