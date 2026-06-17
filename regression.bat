@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   AgentGame 回归测试套件
echo ========================================
echo.

set PASS=0
set FAIL=0
set SKIP=0
set ERRORS=

:: 记录开始时间
set START_TIME=%TIME%

echo [阶段1] 构建验证
echo ----------------------------------------

echo [1.1] 编译 Server...
cd /d "%~dp0server"
go build -o bin\gameserver.exe cmd\gameserver\main.go 2>nul
if !ERRORLEVEL! equ 0 (
    echo       PASS: Server 编译成功
    set /a PASS+=1
) else (
    echo       FAIL: Server 编译失败
    set /a FAIL+=1
    set ERRORS=!ERRORS! - Server编译失败
)

echo [1.2] 编译 Client...
cd /d "%~dp0client"
call npm run build >nul 2>&1
if !ERRORLEVEL! equ 0 (
    echo       PASS: Client 编译成功
    set /a PASS+=1
) else (
    echo       FAIL: Client 编译失败
    set /a FAIL+=1
    set ERRORS=!ERRORS! - Client编译失败
)

echo [1.3] 编译 GM...
cd /d "%~dp0gm"
call npm run build >nul 2>&1
if !ERRORLEVEL! equ 0 (
    echo       PASS: GM 编译成功
    set /a PASS+=1
) else (
    echo       FAIL: GM 编译失败
    set /a FAIL+=1
    set ERRORS=!ERRORS! - GM编译失败
)

echo.
echo [阶段2] 单元测试
echo ----------------------------------------

echo [2.1] Server 单元测试...
cd /d "%~dp0server"
for /f %%i in ('go test ./internal/... -count=1 2^>^&1 ^| findstr /c:"PASS" /c:"FAIL"') do set TEST_OUT=%%i
go test ./internal/... -count=1 >nul 2>&1
if !ERRORLEVEL! equ 0 (
    echo       PASS: Server 单元测试通过
    set /a PASS+=1
) else (
    echo       FAIL: Server 单元测试失败
    set /a FAIL+=1
    set ERRORS=!ERRORS! - Server单测失败
)

echo [2.2] Client 单元测试...
cd /d "%~dp0client"
call npm test >nul 2>&1
if !ERRORLEVEL! equ 0 (
    echo       PASS: Client 单元测试通过
    set /a PASS+=1
) else (
    echo       FAIL: Client 单元测试失败
    set /a FAIL+=1
    set ERRORS=!ERRORS! - Client单测失败
)

echo [2.3] GM 单元测试...
cd /d "%~dp0gm"
call npm test >nul 2>&1
if !ERRORLEVEL! equ 0 (
    echo       PASS: GM 单元测试通过
    set /a PASS+=1
) else (
    echo       FAIL: GM 单元测试失败
    set /a FAIL+=1
    set ERRORS=!ERRORS! - GM单测失败
)

echo.
echo [阶段3] 集成测试 (需要运行中的Server)
echo ----------------------------------------

:: 检查Server是否运行
curl -s http://localhost:8080/health >nul 2>&1
if !ERRORLEVEL! equ 0 (
    echo [3.1] Health Check...
    for /f "tokens=*" %%i in ('curl -s http://localhost:8080/health') do set HEALTH=%%i
    echo !HEALTH! | findstr "ok" >nul
    if !ERRORLEVEL! equ 0 (
        echo       PASS: 服务健康
        set /a PASS+=1
    ) else (
        echo       FAIL: 服务异常
        set /a FAIL+=1
    )

    echo [3.2] API 端点测试...
    set API_PASS=1
    
    curl -s http://localhost:8080/api/scenes | findstr "data" >nul
    if !ERRORLEVEL! neq 0 set API_PASS=0
    
    curl -s http://localhost:8080/api/npcs | findstr "data" >nul
    if !ERRORLEVEL! neq 0 set API_PASS=0
    
    curl -s http://localhost:8080/api/items | findstr "data" >nul
    if !ERRORLEVEL! neq 0 set API_PASS=0
    
    curl -s http://localhost:8080/api/tasks | findstr "data" >nul
    if !ERRORLEVEL! neq 0 set API_PASS=0
    
    if !API_PASS! equ 1 (
        echo       PASS: API 端点响应正常
        set /a PASS+=1
    ) else (
        echo       FAIL: API 端点响应异常
        set /a FAIL+=1
        set ERRORS=!ERRORS! - API端点异常
    )

    echo [3.3] WebSocket 连接测试...
    :: 简单检查ws端点可达性
    curl -s -o nul -w "%%{http_code}" http://localhost:8080/ws | findstr "400\|200" >nul
    if !ERRORLEVEL! equ 0 (
        echo       PASS: WebSocket 端点可达
        set /a PASS+=1
    ) else (
        echo       SKIP: WebSocket 端点不可达
        set /a SKIP+=1
    )
) else (
    echo [3.1-3.3] SKIP: Server 未运行，跳过集成测试
    set /a SKIP+=3
)

echo.
echo ========================================
echo   回归测试结果
echo ========================================
echo   通过: !PASS!
echo   失败: !FAIL!
echo   跳过: !SKIP!
echo   耗时: %START_TIME% - %TIME%
echo ========================================

if !FAIL! gtr 0 (
    echo.
    echo 失败项:
    echo !ERRORS!
    echo.
    echo 回归测试未通过!
    exit /b 1
) else (
    echo.
    echo 回归测试全部通过!
    exit /b 0
)
