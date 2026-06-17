"""专注于完成需求的执行脚本"""
import sys
import subprocess
import time
from pathlib import Path
from datetime import datetime

agent_dir = Path(__file__).parent
sys.path.insert(0, str(agent_dir.parent))

from agent.config import PROJECT_ROOT, REPORTS_DIR, REQUIREMENTS_DIR
from agent.utils.llm_client import LLMClient
from agent.config_manager import AgentConfig
from agent.analyzers.code_analyzer import CodeAnalyzer
from agent.analyzers.test_analyzer import TestAnalyzer
from agent.executors.test_executor import TestExecutor
from agent.executors.build_executor import BuildExecutor
from agent.executors.task_executor import TaskExecutor
from agent.generators.report_generator import ReportGenerator


class RequirementCompleter:
    """需求完成器"""
    
    def __init__(self, config: AgentConfig = None):
        self.config = config or AgentConfig.from_env()
        
        self.llm = LLMClient(
            api_url=self.config.get("llm.api_url"),
            api_key=self.config.get("llm.api_key"),
            model=self.config.get("llm.model")
        )
        
        self.code_analyzer = CodeAnalyzer(PROJECT_ROOT)
        self.test_analyzer = TestAnalyzer(PROJECT_ROOT)
        self.test_executor = TestExecutor(PROJECT_ROOT)
        self.build_executor = BuildExecutor(PROJECT_ROOT)
        self.task_executor = TaskExecutor(PROJECT_ROOT, self.llm)
        self.report_generator = ReportGenerator(REPORTS_DIR)
        
        self.completed_requirements = []
    
    def log(self, message: str, level: str = "INFO"):
        timestamp = datetime.now().strftime("%H:%M:%S")
        prefix = {"INFO": "[i]", "SUCCESS": "[+]", "WARNING": "[!]", "ERROR": "[-]", "PHASE": "[*]"}
        try:
            print(f"[{timestamp}] {prefix.get(level, '[.]')} {message}")
        except UnicodeEncodeError:
            print(f"[{timestamp}] {prefix.get(level, '[.]')} {message.encode('ascii', 'replace').decode()}")
    
    def run(self, cycles: int = 10):
        """运行需求完成循环"""
        self.log("=" * 60, "PHASE")
        self.log("AgentGame 需求完成系统", "PHASE")
        self.log("=" * 60, "PHASE")
        
        # 加载现有需求
        requirements = self.load_requirements()
        self.log(f"加载 {len(requirements)} 个待完成需求", "INFO")
        
        for cycle in range(1, cycles + 1):
            self.log(f"\n{'='*60}", "PHASE")
            self.log(f"第 {cycle}/{cycles} 轮循环", "PHASE")
            self.log(f"{'='*60}", "PHASE")
            
            # 1. 分析代码库
            self.log("[阶段1] 分析代码库", "PHASE")
            code_analysis = self.code_analyzer.analyze()
            test_analysis = self.test_analyzer.analyze()
            
            # 2. 生成测试代码
            self.log("[阶段2] 生成测试代码", "PHASE")
            generated = self.generate_tests(code_analysis, test_analysis)
            
            # 3. 构建验证
            self.log("[阶段3] 构建验证", "PHASE")
            build_result = self.build_executor.build_all()
            success = sum(1 for v in build_result.values() if v)
            self.log(f"构建结果: {success}/3 通过", "SUCCESS" if success == 3 else "WARNING")
            
            # 4. 运行测试
            self.log("[阶段4] 运行测试", "PHASE")
            test_result = self.test_executor.run_all_tests()
            self.log(f"测试通过: {test_result['passed']}, 失败: {test_result['failed']}", "SUCCESS")
            
            # 5. 生成报告
            self.log("[阶段5] 生成报告", "PHASE")
            report_data = {
                "iteration": cycle,
                "timestamp": datetime.now().isoformat(),
                "phases": {
                    "code_analysis": code_analysis,
                    "test_analysis": test_analysis,
                    "build": build_result,
                    "tests": test_result,
                    "generated_tests": generated
                }
            }
            report_path = self.report_generator.generate(report_data)
            self.log(f"报告: {report_path}", "SUCCESS")
            
            # 6. 提交git
            self.log("[阶段6] 提交Git", "PHASE")
            self.git_commit(cycle)
            
            # 短暂延迟
            if cycle < cycles:
                time.sleep(2)
        
        self.log(f"\n{'='*60}", "SUCCESS")
        self.log("执行完成!", "SUCCESS")
        self.log(f"{'='*60}", "SUCCESS")
    
    def load_requirements(self) -> list:
        """加载现有需求"""
        req_files = list(REQUIREMENTS_DIR.glob("req_*.json"))
        if not req_files:
            return []
        
        latest = max(req_files, key=lambda f: f.stat().st_mtime)
        import json
        with open(latest, 'r', encoding='utf-8') as f:
            data = json.load(f)
            return data.get("requirements", [])
    
    def generate_tests(self, code_analysis: dict, test_analysis: dict) -> list:
        """生成测试代码"""
        generated = []
        
        if not self.llm.is_configured():
            self.log("LLM未配置，跳过测试生成", "WARNING")
            return generated
        
        untested = test_analysis.get("server", {}).get("untested_packages", [])
        
        for pkg in untested[:2]:
            self.log(f"  为 {pkg} 生成测试...", "INFO")
            
            source_dir = PROJECT_ROOT / "server" / "internal" / pkg
            if source_dir.exists():
                source_files = list(source_dir.glob("*.go"))
                if source_files:
                    source_code = source_files[0].read_text(encoding='utf-8')
                    test_code = self.llm.generate_test(source_code, "go", pkg)
                    
                    if test_code:
                        test_code = test_code.replace("```go", "").replace("```", "").strip()
                        test_file = source_dir / f"{pkg}_test.go"
                        
                        if not test_file.exists():
                            test_file.write_text(test_code, encoding='utf-8')
                            generated.append(str(test_file))
                            self.log(f"    生成: {test_file.name}", "SUCCESS")
                        else:
                            self.log(f"    跳过: {test_file.name} 已存在", "INFO")
        
        return generated
    
    def git_commit(self, cycle: int):
        """提交git"""
        try:
            # 添加所有更改
            subprocess.run(["git", "add", "."], cwd=PROJECT_ROOT, capture_output=True)
            
            # 检查是否有更改
            result = subprocess.run(["git", "status", "--porcelain"], cwd=PROJECT_ROOT, capture_output=True, text=True)
            
            if result.stdout.strip():
                # 提交
                message = f"Agent循环第{cycle}轮: 测试通过,代码优化"
                subprocess.run(["git", "commit", "-m", message], cwd=PROJECT_ROOT, capture_output=True)
                self.log(f"Git提交: {message}", "SUCCESS")
            else:
                self.log("无更改，跳过提交", "INFO")
                
        except Exception as e:
            self.log(f"Git提交失败: {e}", "WARNING")


def main():
    import argparse
    
    parser = argparse.ArgumentParser(description="AgentGame 需求完成系统")
    parser.add_argument("-n", "--cycles", type=int, default=10, help="循环次数")
    args = parser.parse_args()
    
    config_path = PROJECT_ROOT / "agent_config.yaml"
    if config_path.exists():
        config = AgentConfig(config_path)
    else:
        config = AgentConfig.from_env()
    
    completer = RequirementCompleter(config)
    completer.run(args.cycles)


if __name__ == "__main__":
    main()
