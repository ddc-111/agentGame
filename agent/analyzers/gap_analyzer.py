"""差距分析器 - 识别代码库的改进点和问题"""
from pathlib import Path
from typing import Dict, List


class GapAnalyzer:
    """分析代码库差距和改进点"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
    
    def analyze_gaps(self, code_analysis: Dict, test_analysis: Dict, test_result: Dict) -> Dict:
        """分析改进差距"""
        gaps = {
            "critical": [],      # 关键问题，必须修复
            "improvements": [],  # 改进项
            "opportunities": [], # 优化机会
            "metrics": {}
        }
        
        # 分析测试失败
        gaps["critical"].extend(self._analyze_test_failures(test_result))
        
        # 分析覆盖差距
        gaps["improvements"].extend(self._analyze_coverage_gaps(test_analysis))
        
        # 分析代码质量
        gaps["improvements"].extend(self._analyze_code_quality(code_analysis))
        
        # 分析架构改进
        gaps["opportunities"].extend(self._analyze_architecture(code_analysis))
        
        # 分析现代化改进
        gaps["opportunities"].extend(self._analyze_modernization(code_analysis))
        
        # 计算指标
        gaps["metrics"] = self._calculate_metrics(code_analysis, test_analysis, test_result)
        
        return gaps
    
    def _analyze_test_failures(self, test_result: Dict) -> List[Dict]:
        """分析测试失败"""
        issues = []
        
        if test_result.get("failed", 0) > 0:
            for failure in test_result.get("failures", []):
                issues.append({
                    "type": "test_failure",
                    "severity": "critical",
                    "target": failure.get("target", "unknown"),
                    "test": failure.get("test", "unknown"),
                    "error": failure.get("error", ""),
                    "suggestion": self._suggest_fix(failure)
                })
        
        return issues
    
    def _analyze_coverage_gaps(self, test_analysis: Dict) -> List[Dict]:
        """分析覆盖差距"""
        gaps = []
        
        for gap in test_analysis.get("coverage_gaps", []):
            gaps.append({
                "type": "coverage_gap",
                "severity": gap.get("priority", "medium"),
                "target": gap.get("target"),
                "name": gap.get("name"),
                "suggestion": f"为 {gap.get('target')}/{gap.get('name')} 添加单元测试"
            })
        
        return gaps
    
    def _analyze_code_quality(self, code_analysis: Dict) -> List[Dict]:
        """分析代码质量"""
        issues = []
        
        # 检查大文件
        for target in ["server", "client", "gm"]:
            if code_analysis.get(target, {}).get("lines", 0) > 10000:
                issues.append({
                    "type": "large_codebase",
                    "severity": "low",
                    "target": target,
                    "lines": code_analysis[target]["lines"],
                    "suggestion": f"考虑将 {target} 拆分为更小的模块"
                })
        
        # 检查高复杂度
        complexity = code_analysis.get("complexity", {})
        if complexity.get("server", {}).get("functions", 0) > 100:
            issues.append({
                "type": "high_complexity",
                "severity": "medium",
                "target": "server",
                "functions": complexity["server"]["functions"],
                "suggestion": "Server端函数过多，考虑重构为更小的服务"
            })
        
        # 检查代码重复和复用机会
        issues.extend(self._analyze_code_duplication(code_analysis))
        
        return issues
    
    def _analyze_code_duplication(self, code_analysis: Dict) -> List[Dict]:
        """分析代码重复和复用机会"""
        issues = []
        
        # 检查server端大型文件
        server_lines = code_analysis.get("server", {}).get("lines", 0)
        if server_lines > 15000:
            # 扫描server目录下的大型文件
            server_dir = self.root_dir / "server" / "internal"
            if server_dir.exists():
                for file_path in server_dir.rglob("*.go"):
                    if "_test.go" not in file_path.name:
                        try:
                            lines = sum(1 for _ in open(file_path, encoding='utf-8', errors='ignore'))
                            if lines > 500:
                                issues.append({
                                    "type": "refactor_opportunity",
                                    "severity": "medium",
                                    "target": "server",
                                    "file": str(file_path.relative_to(self.root_dir)),
                                    "lines": lines,
                                    "suggestion": f"大型Go文件 {file_path.name} ({lines}行) 可以拆分或抽象公共逻辑"
                                })
                        except:
                            pass
        
        # 检查client端大型文件
        client_lines = code_analysis.get("client", {}).get("lines", 0)
        if client_lines > 8000:
            client_dir = self.root_dir / "client" / "src"
            if client_dir.exists():
                for file_path in client_dir.rglob("*.js"):
                    if "test" not in file_path.name and "node_modules" not in str(file_path):
                        try:
                            lines = sum(1 for _ in open(file_path, encoding='utf-8', errors='ignore'))
                            if lines > 800:
                                issues.append({
                                    "type": "refactor_opportunity",
                                    "severity": "medium",
                                    "target": "client",
                                    "file": str(file_path.relative_to(self.root_dir)),
                                    "lines": lines,
                                    "suggestion": f"大型JS文件 {file_path.name} ({lines}行) 可以拆分为更小的模块"
                                })
                        except:
                            pass
        
        # 检查gm端大型文件
        gm_lines = code_analysis.get("gm", {}).get("lines", 0)
        if gm_lines > 8000:
            gm_dir = self.root_dir / "gm" / "src"
            if gm_dir.exists():
                for file_path in gm_dir.rglob("*.vue"):
                    try:
                        lines = sum(1 for _ in open(file_path, encoding='utf-8', errors='ignore'))
                        if lines > 500:
                            issues.append({
                                "type": "refactor_opportunity",
                                "severity": "low",
                                "target": "gm",
                                "file": str(file_path.relative_to(self.root_dir)),
                                "lines": lines,
                                "suggestion": f"大型Vue文件 {file_path.name} ({lines}行) 可以拆分或抽取公共组件"
                            })
                    except:
                        pass
        
        return issues
    
    def _analyze_architecture(self, code_analysis: Dict) -> List[Dict]:
        """分析架构改进"""
        opportunities = []
        
        # 检查是否有分层架构
        server_packages = code_analysis.get("server", {}).get("packages", [])
        required_packages = ["network", "database", "game", "agent"]
        
        for pkg in required_packages:
            if pkg not in server_packages:
                opportunities.append({
                    "type": "missing_layer",
                    "target": "server",
                    "package": pkg,
                    "suggestion": f"添加 {pkg} 包以完善分层架构"
                })
        
        # 检查前端是否有状态管理
        if "stores" not in code_analysis.get("gm", {}).get("stores", []):
            opportunities.append({
                "type": "state_management",
                "target": "gm",
                "suggestion": "使用Pinia进行状态管理"
            })
        
        return opportunities
    
    def _analyze_modernization(self, code_analysis: Dict) -> List[Dict]:
        """分析现代化改进"""
        opportunities = []
        
        # 检查依赖
        deps = code_analysis.get("dependencies", {})
        
        # 检查是否使用现代框架
        client_deps = deps.get("client", [])
        if "phaser" in client_deps:
            opportunities.append({
                "type": "framework_upgrade",
                "target": "client",
                "current": "Phaser 3",
                "suggestion": "考虑升级到Phaser最新版本或评估其他框架"
            })
        
        # 检查是否使用TypeScript
        opportunities.append({
            "type": "typescript_migration",
            "target": "client,gm",
            "suggestion": "考虑迁移到TypeScript以获得更好的类型安全"
        })
        
        return opportunities
    
    def _calculate_metrics(self, code_analysis: Dict, test_analysis: Dict, test_result: Dict) -> Dict:
        """计算质量指标"""
        metrics = {
            "test_pass_rate": 0,
            "test_coverage_score": 0,
            "code_health_score": 0,
            "overall_score": 0
        }
        
        # 测试通过率
        total_tests = test_result.get("passed", 0) + test_result.get("failed", 0)
        if total_tests > 0:
            metrics["test_pass_rate"] = (test_result.get("passed", 0) / total_tests) * 100
        
        # 测试覆盖分数
        total_test_files = test_analysis.get("total_test_files", 0)
        total_source_files = code_analysis.get("total_files", 1)
        metrics["test_coverage_score"] = min(100, (total_test_files / total_source_files) * 100 * 10)
        
        # 代码健康分数
        critical_count = 0  # 从传入参数计算
        metrics["code_health_score"] = max(0, 100 - (critical_count * 20))
        
        # 总体分数
        metrics["overall_score"] = (
            metrics["test_pass_rate"] * 0.4 +
            metrics["test_coverage_score"] * 0.3 +
            metrics["code_health_score"] * 0.3
        )
        
        return metrics
    
    def _suggest_fix(self, failure: Dict) -> str:
        """根据测试失败提供建议"""
        error = failure.get("error", "").lower()
        
        if "assertion" in error or "expected" in error:
            return "检查测试断言是否与实际行为匹配"
        elif "timeout" in error:
            return "增加超时时间或优化代码性能"
        elif "import" in error or "module" in error:
            return "检查依赖是否正确安装"
        elif "null" in error or "undefined" in error:
            return "添加空值检查"
        else:
            return "查看错误详情并修复相关代码"
