/**
 * 图片URL处理工具
 * 直接返回原始图片URL
 */
export async function getCloudImageTempUrl(images) {
  if (!images || !Array.isArray(images)) {
    return images || [];
  }
  return images;
}

export async function getSingleCloudImageTempUrl(image) {
  return image || '';
}
