# AgentGame - AI驱动的古风RPG游戏

一个基于AI Agent的古风RPG游戏，具有智能NPC对话系统。

## 项目架构

```
agentGame/
├── client/          # 游戏客户端 (Web前端)
├── server/          # 游戏服务端 (Go)
├── gm/              # GM管理端 (Web前端)
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

### GM管理端 (gm/)
- **框架**: Vue 3
- **UI组件**: Element Plus
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

### 4. 启动GM管理端

```bash
cd gm

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

GM管理端默认运行在 `http://localhost:5174`

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
- `components/` - 通用组件
- `api/` - 后端API调用

## 开发指南

### 添加新NPC

1. 在 `server/internal/agent/` 定义NPC Agent
2. 配置NPC人设和对话逻辑
3. 在客户端添加NPC渲染

### 添加新场景

1. 在 `client/src/assets/` 添加场景资源
2. 在 `client/src/game/scenes/` 创建场景类
3. 注册场景到游戏配置

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
