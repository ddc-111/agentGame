"""代码库分析器 - 分析项目结构、代码质量、依赖关系"""
import os
import re
from pathlib import Path
from typing import Dict, List, Set
from collections import defaultdict


class CodeAnalyzer:
    """分析代码库结构和质量"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.server_dir = root_dir / "server"
        self.client_dir = root_dir / "client"
        self.gm_dir = root_dir / "gm"
    
    def analyze(self) -> Dict:
        """执行完整代码分析"""
        result = {
            "timestamp": self._get_timestamp(),
            "structure": self._analyze_structure(),
            "server": self._analyze_server(),
            "client": self._analyze_client(),
            "gm": self._analyze_gm(),
            "total_files": 0,
            "total_lines": 0,
            "complexity": {},
            "dependencies": self._analyze_dependencies()
        }
        
        # 计算总计
        for key in ["server", "client", "gm"]:
            result["total_files"] += result[key]["files"]
            result["total_lines"] += result[key]["lines"]
        
        # 分析复杂度
        result["complexity"] = self._analyze_complexity()
        
        return result
    
    def _get_timestamp(self) -> str:
        from datetime import datetime
        return datetime.now().isoformat()
    
    def _analyze_structure(self) -> Dict:
        """分析项目目录结构"""
        structure = {
            "directories": [],
            "key_files": []
        }
        
        for item in self.root_dir.iterdir():
            if item.is_dir() and not item.name.startswith('.'):
                structure["directories"].append(item.name)
            elif item.is_file() and item.suffix in ['.md', '.bat', '.json', '.yaml']:
                structure["key_files"].append(item.name)
        
        return structure
    
    def _analyze_server(self) -> Dict:
        """分析Go服务端代码"""
        if not self.server_dir.exists():
            return {"files": 0, "lines": 0, "packages": [], "endpoints": []}
        
        stats = {"files": 0, "lines": 0, "packages": [], "endpoints": [], "models": []}
        
        # 统计Go文件
        for go_file in self.server_dir.rglob("*.go"):
            if "_test.go" not in go_file.name:
                stats["files"] += 1
                stats["lines"] += self._count_lines(go_file)
        
        # 分析包结构
        internal_dir = self.server_dir / "internal"
        if internal_dir.exists():
            for pkg_dir in internal_dir.iterdir():
                if pkg_dir.is_dir():
                    stats["packages"].append(pkg_dir.name)
        
        # 分析API端点
        stats["endpoints"] = self._extract_api_endpoints()
        
        # 分析模型
        stats["models"] = self._extract_models()
        
        return stats
    
    def _analyze_client(self) -> Dict:
        """分析Phaser客户端代码"""
        if not self.client_dir.exists():
            return {"files": 0, "lines": 0, "scenes": [], "systems": []}
        
        stats = {"files": 0, "lines": 0, "scenes": [], "systems": []}
        
        # 统计JS文件
        for js_file in self.client_dir.rglob("*.js"):
            if "test" not in js_file.name and "node_modules" not in str(js_file):
                stats["files"] += 1
                stats["lines"] += self._count_lines(js_file)
        
        # 分析场景
        scenes_dir = self.client_dir / "src" / "game" / "scenes"
        if scenes_dir.exists():
            stats["scenes"] = [f.stem for f in scenes_dir.glob("*.js")]
        
        # 分析系统
        systems_dir = self.client_dir / "src" / "game" / "systems"
        if systems_dir.exists():
            stats["systems"] = [f.stem for f in systems_dir.glob("*.js")]
        
        return stats
    
    def _analyze_gm(self) -> Dict:
        """分析GM管理端代码"""
        if not self.gm_dir.exists():
            return {"files": 0, "lines": 0, "views": [], "stores": [], "components": []}
        
        stats = {"files": 0, "lines": 0, "views": [], "stores": [], "components": []}
        
        # 统计Vue/JS文件
        for vue_file in self.gm_dir.rglob("*.vue"):
            if "node_modules" not in str(vue_file):
                stats["files"] += 1
                stats["lines"] += self._count_lines(vue_file)
        
        for js_file in self.gm_dir.rglob("*.js"):
            if "test" not in js_file.name and "node_modules" not in str(js_file):
                stats["files"] += 1
                stats["lines"] += self._count_lines(js_file)
        
        # 分析视图
        views_dir = self.gm_dir / "src" / "views"
        if views_dir.exists():
            stats["views"] = [f.stem for f in views_dir.glob("*.vue")]
        
        # 分析Store
        stores_dir = self.gm_dir / "src" / "stores"
        if stores_dir.exists():
            stats["stores"] = [f.stem for f in stores_dir.glob("*.js")]
        
        # 分析组件
        components_dir = self.gm_dir / "src" / "components"
        if components_dir.exists():
            stats["components"] = [f.stem for f in components_dir.rglob("*.vue")]
        
        return stats
    
    def _count_lines(self, file_path: Path) -> int:
        """统计文件行数"""
        try:
            with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                return len(f.readlines())
        except:
            return 0
    
    def _extract_api_endpoints(self) -> List[str]:
        """提取API端点"""
        endpoints = []
        router_file = self.server_dir / "internal" / "network" / "router.go"
        
        if router_file.exists():
            content = router_file.read_text(encoding='utf-8')
            # 匹配路由定义
            patterns = [
                r'router\.(GET|POST|PUT|DELETE|PATCH)\s*\(\s*"([^"]+)"',
                r'group\.(GET|POST|PUT|DELETE|PATCH)\s*\(\s*"([^"]+)"',
                r'\.(GET|POST|PUT|DELETE|PATCH)\s*\(\s*"([^"]+)"'
            ]
            for pattern in patterns:
                matches = re.findall(pattern, content)
                for method, path in matches:
                    endpoints.append(f"{method} {path}")
        
        return endpoints
    
    def _extract_models(self) -> List[str]:
        """提取数据模型"""
        models = []
        models_dir = self.server_dir / "internal" / "database" / "models"
        
        if models_dir.exists():
            for go_file in models_dir.glob("*.go"):
                content = go_file.read_text(encoding='utf-8')
                # 匹配结构体定义
                matches = re.findall(r'type\s+(\w+)\s+struct', content)
                models.extend(matches)
        
        return models
    
    def _analyze_dependencies(self) -> Dict:
        """分析项目依赖"""
        deps = {
            "server": [],
            "client": [],
            "gm": []
        }
        
        # Go依赖
        go_mod = self.server_dir / "go.mod"
        if go_mod.exists():
            content = go_mod.read_text(encoding='utf-8')
            for line in content.split('\n'):
                line = line.strip()
                if line and not line.startswith('//') and not line.startswith('module') and not line.startswith('go'):
                    parts = line.split()
                    if len(parts) >= 1:
                        deps["server"].append(parts[0])
        
        # Client依赖
        client_pkg = self.client_dir / "package.json"
        if client_pkg.exists():
            import json
            with open(client_pkg, 'r', encoding='utf-8') as f:
                pkg = json.load(f)
                deps["client"] = list(pkg.get("dependencies", {}).keys())
        
        # GM依赖
        gm_pkg = self.gm_dir / "package.json"
        if gm_pkg.exists():
            import json
            with open(gm_pkg, 'r', encoding='utf-8') as f:
                pkg = json.load(f)
                deps["gm"] = list(pkg.get("dependencies", {}).keys())
        
        return deps
    
    def _analyze_complexity(self) -> Dict:
        """分析代码复杂度"""
        complexity = {
            "server": {"functions": 0, "avg_params": 0},
            "client": {"classes": 0, "methods": 0},
            "gm": {"components": 0, "computed": 0}
        }
        
        # 分析Server函数
        for go_file in self.server_dir.rglob("*.go"):
            if "_test.go" not in go_file.name:
                content = go_file.read_text(encoding='utf-8')
                functions = re.findall(r'func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(([^)]*)\)', content)
                complexity["server"]["functions"] += len(functions)
        
        # 分析Client类和方法
        for js_file in self.client_dir.rglob("*.js"):
            if "test" not in js_file.name and "node_modules" not in str(js_file):
                content = js_file.read_text(encoding='utf-8')
                complexity["client"]["classes"] += len(re.findall(r'class\s+\w+', content))
                complexity["client"]["methods"] += len(re.findall(r'(?:async\s+)?(\w+)\s*\([^)]*\)\s*{', content))
        
        # 分析GM组件
        for vue_file in self.gm_dir.rglob("*.vue"):
            if "node_modules" not in str(vue_file):
                content = vue_file.read_text(encoding='utf-8')
                complexity["gm"]["components"] += 1
                complexity["gm"]["computed"] += len(re.findall(r'computed:', content))
        
        return complexity
