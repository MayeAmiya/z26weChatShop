# 前端API集成说明

## 概述
已将前端从mock数据模式切换到连接真实后端API模式，并添加了完整的安全机制。

## 配置更改

### 1. 关闭Mock模式
文件: `config/index.js`
```javascript
export const config = {
  useMock: false,  // 改为 false 连接后端
};
```

### 2. 环境配置
文件: `config/env.js`
```javascript
// 开发模式：使用本地后端，允许测试 OpenID
// 生产模式：必须真实微信登录
export const IS_DEV_MODE = true;

// 生产环境需要配置正确的 HTTPS API 地址
export const API_CONFIG = {
  development: 'http://localhost:8080/api',
  production: 'https://your-production-api.com/api',
};
```

### 3. HTTP请求工具
文件: `services/_utils/request.js`
- 统一的HTTP请求封装
- 自动添加认证头 (X-OpenID, Authorization)
- 支持 Token 认证
- 生产环境强制认证

## 安全机制

### 1. 用户认证
- 小程序端通过 `X-OpenID` header 发送用户标识
- 后端通过 `MiniProgramAuthMiddleware` 中间件验证
- 生产环境必须提供有效的 OpenID
- 支持 JWT Token 认证（可选）

### 2. 输入验证
文件: `utils/security.js`
- `validateName()` - 姓名验证
- `validatePhone()` - 手机号验证  
- `validateAddress()` - 地址验证
- `validateComment()` - 评论内容验证
- `validateQuantity()` - 数量验证
- `sanitizeInput()` - XSS 防护
- `debounceSubmit()` - 防重复提交

### 3. 防重复提交
关键操作（添加购物车、提交订单、评论）添加了防重复提交保护。

## 已更新的服务

### 商品相关
- ✅ `services/good/fetchGoodsList.js` - 商品列表
- ✅ `services/good/fetchGood.js` - 商品详情
- ✅ `services/good/fetchCategoryList.js` - 分类列表
- ✅ `services/good/fetchSearchResult.js` - 搜索功能
- ✅ `services/sku/sku.js` - SKU详情和列表

### 购物车相关
- ✅ `services/cart/cart.js`
  - fetchCartItems - 获取购物车
  - createCartItem - 添加到购物车（带数量验证）
  - updateCartItemCount - 更新数量（带数量验证）
  - deleteCartItem - 删除购物车项

### 订单相关
- ✅ `services/order/orderList.js` - 订单列表
- ✅ `services/order/orderDetail.js` - 订单详情
- ✅ `services/order/orderConfirm.js` - 提交订单（带防重复提交）

### 地址相关
- ✅ `services/address/address.js`
  - getAllAddress - 获取地址列表
  - createAddress - 创建地址（带完整验证）
  - updateAddress - 更新地址（带完整验证）
  - deleteAddress - 删除地址

### 认证相关
- ✅ `services/auth/auth.js`
  - login - 微信登录
  - silentLogin - 静默登录
  - isLoggedIn - 检查登录状态
  - checkLoginStatus - 检查并跳转登录
  - logout - 登出

### 评论相关
- ✅ `services/comments/comments.js`
  - getGoodsDetailCommentInfo - 获取商品评论
  - createComment - 提交评论（带内容验证）

### 其他
- ✅ `services/home/home.js` - 首页轮播
- ✅ `services/promotion/detail.js` - 促销详情

## API端点映射

| 前端服务 | 后端API | 方法 | 认证要求 |
|---------|---------|------|----------|
| login | /api/wechat/login | POST | 否 |
| fetchGoodsList | /api/goods/list | GET | 否 |
| fetchGood | /api/goods/:id | GET | 否 |
| getCategoryList | /api/goods/category/list | GET | 否 |
| getSearchResult | /api/goods/search | GET | 否 |
| getSkuDetail | /api/sku/:id | GET | 否 |
| getAllSku | /api/sku/list/:spuId | GET | 否 |
| fetchCartItems | /api/cart/items | GET | 是 |
| createCartItem | /api/cart/add | POST | 是 |
| updateCartItemCount | /api/cart/update/:id | PUT | 是 |
| deleteCartItem | /api/cart/remove/:id | DELETE | 是 |
| fetchOrders | /api/order/list | GET | 是 |
| fetchOrderDetail | /api/order/:id | GET | 是 |
| dispatchCommitPay | /api/order/create | POST | 是 |
| getAllAddress | /api/address/list | GET | 是 |
| createAddress | /api/address/create | POST | 是 |
| updateAddress | /api/address/update/:id | PUT | 是 |
| deleteAddress | /api/address/:id | DELETE | 是 |
| getHomeSwiper | /api/home/swiper | GET | 否 |
| getGoodsDetailCommentInfo | /api/comment/list/:spuId | GET | 否 |
| createComment | /api/comment/submit | POST | 是 |
| fetchPromotion | /api/home/promotions | GET | 否 |

## 使用说明

### 1. 启动后端服务
```bash
cd z26b-backend
go run main.go
```
后端运行在: `http://localhost:8080`

### 2. 配置小程序
确保小程序开发工具中已关闭"不校验合法域名"选项（开发阶段）

### 3. 数据格式转换
已处理前端期望的数据格式和后端返回格式的差异：
- SPU字段映射：`primary_image` → `primaryImage`
- SKU字段映射：`min_sale_price` → `minSalePrice`
- 购物车字段：`sku_id` 处理

### 4. 错误处理
所有API调用都包含try-catch错误处理，失败时会在控制台输出错误信息

## 生产环境部署清单

### 小程序端配置
1. 修改 `config/env.js`:
   - 设置 `IS_DEV_MODE = false`
   - 配置正确的生产环境 API 地址（必须 HTTPS）

2. 在小程序管理后台配置:
   - 服务器域名白名单
   - 业务域名

### 后端配置
1. 设置环境变量:
   - `GIN_MODE=release`
   - `WECHAT_APP_ID` - 微信小程序 AppID
   - `WECHAT_APP_SECRET` - 微信小程序 AppSecret
   - `ALLOWED_ORIGINS` - 允许的跨域来源

2. 启用 HTTPS

## 测试建议

### 第一步：测试基础API
```bash
# 健康检查
curl http://localhost:8080/health

# 获取商品列表
curl http://localhost:8080/api/goods/list

# 获取分类列表
curl http://localhost:8080/api/goods/category/list
```

### 第二步：测试认证API
```bash
# 带认证头获取购物车
curl -H "X-OpenID: oTest_dev_openid_001" http://localhost:8080/api/cart/items
```

### 第三步：小程序测试
在小程序中测试各个页面功能

## 待完成功能

- [ ] 微信支付功能对接
- [ ] 图片上传
- [ ] 用户资料编辑

## 回滚到Mock模式

如需切换回mock模式：
```javascript
// config/index.js
export const config = {
  useMock: true,
};
```
