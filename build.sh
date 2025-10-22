#!/bin/bash
# Script de build para Linux/Mac

echo "================================"
echo "  Build - Jogo Multiplayer Go"
echo "================================"
echo ""

echo "1. Baixando dependências..."
go mod tidy
echo "✓ Dependências atualizadas"
echo ""

echo "2. Compilando servidor..."
go build -o cmd/server/server cmd/server/main.go
if [ $? -eq 0 ]; then
    echo "✓ Servidor compilado: cmd/server/server"
else
    echo "✗ Erro ao compilar servidor"
    exit 1
fi
echo ""

echo "3. Compilando cliente..."
go build -o cmd/client/client cmd/client/main.go
if [ $? -eq 0 ]; then
    echo "✓ Cliente compilado: cmd/client/client"
else
    echo "✗ Erro ao compilar cliente"
    exit 1
fi
echo ""

echo "4. Compilando jogo single-player original..."
go build -o game main.go
if [ $? -eq 0 ]; then
    echo "✓ Single-player compilado: ./game"
else
    echo "✗ Erro ao compilar single-player"
    exit 1
fi
echo ""

echo "================================"
echo "  ✓ Build concluído com sucesso!"
echo "================================"
echo ""
echo "Para executar:"
echo "  Servidor:  cd cmd/server && ./server"
echo "  Cliente:   ./cmd/client/client"
echo "  Original:  ./game"
echo ""

