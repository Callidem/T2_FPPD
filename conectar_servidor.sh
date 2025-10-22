#!/bin/bash
# Script para conectar a um servidor (local ou remoto)

echo "================================"
echo "  Cliente - Conectar a Servidor"
echo "================================"
echo
echo "Exemplos de uso:"
echo "  - localhost (padrão)"
echo "  - 192.168.1.100"
echo "  - servidor.local"
echo

read -p "Digite o IP do servidor (ou Enter para localhost): " SERVER_IP

if [ -z "$SERVER_IP" ]; then
    SERVER_IP="localhost"
fi

echo
echo "Conectando ao servidor: $SERVER_IP:8080"
echo

cd cmd/client
./client mapa.txt $SERVER_IP:8080

