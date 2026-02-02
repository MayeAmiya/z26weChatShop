import { config } from '../../config/index';
import { cloudbaseTemplateConfig } from '../../config/index';
import { CART_ITEM, SKU, createId } from '../cloudbaseMock/index';
import { get, post, put, del } from '../_utils/request';
import { validateQuantity, debounceSubmit } from '../../utils/security';

// 防重复提交检查
const canAddToCart = debounceSubmit('addToCart', 500);
const canUpdateCart = debounceSubmit('updateCart', 300);

/** 获取购物车mock数据 */
function mockFetchCartGroupData(params) {
  const { delay } = require('../_utils/delay');
  const { genCartGroupData } = require('../../model/cart');

  return delay().then(() => genCartGroupData(params));
}

/**
 *
 * @param {{id: string}} param0
 * @returns
 */
export async function getCartItem({ id }) {
  if (cloudbaseTemplateConfig.useMock) {
    const cartItem = CART_ITEM.find((x) => x._id === id);
    cartItem.sku = SKU.find((sku) => sku._id === cartItem.sku._id);
    return { data: cartItem };
  }

  try {
    const response = await get(`/cart/items`);
    const storeGoods = response.data?.storeGoods || [];
    const allItems = storeGoods.reduce((acc, store) => acc.concat(store.goodsList || []), []);
    const item = allItems.find(x => x._id === id);
    return { data: item };
  } catch (error) {
    console.error('获取购物车项失败:', error);
    return { data: null };
  }
}

export async function fetchCartItems() {
  if (cloudbaseTemplateConfig.useMock) {
    const items = CART_ITEM.map((cartItem) => {
      const sku = SKU.find((x) => x._id === cartItem.sku._id);
      return {
        ...cartItem,
        sku,
      };
    });
    return items;  // 直接返回数组
  }

  try {
    const response = await get('/cart/items');
    // 后端返回 { data: { isNotEmpty, storeGoods: [{goodsList}] } }
    // 提取所有商品列表
    const storeGoods = response.data?.storeGoods || [];
    const allItems = storeGoods.reduce((acc, store) => {
      return acc.concat(store.goodsList || []);
    }, []);
    return allItems;  // 直接返回数组
  } catch (error) {
    console.error('获取购物车失败:', error);
    return [];  // 返回空数组
  }
}

/**
 * 添加商品到购物车
 * @param {{skuId: string, count: Number}} param0
 */
export async function createCartItem({ skuId, count }) {
  // 验证数量
  const quantityResult = validateQuantity(count, 99);
  if (!quantityResult.valid) {
    throw new Error(quantityResult.message);
  }
  
  // 防重复提交
  if (!canAddToCart()) {
    throw new Error('操作过于频繁，请稍后重试');
  }
  
  if (cloudbaseTemplateConfig.useMock) {
    CART_ITEM.push({ sku: { _id: skuId }, count, _id: createId() });
    return;
  }
  
  try {
    // 注意：后端使用 skuId 和 quantity 字段名
    const response = await post('/cart/add', { skuId: skuId, quantity: count });
    return response.data;
  } catch (error) {
    console.error('添加到购物车失败:', error);
    throw error;
  }
}

/**
 *
 * @param {{cartItemId: string}} param0
 */
export async function deleteCartItem({ cartItemId }) {
  if (cloudbaseTemplateConfig.useMock) {
    CART_ITEM.splice(
      CART_ITEM.findIndex((cartItem) => cartItem._id === cartItemId),
      1,
    );
    return;
  }
  
  try {
    await del(`/cart/remove/${cartItemId}`);
  } catch (error) {
    console.error('删除购物车项失败:', error);
    throw error;
  }
}

/**
 * 更新购物车商品数量
 * @param {{cartItemId: String, count: Number}} param0
 */
export async function updateCartItemCount({ cartItemId, count }) {
  // 验证数量
  const quantityResult = validateQuantity(count, 99);
  if (!quantityResult.valid) {
    throw new Error(quantityResult.message);
  }
  
  // 防重复提交
  if (!canUpdateCart()) {
    return; // 静默忽略频繁的更新请求
  }
  
  if (cloudbaseTemplateConfig.useMock) {
    CART_ITEM.find((x) => x._id === cartItemId).count = count;
    return;
  }
  
  try {
    // 注意：后端使用 quantity 字段名
    const response = await put(`/cart/update/${cartItemId}`, { quantity: count });
    return response.data;
  } catch (error) {
    console.error('更新购物车失败:', error);
    throw error;
  }
}

/** 获取购物车数据 */
export function fetchCartGroupData(params) {
  if (config.useMock) {
    return mockFetchCartGroupData(params);
  }

  return new Promise((resolve) => {
    resolve('real api');
  });
}

/**
 * 添加商品到购物车
 * @param {{skuId: string, quantity: number}} param0
 */
export async function addToCart({ skuId, quantity }) {
  if (cloudbaseTemplateConfig.useMock) {
    const existingItem = CART_ITEM.find(x => x.sku._id === skuId);
    if (existingItem) {
      existingItem.count += quantity;
      return { data: existingItem };
    }
    const newItem = { sku: { _id: skuId }, count: quantity, _id: createId(), isSelected: true };
    CART_ITEM.push(newItem);
    return { data: newItem };
  }
  
  try {
    const response = await post('/cart/add', { skuId, quantity });
    return response;
  } catch (error) {
    console.error('添加到购物车失败:', error);
    throw error;
  }
}

/**
 * 更新购物车项的选中状态
 * @param {{cartItemId: string, isSelected: boolean}} param0
 */
export async function updateCartItemSelected({ cartItemId, isSelected }) {
  if (cloudbaseTemplateConfig.useMock) {
    const item = CART_ITEM.find(x => x._id === cartItemId);
    if (item) {
      item.isSelected = isSelected;
    }
    return;
  }
  
  try {
    const response = await put(`/cart/update/${cartItemId}`, { isSelected });
    return response;
  } catch (error) {
    console.error('更新购物车选中状态失败:', error);
    throw error;
  }
}
