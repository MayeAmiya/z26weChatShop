import Toast from 'tdesign-miniprogram/toast/index';
import { getCates, getAllSpuOfCate } from '../../../services/cate/cate';
import { getPrice } from '../../../services/good/spu';

Page({
  data: {
    cates: [],
    goodsList: [],
    currentCateId: null,
    loading: false,
  },
  async init() {
    try {
      const cates = await getCates();
      this.setData({ cates });
      // 默认加载第一个分类的商品
      if (cates.length > 0) {
        this.loadGoodsByCate(cates[0]._id);
      }
    } catch (e) {
      console.error('获取商品分类列表失败', e);
      Toast({
        context: this,
        selector: '#t-toast',
        message: '获取商品分类列表失败',
        duration: 1000,
        icon: '',
      });
    }
  },

  // 加载分类商品
  async loadGoodsByCate(cateId) {
    if (this.data.currentCateId === cateId) return;
    
    this.setData({ loading: true, currentCateId: cateId });
    
    try {
      const { spu: goodsList } = await getAllSpuOfCate(cateId);
      
      if (goodsList && goodsList.length > 0) {
        // 优先使用 SPU 的 minPrice，如果没有则单独查询
        await Promise.all(goodsList.map(async (spu) => {
          if (spu.minPrice != null && spu.minPrice > 0) {
            spu.price = spu.minPrice;
          } else {
            spu.price = await getPrice(spu._id);
          }
        }));
      }
      
      this.setData({ goodsList: goodsList || [], loading: false });
    } catch (e) {
      console.error('获取分类商品失败', e);
      this.setData({ goodsList: [], loading: false });
    }
  },

  onShow() {
    this.getTabBar().init();
  },
  
  // 切换分类
  onCateChange(e) {
    const cateId = e?.detail?.cateId;
    if (cateId) {
      this.loadGoodsByCate(cateId);
    }
  },
  
  // 点击商品
  onGoodsClick(e) {
    const { goods } = e.detail;
    const spuId = goods._id || goods.spuId;
    if (spuId) {
      wx.navigateTo({
        url: `/pages/goods/details/index?spuId=${spuId}`,
      });
    }
  },
  
  onLoad() {
    this.init(true);
  },
});
