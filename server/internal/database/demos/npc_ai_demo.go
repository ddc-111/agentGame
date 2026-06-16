package demos

// NPCAIDemo showcases AI-powered NPC conversations
var NPCAIDemo = DemoFlow{
	Name:        "AI对话演示",
	Description: "展示AI驱动的NPC对话系统，包括个性化回答、记忆功能、情绪变化和动态任务生成",
	Steps: []DemoStep{
		{
			Name:        "初次对话 - 老村长",
			Description: "与老村长进行初次对话，观察AI的个性化回答",
			Action:      "与老村长对话，说「你好，我是新来的冒险者」",
			Expected:    "老村长用慈祥的语气回应，自我介绍，询问玩家从哪里来。使用「老朽」自称，称呼玩家为「年轻人」",
			NPC:         "npc_chief_chen",
			Duration:    15,
		},
		{
			Name:        "追问村庄历史",
			Description: "继续询问村庄的历史背景",
			Action:      "问「这个村子有什么故事吗？」",
			Expected:    "老村长讲述青石村的由来——因盛产青石而得名，已有数百年历史。提到村口石碑的来历",
			NPC:         "npc_chief_chen",
			Duration:    15,
		},
		{
			Name:        "请求任务指引",
			Description: "询问有什么可以帮忙的",
			Action:      "说「有什么我可以帮忙的吗？」",
			Expected:    "老村长引导玩家去找其他NPC：去杂货铺找李掌柜、去茶摊找王大娘、去铁匠铺找张师傅、帮猎户老周对付狼群",
			NPC:         "npc_chief_chen",
			Duration:    15,
		},
		{
			Name:        "智能商品推荐 - 李掌柜",
			Description: "与商人对话，观察AI的推荐逻辑",
			Action:      "与李掌柜对话，说「我想买点东西，有什么推荐的？」",
			Expected:    "李掌柜根据玩家等级推荐商品。等级1-5推荐基础装备（草药、馒头、新手木剑）。使用「在下」自称，称呼「客官」",
			NPC:         "npc_merchant_li",
			Duration:    15,
		},
		{
			Name:        "询问特定商品",
			Description: "询问特定类型的商品",
			Action:      "说「我需要治疗类的药品」",
			Expected:    "李掌柜详细介绍草药（50文，恢复20HP）、灵芝（200文，恢复100HP）、金疮药（100文，恢复50HP）的区别和性价比",
			NPC:         "npc_merchant_li",
			Duration:    15,
		},
		{
			Name:        "讨价还价",
			Description: "尝试与商人讨价还价",
			Action:      "说「能便宜点吗？我是穷学生」",
			Expected:    "李掌柜犹豫后可能给折扣（「看客官这么有诚意，给打个九折吧」）。展示AI的灵活应变能力",
			NPC:         "npc_merchant_li",
			Duration:    15,
		},
		{
			Name:        "时间感知 - 王大娘",
			Description: "根据游戏时间触发不同对话",
			Action:      "与王大娘对话",
			Expected:    "王大娘根据当前时间讲故事：早上讲早起故事，中午推荐凉茶，下午分享八卦，晚上讲吓人故事",
			NPC:         "npc_tea_wang",
			Duration:    15,
		},
		{
			Name:        "隐藏信息获取",
			Description: "通过特定关键词获取隐藏信息",
			Action:      "问「你知道哪里有宝藏吗？」",
			Expected:    "王大娘透露：「村后的山洞据说有宝贝，但村长不让去。你要去的话，先去找猎户老周问问路。」展示AI知识库检索",
			NPC:         "npc_tea_wang",
			Duration:    15,
		},
		{
			Name:        "NPC情绪变化",
			Description: "观察NPC对不同话题的情绪反应",
			Action:      "与王大娘聊「狼群」",
			Expected:    "王大娘语气变得担忧：「狼群最近越来越凶了，猎户老周都愁坏了。你要是能帮他，他肯定会感谢你的。」",
			NPC:         "npc_tea_wang",
			Duration:    15,
		},
		{
			Name:        "知识问答 - 张铁匠",
			Description: "回答铁匠的兵器知识问题",
			Action:      "与张铁匠对话，接受考验。回答「铁剑和精钢刀的区别是精钢刀更锋利」",
			Expected:    "张铁匠认可回答：「嗯，你小子有两下子！给你打个八折吧。」展示AI判断对错的能力",
			NPC:         "npc_blacksmith_zhang",
			Duration:    20,
		},
		{
			Name:        "动态任务生成 - 猎户老周",
			Description: "根据玩家等级获取不同难度的任务",
			Action:      "与猎户老周对话，询问有什么需要帮忙的",
			Expected:    "老周根据玩家等级给出不同任务：1-3级抓兔子，4-7级驱赶野狼，8+级猎杀头狼。展示AI动态任务生成",
			NPC:         "npc_hunter_zhou",
			Duration:    15,
		},
		{
			Name:        "任务进度反馈",
			Description: "完成任务后NPC的反应",
			Action:      "完成老周的任务后再次对话",
			Expected:    "老周态度变得友善：「干得不错。这是你的报酬。」如果之前对话过，会记得玩家。展示AI记忆功能",
			NPC:         "npc_hunter_zhou",
			Duration:    15,
		},
		{
			Name:        "跨NPC关联",
			Description: "NPC之间的信息关联",
			Action:      "完成猎户任务后与老村长对话",
			Expected:    "老村长知道玩家帮了老周：「听说你帮老周解决了狼群的问题？年轻人，好样的！」展示NPC间信息共享",
			NPC:         "npc_chief_chen",
			Duration:    15,
		},
	},
}

// GetNPCAIDemo returns the NPC AI demo flow
func GetNPCAIDemo() DemoFlow {
	return NPCAIDemo
}
