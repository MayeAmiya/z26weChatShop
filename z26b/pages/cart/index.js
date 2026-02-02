import Dialog from 'tdesign-miniprogram/dialog/index';
import { fetchCartGroupData, fetchCartItems, deleteCartItem, updateCartItemCount, updateCartItemSelected } from '../../services/cart/cart';
import { getSkuDetail } from '../../services/sku/sku';
import { objectToParamString } from '../../utils/util';
import { getSingleCloudImageTempUrl } from '../../utils/cloudImageHandler';

let updateCartItemCountTimeout = null;

Component({
  data: {
    cartGroupData: null,
    cartItems: [],
    selectedCartItemNum: 0,
    allSelected: false,
    totalAmount: 0,
    loading: false,
    inited: false,
  },

  observers: {
    cartItems: function (cartItems) {
      const selectedCartItems = cartItems.filter((x) => x.selected === true);
      const selectedCartItemNum = selectedCartItems.length;
      // 兼容 count 和 quantity 字段，以及 sku 可能为空的情况
      const totalAmount = selectedCartItems.reduce((acc, cur) => {
        const count = cur.count || cur.quantity || 1;
        const price = cur.sku?.price || 0;
        return acc + count * price;
      }, 0);
      const allSelected = cartItems.length > 0 && selectedCartItemNum === cartItems.length;

      this.setData({
        selectedCartItemNum,
        allSelected,
        totalAmount,
      });
    },
  },

  lifetimes: {
    attached: async function () {
      // console.log('called attached');
      // // 调用自定义tabbar的init函数，使页面与tabbar激活状态保持一致
      // this.getTabBar().init();
      // await this.setLoading();
      // await this.setDataPromise({ inited: true });
      // try {
      //   await this.init();
      // } finally {
      //   await this.unsetLoading();
      // }
    },
  },

  pageLifetimes: {
    show: async function () {
      // 调用自定义tabbar的init函数，使页面与tabbar激活状态保持一致
      this.getTabBar().init();

      await this.setLoading();
      try {
        await this.init();
      } finally {
        await this.unsetLoading();
      }
    },
  },

  methods: {
    async init() {
      const cartItems = (await fetchCartItems()).map((x) =>
        Object.assign(x, {
          selected: x.isSelected || false,
        }),
      );
      await Promise.all(
        cartItems.map(async (cartItem) => {
          // 兼容：从 sku 对象或 skuId 字段获取 SKU ID
          const skuId = cartItem.sku?._id || cartItem.skuId;
          if (!skuId) {
            console.warn('购物车项缺少 SKU ID:', cartItem);
            return;
          }
          const sku = await getSkuDetail(skuId);
          if (sku.image) {
            sku.image = await getSingleCloudImageTempUrl(sku.image);
          }
          cartItem.sku = sku;
          // 确保 count 字段存在（兼容 quantity 字段）
          cartItem.count = cartItem.count || cartItem.quantity || 1;
        }),
      );
      await this.setDataPromise({
        cartItems,
      });
    },

    findGoods(spuId, skuId) {
      let currentStore;
      let currentActivity;
      let currentGoods;
      const { storeGoods } = this.data.cartGroupData;
      for (const store of storeGoods) {
        for (const activity of store.promotionGoodsList) {
          for (const goods of activity.goodsPromotionList) {
            if (goods.spuId === spuId && goods.skuId === skuId) {
              currentStore = store;
              currentActivity = currentActivity;
              currentGoods = goods;
              return {
                currentStore,
                currentActivity,
                currentGoods,
              };
            }
          }
        }
      }
      return {
        currentStore,
        currentActivity,
        currentGoods,
      };
    },

    // 注：实际场景时应该调用接口获取购物车数据
    getCartGroupData() {
      const { cartGroupData } = this.data;
      if (!cartGroupData) {
        return fetchCartGroupData();
      }
      return Promise.resolve({ data: cartGroupData });
    },

    // 选择单个商品
    // 注：实际场景时应该调用接口更改选中状态
    selectGoodsService({ spuId, skuId, isSelected }) {
      this.findGoods(spuId, skuId).currentGoods.isSelected = isSelected;
      return Promise.resolve();
    },

    // 全选门店
    // 注：实际场景时应该调用接口更改选中状态
    selectStoreService({ storeId, isSelected }) {
      const currentStore = this.data.cartGroupData.storeGoods.find((s) => s.storeId === storeId);
      currentStore.isSelected = isSelected;
      currentStore.promotionGoodsList.forEach((activity) => {
        activity.goodsPromotionList.forEach((goods) => {
          goods.isSelected = isSelected;
        });
      });
      return Promise.resolve();
    },

    // 加购数量变更
    // 注：实际场景时应该调用接口
    changeQuantityService({ spuId, skuId, quantity }) {
      this.findGoods(spuId, skuId).currentGoods.quantity = quantity;
      return Promise.resolve();
    },

    // 删除加购商品
    // 注：实际场景时应该调用接口
    deleteGoodsService({ spuId, skuId }) {
      function deleteGoods(group) {
        for (const gindex in group) {
          const goods = group[gindex];
          if (goods.spuId === spuId && goods.skuId === skuId) {
            group.splice(gindex, 1);
            return gindex;
          }
        }
        return -1;
      }
      const { storeGoods, invalidGoodItems } = this.data.cartGroupData;
      for (const store of storeGoods) {
        for (const activity of store.promotionGoodsList) {
          if (deleteGoods(activity.goodsPromotionList) > -1) {
            return Promise.resolve();
          }
        }
        if (deleteGoods(store.shortageGoodsList) > -1) {
          return Promise.resolve();
        }
      }
      if (deleteGoods(invalidGoodItems) > -1) {
        return Promise.resolve();
      }
      return Promise.reject();
    },

    onGoodsSelect({ detail: { goods } }) {
      const { cartItems } = this.data;
      const item = cartItems.find((x) => x._id === goods._id);

      if (item == null) {
        console.warn('Cart item not found!');
        return;
      }

      item.selected = !item.selected;
      this.setData({ cartItems });
      
      // 同步选中状态到后端
      updateCartItemSelected({ 
        cartItemId: goods._id, 
        isSelected: item.selected 
      }).catch(e => console.error('同步选中状态失败:', e));
    },

    onQuantityChange({ detail: { cartItemId, count } }) {
      const { cartItems } = this.data;
      const item = cartItems.find((x) => x._id === cartItemId);
      if (item == null) {
        console.warn('Cart item not found');
        return;
      }
      item.count = count;
      this.setData({ cartItems });
      this.debouncedUpdateCartItemCount({ cartItemId, count });
    },

    debouncedUpdateCartItemCount({ cartItemId, count }) {
      clearTimeout(updateCartItemCountTimeout);
      updateCartItemCountTimeout = setTimeout(async () => {
        this.setLoading();
        try {
          await updateCartItemCount({ cartItemId, count });
        } finally {
          this.unsetLoading();
        }
      }, 500);
    },

    goCollect() {
      /** 活动肯定有一个活动ID，用来获取活动banner，活动商品列表等 */
      const promotionID = '123';
      wx.navigateTo({
        url: `/pages/promotion-detail/index?promotion_id=${promotionID}`,
      });
    },

    goGoodsDetail({
      detail: {
        goods: {
          sku: {
            spu: { _id },
          },
        },
      },
    }) {
      wx.navigateTo({
        url: `/pages/goods/details/index?spuId=${_id}`,
      });
    },

    async onGoodsDelete({
      detail: {
        goods: { _id },
      },
    }) {
      try {
        await Dialog.confirm({
          context: this,
          closeOnOverlayClick: true,
          title: '确认删除该商品吗?',
          confirmBtn: '确定',
          cancelBtn: '取消',
        });
        this.setLoading();
        try {
          await deleteCartItem({ cartItemId: _id });
          const { cartItems } = this.data;
          this.setData({
            cartItems: cartItems.filter((x) => x._id !== _id),
          });
        } finally {
          this.unsetLoading();
        }
      } catch {
        console.warn('deletion is cancelled');
      }
    },

    onSelectAll() {
      const { cartItems, allSelected } = this.data;
      const newSelectedState = !allSelected;
      cartItems.forEach((x) => (x.selected = newSelectedState));
      this.setData({ cartItems });
      
      // 同步所有选中状态到后端
      cartItems.forEach(item => {
        updateCartItemSelected({ 
          cartItemId: item._id, 
          isSelected: newSelectedState 
        }).catch(e => console.error('同步选中状态失败:', e));
      });
    },

    async onToSettle() {
      const selectedItems = this.data.cartItems.filter((x) => x.selected === true);
      
      if (selectedItems.length === 0) {
        wx.showToast({
          title: '请选择商品',
          icon: 'none',
        });
        return;
      }
      
      // 确保选中状态已同步到后端
      this.setLoading();
      try {
        await Promise.all(
          selectedItems.map(item => 
            updateCartItemSelected({ 
              cartItemId: item._id, 
              isSelected: true 
            })
          )
        );
      } catch (e) {
        console.error('同步选中状态失败:', e);
      } finally {
        this.unsetLoading();
      }
      
      wx.navigateTo({
        url: `/pages/order/order-confirm/index?${objectToParamString({
          type: 'cart',
          cartIds: selectedItems.map((x) => x._id).join(','),
        })}`,
      });
    },
    onGotoHome() {
      wx.switchTab({ url: '/pages/home/home' });
    },

    setLoading() {
      return this.setDataPromise({
        loading: true,
      });
    },
    unsetLoading() {
      return this.setDataPromise({
        loading: false,
      });
    },
    toggleLoading() {
      this.setData({
        loading: !this.data.loading,
      });
    },

    setDataPromise(data) {
      return new Promise((res) => {
        this.setData(data, () => res());
      });
    },
  },
});
