import { cloudbaseTemplateConfig } from '../../config/index';
import { SKU, SPU, ATTR_VALUE } from '../cloudbaseMock/index';
import { get, put } from '../_utils/request';

/**
 *
 * @param {String} skuId
 */
export async function getSkuDetail(skuId) {
  if (cloudbaseTemplateConfig.useMock) {
    const sku = SKU.find((x) => x._id === skuId);
    sku.attr_value = ATTR_VALUE.filter((x) => x.sku.find((x) => x._id === skuId));
    sku.spu = SPU.find((spu) => spu._id === sku.spu._id);
    return sku;
  }

  try {
    const response = await get(`/sku/${skuId}`);
    const sku = response.data;
    // 映射后端字段
    return {
      ...sku,
      _id: sku._id || sku.id,
      price: sku.price || 0,
      count: sku.count || 0,
      image: sku.image || '',
      description: sku.description || '',
    };
  } catch (error) {
    console.error('获取SKU详情失败:', error);
    throw error;
  }
}

export async function updateSku({ skuId, data }) {
  if (cloudbaseTemplateConfig.useMock) {
    SKU.find((x) => x._id === skuId).count = data.count;
    return;
  }

  try {
    const response = await put(`/sku/${skuId}`, data);
    return response.data;
  } catch (error) {
    console.error('更新SKU失败:', error);
    throw error;
  }
}

export async function getAllSku(spuId) {
  if (cloudbaseTemplateConfig.useMock) {
    return SKU.filter((x) => x.spu._id === spuId);
  }
  
  try {
    const response = await get(`/sku/list/${spuId}`);
    const skus = response.data || [];
    // 确保字段正确映射
    return skus.map(sku => ({
      ...sku,
      _id: sku._id || sku.id,
      price: Number(sku.price) || 0,
      count: Number(sku.count) || 0,
      image: sku.image || '',
      description: sku.description || '',
    }));
  } catch (error) {
    console.error('获取SKU列表失败:', error);
    return [];
  }
}
