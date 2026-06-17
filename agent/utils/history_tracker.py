"""历史追踪器 - 记录和分析迭代历史"""
import json
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Optional


class HistoryTracker:
    """追踪迭代历史和进度"""
    
    def __init__(self, history_dir: Path):
        self.history_dir = history_dir
        self.history_dir.mkdir(parents=True, exist_ok=True)
        self.index_file = history_dir / "index.json"
        self._load_index()
    
    def _load_index(self):
        """加载历史索引"""
        if self.index_file.exists():
            with open(self.index_file, 'r', encoding='utf-8') as f:
                self.index = json.load(f)
        else:
            self.index = {
                "iterations": [],
                "created_at": datetime.now().isoformat(),
                "total_iterations": 0
            }
    
    def _save_index(self):
        """保存历史索引"""
        with open(self.index_file, 'w', encoding='utf-8') as f:
            json.dump(self.index, f, ensure_ascii=False, indent=2)
    
    def record_iteration(self, result: Dict) -> str:
        """记录一次迭代"""
        iteration_id = f"iter_{datetime.now().strftime('%Y%m%d_%H%M%S')}"
        
        # 保存详细结果
        detail_file = self.history_dir / f"{iteration_id}.json"
        with open(detail_file, 'w', encoding='utf-8') as f:
            json.dump(result, f, ensure_ascii=False, indent=2)
        
        # 提取关键指标
        phases = result.get("phases", {})
        test_result = phases.get("tests", {})
        gaps = phases.get("gaps", {})
        
        # 更新索引
        record = {
            "id": iteration_id,
            "timestamp": datetime.now().isoformat(),
            "iteration": result.get("iteration", 0),
            "duration": result.get("duration", 0),
            "metrics": {
                "test_passed": test_result.get("passed", 0),
                "test_failed": test_result.get("failed", 0),
                "critical_issues": len(gaps.get("critical", [])),
                "improvements": len(gaps.get("improvements", [])),
                "tasks_generated": len(phases.get("tasks", []))
            },
            "scores": gaps.get("metrics", {})
        }
        
        self.index["iterations"].append(record)
        self.index["total_iterations"] += 1
        self._save_index()
        
        return iteration_id
    
    def get_history(self, limit: int = 10) -> List[Dict]:
        """获取历史记录"""
        return self.index.get("iterations", [])[-limit:]
    
    def get_trend(self) -> Dict:
        """分析趋势"""
        iterations = self.index.get("iterations", [])
        if len(iterations) < 2:
            return {"trend": "insufficient_data"}
        
        # 提取指标序列
        test_passed = [i["metrics"]["test_passed"] for i in iterations]
        test_failed = [i["metrics"]["test_failed"] for i in iterations]
        critical = [i["metrics"]["critical_issues"] for i in iterations]
        scores = [i.get("scores", {}).get("overall_score", 0) for i in iterations]
        
        # 计算趋势
        trend = {
            "test_passed": self._calculate_trend(test_passed),
            "test_failed": self._calculate_trend(test_failed),
            "critical_issues": self._calculate_trend(critical),
            "overall_score": self._calculate_trend(scores),
            "total_iterations": len(iterations),
            "latest": iterations[-1] if iterations else None
        }
        
        return trend
    
    def _calculate_trend(self, values: List[float]) -> Dict:
        """计算趋势"""
        if len(values) < 2:
            return {"direction": "stable", "change": 0}
        
        recent = values[-1]
        previous = values[-2]
        change = recent - previous
        
        if change > 0:
            direction = "improving"
        elif change < 0:
            direction = "declining"
        else:
            direction = "stable"
        
        return {
            "direction": direction,
            "change": change,
            "current": recent,
            "previous": previous
        }
    
    def get_comparison(self, iter_id1: str, iter_id2: str) -> Dict:
        """比较两次迭代"""
        file1 = self.history_dir / f"{iter_id1}.json"
        file2 = self.history_dir / f"{iter_id2}.json"
        
        if not file1.exists() or not file2.exists():
            return {"error": "迭代记录不存在"}
        
        with open(file1, 'r', encoding='utf-8') as f:
            data1 = json.load(f)
        with open(file2, 'r', encoding='utf-8') as f:
            data2 = json.load(f)
        
        return {
            "iteration1": iter_id1,
            "iteration2": iter_id2,
            "comparison": self._compare_metrics(data1, data2)
        }
    
    def _compare_metrics(self, data1: Dict, data2: Dict) -> Dict:
        """比较指标"""
        def get_metrics(data):
            phases = data.get("phases", {})
            tests = phases.get("tests", {})
            gaps = phases.get("gaps", {})
            return {
                "test_passed": tests.get("passed", 0),
                "test_failed": tests.get("failed", 0),
                "critical": len(gaps.get("critical", [])),
                "score": gaps.get("metrics", {}).get("overall_score", 0)
            }
        
        m1 = get_metrics(data1)
        m2 = get_metrics(data2)
        
        return {
            key: {
                "before": m1[key],
                "after": m2[key],
                "change": m2[key] - m1[key]
            }
            for key in m1
        }
