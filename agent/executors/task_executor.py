"""任务执行反馈循环 - 自动执行任务并验证结果"""
import json
import time
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Optional

from ..utils.llm_client import LLMClient
from ..executors.test_executor import TestExecutor


class TaskExecutor:
    """任务执行器 - 支持LLM自动执行和验证"""
    
    def __init__(self, root_dir: Path, llm_client: LLMClient = None):
        self.root_dir = root_dir
        self.llm = llm_client or LLMClient()
        self.test_executor = TestExecutor(root_dir)
        self.execution_log = []
    
    def execute_with_verification(self, task: Dict, max_retries: int = 3) -> Dict:
        """执行任务并验证结果"""
        result = {
            "task": task,
            "timestamp": datetime.now().isoformat(),
            "success": False,
            "attempts": 0,
            "verification": None,
            "error": None
        }
        
        for attempt in range(max_retries):
            result["attempts"] = attempt + 1
            
            # 执行任务
            exec_result = self._execute_task(task)
            
            if not exec_result["success"]:
                result["error"] = exec_result.get("error", "执行失败")
                continue
            
            # 验证结果
            verify_result = self._verify_task(task)
            result["verification"] = verify_result
            
            if verify_result["passed"]:
                result["success"] = True
                break
            else:
                # 如果验证失败，尝试修复
                if attempt < max_retries - 1:
                    self._fix_after_verification(task, verify_result)
        
        self.execution_log.append(result)
        return result
    
    def _execute_task(self, task: Dict) -> Dict:
        """执行单个任务"""
        task_type = task.get("type", "")
        
        if task_type == "add_test":
            return self._execute_add_test(task)
        elif task_type == "fix_test":
            return self._execute_fix_test(task)
        elif task_type == "add_feature":
            return self._execute_add_feature(task)
        elif task_type == "refactor":
            return self._execute_refactor(task)
        else:
            return {"success": False, "error": f"未知任务类型: {task_type}"}
    
    def _execute_add_test(self, task: Dict) -> Dict:
        """执行添加测试的任务"""
        target = task.get("target", "")
        module = task.get("module", "")
        
        # 获取源代码
        source_code = self._get_source_code(target, module)
        if not source_code:
            return {"success": False, "error": "无法获取源代码"}
        
        # 使用LLM生成测试
        language = self._get_language(target)
        test_code = self.llm.generate_test(source_code, language, module)
        
        if not test_code:
            return {"success": False, "error": "LLM生成测试失败"}
        
        # 清理生成的代码，移除Markdown标记
        test_code = self._clean_llm_output(test_code)
        
        # 写入测试文件
        test_file = self._get_test_file_path(target, module)
        test_file.parent.mkdir(parents=True, exist_ok=True)
        test_file.write_text(test_code, encoding='utf-8')
        
        return {"success": True, "file": str(test_file)}
    
    def _execute_fix_test(self, task: Dict) -> Dict:
        """执行修复测试的任务"""
        target = task.get("target", "")
        test_name = task.get("test", "")
        error = task.get("error", "")
        
        # 获取测试文件内容
        test_file = self._find_test_file(target, test_name)
        if not test_file:
            return {"success": False, "error": "找不到测试文件"}
        
        test_code = test_file.read_text(encoding='utf-8')
        language = self._get_language(target)
        
        # 使用LLM修复
        fixed_code = self.llm.fix_code(test_code, error, language)
        
        if not fixed_code:
            return {"success": False, "error": "LLM修复失败"}
        
        # 写入修复后的代码
        test_file.write_text(fixed_code, encoding='utf-8')
        
        return {"success": True, "file": str(test_file)}
    
    def _execute_add_feature(self, task: Dict) -> Dict:
        """执行添加功能的任务"""
        # 功能添加需要更复杂的逻辑，暂时返回需要手动处理
        return {"success": False, "error": "功能添加需要手动处理"}
    
    def _execute_refactor(self, task: Dict) -> Dict:
        """执行重构任务"""
        target = task.get("target", "")
        file_path = task.get("file", "")
        
        if not file_path:
            return {"success": False, "error": "未指定要重构的文件"}
        
        # 获取文件完整路径
        full_path = self._get_file_path(target, file_path)
        if not full_path or not full_path.exists():
            return {"success": False, "error": f"文件不存在: {file_path}"}
        
        # 读取源代码
        source_code = full_path.read_text(encoding='utf-8')
        language = self._get_language(target)
        
        # 使用LLM分析并重构代码
        refactored_code = self.llm.analyze_and_improve(source_code, language)
        
        if not refactored_code:
            return {"success": False, "error": "LLM重构失败"}
        
        # 清理LLM输出
        refactored_code = self._clean_llm_output(refactored_code)
        
        # 备份原文件
        backup_path = full_path.with_suffix(full_path.suffix + '.backup')
        backup_path.write_text(source_code, encoding='utf-8')
        
        # 写入重构后的代码
        full_path.write_text(refactored_code, encoding='utf-8')
        
        return {
            "success": True,
            "file": str(full_path),
            "backup": str(backup_path),
            "original_lines": len(source_code.split('\n')),
            "refactored_lines": len(refactored_code.split('\n'))
        }
    
    def _get_file_path(self, target: str, file_path: str) -> Path:
        """获取文件完整路径"""
        if target == "server":
            return self.root_dir / "server" / "internal" / file_path
        elif target == "client":
            return self.root_dir / "client" / "src" / file_path
        elif target == "gm":
            return self.root_dir / "gm" / "src" / file_path
        else:
            return Path()
    
    def _verify_task(self, task: Dict) -> Dict:
        """验证任务结果"""
        task_type = task.get("type", "")
        target = task.get("target", "")
        
        verification = {
            "passed": False,
            "tests_run": 0,
            "tests_passed": 0,
            "tests_failed": 0,
            "details": []
        }
        
        # 运行相关测试
        if target == "server":
            result = self.test_executor.run_server_tests()
        elif target == "client":
            result = self.test_executor.run_client_tests()
        elif target == "gm":
            result = self.test_executor.run_gm_tests()
        else:
            return verification
        
        verification["tests_run"] = result.get("passed", 0) + result.get("failed", 0)
        verification["tests_passed"] = result.get("passed", 0)
        verification["tests_failed"] = result.get("failed", 0)
        verification["details"] = result.get("failures", [])
        verification["passed"] = result.get("failed", 0) == 0
        
        return verification
    
    def _fix_after_verification(self, task: Dict, verification: Dict):
        """验证失败后尝试修复"""
        failures = verification.get("details", [])
        
        for failure in failures[:3]:  # 只处理前3个失败
            error = failure.get("error", "")
            if error:
                # 更新任务的错误信息
                task["error"] = error
                # 尝试修复
                self._execute_fix_test(task)
    
    def _get_source_code(self, target: str, module: str) -> Optional[str]:
        """获取源代码"""
        if target == "server":
            # 尝试多个可能的路径
            paths = [
                self.root_dir / "server" / "internal" / f"{module}.go",
                self.root_dir / "server" / "internal" / module / f"{module}.go",
                self.root_dir / "server" / "internal" / module / "main.go",
                self.root_dir / "server" / "internal" / module / "agent.go",
                self.root_dir / "server" / "internal" / module / "config.go",
                self.root_dir / "server" / "internal" / module / "database.go",
                self.root_dir / "server" / "internal" / module / "generator.go",
                self.root_dir / "server" / "internal" / module / "mcp.go",
                self.root_dir / "server" / "internal" / module / "server.go",
                self.root_dir / "server" / "internal" / module / "chat.go",
                self.root_dir / "server" / "internal" / module / "memory.go",
            ]
        elif target == "client":
            paths = [
                self.root_dir / "client" / "src" / "game" / "systems" / f"{module}.js",
                self.root_dir / "client" / "src" / "game" / "scenes" / f"{module}.js",
                self.root_dir / "client" / "src" / "game" / f"{module}.js",
            ]
        elif target == "gm":
            # 处理 stores/xxx.js 或直接 xxx.js
            if "/" in module:
                # stores/achievement.js -> gm/src/stores/achievement.js
                paths = [
                    self.root_dir / "gm" / "src" / module,
                    self.root_dir / "gm" / "src" / module.replace("/", "\\"),
                ]
            else:
                paths = [
                    self.root_dir / "gm" / "src" / "stores" / f"{module}.js",
                    self.root_dir / "gm" / "src" / "components" / f"{module}.vue",
                    self.root_dir / "gm" / "src" / "views" / f"{module}.vue",
                ]
        else:
            return None
        
        for path in paths:
            if path.exists():
                try:
                    return path.read_text(encoding='utf-8')
                except Exception:
                    continue
        
        # 如果都没找到，尝试glob搜索
        if target == "server":
            # 查找目录下的所有.go文件
            module_dir = self.root_dir / "server" / "internal" / module
            if module_dir.exists() and module_dir.is_dir():
                # 合并目录下所有.go文件的内容
                contents = []
                for go_file in module_dir.glob("*.go"):
                    if "_test" not in go_file.name:
                        try:
                            contents.append(go_file.read_text(encoding='utf-8'))
                        except Exception:
                            continue
                if contents:
                    return "\n\n".join(contents)
            
            # 尝试glob搜索
            pattern = f"**/{module}.go"
            for path in (self.root_dir / "server").rglob(pattern):
                if path.is_file() and "_test" not in path.name:
                    try:
                        return path.read_text(encoding='utf-8')
                    except Exception:
                        continue
        elif target == "gm":
            pattern = f"**/{module}" if "/" in module else f"**/{module}.*"
            for path in (self.root_dir / "gm" / "src").rglob(pattern):
                if path.is_file():
                    try:
                        return path.read_text(encoding='utf-8')
                    except Exception:
                        continue
        
        return None
    
    def _get_test_file_path(self, target: str, module: str) -> Path:
        """获取测试文件路径"""
        if target == "server":
            return self.root_dir / "server" / "internal" / f"{module}_test.go"
        elif target == "client":
            return self.root_dir / "client" / "src" / "__tests__" / f"{module}.test.js"
        elif target == "gm":
            return self.root_dir / "gm" / "src" / "__tests__" / "stores" / f"{module}.test.js"
        else:
            return Path()
    
    def _find_test_file(self, target: str, test_name: str) -> Optional[Path]:
        """查找测试文件"""
        if target == "server":
            test_dir = self.root_dir / "server"
        elif target == "client":
            test_dir = self.root_dir / "client" / "src" / "__tests__"
        elif target == "gm":
            test_dir = self.root_dir / "gm" / "src" / "__tests__"
        else:
            return None
        
        # 搜索匹配的测试文件
        for test_file in test_dir.rglob("*test*"):
            if test_file.is_file():
                content = test_file.read_text(encoding='utf-8')
                if test_name in content:
                    return test_file
        
        return None
    
    def _clean_llm_output(self, code: str) -> str:
        """清理LLM输出，移除Markdown标记"""
        import re
        
        # 移除Markdown代码块标记
        # 匹配 ```javascript, ```js, ```vue, ```go 等
        code = re.sub(r'^```\w*\n', '', code, flags=re.MULTILINE)
        code = re.sub(r'\n```$', '', code, flags=re.MULTILINE)
        
        # 移除开头的空行
        code = code.strip()
        
        # 如果代码仍然以```开头，移除第一行
        if code.startswith('```'):
            lines = code.split('\n', 1)
            if len(lines) > 1:
                code = lines[1]
        
        return code
    
    def _get_language(self, target: str) -> str:
        """获取编程语言"""
        if target == "server":
            return "go"
        elif target == "client":
            return "javascript"
        elif target == "gm":
            return "javascript"
        else:
            return "javascript"
    
    def get_execution_summary(self) -> Dict:
        """获取执行摘要"""
        total = len(self.execution_log)
        success = sum(1 for r in self.execution_log if r["success"])
        
        return {
            "total_tasks": total,
            "successful": success,
            "failed": total - success,
            "success_rate": (success / total * 100) if total > 0 else 0,
            "details": self.execution_log
        }
