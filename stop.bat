@echo off
chcp 65001 >nul

echo ========================================
echo   古风RPG Agent游戏 - 停止所有服务
echo ========================================
echo.

echo 正在停止服务端...
taskkill /fi "WindowTitle eq GameServer*" /f >nul 2>&1
taskkill /fi "WindowTitle eq *gameserver*" /f >nul 2>&1

echo 正在停止客户端...
taskkill /fi "WindowTitle eq GameClient*" /f >nul 2>&1
taskkill /fi "WindowTitle eq *vite*" /f >nul 2>&1

echo 正在停止GM管理端...
taskkill /fi "WindowTitle eq GameGM*" /f >nul 2>&1

echo.
echo [完成] 所有服务已停止
pause
