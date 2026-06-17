@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   AgentGame 自循环Agent
echo ========================================
echo.

if "%1"=="--help" (
    echo 用法: run_agent.bat [选项]
    echo.
    echo 选项:
    echo   -n N              运行N次迭代 (默认1)
    echo   -v                详细输出
    echo   -q                安静模式
    echo   -c FILE           使用配置文件
    echo   --enable-llm      启用LLM自动执行
    echo   --api-key KEY     设置LLM API密钥
    echo   --execute-tasks   执行生成的任务
    echo   --report-only     只生成报告
    echo   --show-history    显示历史记录
    echo   --generate-config 生成默认配置文件
    echo.
    echo 示例:
    echo   run_agent.bat -n 3
    echo   run_agent.bat --enable-llm --api-key sk-xxx
    echo   run_agent.bat -c agent_config.yaml
    exit /b 0
)

python agent\main.py %*
