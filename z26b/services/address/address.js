import { cloudbaseTemplateConfig } from '../../config/index';
import { DELIVERY_INFO, createId } from '../cloudbaseMock/index';
import { get, post, put, del } from '../_utils/request';
import { isValidPhone, validateName, validateAddress, sanitizeInput } from '../../utils/security';

export async function getAllAddress() {
  if (cloudbaseTemplateConfig.useMock) {
    return DELIVERY_INFO;
  }
  
  try {
    const response = await get('/address/list');
    const addresses = response.data || [];
    // 兼容字段映射：后端返回 detailAddress，前端组件使用 address
    return addresses.map(addr => ({
      ...addr,
      address: addr.address || addr.detailAddress,
    }));
  } catch (error) {
    console.error('获取地址列表失败:', error);
    return [];
  }
}

/**
 * 创建收货地址
 * @param {{name: String, address: String, phone: String, provinceName?: String, cityName?: String, districtName?: String, isDefault?: Number}} param0
 */
export async function createAddress({ name, address, phone, provinceName, cityName, districtName, isDefault }) {
  // 输入验证
  const nameResult = validateName(name);
  if (!nameResult.valid) {
    throw new Error(nameResult.message);
  }
  
  if (!isValidPhone(phone)) {
    throw new Error('请输入有效的手机号');
  }
  
  const addressResult = validateAddress(address);
  if (!addressResult.valid) {
    throw new Error(addressResult.message);
  }
  
  // 清理输入
  const sanitizedData = {
    name: sanitizeInput(name),
    phone: phone.trim(),
    detailAddress: sanitizeInput(address),
    provinceName: sanitizeInput(provinceName || ''),
    cityName: sanitizeInput(cityName || ''),
    districtName: sanitizeInput(districtName || ''),
    isDefault: isDefault || 0,
  };
  
  if (cloudbaseTemplateConfig.useMock) {
    DELIVERY_INFO.push({
      address: sanitizedData.detailAddress,
      name: sanitizedData.name,
      phone: sanitizedData.phone,
      _id: createId(),
    });
    return;
  }
  
  try {
    const response = await post('/address/create', sanitizedData);
    return response.data;
  } catch (error) {
    console.error('创建地址失败:', error);
    throw error;
  }
}

/**
 * 更新收货地址
 * @param {{name: String, address: String, phone: String, _id: String}} props
 */
export async function updateAddress(props) {
  const { name, address, phone, _id } = props;
  
  // 输入验证
  if (name) {
    const nameResult = validateName(name);
    if (!nameResult.valid) {
      throw new Error(nameResult.message);
    }
  }
  
  if (phone && !isValidPhone(phone)) {
    throw new Error('请输入有效的手机号');
  }
  
  if (address) {
    const addressResult = validateAddress(address);
    if (!addressResult.valid) {
      throw new Error(addressResult.message);
    }
  }
  
  // 清理输入
  const sanitizedData = {
    name: name ? sanitizeInput(name) : undefined,
    detailAddress: address ? sanitizeInput(address) : undefined,
    phone: phone ? phone.trim() : undefined,
  };
  
  if (cloudbaseTemplateConfig.useMock) {
    Object.assign(
      DELIVERY_INFO.find((x) => x._id === _id),
      props,
    );
    return;
  }
  
  try {
    const response = await put(`/address/update/${_id}`, sanitizedData);
    return response.data;
  } catch (error) {
    console.error('更新地址失败:', error);
    throw error;
  }
}

export async function deleteAddress({ id }) {
  if (cloudbaseTemplateConfig.useMock) {
    DELIVERY_INFO.splice(
      DELIVERY_INFO.findIndex((x) => x._id === id),
      1,
    );
    return;
  }
  
  try {
    await del(`/address/${id}`);
  } catch (error) {
    console.error('删除地址失败:', error);
    throw error;
  }
}

export async function getAddress({ id }) {
  if (cloudbaseTemplateConfig.useMock) {
    return DELIVERY_INFO.find((x) => x._id === id);
  }
  
  try {
    const response = await get(`/address/${id}`);
    return response.data;
  } catch (error) {
    console.error('获取地址详情失败:', error);
    throw error;
  }
}
