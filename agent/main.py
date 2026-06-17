"""AgentGame 自循环Agent - 主入口"""
import argparse
import sys
from pathlib import Path

# 添加agent目录到路径
agent_dir = Path(__file__).parent
sys.path.insert(0, str(agent_dir.parent))

from agent.orchestrator import Orchestrator
from agent.executors.agent_executor import AgentExecutor
from agent.config import PROJECT_ROOT
from agent.config_manager import AgentConfig


def main():
    parser = argparse.ArgumentParser(
        description="AgentGame 自循环Agent - 自动分析、测试、改进游戏框架",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  python main.py                        # 运行单次迭代
  python main.py -n 5                   # 运行5次迭代
  python main.py --config agent.yaml    # 使用配置文件
  python main.py --enable-llm           # 启用LLM自动执行
  python main.py --report-only          # 只生成报告
        """
    )
    
    parser.add_argument(
        "-n", "--iterations",
        type=int,
        default=1,
        help="最大迭代次数 (默认: 1)"
    )
    
    parser.add_argument(
        "-v", "--verbose",
        action="store_true",
        default=True,
        help="详细输出"
    )
    
    parser.add_argument(
        "-q", "--quiet",
        action="store_true",
        help="安静模式"
    )
    
    parser.add_argument(
        "-c", "--config",
        type=str,
        help="配置文件路径 (YAML或JSON)"
    )
    
    parser.add_argument(
        "--enable-llm",
        action="store_true",
        help="启用LLM自动执行任务"
    )
    
    parser.add_argument(
        "--api-key",
        type=str,
        help="LLM API密钥"
    )
    
    parser.add_argument(
        "--api-url",
        type=str,
        help="LLM API地址"
    )
    
    parser.add_argument(
        "--model",
        type=str,
        help="LLM模型名称"
    )
    
    parser.add_argument(
        "--execute-tasks",
        action="store_true",
        help="执行生成的任务"
    )
    
    parser.add_argument(
        "--report-only",
        action="store_true",
        help="只生成报告，不执行任务"
    )
    
    parser.add_argument(
        "--show-history",
        action="store_true",
        help="显示历史记录"
    )
    
    parser.add_argument(
        "--generate-config",
        action="store_true",
        help="生成默认配置文件"
    )
    
    args = parser.parse_args()
    
    # 设置verbose
    verbose = args.verbose and not args.quiet
    
    # 加载配置
    config = None
    if args.config:
        config_path = Path(args.config)
        if config_path.exists():
            config = AgentConfig(config_path)
        else:
            print(f"配置文件不存在: {config_path}")
            return 1
    else:
        # 尝试加载默认配置
        default_config = PROJECT_ROOT / "agent_config.yaml"
        if default_config.exists():
            config = AgentConfig(default_config)
        else:
            config = AgentConfig.from_env()
    
    # 命令行参数覆盖配置
    if args.enable_llm:
        config.set("llm.enabled", True)
    if args.api_key:
        config.set("llm.api_key", args.api_key)
    if args.api_url:
        config.set("llm.api_url", args.api_url)
    if args.model:
        config.set("llm.model", args.model)
    if args.execute_tasks:
        config.set("tasks.auto_execute", True)
    
    # 生成配置文件
    if args.generate_config:
        config.save(PROJECT_ROOT / "agent_config.yaml")
        print("配置文件已生成: agent_config.yaml")
        return 0
    
    print("=" * 60)
    print("  AgentGame 自循环Agent")
    print("  目标：不断完善系统成为现代化AgentNpc游戏框架")
    print("=" * 60)
    print()
    
    # 显示历史
    if args.show_history:
        from agent.utils.history_tracker import HistoryTracker
        from agent.config import HISTORY_DIR
        
        tracker = HistoryTracker(HISTORY_DIR)
        history = tracker.get_history(10)
        trend = tracker.get_trend()
        
        print("=== 历史记录 ===")
        for h in history:
            print(f"  迭代 {h['iteration']}: "
                  f"通过={h['metrics']['test_passed']}, "
                  f"失败={h['metrics']['test_failed']}, "
                  f"关键问题={h['metrics']['critical_issues']}")
        
        print("\n=== 趋势 ===")
        score_trend = trend.get("overall_score", {})
        print(f"  总体评分: {score_trend.get('current', 0):.1f} "
              f"({score_trend.get('direction', 'stable')})")
        
        return 0
    
    # 创建调度器
    orchestrator = Orchestrator(
        max_iterations=args.iterations,
        verbose=verbose,
        config=config
    )
    
    # 运行迭代
    history = orchestrator.run()
    
    # 执行任务（如果启用但未自动执行）
    if args.execute_tasks and not config.get("tasks.auto_execute") and history:
        print("\n" + "=" * 60)
        print("  执行生成的任务")
        print("=" * 60)
        
        executor = AgentExecutor(PROJECT_ROOT)
        
        # 获取最后一轮的任务
        last_result = history[-1]
        tasks = last_result.get("phases", {}).get("tasks", [])
        
        if tasks:
            print(f"\n找到 {len(tasks)} 个任务待执行")
            results = executor.execute_tasks(tasks)
            
            # 统计结果
            success_count = sum(1 for r in results if r.get("success", False))
            print(f"\n执行完成: {success_count}/{len(tasks)} 成功")
        else:
            print("\n没有待执行的任务")
    
    # 打印摘要
    print("\n" + "=" * 60)
    print("  运行摘要")
    print("=" * 60)
    
    summary = orchestrator.get_summary()
    print(f"  总迭代次数: {summary.get('total_iterations', 0)}")
    print(f"  关键问题: {summary.get('critical_issues', 0)}")
    print(f"  改进项: {summary.get('improvements', 0)}")
    print(f"  待办任务: {summary.get('pending_tasks', 0)}")
    print(f"  测试通过: {summary.get('test_passed', 0)}")
    print(f"  测试失败: {summary.get('test_failed', 0)}")
    print(f"  LLM状态: {'已启用' if summary.get('llm_enabled') else '未启用'}")
    
    if summary.get('last_report'):
        print(f"\n  最新报告: {summary['last_report']}")
    
    # 显示趋势
    trend = summary.get("trend", {})
    if trend and trend.get("trend") != "insufficient_data":
        score_trend = trend.get("overall_score", {})
        print(f"\n  趋势: {score_trend.get('direction', 'stable')}")
        print(f"  评分变化: {score_trend.get('previous', 0):.1f} -> {score_trend.get('current', 0):.1f}")
    
    print()
    
    # 返回退出码
    if summary.get('test_failed', 0) > 0 or summary.get('critical_issues', 0) > 0:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
