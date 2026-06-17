@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   AgentGame 三端测试运行器
echo ========================================
echo.

set SERVER_RESULT=0
set CLIENT_RESULT=0
set GM_RESULT=0

if "%1"=="" (
    call :run_all
) else if "%1"=="server" (
    call :run_server
) else if "%1"=="client" (
    call :run_client
) else if "%1"=="gm" (
    call :run_gm
) else if "%1"=="coverage" (
    call :run_coverage
) else (
    echo 用法: test.bat [server^|client^|gm^|coverage]
    echo.
    echo   server    - 运行服务端测试
    echo   client    - 运行客户端测试
    echo   gm        - 运行GM管理端测试
    echo   coverage  - 运行所有测试并生成覆盖率报告
    echo   (无参数)  - 运行所有端测试
    exit /b 1
)

echo.
echo ========================================
echo   测试结果汇总
echo ========================================
if !SERVER_RESULT! equ 0 (echo   Server: PASS) else (echo   Server: FAIL)
if !CLIENT_RESULT! equ 0 (echo   Client: PASS) else (echo   Client: FAIL)
if !GM_RESULT! equ 0 (echo   GM:     PASS) else (echo   GM:     FAIL)
echo ========================================

if !SERVER_RESULT! neq 0 exit /b 1
if !CLIENT_RESULT! neq 0 exit /b 1
if !GM_RESULT! neq 0 exit /b 1

exit /b 0

:run_all
    echo [1/3] 运行服务端测试...
    call :run_server
    echo.
    echo [2/3] 运行客户端测试...
    call :run_client
    echo.
    echo [3/3] 运行GM管理端测试...
    call :run_gm
    goto :eof

:run_server
    echo --- Server (Go) ---
    cd /d "%~dp0server"
    go test ./internal/... -v -count=1
    set SERVER_RESULT=!ERRORLEVEL!
    cd /d "%~dp0"
    goto :eof

:run_client
    echo --- Client (Phaser) ---
    cd /d "%~dp0client"
    call npm test
    set CLIENT_RESULT=!ERRORLEVEL!
    cd /d "%~dp0"
    goto :eof

:run_gm
    echo --- GM Editor (Vue) ---
    cd /d "%~dp0gm"
    call npm test
    set GM_RESULT=!ERRORLEVEL!
    cd /d "%~dp0"
    goto :eof

:run_coverage
    echo --- 运行覆盖率测试 ---
    echo.
    echo [1/3] Server 覆盖率...
    cd /d "%~dp0server"
    go test ./internal/... -coverprofile=coverage.out -covermode=atomic
    go tool cover -html=coverage.out -o coverage.html
    echo     报告: server/coverage.html
    echo.
    echo [2/3] Client 覆盖率...
    cd /d "%~dp0client"
    call npm test -- --coverage
    echo.
    echo [3/3] GM 覆盖率...
    cd /d "%~dp0gm"
    call npm test -- --coverage
    cd /d "%~dp0"
    echo.
    echo 覆盖率报告生成完成
    goto :eof
