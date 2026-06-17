"""报告生成器 - 生成分析和测试报告"""
import json
from pathlib import Path
from typing import Dict
from datetime import datetime


class ReportGenerator:
    """生成报告"""
    
    def __init__(self, reports_dir: Path):
        self.reports_dir = reports_dir
        self.reports_dir.mkdir(parents=True, exist_ok=True)
    
    def generate(self, cycle_result: Dict) -> Path:
        """生成完整的周期报告"""
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        
        # 生成JSON报告
        json_path = self.reports_dir / f"report_{timestamp}.json"
        with open(json_path, 'w', encoding='utf-8') as f:
            json.dump(cycle_result, f, ensure_ascii=False, indent=2)
        
        # 生成Markdown报告
        md_path = self.reports_dir / f"report_{timestamp}.md"
        md_content = self._generate_markdown(cycle_result)
        with open(md_path, 'w', encoding='utf-8') as f:
            f.write(md_content)
        
        # 生成摘要
        summary_path = self.reports_dir / f"summary_{timestamp}.txt"
        summary = self._generate_summary(cycle_result)
        with open(summary_path, 'w', encoding='utf-8') as f:
            f.write(summary)
        
        return md_path
    
    def _generate_markdown(self, result: Dict) -> str:
        """生成Markdown格式报告"""
        iteration = result.get("iteration", 0)
        timestamp = result.get("timestamp", "")
        phases = result.get("phases", {})
        
        md = f"""# AgentGame 自循环Agent报告

## 迭代 #{iteration}
- **时间**: {timestamp}
- **耗时**: {result.get('duration', 0):.1f}秒

---

## 📊 代码分析

"""
        # 代码分析部分
        code_analysis = phases.get("code_analysis", {})
        md += f"| 指标 | 数值 |\n|------|------|\n"
        md += f"| 总文件数 | {code_analysis.get('total_files', 0)} |\n"
        md += f"| 总代码行数 | {code_analysis.get('total_lines', 0)} |\n"
        
        # 各端统计
        for target in ["server", "client", "gm"]:
            target_data = code_analysis.get(target, {})
            md += f"| {target.capitalize()} 文件 | {target_data.get('files', 0)} |\n"
            md += f"| {target.capitalize()} 行数 | {target_data.get('lines', 0)} |\n"
        
        md += f"""
---

## 🧪 测试分析

"""
        # 测试分析部分
        test_analysis = phases.get("test_analysis", {})
        md += f"| 指标 | 数值 |\n|------|------|\n"
        md += f"| 测试文件数 | {test_analysis.get('total_test_files', 0)} |\n"
        md += f"| 测试用例数 | {test_analysis.get('total_test_cases', 0)} |\n"
        
        md += f"""
### 测试覆盖差距
"""
        for gap in test_analysis.get("coverage_gaps", []):
            md += f"- **{gap.get('target')}/{gap.get('name')}**: {gap.get('priority')} 优先级\n"
        
        md += f"""
---

## 🔨 构建结果

"""
        # 构建结果
        build_result = phases.get("build", {})
        for target, success in build_result.items():
            status = "✅ 通过" if success else "❌ 失败"
            md += f"| {target.capitalize()} | {status} |\n"
        
        md += f"""
---

## 🎯 测试结果

"""
        # 测试结果
        test_result = phases.get("tests", {})
        md += f"| 指标 | 数值 |\n|------|------|\n"
        md += f"| 通过 | {test_result.get('passed', 0)} |\n"
        md += f"| 失败 | {test_result.get('failed', 0)} |\n"
        md += f"| 跳过 | {test_result.get('skipped', 0)} |\n"
        md += f"| 耗时 | {test_result.get('duration', 0):.1f}秒 |\n"
        
        # 失败详情
        failures = test_result.get("failures", [])
        if failures:
            md += f"""
### 失败详情

| 目标 | 测试 | 错误 |
|------|------|------|
"""
            for failure in failures[:10]:  # 只显示前10个
                error = failure.get("error", "")[:50] + "..." if len(failure.get("error", "")) > 50 else failure.get("error", "")
                md += f"| {failure.get('target', '')} | {failure.get('test', '')} | {error} |\n"
        
        md += f"""
---

## 🔍 差距分析

"""
        # 差距分析
        gaps = phases.get("gaps", {})
        
        # 指标
        metrics = gaps.get("metrics", {})
        md += f"### 质量指标\n\n"
        md += f"| 指标 | 分数 |\n|------|------|\n"
        md += f"| 测试通过率 | {metrics.get('test_pass_rate', 0):.1f}% |\n"
        md += f"| 测试覆盖分 | {metrics.get('test_coverage_score', 0):.1f} |\n"
        md += f"| 代码健康分 | {metrics.get('code_health_score', 0):.1f} |\n"
        md += f"| 总体评分 | {metrics.get('overall_score', 0):.1f} |\n"
        
        # 关键问题
        critical = gaps.get("critical", [])
        if critical:
            md += f"""
### ❗ 关键问题 ({len(critical)})

"""
            for issue in critical:
                md += f"- **{issue.get('type', '')}** ({issue.get('target', '')}): {issue.get('suggestion', '')}\n"
        
        # 改进项
        improvements = gaps.get("improvements", [])
        if improvements:
            md += f"""
### 📈 改进项 ({len(improvements)})

"""
            for imp in improvements[:10]:  # 只显示前10个
                md += f"- **{imp.get('type', '')}** ({imp.get('target', '')}): {imp.get('suggestion', '')}\n"
        
        md += f"""
---

## 📋 生成的任务

"""
        # 任务列表
        tasks = phases.get("tasks", [])
        if tasks:
            md += f"| ID | 类型 | 标题 | 优先级 |\n|-----|------|------|--------|\n"
            for task in tasks:
                md += f"| {task.get('id', '')[:8]} | {task.get('type', '')} | {task.get('title', '')[:30]} | {task.get('priority', '')} |\n"
        else:
            md += "无待办任务\n"
        
        md += f"""
---

## 📝 总结

"""
        # 总结
        if len(critical) == 0 and test_result.get("failed", 0) == 0:
            md += "✅ **本轮迭代结果良好**：无关键问题，所有测试通过。\n"
        elif len(critical) > 0:
            md += f"⚠️ **需要关注**：存在 {len(critical)} 个关键问题需要解决。\n"
        elif test_result.get("failed", 0) > 0:
            md += f"⚠️ **测试失败**：有 {test_result.get('failed', 0)} 个测试失败需要修复。\n"
        
        md += f"""
---

*报告生成时间: {datetime.now().isoformat()}*
"""
        
        return md
    
    def _generate_summary(self, result: Dict) -> str:
        """生成文本摘要"""
        phases = result.get("phases", {})
        test_result = phases.get("tests", {})
        gaps = phases.get("gaps", {})
        tasks = phases.get("tasks", [])
        
        summary = f"""AgentGame 自循环Agent - 迭代 #{result.get('iteration', 0)}
{'='*50}

代码统计:
  - 总文件: {phases.get('code_analysis', {}).get('total_files', 0)}
  - 总行数: {phases.get('code_analysis', {}).get('total_lines', 0)}

测试结果:
  - 通过: {test_result.get('passed', 0)}
  - 失败: {test_result.get('failed', 0)}
  - 跳过: {test_result.get('skipped', 0)}

质量指标:
  - 测试通过率: {gaps.get('metrics', {}).get('test_pass_rate', 0):.1f}%
  - 总体评分: {gaps.get('metrics', {}).get('overall_score', 0):.1f}

关键问题: {len(gaps.get('critical', []))}
待办任务: {len(tasks)}

耗时: {result.get('duration', 0):.1f}秒
"""
        
        return summary
    
    def get_latest_report(self) -> Path:
        """获取最新报告"""
        reports = sorted(self.reports_dir.glob("report_*.md"), reverse=True)
        return reports[0] if reports else None
    
    def get_history(self, limit: int = 10) -> list:
        """获取历史报告列表"""
        reports = sorted(self.reports_dir.glob("report_*.json"), reverse=True)
        return reports[:limit]
