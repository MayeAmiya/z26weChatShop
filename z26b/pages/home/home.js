/* eslint-disable no-param-reassign */
import { getHomeSwiper, getHomeContent } from '../../services/home/home';
import { getCloudImageTempUrl } from '../../utils/cloudImageHandler';
import { getPrice } from '../../services/good/spu';
import { get } from '../../services/_utils/request';

Page({
  data: {
    imgSrcs: [],
    homeContent: null,
    contentParts: [], // 解析后的内容片段
    pageLoading: false,
    current: 1,
    autoplay: true,
    duration: '500',
    interval: 5000,
    navigation: { type: 'dots' },
    swiperImageProps: { mode: 'scaleToFill' }
  },

  onShow() {
    this.getTabBar().init();
  },

  onLoad() {
    this.init();
  },

  onPullDownRefresh() {
    this.init();
  },

  async init() {
    wx.stopPullDownRefresh();

    this.setData({
      pageLoading: true,
    });

    await Promise.all([
      this.loadHomeSwiper(),
      this.loadHomeContent(),
    ]);

    this.setData({
      pageLoading: false,
    });
  },

  async loadHomeSwiper() {
    const { images } = await getHomeSwiper();
    const handledImages = await getCloudImageTempUrl(images);

    this.setData({ imgSrcs: handledImages });
  },

  async loadHomeContent() {
    const content = await getHomeContent('main');
    this.setData({ homeContent: content });
    
    if (content && content.content) {
      await this.parseContent(content.content);
    }
  },

  // 解析内容，提取商品标记
  async parseContent(htmlContent) {
    // 匹配 [商品:ID:名称] 格式
    const productRegex = /\[商品:([a-zA-Z0-9-]+):([^\]]+)\]/g;
    const parts = [];
    let lastIndex = 0;
    let match;

    while ((match = productRegex.exec(htmlContent)) !== null) {
      // 添加商品标记前的HTML内容
      if (match.index > lastIndex) {
        parts.push({
          type: 'html',
          content: htmlContent.substring(lastIndex, match.index)
        });
      }
      
      // 添加商品占位
      parts.push({
        type: 'product',
        productId: match[1],
        productName: match[2],
        loading: true
      });
      
      lastIndex = match.index + match[0].length;
    }

    // 添加剩余的HTML内容
    if (lastIndex < htmlContent.length) {
      parts.push({
        type: 'html',
        content: htmlContent.substring(lastIndex)
      });
    }

    this.setData({ contentParts: parts });

    // 加载商品详情
    await this.loadProductDetails(parts);
  },

  // 加载商品详情
  async loadProductDetails(parts) {
    const productParts = parts.filter(p => p.type === 'product');
    
    for (let i = 0; i < productParts.length; i++) {
      const part = productParts[i];
      try {
        const response = await get(`/goods/${part.productId}`);
        const product = response.data;
        const price = await getPrice(part.productId);
        
        // 找到对应的 part 并更新
        const partIndex = parts.findIndex(p => p.type === 'product' && p.productId === part.productId);
        if (partIndex >= 0) {
          parts[partIndex] = {
            ...parts[partIndex],
            loading: false,
            product: {
              _id: product._id || part.productId,
              name: product.name,
              cover_image: product.cover_image || '',
              price: price
            }
          };
          // 每次更新后立即同步到视图
          this.setData({ contentParts: [...parts] });
        }
      } catch (e) {
        console.error('加载商品失败:', part.productId, e);
        const partIndex = parts.findIndex(p => p.type === 'product' && p.productId === part.productId);
        if (partIndex >= 0) {
          parts[partIndex].loading = false;
          parts[partIndex].error = true;
        }
      }
    }
    
    this.setData({ contentParts: [...parts] });
  },

  // 点击商品卡片
  onProductClick(e) {
    const { productId } = e.currentTarget.dataset;
    if (productId) {
      wx.navigateTo({
        url: `/pages/goods/details/index?spuId=${productId}`,
      });
    }
  },

  navToSearchPage() {
    wx.navigateTo({
      url: '/pages/goods/search/index',
    });
  },
});
