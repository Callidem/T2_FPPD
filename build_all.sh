#!/bin/bash
# Script para compilar para Linux e Windows

echo "================================"
echo "  Build Multiplataforma"
echo "================================"
echo

echo "1. Baixando dependências..."
go mod tidy
echo

echo "2. Compilando para Linux..."
go build -o cmd/server/server cmd/server/main.go
go build -o cmd/client/client cmd/client/main.go
echo "[OK] Linux"

echo
echo "3. Compilando para Windows..."
GOOS=windows GOARCH=amd64 go build -o cmd/server/server.exe cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o cmd/client/client.exe cmd/client/main.go
echo "[OK] Windows"

echo
echo "================================"
echo "  Build concluído!"
echo "================================"
echo
echo "Binários criados:"
echo "  Linux:   cmd/server/server, cmd/client/client"
echo "  Windows: cmd/server/server.exe, cmd/client/client.exe"
echo

