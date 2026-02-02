import { get, post, put } from '../../services/_utils/request';
import { cloudbaseTemplateConfig } from '../../config/index';
import { ORDER, createId, DELIVERY_INFO } from '../cloudbaseMock/index';

const ORDER_STATUS_INFO = {
  TO_PAY: { value: 'TO_PAY', label: '待付款' },
  TO_SEND: { value: 'TO_SEND', label: '待发货' },
  TO_RECEIVE: { value: 'TO_RECEIVE', label: '待收货' },
  FINISHED: { value: 'FINISHED', label: '已完成' },
  CANCELED: { value: 'CANCELED', label: '已取消' },
  RETURN_APPLIED: { value: 'RETURN_APPLIED', label: '申请退货' },
  RETURN_REFUSED: { value: 'RETURN_REFUSED', label: '拒绝退货申请' },
  RETURN_FINISH: { value: 'RETURN_FINISH', label: '退货完成' },
  RETURN_MONEY_REFUSED: { value: 'RETURN_MONEY_REFUSED', label: '拒绝退款' },
};

export const ORDER_STATUS = new Proxy(ORDER_STATUS_INFO, {
  get(target, prop) {
    return target[prop]?.value;
  },
});

export const orderStatusToName = (status) => Object.values(ORDER_STATUS_INFO).find((x) => x.value === status)?.label;

/**
 * 创建订单（已废弃 - 后端在支付时直接创建订单）
 * 现在使用 payWithBalance 直接创建订单并支付
 * @deprecated
 */
export async function createOrder({ status, addressId }) {
  if (cloudbaseTemplateConfig.useMock) {
    const _id = createId();
    ORDER.push({
      status,
      delivery_info: {
        _id: addressId,
      },
      _id,
      createdAt: new Date().getTime()
    });
    return { id: _id };
  }
  // 实际上不再需要单独创建订单，后端在支付时创建
  throw new Error('请使用 payWithBalance 创建订单');
}

/**
 * 获取所有订单
 */
export async function getAllOrder() {
  if (cloudbaseTemplateConfig.useMock) {
    return ORDER;
  }
  const response = await get('/order/list', { pageSize: 100 });
  return response.data?.records || [];
}

/**
 *
 * @param {{
 *   pageSize: Number,
 *   pageNumber: Number,
 *   status?: String
 * }}} param0
 * @returns
 */
export async function listOrder({ pageSize, pageNumber, status }) {
  if (cloudbaseTemplateConfig.useMock) {
    const filteredOrder = status == null ? ORDER : ORDER.filter((x) => x.status === status);
    const startIndex = (pageNumber - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    const records = filteredOrder.slice(startIndex, endIndex);
    const total = filteredOrder.length;
    return {
      records,
      total,
    };
  }

  try {
    const params = { page: pageNumber, pageSize };
    if (status != null) {
      params.status = status;
    }
    const response = await get('/order/list', params);
    const records = (response.data?.records || []).map(order => {
      // 处理日期 - 后端返回的是秒级时间戳，需要转为毫秒
      if (order.createdAt && order.createdAt < 10000000000) {
        order.createdAt = order.createdAt * 1000;
      }
      return order;
    });
    return {
      records,
      total: response.data?.total || 0,
    };
  } catch (error) {
    console.error('获取订单列表失败:', error);
    return { records: [], total: 0 };
  }
}

async function getOrderCountOfStatus(status) {
  if (cloudbaseTemplateConfig.useMock) {
    return ORDER.filter((x) => x.status === status).length;
  }

  try {
    const response = await get('/order/list', { status, pageSize: 1 });
    return response.data?.total || 0;
  } catch (error) {
    console.error('获取订单数量失败:', error);
    return 0;
  }
}

export async function getToPayOrderCount() {
  return getOrderCountOfStatus(ORDER_STATUS.TO_PAY);
}

export async function getToSendOrderCount() {
  return getOrderCountOfStatus(ORDER_STATUS.TO_SEND);
}

export async function getToReceiveOrderCount() {
  return getOrderCountOfStatus(ORDER_STATUS.TO_RECEIVE);
}

/**
 * 获取订单详情
 * @param {String} orderId
 */
export async function getOrder(orderId) {
  if (cloudbaseTemplateConfig.useMock) {
    const order = ORDER.find(o => o._id === orderId);
    order.delivery_info = DELIVERY_INFO.find(i => i._id === order.delivery_info._id)
    return order
  }
  
  try {
    const response = await get(`/order/${orderId}`);
    const order = response.data;
    
    // 处理 delivery_info - 如果是字符串则解析
    if (order.delivery_info && typeof order.delivery_info === 'string') {
      try {
        order.delivery_info = JSON.parse(order.delivery_info);
      } catch (e) {
        console.error('解析 delivery_info 失败:', e);
      }
    }
    
    // 确保 delivery_info 有 address 字段（组合地址）
    if (order.delivery_info) {
      const di = order.delivery_info;
      di.address = [di.provinceName, di.cityName, di.districtName, di.detailAddress].filter(Boolean).join('');
    }
    
    // 处理日期 - 后端返回的是秒级时间戳，需要转为毫秒
    if (order.createdAt && order.createdAt < 10000000000) {
      order.createdAt = order.createdAt * 1000;
    }
    if (order.updatedAt && order.updatedAt < 10000000000) {
      order.updatedAt = order.updatedAt * 1000;
    }
    
    return order;
  } catch (error) {
    console.error('获取订单详情失败:', error);
    throw error;
  }
}

/**
 * 更新订单收货地址
 */
export async function updateOrderDeliveryInfo({ orderId, deliveryInfoId }) {
  if (cloudbaseTemplateConfig.useMock) {
    const order = ORDER.find(x => x._id === orderId);
    if (order) {
      order.delivery_info = { _id: deliveryInfoId };
    }
    return;
  }
  // 后端暂不支持单独更新收货地址
  console.warn('updateOrderDeliveryInfo: 后端暂不支持此功能');
}

/**
 * 更新订单状态
 * @param {{orderId: String, status: String}} param0
 * @returns
 */
export async function updateOrderStatus({ orderId, status }) {
  if (cloudbaseTemplateConfig.useMock) {
    ORDER.find(x => x._id === orderId).status = status
    return;
  }
  
  // 根据状态调用不同的后端接口
  try {
    if (status === ORDER_STATUS.CANCELED) {
      const response = await put(`/order/cancel/${orderId}`);
      return response;
    } else if (status === ORDER_STATUS.FINISHED) {
      const response = await post(`/order/confirm/${orderId}`);
      return response;
    } else {
      console.warn(`updateOrderStatus: 不支持直接更新到状态 ${status}`);
    }
  } catch (error) {
    console.error('更新订单状态失败:', error);
    throw error;
  }
}

/**
 * 取消订单
 * @param {String} orderId
 */
export async function cancelOrder(orderId) {
  try {
    const response = await put(`/order/cancel/${orderId}`);
    return response;
  } catch (error) {
    console.error('取消订单失败:', error);
    throw error;
  }
}

/**
 * 确认收货
 * @param {String} orderId
 */
export async function confirmReceipt(orderId) {
  try {
    const response = await post(`/order/confirm/${orderId}`);
    return response;
  } catch (error) {
    console.error('确认收货失败:', error);
    throw error;
  }
}
