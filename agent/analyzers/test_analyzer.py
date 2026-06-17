"""测试分析器 - 分析测试覆盖、测试质量"""
import re
from pathlib import Path
from typing import Dict, List


class TestAnalyzer:
    """分析测试套件"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.server_dir = root_dir / "server"
        self.client_dir = root_dir / "client"
        self.gm_dir = root_dir / "gm"
    
    def analyze(self) -> Dict:
        """执行完整测试分析"""
        result = {
            "timestamp": self._get_timestamp(),
            "server": self._analyze_server_tests(),
            "client": self._analyze_client_tests(),
            "gm": self._analyze_gm_tests(),
            "total_test_files": 0,
            "total_test_cases": 0,
            "coverage_gaps": []
        }
        
        # 计算总计
        for key in ["server", "client", "gm"]:
            result["total_test_files"] += result[key]["test_files"]
            result["total_test_cases"] += result[key]["test_cases"]
        
        # 分析覆盖差距
        result["coverage_gaps"] = self._identify_coverage_gaps(result)
        
        return result
    
    def _get_timestamp(self) -> str:
        from datetime import datetime
        return datetime.now().isoformat()
    
    def _analyze_server_tests(self) -> Dict:
        """分析服务端测试"""
        stats = {
            "test_files": 0,
            "test_cases": 0,
            "test_types": {
                "unit": 0,
                "integration": 0,
                "api": 0,
                "benchmark": 0
            },
            "files": [],
            "untested_packages": []
        }
        
        if not self.server_dir.exists():
            return stats
        
        # 统计测试文件
        for test_file in self.server_dir.rglob("*_test.go"):
            stats["test_files"] += 1
            stats["files"].append(str(test_file.relative_to(self.server_dir)))
            
            # 统计测试用例
            content = test_file.read_text(encoding='utf-8')
            test_funcs = re.findall(r'func\s+(Test\w+)\s*\(', content)
            stats["test_cases"] += len(test_funcs)
            
            # 分类测试类型
            if "benchmark" in test_file.name.lower() or "Benchmark" in content:
                stats["test_types"]["benchmark"] += len(test_funcs)
            elif "api" in test_file.name.lower():
                stats["test_types"]["api"] += len(test_funcs)
            elif "integration" in test_file.name.lower():
                stats["test_types"]["integration"] += len(test_funcs)
            else:
                stats["test_types"]["unit"] += len(test_funcs)
        
        # 检查未测试的包
        internal_dir = self.server_dir / "internal"
        if internal_dir.exists():
            tested_packages = set()
            for test_file in self.server_dir.rglob("*_test.go"):
                # 从文件路径提取包名
                rel_path = test_file.relative_to(self.server_dir)
                parts = rel_path.parts
                if len(parts) > 2:
                    tested_packages.add(parts[1])
            
            for pkg_dir in internal_dir.iterdir():
                if pkg_dir.is_dir() and pkg_dir.name not in tested_packages:
                    stats["untested_packages"].append(pkg_dir.name)
        
        return stats
    
    def _analyze_client_tests(self) -> Dict:
        """分析客户端测试"""
        stats = {
            "test_files": 0,
            "test_cases": 0,
            "files": [],
            "untested_modules": []
        }
        
        if not self.client_dir.exists():
            return stats
        
        # 查找测试文件
        test_dirs = [
            self.client_dir / "src" / "__tests__",
            self.client_dir / "tests",
            self.client_dir / "test"
        ]
        
        for test_dir in test_dirs:
            if test_dir.exists():
                for test_file in test_dir.rglob("*.test.js"):
                    stats["test_files"] += 1
                    stats["files"].append(str(test_file.relative_to(self.client_dir)))
                    
                    content = test_file.read_text(encoding='utf-8')
                    test_cases = re.findall(r'it\s*\(\s*[\'"](.+?)[\'"]', content)
                    stats["test_cases"] += len(test_cases)
        
        # 检查未测试的模块
        systems_dir = self.client_dir / "src" / "game" / "systems"
        scenes_dir = self.client_dir / "src" / "game" / "scenes"
        
        tested_modules = set()
        for test_file in stats["files"]:
            # 从测试文件名推断被测模块
            module_name = Path(test_file).stem.replace(".test", "")
            tested_modules.add(module_name)
        
        for dir_path in [systems_dir, scenes_dir]:
            if dir_path.exists():
                for js_file in dir_path.glob("*.js"):
                    if js_file.stem not in tested_modules:
                        stats["untested_modules"].append(f"{dir_path.name}/{js_file.name}")
        
        return stats
    
    def _analyze_gm_tests(self) -> Dict:
        """分析GM管理端测试"""
        stats = {
            "test_files": 0,
            "test_cases": 0,
            "files": [],
            "untested_components": [],
            "untested_stores": []
        }
        
        if not self.gm_dir.exists():
            return stats
        
        # 查找测试文件
        test_dir = self.gm_dir / "src" / "__tests__"
        if test_dir.exists():
            for test_file in test_dir.rglob("*.test.js"):
                stats["test_files"] += 1
                stats["files"].append(str(test_file.relative_to(self.gm_dir)))
                
                content = test_file.read_text(encoding='utf-8')
                test_cases = re.findall(r'it\s*\(\s*[\'"](.+?)[\'"]', content)
                stats["test_cases"] += len(test_cases)
        
        # 检查未测试的Store
        stores_dir = self.gm_dir / "src" / "stores"
        if stores_dir.exists():
            tested_stores = set()
            for test_file in stats["files"]:
                if "stores" in test_file:
                    store_name = Path(test_file).stem.replace(".test", "")
                    tested_stores.add(store_name)
            
            for store_file in stores_dir.glob("*.js"):
                if store_file.stem not in tested_stores:
                    stats["untested_stores"].append(store_file.name)
        
        # 检查未测试的组件
        components_dir = self.gm_dir / "src" / "components"
        if components_dir.exists():
            tested_components = set()
            for test_file in stats["files"]:
                if "components" in test_file:
                    comp_name = Path(test_file).stem.replace(".test", "")
                    tested_components.add(comp_name)
            
            for vue_file in components_dir.rglob("*.vue"):
                if vue_file.stem not in tested_components:
                    stats["untested_components"].append(vue_file.name)
        
        return stats
    
    def _identify_coverage_gaps(self, analysis: Dict) -> List[Dict]:
        """识别覆盖差距"""
        gaps = []
        
        # Server未测试的包
        for pkg in analysis["server"].get("untested_packages", []):
            gaps.append({
                "type": "untested_package",
                "target": "server",
                "name": pkg,
                "priority": "high" if pkg in ["network", "agent", "game"] else "medium"
            })
        
        # Client未测试的模块
        for module in analysis["client"].get("untested_modules", []):
            gaps.append({
                "type": "untested_module",
                "target": "client",
                "name": module,
                "priority": "high" if "system" in module else "medium"
            })
        
        # GM未测试的Store
        for store in analysis["gm"].get("untested_stores", []):
            gaps.append({
                "type": "untested_store",
                "target": "gm",
                "name": store,
                "priority": "high"
            })
        
        # GM未测试的组件
        for comp in analysis["gm"].get("untested_components", []):
            gaps.append({
                "type": "untested_component",
                "target": "gm",
                "name": comp,
                "priority": "medium"
            })
        
        return gaps
