"""任务生成器 - 根据分析结果生成改进任务"""
from pathlib import Path
from typing import Dict, List
from datetime import datetime

from ..config import PRIORITY_CRITICAL, PRIORITY_HIGH, PRIORITY_MEDIUM, PRIORITY_LOW


class TaskGenerator:
    """生成改进任务"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
    
    def generate_tasks(self, gaps: Dict, code_analysis: Dict, test_analysis: Dict) -> List[Dict]:
        """根据差距分析生成任务"""
        tasks = []
        
        # 处理关键问题
        for issue in gaps.get("critical", []):
            task = self._create_task_from_issue(issue)
            if task:
                tasks.append(task)
        
        # 处理改进项
        for improvement in gaps.get("improvements", []):
            task = self._create_task_from_improvement(improvement)
            if task:
                tasks.append(task)
        
        # 处理优化机会
        for opportunity in gaps.get("opportunities", []):
            task = self._create_task_from_opportunity(opportunity)
            if task:
                tasks.append(task)
        
        # 生成测试补充任务
        tasks.extend(self._generate_test_tasks(test_analysis))
        
        # 生成文档任务
        tasks.extend(self._generate_doc_tasks(code_analysis))
        
        # 按优先级排序
        tasks.sort(key=lambda t: self._priority_score(t.get("priority", PRIORITY_LOW)))
        
        return tasks
    
    def _priority_score(self, priority: str) -> int:
        """优先级分数（越小越优先）"""
        scores = {
            PRIORITY_CRITICAL: 0,
            PRIORITY_HIGH: 1,
            PRIORITY_MEDIUM: 2,
            PRIORITY_LOW: 3
        }
        return scores.get(priority, 4)
    
    def _create_task_from_issue(self, issue: Dict) -> Dict:
        """从问题创建任务"""
        issue_type = issue.get("type", "")
        
        if issue_type == "test_failure":
            return {
                "id": self._generate_task_id(),
                "type": "fix_test",
                "title": f"修复失败的测试: {issue.get('test', '')}",
                "description": issue.get("error", ""),
                "target": issue.get("target", ""),
                "priority": PRIORITY_CRITICAL,
                "requirements": [
                    "分析测试失败原因",
                    "修复测试或被测代码",
                    "确保测试通过"
                ],
                "created_at": datetime.now().isoformat()
            }
        
        return None
    
    def _create_task_from_improvement(self, improvement: Dict) -> Dict:
        """从改进项创建任务"""
        imp_type = improvement.get("type", "")
        
        if imp_type == "coverage_gap":
            return {
                "id": self._generate_task_id(),
                "type": "add_test",
                "title": f"添加测试: {improvement.get('name', '')}",
                "description": improvement.get("suggestion", ""),
                "target": improvement.get("target", ""),
                "module": improvement.get("name", ""),
                "priority": improvement.get("priority", PRIORITY_MEDIUM),
                "requirements": [
                    f"为 {improvement.get('name', '')} 编写单元测试",
                    "覆盖主要功能点",
                    "测试边界条件",
                    "确保测试可重复运行"
                ],
                "created_at": datetime.now().isoformat()
            }
        
        return None
    
    def _create_task_from_opportunity(self, opportunity: Dict) -> Dict:
        """从优化机会创建任务"""
        opp_type = opportunity.get("type", "")
        
        task_map = {
            "missing_layer": {
                "type": "add_feature",
                "title": f"添加 {opportunity.get('package', '')} 模块",
                "requirements": [
                    "设计模块接口",
                    "实现核心功能",
                    "添加单元测试",
                    "更新文档"
                ]
            },
            "typescript_migration": {
                "type": "refactor",
                "title": "迁移到TypeScript",
                "requirements": [
                    "添加tsconfig配置",
                    "逐步迁移JS文件到TS",
                    "添加类型定义",
                    "更新构建配置"
                ]
            },
            "framework_upgrade": {
                "type": "refactor",
                "title": f"升级 {opportunity.get('target', '')} 框架",
                "requirements": [
                    "评估新版本兼容性",
                    "更新依赖",
                    "修复破坏性变更",
                    "运行回归测试"
                ]
            }
        }
        
        template = task_map.get(opp_type, {})
        if template:
            return {
                "id": self._generate_task_id(),
                "type": template.get("type", "unknown"),
                "title": template.get("title", ""),
                "description": opportunity.get("suggestion", ""),
                "target": opportunity.get("target", ""),
                "priority": PRIORITY_LOW,
                "requirements": template.get("requirements", []),
                "created_at": datetime.now().isoformat()
            }
        
        return None
    
    def _generate_test_tasks(self, test_analysis: Dict) -> List[Dict]:
        """生成测试补充任务"""
        tasks = []
        
        # 为未测试的包生成任务
        for pkg in test_analysis.get("server", {}).get("untested_packages", []):
            tasks.append({
                "id": self._generate_task_id(),
                "type": "add_test",
                "title": f"为 server/{pkg} 添加测试",
                "description": f"包 {pkg} 缺少单元测试",
                "target": "server",
                "module": pkg,
                "priority": PRIORITY_HIGH,
                "requirements": [
                    f"分析 {pkg} 包的功能",
                    "编写单元测试覆盖主要功能",
                    "测试错误处理",
                    "确保测试独立可运行"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        # 为未测试的模块生成任务
        for module in test_analysis.get("client", {}).get("untested_modules", []):
            tasks.append({
                "id": self._generate_task_id(),
                "type": "add_test",
                "title": f"为 client/{module} 添加测试",
                "description": f"模块 {module} 缺少单元测试",
                "target": "client",
                "module": module.replace("/", "_"),
                "priority": PRIORITY_MEDIUM,
                "requirements": [
                    f"分析 {module} 模块的功能",
                    "使用Vitest编写单元测试",
                    "Mock外部依赖",
                    "测试边界条件"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        # 为未测试的Store生成任务
        for store in test_analysis.get("gm", {}).get("untested_stores", []):
            tasks.append({
                "id": self._generate_task_id(),
                "type": "add_test",
                "title": f"为 gm/stores/{store} 添加测试",
                "description": f"Store {store} 缺少单元测试",
                "target": "gm",
                "module": f"stores/{store}",
                "priority": PRIORITY_HIGH,
                "requirements": [
                    f"分析 {store} Store的状态管理",
                    "使用Pinia测试工具编写测试",
                    "测试actions和getters",
                    "测试状态变更"
                ],
                "created_at": datetime.now().isoformat()
            })
        
        return tasks
    
    def _generate_doc_tasks(self, code_analysis: Dict) -> List[Dict]:
        """生成文档任务"""
        tasks = []
        
        # 如果API端点较多但没有API文档
        endpoints = code_analysis.get("server", {}).get("endpoints", [])
        if len(endpoints) > 5:
            tasks.append({
                "id": self._generate_task_id(),
                "type": "add_documentation",
                "title": "生成API文档",
                "description": "项目有多个API端点，需要完善API文档",
                "target": "server",
                "priority": PRIORITY_MEDIUM,
                "requirements": [
                    "列出所有API端点",
                    "描述请求参数",
                    "描述响应格式",
                    "添加使用示例"
                ],
                "file": "docs/API.md",
                "content": self._generate_api_doc_template(endpoints),
                "created_at": datetime.now().isoformat()
            })
        
        return tasks
    
    def _generate_api_doc_template(self, endpoints: List[str]) -> str:
        """生成API文档模板"""
        doc = """# AgentGame API 文档

## 概述
本文档描述了AgentGame游戏服务器提供的REST API。

## 基础信息
- Base URL: `http://localhost:8080/api`
- Content-Type: `application/json`

## API 端点

"""
        for endpoint in endpoints:
            parts = endpoint.split(" ", 1)
            if len(parts) == 2:
                method, path = parts
                doc += f"### {method} {path}\n\n"
                doc += f"**描述**: TODO\n\n"
                doc += f"**请求参数**:\n```json\n{{}}\n```\n\n"
                doc += f"**响应示例**:\n```json\n{{\"data\": {{}}}}\n```\n\n"
                doc += "---\n\n"
        
        return doc
    
    def _generate_task_id(self) -> str:
        """生成任务ID"""
        import uuid
        return f"task_{uuid.uuid4().hex[:8]}"
