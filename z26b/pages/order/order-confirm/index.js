import Toast from 'tdesign-miniprogram/toast/index';
import { getSkuDetail, updateSku } from '../../../services/sku/sku';
import { getAddressPromise } from '../../usercenter/address/list/util';
import { getSingleCloudImageTempUrl } from '../../../utils/cloudImageHandler';
import { cartShouldFresh } from '../../../utils/cartFresh';
import { payWithBalance } from '../../../services/pay/pay';
import { addToCart, deleteCartItem, fetchCartItems, updateCartItemSelected } from '../../../services/cart/cart';
import { getAllAddress } from '../../../services/address/address';

const stripeImg = `https://cdn-we-retail.ym.tencent.com/miniapp/order/stripe.png`;

/**
 * 将购物车项转换为商品列表用于展示
 */
function cartItemToGoodList(cartItems) {
  return cartItems.map((item) => ({
    thumb: item.sku?.image || item.image || '',
    title: item.sku?.spu?.name || item.title || '商品',
    specs: item.sku?.attr_value?.map((v) => v.value).join('，') || '',
    price: item.sku?.price || item.price || 0,
    num: item.quantity || item.count || 1,
  }));
}

Page({
  data: {
    placeholder: '备注信息',
    stripeImg,
    loading: true,
    orderCardList: [], // 仅用于商品卡片展示
    goodsRequestList: [],
    userAddressReq: null,
    storeInfoList: [],
    promotionGoodsList: [], //当前门店商品列表(优惠券)
    currentStoreId: null, //当前优惠券storeId
    userAddress: null,
    goodsList: [],
    cartItems: [],
    totalSalePrice: 0,
    directSku: null,
    directCount: 1,
  },

  payLock: false,

  type: null,

  /**
   * 从购物车加载商品
   * 后端会自动获取选中的购物车项，这里只是用于展示
   */
  async onLoadFromCart() {
    try {
      // 获取购物车中已选中的商品用于展示
      // fetchCartItems 直接返回数组
      let cartItems = await fetchCartItems();
      
      // 过滤选中的商品
      cartItems = cartItems.filter(item => item.isSelected || item.is_selected || item.selected);
      
      if (cartItems.length === 0) {
        this.failedAndBack('购物车中没有选中的商品');
        return;
      }
      
      // 处理图片和数据格式
      for (const item of cartItems) {
        // 获取完整的 SKU 信息
        const skuId = item.sku?._id || item.skuId;
        if (skuId) {
          const sku = await getSkuDetail(skuId);
          if (sku.image) {
            sku.image = await getSingleCloudImageTempUrl(sku.image);
          }
          item.sku = sku;
        }
        // 兼容字段名
        item.quantity = item.quantity || item.count || 1;
      }
      
      const goodsList = cartItemToGoodList(cartItems);
      const totalSalePrice = goodsList.reduce((acc, cur) => acc + cur.price * cur.num, 0);
      
      this.setData({
        goodsList,
        totalSalePrice,
        cartItems,
      });
    } catch (e) {
      this.failedAndBack('获取购物车信息失败', e);
    }
  },
  
  /**
   * 直接购买模式 - 需要先将商品加入购物车
   */
  async onLoadFromDirect(countStr, skuId) {
    const count = parseInt(countStr);
    if (typeof count !== 'number' || isNaN(count) || typeof skuId !== 'string') {
      console.error('invalid count or skuId', count, skuId);
      this.failedAndBack('初始化信息有误');
      return;
    }

    try {
      const sku = await getSkuDetail(skuId);
      sku.image = await getSingleCloudImageTempUrl(sku.image);

      const goodsList = [
        {
          thumb: sku.image,
          title: sku.spu?.name || '商品',
          specs: sku.attr_value?.map((v) => v.value).join('，') || '',
          price: sku.price,
          num: count,
        },
      ];

      const totalSalePrice = goodsList.reduce((acc, cur) => acc + cur.price * cur.num, 0);
      
      this.setData({
        goodsList,
        totalSalePrice,
        directSku: sku,
        directCount: count,
      });
    } catch (e) {
      this.failedAndBack('获取商品信息失败', e);
    }
  },

  async onLoad(options) {
    this.type = options?.type;
    if (this.type === 'cart') {
      await this.onLoadFromCart();
    } else if (this.type === 'direct') {
      await this.onLoadFromDirect(options?.count, options?.skuId);
    } else {
      this.failedAndBack('初始化信息有误', 'invalid type');
    }

    // 获取默认地址
    try {
      const addresses = await getAllAddress();
      if (addresses.length > 0) {
        // 优先选择默认地址，否则选第一个
        const defaultAddr = addresses.find(a => a.isDefault === 1) || addresses[0];
        this.setData({
          userAddress: {
            ...defaultAddr,
            detailAddress: defaultAddr.address || defaultAddr.detailAddress,
          },
        });
      }
    } catch (e) {
      console.error('获取地址失败:', e);
    }

    this.setData({
      loading: false,
    });
  },

  init() {
    this.setData({
      loading: false,
    });
    const { goodsRequestList } = this;
    this.handleOptionsParams({ goodsRequestList });
  },

  toast(message) {
    Toast({
      context: this,
      selector: '#t-toast',
      message,
      duration: 1000,
      icon: '',
    });
  },

  onGotoAddress() {
    /** 获取一个Promise */
    getAddressPromise()
      .then((address) => {
        this.setData({
          userAddress: {
            ...address,
            detailAddress: address.address || address.detailAddress,
          },
        });
      })
      .catch(() => {});

    wx.navigateTo({
      url: `/pages/usercenter/address/list/index?selectMode=true`,
    });
  },
  onTap() {
    this.setData({
      placeholder: '',
    });
  },

  /**
   * 从购物车提交订单 - 直接下单（暂不需要支付）
   * 后端会自动：获取选中的购物车项 -> 创建订单 -> 清除购物车
   */
  async submitOrderFromCart() {
    const { userAddress, totalSalePrice } = this.data;

    try {
      // 调用后端创建订单
      const response = await payWithBalance({
        addressId: userAddress._id,
        remarks: this.data.placeholder || '',
      });
      
      // 通知购物车刷新
      cartShouldFresh();
      
      this.toast('下单成功');
      setTimeout(() => {
        wx.navigateTo({
          url: `/pages/order/pay-result/index?orderId=${response.data?.order?._id || ''}&paidAmount=${response.data?.paidAmount || totalSalePrice}`,
        });
      }, 1000);
    } catch (e) {
      console.error('下单失败:', e);
      this.failedAndBack('下单失败', e);
    }
  },

  /**
   * 直接购买 - 需要先将商品加入购物车，选中，然后下单
   */
  async submitOrderFromDirect() {
    const { directSku, userAddress, totalSalePrice, directCount } = this.data;

    try {
      // 1. 将商品加入购物车并选中
      const addResult = await addToCart({
        skuId: directSku._id,
        quantity: directCount || 1,
      });
      
      // 2. 确保商品被选中
      if (addResult.data?._id) {
        await updateCartItemSelected({ 
          cartItemId: addResult.data._id, 
          isSelected: true 
        });
      }
      
      // 3. 调用后端创建订单
      const response = await payWithBalance({
        addressId: userAddress._id,
        remarks: this.data.placeholder || '',
      });
      
      // 通知购物车刷新
      cartShouldFresh();
      
      this.toast('下单成功');
      setTimeout(() => {
        wx.navigateTo({
          url: `/pages/order/pay-result/index?orderId=${response.data?.order?._id || ''}&paidAmount=${response.data?.paidAmount || totalSalePrice}`,
        });
      }, 1000);
    } catch (e) {
      console.error('下单失败:', e);
      this.failedAndBack('下单失败', e);
    }
  },

  failedAndBack(message, e) {
    e && console.error(e);
    this.toast(message);
    setTimeout(() => {
      wx.navigateBack();
    }, 1000);
  },

  // 提交订单
  async submitOrder() {
    const { userAddress } = this.data;
    if (!userAddress) {
      Toast({
        context: this,
        selector: '#t-toast',
        message: '请添加收货地址',
        duration: 2000,
        icon: 'help-circle',
      });
      return;
    }

    if (this.type === 'cart') {
      this.submitOrderFromCart();
    } else if (this.type === 'direct') {
      this.submitOrderFromDirect();
    } else {
      console.error('invalid type', this.type);
      this.failedAndBack('初始化信息有误');
    }
  },
});
