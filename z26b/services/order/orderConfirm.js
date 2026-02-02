import { config } from '../../config/index';
import { post } from '../_utils/request';
import { mockIp, mockReqId } from '../../utils/mock';
import { sanitizeInput, debounceSubmit } from '../../utils/security';

// 防重复提交
const canSubmitOrder = debounceSubmit('submitOrder', 3000);

/** 获取结算mock数据 */
function mockFetchSettleDetail(params) {
  const { delay } = require('../_utils/delay');
  const { genSettleDetail } = require('../../model/order/orderConfirm');

  return delay().then(() => genSettleDetail(params));
}

/** 提交mock订单 */
function mockDispatchCommitPay() {
  const { delay } = require('../_utils/delay');

  return delay().then(() => ({
    data: {
      isSuccess: true,
      tradeNo: '350930961469409099',
      payInfo: '{}',
      code: null,
      transactionId: 'E-200915180100299000',
      msg: null,
      interactId: '15145',
      channel: 'wechat',
      limitGoodsList: null,
    },
    code: 'Success',
    msg: null,
    requestId: mockReqId(),
    clientIp: mockIp(),
    rt: 891,
    success: true,
  }));
}

/** 获取结算数据 */
export function fetchSettleDetail(params) {
  if (config.useMock) {
    return mockFetchSettleDetail(params);
  }

  return new Promise((resolve) => {
    resolve('real api');
  });
}

/* 提交订单 */
export async function dispatchCommitPay(params) {
  if (config.useMock) {
    return mockDispatchCommitPay(params);
  }

  // 防重复提交
  if (!canSubmitOrder()) {
    throw new Error('订单提交中，请勿重复操作');
  }
  
  // 验证必要参数
  if (!params.addressId) {
    throw new Error('请选择收货地址');
  }
  
  // 清理备注
  const sanitizedParams = {
    ...params,
    remarks: params.remarks ? sanitizeInput(params.remarks) : '',
  };

  try {
    const response = await post('/order/create', sanitizedParams);
    return response.data;
  } catch (error) {
    console.error('提交订单失败:', error);
    throw error;
  }
}

/** 开发票 */
export function dispatchSupplementInvoice() {
  if (config.useMock) {
    const { delay } = require('../_utils/delay');
    return delay();
  }

  return new Promise((resolve) => {
    resolve('real api');
  });
}
