import sys
import os
import json
from datetime import datetime

# 设置编码
os.environ["PYTHONIOENCODING"] = "utf-8"

sys.path.insert(0, '.')
from agent.orchestrator import Orchestrator
from agent.config_manager import AgentConfig
from agent.config import PROJECT_ROOT, REPORTS_DIR

def run_agent():
    """运行agent执行所有任务"""
    config = AgentConfig(PROJECT_ROOT / 'agent_config.yaml')
    config.set('tasks.auto_execute', True)
    config.set('tasks.priority_filter', ['critical', 'high', 'medium'])
    
    # 创建orchestrator
    orch = Orchestrator(max_iterations=10, verbose=True, config=config)
    
    print("=" * 60)
    print("  AgentGame Agent - 执行所有需求")
    print("=" * 60)
    print()
    
    # 运行
    history = orch.run()
    
    # 生成最终报告
    generate_final_report(history)
    
    return history

def generate_final_report(history):
    """生成最终报告"""
    if not history:
        print("没有历史记录")
        return
    
    # 统计信息
    total_iterations = len(history)
    last_result = history[-1]
    
    # 收集所有执行的任务
    all_executed = []
    all_success = []
    all_failed = []
    
    for result in history:
        exec_results = result.get('phases', {}).get('task_execution', {})
        all_executed.extend(exec_results.get('details', []))
        all_success.extend([d for d in exec_results.get('details', []) if d.get('success')])
        all_failed.extend([d for d in exec_results.get('details', []) if not d.get('success')])
    
    # 生成报告
    report = f"""# AgentGame Agent 执行报告

生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

## 执行概要

- 总迭代次数: {total_iterations}
- 总执行任务数: {len(all_executed)}
- 成功任务数: {len(all_success)}
- 失败任务数: {len(all_failed)}
- 成功率: {len(all_success)/len(all_executed)*100 if all_executed else 0:.1f}%

## 最后一轮状态

- 测试通过: {last_result.get('phases', {}).get('tests', {}).get('passed', 0)}
- 测试失败: {last_result.get('phases', {}).get('tests', {}).get('failed', 0)}
- 关键问题: {len(last_result.get('phases', {}).get('gaps', {}).get('critical', []))}
- 改进点: {len(last_result.get('phases', {}).get('gaps', {}).get('improvements', []))}

## 成功执行的任务

"""
    
    # 列出成功执行的任务
    for i, task in enumerate(all_success[:50]):  # 最多显示50个
        task_info = task.get('task', {})
        report += f"{i+1}. {task_info.get('title', 'Unknown')}\n"
        if task.get('file'):
            report += f"   - 文件: {task.get('file')}\n"
    
    if len(all_success) > 50:
        report += f"\n... 还有 {len(all_success) - 50} 个成功任务\n"
    
    report += "\n## 失败的任务\n\n"
    
    # 列出失败的任务
    for i, task in enumerate(all_failed[:20]):  # 最多显示20个
        task_info = task.get('task', {})
        report += f"{i+1}. {task_info.get('title', 'Unknown')}\n"
        if task.get('error'):
            report += f"   - 错误: {task.get('error')}\n"
    
    if len(all_failed) > 20:
        report += f"\n... 还有 {len(all_failed) - 20} 个失败任务\n"
    
    report += f"""
## 历史迭代

| 迭代 | 测试通过 | 测试失败 | 关键问题 | 改进点 |
|------|----------|----------|----------|--------|
"""
    
    for result in history:
        tests = result.get('phases', {}).get('tests', {})
        gaps = result.get('phases', {}).get('gaps', {})
        report += f"| {result.get('iteration')} | {tests.get('passed', 0)} | {tests.get('failed', 0)} | {len(gaps.get('critical', []))} | {len(gaps.get('improvements', []))} |\n"
    
    report += """
## 建议后续工作

1. 修复所有失败的测试
2. 继续执行未完成的任务
3. 定期运行agent进行代码质量检查
"""
    
    # 保存报告
    report_file = REPORTS_DIR / f"final_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
    report_file.parent.mkdir(parents=True, exist_ok=True)
    report_file.write_text(report, encoding='utf-8')
    
    print("\n" + "=" * 60)
    print(f"  最终报告已保存: {report_file}")
    print("=" * 60)
    print()

if __name__ == "__main__":
    run_agent()
