package database

import (
	"log"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	log.Println("Seeding database...")

	// 检查是否已有数据
	var count int64
	db.Model(&models.Scene{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// 创建场景
	scenes := []models.Scene{
		{
			Name:        "小镇中心",
			Code:        "scene_town_center",
			Description: "古风小镇的中心广场，人来人往",
			Background:  "town_center.png",
			Width:       1920,
			Height:      1080,
		},
		{
			Name:        "杂货铺",
			Code:        "scene_general_store",
			Description: "售卖各种日常用品的店铺",
			Background:  "shop_general.png",
			Width:       1280,
			Height:      720,
		},
		{
			Name:        "铁匠铺",
			Code:        "scene_blacksmith",
			Description: "打造和售卖兵器的地方",
			Background:  "shop_blacksmith.png",
			Width:       1280,
			Height:      720,
		},
		{
			Name:        "茶摊",
			Code:        "scene_tea_stand",
			Description: "路边的茶摊，可以休息喝茶",
			Background:  "tea_stand.png",
			Width:       800,
			Height:      600,
		},
	}

	for i := range scenes {
		if err := db.Create(&scenes[i]).Error; err != nil {
			return err
		}
	}

	// 创建NPC
	npcs := []models.NPC{
		{
			Name:        "李掌柜",
			Code:        "npc_merchant_li",
			Title:       "杂货铺老板",
			Description: "一位精明的中年商人，经营着镇上最大的杂货铺",
			Avatar:      "merchant_li.png",
			Sprite:      "merchant_li_sprite.png",
			Behaviors:   `["idle","greet","sell"]`,
			Schedule:    `[{"time":"06:00","action":"open_shop","scene":"scene_general_store"},{"time":"22:00","action":"close_shop","scene":"scene_town_center"}]`,
		},
		{
			Name:        "王大娘",
			Code:        "npc_tea_wang",
			Title:       "茶摊老板娘",
			Description: "热情好客的中年妇人，卖的茶远近闻名",
			Avatar:      "tea_wang.png",
			Sprite:      "tea_wang_sprite.png",
			Behaviors:   `["idle","chat","serve_tea"]`,
			Schedule:    `[]`,
		},
		{
			Name:        "张铁匠",
			Code:        "npc_blacksmith_zhang",
			Title:       "铁匠铺师傅",
			Description: "技艺精湛的老铁匠，打造的兵器远近闻名",
			Avatar:      "blacksmith_zhang.png",
			Sprite:      "blacksmith_zhang_sprite.png",
			Behaviors:   `["idle","forge","sell"]`,
			Schedule:    `[{"time":"08:00","action":"open_shop","scene":"scene_blacksmith"},{"time":"20:00","action":"close_shop","scene":"scene_town_center"}]`,
		},
	}

	for i := range npcs {
		if err := db.Create(&npcs[i]).Error; err != nil {
			return err
		}
	}

	// 创建智能体
	agents := []models.Agent{
		{
			Name:           "李掌柜智能体",
			Code:           "agent_merchant_li",
			Description:    "杂货铺老板的AI智能体，负责与玩家进行智能对话",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.7,
			MaxTokens:      500,
			SystemPrompt:   `你是李掌柜，一位古风小镇杂货铺的老板。你性格精明但诚实，对顾客热情周到。你经营各种日常用品、药材、食材等。当顾客询问商品时，你会详细介绍商品的用途和价格。你偶尔会讲一些镇上的趣事。当前店铺状态：营业中，库存充足。与顾客对话时，保持古风语气，使用"客官"、"小店"等称呼。`,
			MemoryType:     "sliding_window",
			MaxMessages:    20,
			SummaryEnabled: true,
			KnowledgeBase:  `[{"id":"kb_001","name":"商品目录","type":"text"},{"id":"kb_002","name":"价格表","type":"table"}]`,
			Tools:          `[{"name":"query_inventory","description":"查询库存"},{"name":"make_deal","description":"进行交易"}]`,
		},
		{
			Name:           "王大娘智能体",
			Code:           "agent_tea_wang",
			Description:    "茶摊老板娘的AI智能体，喜欢聊天",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.8,
			MaxTokens:      300,
			SystemPrompt:   `你是王大娘，一位古风小镇茶摊的老板娘。你性格热情开朗，喜欢与人聊天。你卖各种茶水，也提供一些小点心。你知道镇上很多八卦和故事。经常主动与路人搭话，分享镇上的新鲜事。与客人对话时，使用亲切的语气，如"哎呀"、"客官"等。`,
			MemoryType:     "sliding_window",
			MaxMessages:    15,
			SummaryEnabled: false,
			KnowledgeBase:  `[{"id":"kb_003","name":"茶品目录","type":"text"},{"id":"kb_004","name":"镇上故事","type":"text"}]`,
			Tools:          `[]`,
		},
		{
			Name:           "张铁匠智能体",
			Code:           "agent_blacksmith_zhang",
			Description:    "铁匠铺师傅的AI智能体，精通锻造",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.6,
			MaxTokens:      400,
			SystemPrompt:   `你是张铁匠，一位古风小镇铁匠铺的老师傅。你性格沉稳，做事认真。你精通各种兵器的打造，对每件作品都精益求精。你说话直接，不喜欢拐弯抹角。当客人询问兵器时，你会详细介绍材质和工艺。与客人对话时，使用朴实的语气。`,
			MemoryType:     "sliding_window",
			MaxMessages:    20,
			SummaryEnabled: true,
			KnowledgeBase:  `[{"id":"kb_005","name":"兵器谱","type":"text"},{"id":"kb_006","name":"锻造材料","type":"text"}]`,
			Tools:          `[{"name":"forge_weapon","description":"打造兵器"},{"name":"repair_item","description":"修理装备"}]`,
		},
	}

	for i := range agents {
		if err := db.Create(&agents[i]).Error; err != nil {
			return err
		}
	}

	// 更新NPC关联的智能体
	db.Model(&npcs[0]).Update("agent_id", agents[0].ID)
	db.Model(&npcs[1]).Update("agent_id", agents[1].ID)
	db.Model(&npcs[2]).Update("agent_id", agents[2].ID)

	// 创建道具
	items := []models.Item{
		{Name: "草药", Code: "item_herb", Category: "medicine", Description: "普通的草药，可恢复少量生命", Effect: `{"hp":20}`},
		{Name: "灵芝", Code: "item_lingzhi", Category: "medicine", Description: "珍贵的药材，可恢复大量生命", Effect: `{"hp":100}`},
		{Name: "馒头", Code: "item_mantou", Category: "food", Description: "普通的干粮，可恢复少量体力", Effect: `{"stamina":10}`},
		{Name: "烧酒", Code: "item_wine", Category: "food", Description: "烈性白酒，可驱寒保暖", Effect: `{"cold_resist":30}`},
		{Name: "麻绳", Code: "item_rope", Category: "tool", Description: "结实的麻绳，可用于攀爬", Effect: `{}`},
		{Name: "铁剑", Code: "item_iron_sword", Category: "weapon", Description: "普通的铁制长剑", Effect: `{"attack":10}`},
		{Name: "精钢刀", Code: "item_steel_blade", Category: "weapon", Description: "锻造精良的钢刀", Effect: `{"attack":25}`},
		{Name: "铁甲", Code: "item_iron_armor", Category: "armor", Description: "普通的铁制铠甲", Effect: `{"defense":15}`},
	}

	for i := range items {
		if err := db.Create(&items[i]).Error; err != nil {
			return err
		}
	}

	// 创建商店
	shops := []models.Shop{
		{
			Name:        "李记杂货铺",
			Code:        "shop_general_store",
			Type:        "general",
			Description: "售卖各种日常用品、药材、食材",
			OwnerNPC:    "npc_merchant_li",
			SceneCode:   "scene_general_store",
			OpenTime:    "06:00",
			CloseTime:   "22:00",
			Discount:    `{"threshold":3,"rate":0.9}`,
		},
		{
			Name:        "张记铁匠铺",
			Code:        "shop_blacksmith",
			Type:        "blacksmith",
			Description: "打造和售卖各种兵器、工具",
			OwnerNPC:    "npc_blacksmith_zhang",
			SceneCode:   "scene_blacksmith",
			OpenTime:    "08:00",
			CloseTime:   "20:00",
			Discount:    `{"threshold":2,"rate":0.95}`,
		},
	}

	for i := range shops {
		if err := db.Create(&shops[i]).Error; err != nil {
			return err
		}
	}

	// 创建商店商品
	shopItems := []models.ShopItem{
		{ShopID: shops[0].ID, ItemID: items[0].ID, Price: 100, Stock: 50},
		{ShopID: shops[0].ID, ItemID: items[1].ID, Price: 200, Stock: 30},
		{ShopID: shops[0].ID, ItemID: items[2].ID, Price: 50, Stock: 100},
		{ShopID: shops[0].ID, ItemID: items[3].ID, Price: 80, Stock: 40},
		{ShopID: shops[0].ID, ItemID: items[4].ID, Price: 30, Stock: 60},
		{ShopID: shops[1].ID, ItemID: items[5].ID, Price: 500, Stock: 10},
		{ShopID: shops[1].ID, ItemID: items[6].ID, Price: 800, Stock: 5},
		{ShopID: shops[1].ID, ItemID: items[7].ID, Price: 300, Stock: 20},
	}

	for i := range shopItems {
		if err := db.Create(&shopItems[i]).Error; err != nil {
			return err
		}
	}

	// 更新NPC关联的商店
	db.Model(&npcs[0]).Update("shop_id", shops[0].ID)
	db.Model(&npcs[2]).Update("shop_id", shops[1].ID)

	// 创建提示词模板
	templates := []models.PromptTemplate{
		{
			Name:     "NPC基础人设",
			Code:     "template_npc_basic",
			Category: "system",
			Content:  `你是{{npc_name}}，{{npc_title}}。{{npc_description}}当前状态：位置：{{current_scene}}，时间：{{current_time}}，心情：{{mood}}。{{#if has_shop}}你经营着一家{{shop_type}}店。{{/if}}{{#if has_task}}你有任务要交给冒险者：{{task_description}}{{/if}}请保持角色设定，用古风语气与玩家对话。`,
			Variables: `[{"name":"npc_name","type":"string","description":"NPC名称"},{"name":"npc_title","type":"string","description":"NPC称号"},{"name":"npc_description","type":"string","description":"NPC描述"},{"name":"current_scene","type":"string","description":"当前场景"},{"name":"current_time","type":"string","description":"当前时间"},{"name":"mood","type":"string","description":"NPC心情"},{"name":"has_shop","type":"boolean","description":"是否有商店"},{"name":"shop_type","type":"string","description":"商店类型"},{"name":"has_task","type":"boolean","description":"是否有任务"},{"name":"task_description","type":"string","description":"任务描述"}]`,
		},
		{
			Name:     "商人对话模板",
			Code:     "template_merchant",
			Category: "system",
			Content:  `你是{{npc_name}}，一位{{personality}}的商人。你售卖以下商品：{{#each items}}- {{this.name}}: {{this.description}} ({{this.price}}文){{/each}}交易规则：1. 可以适当还价，但不能低于成本价 2. 购买超过3件可打9折 3. 老顾客可以赊账。当前库存状态：{{#each items}}- {{this.name}}: {{this.stock}}件{{/each}}`,
			Variables: `[{"name":"npc_name","type":"string","description":"NPC名称"},{"name":"personality","type":"string","description":"性格特点"},{"name":"items","type":"array","description":"商品列表"}]`,
		},
	}

	for i := range templates {
		if err := db.Create(&templates[i]).Error; err != nil {
			return err
		}
	}

	// 创建任务
	tasks := []models.Task{
		{
			Name:        "初来乍到",
			Code:        "task_first_arrival",
			Type:        "main",
			Description: "新来的冒险者，先去杂货铺买些必需品吧",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[{"type":"player_level","operator":"==","value":1}]}`,
			Objectives:  `[{"id":"obj_001","type":"dialogue","target":"npc_merchant_li","description":"与李掌柜对话","completed":false},{"id":"obj_002","type":"collect","target":"item_herb","count":3,"description":"购买3份草药","completed":false},{"id":"obj_003","type":"collect","target":"item_mantou","count":5,"description":"购买5个馒头","completed":false}]`,
			Rewards:     `{"exp":100,"gold":500,"items":[]}`,
			NextTask:    "task_equip_yourself",
		},
		{
			Name:        "装备自己",
			Code:        "task_equip_yourself",
			Type:        "main",
			Description: "有了补给，接下来去铁匠铺买把武器吧",
			Status:      "locked",
			Trigger:     `{"type":"task_complete","conditions":[{"type":"task_id","value":"task_first_arrival"}]}`,
			Objectives:  `[{"id":"obj_004","type":"dialogue","target":"npc_blacksmith_zhang","description":"与张铁匠对话","completed":false},{"id":"obj_005","type":"collect","target":"item_iron_sword","count":1,"description":"购买一把铁剑","completed":false}]`,
			Rewards:     `{"exp":200,"gold":300,"items":[]}`,
			NextTask:    "",
		},
	}

	for i := range tasks {
		if err := db.Create(&tasks[i]).Error; err != nil {
			return err
		}
	}

	// 创建NPC购物流程
	flows := []models.Flow{
		{
			Name:        "NPC出门购物流程",
			Code:        "flow_npc_shopping",
			Description: "NPC从家里出发，前往商店购买物品的完整流程",
			Nodes:       `[{"id":"node_1","type":"start","position":{"x":100,"y":200},"data":{"label":"开始"}},{"id":"node_2","type":"action","position":{"x":250,"y":200},"data":{"label":"NPC离开家","action":"move","params":{"npcId":"npc_merchant_li","from":"home","to":"scene_town_center"}}},{"id":"node_3","type":"action","position":{"x":400,"y":200},"data":{"label":"前往商店","action":"move","params":{"npcId":"npc_merchant_li","from":"scene_town_center","to":"scene_general_store"}}},{"id":"node_4","type":"condition","position":{"x":550,"y":200},"data":{"label":"商店开门?","condition":"shop.isOpen(npc.shopId)"}},{"id":"node_5","type":"action","position":{"x":700,"y":100},"data":{"label":"等待商店开门","action":"wait","params":{"duration":60}}},{"id":"node_6","type":"action","position":{"x":700,"y":300},"data":{"label":"进入商店","action":"enter","params":{"npcId":"npc_merchant_li","shopId":"shop_general_store"}}},{"id":"node_7","type":"action","position":{"x":850,"y":200},"data":{"label":"查看商品","action":"browse","params":{"duration":30}}},{"id":"node_8","type":"action","position":{"x":1000,"y":200},"data":{"label":"购买物品","action":"purchase","params":{"items":[{"itemId":"item_herb","count":5},{"itemId":"item_mantou","count":10}]}}},{"id":"node_9","type":"action","position":{"x":1150,"y":200},"data":{"label":"离开商店","action":"leave"}},{"id":"node_10","type":"action","position":{"x":1300,"y":200},"data":{"label":"返回家中","action":"move","params":{"to":"home"}}},{"id":"node_11","type":"end","position":{"x":1450,"y":200},"data":{"label":"结束"}}]`,
			Edges:       `[{"id":"e1-2","source":"node_1","target":"node_2"},{"id":"e2-3","source":"node_2","target":"node_3"},{"id":"e3-4","source":"node_3","target":"node_4"},{"id":"e4-5","source":"node_4","target":"node_5","label":"否"},{"id":"e4-6","source":"node_4","target":"node_6","label":"是"},{"id":"e5-6","source":"node_5","target":"node_6"},{"id":"e6-7","source":"node_6","target":"node_7"},{"id":"e7-8","source":"node_7","target":"node_8"},{"id":"e8-9","source":"node_8","target":"node_9"},{"id":"e9-10","source":"node_9","target":"node_10"},{"id":"e10-11","source":"node_10","target":"node_11"}]`,
		},
	}

	for i := range flows {
		if err := db.Create(&flows[i]).Error; err != nil {
			return err
		}
	}

	log.Println("Database seeded successfully!")
	return nil
}
