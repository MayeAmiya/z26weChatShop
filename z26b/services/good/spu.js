import { getCloudImageTempUrl } from '../../utils/cloudImageHandler';
import { SPU_SELLING_STATUS } from '../../utils/spuStatus';
import { cloudbaseTemplateConfig } from '../../config/index';
import { SPU, SKU } from '../cloudbaseMock/index';
import { get as apiGet } from '../_utils/request';

/**
 *
 * @param {{
 *   pageSize: Number,
 *   pageNumber: Number,
 *   cateId: String,
 *   search: String
 * }}} param0
 * @returns
 */
export async function listGood({ pageSize, pageNumber, search }) {
  if (cloudbaseTemplateConfig.useMock) {
    const records = search ? SPU.filter((x) => x.name.includes(search)) : SPU;
    const startIndex = (pageNumber - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return {
      records: records.slice(startIndex, endIndex),
      total: records.length,
    };
  }
  
  // 使用后端API
  try {
    const params = {
      page: pageNumber,
      pageSize: pageSize,
    };
    if (search) {
      params.keyword = search;
    }
    const response = await apiGet('/goods/list', params);
    return {
      records: response.data?.records || [],
      total: response.data?.total || 0,
    };
  } catch (error) {
    console.error('获取商品列表失败:', error);
    return { records: [], total: 0 };
  }
}

export async function getPrice(spuId) {
  if (cloudbaseTemplateConfig.useMock) {
    return SKU.find((x) => x.spu._id === spuId).price;
  }
  try {
    // 先尝试从 SPU 获取 minPrice（如果后端支持）
    const spuResponse = await apiGet(`/goods/${spuId}`);
    if (spuResponse.data?.minPrice) {
      return spuResponse.data.minPrice;
    }
    // 降级：从 SKU 列表获取最低价
    const response = await apiGet(`/sku/list/${spuId}`);
    const skus = response.data || [];
    if (skus.length > 0) {
      return Math.min(...skus.map(s => s.price));
    }
    return 0;
  } catch (error) {
    console.error('获取商品价格失败:', error);
    return 0;
  }
}

export async function handleSpuCloudImage(spu) {
  if (spu == null) {
    return;
  }
  
  // Safely handle undefined/null fields and potential JSON strings
  let swiperImages = spu.swiper_images;
  if (typeof swiperImages === 'string') {
    try {
      swiperImages = JSON.parse(swiperImages);
    } catch (e) {
      swiperImages = [];
    }
  }
  swiperImages = Array.isArray(swiperImages) ? swiperImages : [];
  
  // 确保 swiper_images 包含封面图作为第一张
  const coverImage = spu.cover_image;
  if (coverImage && !swiperImages.includes(coverImage)) {
    swiperImages = [coverImage, ...swiperImages];
  }
  
  const allImages = swiperImages.filter(img => img);
  
  if (allImages.length === 0) {
    spu.swiper_images = [];
    return;
  }
  
  const handledImages = await getCloudImageTempUrl(allImages);
  
  // 更新封面图（第一张）
  if (handledImages.length > 0) {
    spu.cover_image = handledImages[0];
  }
  
  // 更新轮播图数组
  spu.swiper_images = handledImages;
}

export async function getSpu(spuId) {
  if (cloudbaseTemplateConfig.useMock) {
    return SPU.find((x) => x._id === spuId);
  }
  try {
    const response = await apiGet(`/goods/${spuId}`);
    return response.data;
  } catch (error) {
    console.error('获取商品详情失败:', error);
    return null;
  }
}

