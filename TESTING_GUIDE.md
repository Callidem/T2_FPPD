# Guia de Testes - Jogo Multiplayer

## 🧪 Objetivo

Este guia fornece instruções detalhadas para testar todas as funcionalidades implementadas no jogo multiplayer.

---

## 📋 Pré-requisitos

1. Go 1.21+ instalado
2. Terminal com suporte a UTF-8
3. Binários compilados (use `build_windows.bat` ou `build.sh`)

---

## 🚀 Cenários de Teste

### Teste 1: Compilação e Build

**Objetivo**: Verificar se o projeto compila sem erros

**Passos**:
```bash
# Windows
build_windows.bat

# Linux/Mac
./build.sh
```

**Resultado Esperado**:
```
✓ Servidor compilado: cmd/server/server.exe
✓ Cliente compilado: cmd/client/client.exe
✓ Single-player compilado: game.exe
```

**Status**: [ ] Passou [ ] Falhou

---

### Teste 2: Iniciar Servidor

**Objetivo**: Verificar se o servidor inicia corretamente

**Passos**:
```bash
cd cmd/server
./server.exe    # Windows
./server        # Linux/Mac
```

**Resultado Esperado**:
```
====================================
  SERVIDOR DE JOGO MULTIPLAYER
====================================
Servidor iniciado na porta :8080
Aguardando conexões de clientes...
====================================
```

**Status**: [ ] Passou [ ] Falhou

---

### Teste 3: Conectar Cliente Único

**Objetivo**: Verificar conexão de um cliente ao servidor

**Passos**:
1. Iniciar servidor (Teste 2)
2. Em outro terminal:
```bash
cd ../../
./cmd/client/client.exe
```

**Resultado Esperado (Servidor)**:
```
[CONEXÃO] Nova conexão estabelecida de 127.0.0.1:XXXXX
[REQUISICAO] Conectar - Nome: Jogador_XXX, Posição: (X, Y)
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_1
```

**Resultado Esperado (Cliente)**:
- Interface gráfica aparece
- Personagem '☺' visível no mapa
- Mensagem de status: "Jogador XXX conectado com sucesso!"

**Status**: [ ] Passou [ ] Falhou

---

### Teste 4: Movimentação Básica

**Objetivo**: Verificar movimentação do personagem e sincronização com servidor

**Passos**:
1. Com cliente conectado (Teste 3)
2. Pressione teclas WASD para mover
3. Observe logs do servidor

**Resultado Esperado (Servidor)**:
```
[REQUISICAO] ProcessarComando - JogadorID: jogador_1, Tipo: mover, Direção: d, SeqNum: 1
[RESPOSTA] ProcessarComando - Sucesso: true, Mensagem: Jogador movido para (X, Y)
[REQUISICAO] ObterEstado - JogadorID: jogador_1
[RESPOSTA] ObterEstado - Total de jogadores ativos: 1
```

**Resultado Esperado (Cliente)**:
- Personagem se move suavemente
- Colisões com paredes são detectadas (personagem não atravessa '▤')
- Barra de status atualiza a mensagem

**Status**: [ ] Passou [ ] Falhou

---

### Teste 5: Múltiplos Clientes Simultâneos

**Objetivo**: Verificar suporte a múltiplos jogadores

**Passos**:
1. Iniciar servidor
2. Abrir 3 terminais
3. Executar cliente em cada terminal:
```bash
Terminal 1: ./cmd/client/client.exe
Terminal 2: ./cmd/client/client.exe
Terminal 3: ./cmd/client/client.exe
```

**Resultado Esperado (Servidor)**:
```
[CONEXÃO] Nova conexão estabelecida de ...
[REQUISICAO] Conectar - Nome: Jogador_XXX, ...
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_1

[CONEXÃO] Nova conexão estabelecida de ...
[REQUISICAO] Conectar - Nome: Jogador_YYY, ...
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_2

[CONEXÃO] Nova conexão estabelecida de ...
[REQUISICAO] Conectar - Nome: Jogador_ZZZ, ...
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_3

[RESPOSTA] ObterEstado - Total de jogadores ativos: 3
```

**Resultado Esperado (Cada Cliente)**:
- Vê o próprio personagem como '☺' (cinza)
- Vê outros jogadores como '◉' (ciano)
- Atualização em tempo real quando outros jogadores se movem

**Status**: [ ] Passou [ ] Falhou

---

### Teste 6: Sincronização em Tempo Real

**Objetivo**: Verificar atualização periódica do estado

**Passos**:
1. Com 2 clientes conectados (Teste 5)
2. No Cliente 1: Mover com WASD
3. Observar Cliente 2

**Resultado Esperado**:
- Cliente 2 vê o personagem do Cliente 1 se mover
- Atualização ocorre em ~500ms (pode ter pequeno delay)
- Posição é consistente entre cliente e servidor

**Status**: [ ] Passou [ ] Falhou

---

### Teste 7: Exactly-Once Semantics

**Objetivo**: Verificar que comandos não são executados duas vezes

**Passos**:
1. Conectar 1 cliente
2. Simular latência de rede (opcional: desabilitar Wi-Fi temporariamente)
3. Mover personagem rapidamente
4. Observar logs do servidor

**Resultado Esperado (Servidor)**:
Se houver retry, deve aparecer:
```
[RESPOSTA] ProcessarComando - Comando duplicado ignorado (SeqNum: X)
```

**Verificação**:
- Personagem não "pula" posições
- Cada SeqNum é processado apenas uma vez

**Status**: [ ] Passou [ ] Falhou

---

### Teste 8: Tratamento de Erro com Retry

**Objetivo**: Verificar reconexão automática

**Passos**:
1. Iniciar servidor e cliente
2. **Fechar o servidor** (Ctrl+C)
3. No cliente, tentar mover
4. **Reiniciar o servidor**

**Resultado Esperado (Cliente)**:
```
Erro na chamada RPC ProcessarComando: ...
Tentativa 2/3 para ServidorJogo.ProcessarComando...
Tentativa 3/3 para ServidorJogo.ProcessarComando...
Falha após 3 tentativas: ...
```

**Após reiniciar servidor**:
- Cliente não conecta automaticamente (precisa reiniciar cliente)
- Isso é esperado nesta versão

**Status**: [ ] Passou [ ] Falhou

---

### Teste 9: Interação

**Objetivo**: Verificar comando de interação

**Passos**:
1. Conectar 1 cliente
2. Pressionar tecla **E**

**Resultado Esperado (Servidor)**:
```
[REQUISICAO] ProcessarComando - JogadorID: jogador_1, Tipo: interagir, Direção: , SeqNum: X
[RESPOSTA] ProcessarComando - Sucesso: true, Mensagem: Interação registrada em (X, Y)
```

**Resultado Esperado (Cliente)**:
- Barra de status mostra: "Interagindo em (X, Y)"

**Status**: [ ] Passou [ ] Falhou

---

### Teste 10: Desconexão Limpa

**Objetivo**: Verificar desconexão correta do cliente

**Passos**:
1. Conectar 2 clientes
2. No Cliente 1, pressionar **ESC**
3. Observar servidor e Cliente 2

**Resultado Esperado (Servidor)**:
```
[REQUISICAO] ProcessarComando - JogadorID: jogador_1, Tipo: desconectar, ...
[RESPOSTA] ProcessarComando - Sucesso: true, Mensagem: Jogador desconectado
```

**Resultado Esperado (Cliente 2)**:
- Jogador 1 ainda aparece temporariamente (marcado como inativo)
- Após próxima atualização, pode desaparecer (depende da implementação)

**Status**: [ ] Passou [ ] Falhou

---

### Teste 11: Mapas Diferentes

**Objetivo**: Verificar que cada cliente pode ter mapa diferente

**Passos**:
1. Iniciar servidor
2. Cliente 1:
```bash
./cmd/client/client.exe mapa.txt
```
3. Cliente 2:
```bash
./cmd/client/client.exe maze.txt
```

**Resultado Esperado**:
- Cliente 1 vê mapa de `mapa.txt`
- Cliente 2 vê mapa de `maze.txt`
- Ambos se conectam ao mesmo servidor
- Posições são compartilhadas, mas visualização é diferente

**Nota**: Isso pode causar inconsistências (ex: jogador 1 atravessando parede do jogador 2). É esperado nesta versão, pois servidor NÃO valida mapa.

**Status**: [ ] Passou [ ] Falhou

---

### Teste 12: Colisão com Parede

**Objetivo**: Verificar que cliente valida colisões localmente

**Passos**:
1. Conectar 1 cliente
2. Tentar mover para uma parede '▤'

**Resultado Esperado (Cliente)**:
- Personagem NÃO se move
- Nenhum comando é enviado ao servidor

**Resultado Esperado (Servidor)**:
- Nenhuma requisição `ProcessarComando` recebida

**Verificação**: Servidor NÃO valida colisões (correto!)

**Status**: [ ] Passou [ ] Falhou

---

### Teste 13: Performance com Muitos Jogadores

**Objetivo**: Verificar comportamento com 10+ clientes

**Passos**:
1. Iniciar servidor
2. Abrir 10 terminais e executar cliente em cada um

**Resultado Esperado**:
- Todos os clientes conectam
- Servidor mantém 10 jogadores ativos
- Atualização continua funcionando (pode ter mais latência)

**Observações**:
- Renderização pode ficar lenta com muitos jogadores sobrepostos
- Servidor deve continuar responsivo

**Status**: [ ] Passou [ ] Falhou

---

### Teste 14: Logging Completo

**Objetivo**: Verificar que servidor loga todas as operações

**Passos**:
1. Executar Testes 3-10
2. Revisar logs do servidor

**Resultado Esperado**:
- Toda requisição tem log `[REQUISICAO]`
- Toda resposta tem log `[RESPOSTA]`
- Conexões são logadas com `[CONEXÃO]`
- Comandos duplicados são identificados

**Status**: [ ] Passou [ ] Falhou

---

### Teste 15: Demo Rápida (Windows)

**Objetivo**: Testar script de demonstração automática

**Passos**:
```bash
run_demo.bat
```

**Resultado Esperado**:
- 3 janelas abrem automaticamente
- 1 servidor + 2 clientes
- Funcionamento normal

**Status**: [ ] Passou [ ] Falhou

---

## 📊 Resumo dos Testes

| # | Teste | Passou | Falhou | N/A |
|---|-------|--------|--------|-----|
| 1 | Compilação | [ ] | [ ] | [ ] |
| 2 | Iniciar Servidor | [ ] | [ ] | [ ] |
| 3 | Conectar Cliente | [ ] | [ ] | [ ] |
| 4 | Movimentação | [ ] | [ ] | [ ] |
| 5 | Múltiplos Clientes | [ ] | [ ] | [ ] |
| 6 | Sincronização | [ ] | [ ] | [ ] |
| 7 | Exactly-Once | [ ] | [ ] | [ ] |
| 8 | Retry | [ ] | [ ] | [ ] |
| 9 | Interação | [ ] | [ ] | [ ] |
| 10 | Desconexão | [ ] | [ ] | [ ] |
| 11 | Mapas Diferentes | [ ] | [ ] | [ ] |
| 12 | Colisão | [ ] | [ ] | [ ] |
| 13 | Performance | [ ] | [ ] | [ ] |
| 14 | Logging | [ ] | [ ] | [ ] |
| 15 | Demo | [ ] | [ ] | [ ] |

**Total**: ___/15 testes passaram

---

## 🐛 Bugs Conhecidos

Nenhum bug conhecido no momento. Reporte bugs encontrados durante os testes.

---

## 📝 Notas de Teste

### Limitações Esperadas

1. **Reconexão**: Cliente não reconecta automaticamente após queda do servidor
2. **Mapas diferentes**: Podem causar inconsistências visuais
3. **Colisão entre jogadores**: Não implementada (jogadores podem sobrepor)
4. **Escalabilidade**: Testado com até 10 clientes simultâneos

### Comportamentos Corretos

1. **Servidor NÃO valida mapa**: Correto conforme requisito
2. **Cliente valida colisões**: Correto conforme requisito
3. **Pequeno delay na sincronização**: Normal devido a intervalo de 500ms

---

## 📧 Reportar Problemas

Se encontrar bugs durante os testes:

1. Anote o número do teste
2. Descreva o comportamento esperado vs obtido
3. Inclua logs do servidor e cliente
4. Informe sistema operacional e versão do Go

---

**Última atualização**: 22/10/2025  
**Desenvolvido para**: T2 - FPPD

