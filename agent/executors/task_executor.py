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
        elif task_type == "add_documentation":
            return self._execute_add_documentation(task)
        else:
            # 尝试作为通用任务处理
            return self._execute_generic_task(task)
    
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
        target = task.get("target", "")
        title = task.get("title", "")
        description = task.get("description", "")
        requirements = task.get("requirements", [])
        
        # 获取相关上下文
        context = self._get_feature_context(target, task)
        
        # 确定语言
        language = self._get_language(target)
        
        # 使用LLM生成功能代码
        feature_code = self.llm.generate_feature(
            task_description=f"{title}: {description}",
            requirements=requirements,
            language=language,
            context=context
        )
        
        if not feature_code:
            return {"success": False, "error": "LLM生成功能代码失败"}
        
        # 清理LLM输出
        feature_code = self._clean_llm_output(feature_code)
        
        # 确定文件路径
        file_path = self._determine_feature_file(target, task, feature_code)
        if not file_path:
            return {"success": False, "error": "无法确定文件路径"}
        
        # 创建目录
        file_path.parent.mkdir(parents=True, exist_ok=True)
        
        # 备份原文件（如果存在）
        backup_path = None
        if file_path.exists():
            backup_path = file_path.with_suffix(file_path.suffix + '.backup')
            backup_path.write_text(file_path.read_text(encoding='utf-8'), encoding='utf-8')
        
        # 写入文件
        file_path.write_text(feature_code, encoding='utf-8')
        
        result = {
            "success": True,
            "file": str(file_path),
            "lines": len(feature_code.split('\n'))
        }
        
        if backup_path:
            result["backup"] = str(backup_path)
        
        return result
    
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
    
    def _execute_add_documentation(self, task: Dict) -> Dict:
        """执行添加文档的任务"""
        target = task.get("target", "")
        file_path = task.get("file", "")
        content = task.get("content", "")
        title = task.get("title", "")
        
        if not file_path:
            # 自动生成文档路径
            if target == "server":
                file_path = "docs/API.md"
            else:
                file_path = f"docs/{title.replace(' ', '_')}.md"
        
        # 确定完整路径
        if target == "server":
            full_path = self.root_dir / "server" / file_path
        elif target == "gm":
            full_path = self.root_dir / "gm" / file_path
        elif target == "client":
            full_path = self.root_dir / "client" / file_path
        else:
            full_path = self.root_dir / file_path
        
        # 如果没有预设内容，使用LLM生成
        if not content:
            context = self._get_feature_context(target, task)
            content = self.llm.generate_feature(
                task_description=f"{title}: {task.get('description', '')}",
                requirements=task.get("requirements", []),
                language="markdown",
                context=context,
                file_type="documentation"
            )
            
            if not content:
                return {"success": False, "error": "LLM生成文档失败"}
            
            content = self._clean_llm_output(content)
        
        # 创建目录
        full_path.parent.mkdir(parents=True, exist_ok=True)
        
        # 写入文件
        full_path.write_text(content, encoding='utf-8')
        
        return {
            "success": True,
            "file": str(full_path),
            "lines": len(content.split('\n'))
        }
    
    def _execute_generic_task(self, task: Dict) -> Dict:
        """执行通用任务"""
        target = task.get("target", "")
        title = task.get("title", "")
        description = task.get("description", "")
        requirements = task.get("requirements", [])
        
        # 获取上下文
        context = self._get_feature_context(target, task)
        
        # 确定语言
        language = self._get_language(target)
        
        # 使用LLM生成代码
        code = self.llm.generate_feature(
            task_description=f"{title}: {description}",
            requirements=requirements,
            language=language,
            context=context
        )
        
        if not code:
            return {"success": False, "error": "LLM生成代码失败"}
        
        # 清理输出
        code = self._clean_llm_output(code)
        
        # 确定文件路径
        file_path = self._determine_feature_file(target, task, code)
        if not file_path:
            return {"success": False, "error": "无法确定文件路径"}
        
        # 创建目录
        file_path.parent.mkdir(parents=True, exist_ok=True)
        
        # 备份原文件（如果存在）
        backup_path = None
        if file_path.exists():
            backup_path = file_path.with_suffix(file_path.suffix + '.backup')
            backup_path.write_text(file_path.read_text(encoding='utf-8'), encoding='utf-8')
        
        # 写入文件
        file_path.write_text(code, encoding='utf-8')
        
        result = {
            "success": True,
            "file": str(file_path),
            "lines": len(code.split('\n'))
        }
        
        if backup_path:
            result["backup"] = str(backup_path)
        
        return result
    
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
                # stores/achievement -> gm/src/stores/achievement.js
                module_with_ext = f"{module}.js"
                paths = [
                    self.root_dir / "gm" / "src" / module_with_ext,
                    self.root_dir / "gm" / "src" / module,
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
    
    def _get_feature_context(self, target: str, task: Dict) -> str:
        """获取功能开发的上下文"""
        context_parts = []
        
        # 获取相关现有代码
        if target == "server":
            # 获取主要的模型定义
            models_file = self.root_dir / "server" / "internal" / "database" / "models" / "models.go"
            if models_file.exists():
                try:
                    content = models_file.read_text(encoding='utf-8')
                    context_parts.append(f"=== 数据模型 ===\n{content[:3000]}")
                except:
                    pass
            
            # 获取路由定义
            server_file = self.root_dir / "server" / "internal" / "network" / "server.go"
            if server_file.exists():
                try:
                    content = server_file.read_text(encoding='utf-8')
                    context_parts.append(f"=== 路由定义 ===\n{content[:2000]}")
                except:
                    pass
            
            # 获取repository接口
            repo_file = self.root_dir / "server" / "internal" / "database" / "repository" / "repository.go"
            if repo_file.exists():
                try:
                    content = repo_file.read_text(encoding='utf-8')
                    context_parts.append(f"=== Repository ===\n{content[:2000]}")
                except:
                    pass
        
        elif target == "gm":
            # 获取API客户端
            api_dir = self.root_dir / "gm" / "src" / "api"
            if api_dir.exists():
                for api_file in api_dir.glob("*.js"):
                    try:
                        content = api_file.read_text(encoding='utf-8')
                        context_parts.append(f"=== API {api_file.name} ===\n{content[:1500]}")
                    except:
                        pass
            
            # 获取现有Store示例
            stores_dir = self.root_dir / "gm" / "src" / "stores"
            if stores_dir.exists():
                for store_file in list(stores_dir.glob("*.js"))[:2]:
                    try:
                        content = store_file.read_text(encoding='utf-8')
                        context_parts.append(f"=== Store {store_file.name} ===\n{content[:1500]}")
                    except:
                        pass
        
        elif target == "client":
            # 获取游戏系统示例
            systems_dir = self.root_dir / "client" / "src" / "game" / "systems"
            if systems_dir.exists():
                for sys_file in list(systems_dir.glob("*.js"))[:2]:
                    try:
                        content = sys_file.read_text(encoding='utf-8')
                        context_parts.append(f"=== System {sys_file.name} ===\n{content[:1500]}")
                    except:
                        pass
        
        return "\n\n".join(context_parts) if context_parts else ""
    
    def _determine_feature_file(self, target: str, task: Dict, code: str) -> Optional[Path]:
        """确定功能代码的文件路径"""
        title = task.get("title", "").lower()
        task_type = task.get("type", "")
        
        # 从代码中尝试提取package名或模块名
        package_hint = self._extract_package_from_code(code, target)
        
        if target == "server":
            # 根据任务标题推断目录
            if "handler" in title or "api" in title or "endpoint" in title:
                base_dir = self.root_dir / "server" / "internal" / "network"
            elif "model" in title or "entity" in title or "database" in title:
                base_dir = self.root_dir / "server" / "internal" / "database" / "models"
            elif "repository" in title or "repo" in title or "data access" in title:
                base_dir = self.root_dir / "server" / "internal" / "database" / "repository"
            elif "agent" in title or "ai" in title or "llm" in title:
                base_dir = self.root_dir / "server" / "internal" / "agent"
            elif "game" in title or "logic" in title or "combat" in title or "npc" in title:
                base_dir = self.root_dir / "server" / "internal" / "game"
            elif "config" in title or "setting" in title:
                base_dir = self.root_dir / "server" / "internal" / "config"
            elif "test" in title:
                base_dir = self.root_dir / "server" / "internal" / "tests"
            else:
                base_dir = self.root_dir / "server" / "internal" / package_hint if package_hint else self.root_dir / "server" / "internal" / "game"
            
            # 生成文件名
            filename = self._generate_filename(title, "go")
            return base_dir / filename
        
        elif target == "gm":
            if "store" in title or "state" in title:
                base_dir = self.root_dir / "gm" / "src" / "stores"
                filename = self._generate_filename(title, "js")
            elif "component" in title or "vue" in title:
                base_dir = self.root_dir / "gm" / "src" / "components"
                filename = self._generate_filename(title, "vue")
            elif "view" in title or "page" in title:
                base_dir = self.root_dir / "gm" / "src" / "views"
                filename = self._generate_filename(title, "vue")
            elif "api" in title or "service" in title:
                base_dir = self.root_dir / "gm" / "src" / "api"
                filename = self._generate_filename(title, "js")
            else:
                base_dir = self.root_dir / "gm" / "src" / "stores"
                filename = self._generate_filename(title, "js")
            
            return base_dir / filename
        
        elif target == "client":
            if "scene" in title:
                base_dir = self.root_dir / "client" / "src" / "game" / "scenes"
            elif "system" in title or "manager" in title:
                base_dir = self.root_dir / "client" / "src" / "game" / "systems"
            else:
                base_dir = self.root_dir / "client" / "src" / "game" / "systems"
            
            filename = self._generate_filename(title, "js")
            return base_dir / filename
        
        return None
    
    def _extract_package_from_code(self, code: str, target: str) -> str:
        """从代码中提取包名"""
        import re
        
        if target == "server":
            # Go: package xxx
            match = re.search(r'^package\s+(\w+)', code, re.MULTILINE)
            if match:
                return match.group(1)
        
        return ""
    
    def _generate_filename(self, title: str, extension: str) -> str:
        """根据标题生成文件名"""
        import re
        
        # 移除中文和特殊字符，保留英文
        english_parts = re.findall(r'[a-zA-Z]+', title.lower())
        
        if not english_parts:
            # 如果没有英文，使用默认名
            return f"new_feature.{extension}"
        
        # 组合文件名
        filename = "_".join(english_parts[:3])  # 最多取3个词
        
        # 添加常见后缀
        if extension == "go":
            if not filename.endswith("_handler") and not filename.endswith("_test"):
                if "handler" in title.lower() or "api" in title.lower():
                    filename += "_handlers"
                elif "test" in title.lower():
                    filename += "_test"
        
        return f"{filename}.{extension}"
    
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
