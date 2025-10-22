@echo off
REM Script para demonstração rápida - Windows

echo ================================
echo   Demo - Jogo Multiplayer Go
echo ================================
echo.

echo Verificando binarios...
if not exist "cmd\server\server.exe" (
    echo Servidor nao encontrado. Executando build...
    call build_windows.bat
    if %ERRORLEVEL% NEQ 0 (
        echo [ERRO] Falha no build
        pause
        exit /b 1
    )
)

echo.
echo Iniciando servidor em nova janela...
start "Servidor - Jogo Multiplayer" cmd /k "cd cmd\server && server.exe"

timeout /t 2 /nobreak > nul

echo Iniciando cliente 1 em nova janela...
start "Cliente 1" cmd /k "cd cmd\client && client.exe"

timeout /t 1 /nobreak > nul

echo Iniciando cliente 2 em nova janela...
start "Cliente 2" cmd /k "cd cmd\client && client.exe"

echo.
echo ================================
echo   Demo iniciada!
echo ================================
echo.
echo Foram abertas 3 janelas:
echo   1. Servidor (porta 8080)
echo   2. Cliente 1
echo   3. Cliente 2
echo.
echo Use WASD para mover e ESC para sair
echo.
pause

