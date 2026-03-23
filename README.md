# 深圳房产信息智能分析平台

基于 **Go + Vue 3 + Element Plus + ECharts** 的深圳房地产信息平台智能分析系统。

该项目从深圳市住房和建设局公开平台抓取预售房源信息，对房源价格、楼层、销售状态、面积、性价比等指标进行量化分析，并通过前端可视化界面展示。

---

## 项目特性

- 支持 **楼盘关键字搜索**
- 支持 **区域筛选**（福田、罗湖、南山、宝安等）
- 支持 **楼栋联动选择**
- 支持 **户型筛选**
- 支持 **房源价格分析**
- 支持 **销售状态统计**
- 支持 **楼层热力图展示**
- 支持 **历史趋势记录**
- 支持 **多楼盘对比分析**
- 支持 **后端缓存与定时刷新**
- 支持 **房源数据 CSV 导出**
- 支持 **收藏楼盘管理**
- 支持 **企业微信价格/销售状态推送**
- 支持 **Docker 容器化部署**

---

## 项目结构

```text
.
├── backend/                        # Go 后端服务
│   ├── main.go                     # 启动 HTTP 服务、定时刷新任务
│   ├── api.go                      # REST API
│   ├── cache.go                    # 缓存与历史记录持久化
│   ├── config.go                   # 全局配置
│   ├── types.go                    # 数据结构定义
│   ├── http.go                     # 请求工具
│   ├── cookie.go                   # Cookie 自动获取与刷新
│   ├── resolver.go                 # 楼盘/楼栋参数解析
│   ├── analyzer.go                 # 数据分析逻辑
│   ├── favorite.go                 # 收藏楼盘管理
│   ├── export.go                   # CSV 导出功能
│   ├── wechat_push.go              # 企业微信推送
│   ├── push_manager.go             # 推送监控管理
│   ├── Dockerfile                  # 后端容器化配置
│   ├── go.mod
│   └── .gitignore
│
├── frontend/                       # Vue 3 前端
│   ├── src/
│   │   ├── api/                    # 前端接口封装
│   │   ├── components/             # 搜索、图表、表格、历史、对比组件
│   │   ├── types/                  # TypeScript 类型定义
│   │   ├── views/
│   │   │   └── Dashboard.vue       # 主页面
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css
│   ├── Dockerfile                  # 前端容器化配置
│   ├── package.json
│   ├── vite.config.ts
│   └── .gitignore
│
├── docker-compose.yml              # Docker 编排配置
├── README.md                       # 根目录项目说明
└── lechenginfo_fullstack_memory.md # 根目录项目记忆文件
```

---

## 技术栈

### 后端
- Go 1.23.6
- net/http
- encoding/json
- 原生文件存储（history.json）

### 前端
- Vue 3
- TypeScript
- Vite
- Element Plus
- ECharts
- Axios

---

## 核心能力

### 1. 智能搜索
前端输入楼盘关键字后：
- 调用后端搜索楼盘接口
- 自动展示候选楼盘
- 选择楼盘后再拉取楼栋列表
- 选择楼栋与户型后执行分析

### 2. 数据分析
后端对抓取到的房源数据进行分析，包括：
- 最低/最高/中位挂牌单价
- 最低/最高/平均总价
- 建筑面积 / 套内面积 / 拓展面积
- 得房率
- 低中高楼层分布
- 销售状态统计
- 性价比衍生指标

### 3. 历史趋势
每次分析完成后，会将结果写入历史记录文件：
- 默认持久化到 `backend/history.json`
- 前端可展示价格趋势图
- 支持查看与删除历史记录

### 4. 多楼盘对比
支持最多 5 个楼盘同时对比，展示：
- 综合雷达图
- 价格柱状图
- 销售状态对比图
- 指标数据表

### 5. 缓存与定时刷新
- 后端分析结果会缓存 30 分钟
- 默认定时刷新默认楼盘数据
- 可手动清除缓存

### 6. 数据导出与收藏
- 支持房源数据 CSV 导出
- 支持收藏楼盘管理
- 支持推送配置开关

### 7. 企业微信推送
- 支持价格变化提醒（≥3%变化）
- 支持销售状态提醒（≥5%变化）
- 支持新收藏楼盘通知
- 每30分钟自动检查变化

---

## 后端 API

### 楼盘搜索
```http
GET /api/search?keyword=乐宸&pageIndex=1&pageSize=12&zone=
```

### 获取楼栋列表
```http
GET /api/buildings?ysProjectId=xxx&preSellId=xxx
```

### 执行分析
```http
POST /api/analyze
Content-Type: application/json
```

请求体示例：
```json
{
  "keyword": "乐宸",
  "buildingName": "1栋",
  "houseType": "三房"
}
```

### 多楼盘对比
```http
POST /api/compare
Content-Type: application/json
```

### 获取历史记录
```http
GET /api/history
```

### 删除历史记录
```http
DELETE /api/history?id=记录ID
```

### 清除缓存
```http
POST /api/cache/clear
```

### 收藏楼盘管理
```http
GET    /api/favorites                    # 获取收藏列表
POST   /api/favorites                    # 新增收藏
DELETE /api/favorites?id=xxx             # 删除收藏
PUT    /api/favorites/update?id=xxx      # 更新推送配置
```

### 导出 CSV
```http
GET /api/export/csv?keyword=xxx&buildingName=xxx&houseType=xxx&zone=
```

---

## 启动方式

### 方式一：传统启动（开发环境）

#### 1. 启动后端

```powershell
cd backend
go run .
```

后端默认监听：

```text
http://localhost:8080
```

---

## 2. 启动前端

```powershell
cd frontend
npm install
npm run dev
```

前端默认地址通常为：

```text
http://127.0.0.1:5173
```

如果端口被占用，Vite 会自动切换到其它端口，比如 `5174`。

### 方式二：Docker 容器化启动

#### 1. 构建并启动所有服务

```bash
docker-compose up --build
```

#### 2. 访问服务

- 前端访问：`http://localhost:8081`
- 后端 API：`http://localhost:8080`

#### 3. 停止服务

```bash
docker-compose down
```

---

## 默认演示配置

后端默认配置：
- 楼盘关键字：`乐宸`
- 楼栋：`1栋`
- 户型：`三房`

对应配置文件：

```text
backend/config.go
```

---

## 重要说明

### Cookie 机制
由于目标站点存在会话限制，后端会：
- 启动时自动访问入口页获取 Cookie
- 每 30 分钟自动刷新 Cookie
- 如果自动获取失败，则降级使用静态 Cookie

### 数据文件
运行时数据文件默认保存在：

```text
backend/history.json          # 历史分析记录
backend/favorites.json        # 收藏楼盘数据
backend/push_records.json     # 推送记录
```

这些文件建议加入 `.gitignore`，避免把运行时数据提交到 GitHub。

### Windows 环境说明
你当前是 Windows 环境，建议使用 PowerShell 启动：

```powershell
cd backend
go run .
```

```powershell
cd frontend
npm run dev
```

### Docker 环境说明
Docker 部署支持跨平台运行，Windows 环境下需要安装 Docker Desktop。

---

## 企业微信推送配置

### 1. 配置推送
1. 在收藏夹中启用推送开关
2. 配置价格变化提醒（≥3%）
3. 配置销售状态提醒（≥5%）

### 2. 推送内容
- **价格变化提醒**：中位单价变化超过3%
- **销售状态提醒**：在售率变化超过5%
- **新收藏通知**：新增收藏楼盘信息

### 3. 推送频率
- 每30分钟自动检查一次
- 只推送显著变化，避免频繁打扰

### 4. Webhook 配置
推送使用企业微信机器人 Webhook：
```
https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=f5e8ee74-efbd-4e02-9011-20ab8decb475
```

---

## 适用场景

- 深圳新房预售房源观察
- 楼盘横向对比
- 购房数据辅助分析
- 销售进度趋势跟踪
- 房产信息可视化展示
- 价格变化实时监控
- 投资决策数据支持

---

## 🎯 项目进展

✅ **已完成功能**：
- 基础楼盘搜索与分析
- 区域筛选支持  
- CSV 数据导出
- 收藏楼盘管理
- 企业微信推送
- Docker 容器化部署
- 价格/销售变化预警

🚀 **下一步可扩展方向**：
- 登录页与用户体系
- 导出 Excel 格式报表
- 移动端适配
- 云服务器部署
- Nginx 反向代理配置
- 更多数据源集成
