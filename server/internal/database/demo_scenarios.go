package database

import (
	"log"

	"github.com/ddc-111/agentGame/server/internal/database/models"
	"gorm.io/gorm"
)

// SeedDemoScenarios 添加AI对话演示场景
func SeedDemoScenarios(db *gorm.DB) error {
	log.Println("Seeding demo scenarios...")

	// 检查是否已有演示数据
	var count int64
	db.Model(&models.Task{}).Where("code LIKE ?", "demo_%").Count(&count)
	if count > 0 {
		log.Println("Demo scenarios already seeded, skipping...")
		return nil
	}

	// ===== 场景1: 李掌柜的推荐 =====
	// 当玩家第一次和商人对话时，他会询问需要什么
	// AI根据玩家等级/职业提供个性化推荐
	// 记住玩家之前是否购买过东西
	merchantDemoTasks := []models.Task{
		{
			Name:        "李掌柜的推荐",
			Code:        "demo_merchant_recommend",
			Type:        "side",
			Description: "李掌柜想根据你的需求推荐商品。和他聊聊你需要什么装备。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives: `[{"id":"obj_demo_01","type":"dialogue","target":"npc_merchant_li","description":"与李掌柜对话，告诉他你的需求","completed":false},
{"id":"obj_demo_02","type":"dialogue","target":"npc_merchant_li","description":"听取李掌柜的推荐","completed":false}]`,
			Rewards:  `{"exp":50,"gold":0,"items":[{"code":"item_herb","count":3}]}`,
			NextTask: "",
		},
	}

	for i := range merchantDemoTasks {
		if err := db.Create(&merchantDemoTasks[i]).Error; err != nil {
			return err
		}
	}

	// 李掌柜智能体增强系统提示词（在seed.go中已创建，这里通过更新Agent来增强）
	var merchantAgent models.Agent
	if err := db.Where("code = ?", "agent_merchant_li").First(&merchantAgent).Error; err == nil {
		enhancedPrompt := merchantAgent.SystemPrompt + `

【AI演示功能】
- 当玩家第一次来时，主动询问他们需要什么类型的装备（武器/防具/药品/食物）
- 根据玩家等级推荐合适的商品：
  - 等级1-5：推荐基础装备（草药、馒头、新手木剑）
  - 等级6-10：推荐中级装备（金疮药、铁剑、皮甲）
  - 等级10+：推荐高级装备（灵芝、精钢刀、铁甲）
- 如果玩家之前买过东西，说"客官又来了！上次买的用着还顺手吗？"
- 如果玩家犹豫不决，可以给个小折扣（"看客官这么有诚意，给打个九折吧"）
- 用古风语气，但要让玩家感受到AI的智能`
		db.Model(&merchantAgent).Update("system_prompt", enhancedPrompt)
	}

	// ===== 场景2: 王大娘的故事 =====
	// 茶摊老板娘根据时间讲不同的故事
	// 能给出隐藏任务的提示
	// 对玩家的任务进度有反应
	teaDemoTasks := []models.Task{
		{
			Name:        "王大娘的故事",
			Code:        "demo_tea_stories",
			Type:        "side",
			Description: "王大娘有讲不完的故事。去茶摊坐坐，听听她最近在说什么。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives: `[{"id":"obj_demo_03","type":"dialogue","target":"npc_tea_wang","description":"听王大娘讲故事","completed":false},
{"id":"obj_demo_04","type":"dialogue","target":"npc_tea_wang","description":"问问村里的八卦","completed":false}]`,
			Rewards:  `{"exp":30,"gold":0,"items":[{"code":"item_herb","count":2}]}`,
			NextTask: "",
		},
	}

	for i := range teaDemoTasks {
		if err := db.Create(&teaDemoTasks[i]).Error; err != nil {
			return err
		}
	}

	// 王大娘智能体增强
	var teaAgent models.Agent
	if err := db.Where("code = ?", "agent_tea_wang").First(&teaAgent).Error; err == nil {
		enhancedPrompt := teaAgent.SystemPrompt + `

【AI演示功能】
- 根据当前时间讲故事：
  - 早上（6-10点）："早起的鸟儿有虫吃！我跟你讲个早起的故事..."
  - 中午（10-14点）："日头正毒呢，来杯凉茶。你知道吗..."
  - 下午（14-18点）："下午茶时间！我听说..."
  - 晚上（18-6点）："天黑了，讲个有点吓人的故事..."
- 隐藏任务提示：
  - 如果玩家问"宝藏"或"山洞"，透露："村后的山洞据说有宝贝，但村长不让去。你要去的话，先去找猎户老周问问路。"
  - 如果玩家问"狼群"，说："狼群最近越来越凶了，猎户老周都愁坏了。你要是能帮他，他肯定会感谢你的。"
- 对玩家任务进度的反应：
  - 如果玩家完成了"初来乍到"，说："哟，你见过村长了？他老人家可喜欢你这样的年轻人。"
  - 如果玩家有了武器，说："哟，都买上兵器了？看来你是要出去闯荡啊！"
- 保持啰嗦但亲切的风格`
		db.Model(&teaAgent).Update("system_prompt", enhancedPrompt)
	}

	// ===== 场景3: 张铁匠的考验 =====
	// 铁匠测试玩家的武器知识
	// 回答正确给折扣
	// 玩家带材料可以打造特殊物品
	blacksmithDemoTasks := []models.Task{
		{
			Name:        "张铁匠的考验",
			Code:        "demo_blacksmith_test",
			Type:        "side",
			Description: "张铁匠想考考你对兵器的了解。通过考验可以获得折扣。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives: `[{"id":"obj_demo_05","type":"dialogue","target":"npc_blacksmith_zhang","description":"接受张铁匠的考验","completed":false},
{"id":"obj_demo_06","type":"dialogue","target":"npc_blacksmith_zhang","description":"回答铁匠的问题","completed":false}]`,
			Rewards:  `{"exp":80,"gold":0,"items":[]}`,
			NextTask: "",
		},
	}

	for i := range blacksmithDemoTasks {
		if err := db.Create(&blacksmithDemoTasks[i]).Error; err != nil {
			return err
		}
	}

	// 张铁匠智能体增强
	var blacksmithAgent models.Agent
	if err := db.Where("code = ?", "agent_blacksmith_zhang").First(&blacksmithAgent).Error; err == nil {
		enhancedPrompt := blacksmithAgent.SystemPrompt + `

【AI演示功能】
- 铁匠的考验（知识问答）：
  - 问题1："你知道铁剑和精钢刀有什么区别吗？"（正确答案：精钢刀更锋利，攻击+25）
  - 问题2："打造一把好剑需要什么材料？"（正确答案：好铁、炭火、淬火用水）
  - 问题3："兵器用久了会怎样？"（正确答案：会钝，需要磨砺）
- 如果玩家回答正确，给予折扣："嗯，你小子有两下子！给你打个八折吧。"
- 材料锻造：
  - 如果玩家说有"狼皮"或"材料"，回应："哦？你有材料？狼皮可以做皮甲，比店里卖的还结实。拿来我看看。"
  - 如果玩家说想"定制"武器，回应："定制兵器？好啊！不过得有好材料。你有精钢吗？有的话俺可以给你打一把好刀。"
- 保持朴实直接的风格`
		db.Model(&blacksmithAgent).Update("system_prompt", enhancedPrompt)
	}

	// ===== 场景4: 猎户老周的委托 =====
	// 猎人根据狼群状态给出动态任务
	// AI生成独特的任务描述
	// 奖励随任务难度变化
	hunterDemoTasks := []models.Task{
		{
			Name:        "猎户老周的委托",
			Code:        "demo_hunter_quest",
			Type:        "side",
			Description: "猎户老周需要帮手对付狼群。去找他问问情况吧。",
			Status:      "active",
			Trigger:     `{"type":"auto","conditions":[]}`,
			Objectives: `[{"id":"obj_demo_07","type":"dialogue","target":"npc_hunter_zhou","description":"了解狼群的情况","completed":false},
{"id":"obj_demo_08","type":"dialogue","target":"npc_hunter_zhou","description":"接受猎人的委托","completed":false}]`,
			Rewards:  `{"exp":100,"gold":200,"items":[{"code":"item_wolf_pelt","count":1}]}`,
			NextTask: "",
		},
	}

	for i := range hunterDemoTasks {
		if err := db.Create(&hunterDemoTasks[i]).Error; err != nil {
			return err
		}
	}

	// 猎户老周智能体增强
	var hunterAgent models.Agent
	if err := db.Where("code = ?", "agent_hunter_zhou").First(&hunterAgent).Error; err == nil {
		enhancedPrompt := hunterAgent.SystemPrompt + `

【AI演示功能】
- 动态任务生成：
  - 简单任务："村外有几只野兔，你去帮我抓两只回来。兔肉可以卖钱，皮毛也值点银子。"
  - 中等任务："最近狼群出没越来越频繁了。你去村外小路查看一下，看看有多少只狼。"
  - 困难任务："有一只头狼特别大，带着一群狼在村外游荡。你要是能把它赶走，我重重有赏！"
- 根据玩家等级调整任务难度：
  - 等级1-3：给简单任务（抓兔子、采草药）
  - 等级4-7：给中等任务（驱赶野狼、设置陷阱）
  - 等级8+：给困难任务（猎杀头狼、探索山洞）
- 任务描述示例（AI会动态生成）：
  - "俺昨天在村东头看到一只大灰狼，比一般的狼大一倍。它的左耳有个缺口，是以前被猎人伤过的。你小心点。"
  - "村外的树林里有个狼窝，里面有三只小狼崽。母狼肯定在附近，你去的时候要特别小心。"
- 奖励机制：
  - 简单任务：50经验 + 100金币
  - 中等任务：100经验 + 200金币 + 狼皮
  - 困难任务：200经验 + 500金币 + 稀有材料
- 保持沉默寡言但实在的风格`
		db.Model(&hunterAgent).Update("system_prompt", enhancedPrompt)
	}

	// ===== 更新老村长智能体，让他知道演示任务 =====
	var chiefAgent models.Agent
	if err := db.Where("code = ?", "agent_chief_chen").First(&chiefAgent).Error; err == nil {
		enhancedPrompt := chiefAgent.SystemPrompt + `

【AI演示功能】
- 当玩家问"有什么任务"或"我能帮什么忙"时，引导他们去找其他NPC：
  - "你去杂货铺找李掌柜，他最近进了一批新货，正愁没人买呢。"
  - "去茶摊找王大娘聊聊，她知道村里的大小事情。"
  - "铁匠张师傅最近在研究新兵器，你去跟他聊聊。"
  - "猎户老周最近为狼群的事发愁，你去帮帮他吧。"
- 根据玩家完成的任务给予反馈：
  - 如果玩家完成了"初来乍到"，说："呵呵，年轻人，你已经认识村里的人了。接下来准备做什么？"
  - 如果玩家问"外面的世界"，说："外面的世界很大，也很危险。你先在村里练练手，等有了本事再出去闯荡。"
- 保持慈祥和蔼的风格`
		db.Model(&chiefAgent).Update("system_prompt", enhancedPrompt)
	}

	log.Println("Demo scenarios seeded successfully!")
	return nil
}
