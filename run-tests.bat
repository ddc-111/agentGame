@echo off
chcp 65001 >nul
echo ========================================
echo    AgentGame 自动化测试套件
echo ========================================
echo.

setlocal enabledelayedexpansion
set "ROOT=%~dp0"
set ERROR_COUNT=0

echo [1/5] 检查服务端环境...
if not exist "%ROOT%server\go.mod" (
    echo 错误: 未找到 go.mod 文件
    set /a ERROR_COUNT+=1
    goto :check_client
)
echo      服务端环境正常

:check_client
echo [2/5] 检查客户端环境...
if not exist "%ROOT%client\package.json" (
    echo 错误: 未找到 package.json 文件
    set /a ERROR_COUNT+=1
    goto :run_tests
)
echo      客户端环境正常

:run_tests

echo.
echo [3/5] 运行Go单元测试...
echo ----------------------------------------
cd /d "%ROOT%server"
go test ./internal/tests/... -v -count=1
if %ERRORLEVEL% NEQ 0 (
    echo Go测试失败!
    set /a ERROR_COUNT+=1
) else (
    echo Go测试通过!
)
cd /d "%ROOT%"

echo.
echo [4/5] 运行Go性能测试...
echo ----------------------------------------
cd /d "%ROOT%server"
go test ./internal/tests/... -bench=. -benchmem -count=1
if %ERRORLEVEL% NEQ 0 (
    echo 性能测试失败!
    set /a ERROR_COUNT+=1
) else (
    echo 性能测试完成!
)
cd /d "%ROOT%"

echo.
echo [5/5] 检查构建...
echo ----------------------------------------
cd /d "%ROOT%server"
go build -o bin/gameserver.exe cmd/gameserver/main.go
if %ERRORLEVEL% NEQ 0 (
    echo 服务端构建失败!
    set /a ERROR_COUNT+=1
) else (
    echo 服务端构建成功
)
cd /d "%ROOT%"

cd /d "%ROOT%client"
call npm run build 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 客户端构建失败或未安装依赖
) else (
    echo 客户端构建成功
)
cd /d "%ROOT%"

echo.
echo ========================================
echo    测试结果汇总
echo ========================================

if %ERROR_COUNT% EQU 0 (
    echo 所有测试通过!
) else (
    echo 发现 %ERROR_COUNT% 个错误
)

echo.
echo 浏览器测试工具:
echo   - API测试: client/test/browser-test.html
echo   - 视觉测试: client/test/game-test.html
echo.
echo 使用方法:
echo   1. 启动服务端: cd server ^&^& go run cmd/gameserver/main.go
echo   2. 启动客户端: cd client ^&^& npm run dev
echo   3. 在浏览器中打开测试工具
echo ========================================

pause
