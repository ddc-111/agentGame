package demos

// CombatDemo showcases the combat system
var CombatDemo = DemoFlow{
	Name:        "战斗系统演示",
	Description: "展示回合制战斗的完整机制，包括普通攻击、技能使用、道具使用、逃跑和胜利奖励",
	Steps: []DemoStep{
		{
			Name:        "准备阶段",
			Description: "确保角色已装备武器，携带草药",
			Action:      "打开背包确认装备和道具",
			Expected:    "角色装备铁剑（攻击+10），背包有草药3份、金疮药2份",
			Duration:    5,
		},
		{
			Name:        "进入战斗",
			Description: "在村外小路遭遇野狼",
			Action:      "在村外小路移动，触发随机战斗",
			Expected:    "进入战斗界面，显示敌方「野狼」信息：HP 80，攻击 12，防御 3。显示己方回合",
			Scene:       "scene_village_path",
			Duration:    5,
		},
		{
			Name:        "普通攻击",
			Description: "使用普通攻击造成基础伤害",
			Action:      "点击「攻击」按钮",
			Expected:    "角色挥剑攻击，造成伤害 = 攻击力(20) - 敌方防御(3) = 17点伤害。野狼HP从80变为63",
			Duration:    5,
		},
		{
			Name:        "敌方回合",
			Description: "野狼进行反击",
			Action:      "等待敌方行动",
			Expected:    "野狼扑咬，造成伤害 = 12 - 己方防御(5) = 7点。己方HP减少",
			Duration:    5,
		},
		{
			Name:        "使用技能",
			Description: "使用「基础斩击」技能造成额外伤害",
			Action:      "点击技能按钮，选择「基础斩击」（消耗5MP）",
			Expected:    "技能释放，造成1.5倍武器伤害 = 25点伤害。MP减少5点。野狼HP从63变为38。显示技能特效",
			Duration:    8,
		},
		{
			Name:        "使用增益技能",
			Description: "使用「战吼」提升攻击力",
			Action:      "点击技能按钮，选择「战吼」（消耗8MP）",
			Expected:    "攻击力提升30%，持续2回合。显示buff图标。攻击力从20变为26",
			Duration:    8,
		},
		{
			Name:        "敌方攻击",
			Description: "野狼再次攻击",
			Action:      "等待敌方行动",
			Expected:    "野狼攻击，造成7点伤害",
			Duration:    5,
		},
		{
			Name:        "使用道具",
			Description: "生命值较低时使用草药恢复",
			Action:      "点击道具按钮，选择「草药」使用",
			Expected:    "HP恢复20点。草药数量减少1份",
			Duration:    5,
		},
		{
			Name:        "暴击演示",
			Description: "触发暴击造成高额伤害",
			Action:      "继续攻击（有概率触发暴击）",
			Expected:    "触发暴击（10%概率），伤害 = 26 * 1.5 = 39点。显示暴击特效和数字",
			Duration:    5,
		},
		{
			Name:        "击败敌人",
			Description: "将野狼HP降为0",
			Action:      "继续攻击直到击败野狼",
			Expected:    "野狼被击败，显示胜利界面。获得：经验值100，金币50，狼皮×1。触发成就「初试牛刀」",
			Duration:    10,
		},
		{
			Name:        "逃跑演示",
			Description: "展示逃跑机制（新战斗）",
			Action:      "再次遭遇野狼，选择「逃跑」",
			Expected:    "有概率逃跑成功（基础60%）。成功则脱离战斗，失败则敌方攻击一次后可再次选择",
			Duration:    10,
		},
		{
			Name:        "战斗失败演示",
			Description: "展示HP归零的处理",
			Action:      "不使用道具，让HP降为0",
			Expected:    "角色倒下，显示「战斗失败」。自动使用回城符回到村庄，HP恢复为50%",
			Duration:    10,
		},
	},
}

// GetCombatDemo returns the combat demo flow
func GetCombatDemo() DemoFlow {
	return CombatDemo
}
