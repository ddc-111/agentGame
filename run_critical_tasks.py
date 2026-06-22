import sys
import os
import json
from datetime import datetime

# 设置编码
os.environ["PYTHONIOENCODING"] = "utf-8"

sys.path.insert(0, '.')
from agent.executors.task_executor import TaskExecutor
from agent.utils.llm_client import LLMClient
from agent.config_manager import AgentConfig
from agent.config import PROJECT_ROOT, REPORTS_DIR

def run_critical_tasks():
    """只执行关键任务"""
    config = AgentConfig(PROJECT_ROOT / 'agent_config.yaml')
    
    # 初始化LLM
    llm = LLMClient(
        api_url=config.get('llm.api_url'),
        api_key=config.get('llm.api_key'),
        model=config.get('llm.model')
    )
    
    # 初始化任务执行器
    executor = TaskExecutor(PROJECT_ROOT, llm)
    
    # 加载任务
    tasks = load_tasks()
    
    # 只执行critical和high优先级的测试修复任务
    critical_tasks = [t for t in tasks if t.get('priority') in ['critical'] and t.get('type') == 'fix_test']
    high_tasks = [t for t in tasks if t.get('priority') == 'high' and t.get('type') in ['fix_test', 'add_test']]
    
    # 限制数量
    tasks_to_execute = critical_tasks[:10] + high_tasks[:20]
    
    print("=" * 60)
    print("  执行关键任务")
    print("=" * 60)
    print(f"  Critical任务: {len(critical_tasks[:10])}")
    print(f"  High任务: {len(high_tasks[:20])}")
    print(f"  总执行: {len(tasks_to_execute)}")
    print()
    
    # 执行任务
    results = []
    for i, task in enumerate(tasks_to_execute):
        print(f"[{i+1}/{len(tasks_to_execute)}] 执行: {task.get('title', '')[:50]}...")
        result = executor.execute_with_verification(task)
        results.append(result)
        
        if result.get('success'):
            print(f"  ✓ 成功")
        else:
            error = result.get('error') or ''
            print(f"  ✗ 失败: {error[:50]}")
    
    # 生成报告
    generate_report(results, tasks_to_execute)
    
    return results

def load_tasks():
    """加载任务文件"""
    tasks_dir = PROJECT_ROOT / 'agent' / 'tasks'
    
    # 找到最新的任务文件
    task_files = sorted(tasks_dir.glob('tasks_*.json'), reverse=True)
    if not task_files:
        print("未找到任务文件")
        return []
    
    latest_file = task_files[0]
    print(f"加载任务文件: {latest_file.name}")
    
    with open(latest_file, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    return data.get('tasks', [])

def generate_report(results, tasks):
    """生成报告"""
    success_count = sum(1 for r in results if r.get('success'))
    failed_count = len(results) - success_count
    
    report = f"""# 关键任务执行报告

生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

## 执行概要

- 总任务数: {len(results)}
- 成功: {success_count}
- 失败: {failed_count}
- 成功率: {success_count/len(results)*100 if results else 0:.1f}%

## 成功执行的任务

"""
    
    for i, result in enumerate(results):
        if result.get('success'):
            task = result.get('task', {})
            report += f"{i+1}. {task.get('title', 'Unknown')}\n"
            if result.get('file'):
                report += f"   - 文件: {result.get('file')}\n"
    
    report += "\n## 失败的任务\n\n"
    
    for i, result in enumerate(results):
        if not result.get('success'):
            task = result.get('task', {})
            report += f"{i+1}. {task.get('title', 'Unknown')}\n"
            if result.get('error'):
                report += f"   - 错误: {result.get('error')}\n"
    
    # 保存报告
    report_file = REPORTS_DIR / f"critical_tasks_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
    report_file.parent.mkdir(parents=True, exist_ok=True)
    report_file.write_text(report, encoding='utf-8')
    
    print("\n" + "=" * 60)
    print(f"  报告已保存: {report_file}")
    print("=" * 60)

if __name__ == "__main__":
    run_critical_tasks()
