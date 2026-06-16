package database

import (
	"log"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	log.Println("Seeding database...")

	var count int64
	db.Model(&models.Scene{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// ===== 新手村场景 =====
	scenes := []models.Scene{
		{
			Name:        "村口",
			Code:        "scene_village_entrance",
			Description: "青石村的入口，一块刻着「青石村」三个大字的石碑立在路边。远处可见炊烟袅袅，偶尔传来几声鸡鸣。",
			Background:  "#4a7c59",
			Width:       1600,
			Height:      900,
		},
		{
			Name:        "村中心广场",
			Code:        "scene_village_center",
			Description: "村子的中心地带，一棵大榕树下是村民们聚集闲聊的地方。广场中央有一口老井。",
			Background:  "#8b7355",
			Width:       1600,
			Height:      900,
		},
		{
			Name:        "李记杂货铺",
			Code:        "scene_general_store",
			Description: "门口挂着「李记杂货」的招牌，店里摆满了各种日用品和药材。老板李掌柜正在柜台后算账。",
			Background:  "#6b4226",
			Width:       1200,
			Height:      800,
		},
		{
			Name:        "张记铁匠铺",
			Code:        "scene_blacksmith",
			Description: "炉火通红，锤声叮当。铁匠铺里摆满了各式兵器和农具，空气中弥漫着铁锈和炭火的气息。",
			Background:  "#4a4a4a",
			Width:       1200,
			Height:      800,
		},
		{
			Name:        "王大娘茶摊",
			Code:        "scene_tea_stand",
			Description: "路边的一个小茶摊，几把竹椅一张木桌，王大娘热情地招呼着过往的行人。",
			Background:  "#6b8e23",
			Width:       1000,
			Height:      700,
		},
		{
			Name:        "村外小路",
			Code:        "scene_village_path",
			Description: "通往外界的小路，两旁是茂密的树林。偶尔能看到野兔窜过，远处似乎有狼嚎声。",
			Background:  "#556b2f",
			Width:       2000,
			Height:      900,
		},
	}

	for i := range scenes {
		if err := db.Create(&scenes[i]).Error; err != nil {
			return err
		}
	}

	// ===== NPC =====
	npcs := []models.NPC{
		{
			Name:        "老村长",
			Code:        "npc_chief_chen",
			Title:       "青石村村长",
			Description: "一位须发皆白的老人，目光慈祥而深邃。他在青石村生活了一辈子，对村子的一切了如指掌。",
			Avatar:      "👴",
			Sprite:      "chief",
			Behaviors:   `["idle","greet","guide"]`,
			Schedule:    `[{"time":"06:00","action":"stand_at_tree","scene":"scene_village_center"},{"time":"22:00","action":"go_home"}]`,
		},
		{
			Name:        "李掌柜",
			Code:        "npc_merchant_li",
			Title:       "杂货铺老板",
			Description: "一位精明的中年商人，笑起来眼睛眯成一条缝。他的杂货铺是村里最热闹的地方之一。",
			Avatar:      "🧔",
			Sprite:      "merchant",
			Behaviors:   `["idle","greet","sell"]`,
			Schedule:    `[{"time":"06:00","action":"open_shop","scene":"scene_general_store"},{"time":"22:00","action":"close_shop","scene":"scene_village_center"}]`,
		},
		{
			Name:        "王大娘",
			Code:        "npc_tea_wang",
			Title:       "茶摊老板娘",
			Description: "热情开朗的中年妇人，她的茶远近闻名。她知道村里所有的八卦和故事。",
			Avatar:      "👩",
			Sprite:      "tealady",
			Behaviors:   `["idle","chat","serve_tea"]`,
			Schedule:    `[]`,
		},
		{
			Name:        "张铁匠",
			Code:        "npc_blacksmith_zhang",
			Title:       "铁匠铺师傅",
			Description: "膀大腰圆的壮汉，手臂上满是烫伤的疤痕。他打造的兵器在方圆百里都小有名气。",
			Avatar:      "👨‍🔧",
			Sprite:      "blacksmith",
			Behaviors:   `["idle","forge","sell"]`,
			Schedule:    `[{"time":"08:00","action":"open_shop","scene":"scene_blacksmith"},{"time":"20:00","action":"close_shop","scene":"scene_village_center"}]`,
		},
		{
			Name:        "猎户老周",
			Code:        "npc_hunter_zhou",
			Title:       "猎人",
			Description: "沉默寡言的中年猎人，腰间别着一把猎弓。他经常在村外打猎，对野外的危险了如指掌。",
			Avatar:      "🏹",
			Sprite:      "hunter",
			Behaviors:   `["idle","patrol","hunt"]`,
			Schedule:    `[{"time":"05:00","action":"go_hunt","scene":"scene_village_path"},{"time":"18:00","action":"return_village","scene":"scene_village_center"}]`,
		},
		{
			Name:        "小石头",
			Code:        "npc_kid_stone",
			Title:       "村里的孩子",
			Description: "一个调皮的小男孩，整天在村里跑来跑去。他对外面的世界充满好奇。",
			Avatar:      "👦",
			Sprite:      "kid",
			Behaviors:   `["idle","run","play"]`,
			Schedule:    `[]`,
		},
	}

	for i := range npcs {
		if err := db.Create(&npcs[i]).Error; err != nil {
			return err
		}
	}

	// ===== 智能体 =====
	agents := []models.Agent{
		{
			Name:           "老村长智能体",
			Code:           "agent_chief_chen",
			Description:    "青石村村长的AI智能体，负责引导新手冒险者",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.7,
			MaxTokens:      500,
			SystemPrompt:   `你是陈老村长，青石村最受尊敬的长者。你慈祥和蔼，说话慢条斯理，总是带着"呵呵"的笑声。你对村子的历史了如指掌，对每一位村民都很熟悉。

当遇到新来的冒险者时，你会热情地欢迎他们，介绍村子的情况，并指引他们去认识村里的其他人。

村子背景：青石村是一个偏远的小村庄，以盛产青石而得名。村子虽然不大，但五脏俱全，有杂货铺、铁匠铺、茶摊等。最近村子外面出现了不少野兽，村民们有些担忧。

对话风格：
- 使用"老朽"自称
- 称呼对方为"年轻人"或"小友"
- 语气慈祥温和
- 偶尔讲讲村子的故事`,
			MemoryType:     "sliding_window",
			MaxMessages:    20,
			SummaryEnabled: true,
			KnowledgeBase:  `[{"id":"kb_village","name":"青石村历史","type":"text"},{"id":"kb_npcs","name":"村民信息","type":"text"}]`,
			Tools:          `[{"name":"give_quest","description":"发布任务"},{"name":"give_reward","description":"发放奖励"}]`,
		},
		{
			Name:           "李掌柜智能体",
			Code:           "agent_merchant_li",
			Description:    "杂货铺老板的AI智能体，负责与玩家进行买卖对话",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.7,
			MaxTokens:      500,
			SystemPrompt:   `你是李掌柜，青石村杂货铺的老板。你精明但诚实，总是笑脸迎人。你经营各种日用品、药材和食材。

商品列表：
- 草药：50文，可恢复20点生命
- 灵芝：200文，可恢复100点生命
- 馒头：20文，恢复体力
- 烧酒：80文，驱寒保暖
- 麻绳：30文，攀爬工具
- 回城符：100文，瞬间回村

对话风格：
- 使用"在下"自称
- 称呼对方为"客官"
- 热情推荐商品
- 偶尔给老顾客打折`,
			MemoryType:     "sliding_window",
			MaxMessages:    20,
			SummaryEnabled: true,
			KnowledgeBase:  `[{"id":"kb_items","name":"商品目录","type":"text"},{"id":"kb_prices","name":"价格表","type":"table"}]`,
			Tools:          `[{"name":"query_inventory","description":"查询库存"},{"name":"make_deal","description":"进行交易"}]`,
		},
		{
			Name:           "王大娘智能体",
			Code:           "agent_tea_wang",
			Description:    "茶摊老板娘的AI智能体，喜欢聊天八卦",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.8,
			MaxTokens:      400,
			SystemPrompt:   `你是王大娘，青石村茶摊的老板娘。你热情开朗，最喜欢聊天。你知道村里所有的八卦和故事。

你知道的事情：
1. 村外最近出现了狼群，猎户老周正在发愁
2. 铁匠张师傅最近在打造一批新兵器
3. 老村长说最近会有大人物来村里
4. 村东头的古井据说有几百年的历史
5. 传说村子后面的山上有个山洞，里面藏着宝贝

对话风格：
- 使用"哎呀"、"我跟你说"等口头禅
- 称呼对方为"客官"或"孩子"
- 爱分享八卦和故事
- 说话比较啰嗦但很亲切`,
			MemoryType:     "sliding_window",
			MaxMessages:    15,
			SummaryEnabled: false,
			KnowledgeBase:  `[{"id":"kb_gossip","name":"村中八卦","type":"text"},{"id":"kb_history","name":"传说故事","type":"text"}]`,
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
			SystemPrompt:   `你是张铁匠，青石村铁匠铺的师傅。你性格沉稳，说话直接，做事认真。你打造的兵器远近闻名。

商品列表：
- 铁剑：500文，攻击+10
- 精钢刀：1000文，攻击+25
- 铁甲：800文，防御+15
- 皮甲：300文，防御+8
- 铁盾：600文，防御+20
- 猎弓：400文，远程攻击

对话风格：
- 使用"俺"自称
- 说话朴实直接
- 对兵器很自豪
- 不喜欢拐弯抹角`,
			MemoryType:     "sliding_window",
			MaxMessages:    20,
			SummaryEnabled: true,
			KnowledgeBase:  `[{"id":"kb_weapons","name":"兵器谱","type":"text"},{"id":"kb_materials","name":"锻造材料","type":"text"}]`,
			Tools:          `[{"name":"forge_weapon","description":"打造兵器"},{"name":"repair_item","description":"修理装备"}]`,
		},
		{
			Name:           "猎户老周智能体",
			Code:           "agent_hunter_zhou",
			Description:    "猎人的AI智能体，熟悉野外生存",
			LLMProvider:    "openai",
			LLMModel:       "gpt-4",
			Temperature:    0.6,
			MaxTokens:      400,
			SystemPrompt:   `你是猎户老周，青石村的猎人。你沉默寡言但心地善良，对野外的危险了如指掌。

你最近很担心：村外的狼群越来越多，而且有几只特别大的头狼。你一个人应付不来，希望有冒险者能帮忙。

你可以教冒险者：
1. 如何设置陷阱捕捉野兔
2. 如何识别野兽的踪迹
3. 如何在野外生火
4. 哪些草药可以治伤

对话风格：
- 话不多，但每句都很实在
- 使用"嗯"、"哦"等简短回应
- 谈到打猎时会变得健谈
- 对勇敢的人比较友善`,
			MemoryType:     "sliding_window",
			MaxMessages:    15,
			SummaryEnabled: true,
			KnowledgeBase:  `[{"id":"kb_hunting","name":"狩猎技巧","type":"text"},{"id":"kb_danger","name":"野外危险","type":"text"}]`,
			Tools:          `[{"name":"set_trap","description":"设置陷阱"},{"name":"track_beast","description":"追踪野兽"}]`,
		},
	}

	for i := range agents {
		if err := db.Create(&agents[i]).Error; err != nil {
			return err
		}
	}

	// 关联NPC和智能体
	db.Model(&npcs[0]).Update("agent_id", agents[0].ID)
	db.Model(&npcs[1]).Update("agent_id", agents[1].ID)
	db.Model(&npcs[2]).Update("agent_id", agents[2].ID)
	db.Model(&npcs[3]).Update("agent_id", agents[3].ID)
	db.Model(&npcs[4]).Update("agent_id", agents[4].ID)

	// ===== 道具 =====
	items := []models.Item{
		{Name: "草药", Code: "item_herb", Category: "medicine", Description: "普通的草药，可恢复少量生命", Effect: `{"hp":20}`},
		{Name: "灵芝", Code: "item_lingzhi", Category: "medicine", Description: "珍贵的药材，可恢复大量生命", Effect: `{"hp":100}`},
		{Name: "金疮药", Code: "item_heal_paste", Category: "medicine", Description: "特效伤药，快速止血愈伤", Effect: `{"hp":50}`},
		{Name: "馒头", Code: "item_mantou", Category: "food", Description: "热腾腾的馒头，恢复体力", Effect: `{"stamina":30}`},
		{Name: "烧酒", Code: "item_wine", Category: "food", Description: "烈性白酒，驱寒保暖", Effect: `{"cold_resist":30}`},
		{Name: "麻绳", Code: "item_rope", Category: "tool", Description: "结实的麻绳，可用于攀爬", Effect: `{}`},
		{Name: "回城符", Code: "item_town_portal", Category: "tool", Description: "神奇的符纸，使用后可瞬间回到村口", Effect: `{"teleport":"scene_village_entrance"}`},
		{Name: "铁剑", Code: "item_iron_sword", Category: "weapon", Description: "普通的铁制长剑，适合新手使用", Effect: `{"attack":10}`},
		{Name: "精钢刀", Code: "item_steel_blade", Category: "weapon", Description: "锻造精良的钢刀，锋利无比", Effect: `{"attack":25}`},
		{Name: "猎弓", Code: "item_hunting_bow", Category: "weapon", Description: "猎人常用的弓，适合远程攻击", Effect: `{"attack":15,"range":200}`},
		{Name: "铁甲", Code: "item_iron_armor", Category: "armor", Description: "铁制铠甲，提供不错的防护", Effect: `{"defense":15}`},
		{Name: "皮甲", Code: "item_leather_armor", Category: "armor", Description: "轻便的皮甲，不影响行动", Effect: `{"defense":8}`},
		{Name: "铁盾", Code: "item_iron_shield", Category: "armor", Description: "坚固的铁盾，可格挡攻击", Effect: `{"defense":20}`},
		{Name: "新手木剑", Code: "item_wooden_sword", Category: "weapon", Description: "练习用的木剑，聊胜于无", Effect: `{"attack":3}`},
		{Name: "粗布衣", Code: "item_cloth_armor", Category: "armor", Description: "普通的粗布衣服，稍微有点防护", Effect: `{"defense":2}`},
		{Name: "兔肉", Code: "item_rabbit_meat", Category: "food", Description: "新鲜的兔肉，烤熟后很美味", Effect: `{"hp":15,"stamina":20}`},
		{Name: "狼皮", Code: "item_wolf_pelt", Category: "material", Description: "狼的皮毛，可以卖个好价钱", Effect: `{}`},
		{Name: "草药种子", Code: "item_herb_seed", Category: "material", Description: "可以种出草药的种子", Effect: `{}`},
	}

	for i := range items {
		if err := db.Create(&items[i]).Error; err != nil {
			return err
		}
	}

	// ===== 商店 =====
	shops := []models.Shop{
		{
			Name:        "李记杂货铺",
			Code:        "shop_general_store",
			Type:        "general",
			Description: "售卖各种日用品、药材和食材",
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
			Description: "打造和售卖各种兵器、防具",
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

	// 商店商品
	shopItems := []models.ShopItem{
		// 杂货铺
		{ShopID: shops[0].ID, ItemID: items[0].ID, Price: 50, Stock: 100},   // 草药
		{ShopID: shops[0].ID, ItemID: items[1].ID, Price: 200, Stock: 30},    // 灵芝
		{ShopID: shops[0].ID, ItemID: items[2].ID, Price: 100, Stock: 50},    // 金疮药
		{ShopID: shops[0].ID, ItemID: items[3].ID, Price: 20, Stock: 200},    // 馒头
		{ShopID: shops[0].ID, ItemID: items[4].ID, Price: 80, Stock: 60},     // 烧酒
		{ShopID: shops[0].ID, ItemID: items[5].ID, Price: 30, Stock: 40},     // 麻绳
		{ShopID: shops[0].ID, ItemID: items[6].ID, Price: 100, Stock: 20},    // 回城符
		// 铁匠铺
		{ShopID: shops[1].ID, ItemID: items[7].ID, Price: 500, Stock: 15},    // 铁剑
		{ShopID: shops[1].ID, ItemID: items[8].ID, Price: 1000, Stock: 5},    // 精钢刀
		{ShopID: shops[1].ID, ItemID: items[9].ID, Price: 400, Stock: 10},    // 猎弓
		{ShopID: shops[1].ID, ItemID: items[10].ID, Price: 800, Stock: 8},    // 铁甲
		{ShopID: shops[1].ID, ItemID: items[11].ID, Price: 300, Stock: 20},   // 皮甲
		{ShopID: shops[1].ID, ItemID: items[12].ID, Price: 600, Stock: 10},   // 铁盾
	}

	for i := range shopItems {
		if err := db.Create(&shopItems[i]).Error; err != nil {
			return err
		}
	}

	// 更新NPC关联的商店
	db.Model(&npcs[1]).Update("shop_id", shops[0].ID)
	db.Model(&npcs[3]).Update("shop_id", shops[1].ID)

	// ===== 场景NPC放置 =====
	sceneNPCs := []models.SceneNPC{
		// 村口
		{SceneID: scenes[0].ID, NPCID: npcs[0].ID, X: 400, Y: 400},  // 老村长在村口迎接
		{SceneID: scenes[0].ID, NPCID: npcs[5].ID, X: 800, Y: 500},  // 小石头在村口玩
		// 村中心
		{SceneID: scenes[1].ID, NPCID: npcs[0].ID, X: 800, Y: 300},  // 老村长在榕树下
		{SceneID: scenes[1].ID, NPCID: npcs[2].ID, X: 400, Y: 600},  // 王大娘在茶摊
		{SceneID: scenes[1].ID, NPCID: npcs[5].ID, X: 1000, Y: 700}, // 小石头在玩耍
		// 杂货铺
		{SceneID: scenes[2].ID, NPCID: npcs[1].ID, X: 600, Y: 400},  // 李掌柜在柜台
		// 铁匠铺
		{SceneID: scenes[3].ID, NPCID: npcs[3].ID, X: 600, Y: 400},  // 张铁匠在打铁
		// 茶摊
		{SceneID: scenes[4].ID, NPCID: npcs[2].ID, X: 500, Y: 350},  // 王大娘在泡茶
		// 村外小路
		{SceneID: scenes[5].ID, NPCID: npcs[4].ID, X: 300, Y: 450},  // 猎户老周在巡逻
	}

	for i := range sceneNPCs {
		if err := db.Create(&sceneNPCs[i]).Error; err != nil {
			return err
		}
	}

	// ===== 传送点 =====
	portals := []models.Portal{
		// 村口 -> 村中心
		{SceneID: scenes[0].ID, X: 1400, Y: 450, TargetScene: "scene_village_center", TargetX: 100, TargetY: 450},
		// 村中心 -> 村口
		{SceneID: scenes[1].ID, X: 50, Y: 450, TargetScene: "scene_village_entrance", TargetX: 1300, TargetY: 450},
		// 村中心 -> 杂货铺
		{SceneID: scenes[1].ID, X: 600, Y: 100, TargetScene: "scene_general_store", TargetX: 600, TargetY: 700},
		// 村中心 -> 铁匠铺
		{SceneID: scenes[1].ID, X: 1100, Y: 100, TargetScene: "scene_blacksmith", TargetX: 600, TargetY: 700},
		// 村中心 -> 茶摊
		{SceneID: scenes[1].ID, X: 300, Y: 700, TargetScene: "scene_tea_stand", TargetX: 500, TargetY: 600},
		// 村中心 -> 村外小路
		{SceneID: scenes[1].ID, X: 1500, Y: 450, TargetScene: "scene_village_path", TargetX: 100, TargetY: 450},
		// 杂货铺 -> 村中心
		{SceneID: scenes[2].ID, X: 600, Y: 750, TargetScene: "scene_village_center", TargetX: 600, TargetY: 200},
		// 铁匠铺 -> 村中心
		{SceneID: scenes[3].ID, X: 600, Y: 750, TargetScene: "scene_village_center", TargetX: 1100, TargetY: 200},
		// 茶摊 -> 村中心
		{SceneID: scenes[4].ID, X: 500, Y: 650, TargetScene: "scene_village_center", TargetX: 300, TargetY: 600},
		// 村外小路 -> 村中心
		{SceneID: scenes[5].ID, X: 50, Y: 450, TargetScene: "scene_village_center", TargetX: 1400, TargetY: 450},
	}

	for i := range portals {
		if err := db.Create(&portals[i]).Error; err != nil {
			return err
		}
	}

	// ===== 提示词模板 =====
	templates := []models.PromptTemplate{
		{
			Name:     "NPC基础人设",
			Code:     "template_npc_basic",
			Category: "system",
			Content:  `你是{{npc_name}}，{{npc_title}}。{{npc_description}}当前状态：位置：{{current_scene}}，时间：{{current_time}}，心情：{{mood}}。请保持角色设定，用古风语气与玩家对话。`,
			Variables: `[{"name":"npc_name","type":"string","description":"NPC名称"},{"name":"npc_title","type":"string","description":"NPC称号"},{"name":"npc_description","type":"string","description":"NPC描述"},{"name":"current_scene","type":"string","description":"当前场景"},{"name":"current_time","type":"string","description":"当前时间"},{"name":"mood","type":"string","description":"NPC心情"}]`,
		},
	}

	for i := range templates {
		if err := db.Create(&templates[i]).Error; err != nil {
			return err
		}
	}

	// ===== 任务 - 新手村主线 =====
	tasks := []models.Task{
		{
			Name:        "初来乍到",
			Code:        "task_first_arrival",
			Type:        "main",
			Description: "你来到了青石村，先和村长打个招呼吧。也许他能告诉你一些关于这个地方的事情。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives:  `[{"id":"obj_001","type":"dialogue","target":"npc_chief_chen","description":"与老村长对话","completed":false}]`,
			Rewards:     `{"exp":50,"gold":100,"items":[{"code":"item_wooden_sword","count":1},{"code":"item_cloth_armor","count":1}]}`,
			NextTask:    "task_explore_village",
		},
		{
			Name:        "熟悉村庄",
			Code:        "task_explore_village",
			Type:        "main",
			Description: "村长让你去认识村里的人，逛逛各个地方。去杂货铺、铁匠铺和茶摊看看吧。",
			Status:      "locked",
			Trigger:     `{"type":"task_complete","conditions":[{"type":"task_id","value":"task_first_arrival"}]}`,
			Objectives:  `[{"id":"obj_002","type":"visit","target":"scene_general_store","description":"拜访杂货铺","completed":false},{"id":"obj_003","type":"visit","target":"scene_blacksmith","description":"拜访铁匠铺","completed":false},{"id":"obj_004","type":"visit","target":"scene_tea_stand","description":"去茶摊喝杯茶","completed":false}]`,
			Rewards:     `{"exp":100,"gold":200,"items":[{"code":"item_herb","count":5},{"code":"item_mantou","count":3}]}`,
			NextTask:    "task_stock_up",
		},
		{
			Name:        "采购补给",
			Code:        "task_stock_up",
			Type:        "main",
			Description: "冒险需要准备，去杂货铺买些草药和食物吧。记住，有备无患！",
			Status:      "locked",
			Trigger:     `{"type":"task_complete","conditions":[{"type":"task_id","value":"task_explore_village"}]}`,
			Objectives:  `[{"id":"obj_005","type":"collect","target":"item_herb","count":3,"description":"购买3份草药","completed":false},{"id":"obj_006","type":"collect","target":"item_mantou","count":5,"description":"购买5个馒头","completed":false}]`,
			Rewards:     `{"exp":80,"gold":150,"items":[]}`,
			NextTask:    "task_gear_up",
		},
		{
			Name:        "装备升级",
			Code:        "task_gear_up",
			Type:        "main",
			Description: "有了补给，接下来去铁匠铺买把像样的武器吧。村外可不太平。",
			Status:      "locked",
			Trigger:     `{"type":"task_complete","conditions":[{"type":"task_id","value":"task_stock_up"}]}`,
			Objectives:  `[{"id":"obj_007","type":"dialogue","target":"npc_blacksmith_zhang","description":"与张铁匠对话","completed":false},{"id":"obj_008","type":"collect","target":"item_iron_sword","count":1,"description":"购买一把铁剑","completed":false}]`,
			Rewards:     `{"exp":120,"gold":300,"items":[{"code":"item_heal_paste","count":3}]}`,
			NextTask:    "task_first_test",
		},
		{
			Name:        "初试身手",
			Code:        "task_first_test",
			Type:        "main",
			Description: "猎户老周说村外有狼群出没，去村外小路帮他查看一下情况吧。",
			Status:      "locked",
			Trigger:     `{"type":"task_complete","conditions":[{"type":"task_id","value":"task_gear_up"}]}`,
			Objectives:  `[{"id":"obj_009","type":"dialogue","target":"npc_hunter_zhou","description":"找到猎户老周","completed":false},{"id":"obj_010","type":"visit","target":"scene_village_path","description":"前往村外小路","completed":false}]`,
			Rewards:     `{"exp":200,"gold":500,"items":[{"code":"item_town_portal","count":3}]}`,
			NextTask:    "",
		},
		// 支线任务
		{
			Name:        "王大娘的茶",
			Code:        "task_tea_time",
			Type:        "side",
			Description: "王大娘想请你喝杯茶，顺便聊聊村里的趣事。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives:  `[{"id":"obj_s01","type":"dialogue","target":"npc_tea_wang","description":"和王大娘聊天","completed":false}]`,
			Rewards:     `{"exp":30,"gold":0,"items":[{"code":"item_herb","count":2}]}`,
			NextTask:    "",
		},
		{
			Name:        "调皮的小石头",
			Code:        "task_kid_stone",
			Type:        "side",
			Description: "小石头在村口玩，他说想看看你的武器。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives:  `[{"id":"obj_s02","type":"dialogue","target":"npc_kid_stone","description":"和小石头说话","completed":false}]`,
			Rewards:     `{"exp":20,"gold":50,"items":[]}`,
			NextTask:    "",
		},
	}

	for i := range tasks {
		if err := db.Create(&tasks[i]).Error; err != nil {
			return err
		}
	}

	// ===== 新手教程流程 =====
	flows := []models.Flow{
		{
			Name:        "新手引导流程",
			Code:        "flow_tutorial",
			Description: "新玩家进入游戏后的引导教程",
			Nodes: `[{"id":"node_1","type":"start","position":{"x":100,"y":200},"data":{"label":"开始","text":"欢迎来到青石村！"}},{"id":"node_2","type":"tutorial","position":{"x":300,"y":200},"data":{"label":"移动教程","text":"使用 WASD 或方向键移动角色","action":"show_move_tutorial"}},{"id":"node_3","type":"tutorial","position":{"x":500,"y":200},"data":{"label":"对话教程","text":"靠近NPC时，点击NPC可以与其对话","action":"show_talk_tutorial"}},{"id":"node_4","type":"tutorial","position":{"x":700,"y":200},"data":{"label":"任务教程","text":"完成NPC交给你的任务可以获得奖励","action":"show_quest_tutorial"}},{"id":"node_5","type":"tutorial","position":{"x":900,"y":200},"data":{"label":"商店教程","text":"在商店可以购买道具和装备","action":"show_shop_tutorial"}},{"id":"node_6","type":"action","position":{"x":1100,"y":200},"data":{"label":"第一个任务","text":"去找老村长聊聊吧！","action":"activate_task","params":{"task":"task_first_arrival"}}},{"id":"node_7","type":"end","position":{"x":1300,"y":200},"data":{"label":"教程完成","text":"祝你冒险愉快！"}}]`,
			Edges: `[{"id":"e1-2","source":"node_1","target":"node_2"},{"id":"e2-3","source":"node_2","target":"node_3"},{"id":"e3-4","source":"node_3","target":"node_4"},{"id":"e4-5","source":"node_4","target":"node_5"},{"id":"e5-6","source":"node_5","target":"node_6"},{"id":"e6-7","source":"node_6","target":"node_7"}]`,
		},
	}

	for i := range flows {
		if err := db.Create(&flows[i]).Error; err != nil {
			return err
		}
	}

	// ===== 游戏配置 =====
	gameConfigs := []models.GameConfig{
		{Key: "game_name", Value: "青石村传说"},
		{Key: "game_version", Value: "0.1.0"},
		{Key: "start_scene", Value: "scene_village_entrance"},
		{Key: "start_x", Value: "200"},
		{Key: "start_y", Value: "450"},
		{Key: "player_speed", Value: "200"},
		{Key: "default_hp", Value: "100"},
		{Key: "default_mp", Value: "50"},
		{Key: "default_gold", Value: "500"},
	}

	for i := range gameConfigs {
		if err := db.Create(&gameConfigs[i]).Error; err != nil {
			return err
		}
	}

	// ===== 技能 =====
	skills := []models.Skill{
		{
			Name:        "基础斩击",
			Code:        "skill_basic_slash",
			Description: "用武器进行强力斩击，造成1.5倍武器伤害",
			Type:        "attack",
			MPCost:      5,
			Damage:      15,
			Heal:        0,
			Cooldown:    0,
			Level:       1,
			Effect:      `{}`,
		},
		{
			Name:        "治疗术",
			Code:        "skill_heal",
			Description: "使用法力恢复生命值",
			Type:        "heal",
			MPCost:      8,
			Damage:      0,
			Heal:        30,
			Cooldown:    1,
			Level:       3,
			Effect:      `{}`,
		},
		{
			Name:        "火球术",
			Code:        "skill_fireball",
			Description: "发射火球造成魔法伤害",
			Type:        "attack",
			MPCost:      15,
			Damage:      25,
			Heal:        0,
			Cooldown:    1,
			Level:       5,
			Effect:      `{"type":"burn","duration":2}`,
		},
		{
			Name:        "铁壁",
			Code:        "skill_iron_wall",
			Description: "提升防御力3回合",
			Type:        "buff",
			MPCost:      10,
			Damage:      0,
			Heal:        0,
			Cooldown:    2,
			Level:       4,
			Effect:      `{"type":"defense_up","duration":3,"value":50}`,
		},
		{
			Name:        "疾风步",
			Code:        "skill_wind_step",
			Description: "提升速度3回合",
			Type:        "buff",
			MPCost:      12,
			Damage:      0,
			Heal:        0,
			Cooldown:    2,
			Level:       6,
			Effect:      `{"type":"speed_up","duration":3,"value":50}`,
		},
		{
			Name:        "毒击",
			Code:        "skill_poison_strike",
			Description: "附带毒素的攻击，使敌人持续掉血",
			Type:        "debuff",
			MPCost:      10,
			Damage:      10,
			Heal:        0,
			Cooldown:    2,
			Level:       3,
			Effect:      `{"type":"poison","duration":3,"value":5}`,
		},
		{
			Name:        "战吼",
			Code:        "skill_warcry",
			Description: "发出战吼，提升攻击力2回合",
			Type:        "buff",
			MPCost:      8,
			Damage:      0,
			Heal:        0,
			Cooldown:    2,
			Level:       2,
			Effect:      `{"type":"attack_up","duration":2,"value":30}`,
		},
	}

	for i := range skills {
		if err := db.Create(&skills[i]).Error; err != nil {
			return err
		}
	}

	// ===== 成就 =====
	achievements := []models.Achievement{
		{
			Name:        "初来乍到",
			Code:        "ach_first_quest",
			Description: "完成第一个任务",
			Condition:   `{"type":"quest_complete","value":1}`,
			Reward:      `{"exp":50,"gold":100}`,
			Icon:        "⭐",
		},
		{
			Name:        "村庄之友",
			Code:        "ach_talk_all_npcs",
			Description: "与所有NPC对话过",
			Condition:   `{"type":"talk_all_npcs","value":1}`,
			Reward:      `{"exp":200,"gold":300}`,
			Icon:        "👥",
		},
		{
			Name:        "富甲一方",
			Code:        "ach_rich",
			Description: "累计获得10000金币",
			Condition:   `{"type":"gold","value":10000}`,
			Reward:      `{"exp":500,"gold":1000}`,
			Icon:        "💰",
		},
		{
			Name:        "百战百胜",
			Code:        "ach_combat_100",
			Description: "赢得100场战斗",
			Condition:   `{"type":"combat_win","value":100}`,
			Reward:      `{"exp":1000,"gold":2000}`,
			Icon:        "⚔️",
		},
		{
			Name:        "探索者",
			Code:        "ach_explorer",
			Description: "探索所有场景",
			Condition:   `{"type":"explore","value":6}`,
			Reward:      `{"exp":300,"gold":500}`,
			Icon:        "🗺️",
		},
		{
			Name:        "收藏家",
			Code:        "ach_collector",
			Description: "拥有50种不同的道具",
			Condition:   `{"type":"collect","value":50}`,
			Reward:      `{"exp":800,"gold":1500}`,
			Icon:        "🎒",
		},
		{
			Name:        "初试牛刀",
			Code:        "ach_first_combat",
			Description: "赢得第一场战斗",
			Condition:   `{"type":"combat_win","value":1}`,
			Reward:      `{"exp":30,"gold":50}`,
			Icon:        "🗡️",
		},
		{
			Name:        "技能大师",
			Code:        "ach_skill_master",
			Description: "使用100次技能",
			Condition:   `{"type":"skill_use","value":100}`,
			Reward:      `{"exp":600,"gold":800}`,
			Icon:        "✨",
		},
		{
			Name:        "等级10",
			Code:        "ach_level_10",
			Description: "达到10级",
			Condition:   `{"type":"level","value":10}`,
			Reward:      `{"exp":200,"gold":500}`,
			Icon:        "🏆",
		},
		{
			Name:        "等级20",
			Code:        "ach_level_20",
			Description: "达到20级",
			Condition:   `{"type":"level","value":20}`,
			Reward:      `{"exp":500,"gold":1000}`,
			Icon:        "👑",
		},
	}

	for i := range achievements {
		if err := db.Create(&achievements[i]).Error; err != nil {
			return err
		}
	}

	log.Println("Database seeded successfully!")
	return nil
}
