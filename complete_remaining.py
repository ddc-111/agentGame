import sys
import os
import json
from datetime import datetime

os.environ["PYTHONIOENCODING"] = "utf-8"

sys.path.insert(0, '.')
from agent.executors.task_executor import TaskExecutor
from agent.utils.llm_client import LLMClient
from agent.config_manager import AgentConfig
from agent.config import PROJECT_ROOT, REPORTS_DIR

def complete_remaining():
    """完成剩余任务"""
    config = AgentConfig(PROJECT_ROOT / 'agent_config.yaml')
    
    llm = LLMClient(
        api_url=config.get('llm.api_url'),
        api_key=config.get('llm.api_key'),
        model=config.get('llm.model')
    )
    
    executor = TaskExecutor(PROJECT_ROOT, llm)
    
    # 加载任务
    tasks = load_tasks()
    
    # 获取已执行的任务（从之前的报告中）
    executed_titles = get_executed_titles()
    
    # 过滤未执行的任务
    remaining_tasks = [t for t in tasks if t.get('title') not in executed_titles]
    
    # 按优先级排序
    priority_order = {'critical': 0, 'high': 1, 'medium': 2, 'low': 3}
    remaining_tasks.sort(key=lambda t: priority_order.get(t.get('priority'), 4))
    
    print("=" * 60)
    print("  完成剩余任务")
    print("=" * 60)
    print(f"  剩余任务数: {len(remaining_tasks)}")
    print()
    
    # 执行任务
    results = []
    for i, task in enumerate(remaining_tasks[:15]):  # 执行前15个
        task_type = task.get('type')
        title = task.get('title', '')[:50]
        print(f"[{i+1}/{min(15, len(remaining_tasks))}] 执行: {title}...")
        
        try:
            if task_type == 'fix_test':
                exec_result = executor._execute_fix_test(task)
            elif task_type == 'add_test':
                exec_result = executor._execute_add_test(task)
            elif task_type == 'add_feature':
                exec_result = executor._execute_add_feature(task)
            elif task_type == 'refactor':
                exec_result = executor._execute_refactor(task)
            else:
                exec_result = executor._execute_generic_task(task)
            
            results.append({
                "task": task,
                "success": exec_result.get("success", False),
                "error": exec_result.get("error"),
                "file": exec_result.get("file")
            })
            
            if exec_result.get('success'):
                print(f"  ✓ 成功")
            else:
                error = exec_result.get('error') or ''
                print(f"  ✗ 失败: {error[:40]}")
        except Exception as e:
            print(f"  ✗ 异常: {str(e)[:40]}")
            results.append({
                "task": task,
                "success": False,
                "error": str(e)
            })
    
    # 生成报告
    generate_report(results)
    
    return results

def load_tasks():
    """加载任务"""
    tasks_dir = PROJECT_ROOT / 'agent' / 'tasks'
    task_files = sorted(tasks_dir.glob('tasks_*.json'), reverse=True)
    
    if not task_files:
        return []
    
    # 找到包含最多任务的文件
    best_file = task_files[0]
    best_count = 0
    
    for f in task_files[:3]:
        with open(f, 'r', encoding='utf-8') as fh:
            data = json.load(fh)
        count = len(data.get('tasks', []))
        if count > best_count:
            best_count = count
            best_file = f
    
    with open(best_file, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    return data.get('tasks', [])

def get_executed_titles():
    """获取已执行的任务标题"""
    executed = set()
    
    # 从之前的报告中提取
    reports_dir = PROJECT_ROOT / 'agent' / 'reports'
    for report_file in reports_dir.glob('*.md'):
        try:
            content = report_file.read_text(encoding='utf-8')
            # 简单提取已执行的任务
            if '✓' in content:
                lines = content.split('\n')
                for line in lines:
                    if '✓' in line and '-' in line:
                        # 提取任务标题
                        parts = line.split('-', 1)
                        if len(parts) > 1:
                            title = parts[1].strip()
                            if title:
                                executed.add(title)
        except:
            pass
    
    return executed

def generate_report(results):
    """生成报告"""
    success_count = sum(1 for r in results if r.get('success'))
    failed_count = len(results) - success_count
    
    report = f"""# 剩余任务执行报告

生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}

## 执行概要

- 总任务数: {len(results)}
- 成功: {success_count}
- 失败: {failed_count}
- 成功率: {success_count/len(results)*100 if results else 0:.1f}%

## 执行详情

"""
    
    for i, result in enumerate(results):
        task = result.get('task', {})
        status = "✓" if result.get('success') else "✗"
        report += f"{i+1}. {status} {task.get('title', 'Unknown')}\n"
        if result.get('file'):
            report += f"   - 文件: {result.get('file')}\n"
        if result.get('error'):
            report += f"   - 错误: {result.get('error')}\n"
    
    report_file = REPORTS_DIR / f"remaining_tasks_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
    report_file.parent.mkdir(parents=True, exist_ok=True)
    report_file.write_text(report, encoding='utf-8')
    
    print("\n" + "=" * 60)
    print(f"  报告已保存: {report_file}")
    print("=" * 60)

if __name__ == "__main__":
    complete_remaining()
