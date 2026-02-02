/* eslint-disable eqeqeq */
import { get } from '../_utils/request';
import { cloudbaseTemplateConfig } from '../../config/index';
import { CATEGORY, SPU } from '../cloudbaseMock/index';

// TODO: we should do pagination
export async function getAllSpuOfCate(cateId) {
  if (cloudbaseTemplateConfig.useMock) {
    return { spu: CATEGORY.find((x) => x._id === cateId).spu.map(({ _id }) => SPU.find((x) => x._id === _id)) };
  }

  try {
    // "全部"分类：获取所有商品
    if (cateId === 'all') {
      const goodsResp = await get('/goods/list?pageSize=100');
      const goods = goodsResp.data?.records || [];
      return { spu: goods };
    }
    
    // 直接查询该分类的商品
    const goodsResp = await get(`/goods/list?categoryId=${cateId}`);
    const goods = goodsResp.data?.records || [];
    
    // 如果没有商品，可能是父分类，尝试获取子分类商品
    if (goods.length === 0) {
      const catResp = await get('/goods/category/list');
      const allCategories = catResp.data || [];
      
      const parentCat = allCategories.find(cat => cat._id === cateId);
      if (parentCat && parentCat.child_cate && parentCat.child_cate.length > 0) {
        let allGoods = [];
        for (const child of parentCat.child_cate) {
          const childGoodsResp = await get(`/goods/list?categoryId=${child._id}`);
          const childGoods = childGoodsResp.data?.records || [];
          allGoods = allGoods.concat(childGoods);
        }
        return { spu: allGoods };
      }
    }
    
    return { spu: goods };
  } catch (error) {
    console.error('获取分类商品失败:', error);
    return { spu: [] };
  }
}

export async function getCates() {
  if (cloudbaseTemplateConfig.useMock) {
    return CATEGORY.filter((x) => x.child_cate?.length > 0);
  }

  const resp = await get('/goods/category/list');
  const categories = resp.data || [];
  
  // 添加"全部"选项到开头
  const allCategory = {
    _id: 'all',
    name: '全部',
    icon: '📦',
    image: '',
  };
  
  return [allCategory, ...categories];
}

