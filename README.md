# AgentGame - AI驱动的古风RPG游戏框架

一个基于AI Agent的古风RPG游戏框架，具有智能NPC对话系统、完整的游戏编辑器和自循环开发系统。

## 🎯 项目特色

- **AI驱动NPC**: 基于大语言模型的智能NPC对话系统
- **可视化编辑器**: 完整的GM管理端，支持场景、NPC、任务等配置
- **自循环开发**: 内置AI驱动的自循环开发系统，自动分析、测试、重构
- **模块化架构**: 清晰的三层架构，易于扩展和维护

---

## 📐 完整项目架构

```
agentGame/
├── 📁 agent/                    # 自循环Agent系统 (Python)
│   ├── analyzers/               # 代码分析器
│   │   ├── code_analyzer.py     # 代码结构分析
│   │   ├── test_analyzer.py     # 测试覆盖分析
│   │   └── gap_analyzer.py      # 差距识别
│   ├── executors/               # 任务执行器
│   │   ├── test_executor.py     # 测试执行
│   │   ├── build_executor.py    # 构建执行
│   │   └── task_executor.py     # 任务执行
│   ├── generators/              # 生成器
│   │   ├── task_generator.py    # 任务生成
│   │   └── report_generator.py  # 报告生成
│   ├── agents/                  # 子Agent
│   │   └── requirement_agent.py # 需求生成Agent
│   ├── utils/                   # 工具类
│   │   ├── llm_client.py        # LLM客户端
│   │   └── history_tracker.py   # 历史追踪
│   ├── main.py                  # 主入口
│   ├── orchestrator.py          # 调度器核心
│   └── config_manager.py        # 配置管理
│
├── 📁 server/                   # 游戏服务端 (Go)
│   ├── cmd/gameserver/          # 入口
│   │   └── main.go
│   └── internal/
│       ├── agent/               # AI对话系统
│       │   ├── chat.go          # 对话管理
│       │   └── memory.go        # 记忆系统
│       ├── config/              # 配置管理
│       ├── database/            # 数据层
│       │   ├── models/          # 数据模型
│       │   └── repository/      # 数据访问
│       ├── game/                # 游戏逻辑
│       │   ├── combat.go        # 战斗系统
│       │   ├── inventory.go     # 背包系统
│       │   └── savegame.go      # 存档系统
│       ├── generator/           # AI配置生成
│       ├── mcp/                 # MCP服务器
│       └── network/             # 网络层
│           ├── server.go        # HTTP服务器
│           ├── websocket.go     # WebSocket
│           └── *_handlers.go    # API处理器
│
├── 📁 client/                   # 游戏客户端 (Phaser 3)
│   └── src/
│       ├── main.js              # 入口
│       └── game/
│           ├── scenes/          # 游戏场景
│           │   ├── BootScene.js # 启动场景
│           │   └── GameScene.js # 主游戏场景
│           ├── systems/         # 游戏系统
│           │   ├── CombatManager.js    # 战斗管理
│           │   └── InventoryManager.js # 背包管理
│           └── ui/              # UI组件
│               ├── CombatUI.js     # 战斗UI
│               ├── InventoryUI.js  # 背包UI
│               └── MiniMap.js      # 小地图
│
├── 📁 gm/                       # GM管理端 (Vue 3)
│   └── src/
│       ├── main.js              # 入口
│       ├── App.vue              # 根组件
│       ├── router/              # 路由
│       ├── stores/              # 状态管理
│       │   ├── scene.js         # 场景Store
│       │   ├── npc.js           # NPC Store
│       │   ├── agent.js         # 智能体Store
│       │   ├── shop.js          # 商店Store
│       │   ├── task.js          # 任务Store
│       │   └── config.js        # 配置Store
│       ├── views/               # 页面
│       │   ├── scene/           # 场景编辑
│       │   ├── npc/             # NPC编辑
│       │   ├── agent/           # 智能体配置
│       │   ├── shop/            # 商店配置
│       │   ├── task/            # 任务系统
│       │   └── config/          # 系统配置
│       ├── components/          # 通用组件
│       └── api/                 # API调用
│
├── 📁 docs/                     # 项目文档
├── 📁 deploy/                   # 部署配置
│
├── agent_config.yaml            # 自循环Agent配置
├── AGENT.md                     # Agent系统文档
├── SELF_LOOP_REPORT.md          # 自循环报告
└── IMPROVEMENT_BACKLOG.json     # 改进待办
```

---

## 🤖 自循环Agent系统

### 系统概述

自循环Agent是一个AI驱动的开发系统，能够自动：
- **分析代码**: 识别代码结构、复杂度、依赖关系
- **运行测试**: 自动执行三端测试套件
- **识别问题**: 发现代码质量问题和改进机会
- **生成任务**: 自动生成改进任务
- **执行重构**: 使用LLM辅助代码重构
- **生成报告**: 自动生成开发报告

### 运行流程

```
┌─────────────────────────────────────────────────────────────┐
│                    自循环Agent运行流程                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ 代码分析 │───▶│ 测试分析 │───▶│ 构建验证 │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│       │               │               │                    │
│       ▼               ▼               ▼                    │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ 差距分析 │◀───│ 测试执行 │◀───│ 任务生成 │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│       │                                                     │
│       ▼                                                     │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ 需求生成 │───▶│ 任务执行 │───▶│ 报告生成 │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 使用方法

```bash
# 运行单次迭代
python agent/main.py -n 1

# 运行5次迭代，启用LLM自动执行
python agent/main.py -n 5 --enable-llm

# 只生成报告，不执行任务
python agent/main.py --report-only

# 显示历史记录
python agent/main.py --show-history
```

### 配置说明

`agent_config.yaml`:
```yaml
llm:
  enabled: true
  api_url: https://api.openai.com/v1
  api_key: your-api-key
  model: gpt-4
  temperature: 0.3

iterations:
  max_iterations: 10
  stop_on_stable: true

tasks:
  auto_execute: true
  priority_filter:
    - critical
    - high
    - medium
```

---

## 🖥️ 技术栈

### 服务端 (server/)
| 技术 | 用途 |
|------|------|
| Go 1.21+ | 主语言 |
| Gin | Web框架 |
| GORM | ORM框架 |
| gorilla/websocket | WebSocket |
| OpenAI API | AI集成 |

### 客户端 (client/)
| 技术 | 用途 |
|------|------|
| Phaser 3 | 游戏引擎 |
| Vite | 构建工具 |
| Socket.IO | 网络通信 |

### GM管理端 (gm/)
| 技术 | 用途 |
|------|------|
| Vue 3 | 前端框架 |
| Element Plus | UI组件库 |
| Pinia | 状态管理 |
| Vue Router | 路由管理 |
| Vitest | 测试框架 |

### 自循环Agent (agent/)
| 技术 | 用途 |
|------|------|
| Python 3.10+ | 主语言 |
| requests | HTTP客户端 |
| PyYAML | 配置解析 |

---

## 📊 测试覆盖

### 当前测试状态

| 模块 | 测试文件 | 测试用例 | 通过率 |
|------|----------|----------|--------|
| Server | 8个 | 324个 | 100% |
| Client | 4个 | 124个 | 100% |
| GM | 11个 | 125个 | 96% |
| **总计** | **23个** | **573个** | **99.1%** |

### 运行测试

```bash
# Server端测试
cd server && go test ./internal/... -v

# Client端测试
cd client && npm test

# GM端测试
cd gm && npm test
```

---

## 🚀 快速开始

### 环境要求

- Node.js >= 18.0.0
- Go >= 1.21
- Python >= 3.10 (可选，用于自循环Agent)
- npm >= 9.0.0

### 1. 克隆项目

```bash
git clone https://github.com/ddc-111/agentGame.git
cd agentGame
```

### 2. 启动服务端

```bash
cd server
go mod tidy
go run cmd/gameserver/main.go
```

服务端默认运行在 `http://localhost:8080`

### 3. 启动客户端

```bash
cd client
npm install
npm run dev
```

客户端默认运行在 `http://localhost:5173`

### 4. 启动GM管理端

```bash
cd gm
npm install
npm run dev
```

GM管理端默认运行在 `http://localhost:5174`

---

## 🎮 功能模块

### 游戏系统

- **场景系统**: 多场景切换、传送点、NPC分布
- **战斗系统**: 回合制战斗、技能系统、装备系统
- **背包系统**: 道具管理、装备穿戴、物品使用
- **任务系统**: 主线任务、支线任务、任务追踪
- **存档系统**: 游戏进度保存、多存档支持

### AI系统

- **NPC对话**: 基于LLM的智能对话
- **记忆系统**: 对话历史、玩家信息记忆
- **行为系统**: NPC日程、状态机
- **配置生成**: AI自动生成游戏配置

### 编辑器功能

- **场景编辑**: 背景、尺寸、NPC、传送点
- **NPC编辑**: 头像、精灵图、行为、对话
- **智能体配置**: 模型、提示词、知识库
- **商店配置**: 商品、价格、库存
- **任务编辑**: 触发条件、目标、奖励

---

## 📁 核心模块说明

### Server端模块 (server/internal/)

| 模块 | 功能 | 关键文件 |
|------|------|----------|
| `agent/` | AI对话系统 | chat.go, memory.go |
| `database/` | 数据层 | models/, repository/ |
| `game/` | 游戏逻辑 | combat.go, inventory.go |
| `network/` | 网络层 | server.go, websocket.go |
| `generator/` | AI生成 | generator.go |
| `mcp/` | MCP协议 | server.go |

### Client端模块 (client/src/game/)

| 模块 | 功能 | 关键文件 |
|------|------|----------|
| `scenes/` | 游戏场景 | GameScene.js |
| `systems/` | 游戏系统 | CombatManager.js |
| `ui/` | UI组件 | CombatUI.js, InventoryUI.js |

### GM端模块 (gm/src/)

| 模块 | 功能 | 关键文件 |
|------|------|----------|
| `stores/` | 状态管理 | scene.js, npc.js, agent.js |
| `views/` | 页面视图 | SceneEdit.vue, NPCEdit.vue |
| `components/` | 通用组件 | GeneratorPanel.vue |
| `api/` | API调用 | index.js |

---

## 🔧 配置说明

### 服务端配置 (server/config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  driver: sqlite
  dsn: game.db

ai:
  provider: openai
  api_key: your-api-key-here
  model: gpt-4

generator:
  enabled: true
  provider: openai
  api_key: your-generator-api-key
  model: gpt-4-turbo

game:
  max_players: 100
  tick_rate: 20
```

---

## 📈 开发流程

### 使用自循环Agent开发

```bash
# 1. 分析当前状态
python agent/main.py --show-history

# 2. 运行一次迭代
python agent/main.py -n 1 --enable-llm

# 3. 查看生成的报告
cat agent/reports/report_*.md

# 4. 运行多次迭代
python agent/main.py -n 5 --enable-llm
```

### 手动开发流程

```bash
# 1. 创建功能分支
git checkout -b feature/new-feature

# 2. 开发功能
# ...

# 3. 运行测试
cd server && go test ./internal/...
cd client && npm test
cd gm && npm test

# 4. 提交代码
git add .
git commit -m "feat: add new feature"

# 5. 推送并创建PR
git push origin feature/new-feature
```

---

## 📊 项目统计

### 代码规模

| 模块 | 文件数 | 代码行数 |
|------|--------|----------|
| Server | 54个 | 17,170行 |
| Client | 12个 | 10,007行 |
| GM | 85个 | 9,274行 |
| Agent | 20+个 | 3,000+行 |
| **总计** | **151+个** | **36,451+行** |

### 功能统计

- **API端点**: 40+个
- **MCP工具**: 37个
- **数据模型**: 20+个
- **测试用例**: 573个
- **Vue组件**: 30+个
- **Pinia Store**: 12个

---

## 🛠️ 部署

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

# 运行服务端
./gameserver
```

---

## 📚 相关文档

- [AGENT.md](AGENT.md) - Agent系统详细文档
- [SELF_LOOP_REPORT.md](SELF_LOOP_REPORT.md) - 自循环开发报告
- [IMPROVEMENT_BACKLOG.json](IMPROVEMENT_BACKLOG.json) - 改进待办列表
- [agent/requirements/](agent/requirements/) - 需求文档

---

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 运行测试确保通过
4. 提交更改 (`git commit -m 'feat: add AmazingFeature'`)
5. 推送到分支 (`git push origin feature/AmazingFeature`)
6. 创建 Pull Request

---

## 📄 许可证

MIT License

---

## 👥 联系方式

- GitHub: [@ddc-111](https://github.com/ddc-111)
- 项目链接: [https://github.com/ddc-111/agentGame](https://github.com/ddc-111/agentGame)

---

## 🙏 致谢

- [Phaser](https://phaser.io/) - 游戏引擎
- [Vue.js](https://vuejs.org/) - 前端框架
- [Gin](https://gin-gonic.com/) - Go Web框架
- [GORM](https://gorm.io/) - Go ORM框架
- [Element Plus](https://element-plus.org/) - Vue UI组件库

---

*最后更新: 2026-06-18*