"""完整的需求-开发-测试循环"""
import sys
import time
from pathlib import Path
from datetime import datetime

# 添加路径
agent_dir = Path(__file__).parent
sys.path.insert(0, str(agent_dir.parent))

from agent.config import PROJECT_ROOT, REPORTS_DIR, REQUIREMENTS_DIR
from agent.utils.llm_client import LLMClient
from agent.config_manager import AgentConfig
from agent.agents.requirement_agent import RequirementAgent
from agent.analyzers.code_analyzer import CodeAnalyzer
from agent.analyzers.test_analyzer import TestAnalyzer
from agent.executors.test_executor import TestExecutor
from agent.executors.build_executor import BuildExecutor
from agent.executors.task_executor import TaskExecutor
from agent.generators.report_generator import ReportGenerator


class FullCycleRunner:
    """完整循环执行器"""
    
    def __init__(self, config: AgentConfig = None):
        self.config = config or AgentConfig.from_env()
        
        # 初始化组件
        self.llm = LLMClient(
            api_url=self.config.get("llm.api_url"),
            api_key=self.config.get("llm.api_key"),
            model=self.config.get("llm.model")
        )
        
        self.requirement_agent = RequirementAgent(self.llm)
        self.code_analyzer = CodeAnalyzer(PROJECT_ROOT)
        self.test_analyzer = TestAnalyzer(PROJECT_ROOT)
        self.test_executor = TestExecutor(PROJECT_ROOT)
        self.build_executor = BuildExecutor(PROJECT_ROOT)
        self.task_executor = TaskExecutor(PROJECT_ROOT, self.llm)
        self.report_generator = ReportGenerator(REPORTS_DIR)
        
        self.results = []
    
    def log(self, message: str, level: str = "INFO"):
        """日志输出"""
        timestamp = datetime.now().strftime("%H:%M:%S")
        prefix = {"INFO": "[i]", "SUCCESS": "[+]", "WARNING": "[!]", "ERROR": "[-]", "PHASE": "[*]"}
        try:
            print(f"[{timestamp}] {prefix.get(level, '[.]')} {message}")
        except UnicodeEncodeError:
            print(f"[{timestamp}] {prefix.get(level, '[.]')} {message.encode('ascii', 'replace').decode()}")
    
    def run_full_cycle(self, cycles: int = 3):
        """运行完整循环"""
        self.log("=" * 60, "PHASE")
        self.log("AgentGame 自循环系统 - 需求/开发/测试循环", "PHASE")
        self.log("=" * 60, "PHASE")
        
        for cycle in range(1, cycles + 1):
            self.log(f"\n{'='*60}", "PHASE")
            self.log(f"第 {cycle}/{cycles} 轮循环", "PHASE")
            self.log(f"{'='*60}", "PHASE")
            
            cycle_result = {
                "cycle": cycle,
                "timestamp": datetime.now().isoformat(),
                "phases": {}
            }
            
            # 阶段1: 需求生成
            self.log("\n[阶段1] 需求生成", "PHASE")
            requirements = self.phase_requirements()
            cycle_result["phases"]["requirements"] = requirements
            
            # 阶段2: 开发实现
            self.log("\n[阶段2] 开发实现", "PHASE")
            development = self.phase_development(requirements)
            cycle_result["phases"]["development"] = development
            
            # 阶段3: 测试验证
            self.log("\n[阶段3] 测试验证", "PHASE")
            testing = self.phase_testing()
            cycle_result["phases"]["testing"] = testing
            
            # 阶段4: 报告生成
            self.log("\n[阶段4] 报告生成", "PHASE")
            report = self.phase_report(cycle_result)
            cycle_result["report"] = report
            
            self.results.append(cycle_result)
            
            # 短暂延迟
            if cycle < cycles:
                self.log(f"\n等待5秒后开始下一轮循环...", "INFO")
                time.sleep(5)
        
        # 生成最终报告
        self.generate_final_report()
    
    def phase_requirements(self) -> dict:
        """阶段1: 需求生成"""
        self.log("分析代码库结构...", "INFO")
        code_analysis = self.code_analyzer.analyze()
        
        self.log("分析测试覆盖...", "INFO")
        test_analysis = self.test_analyzer.analyze()
        
        self.log("生成新需求...", "INFO")
        current_features = self._extract_features(code_analysis)
        
        requirements = self.requirement_agent.generate_requirements(
            code_analysis, test_analysis, current_features
        )
        
        # 保存需求
        req_file = REQUIREMENTS_DIR / f"req_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        self.requirement_agent.save_requirements(req_file)
        
        summary = self.requirement_agent.get_requirements_summary()
        self.log(f"生成 {summary['total']} 个需求", "SUCCESS")
        self.log(f"  高优先级: {summary['by_priority']['high']}", "INFO")
        self.log(f"  中优先级: {summary['by_priority']['medium']}", "INFO")
        
        return {
            "total": summary['total'],
            "high_priority": summary['by_priority']['high'],
            "medium_priority": summary['by_priority']['medium'],
            "top_requirements": [r["title"] for r in summary["top_requirements"][:5]],
            "file": str(req_file)
        }
    
    def phase_development(self, requirements: dict) -> dict:
        """阶段2: 开发实现"""
        self.log("执行构建验证...", "INFO")
        build_result = self.build_executor.build_all()
        
        success_count = sum(1 for v in build_result.values() if v)
        self.log(f"构建结果: {success_count}/3 通过", "SUCCESS" if success_count == 3 else "WARNING")
        
        # 使用LLM生成测试代码
        generated_tests = []
        if self.llm.is_configured():
            self.log("使用LLM生成测试代码...", "INFO")
            generated_tests = self._generate_tests_with_llm()
        
        return {
            "build": build_result,
            "generated_tests": len(generated_tests)
        }
    
    def phase_testing(self) -> dict:
        """阶段3: 测试验证"""
        self.log("运行完整测试套件...", "INFO")
        test_result = self.test_executor.run_all_tests()
        
        self.log(f"测试结果:", "SUCCESS")
        self.log(f"  通过: {test_result['passed']}", "INFO")
        self.log(f"  失败: {test_result['failed']}", "INFO")
        self.log(f"  跳过: {test_result['skipped']}", "INFO")
        self.log(f"  耗时: {test_result['duration']:.1f}秒", "INFO")
        
        return {
            "passed": test_result['passed'],
            "failed": test_result['failed'],
            "skipped": test_result['skipped'],
            "duration": test_result['duration'],
            "pass_rate": (test_result['passed'] / (test_result['passed'] + test_result['failed']) * 100) if (test_result['passed'] + test_result['failed']) > 0 else 0
        }
    
    def phase_report(self, cycle_result: dict) -> str:
        """阶段4: 报告生成"""
        self.log("生成报告...", "INFO")
        report_path = self.report_generator.generate(cycle_result)
        self.log(f"报告已保存: {report_path}", "SUCCESS")
        return str(report_path)
    
    def _extract_features(self, code_analysis: dict) -> list:
        """提取当前特性"""
        features = []
        
        for pkg in code_analysis.get("server", {}).get("packages", []):
            features.append(f"server/{pkg}")
        
        for model in code_analysis.get("server", {}).get("models", [])[:10]:
            features.append(model)
        
        for scene in code_analysis.get("client", {}).get("scenes", []):
            features.append(scene)
        
        return features
    
    def _generate_tests_with_llm(self) -> list:
        """使用LLM生成测试"""
        generated = []
        
        # 获取需要测试的模块
        test_analysis = self.test_analyzer.analyze()
        untested = test_analysis.get("server", {}).get("untested_packages", [])
        
        for pkg in untested[:2]:  # 每次只处理2个
            self.log(f"  为 {pkg} 生成测试...", "INFO")
            
            # 获取源代码
            source_dir = PROJECT_ROOT / "server" / "internal" / pkg
            if source_dir.exists():
                source_files = list(source_dir.glob("*.go"))
                if source_files:
                    source_code = source_files[0].read_text(encoding='utf-8')
                    test_code = self.llm.generate_test(source_code, "go", pkg)
                    
                    if test_code:
                        # 清理markdown标记
                        test_code = test_code.replace("```go", "").replace("```", "").strip()
                        
                        test_file = source_dir / f"{pkg}_test.go"
                        test_file.write_text(test_code, encoding='utf-8')
                        generated.append(str(test_file))
                        self.log(f"    生成: {test_file.name}", "SUCCESS")
        
        return generated
    
    def generate_final_report(self):
        """生成最终报告"""
        self.log("\n" + "=" * 60, "PHASE")
        self.log("生成最终报告", "PHASE")
        self.log("=" * 60, "PHASE")
        
        # 汇总数据
        total_requirements = sum(r["phases"]["requirements"]["total"] for r in self.results)
        total_tests_passed = max(r["phases"]["testing"]["passed"] for r in self.results)
        total_generated_tests = sum(r["phases"]["development"]["generated_tests"] for r in self.results)
        
        final_report = f"""# AgentGame 自循环系统 - 最终报告

**生成时间**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
**循环次数**: {len(self.results)}

---

## 执行摘要

| 指标 | 结果 |
|------|------|
| 总循环次数 | {len(self.results)} |
| 生成需求数 | {total_requirements} |
| 生成测试文件 | {total_generated_tests} |
| 最终测试通过 | {total_tests_passed} |

---

## 各轮循环详情

"""
        for r in self.results:
            req = r["phases"]["requirements"]
            test = r["phases"]["testing"]
            final_report += f"""### 第 {r['cycle']} 轮

- **需求**: {req['total']}个 (高优先级: {req['high_priority']})
- **测试通过**: {test['passed']}个
- **测试失败**: {test['failed']}个
- **通过率**: {test['pass_rate']:.1f}%

"""
        
        final_report += """---

## 生成的需求分类

### 核心功能
- 完整的NPC AI对话系统
- 多Agent协作机制
- 动态任务生成与执行
- 玩家行为分析与个性化
- 实时状态同步

### AI集成
- 多LLM支持 (OpenAI/Claude/本地模型)
- Prompt模板管理
- 上下文窗口优化
- 流式响应
- AI决策链

### 现代化架构
- 微服务架构
- 容器化部署
- CI/CD自动化
- 监控与日志
- 自动扩缩容

### 工作流支持
- 可视化工作流编辑器
- 条件分支与循环
- 事件触发机制
- 定时任务
- 异步任务队列

### 竞争力特性
- 高性能（低延迟）
- 可扩展性
- 安全性
- 多语言支持
- 插件系统

---

## 下一步建议

1. **优先实现高优先级需求** - AI对话系统、多Agent协作
2. **完善测试覆盖** - 目标80%覆盖率
3. **优化LLM集成** - 支持更多模型、降低延迟
4. **建立CI/CD流程** - 自动化构建、测试、部署
5. **编写开发者文档** - API文档、SDK、示例代码

---

*报告由AgentGame自循环系统自动生成*
"""
        
        # 保存报告
        report_file = REPORTS_DIR / f"final_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
        report_file.write_text(final_report, encoding='utf-8')
        
        self.log(f"最终报告已保存: {report_file}", "SUCCESS")
        
        # 打印摘要
        self.log("\n" + "=" * 60, "SUCCESS")
        self.log("执行完成!", "SUCCESS")
        self.log(f"  总循环: {len(self.results)}", "INFO")
        self.log(f"  生成需求: {total_requirements}", "INFO")
        self.log(f"  测试通过: {total_tests_passed}", "INFO")
        self.log(f"  报告: {report_file}", "INFO")
        self.log("=" * 60, "SUCCESS")


def main():
    """主函数"""
    import argparse
    
    parser = argparse.ArgumentParser(description="AgentGame 完整循环执行器")
    parser.add_argument("-n", "--cycles", type=int, default=3, help="循环次数")
    args = parser.parse_args()
    
    # 加载配置
    config_path = PROJECT_ROOT / "agent_config.yaml"
    if config_path.exists():
        config = AgentConfig(config_path)
    else:
        config = AgentConfig.from_env()
    
    # 运行循环
    runner = FullCycleRunner(config)
    runner.run_full_cycle(args.cycles)


if __name__ == "__main__":
    main()
