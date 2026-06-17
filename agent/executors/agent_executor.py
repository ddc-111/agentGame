"""Agent执行器 - 管理子Agent任务"""
import json
import subprocess
from pathlib import Path
from typing import Dict, List, Optional
from datetime import datetime


class AgentExecutor:
    """执行子Agent任务"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.tasks_dir = root_dir / "agent" / "tasks"
        self.tasks_dir.mkdir(parents=True, exist_ok=True)
    
    def execute_task(self, task: Dict) -> Dict:
        """执行单个任务"""
        task_type = task.get("type", "")
        target = task.get("target", "")
        
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "output": "",
            "error": ""
        }
        
        try:
            if task_type == "add_test":
                result = self._execute_add_test(task)
            elif task_type == "fix_test":
                result = self._execute_fix_test(task)
            elif task_type == "add_feature":
                result = self._execute_add_feature(task)
            elif task_type == "refactor":
                result = self._execute_refactor(task)
            elif task_type == "add_documentation":
                result = self._execute_add_docs(task)
            else:
                result["error"] = f"未知任务类型: {task_type}"
        except Exception as e:
            result["error"] = str(e)
        
        return result
    
    def execute_tasks(self, tasks: List[Dict]) -> List[Dict]:
        """执行多个任务"""
        results = []
        
        for task in tasks:
            print(f"执行任务: {task.get('title', task.get('type', 'unknown'))}")
            result = self.execute_task(task)
            results.append(result)
            
            status = "✓" if result["success"] else "✗"
            print(f"  {status} {'完成' if result['success'] else '失败'}")
        
        return results
    
    def _execute_add_test(self, task: Dict) -> Dict:
        """执行添加测试的任务"""
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "output": "",
            "error": ""
        }
        
        target = task.get("target", "")
        module = task.get("module", "")
        test_content = task.get("content", "")
        
        if not test_content:
            result["error"] = "未提供测试内容"
            return result
        
        # 确定测试文件路径
        if target == "server":
            test_file = self.root_dir / "server" / "internal" / f"{module}_test.go"
        elif target == "client":
            test_file = self.root_dir / "client" / "src" / "__tests__" / f"{module}.test.js"
        elif target == "gm":
            test_file = self.root_dir / "gm" / "src" / "__tests__" / f"{module}.test.js"
        else:
            result["error"] = f"未知目标: {target}"
            return result
        
        try:
            # 确保目录存在
            test_file.parent.mkdir(parents=True, exist_ok=True)
            
            # 写入测试文件
            test_file.write_text(test_content, encoding='utf-8')
            
            result["success"] = True
            result["output"] = f"测试文件已创建: {test_file}"
        except Exception as e:
            result["error"] = f"创建测试文件失败: {e}"
        
        return result
    
    def _execute_fix_test(self, task: Dict) -> Dict:
        """执行修复测试的任务"""
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "output": "",
            "error": ""
        }
        
        # 这里可以集成LLM来自动修复测试
        # 目前返回需要手动修复的提示
        result["error"] = "自动修复需要LLM集成，请手动修复"
        return result
    
    def _execute_add_feature(self, task: Dict) -> Dict:
        """执行添加功能的任务"""
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "output": "",
            "error": ""
        }
        
        # 这里可以集成LLM来生成代码
        result["error"] = "功能添加需要LLM集成，请手动实现"
        return result
    
    def _execute_refactor(self, task: Dict) -> Dict:
        """执行重构任务"""
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "output": "",
            "error": ""
        }
        
        result["error"] = "重构需要LLM集成，请手动重构"
        return result
    
    def _execute_add_docs(self, task: Dict) -> Dict:
        """执行添加文档的任务"""
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "output": "",
            "error": ""
        }
        
        doc_content = task.get("content", "")
        doc_file = task.get("file", "")
        
        if not doc_content or not doc_file:
            result["error"] = "未提供文档内容或文件路径"
            return result
        
        try:
            file_path = self.root_dir / doc_file
            file_path.parent.mkdir(parents=True, exist_ok=True)
            file_path.write_text(doc_content, encoding='utf-8')
            
            result["success"] = True
            result["output"] = f"文档已创建: {file_path}"
        except Exception as e:
            result["error"] = f"创建文档失败: {e}"
        
        return result
    
    def generate_sub_agent_prompt(self, task: Dict) -> str:
        """生成子Agent的提示词"""
        task_type = task.get("type", "")
        target = task.get("target", "")
        description = task.get("description", "")
        requirements = task.get("requirements", [])
        
        prompt = f"""你是一个专业的软件工程师Agent，负责为AgentGame项目执行以下任务：

## 任务类型
{task_type}

## 目标端
{target}

## 任务描述
{description}

## 具体要求
"""
        
        for i, req in enumerate(requirements, 1):
            prompt += f"{i}. {req}\n"
        
        prompt += f"""
## 项目结构
- server/: Go 1.21+ 后端 (Gin + GORM + WebSocket)
- client/: Phaser 3 游戏客户端 (Vite)
- gm/: Vue 3 GM编辑器 (Vite + Element Plus + Pinia)

## 测试要求
- 添加完整的单元测试
- 确保所有测试通过
- 测试覆盖率要高

请根据以上要求完成任务，并返回：
1. 修改的文件列表
2. 新增的测试文件
3. 测试结果
"""
        
        return prompt
