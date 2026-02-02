import { cloudbaseTemplateConfig } from "../../config/index"
import { ATTR_VALUE } from "../cloudbaseMock/index"
import { get } from '../_utils/request';

export async function getAllAttrValues(skuId) {
  if (cloudbaseTemplateConfig.useMock) {
    return ATTR_VALUE.filter(x => x.sku.find(x => x._id === skuId))
  }
  try {
    // 从SKU详情中获取属性值
    const response = await get(`/sku/${skuId}`);
    const sku = response.data;
    
    // 如果有 description，生成一个虚拟的属性值
    // 这样可以在规格选择弹窗中展示
    if (sku?.description) {
      return [{
        _id: `attr_${skuId}`,
        value: sku.description,
        attr_name: {
          _id: 'attr_spec',
          name: '规格'
        }
      }];
    }
    
    return sku?.attr_values || [];
  } catch (error) {
    console.error('获取属性值失败:', error);
    return [];
  }
}
