"""构建执行器 - 验证三端构建"""
import subprocess
from pathlib import Path
from typing import Dict
from datetime import datetime


class BuildExecutor:
    """执行构建验证"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.server_dir = root_dir / "server"
        self.client_dir = root_dir / "client"
        self.gm_dir = root_dir / "gm"
    
    def build_all(self) -> Dict[str, bool]:
        """构建所有端"""
        return {
            "server": self.build_server(),
            "client": self.build_client(),
            "gm": self.build_gm()
        }
    
    def build_server(self) -> bool:
        """构建Go服务端"""
        try:
            proc = subprocess.run(
                ["go", "build", "-o", "bin\\gameserver.exe", "cmd\\gameserver\\main.go"],
                cwd=self.server_dir,
                capture_output=True,
                text=True,
                timeout=120
            )
            return proc.returncode == 0
        except Exception as e:
            print(f"Server构建失败: {e}")
            return False
    
    def build_client(self) -> bool:
        """构建客户端"""
        try:
            proc = subprocess.run(
                ["npm", "run", "build"],
                cwd=self.client_dir,
                capture_output=True,
                text=True,
                timeout=120,
                shell=True
            )
            return proc.returncode == 0
        except Exception as e:
            print(f"Client构建失败: {e}")
            return False
    
    def build_gm(self) -> bool:
        """构建GM管理端"""
        try:
            proc = subprocess.run(
                ["npm", "run", "build"],
                cwd=self.gm_dir,
                capture_output=True,
                text=True,
                timeout=120,
                shell=True
            )
            return proc.returncode == 0
        except Exception as e:
            print(f"GM构建失败: {e}")
            return False
