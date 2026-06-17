"""Web报告生成器 - 生成可交互的HTML报告"""
import json
from pathlib import Path
from typing import Dict
from datetime import datetime


class WebReportGenerator:
    """生成Web格式的报告"""
    
    def __init__(self, reports_dir: Path):
        self.reports_dir = reports_dir
        self.reports_dir.mkdir(parents=True, exist_ok=True)
    
    def generate_dashboard(self, history: list, trend: Dict) -> Path:
        """生成仪表板HTML"""
        html = self._generate_dashboard_html(history, trend)
        
        output_file = self.reports_dir / "dashboard.html"
        output_file.write_text(html, encoding='utf-8')
        
        return output_file
    
    def _generate_dashboard_html(self, history: list, trend: Dict) -> str:
        """生成仪表板HTML内容"""
        # 提取数据
        iterations = [h.get("iteration", 0) for h in history]
        test_passed = [h.get("phases", {}).get("tests", {}).get("passed", 0) for h in history]
        test_failed = [h.get("phases", {}).get("tests", {}).get("failed", 0) for h in history]
        critical = [len(h.get("phases", {}).get("gaps", {}).get("critical", [])) for h in history]
        scores = [h.get("phases", {}).get("gaps", {}).get("metrics", {}).get("overall_score", 0) for h in history]
        
        # 最新数据
        latest = history[-1] if history else {}
        latest_tests = latest.get("phases", {}).get("tests", {})
        latest_gaps = latest.get("phases", {}).get("gaps", {})
        
        html = f"""<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AgentGame 自循环Agent 仪表板</title>
    <style>
        * {{
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }}
        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            color: #fff;
            min-height: 100vh;
            padding: 20px;
        }}
        .container {{
            max-width: 1400px;
            margin: 0 auto;
        }}
        .header {{
            text-align: center;
            padding: 30px 0;
        }}
        .header h1 {{
            font-size: 2.5em;
            background: linear-gradient(90deg, #00d2ff, #3a7bd5);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }}
        .header p {{
            color: #8892b0;
            margin-top: 10px;
        }}
        .grid {{
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-top: 30px;
        }}
        .card {{
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            padding: 25px;
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.1);
            transition: transform 0.3s ease;
        }}
        .card:hover {{
            transform: translateY(-5px);
        }}
        .card h2 {{
            font-size: 1.2em;
            color: #8892b0;
            margin-bottom: 15px;
        }}
        .metric {{
            font-size: 3em;
            font-weight: bold;
            background: linear-gradient(90deg, #00d2ff, #3a7bd5);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }}
        .metric.success {{
            background: linear-gradient(90deg, #00b09b, #96c93d);
            -webkit-background-clip: text;
        }}
        .metric.danger {{
            background: linear-gradient(90deg, #ff416c, #ff4b2b);
            -webkit-background-clip: text;
        }}
        .metric.warning {{
            background: linear-gradient(90deg, #f7971e, #ffd200);
            -webkit-background-clip: text;
        }}
        .chart-container {{
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            padding: 25px;
            margin-top: 20px;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }}
        .chart-container h2 {{
            color: #8892b0;
            margin-bottom: 20px;
        }}
        .bar-chart {{
            display: flex;
            align-items: flex-end;
            height: 200px;
            gap: 10px;
            padding: 20px 0;
        }}
        .bar {{
            flex: 1;
            background: linear-gradient(180deg, #00d2ff, #3a7bd5);
            border-radius: 5px 5px 0 0;
            position: relative;
            min-width: 30px;
            transition: all 0.3s ease;
        }}
        .bar:hover {{
            opacity: 0.8;
        }}
        .bar-label {{
            position: absolute;
            bottom: -25px;
            left: 50%;
            transform: translateX(-50%);
            font-size: 0.8em;
            color: #8892b0;
        }}
        .bar-value {{
            position: absolute;
            top: -25px;
            left: 50%;
            transform: translateX(-50%);
            font-size: 0.8em;
            color: #fff;
        }}
        .tasks-list {{
            max-height: 400px;
            overflow-y: auto;
        }}
        .task-item {{
            background: rgba(255, 255, 255, 0.03);
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 10px;
            border-left: 4px solid #3a7bd5;
        }}
        .task-item.critical {{
            border-left-color: #ff416c;
        }}
        .task-item.high {{
            border-left-color: #f7971e;
        }}
        .task-item.medium {{
            border-left-color: #00d2ff;
        }}
        .task-title {{
            font-weight: bold;
            margin-bottom: 5px;
        }}
        .task-meta {{
            font-size: 0.9em;
            color: #8892b0;
        }}
        .trend {{
            display: flex;
            align-items: center;
            gap: 10px;
            margin-top: 10px;
        }}
        .trend-arrow {{
            font-size: 1.5em;
        }}
        .trend-arrow.up {{
            color: #96c93d;
        }}
        .trend-arrow.down {{
            color: #ff416c;
        }}
        .trend-arrow.stable {{
            color: #8892b0;
        }}
        .footer {{
            text-align: center;
            padding: 30px 0;
            color: #8892b0;
            font-size: 0.9em;
        }}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>AgentGame 自循环Agent</h1>
            <p>自动化分析、测试、改进游戏框架</p>
        </div>
        
        <div class="grid">
            <div class="card">
                <h2>测试通过</h2>
                <div class="metric success">{latest_tests.get('passed', 0)}</div>
                <div class="trend">
                    <span class="trend-arrow {'up' if len(test_passed) > 1 and test_passed[-1] > test_passed[-2] else 'stable'}">{'↑' if len(test_passed) > 1 and test_passed[-1] > test_passed[-2] else '→'}</span>
                    <span>较上次迭代</span>
                </div>
            </div>
            
            <div class="card">
                <h2>测试失败</h2>
                <div class="metric danger">{latest_tests.get('failed', 0)}</div>
                <div class="trend">
                    <span class="trend-arrow {'down' if len(test_failed) > 1 and test_failed[-1] > test_failed[-2] else 'stable'}">{'↑' if len(test_failed) > 1 and test_failed[-1] > test_failed[-2] else '→'}</span>
                    <span>较上次迭代</span>
                </div>
            </div>
            
            <div class="card">
                <h2>关键问题</h2>
                <div class="metric warning">{len(latest_gaps.get('critical', []))}</div>
                <div class="trend">
                    <span class="trend-arrow {'down' if len(critical) > 1 and critical[-1] > critical[-2] else 'stable'}">{'↑' if len(critical) > 1 and critical[-1] > critical[-2] else '→'}</span>
                    <span>需要解决</span>
                </div>
            </div>
            
            <div class="card">
                <h2>总体评分</h2>
                <div class="metric">{latest_gaps.get('metrics', {}).get('overall_score', 0):.1f}</div>
                <div class="trend">
                    <span class="trend-arrow {'up' if len(scores) > 1 and scores[-1] > scores[-2] else 'down' if len(scores) > 1 and scores[-1] < scores[-2] else 'stable'}">{'↑' if len(scores) > 1 and scores[-1] > scores[-2] else '↓' if len(scores) > 1 and scores[-1] < scores[-2] else '→'}</span>
                    <span>{'提升中' if len(scores) > 1 and scores[-1] > scores[-2] else '下降中' if len(scores) > 1 and scores[-1] < scores[-2] else '稳定'}</span>
                </div>
            </div>
        </div>
        
        <div class="chart-container">
            <h2>测试趋势</h2>
            <div class="bar-chart">
                {self._generate_bars(iterations[-10:], test_passed[-10:], test_failed[-10:])}
            </div>
        </div>
        
        <div class="grid">
            <div class="card">
                <h2>待办任务 ({len(latest.get('phases', {}).get('tasks', []))})</h2>
                <div class="tasks-list">
                    {self._generate_tasks_html(latest.get('phases', {}).get('tasks', [])[:10])}
                </div>
            </div>
            
            <div class="card">
                <h2>代码统计</h2>
                {self._generate_code_stats_html(latest.get('phases', {}).get('code_analysis', {}))}
            </div>
        </div>
        
        <div class="footer">
            <p>生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')} | 迭代次数: {len(history)}</p>
        </div>
    </div>
</body>
</html>"""
        
        return html
    
    def _generate_bars(self, iterations: list, passed: list, failed: list) -> str:
        """生成柱状图"""
        bars = []
        max_val = max(max(passed, default=1), max(failed, default=1))
        
        for i, (iter_num, p, f) in enumerate(zip(iterations, passed, failed)):
            height_p = (p / max_val * 180) if max_val > 0 else 0
            height_f = (f / max_val * 180) if max_val > 0 else 0
            
            bars.append(f'''
                <div class="bar" style="height: {height_p}px; background: linear-gradient(180deg, #00b09b, #96c93d);">
                    <span class="bar-value">{p}</span>
                    <span class="bar-label">#{iter_num}</span>
                </div>
                <div class="bar" style="height: {height_f}px; background: linear-gradient(180deg, #ff416c, #ff4b2b);">
                    <span class="bar-value">{f}</span>
                </div>
            ''')
        
        return ''.join(bars)
    
    def _generate_tasks_html(self, tasks: list) -> str:
        """生成任务列表HTML"""
        html = ""
        for task in tasks:
            priority = task.get("priority", "medium")
            html += f'''
                <div class="task-item {priority}">
                    <div class="task-title">{task.get('title', '')}</div>
                    <div class="task-meta">
                        类型: {task.get('type', '')} | 优先级: {priority}
                    </div>
                </div>
            '''
        return html
    
    def _generate_code_stats_html(self, code_analysis: Dict) -> str:
        """生成代码统计HTML"""
        if not code_analysis:
            return "<p>无数据</p>"
        
        return f'''
            <div style="margin-top: 15px;">
                <p>总文件数: <strong>{code_analysis.get('total_files', 0)}</strong></p>
                <p>总代码行: <strong>{code_analysis.get('total_lines', 0)}</strong></p>
                <hr style="border-color: rgba(255,255,255,0.1); margin: 15px 0;">
                <p>Server: {code_analysis.get('server', {}).get('files', 0)} 文件, {code_analysis.get('server', {}).get('lines', 0)} 行</p>
                <p>Client: {code_analysis.get('client', {}).get('files', 0)} 文件, {code_analysis.get('client', {}).get('lines', 0)} 行</p>
                <p>GM: {code_analysis.get('gm', {}).get('files', 0)} 文件, {code_analysis.get('gm', {}).get('lines', 0)} 行</p>
            </div>
        '''
