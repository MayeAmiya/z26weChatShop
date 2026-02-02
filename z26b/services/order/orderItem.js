import { getSkuDetail } from '../sku/sku';
import { cloudbaseTemplateConfig } from '../../config/index';
import { ORDER_ITEM, createId } from '../cloudbaseMock/index';
import { get } from '../_utils/request';

export async function getOrderItem(id) {
  if (cloudbaseTemplateConfig.useMock) {
    return ORDER_ITEM.find(x => x._id === id);
  }
  
  try {
    const response = await get(`/order/${id}`);
    const order = response.data;
    if (order && order.items && order.items.length > 0) {
      return order.items[0];
    }
    return null;
  } catch (error) {
    console.error('获取订单项失败:', error);
    return null;
  }
}

export async function createOrderItem({ count, skuId, orderId }) {
  if (cloudbaseTemplateConfig.useMock) {
    ORDER_ITEM.push({
      _id: createId(),
      count,
      order: {
        _id: orderId,
      },
      sku: {
        _id: skuId,
      },
    });
    return;
  }
  // 订单项通过创建订单时自动创建，这里不需要单独调用
  return;
}

/**
 *
 * @param {{orderId: String}} param0
 */
export async function getAllOrderItemsOfAnOrder({ orderId }) {
  if (cloudbaseTemplateConfig.useMock) {
    const orderItems = ORDER_ITEM.filter((orderItem) => orderItem.order._id === orderId);
    await Promise.all(
      orderItems.map(async (orderItem) => {
        const skuId = orderItem.sku._id;
        const sku = await getSkuDetail(skuId);
        orderItem.sku = sku;
      }),
    );
    return orderItems;
  }

  try {
    const response = await get(`/order/${orderId}`);
    const order = response.data;
    // 后端返回的 items 格式转换为前端期望的格式
    const items = order?.items || [];
    return items.map(item => {
      // 确保 sku 对象有完整的数据
      const sku = item.sku || {};
      const spu = sku.spu || {};
      
      // 将 description 转换为 attr_value 格式供组件显示规格
      const attrValue = sku.description 
        ? [{ value: sku.description }]  // 组件会用 map(v => v.value).join('，') 显示
        : [];
      
      // 构建适配 pages/order/components/goods-card 组件的数据格式
      // 该组件期望: goods.sku.image, goods.sku.spu.name, goods.sku.price, goods.sku.attr_value, goods.count
      return {
        _id: item._id,
        id: item._id,
        count: item.quantity,
        quantity: item.quantity,
        order: { _id: orderId },
        // sku 数据结构（goods-card 需要的格式）
        sku: {
          _id: sku._id || item.skuId,
          spuId: sku.spuId || '',
          price: item.price || sku.price || 0,
          image: sku.image || spu.cover_image || '',
          description: sku.description || '',
          attr_value: attrValue, // 规格值数组，用于显示商品规格
          spu: {
            _id: spu._id || '',
            name: spu.name || '商品',
            cover_image: spu.cover_image || '',
          },
        },
      };
    });
  } catch (error) {
    console.error('获取订单项列表失败:', error);
    return [];
  }
}
