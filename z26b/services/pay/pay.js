/**
 * 支付服务 - 使用钱包余额支付
 */

import { post, get, put } from '../_utils/request';

/**
 * 创建订单并支付（使用余额）
 * 后端会从购物车中获取已选中的商品，计算总价，检查余额，创建订单并扣款
 * 
 * @param {Object} params
 * @param {string} params.addressId - 收货地址ID
 * @param {string} [params.remarks] - 订单备注
 * @returns {Promise<{data: {order: Object, remainBalance: number, paidAmount: number, paymentMethod: string}}>}
 */
export async function payWithBalance({ addressId, remarks = '' }) {
  try {
    const response = await post('/order/create', {
      addressId,
      remarks,
    });
    return response;
  } catch (error) {
    console.error('余额支付失败:', error);
    throw error;
  }
}

/**
 * 获取用户余额
 * @returns {Promise<{data: {balance: number}}>}
 */
export async function getBalance() {
  try {
    const response = await get('/user/balance');
    return response;
  } catch (error) {
    console.error('获取余额失败:', error);
    throw error;
  }
}

/**
 * 申请退款（取消订单）
 * @param {string} orderId - 订单ID
 * @returns {Promise<Object>}
 */
export async function refund(orderId) {
  try {
    const response = await put(`/order/cancel/${orderId}`);
    return response;
  } catch (error) {
    console.error('取消订单失败:', error);
    throw error;
  }
}

