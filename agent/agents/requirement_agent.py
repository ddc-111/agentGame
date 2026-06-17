"""需求生成Agent - 负责分析市场趋势和提出新需求"""
import json
import requests
from typing import Dict, List, Optional
from datetime import datetime
from pathlib import Path


class RequirementAgent:
    """需求生成Agent - 分析市场、竞品、用户需求，提出新功能需求"""
    
    def __init__(self, llm_client=None):
        self.llm = llm_client
        self.requirements = []
        self.market_trends = []
        self.competitor_features = []
        
        # AgentNPC框架的核心目标
        self.framework_goals = {
            "core_features": [
                "完整的NPC AI对话系统",
                "多Agent协作机制",
                "动态任务生成与执行",
                "玩家行为分析与个性化",
                "实时状态同步"
            ],
            "modernization": [
                "微服务架构",
                "容器化部署",
                "CI/CD自动化",
                "监控与日志",
                "自动扩缩容"
            ],
            "ai_integration": [
                "多LLM支持（OpenAI/Claude/本地模型）",
                "Prompt模板管理",
                "上下文窗口优化",
                "流式响应",
                "AI决策链"
            ],
            "workflow": [
                "可视化工作流编辑器",
                "条件分支与循环",
                "事件触发机制",
                "定时任务",
                "异步任务队列"
            ],
            "developer_experience": [
                "完善的API文档",
                "SDK与示例代码",
                "本地开发环境",
                "调试工具",
                "性能分析工具"
            ],
            "competitiveness": [
                "高性能（低延迟）",
                "可扩展性",
                "安全性",
                "多语言支持",
                "插件系统"
            ]
        }
    
    def generate_requirements(self, code_analysis: Dict, test_analysis: Dict, current_features: List[str]) -> List[Dict]:
        """基于分析结果生成新需求"""
        requirements = []
        
        # 1. 基于框架目标生成需求
        requirements.extend(self._generate_goal_based_requirements(current_features))
        
        # 2. 基于代码分析生成需求
        requirements.extend(self._generate_analysis_based_requirements(code_analysis))
        
        # 3. 基于测试覆盖生成需求
        requirements.extend(self._generate_test_based_requirements(test_analysis))
        
        # 4. 基于市场趋势生成需求
        requirements.extend(self._generate_market_requirements())
        
        # 5. 使用LLM生成创新需求
        if self.llm and self.llm.is_configured():
            llm_requirements = self._generate_llm_requirements(code_analysis, current_features)
            requirements.extend(llm_requirements)
        
        # 去重并排序
        requirements = self._deduplicate_requirements(requirements)
        requirements.sort(key=lambda r: self._priority_score(r.get("priority", "low")))
        
        self.requirements = requirements
        return requirements
    
    def _generate_goal_based_requirements(self, current_features: List[str]) -> List[Dict]:
        """基于框架目标生成需求"""
        requirements = []
        
        for category, goals in self.framework_goals.items():
            for goal in goals:
                # 检查是否已实现
                if not self._is_feature_implemented(goal, current_features):
                    requirements.append({
                        "id": self._generate_id(),
                        "type": "feature",
                        "category": category,
                        "title": f"实现{goal}",
                        "description": f"作为AgentNPC框架的核心目标，需要实现{goal}",
                        "priority": "high" if category in ["core_features", "ai_integration"] else "medium",
                        "source": "framework_goals",
                        "acceptance_criteria": [
                            f"功能完整实现",
                            f"有单元测试覆盖",
                            f"有使用文档",
                            f"性能达标"
                        ],
                        "created_at": datetime.now().isoformat()
                    })
        
        return requirements
    
    def _generate_analysis_based_requirements(self, code_analysis: Dict) -> List[Dict]:
        """基于代码分析生成需求"""
        requirements = []
        
        # 检查代码复杂度
        complexity = code_analysis.get("complexity", {})
        if complexity.get("server", {}).get("functions", 0) > 100:
            requirements.append({
                "id": self._generate_id(),
                "type": "refactor",
                "category": "code_quality",
                "title": "重构Server端为微服务架构",
                "description": "当前Server端函数过多，建议拆分为独立的微服务",
                "priority": "high",
                "source": "code_analysis",
                "acceptance_criteria": [
                    "识别核心服务边界",
                    "实现服务间通信",
                    "保持向后兼容",
                    "添加服务发现机制"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        # 检查依赖
        deps = code_analysis.get("dependencies", {})
        if "phaser" in deps.get("client", []):
            requirements.append({
                "id": self._generate_id(),
                "type": "upgrade",
                "category": "modernization",
                "title": "评估并升级游戏引擎",
                "description": "当前使用Phaser 3，评估是否升级或迁移到更现代的引擎",
                "priority": "medium",
                "source": "code_analysis",
                "acceptance_criteria": [
                    "对比主流游戏引擎",
                    "评估迁移成本",
                    "原型验证",
                    "性能基准测试"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        return requirements
    
    def _generate_test_based_requirements(self, test_analysis: Dict) -> List[Dict]:
        """基于测试覆盖生成需求"""
        requirements = []
        
        # 检查测试覆盖率
        total_tests = test_analysis.get("total_test_cases", 0)
        if total_tests < 1000:
            requirements.append({
                "id": self._generate_id(),
                "type": "testing",
                "category": "quality",
                "title": "提升测试覆盖率到80%",
                "description": f"当前测试用例{total_tests}个，需要大幅增加",
                "priority": "high",
                "source": "test_analysis",
                "acceptance_criteria": [
                    "单元测试覆盖所有核心模块",
                    "集成测试覆盖所有API端点",
                    "E2E测试覆盖核心用户流程",
                    "性能测试基准建立"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        return requirements
    
    def _generate_market_requirements(self) -> List[Dict]:
        """基于市场趋势生成需求"""
        requirements = []
        
        # 当前AI游戏框架的市场趋势
        market_trends = [
            {
                "trend": "多模态AI交互",
                "description": "支持语音、图像、文本多模态输入输出",
                "priority": "high"
            },
            {
                "trend": "实时AI生成内容",
                "description": "动态生成对话、任务、场景等内容",
                "priority": "high"
            },
            {
                "trend": "个性化NPC",
                "description": "NPC具有记忆、情感、学习能力",
                "priority": "high"
            },
            {
                "trend": "分布式游戏架构",
                "description": "支持大规模并发和跨服交互",
                "priority": "medium"
            },
            {
                "trend": "低代码开发",
                "description": "可视化配置NPC行为和对话",
                "priority": "medium"
            },
            {
                "trend": "边缘计算",
                "description": "将AI推理部署到边缘节点降低延迟",
                "priority": "low"
            }
        ]
        
        for trend in market_trends:
            requirements.append({
                "id": self._generate_id(),
                "type": "innovation",
                "category": "market_trend",
                "title": f"支持{trend['trend']}",
                "description": trend["description"],
                "priority": trend["priority"],
                "source": "market_analysis",
                "acceptance_criteria": [
                    "完成技术调研",
                    "实现PoC原型",
                    "性能评估",
                    "用户反馈收集"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        return requirements
    
    def _generate_llm_requirements(self, code_analysis: Dict, current_features: List[str]) -> List[Dict]:
        """使用LLM生成创新需求"""
        if not self.llm or not self.llm.is_configured():
            return []
        
        prompt = f"""作为AgentNPC游戏框架的产品经理，基于以下信息提出新的功能需求：

当前特性：{json.dumps(current_features[:10], ensure_ascii=False)}

框架目标：
- 完整的NPC AI对话系统
- 多Agent协作机制
- 现代化架构
- AI友好接口
- 工作流支持

请提出5个创新的功能需求，要求：
1. 有竞争力
2. 技术可行
3. 用户价值高

返回JSON数组格式：
[{{"title": "需求标题", "description": "详细描述", "priority": "high/medium/low", "category": "分类"}}]"""
        
        try:
            response = self.llm.chat([{"role": "user", "content": prompt}])
            if response:
                # 提取JSON
                import re
                json_match = re.search(r'\[.*\]', response, re.DOTALL)
                if json_match:
                    llm_requirements = json.loads(json_match.group())
                    return [{
                        "id": self._generate_id(),
                        "type": "innovation",
                        "source": "llm_generation",
                        "acceptance_criteria": ["功能实现", "测试通过", "文档完善"],
                        "created_at": datetime.now().isoformat(),
                        **req
                    } for req in llm_requirements]
        except Exception as e:
            print(f"  [!] LLM需求生成失败: {e}")
        
        return []
    
    def _is_feature_implemented(self, feature: str, current_features: List[str]) -> bool:
        """检查特性是否已实现"""
        feature_lower = feature.lower()
        for current in current_features:
            if feature_lower in current.lower() or current.lower() in feature_lower:
                return True
        return False
    
    def _deduplicate_requirements(self, requirements: List[Dict]) -> List[Dict]:
        """去重需求"""
        seen_titles = set()
        unique = []
        for req in requirements:
            title = req.get("title", "")
            if title not in seen_titles:
                seen_titles.add(title)
                unique.append(req)
        return unique
    
    def _priority_score(self, priority: str) -> int:
        """优先级分数"""
        scores = {"critical": 0, "high": 1, "medium": 2, "low": 3}
        return scores.get(priority, 4)
    
    def _generate_id(self) -> str:
        """生成ID"""
        import uuid
        return f"req_{uuid.uuid4().hex[:8]}"
    
    def get_requirements_summary(self) -> Dict:
        """获取需求摘要"""
        return {
            "total": len(self.requirements),
            "by_priority": {
                "critical": len([r for r in self.requirements if r.get("priority") == "critical"]),
                "high": len([r for r in self.requirements if r.get("priority") == "high"]),
                "medium": len([r for r in self.requirements if r.get("priority") == "medium"]),
                "low": len([r for r in self.requirements if r.get("priority") == "low"])
            },
            "by_category": self._count_by_category(),
            "top_requirements": self.requirements[:5]
        }
    
    def _count_by_category(self) -> Dict:
        """按分类统计"""
        categories = {}
        for req in self.requirements:
            cat = req.get("category", "unknown")
            categories[cat] = categories.get(cat, 0) + 1
        return categories
    
    def save_requirements(self, file_path: Path):
        """保存需求到文件"""
        with open(file_path, 'w', encoding='utf-8') as f:
            json.dump({
                "generated_at": datetime.now().isoformat(),
                "summary": self.get_requirements_summary(),
                "requirements": self.requirements
            }, f, ensure_ascii=False, indent=2)
    
    def load_requirements(self, file_path: Path):
        """从文件加载需求"""
        if file_path.exists():
            with open(file_path, 'r', encoding='utf-8') as f:
                data = json.load(f)
                self.requirements = data.get("requirements", [])
