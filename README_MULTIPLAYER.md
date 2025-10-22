# Jogo Multiplayer - Guia Rápido

## 🎮 Visão Geral

Jogo multiplayer em Go com arquitetura cliente-servidor, onde cada cliente mantém seu próprio mapa e o servidor gerencia o estado compartilhado dos jogadores.

## 🚀 Como Executar

### 1. Compilar

```bash
# No diretório raiz do projeto
go mod tidy
go build -o cmd/server/server.exe cmd/server/main.go
go build -o cmd/client/client.exe cmd/client/main.go
```

### 2. Iniciar o Servidor

```bash
cd cmd/server
./server.exe       # Windows
./server           # Linux/Mac
```

O servidor iniciará na porta `:8080` e exibirá:
```
====================================
  SERVIDOR DE JOGO MULTIPLAYER
====================================
Servidor iniciado na porta :8080
Aguardando conexões de clientes...
====================================
```

### 3. Iniciar Clientes

Em terminais separados (você pode abrir quantos quiser):

```bash
# Terminal 1 - Cliente 1
cd ../..
./cmd/client/client.exe

# Terminal 2 - Cliente 2
./cmd/client/client.exe

# Terminal 3 - Cliente 3
./cmd/client/client.exe maze.txt  # Com mapa diferente
```

## 🎯 Controles

- **W** - Mover para cima
- **A** - Mover para esquerda  
- **S** - Mover para baixo
- **D** - Mover para direita
- **E** - Interagir
- **ESC** - Sair

## 📊 O que Você Verá

- **Seu personagem**: `☺` (cinza)
- **Outros jogadores**: `◉` (ciano)
- **Paredes**: `▤` (preto)
- **Vegetação**: `♣` (verde)
- **Inimigos**: `☠` (vermelho)

## 🔧 Características Implementadas

✅ **Servidor sem mapa**: Apenas gerencia posições e estado dos jogadores  
✅ **Cliente com mapa local**: Cada cliente carrega seu próprio arquivo de mapa  
✅ **Exactly-once semantics**: Comandos não são executados mais de uma vez  
✅ **Retry automático**: Reconexão automática em caso de falha  
✅ **Atualização em tempo real**: Estado sincronizado a cada 500ms  
✅ **Logging detalhado**: Servidor exibe todas as requisições para debug  

## 📁 Estrutura

```
T2_FPPD/
├── cmd/
│   ├── server/       # Servidor multiplayer
│   └── client/       # Cliente multiplayer
├── pkg/game/         # Código compartilhado
├── mapa.txt         # Mapa principal
├── maze.txt         # Mapa alternativo
└── IMPLEMENTATION.md # Documentação completa
```

## 🐛 Debug

O servidor imprime todas as requisições e respostas:

```
[REQUISICAO] Conectar - Nome: Jogador_123, Posição: (4, 11)
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_1
[REQUISICAO] ProcessarComando - JogadorID: jogador_1, Tipo: mover, Direção: d, SeqNum: 1
[RESPOSTA] ProcessarComando - Sucesso: true, Mensagem: Jogador movido para (5, 11)
[REQUISICAO] ObterEstado - JogadorID: jogador_1
[RESPOSTA] ObterEstado - Total de jogadores ativos: 3
```

## 📚 Documentação Completa

Para entender em profundidade a implementação, veja:
- **`IMPLEMENTATION.md`**: Documentação técnica detalhada
  - Arquitetura do sistema
  - Protocolo de comunicação
  - Garantia de exactly-once
  - Fluxos e diagramas

## ⚙️ Configurações

- **Porta do servidor**: `8080` (modificar em `cmd/server/main.go`)
- **Endereço do servidor**: `localhost:8080` (modificar em `cmd/client/main.go`)
- **Intervalo de atualização**: `500ms` (modificar em `cmd/client/main.go`)
- **Retry timeout**: `500ms` com máximo de 3 tentativas

## 🎓 Requisitos do Projeto Atendidos

### Servidor
✅ Gerencia sessão de jogo  
✅ Mantém estado dos jogadores (posições, vidas)  
✅ NÃO mantém cópia do mapa  
✅ NÃO tem lógica de movimentação  
✅ Sem interface gráfica  
✅ Imprime requisições e respostas no terminal  

### Cliente
✅ Interface onde jogador interage  
✅ Controla lógica de movimentação  
✅ Conecta ao servidor para obter estado  
✅ Goroutine dedicada para atualização periódica  

### Comunicação
✅ Iniciada apenas pelos clientes  
✅ Servidor apenas responde  
✅ Tratamento de erro com retry automático  
✅ Garantia de exactly-once com sequenceNumber  

## 🛠️ Tecnologias

- **Go 1.21+**
- **RPC nativo** (net/rpc)
- **termbox-go** (interface de terminal)
- **Concorrência**: goroutines e mutexes

## 📝 Licença

Projeto acadêmico - FPPD (Fundamentos de Processamento Paralelo e Distribuído)

---

**Autor**: Desenvolvido para T2_FPPD  
**Data**: Outubro 2025

