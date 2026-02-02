# 项目部署与启动说明

本仓库包含三个子项目：

- `z26a/`：前端 Web（Vite + Vue）。
- `z26b/`：微信小程序前端。
- `z26b-backend/`：后端服务（Go）。

## 统一前置要求

- Node.js 18+（用于 `z26a/`、`z26b/`）
- Go 1.24+（用于 `z26b-backend/`）
- 微信开发者工具（用于 `z26b/`）

---

## 1) 构建 `z26a`（Web 前端）

```bash
cd z26a
npm install
npm run build
```

---

## 2) 启动 `z26b`（微信小程序）

使用微信开发者工具：

1. 选择“导入项目”，目录选择 `z26b/`。

> 如需配置环境或接口地址，请查看 `z26b/config/` 下的配置文件。

---

## 3) 启动 `z26b-backend`（Go 后端）

```bash
cd z26b-backend
go mod download
go run main.go
```

需要使用minio和pgsql 均可本地部署
---
