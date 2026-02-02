import { config } from '../../config/index';
import { get } from '../_utils/request';

/** 获取订单列表mock数据 */
function mockFetchOrders(params) {
  const { delay } = require('../_utils/delay');
  const { genOrders } = require('../../model/order/orderList');

  return delay(200).then(() => genOrders(params));
}

/** 获取订单列表数据 */
export async function fetchOrders(params) {
  if (config.useMock) {
    return mockFetchOrders(params);
  }

  try {
    const response = await get('/order/list', params);
    // 后端返回格式: { data: { records, total, page, pageSize } }
    return {
      orders: response.data?.records || [],
      totalCount: response.data?.total || 0,
    };
  } catch (error) {
    console.error('获取订单列表失败:', error);
    throw error;
  }
}

/** 获取订单列表mock数据 */
function mockFetchOrdersCount(params) {
  const { delay } = require('../_utils/delay');
  const { genOrdersCount } = require('../../model/order/orderList');

  return delay().then(() => genOrdersCount(params));
}

/** 获取订单列表统计 - 通过分别查询各状态数量实现 */
export async function fetchOrdersCount(params) {
  if (config.useMock) {
    return mockFetchOrdersCount(params);
  }

  try {
    // 分别获取各状态的订单数量
    const statuses = ['TO_PAY', 'TO_SEND', 'TO_RECEIVE', 'COMPLETED', 'CANCELED'];
    const counts = {};
    
    for (const status of statuses) {
      const response = await get('/order/list', { status, pageSize: 1 });
      counts[status] = response.data?.total || 0;
    }
    
    return {
      orderStatusCounts: {
        unpaid: counts['TO_PAY'],
        undelivered: counts['TO_SEND'],
        unreceived: counts['TO_RECEIVE'],
        completed: counts['COMPLETED'],
        canceled: counts['CANCELED'],
      }
    };
  } catch (error) {
    console.error('获取订单统计失败:', error);
    return {
      orderStatusCounts: {
        unpaid: 0,
        undelivered: 0,
        unreceived: 0,
        completed: 0,
        canceled: 0,
      }
    };
  }
}
