import sys
sys.path.insert(0, '.')
from agent.analyzers.code_analyzer import CodeAnalyzer
from agent.analyzers.test_analyzer import TestAnalyzer
from agent.executors.build_executor import BuildExecutor
from agent.config import PROJECT_ROOT

print('=' * 60)
print('  框架现状检查')
print('=' * 60)

# 代码分析
print('\n[1] 代码分析...')
code_analyzer = CodeAnalyzer(PROJECT_ROOT)
code = code_analyzer.analyze()
print(f'  总文件数: {code.get("total_files")}')
print(f'  总代码行: {code.get("total_lines")}')
print(f'  Server: {code.get("server", {}).get("lines", 0)} 行')
print(f'  Client: {code.get("client", {}).get("lines", 0)} 行')
print(f'  GM: {code.get("gm", {}).get("lines", 0)} 行')

# 测试分析
print('\n[2] 测试分析...')
test_analyzer = TestAnalyzer(PROJECT_ROOT)
tests = test_analyzer.analyze()
print(f'  测试文件: {tests.get("total_test_files")}')
print(f'  测试用例: {tests.get("total_test_cases")}')

# 构建检查
print('\n[3] 构建检查...')
build_executor = BuildExecutor(PROJECT_ROOT)
builds = build_executor.build_all()
for target, success in builds.items():
    status = '✓ 通过' if success else '✗ 失败'
    print(f'  {target}: {status}')

print('\n' + '=' * 60)
