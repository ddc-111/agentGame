@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   古风RPG Agent游戏 - 一键启动脚本
echo ========================================
echo.

set ROOT_DIR=%~dp0
set SERVER_DIR=%ROOT_DIR%server
set CLIENT_DIR=%ROOT_DIR%client
set GM_DIR=%ROOT_DIR%gm

:: 检查是否首次运行
if not exist "%SERVER_DIR%\config.yaml" (
    echo [提示] 未检测到服务端配置文件，正在从模板创建...
    copy "%SERVER_DIR%\config.example.yaml" "%SERVER_DIR%\config.yaml" >nul
    echo [提示] 已创建 config.yaml，请根据需要修改配置
    echo.
)

:: 菜单
echo 请选择操作：
echo.
echo   [1] 一键启动全部（服务端 + 客户端 + GM）
echo   [2] 仅启动服务端
echo   [3] 仅启动客户端
echo   [4] 仅启动GM管理端
echo   [5] 编译服务端
echo   [6] 安装前端依赖
echo   [7] 退出
echo.
set /p choice=请输入选择 (1-7): 

if "%choice%"=="1" goto start_all
if "%choice%"=="2" goto start_server
if "%choice%"=="3" goto start_client
if "%choice%"=="4" goto start_gm
if "%choice%"=="5" goto build_server
if "%choice%"=="6" goto install_deps
if "%choice%"=="7" goto end
echo [错误] 无效选择
pause
goto end

:start_all
echo.
echo [1/3] 启动服务端...
start "GameServer" cmd /k "cd /d %SERVER_DIR% && go run cmd/gameserver/main.go"
timeout /t 3 /nobreak >nul

echo [2/3] 启动客户端...
start "GameClient" cmd /k "cd /d %CLIENT_DIR% && npm run dev"
timeout /t 2 /nobreak >nul

echo [3/3] 启动GM管理端...
start "GameGM" cmd /k "cd /d %GM_DIR% && npm run dev"

echo.
echo ========================================
echo   全部服务已启动！
echo   服务端: http://localhost:8080
echo   客户端: http://localhost:5173
echo   GM管理端: http://localhost:5174
echo ========================================
pause
goto end

:start_server
echo.
echo 正在启动服务端...
cd /d %SERVER_DIR%
go run cmd/gameserver/main.go
pause
goto end

:start_client
echo.
echo 正在启动客户端...
cd /d %CLIENT_DIR%
npm run dev
pause
goto end

:start_gm
echo.
echo 正在启动GM管理端...
cd /d %GM_DIR%
npm run dev
pause
goto end

:build_server
echo.
echo 正在编译服务端...
cd /d %SERVER_DIR%

:: 检查是否有go.mod
if not exist "go.mod" (
    echo [错误] 未找到 go.mod
    pause
    goto end
)

:: 下载依赖
echo [1/2] 下载依赖...
go mod tidy
if errorlevel 1 (
    echo [错误] 依赖下载失败
    pause
    goto end
)

:: 编译
echo [2/2] 编译中...
go build -o bin/gameserver.exe cmd/gameserver/main.go
if errorlevel 1 (
    echo [错误] 编译失败
    pause
    goto end
)

echo.
echo [成功] 编译完成: %SERVER_DIR%\bin\gameserver.exe
pause
goto end

:install_deps
echo.
echo 正在安装前端依赖...

echo [1/2] 安装客户端依赖...
cd /d %CLIENT_DIR%
call npm install
if errorlevel 1 (
    echo [警告] 客户端依赖安装可能有问题
)

echo.
echo [2/2] 安装GM管理端依赖...
cd /d %GM_DIR%
call npm install
if errorlevel 1 (
    echo [警告] GM管理端依赖安装可能有问题
)

echo.
echo [完成] 前端依赖安装完成
pause
goto end

:end
