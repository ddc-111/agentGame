# AgentGame - AI驱动的古风RPG游戏

一个基于AI Agent的古风RPG游戏，具有智能NPC对话系统和完整的游戏编辑器。

## 项目架构

```
agentGame/
├── client/          # 游戏客户端 (Web前端)
├── server/          # 游戏服务端 (Go)
├── gm/              # GM管理端/游戏编辑器 (Web前端)
├── docs/            # 项目文档
└── deploy/          # 部署配置
```

## 技术栈

### 客户端 (client/)
- **游戏引擎**: Phaser 3
- **构建工具**: Vite
- **网络通信**: Socket.IO Client

### 服务端 (server/)
- **语言**: Go 1.21+
- **Web框架**: Gin
- **WebSocket**: gorilla/websocket
- **AI集成**: OpenAI API

### GM管理端/编辑器 (gm/)
- **框架**: Vue 3
- **UI组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP客户端**: Axios

## 环境要求

- Node.js >= 18.0.0
- Go >= 1.21
- npm >= 9.0.0

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/ddc-111/agentGame.git
cd agentGame
```

### 2. 启动服务端

```bash
cd server

# 下载依赖
go mod tidy

# 运行服务
go run cmd/gameserver/main.go
```

服务端默认运行在 `http://localhost:8080`

### 3. 启动客户端

```bash
cd client

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

客户端默认运行在 `http://localhost:5173`

### 4. 启动GM管理端/编辑器

```bash
cd gm

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

GM管理端默认运行在 `http://localhost:5174`

## 编辑器功能模块

### 场景编辑
- 场景列表管理
- 场景编辑器（背景、尺寸、NPC、传送点）
- 图块管理

### NPC编辑
- NPC列表管理
- NPC编辑器（头像、精灵图、位置、行为）
- 对话树编辑

### 智能体配置
- 智能体列表管理
- 智能体编辑器（模型、提示词、知识库、工具）
- 记忆配置（短期记忆、长期摘要、向量检索）

### 大模型配置
- 模型提供商管理（OpenAI、Anthropic、本地模型）
- 模型配置（场景化模型选择）
- 连接测试

### 提示词配置
- 提示词模板管理
- 变量管理
- 提示词测试

### 商店配置
- 商店列表管理
- 道具管理
- 商店编辑（商品、价格、库存、折扣）

### 任务系统
- 任务列表管理
- 任务编辑（触发条件、目标、奖励）
- 流程编排（可视化流程编辑器）

### 系统配置
- 游戏配置（基础、玩家、世界、战斗、经济）
- 数据导出
- 数据导入

## NPC购物流程示例

编辑器内置了NPC出门购物的完整流程：

1. **NPC在家** - 李掌柜在家中准备出门
2. **前往商店** - 离开家，前往杂货铺
3. **检查商店** - 确认商店是否开门
4. **购买物品** - 购买草药、馒头等物资
5. **返回家中** - 带着购买的物品回家
6. **流程结束** - NPC完成购物，更新库存

## 配置说明

### 服务端配置

创建 `server/config.yaml`:

```yaml
server:
  port: 8080
  mode: debug

ai:
  provider: openai
  api_key: your-api-key-here
  model: gpt-4

game:
  max_players: 100
  tick_rate: 20
```

## 项目模块说明

### 客户端模块 (client/src/)
- `game/` - 游戏核心逻辑、场景管理
- `components/` - UI组件
- `network/` - WebSocket通信
- `assets/` - 图片、音频等资源
- `utils/` - 工具函数

### 服务端模块 (server/internal/)
- `agent/` - AI Agent系统，NPC智能对话
- `game/` - 游戏逻辑、状态管理
- `network/` - 网络层、消息处理
- `config/` - 配置管理

### GM管理端模块 (gm/src/)
- `views/` - 页面视图
  - `scene/` - 场景编辑相关页面
  - `npc/` - NPC编辑相关页面
  - `agent/` - 智能体配置相关页面
  - `llm/` - 大模型配置相关页面
  - `prompt/` - 提示词配置相关页面
  - `shop/` - 商店配置相关页面
  - `task/` - 任务系统相关页面
  - `config/` - 系统配置相关页面
- `stores/` - Pinia状态管理
- `api/` - 后端API调用
- `components/` - 通用组件

## 开发指南

### 添加新NPC

1. 在编辑器的NPC管理页面创建新NPC
2. 配置NPC的基本信息、头像、精灵图
3. 关联智能体和对话树
4. 设置NPC的行为模式和日程

### 添加新场景

1. 在编辑器的场景管理页面创建新场景
2. 上传背景图片
3. 配置场景中的NPC和传送点
4. 使用图块编辑器添加装饰

### 配置智能体

1. 在智能体管理页面创建新智能体
2. 选择大模型提供商和模型
3. 编写系统提示词（定义NPC人设）
4. 配置记忆参数（短期记忆、摘要）
5. 添加知识库和工具

## 部署

### Docker部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d
```

### 手动部署

```bash
# 构建客户端
cd client && npm run build

# 构建GM管理端
cd gm && npm run build

# 编译服务端
cd server && go build -o gameserver cmd/gameserver/main.go
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- GitHub: [@ddc-111](https://github.com/ddc-111)
- 项目链接: [https://github.com/ddc-111/agentGame](https://github.com/ddc-111/agentGame)
