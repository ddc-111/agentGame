package demos

// DemoFlow defines a scripted demo sequence
type DemoFlow struct {
	Name        string
	Description string
	Steps       []DemoStep
}

// DemoStep is a single step in a demo
type DemoStep struct {
	Name        string
	Description string
	Action      string // what the demo does
	Expected    string // what should happen
	Scene       string // scene code
	NPC         string // npc code (if applicable)
	Duration    int    // estimated seconds
}

// BeginnerVillageDemo showcases the full game loop
var BeginnerVillageDemo = DemoFlow{
	Name:        "新手村完整流程",
	Description: "展示从角色创建到完成第一个任务的完整游戏流程，包含对话、购物、战斗、成就等核心系统",
	Steps: []DemoStep{
		{
			Name:        "进入游戏",
			Description: "玩家创建角色，出现在青石村村口",
			Action:      "创建新角色，选择默认外观，点击开始游戏",
			Expected:    "角色出现在村口场景，显示欢迎提示，触发「初来乍到」任务",
			Scene:       "scene_village_entrance",
			Duration:    5,
		},
		{
			Name:        "认识老村长",
			Description: "与村口的老村长对话，了解村庄情况",
			Action:      "靠近老村长NPC，点击对话按钮",
			Expected:    "老村长用慈祥的语气欢迎玩家，介绍青石村的情况，指引玩家去认识其他村民。任务「初来乍到」目标完成",
			Scene:       "scene_village_entrance",
			NPC:         "npc_chief_chen",
			Duration:    15,
		},
		{
			Name:        "前往村中心",
			Description: "通过传送点前往村中心广场",
			Action:      "走到村口右侧的传送点",
			Expected:    "场景切换到村中心广场，可以看到大榕树、老井和来往的村民",
			Scene:       "scene_village_center",
			Duration:    5,
		},
		{
			Name:        "拜访王大娘茶摊",
			Description: "在茶摊与王大娘聊天，了解村里的八卦",
			Action:      "走到茶摊区域，与王大娘对话",
			Expected:    "王大娘热情招呼，分享村里的故事和八卦。「熟悉村庄」任务进度更新",
			Scene:       "scene_tea_stand",
			NPC:         "npc_tea_wang",
			Duration:    15,
		},
		{
			Name:        "前往杂货铺",
			Description: "通过传送点进入杂货铺",
			Action:      "从村中心走到杂货铺入口传送点",
			Expected:    "场景切换到杂货铺，看到李掌柜在柜台后",
			Scene:       "scene_general_store",
			Duration:    5,
		},
		{
			Name:        "购买补给品",
			Description: "与李掌柜对话，购买草药和馒头",
			Action:      "与李掌柜对话，选择购买3份草药和5个馒头",
			Expected:    "李掌柜热情推荐商品，完成交易。金币减少，背包增加道具。「采购补给」任务进度更新。触发成就「初来乍到」",
			Scene:       "scene_general_store",
			NPC:         "npc_merchant_li",
			Duration:    20,
		},
		{
			Name:        "前往铁匠铺",
			Description: "进入铁匠铺购买武器",
			Action:      "从杂货铺返回村中心，走到铁匠铺入口",
			Expected:    "场景切换到铁匠铺，炉火通红，张铁匠在打铁",
			Scene:       "scene_blacksmith",
			Duration:    5,
		},
		{
			Name:        "购买铁剑",
			Description: "与张铁匠对话，购买一把铁剑",
			Action:      "与张铁匠对话，选择购买铁剑（500文）",
			Expected:    "张铁匠展示铁剑，完成交易。装备铁剑后攻击力提升。「装备升级」任务完成，获得奖励",
			Scene:       "scene_blacksmith",
			NPC:         "npc_blacksmith_zhang",
			Duration:    20,
		},
		{
			Name:        "前往村外小路",
			Description: "前往村外寻找猎户老周",
			Action:      "从村中心右侧传送点前往村外小路",
			Expected:    "场景切换到村外小路，背景音乐变化，可以看到猎户老周在巡逻",
			Scene:       "scene_village_path",
			Duration:    5,
		},
		{
			Name:        "与猎户老周对话",
			Description: "了解狼群的情况，接受委托",
			Action:      "与猎户老周对话，了解狼群威胁",
			Expected:    "老周简短地说明情况，玩家接受任务。「初试身手」任务目标更新",
			Scene:       "scene_village_path",
			NPC:         "npc_hunter_zhou",
			Duration:    15,
		},
		{
			Name:        "遭遇野狼",
			Description: "在村外小路遭遇野狼，进入战斗",
			Action:      "在小路上探索，触发随机遭遇",
			Expected:    "进入回合制战斗界面。可以使用普通攻击和基础斩击技能。击败野狼获得经验和金币奖励。触发成就「初试牛刀」",
			Scene:       "scene_village_path",
			Duration:    30,
		},
		{
			Name:        "使用道具恢复",
			Description: "战斗后使用草药恢复生命值",
			Action:      "打开背包，使用草药",
			Expected:    "生命值恢复20点，道具消耗1份",
			Duration:    5,
		},
		{
			Name:        "返回村庄交任务",
			Description: "回到村中心向村长汇报",
			Action:      "使用回城符或走回村庄",
			Expected:    "回到村口，与村长对话完成任务链。获得大量经验、金币和道具奖励。等级提升",
			Scene:       "scene_village_entrance",
			NPC:         "npc_chief_chen",
			Duration:    15,
		},
		{
			Name:        "成就展示",
			Description: "查看获得的成就",
			Action:      "打开成就面板",
			Expected:    "显示已获得的成就：「初来乍到」「初试牛刀」「村庄之友」等，展示成就系统",
			Duration:    10,
		},
	},
}

// GetBeginnerVillageDemo returns the beginner village demo flow
func GetBeginnerVillageDemo() DemoFlow {
	return BeginnerVillageDemo
}
