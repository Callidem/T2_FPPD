@echo off
REM Script para compilar para Windows e Linux

echo ================================
echo   Build Multiplataforma
echo ================================
echo.

echo 1. Baixando dependencias...
go mod tidy
echo.

echo 2. Compilando para Windows...
go build -o cmd\server\server.exe cmd\server\main.go
go build -o cmd\client\client.exe cmd\client\main.go
echo [OK] Windows

echo.
echo 3. Compilando para Linux...
set GOOS=linux
set GOARCH=amd64
go build -o cmd\server\server cmd\server\main.go
go build -o cmd\client\client cmd\client\main.go
set GOOS=windows
echo [OK] Linux

echo.
echo ================================
echo   Build concluido!
echo ================================
echo.
echo Binarios criados:
echo   Windows: cmd\server\server.exe, cmd\client\client.exe
echo   Linux:   cmd\server\server, cmd\client\client
echo.
pause

