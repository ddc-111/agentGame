import sys
import json
sys.path.insert(0, '.')
from agent.config import PROJECT_ROOT

tasks_dir = PROJECT_ROOT / 'agent' / 'tasks'
task_files = sorted(tasks_dir.glob('tasks_*.json'), reverse=True)
latest_file = task_files[0]

with open(latest_file, 'r', encoding='utf-8') as f:
    data = json.load(f)

tasks = data.get('tasks', [])

# 统计
from collections import Counter
type_priority = Counter((t.get('type'), t.get('priority')) for t in tasks)

print('Task types and priorities:')
for (t, p), count in sorted(type_priority.items()):
    print(f'  {t} ({p}): {count}')

# 显示add_feature任务
feature_tasks = [t for t in tasks if t.get('type') == 'add_feature']
print(f'\nTotal add_feature tasks: {len(feature_tasks)}')
for t in feature_tasks[:10]:
    priority = t.get('priority')
    title = t.get('title')[:60]
    print(f'  - [{priority}] {title}')
