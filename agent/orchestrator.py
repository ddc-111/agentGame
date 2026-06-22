"""Agent调度器核心 - 协调分析、测试、任务生成的循环"""
import json
import time
import sys
import io
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Optional

from .config import (
    PROJECT_ROOT, REPORTS_DIR, TASKS_DIR, HISTORY_DIR, REQUIREMENTS_DIR,
    PRIORITY_CRITICAL, PRIORITY_HIGH, PRIORITY_MEDIUM
)
from .analyzers.code_analyzer import CodeAnalyzer
from .analyzers.test_analyzer import TestAnalyzer
from .analyzers.gap_analyzer import GapAnalyzer
from .executors.test_executor import TestExecutor
from .executors.build_executor import BuildExecutor
from .executors.task_executor import TaskExecutor
from .generators.task_generator import TaskGenerator
from .generators.report_generator import ReportGenerator
from .utils.llm_client import LLMClient
from .utils.history_tracker import HistoryTracker
from .config_manager import AgentConfig
from .agents.requirement_agent import RequirementAgent


class Orchestrator:
    """自循环Agent调度器"""
    
    def __init__(self, max_iterations: int = 10, verbose: bool = True, config: AgentConfig = None):
        self.max_iterations = max_iterations
        self.verbose = verbose
        self.iteration = 0
        self.history: List[Dict] = []
        
        # 加载配置
        self.config = config or AgentConfig.from_env()
        
        # 修复Windows编码
        self._fix_encoding()
        
        # 初始化各模块
        self.code_analyzer = CodeAnalyzer(PROJECT_ROOT)
        self.test_analyzer = TestAnalyzer(PROJECT_ROOT)
        self.gap_analyzer = GapAnalyzer(PROJECT_ROOT)
        self.test_executor = TestExecutor(PROJECT_ROOT)
        self.build_executor = BuildExecutor(PROJECT_ROOT)
        self.task_generator = TaskGenerator(PROJECT_ROOT)
        self.report_generator = ReportGenerator(REPORTS_DIR)
        self.history_tracker = HistoryTracker(HISTORY_DIR)
        
        # LLM集成
        self.llm = LLMClient(
            api_url=self.config.get("llm.api_url"),
            api_key=self.config.get("llm.api_key"),
            model=self.config.get("llm.model")
        )
        self.task_executor = TaskExecutor(PROJECT_ROOT, self.llm)
        
        # 需求生成Agent
        self.requirement_agent = RequirementAgent(self.llm)
    
    def _fix_encoding(self):
        """修复Windows编码问题"""
        if sys.platform == 'win32':
            # 设置stdout编码
            if hasattr(sys.stdout, 'reconfigure'):
                try:
                    sys.stdout.reconfigure(encoding='utf-8', errors='replace')
                except:
                    pass
            # 设置环境变量
            import os
            os.environ["PYTHONIOENCODING"] = "utf-8"
    
    def log(self, message: str, level: str = "INFO"):
        """日志输出"""
        if self.verbose:
            timestamp = datetime.now().strftime("%H:%M:%S")
            prefix = {"INFO": "[i]", "SUCCESS": "[+]", "WARNING": "[!]", "ERROR": "[-]"}
            try:
                print(f"[{timestamp}] {prefix.get(level, '[.]')} {message}")
            except UnicodeEncodeError:
                # 如果还是有编码问题，使用ASCII
                safe_message = message.encode('ascii', 'replace').decode('ascii')
                print(f"[{timestamp}] {prefix.get(level, '[.]')} {safe_message}")
    
    def run_cycle(self) -> Dict:
        """运行一个完整的自循环周期"""
        self.iteration += 1
        cycle_start = time.time()
        
        self.log(f"=== 开始第 {self.iteration} 轮迭代 ===", "INFO")
        
        result = {
            "iteration": self.iteration,
            "timestamp": datetime.now().isoformat(),
            "phases": {}
        }
        
        # 阶段1: 代码分析
        self.log("阶段1: 分析代码库结构...", "INFO")
        code_analysis = self.code_analyzer.analyze()
        result["phases"]["code_analysis"] = code_analysis
        self.log(f"  发现 {code_analysis['total_files']} 个文件, {code_analysis['total_lines']} 行代码", "SUCCESS")
        
        # 阶段2: 测试分析
        self.log("阶段2: 分析测试覆盖...", "INFO")
        test_analysis = self.test_analyzer.analyze()
        result["phases"]["test_analysis"] = test_analysis
        self.log(f"  测试文件: {test_analysis['total_test_files']}, 测试用例: {test_analysis['total_test_cases']}", "SUCCESS")
        
        # 阶段3: 构建验证
        self.log("阶段3: 执行构建验证...", "INFO")
        build_result = self.build_executor.build_all()
        result["phases"]["build"] = build_result
        for target, success in build_result.items():
            status = "SUCCESS" if success else "ERROR"
            self.log(f"  {target}: {'通过' if success else '失败'}", status)
        
        # 阶段4: 执行测试
        self.log("阶段4: 执行测试套件...", "INFO")
        test_result = self.test_executor.run_all_tests()
        result["phases"]["tests"] = test_result
        self.log(f"  通过: {test_result['passed']}, 失败: {test_result['failed']}, 跳过: {test_result['skipped']}", "SUCCESS")
        
        # 阶段5: 差距分析
        self.log("阶段5: 分析改进差距...", "INFO")
        gaps = self.gap_analyzer.analyze_gaps(code_analysis, test_analysis, test_result)
        result["phases"]["gaps"] = gaps
        self.log(f"  发现 {len(gaps['critical'])} 个关键问题, {len(gaps['improvements'])} 个改进点", "WARNING" if gaps['critical'] else "SUCCESS")
        
        # 阶段6: 生成任务（从现有需求加载）
        self.log("阶段6: 从现有需求加载任务...", "INFO")
        
        # 首先从gap分析生成任务
        tasks = self.task_generator.generate_tasks(gaps, code_analysis, test_analysis)
        
        # 然后加载现有需求文件中的任务
        requirement_tasks = self._load_requirements_as_tasks()
        tasks.extend(requirement_tasks)
        
        # 去重（基于title）
        seen_titles = set()
        unique_tasks = []
        for task in tasks:
            title = task.get("title", "")
            if title not in seen_titles:
                seen_titles.add(title)
                unique_tasks.append(task)
        
        result["phases"]["tasks"] = unique_tasks
        self.log(f"  加载 {len(unique_tasks)} 个任务（需求: {len(requirement_tasks)}, 分析: {len(tasks) - len(requirement_tasks)}）", "SUCCESS")
        
        # 阶段6.5: 跳过需求生成（用户要求完成现有需求，不生成新需求）
        result["phases"]["requirements"] = []
        self.log("阶段6.5: 跳过需求生成（完成现有需求模式）", "INFO")
        
        # 阶段7: 执行任务（如果启用LLM）
        if self.config.is_llm_enabled() and self.config.get("tasks.auto_execute"):
            self.log("阶段7: 执行任务...", "INFO")
            exec_results = self._execute_tasks(tasks)
            result["phases"]["task_execution"] = exec_results
        else:
            self.log("阶段7: 跳过任务执行 (LLM未启用或未配置自动执行)", "INFO")
        
        # 阶段8: 生成报告
        self.log("阶段8: 生成报告...", "INFO")
        report_path = self.report_generator.generate(result)
        result["report_path"] = str(report_path)
        self.log(f"  报告已保存: {report_path}", "SUCCESS")
        
        # 计算总耗时
        result["duration"] = time.time() - cycle_start
        self.history.append(result)
        
        # 记录历史
        iter_id = self.history_tracker.record_iteration(result)
        self.log(f"  历史记录: {iter_id}", "INFO")
        
        self.log(f"=== 第 {self.iteration} 轮迭代完成 (耗时: {result['duration']:.1f}秒) ===", "SUCCESS")
        
        return result
    
    def _execute_tasks(self, tasks: List[Dict]) -> Dict:
        """执行任务"""
        results = {
            "total": len(tasks),
            "executed": 0,
            "success": 0,
            "failed": 0,
            "details": []
        }
        
        # 只执行高优先级任务
        priority_filter = self.config.get("tasks.priority_filter", ["critical", "high"])
        
        for task in tasks:
            if task.get("priority") not in priority_filter:
                continue
            
            self.log(f"    执行: {task.get('title', '')[:50]}...", "INFO")
            result = self.task_executor.execute_with_verification(task)
            results["executed"] += 1
            
            if result["success"]:
                results["success"] += 1
                self.log(f"      成功", "SUCCESS")
            else:
                results["failed"] += 1
                self.log(f"      失败: {result.get('error', '')[:50]}", "ERROR")
            
            results["details"].append(result)
            
            # 避免过于频繁的执行
            time.sleep(0.5)
        
        return results
    
    def run(self) -> List[Dict]:
        """运行完整的自循环迭代"""
        self.log("AgentGame 自循环Agent启动", "INFO")
        self.log(f"最大迭代次数: {self.max_iterations}", "INFO")
        self.log(f"LLM状态: {'已配置' if self.config.is_llm_enabled() else '未配置'}", "INFO")
        
        # 用于检测循环是否有效果
        previous_metrics = None
        stable_iterations = 0
        max_stable_iterations = 3  # 连续3轮无变化则停止
        
        while self.iteration < self.max_iterations:
            try:
                result = self.run_cycle()
                
                # 获取当前指标
                current_metrics = {
                    "critical": len(result["phases"]["gaps"]["critical"]),
                    "improvements": len(result["phases"]["gaps"]["improvements"]),
                    "tasks": len(result["phases"]["tasks"]),
                    "test_passed": result["phases"]["tests"]["passed"],
                    "test_failed": result["phases"]["tests"]["failed"]
                }
                
                # 检查是否还有关键问题需要解决
                if current_metrics["critical"] == 0:
                    self.log("所有关键问题已解决！", "SUCCESS")
                    
                    # 检查是否有高优先级任务
                    high_priority_tasks = [t for t in result["phases"]["tasks"] 
                                          if t.get("priority") in [PRIORITY_CRITICAL, PRIORITY_HIGH]]
                    
                    # 如果没有高优先级任务，系统已稳定
                    if len(high_priority_tasks) == 0 and self.config.get("iterations.stop_on_stable"):
                        self.log("系统已达到稳定状态，结束迭代", "SUCCESS")
                        break
                    
                    # 检查指标是否变化
                    if previous_metrics and self._metrics_equal(current_metrics, previous_metrics):
                        stable_iterations += 1
                        self.log(f"指标无变化 ({stable_iterations}/{max_stable_iterations})", "WARNING")
                        
                        if stable_iterations >= max_stable_iterations:
                            self.log(f"连续 {max_stable_iterations} 轮指标无变化，结束迭代", "SUCCESS")
                            break
                    else:
                        stable_iterations = 0
                    
                    previous_metrics = current_metrics.copy()
                
                # 保存任务到文件供子Agent执行
                self._save_tasks(result["phases"]["tasks"])
                
                # 显示趋势
                self._show_trend()
                
                delay = self.config.get("iterations.delay_between", 1)
                self.log(f"等待 {delay} 秒后进行下一轮迭代...", "INFO")
                time.sleep(delay)
                
            except KeyboardInterrupt:
                self.log("用户中断，停止迭代", "WARNING")
                break
            except Exception as e:
                self.log(f"迭代出错: {e}", "ERROR")
                if self.config.get("iterations.stop_on_critical"):
                    break
                continue
        
        self.log(f"自循环完成，共 {self.iteration} 轮迭代", "SUCCESS")
        return self.history
    
    def _metrics_equal(self, metrics1: Dict, metrics2: Dict) -> bool:
        """比较两组指标是否相等"""
        return (
            metrics1.get("critical") == metrics2.get("critical") and
            metrics1.get("improvements") == metrics2.get("improvements") and
            metrics1.get("tasks") == metrics2.get("tasks") and
            metrics1.get("test_passed") == metrics2.get("test_passed") and
            metrics1.get("test_failed") == metrics2.get("test_failed")
        )
    
    def _show_trend(self):
        """显示趋势"""
        trend = self.history_tracker.get_trend()
        
        if trend.get("trend") == "insufficient_data":
            return
        
        score_trend = trend.get("overall_score", {})
        direction = score_trend.get("direction", "stable")
        
        if direction == "improving":
            self.log(f"  趋势: 质量提升中 (+{score_trend.get('change', 0):.1f})", "SUCCESS")
        elif direction == "declining":
            self.log(f"  趋势: 质量下降 ({score_trend.get('change', 0):.1f})", "WARNING")
        else:
            self.log(f"  趋势: 稳定", "INFO")
    
    def _save_tasks(self, tasks: List[Dict]):
        """保存任务到文件"""
        task_file = TASKS_DIR / f"tasks_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        with open(task_file, 'w', encoding='utf-8') as f:
            json.dump({
                "iteration": self.iteration,
                "generated_at": datetime.now().isoformat(),
                "tasks": tasks
            }, f, ensure_ascii=False, indent=2)
        self.log(f"  任务已保存: {task_file}", "INFO")
    
    def _load_requirements_as_tasks(self) -> List[Dict]:
        """加载需求文件并转换为任务格式"""
        tasks = []
        
        # 查找最新的需求文件
        req_files = sorted(REQUIREMENTS_DIR.glob("requirements_*.json"), reverse=True)
        if not req_files:
            self.log("  未找到需求文件", "WARNING")
            return tasks
        
        # 加载最新的需求文件
        latest_req_file = req_files[0]
        self.log(f"  加载需求文件: {latest_req_file.name}", "INFO")
        
        try:
            with open(latest_req_file, 'r', encoding='utf-8') as f:
                req_data = json.load(f)
            
            requirements = req_data.get("requirements", [])
            
            for req in requirements:
                # 跳过已完成的需求（如果有状态字段）
                if req.get("status") == "completed":
                    continue
                
                # 将需求转换为任务格式
                task = {
                    "id": req.get("id", f"req_{len(tasks)}"),
                    "type": self._map_req_type_to_task_type(req.get("type", "feature")),
                    "title": req.get("title", ""),
                    "description": req.get("description", ""),
                    "target": self._determine_target_from_req(req),
                    "priority": req.get("priority", "medium"),
                    "requirements": req.get("acceptance_criteria", []),
                    "source": "requirement",
                    "category": req.get("category", ""),
                    "created_at": req.get("created_at", datetime.now().isoformat())
                }
                
                # 根据类别设置更具体的target
                category = req.get("category", "")
                if category in ["ai_integration", "NPC智能基础", "AI对话系统深化", "多Agent协作机制"]:
                    task["target"] = "server"
                elif category in ["developer_experience", "workflow", "开发者工具与易用性"]:
                    # 可能涉及多个目标
                    if "前端" in req.get("title", "") or "GM" in req.get("title", ""):
                        task["target"] = "gm"
                    elif "客户端" in req.get("title", "") or "Phaser" in req.get("title", ""):
                        task["target"] = "client"
                    else:
                        task["target"] = "server"
                elif category in ["core_features", "quality", "code_quality"]:
                    task["target"] = "server"
                elif category in ["modernization", "competitiveness"]:
                    task["target"] = "server"
                
                tasks.append(task)
            
            # 按优先级排序
            priority_order = {"critical": 0, "high": 1, "medium": 2, "low": 3}
            tasks.sort(key=lambda t: priority_order.get(t.get("priority", "low"), 4))
            
        except Exception as e:
            self.log(f"  加载需求文件失败: {e}", "ERROR")
        
        return tasks
    
    def _map_req_type_to_task_type(self, req_type: str) -> str:
        """将需求类型映射为任务类型"""
        mapping = {
            "feature": "add_feature",
            "refactor": "refactor",
            "testing": "add_test",
            "innovation": "add_feature",
            "upgrade": "refactor",
            "documentation": "add_documentation"
        }
        return mapping.get(req_type, "add_feature")
    
    def _determine_target_from_req(self, req: Dict) -> str:
        """根据需求确定目标（server/client/gm）"""
        title = req.get("title", "").lower()
        description = req.get("description", "").lower()
        category = req.get("category", "")
        
        # 前端相关
        if any(kw in title or kw in description for kw in ["前端", "gm", "编辑器", "vue", "组件", "可视化"]):
            return "gm"
        
        # 客户端相关
        if any(kw in title or kw in description for kw in ["客户端", "client", "phaser", "游戏场景", "渲染"]):
            return "client"
        
        # 默认服务器端
        return "server"
    
    def _extract_current_features(self, code_analysis: Dict) -> List[str]:
        """提取当前已实现的功能"""
        features = []
        
        # 从包结构提取
        server_packages = code_analysis.get("server", {}).get("packages", [])
        for pkg in server_packages:
            features.append(f"server/{pkg}模块")
        
        # 从模型提取
        models = code_analysis.get("server", {}).get("models", [])
        for model in models[:10]:  # 只取前10个
            features.append(f"{model}数据模型")
        
        # 从场景提取
        scenes = code_analysis.get("client", {}).get("scenes", [])
        for scene in scenes:
            features.append(f"{scene}游戏场景")
        
        # 从系统提取
        systems = code_analysis.get("client", {}).get("systems", [])
        for system in systems:
            features.append(f"{system}游戏系统")
        
        return features
    
    def get_summary(self) -> Dict:
        """获取运行摘要"""
        if not self.history:
            return {"status": "未运行"}
        
        last = self.history[-1]
        trend = self.history_tracker.get_trend()
        
        return {
            "total_iterations": self.iteration,
            "last_iteration": last["iteration"],
            "last_report": last.get("report_path"),
            "critical_issues": len(last["phases"]["gaps"]["critical"]),
            "improvements": len(last["phases"]["gaps"]["improvements"]),
            "pending_tasks": len(last["phases"]["tasks"]),
            "test_passed": last["phases"]["tests"]["passed"],
            "test_failed": last["phases"]["tests"]["failed"],
            "trend": trend,
            "llm_enabled": self.config.is_llm_enabled()
        }
    
    def get_history_comparison(self, limit: int = 5) -> Dict:
        """获取历史对比"""
        history = self.history_tracker.get_history(limit)
        
        if len(history) < 2:
            return {"message": "历史数据不足"}
        
        return {
            "iterations": history,
            "trend": self.history_tracker.get_trend()
        }
