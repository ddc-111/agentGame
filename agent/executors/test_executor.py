"""测试执行器 - 运行三端测试并收集结果"""
import subprocess
import json
import re
import os
from pathlib import Path
from typing import Dict, List
from datetime import datetime


class TestExecutor:
    """执行测试套件"""
    
    def __init__(self, root_dir: Path):
        self.root_dir = root_dir
        self.server_dir = root_dir / "server"
        self.client_dir = root_dir / "client"
        self.gm_dir = root_dir / "gm"
        
        # 设置环境变量以解决Windows编码问题
        self.env = os.environ.copy()
        self.env["PYTHONIOENCODING"] = "utf-8"
        if os.name == 'nt':  # Windows
            self.env["CHCP"] = "65001"
    
    def run_all_tests(self) -> Dict:
        """运行所有测试"""
        result = {
            "timestamp": datetime.now().isoformat(),
            "server": self.run_server_tests(),
            "client": self.run_client_tests(),
            "gm": self.run_gm_tests(),
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "failures": [],
            "duration": 0
        }
        
        # 汇总结果
        for target in ["server", "client", "gm"]:
            result["passed"] += result[target].get("passed", 0)
            result["failed"] += result[target].get("failed", 0)
            result["skipped"] += result[target].get("skipped", 0)
            result["failures"].extend(result[target].get("failures", []))
            result["duration"] += result[target].get("duration", 0)
        
        return result
    
    def run_server_tests(self) -> Dict:
        """运行服务端Go测试"""
        result = {
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "failures": [],
            "duration": 0,
            "output": ""
        }
        
        try:
            start_time = datetime.now()
            
            # 运行Go测试
            proc = subprocess.run(
                ["go", "test", "./internal/...", "-v", "-count=1", "-json"],
                cwd=self.server_dir,
                capture_output=True,
                text=True,
                timeout=300,
                env=self.env,
                encoding='utf-8',
                errors='replace'
            )
            
            result["output"] = proc.stdout
            result["duration"] = (datetime.now() - start_time).total_seconds()
            
            # 解析JSON输出
            test_packages = {}
            for line in proc.stdout.strip().split('\n'):
                if not line:
                    continue
                try:
                    event = json.loads(line)
                    package = event.get("Package", "unknown")
                    test = event.get("Test", "")
                    action = event.get("Action", "")
                    
                    if package not in test_packages:
                        test_packages[package] = {"passed": 0, "failed": 0, "skipped": 0}
                    
                    if action == "pass" and test:
                        test_packages[package]["passed"] += 1
                    elif action == "fail" and test:
                        test_packages[package]["failed"] += 1
                        result["failures"].append({
                            "target": "server",
                            "package": package,
                            "test": test,
                            "error": event.get("Output", "").strip()
                        })
                    elif action == "skip" and test:
                        test_packages[package]["skipped"] += 1
                        
                except json.JSONDecodeError:
                    continue
            
            # 汇总结果
            for pkg, stats in test_packages.items():
                result["passed"] += stats["passed"]
                result["failed"] += stats["failed"]
                result["skipped"] += stats["skipped"]
            
            # 如果没有解析到测试，使用简单解析
            if result["passed"] == 0 and result["failed"] == 0:
                if proc.returncode == 0:
                    # 尝试匹配 "ok" 或 "PASS"
                    ok_matches = re.findall(r'^ok\s+', proc.stdout, re.MULTILINE)
                    pass_matches = re.findall(r'^PASS', proc.stdout, re.MULTILINE)
                    result["passed"] = max(len(ok_matches), len(pass_matches), 1)
                else:
                    result["failed"] = 1
                    # 提取错误信息
                    fail_matches = re.findall(r'^FAIL\s+(.+)', proc.stdout, re.MULTILINE)
                    error_msg = proc.stderr or proc.stdout or "未知错误"
                    result["failures"].append({
                        "target": "server",
                        "test": fail_matches[0] if fail_matches else "unknown",
                        "error": error_msg[:500]
                    })
            
        except subprocess.TimeoutExpired:
            result["failed"] = 1
            result["failures"].append({
                "target": "server",
                "test": "timeout",
                "error": "测试执行超时(300秒)"
            })
        except FileNotFoundError:
            result["failed"] = 1
            result["failures"].append({
                "target": "server",
                "test": "environment",
                "error": "找不到go命令，请确保Go已安装并在PATH中"
            })
        except Exception as e:
            result["failed"] = 1
            result["failures"].append({
                "target": "server",
                "test": "error",
                "error": str(e)
            })
        
        return result
    
    def run_client_tests(self) -> Dict:
        """运行客户端测试"""
        return self._run_vitest_tests("client", self.client_dir)
    
    def run_gm_tests(self) -> Dict:
        """运行GM管理端测试"""
        return self._run_vitest_tests("gm", self.gm_dir)
    
    def _run_vitest_tests(self, target: str, work_dir: Path) -> Dict:
        """运行Vitest测试"""
        result = {
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "failures": [],
            "duration": 0,
            "output": ""
        }
        
        try:
            start_time = datetime.now()
            
            # 运行Vitest
            proc = subprocess.run(
                ["npm", "test", "--", "--reporter=json"],
                cwd=work_dir,
                capture_output=True,
                text=True,
                timeout=120,
                shell=True,
                env=self.env,
                encoding='utf-8',
                errors='replace'
            )
            
            result["output"] = proc.stdout
            result["duration"] = (datetime.now() - start_time).total_seconds()
            
            # 尝试解析JSON输出
            json_match = re.search(r'\{[\s\S]*\}', proc.stdout)
            if json_match:
                try:
                    test_data = json.loads(json_match.group())
                    result["passed"] = test_data.get("numPassedTests", 0)
                    result["failed"] = test_data.get("numFailedTests", 0)
                    result["skipped"] = test_data.get("numPendingTests", 0)
                    
                    # 收集失败详情
                    for test_result in test_data.get("testResults", []):
                        for assertion in test_result.get("assertionResults", []):
                            if assertion.get("status") == "failed":
                                failure_msgs = assertion.get("failureMessages", [])
                                error_msg = failure_msgs[0] if failure_msgs else "未知错误"
                                result["failures"].append({
                                    "target": target,
                                    "test": assertion.get("fullName", ""),
                                    "error": error_msg[:500]
                                })
                except json.JSONDecodeError:
                    pass
            
            # 如果JSON解析失败，使用简单解析
            if result["passed"] == 0 and result["failed"] == 0:
                if proc.returncode == 0:
                    # 匹配 "X passed"
                    match = re.search(r'(\d+)\s+passed', proc.stdout)
                    if match:
                        result["passed"] = int(match.group(1))
                    else:
                        result["passed"] = 1
                else:
                    result["failed"] = 1
                    # 提取失败信息
                    fail_match = re.search(r'(\d+)\s+failed', proc.stdout)
                    if fail_match:
                        result["failed"] = int(fail_match.group(1))
                    result["failures"].append({
                        "target": target,
                        "test": "unknown",
                        "error": (proc.stderr or proc.stdout)[:500]
                    })
            
        except subprocess.TimeoutExpired:
            result["failed"] = 1
            result["failures"].append({
                "target": target,
                "test": "timeout",
                "error": "测试执行超时(120秒)"
            })
        except FileNotFoundError:
            result["failed"] = 1
            result["failures"].append({
                "target": target,
                "test": "environment",
                "error": "找不到npm命令，请确保Node.js已安装并在PATH中"
            })
        except Exception as e:
            result["failed"] = 1
            result["failures"].append({
                "target": target,
                "test": "error",
                "error": str(e)
            })
        
        return result
    
    def run_specific_test(self, target: str, test_pattern: str = None) -> Dict:
        """运行特定测试"""
        if target == "server":
            cmd = ["go", "test", "./internal/...", "-v", "-count=1"]
            if test_pattern:
                cmd.extend(["-run", test_pattern])
            return self._run_command("server", self.server_dir, cmd)
        elif target in ["client", "gm"]:
            work_dir = self.client_dir if target == "client" else self.gm_dir
            cmd = ["npm", "test"]
            if test_pattern:
                cmd.extend(["--", "--testNamePattern", test_pattern])
            return self._run_command(target, work_dir, cmd, shell=True)
        else:
            return {"passed": 0, "failed": 1, "failures": [{"error": f"未知目标: {target}"}]}
    
    def _run_command(self, target: str, work_dir: Path, cmd: List[str], shell: bool = False) -> Dict:
        """运行命令"""
        result = {
            "passed": 0,
            "failed": 0,
            "skipped": 0,
            "failures": [],
            "duration": 0,
            "output": ""
        }
        
        try:
            start_time = datetime.now()
            proc = subprocess.run(
                cmd,
                cwd=work_dir,
                capture_output=True,
                text=True,
                timeout=300,
                shell=shell,
                env=self.env,
                encoding='utf-8',
                errors='replace'
            )
            result["output"] = proc.stdout
            result["duration"] = (datetime.now() - start_time).total_seconds()
            
            if proc.returncode == 0:
                result["passed"] = 1
            else:
                result["failed"] = 1
                result["failures"].append({
                    "target": target,
                    "test": "command",
                    "error": (proc.stderr or proc.stdout)[:500]
                })
                
        except Exception as e:
            result["failed"] = 1
            result["failures"].append({"target": target, "test": "error", "error": str(e)})
        
        return result
