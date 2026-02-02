# 数据库初始化工具

## 文件结构

```
cmd/initdb/
├── main.go     # 主入口文件，负责连接数据库和调用初始化函数
└── init.go     # 初始化逻辑文件，包含所有数据库操作函数
```

## 代码组织

### main.go
- **职责**: 程序入口，负责配置加载和数据库连接
- **功能**:
  - 加载环境变量
  - 连接数据库
  - 按顺序调用初始化函数
  - 错误处理

### init.go
- **职责**: 数据库初始化逻辑
- **主要函数**:
  - `DropAllTables(db)` - 删除所有表
  - `CreateTables(db)` - 创建所有表结构
  - `InsertSampleData(db)` - 插入示例数据
  - `PrintSummary()` - 打印初始化摘要

- **辅助函数**:
  - `generateUUID()` - 生成UUID
  - `hashPassword(password)` - 加密密码
  - `insertAdmin()` - 插入管理员
  - `insertTestUser()` - 插入测试用户
  - `insertCategories()` - 插入分类
  - `insertTags()` - 插入标签
  - `insertSampleProduct()` - 插入示例商品
  - `insertSwiper()` - 插入轮播图
  - `insertHomeContent()` - 插入首页内容
  - `insertCoupon()` - 插入优惠券
  - `insertPromotion()` - 插入促销活动

## 使用方法

### 方式1: 使用 Makefile
```bash
make init-db
```

### 方式2: 直接运行
```bash
go run cmd/initdb/main.go
```

### 方式3: 编译后运行
```bash
go build -o initdb.exe cmd/initdb/*.go
./initdb.exe
```

## 优势

1. **代码分离**: 主逻辑和数据操作分离，更易维护
2. **可复用**: init.go 中的函数可以在其他地方复用
3. **易测试**: 每个函数独立，便于单元测试
4. **清晰结构**: main.go 简洁明了，只关注流程控制
5. **易扩展**: 需要添加新功能只需在 init.go 中添加新函数

## 初始化流程

1. **连接数据库** (main.go)
2. **删除旧表** → `DropAllTables()`
3. **创建新表** → `CreateTables()`
4. **插入数据** → `InsertSampleData()`
   - 管理员账户
   - 测试用户
   - 分类数据
   - 标签数据
   - 示例商品
   - 轮播图
   - 首页内容
   - 优惠券
   - 促销活动
5. **显示摘要** → `PrintSummary()`

## 自定义数据

要自定义初始数据，编辑 [init.go](init.go) 中对应的函数：

- 修改管理员: `insertAdmin()`
- 修改分类: `insertCategories()`
- 修改标签: `insertTags()`
- 添加更多商品: `insertSampleProduct()`

## 安全提示

⚠️ **此工具会删除所有现有表和数据！**
- 仅用于首次部署或开发环境
- 生产环境使用前请备份数据

## 错误处理

- 数据库连接失败: 程序立即退出
- 表创建失败: 程序立即退出
- 数据插入失败: 显示警告但继续执行
