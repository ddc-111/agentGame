"""自循环Agent配置"""
import os
from pathlib import Path

# 项目根目录
PROJECT_ROOT = Path(__file__).parent.parent

# 三端目录
SERVER_DIR = PROJECT_ROOT / "server"
CLIENT_DIR = PROJECT_ROOT / "client"
GM_DIR = PROJECT_ROOT / "gm"

# Agent工作目录
AGENT_DIR = PROJECT_ROOT / "agent"
REPORTS_DIR = AGENT_DIR / "reports"
TASKS_DIR = AGENT_DIR / "tasks"
HISTORY_DIR = AGENT_DIR / "history"
REQUIREMENTS_DIR = AGENT_DIR / "requirements"

# 测试配置
TEST_TIMEOUT = 300  # 测试超时时间（秒）
MAX_RETRY = 3       # 最大重试次数

# LLM配置（用于需求生成）
LLM_API_URL = os.getenv("LLM_API_URL", "https://api.openai.com/v1")
LLM_API_KEY = os.getenv("LLM_API_KEY", "")
LLM_MODEL = os.getenv("LLM_MODEL", "gpt-4")

# 质量阈值
COVERAGE_THRESHOLD = {
    "server": 60,
    "client": 50,
    "gm": 50
}

# 任务优先级
PRIORITY_CRITICAL = "critical"
PRIORITY_HIGH = "high"
PRIORITY_MEDIUM = "medium"
PRIORITY_LOW = "low"

# 测试类型
TEST_UNIT = "unit"
TEST_INTEGRATION = "integration"
TEST_E2E = "e2e"
TEST_REGRESSION = "regression"

# 确保目录存在
REPORTS_DIR.mkdir(parents=True, exist_ok=True)
TASKS_DIR.mkdir(parents=True, exist_ok=True)
HISTORY_DIR.mkdir(parents=True, exist_ok=True)
REQUIREMENTS_DIR.mkdir(parents=True, exist_ok=True)
