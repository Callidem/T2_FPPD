@echo off
REM Script para conectar a um servidor (local ou remoto)

echo ================================
echo   Cliente - Conectar a Servidor
echo ================================
echo.
echo Exemplos de uso:
echo   - localhost (padrao)
echo   - 192.168.1.100
echo   - servidor.local
echo.

set /p SERVER_IP="Digite o IP do servidor (ou Enter para localhost): "

if "%SERVER_IP%"=="" (
    set SERVER_IP=localhost
)

echo.
echo Conectando ao servidor: %SERVER_IP%:8080
echo.

cd cmd\client
client.exe mapa.txt %SERVER_IP%:8080

pause

