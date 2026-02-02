import { cloudbaseTemplateConfig } from '../../config/index';
import { COMMENTS, createId } from '../cloudbaseMock/index';
import { get, post } from '../_utils/request';
import { validateComment, sanitizeInput, debounceSubmit } from '../../utils/security';

// 防重复提交
const canSubmitComment = debounceSubmit('submitComment', 2000);

export async function getGoodsDetailCommentInfo(spuId) {
  if (cloudbaseTemplateConfig.useMock) {
    const all = COMMENTS.filter((x) => x.spu._id === spuId);
    const good = all.filter((x) => x.rating > 3);
    const firstAndTotal = (x) => ({
      data: {
        records: x.length === 0 ? [] : [x[0]],
        total: x.length,
      },
    });
    return Promise.resolve([firstAndTotal(all), firstAndTotal(good)]);
  }

  try {
    const response = await get(`/comment/list/${spuId}`);
    // Backend returns: { data: { pageNum, pageSize, totalCount, pageList } }
    const commentList = response.data && response.data.pageList ? response.data.pageList : [];
    const total = response.data && response.data.totalCount ? response.data.totalCount : 0;
    
    const goodComments = commentList.filter(x => x.commentScore > 3);
    const goodTotal = goodComments.length;
    
    const firstAndTotal = (list, count) => ({
      data: {
        records: list.length === 0 ? [] : [list[0]],
        total: count,
      },
    });
    
    return [firstAndTotal(commentList, total), firstAndTotal(goodComments, goodTotal)];
  } catch (error) {
    console.error('获取评论失败:', error);
    return [{ data: { records: [], total: 0 } }, { data: { records: [], total: 0 } }];
  }
}

export async function getCommentsOfSpu({ spuId, pageNumber, pageSize }) {
  if (cloudbaseTemplateConfig.useMock) {
    const all = COMMENTS.filter((x) => x.spu._id === spuId);
    const startIndex = (pageNumber - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    const records = all.slice(startIndex, endIndex);
    return {
      records,
      total: all.length,
    };
  }
  
  try {
    const response = await get(`/comment/list/${spuId}?pageNumber=${pageNumber}&pageSize=${pageSize}`);
    return response.data;
  } catch (error) {
    console.error('获取评论列表失败:', error);
    return { records: [], total: 0 };
  }
}

/**
 * 提交商品评论
 * @param {{orderItemId: string, content: string, rating: number, spuId: string}} param0
 */
export async function createComment({ orderItemId, content, rating, spuId }) {
  // 验证评论内容
  const contentResult = validateComment(content);
  if (!contentResult.valid) {
    throw new Error(contentResult.message);
  }
  
  // 验证评分
  if (typeof rating !== 'number' || rating < 1 || rating > 5) {
    throw new Error('请选择有效的评分（1-5星）');
  }
  
  // 防重复提交
  if (!canSubmitComment()) {
    throw new Error('评论提交过于频繁，请稍后重试');
  }
  
  // 清理输入
  const sanitizedContent = sanitizeInput(content);
  
  if (cloudbaseTemplateConfig.useMock) {
    COMMENTS.push({
      _id: createId(),
      content: sanitizedContent,
      rating,
      order_item: { _id: orderItemId },
      spu: { _id: spuId },
    });
    return;
  }
  
  try {
    const response = await post('/comment/submit', {
      orderItemId,
      commentContent: sanitizedContent,
      commentScore: rating,
      spuId,
    });
    return response.data;
  } catch (error) {
    console.error('提交评论失败:', error);
    throw error;
  }
}
