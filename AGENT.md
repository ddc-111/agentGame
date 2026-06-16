# AgentGame - AI可操作的游戏编辑系统

> 本文档为自描述文件，供AI模型理解和操作系统

## 系统概述

这是一个古风RPG游戏编辑系统，提供完整的MCP接口，AI可以直接操作游戏配置。

## MCP服务端点

```
POST http://localhost:8080/mcp
```

## 可用工具列表

### 场景管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_scenes` | 获取所有场景 | - |
| `get_scene` | 获取场景详情 | `id` |
| `create_scene` | 创建场景 | `name`, `code` |
| `update_scene` | 更新场景 | `id` |
| `delete_scene` | 删除场景 | `id` |

### NPC管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_npcs` | 获取所有NPC | - |
| `get_npc` | 获取NPC详情 | `id` |
| `create_npc` | 创建NPC | `name`, `code` |
| `update_npc` | 更新NPC | `id` |
| `delete_npc` | 删除NPC | `id` |

### 智能体管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_agents` | 获取所有智能体 | - |
| `get_agent` | 获取智能体详情 | `id` |
| `create_agent` | 创建智能体 | `name`, `code`, `system_prompt` |
| `update_agent` | 更新智能体 | `id` |
| `delete_agent` | 删除智能体 | `id` |

### 商店管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_shops` | 获取所有商店 | - |
| `get_shop` | 获取商店详情 | `id` |
| `create_shop` | 创建商店 | `name`, `code` |
| `update_shop` | 更新商店 | `id` |
| `delete_shop` | 删除商店 | `id` |

### 道具管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_items` | 获取所有道具 | - |
| `get_item` | 获取道具详情 | `id` |
| `create_item` | 创建道具 | `name`, `code` |
| `update_item` | 更新道具 | `id` |
| `delete_item` | 删除道具 | `id` |

### 任务管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_tasks` | 获取所有任务 | - |
| `get_task` | 获取任务详情 | `id` |
| `create_task` | 创建任务 | `name`, `code` |
| `update_task` | 更新任务 | `id` |
| `delete_task` | 删除任务 | `id` |

### 流程管理

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_flows` | 获取所有流程 | - |
| `create_flow` | 创建流程 | `name`, `code` |

### 提示词模板

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `list_templates` | 获取所有模板 | - |
| `create_template` | 创建模板 | `name`, `code`, `content` |

### AI生成

| 工具名 | 描述 | 必需参数 |
|--------|------|----------|
| `generate_config` | AI生成配置 | `type`, `description` |

生成类型: `npc`, `scene`, `task`, `shop`, `item`, `agent`, `dialogue`, `flow`

### 数据操作

| 工具名 | 描述 |
|--------|------|
| `export_data` | 导出所有数据 |
| `get_game_stats` | 获取数据统计 |

## MCP调用示例

### 初始化连接

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {}
  }
}
```

### 获取工具列表

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list"
}
```

### 调用工具

```json
{
  "jsonrpc": "2.0",
  "id": 3,
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

### AI生成配置

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "generate_config",
    "arguments": {
      "type": "npc",
      "description": "一个卖包子的老大爷，性格豪爽，喜欢讲故事"
    }
  }
}
```

## 数据模型

### Scene (场景)
```json
{
  "id": 1,
  "name": "小镇中心",
  "code": "scene_town_center",
  "description": "古风小镇的中心广场",
  "background": "town_center.png",
  "width": 1920,
  "height": 1080
}
```

### NPC (非玩家角色)
```json
{
  "id": 1,
  "name": "李掌柜",
  "code": "npc_merchant_li",
  "title": "杂货铺老板",
  "description": "一位精明的中年商人",
  "agent_id": 1,
  "shop_id": 1
}
```

### Agent (智能体)
```json
{
  "id": 1,
  "name": "李掌柜智能体",
  "code": "agent_merchant_li",
  "system_prompt": "你是李掌柜，一位古风小镇杂货铺的老板...",
  "llm_provider": "openai",
  "llm_model": "gpt-4",
  "temperature": 0.7
}
```

### Shop (商店)
```json
{
  "id": 1,
  "name": "李记杂货铺",
  "code": "shop_general_store",
  "type": "general",
  "owner_npc": "npc_merchant_li",
  "open_time": "06:00",
  "close_time": "22:00"
}
```

### Item (道具)
```json
{
  "id": 1,
  "name": "草药",
  "code": "item_herb",
  "category": "medicine",
  "description": "普通的草药，可恢复少量生命",
  "effect": "{\"hp\":20}"
}
```

### Task (任务)
```json
{
  "id": 1,
  "name": "初来乍到",
  "code": "task_first_arrival",
  "type": "main",
  "description": "新来的冒险者，先去杂货铺买些必需品吧",
  "status": "active"
}
```

## 现有数据概览

系统已预置以下数据:

- **场景**: 小镇中心、杂货铺、铁匠铺、茶摊
- **NPC**: 李掌柜、王大娘、张铁匠
- **智能体**: 对应NPC的AI智能体
- **商店**: 李记杂货铺、张记铁匠铺
- **道具**: 草药、灵芝、馒头、烧酒、铁剑等
- **任务**: 初来乍到、装备自己
- **流程**: NPC出门购物流程

## 快速开始

1. 启动服务端: `cd server && go run cmd/gameserver/main.go`
2. MCP端点: `POST http://localhost:8080/mcp`
3. 发送工具调用请求

## 注意事项

- 所有ID为数字类型
- JSON字段使用字符串存储
- 删除操作为软删除
- 生成智能体需要配置API Key
