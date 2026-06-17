"""配置管理 - 支持YAML配置文件"""
import json
import os
from pathlib import Path
from typing import Dict, Any, Optional

try:
    import yaml
    HAS_YAML = True
except ImportError:
    HAS_YAML = False


class AgentConfig:
    """Agent配置管理"""
    
    DEFAULT_CONFIG = {
        "version": "1.0.0",
        "project_root": None,  # 自动检测
        
        "iterations": {
            "max_iterations": 10,
            "delay_between": 1,  # 秒
            "stop_on_critical": False,
            "stop_on_stable": True
        },
        
        "analysis": {
            "include_patterns": ["*.go", "*.js", "*.vue", "*.ts"],
            "exclude_patterns": ["node_modules", "vendor", ".git", "__pycache__"],
            "complexity_threshold": 10
        },
        
        "testing": {
            "timeout": 300,
            "retry_on_failure": True,
            "max_retries": 3,
            "coverage_threshold": {
                "server": 60,
                "client": 50,
                "gm": 50
            }
        },
        
        "llm": {
            "enabled": False,
            "api_url": "https://api.openai.com/v1",
            "api_key": "",
            "model": "gpt-4",
            "temperature": 0.3,
            "max_tokens": 4096
        },
        
        "tasks": {
            "auto_execute": False,
            "verify_after_execute": True,
            "max_concurrent": 1,
            "priority_filter": ["critical", "high", "medium", "low"]
        },
        
        "reports": {
            "format": ["markdown", "json"],
            "include_code_snippets": True,
            "include_history": True,
            "history_limit": 10
        },
        
        "notifications": {
            "enabled": False,
            "on_critical": True,
            "on_failure": True,
            "webhook_url": ""
        }
    }
    
    def __init__(self, config_path: Path = None):
        self.config_path = config_path
        self.config = self.DEFAULT_CONFIG.copy()
        
        if config_path and config_path.exists():
            self.load()
    
    def load(self, path: Path = None):
        """加载配置文件"""
        path = path or self.config_path
        if not path or not path.exists():
            return
        
        try:
            content = path.read_text(encoding='utf-8')
            
            if path.suffix in ['.yaml', '.yml'] and HAS_YAML:
                loaded = yaml.safe_load(content)
            elif path.suffix == '.json':
                loaded = json.loads(content)
            else:
                return
            
            if loaded:
                self._merge_config(self.config, loaded)
                
        except Exception as e:
            print(f"  [!] 加载配置失败: {e}")
    
    def save(self, path: Path = None):
        """保存配置文件"""
        path = path or self.config_path
        if not path:
            path = Path("agent_config.yaml")
        
        try:
            path.parent.mkdir(parents=True, exist_ok=True)
            
            if path.suffix in ['.yaml', '.yml'] and HAS_YAML:
                content = yaml.dump(self.config, allow_unicode=True, default_flow_style=False)
            else:
                content = json.dumps(self.config, ensure_ascii=False, indent=2)
            
            path.write_text(content, encoding='utf-8')
            
        except Exception as e:
            print(f"  [!] 保存配置失败: {e}")
    
    def get(self, key: str, default: Any = None) -> Any:
        """获取配置值"""
        keys = key.split('.')
        value = self.config
        
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k)
            else:
                return default
        
        return value if value is not None else default
    
    def set(self, key: str, value: Any):
        """设置配置值"""
        keys = key.split('.')
        config = self.config
        
        for k in keys[:-1]:
            if k not in config:
                config[k] = {}
            config = config[k]
        
        config[keys[-1]] = value
    
    def _merge_config(self, base: Dict, override: Dict):
        """合并配置"""
        for key, value in override.items():
            if key in base and isinstance(base[key], dict) and isinstance(value, dict):
                self._merge_config(base[key], value)
            else:
                base[key] = value
    
    def get_llm_config(self) -> Dict:
        """获取LLM配置"""
        return self.config.get("llm", {})
    
    def is_llm_enabled(self) -> bool:
        """检查LLM是否启用"""
        llm = self.get_llm_config()
        return llm.get("enabled", False) and bool(llm.get("api_key"))
    
    def get_project_root(self) -> Path:
        """获取项目根目录"""
        if self.config.get("project_root"):
            return Path(self.config["project_root"])
        return Path(__file__).parent.parent
    
    def to_dict(self) -> Dict:
        """转换为字典"""
        return self.config.copy()
    
    @classmethod
    def from_env(cls) -> 'AgentConfig':
        """从环境变量创建配置"""
        config = cls()
        
        # LLM配置
        if os.getenv("LLM_API_KEY"):
            config.set("llm.enabled", True)
            config.set("llm.api_key", os.getenv("LLM_API_KEY"))
        if os.getenv("LLM_API_URL"):
            config.set("llm.api_url", os.getenv("LLM_API_URL"))
        if os.getenv("LLM_MODEL"):
            config.set("llm.model", os.getenv("LLM_MODEL"))
        
        return config
