import sys
import os

# 设置编码
os.environ["PYTHONIOENCODING"] = "utf-8"

sys.path.insert(0, '.')
from agent.orchestrator import Orchestrator
from agent.config_manager import AgentConfig
from agent.config import PROJECT_ROOT

config = AgentConfig(PROJECT_ROOT / 'agent_config.yaml')
config.set('tasks.auto_execute', False)  # 不执行任务

orch = Orchestrator(max_iterations=1, verbose=True, config=config)
print('Starting cycle...')
result = orch.run_cycle()
print('Cycle complete!')

# 保存任务到文件
orch._save_tasks(result.get('phases', {}).get('tasks', []))

# 显示任务
tasks = result.get('phases', {}).get('tasks', [])
print(f'Total tasks: {len(tasks)}')

# 按优先级统计
from collections import Counter
priority_count = Counter(t.get('priority') for t in tasks)
print('By priority:', dict(priority_count))

# 按类型统计
type_count = Counter(t.get('type') for t in tasks)
print('By type:', dict(type_count))

# 显示前5个任务
print('\nFirst 5 tasks:')
for i, t in enumerate(tasks[:5]):
    title = t.get('title', '')[:60]
    priority = t.get('priority')
    print(f'{i+1}. [{priority}] {title}')
