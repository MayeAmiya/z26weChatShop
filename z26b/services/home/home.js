import { cloudbaseTemplateConfig } from '../../config/index';
import { HOME_SWIPER } from '../cloudbaseMock/index';
import { get } from '../_utils/request';

export async function getHomeSwiper() {
  if (cloudbaseTemplateConfig.useMock) {
    return HOME_SWIPER[0];
  }

  try {
    const response = await get('/home/swiper');
    const data = response.data;
    // 兼容处理：确保返回的数据格式正确
    if (data && data.images) {
      return data;
    }
    // 如果后端返回的是数组格式的图片
    if (Array.isArray(data)) {
      return { images: data };
    }
    return { images: [] };
  } catch (error) {
    console.error('获取首页轮播失败:', error);
    return { images: [] };
  }
}

// 获取首页富文本内容
export async function getHomeContent(key = 'main') {
  if (cloudbaseTemplateConfig.useMock) {
    return { title: '欢迎', content: '<p>欢迎使用商城</p>' };
  }

  try {
    const response = await get(`/home/content?key=${key}`);
    return response.data || null;
  } catch (error) {
    console.error('获取首页内容失败:', error);
    return null;
  }
}
