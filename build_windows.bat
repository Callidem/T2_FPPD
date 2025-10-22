@echo off
REM Script de build para Windows

echo ================================
echo   Build - Jogo Multiplayer Go
echo ================================
echo.

echo 1. Baixando dependencias...
go mod tidy
if %ERRORLEVEL% EQU 0 (
    echo [OK] Dependencias atualizadas
) else (
    echo [ERRO] Falha ao atualizar dependencias
    exit /b 1
)
echo.

echo 2. Compilando servidor...
go build -o cmd\server\server.exe cmd\server\main.go
if %ERRORLEVEL% EQU 0 (
    echo [OK] Servidor compilado: cmd\server\server.exe
) else (
    echo [ERRO] Falha ao compilar servidor
    exit /b 1
)
echo.

echo 3. Compilando cliente...
go build -o cmd\client\client.exe cmd\client\main.go
if %ERRORLEVEL% EQU 0 (
    echo [OK] Cliente compilado: cmd\client\client.exe
) else (
    echo [ERRO] Falha ao compilar cliente
    exit /b 1
)
echo.

echo 4. Compilando jogo single-player original...
go build -o game.exe main.go
if %ERRORLEVEL% EQU 0 (
    echo [OK] Single-player compilado: game.exe
) else (
    echo [ERRO] Falha ao compilar single-player
    exit /b 1
)
echo.

echo ================================
echo   Build concluido com sucesso!
echo ================================
echo.
echo Para executar:
echo   Servidor:  cd cmd\server ^&^& server.exe
echo   Cliente:   cmd\client\client.exe
echo   Original:  game.exe
echo.
pause

