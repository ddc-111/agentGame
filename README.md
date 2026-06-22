# AgentGame - AI驱动古风RPG游戏框架

一个基于AI Agent的古风RPG游戏编辑系统，提供完整的游戏服务端、Phaser 3游戏客户端、Vue 3可视化GM编辑器，以及MCP协议接口供AI直接操作游戏配置。

---

## 目录

- [项目架构](#项目架构)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [游戏系统详解](#游戏系统详解)
- [服务端API](#服务端api)
- [MCP协议接口](#mcp协议接口)
- [GM编辑器](#gm编辑器)
- [游戏客户端](#游戏客户端)
- [自循环Agent系统](#自循环agent系统)
- [部署指南](#部署指南)
- [配置说明](#配置说明)
- [测试](#测试)
- [项目统计](#项目统计)
- [贡献指南](#贡献指南)

---

## 项目架构

```
agentGame/
├── server/                    # 游戏服务端 (Go)
│   ├── cmd/gameserver/        # 入口
│   └── internal/
│       ├── agent/             # AI对话系统 (OpenAI Chat)
│       ├── config/            # YAML配置加载
│       ├── database/          # GORM数据层
│       │   ├── models/        # 18个数据模型
│       │   ├── repository/    # 数据访问层 (~80个方法)
│       │   ├── migrations/    # 数据库迁移
│       │   ├── demos/         # 演示场景数据
│       │   └── seed.go        # 数据库种子
│       ├── game/              # 游戏逻辑
│       │   ├── combat.go      # 回合制战斗
│       │   ├── inventory.go   # 背包/装备
│       │   ├── skills.go      # 技能系统
│       │   ├── achievements.go# 成就系统
│       │   ├── npc_behavior.go# NPC自主行为
│       │   └── savegame.go    # 存档系统
│       ├── generator/         # AI配置生成
│       ├── mcp/               # MCP JSON-RPC服务 (38个工具)
│       ├── network/           # HTTP/WebSocket网络层
│       └── tests/             # 集成测试
│
├── client/                    # 游戏客户端 (Phaser 3)
│   └── src/
│       ├── main.js            # 入口
│       └── game/
│           ├── scenes/        # BootScene, GameScene
│           ├── systems/       # CombatManager, InventoryManager
│           └── ui/            # CombatUI, InventoryUI, MiniMap, SaveLoadUI, SkillBar
│
├── gm/                        # GM管理端 (Vue 3)
│   └── src/
│       ├── stores/            # 11个Pinia Store
│       ├── views/             # 11个功能模块 (30个页面)
│       ├── components/        # 通用组件
│       ├── api/               # API封装
│       └── router.js          # 30条路由
│
├── agent/                     # 自循环Agent (Python)
│   ├── analyzers/             # 代码/测试/差距分析器
│   ├── executors/             # 任务/构建/测试执行器
│   ├── generators/            # 任务/报告生成器
│   ├── utils/                 # LLM客户端、历史追踪
│   └── main.py                # CLI入口
│
├── deploy/                    # Docker部署配置
├── workflows/                 # 自动化工作流
└── .github/workflows/         # CI/CD
```

---

## 技术栈

| 层级 | 技术 | 用途 |
|------|------|------|
| **服务端** | Go 1.23+, Gin, GORM, gorilla/websocket | HTTP API, WebSocket, ORM |
| **客户端** | Phaser 3, Vite | 2D游戏引擎, 构建工具 |
| **GM编辑器** | Vue 3, Element Plus, Pinia, Vue Router | 可视化配置编辑 |
| **AI集成** | OpenAI API (gpt-4/gpt-4-turbo) | NPC对话, 配置生成 |
| **数据库** | SQLite (默认) / MySQL | 数据持久化 |
| **Agent** | Python 3.10+, requests, PyYAML | 自循环开发系统 |
| **部署** | Docker, Nginx | 容器化部署 |

---

## 快速开始

### 环境要求

- Go >= 1.23
- Node.js >= 18
- Python >= 3.10 (可选，用于自循环Agent)

### 一键启动 (Windows)

```bash
# 启动全部服务
start.bat

# 停止全部服务
stop.bat
```

### 手动启动

```bash
# 1. 配置服务端
cp server/config.example.yaml server/config.yaml
# 编辑 server/config.yaml 填入AI API Key

# 2. 启动服务端 (自动建表 + 种子数据)
cd server && go run cmd/gameserver/main.go
# -> http://localhost:8080

# 3. 启动客户端
cd client && npm install && npm run dev
# -> http://localhost:5173

# 4. 启动GM编辑器
cd gm && npm install && npm run dev
# -> http://localhost:5174
```

### 构建

```bash
cd server && go build -o bin/gameserver.exe cmd/gameserver/main.go
cd client && npm run build
cd gm && npm run build
```

---

## 游戏系统详解

### 1. 回合制战斗系统

文件: `server/internal/game/combat.go`

回合制战斗引擎，支持5种敌人模板（狼、山贼、熊、虎、幽灵），每种有不同的属性和掉落。

**核心机制:**
- **伤害公式:** `(攻击力 - 防御力) × 随机(0.8~1.2)`，最低1点
- **逃跑概率:** 50%基础 + 5%/等级，上限90%
- **状态效果:** 眩晕(跳过回合)、中毒(持续伤害)、防御提升、攻击提升
- **奖励:** 经验值 + 金币 + 随机道具掉落

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/combat/start` | 发起战斗 |
| POST | `/api/combat/action` | 执行战斗动作 (攻击/使用道具/逃跑) |

### 2. 背包与装备系统

文件: `server/internal/game/inventory.go`

JSON背包管理，支持道具堆叠、装备穿脱、道具使用。

**核心机制:**
- **背包:** `{item_id: count}` 映射，支持添加/移除/使用
- **装备槽:** 武器 + 护甲，穿新装备自动卸旧装备回背包
- **属性计算:** 基础属性 + 装备加成 = 总属性

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/inventory/:player_id` | 获取背包 |
| POST | `/api/inventory/equip` | 装备道具 |
| POST | `/api/inventory/unequip` | 卸下装备 |
| POST | `/api/inventory/use` | 使用道具 |

### 3. 技能系统

文件: `server/internal/game/skills.go`

4种技能类型: 攻击(造成伤害)、治疗(恢复HP)、增益(自身buff)、减益(敌方debuff)。

**核心机制:**
- **MP消耗:** 每个技能有独立MP消耗
- **等级限制:** 技能有最低等级要求
- **冷却:** 技能使用后有冷却回合数
- **伤害公式:** `(技能伤害 + 攻击力 - 防御力) × 随机(0.8~1.2)`

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/skills` | 技能列表 (分页) |
| POST | `/api/skills/use` | 使用技能 |

### 4. 成就系统

文件: `server/internal/game/achievements.go`

8种条件类型，10个预置成就。

**条件类型:**
| 类型 | 说明 |
|------|------|
| `combat_win` | 战斗胜利次数 |
| `gold` | 累计金币 |
| `level` | 达到等级 |
| `quest_complete` | 完成任务 (可指定任务code) |
| `explore` | 探索场景数 |
| `collect` | 收集道具数 |
| `talk_all_npcs` | 与所有NPC对话 |
| `skill_use` | 使用技能次数 |

**预置成就:** 初次任务、小镇之友、富有、战斗大师、探索者、收藏家、初次战斗、技能大师、10级、20级

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/achievements/:player_id` | 获取玩家成就 |
| POST | `/api/achievements/check` | 检查/授予成就 |

### 5. NPC自主行为系统

文件: `server/internal/game/npc_behavior.go`

NPC根据时间表自主切换状态，对玩家行为做出情绪反应。

**NPC状态:** 空闲、巡逻、交谈、交易、逃离

**情绪系统:** 开心、中性、生气、害怕

**核心机制:**
- **时间驱动:** 游戏循环每分钟tick，每小时检查NPC时间表
- **状态映射:** `open_shop` → 交易, `patrol` → 巡逻, `rest` → 空闲
- **玩家交互:** 对话→交谈, 赠礼→开心, 攻击→生气+逃离
- **对话上下文:** 根据NPC状态和情绪生成AI对话提示词

**WebSocket广播:** NPC状态变化时自动广播给同场景玩家

### 6. 存档系统

文件: `server/internal/game/savegame.go`

最多11个存档槽 (0-10)，槽0为自动存档。

**自动存档触发:**
- 切换场景
- 完成任务
- 升级
- 战斗胜利

**存档内容:** 玩家完整状态快照 (属性、背包、装备、位置、进度)

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/save` | 保存游戏 |
| GET | `/api/saves/:player_id` | 获取存档列表 |
| POST | `/api/load/:save_id` | 加载存档 |

### 7. AI对话系统

文件: `server/internal/agent/chat.go`, `memory.go`

基于OpenAI的NPC智能对话，支持滑动窗口记忆和对话摘要。

**记忆模式:**
- **滑动窗口:** 保留最近N条消息
- **摘要模式:** 超出窗口时自动生成对话摘要
- **数据库持久化:** 对话历史存储在Conversation表

**对话上下文:** NPC状态、情绪、玩家信息自动注入AI提示词

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/npc/chat` | 与NPC对话 |

### 8. AI配置生成器

文件: `server/internal/generator/generator.go`

使用GPT-4自动生成游戏配置，支持8种类型: NPC、场景、任务、商店、道具、智能体、对话、流程。

**支持操作:** create (创建)、complete (补全)、expand (扩展)

**API端点:**
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/generator/generate` | 生成配置 |
| GET | `/api/generator/status` | 生成器状态 |
| POST | `/api/generator/test` | 测试生成器 |

---

## 服务端API

### 认证

GM管理端使用JWT认证:
```
POST /api/gm/login          # 登录获取token
GET  /api/gm/me             # 获取当前用户 (需Authorization头)
```

### 全局CRUD端点

所有实体提供标准CRUD:

| 实体 | 端点前缀 | 说明 |
|------|----------|------|
| 场景 | `/api/scenes` | 场景管理 |
| NPC | `/api/npcs` | NPC管理 |
| 智能体 | `/api/agents` | AI智能体配置 |
| LLM提供商 | `/api/llm/providers` | 大模型接入配置 |
| 提示词模板 | `/api/prompts` | 提示词管理 |
| 商店 | `/api/shops` | 商店管理 |
| 道具 | `/api/items` | 道具管理 |
| 任务 | `/api/tasks` | 任务管理 |
| 流程 | `/api/flows` | 对话/任务流程 |
| 玩家 | `/api/players` | 玩家管理 |
| 对话 | `/api/conversations` | 对话记录 |

### 游戏API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/game/init` | 初始化游戏 |
| GET | `/api/game/scene/:code` | 按code获取场景 |
| GET | `/api/game/npc/:code` | 按code获取NPC |
| GET | `/api/game/shop/:code/items` | 获取商店商品 |
| POST | `/api/game/tick` | 推进游戏时钟 |
| POST | `/api/shop/buy` | 购买道具 |
| POST | `/api/npc/chat` | 与NPC对话 |

### 数据导入导出

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/export` | 导出全部数据 |
| POST | `/api/import` | 导入数据 |

### WebSocket

```
GET /api/ws
```

支持房间制广播，玩家加入场景房间后接收该场景的实时事件 (NPC状态变化、战斗结果等)。

### 中间件

- **速率限制:** 普通API 100/200 req/min, NPC对话 10/20, 战斗 30/50
- **CORS:** 支持 localhost:5173 和 localhost:5174
- **JWT:** GM端认证
- **超时:** 30秒请求超时
- **结构化日志:** JSON格式日志输出

---

## MCP协议接口

### 端点

```
POST http://localhost:8080/mcp          # JSON-RPC 2.0
GET  http://localhost:8080/mcp/sse      # SSE流式
```

REST封装:
```
GET  /api/mcp/tools       # 列出工具
POST /api/mcp/call        # 调用工具
GET  /api/mcp/resources   # 列出资源
GET  /api/mcp/prompts     # 列出提示词
```

### 38个MCP工具

#### 场景管理 (5)
| 工具 | 说明 |
|------|------|
| `list_scenes` | 获取所有场景 |
| `get_scene` | 获取场景详情 |
| `create_scene` | 创建场景 |
| `update_scene` | 更新场景 |
| `delete_scene` | 删除场景 |

#### NPC管理 (5)
| 工具 | 说明 |
|------|------|
| `list_npcs` | 获取所有NPC |
| `get_npc` | 获取NPC详情 |
| `create_npc` | 创建NPC |
| `update_npc` | 更新NPC |
| `delete_npc` | 删除NPC |

#### 智能体管理 (5)
| 工具 | 说明 |
|------|------|
| `list_agents` | 获取所有智能体 |
| `get_agent` | 获取智能体详情 |
| `create_agent` | 创建智能体 |
| `update_agent` | 更新智能体 |
| `delete_agent` | 删除智能体 |

#### 商店管理 (5)
| 工具 | 说明 |
|------|------|
| `list_shops` | 获取所有商店 |
| `get_shop` | 获取商店详情 |
| `create_shop` | 创建商店 |
| `update_shop` | 更新商店 |
| `delete_shop` | 删除商店 |

#### 道具管理 (5)
| 工具 | 说明 |
|------|------|
| `list_items` | 获取所有道具 |
| `get_item` | 获取道具详情 |
| `create_item` | 创建道具 |
| `update_item` | 更新道具 |
| `delete_item` | 删除道具 |

#### 任务管理 (5)
| 工具 | 说明 |
|------|------|
| `list_tasks` | 获取所有任务 |
| `get_task` | 获取任务详情 |
| `create_task` | 创建任务 |
| `update_task` | 更新任务 |
| `delete_task` | 删除任务 |

#### 流程管理 (2)
| 工具 | 说明 |
|------|------|
| `list_flows` | 获取所有流程 |
| `create_flow` | 创建流程 |

#### 提示词模板 (2)
| 工具 | 说明 |
|------|------|
| `list_templates` | 获取所有模板 |
| `create_template` | 创建模板 |

#### 工具函数 (4)
| 工具 | 说明 |
|------|------|
| `generate_config` | AI生成游戏配置 (支持npc/scene/task/shop/item/agent/dialogue/flow) |
| `export_data` | 导出全部数据 |
| `get_game_stats` | 获取数据统计 |

### MCP资源 (7)

| URI | 说明 |
|-----|------|
| `game_state://scenes` | 场景列表 |
| `game_state://npcs` | NPC列表 |
| `game_state://agents` | 智能体列表 |
| `game_state://shops` | 商店列表 |
| `game_state://items` | 道具列表 |
| `game_state://tasks` | 任务列表 |
| `game_state://overview` | 数据统计概览 |

### MCP提示词 (4)

| 名称 | 说明 | 参数 |
|------|------|------|
| `npc_personality` | 生成NPC性格设定 | name, title, background |
| `npc_dialogue` | 生成NPC对话风格 | name, personality, scenario |
| `scene_description` | 生成场景描述 | name, type, atmosphere |
| `quest_design` | 生成任务设计 | theme, difficulty, npc_involved |

### MCP调用示例

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_npc",
    "arguments": {
      "name": "李铁匠",
      "code": "npc_blacksmith_li",
      "title": "铁匠铺老板",
      "description": "一位技艺精湛的铁匠"
    }
  }
}
```

---

## GM编辑器

### 功能模块 (11个, 30个页面)

#### 场景编辑器
| 页面 | 路由 | 功能 |
|------|------|------|
| 场景列表 | `/scene/list` | 场景CRUD |
| 场景编辑 | `/scene/edit/:id` | 场景详情、NPC放置、传送点配置 |
| 地图素材 | `/scene/tileset` | 地形/建筑/装饰素材管理 |

#### NPC编辑器
| 页面 | 路由 | 功能 |
|------|------|------|
| NPC列表 | `/npc/list` | NPC CRUD |
| NPC编辑 | `/npc/edit/:id` | NPC属性、位置、智能体绑定、商店绑定 |
| 对话树 | `/npc/dialogue/:id` | 节点式对话编辑器 |
| 行为编辑 | `/npc/behavior/:id` | NPC时间表、情绪、行为树 |

#### 智能体配置
| 页面 | 路由 | 功能 |
|------|------|------|
| 智能体列表 | `/agent/list` | 智能体CRUD |
| 智能体编辑 | `/agent/edit/:id` | LLM配置、系统提示词、知识库、工具 |
| 记忆配置 | `/agent/memory/:id` | 滑动窗口、摘要、向量搜索配置 |

#### LLM配置
| 页面 | 路由 | 功能 |
|------|------|------|
| 提供商配置 | `/llm/provider` | 添加/编辑/删除LLM提供商 |
| 模型配置 | `/llm/model` | 默认模型、按场景分配模型 |
| 连接测试 | `/llm/test` | 测试LLM连接 |

#### 提示词管理
| 页面 | 路由 | 功能 |
|------|------|------|
| 模板管理 | `/prompt/template` | 提示词模板CRUD |
| 变量管理 | `/prompt/variable` | 模板变量定义 |
| 模板测试 | `/prompt/test` | 模板渲染预览 |

#### 商店配置
| 页面 | 路由 | 功能 |
|------|------|------|
| 商店列表 | `/shop/list` | 商店CRUD |
| 商店编辑 | `/shop/edit/:id` | 商店详情、商品配置 |
| 道具管理 | `/shop/items` | 道具CRUD |

#### 任务系统
| 页面 | 路由 | 功能 |
|------|------|------|
| 任务列表 | `/task/list` | 任务CRUD |
| 任务编辑 | `/task/edit/:id` | 触发条件、目标、奖励 |
| 流程编辑 | `/task/flow/:id` | 节点式流程编辑器 |

#### 技能系统
| 页面 | 路由 | 功能 |
|------|------|------|
| 技能编辑 | `/skill/list` | 技能CRUD |
| 技能树 | `/skill/tree` | 技能树可视化 |

#### 成就系统
| 页面 | 路由 | 功能 |
|------|------|------|
| 成就编辑 | `/achievement/list` | 成就CRUD、条件配置 |

#### 系统配置
| 页面 | 路由 | 功能 |
|------|------|------|
| 游戏配置 | `/config/game` | 基础/玩家/世界/战斗/经济设置 |
| 战斗配置 | `/config/combat` | 敌人编辑、公式配置、战斗模拟 |
| 数据导出 | `/config/export` | 导出JSON |
| 数据导入 | `/config/import` | 导入JSON |

#### 演示
| 页面 | 路由 | 功能 |
|------|------|------|
| 演示播放 | `/demo/player` | 分步演示播放器 |

### AI生成助手

GM编辑器内置浮动AI助手面板，可调用后端Generator API自动生成:
- NPC配置 (性格、对话、行为)
- 场景配置 (描述、布局)
- 任务配置 (目标、奖励)
- 商店/道具配置
- 对话流程

---

## 游戏客户端

### 场景系统

- **BootScene:** 资源加载、程序化纹理生成
- **GameScene:** 主游戏场景 (1280行)
  - 玩家移动 (WASD/方向键)
  - NPC交互 (点击对话)
  - 传送点 (场景切换)
  - 背景渲染 (程序化生成)
  - 小地图
  - 随机遭遇

### UI组件

| 组件 | 功能 |
|------|------|
| `CombatUI` | 战斗覆盖层 (敌人信息、动作菜单、伤害动画) |
| `InventoryUI` | 背包面板 (道具列表、装备穿脱、详情) |
| `MiniMap` | 小地图 (场景概览、NPC/传送点标记) |
| `SaveLoadUI` | 存档/读档界面 (11个槽位) |
| `SkillBar` | 技能快捷栏 |

### 游戏系统 (客户端)

| 系统 | 文件 | 功能 |
|------|------|------|
| 战斗管理 | `CombatManager.js` | 回合制战斗逻辑、状态效果、奖励 |
| 背包管理 | `InventoryManager.js` | 道具管理、装备穿脱、属性计算 |

---

## 自循环Agent系统

### 概述

Python实现的AI驱动开发系统，可自动分析代码、运行测试、生成改进任务、执行重构。

### 模块

| 模块 | 功能 |
|------|------|
| `analyzers/code_analyzer.py` | 代码结构分析 |
| `analyzers/test_analyzer.py` | 测试覆盖分析 |
| `analyzers/gap_analyzer.py` | 差距识别 |
| `executors/task_executor.py` | LLM辅助任务执行 |
| `executors/test_executor.py` | Go/Vitest测试运行 |
| `executors/build_executor.py` | 构建验证 |
| `generators/task_generator.py` | 任务生成 |
| `generators/report_generator.py` | 报告生成 |
| `utils/llm_client.py` | OpenAI兼容API客户端 |
| `utils/history_tracker.py` | 迭代历史记录 |

### 使用

```bash
# 安装依赖
pip install -r agent/requirements.txt

# 运行单次迭代
python agent/main.py -n 1

# 运行5次迭代，启用LLM
python agent/main.py -n 5 --enable-llm

# 只生成报告
python agent/main.py --report-only

# 查看历史
python agent/main.py --show-history
```

### 配置

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
  priority_filter: [critical, high, medium]
```

---

## 部署指南

### Docker部署

```bash
# 生产环境
cd deploy
docker-compose up -d

# 开发环境 (热重载)
cd deploy
docker-compose -f docker-compose.dev.yml up -d
```

### 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Server | 8080 | Go API + WebSocket + MCP |
| Client | 5173 | Phaser游戏客户端 |
| GM | 5174 | Vue管理编辑器 |

### Nginx配置

`deploy/` 目录包含两个Nginx配置:
- `nginx-client.conf` - 客户端 (代理API/WebSocket)
- `nginx-gm.conf` - GM编辑器 (代理API)

### 手动部署

```bash
# 构建
cd server && go build -o bin/gameserver cmd/gameserver/main.go
cd client && npm run build
cd gm && npm run build

# 运行 (需要config.yaml)
./server/bin/gameserver
# 静态文件: client/dist/, gm/dist/
```

---

## 配置说明

### 服务端配置 (server/config.yaml)

```yaml
server:
  port: 8080
  mode: debug           # debug/release
  log_level: info

database:
  driver: sqlite         # sqlite/mysql
  dsn: game.db

# 游戏内NPC AI配置
ai:
  provider: openai
  base_url: https://api.openai.com/v1
  api_key: your-api-key
  model: gpt-4
  temperature: 0.7
  max_tokens: 500

# 配置生成AI (可指向不同模型)
generator:
  enabled: true
  provider: openai
  base_url: https://api.openai.com/v1
  api_key: your-generator-api-key
  model: gpt-4-turbo
  temperature: 0.7
  max_tokens: 4000

auth:
  jwt_secret: change-me-in-production
  token_expiry: 24
  gm_username: admin
  gm_password: admin123

cors:
  allowed_origins:
    - http://localhost:5173
    - http://localhost:5174
```

### 预置数据

首次启动自动种子:
- **4个场景:** 小镇中心、杂货铺、铁匠铺、茶摊
- **3个NPC:** 李掌柜、王大娘、张铁匠
- **3个智能体:** 对应NPC的AI配置
- **2个商店:** 李记杂货铺、张记铁匠铺
- **8个道具:** 草药、灵芝、馒头、烧酒、铁剑等
- **2个任务:** 初来乍到、装备自己
- **1个流程:** NPC出门购物流程
- **2个提示词模板:** NPC对话、场景描述

---

## 测试

### 运行测试

```bash
# 服务端
cd server && go test ./internal/... -v

# 客户端
cd client && npm test

# GM编辑器
cd gm && npm test

# 全部 (Windows)
test.bat
```

### 测试覆盖

| 模块 | 测试文件 | 说明 |
|------|----------|------|
| Server | `internal/tests/` | API集成、认证、WebSocket、MCP、内存、迁移、NPC行为、分页、技能成就、验证、性能基准 |
| Client | `src/__tests__/` | BootScene、GameScene、CombatManager、InventoryManager |
| GM | `src/__tests__/` | Store测试 (scene/npc/agent/shop/llm/prompt/skill/task/config/demo/achievement)、组件测试 |

---

## 项目统计

### 代码规模

| 模块 | 文件数 | 主要语言 |
|------|--------|----------|
| Server | ~50个Go文件 | Go |
| Client | ~14个源文件 | JavaScript |
| GM | ~58个源文件 | Vue/JavaScript |
| Agent | ~22个Python文件 | Python |

### 数据模型 (18个)

Scene, SceneNPC, Portal, NPC, Agent, LLMProvider, PromptTemplate, Shop, ShopItem, Item, Task, Flow, GameConfig, Player, Conversation, SaveGame, Skill, Achievement, PlayerAchievement, PlayerConversationContext, GMUser

### API端点

- **CRUD:** 11个实体 × 5个端点 = 55个
- **游戏API:** ~15个专用端点
- **系统API:** ~10个 (健康检查、Swagger、认证、导入导出)
- **总计:** ~70+个HTTP端点

### MCP接口

- **工具:** 38个
- **资源:** 7个
- **提示词:** 4个

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [AGENT.md](AGENT.md) | MCP协议完整文档 (AI可读) |
| [AGENTS.md](AGENTS.md) | 多Agent协作文档 |

---

## 贡献指南

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 运行测试确保通过
4. 提交更改 (`git commit -m 'feat: add AmazingFeature'`)
5. 推送到分支 (`git push origin feature/AmazingFeature`)
6. 创建Pull Request

### 代码规范

- 实体code使用snake_case，带前缀: `scene_*`, `npc_*`, `agent_*`, `shop_*`, `item_*`, `task_*`, `flow_*`
- 所有模型继承BaseModel (ID, CreatedAt, UpdatedAt, DeletedAt)
- 删除操作使用软删除 (GORM DeletedAt)
- JSON字段在数据库中以字符串存储

---

## 许可证

MIT License

---

## 致谢

- [Phaser](https://phaser.io/) - 游戏引擎
- [Vue.js](https://vuejs.org/) - 前端框架
- [Element Plus](https://element-plus.org/) - UI组件库
- [Gin](https://gin-gonic.com/) - Go Web框架
- [GORM](https://gorm.io/) - Go ORM框架

---

*最后更新: 2026-06-22*
