# ⚡ Início Rápido - 5 Minutos

## 🎯 Objetivo

Colocar o jogo multiplayer rodando em **menos de 5 minutos**.

---

## 📋 Pré-requisitos

- [x] Go 1.21+ instalado
- [x] Terminal com UTF-8
- [x] Windows/Linux/Mac

---

## 🚀 3 Passos Simples

### 1️⃣ Compilar (30 segundos)

**Windows:**
```bash
build_windows.bat
```

**Linux/Mac:**
```bash
chmod +x build.sh
./build.sh
```

**Resultado esperado:**
```
✓ Servidor compilado
✓ Cliente compilado
✓ Single-player compilado
```

---

### 2️⃣ Iniciar Servidor (10 segundos)

**Abra um terminal:**

**Windows:**
```bash
cd cmd\server
server.exe
```

**Linux/Mac:**
```bash
cd cmd/server
./server
```

**Resultado esperado:**
```
====================================
  SERVIDOR DE JOGO MULTIPLAYER
====================================
Servidor iniciado na porta :8080
Aguardando conexões de clientes...
====================================
```

✅ Deixe este terminal aberto!

---

### 3️⃣ Iniciar Cliente(s) (10 segundos)

**Abra OUTRO terminal (ou vários):**

**Windows:**
```bash
cmd\client\client.exe
```

**Linux/Mac:**
```bash
./cmd/client/client
```

**Resultado esperado:**
- Interface gráfica do jogo aparece
- Personagem '☺' visível no mapa
- Mensagem: "Jogador XXX conectado com sucesso!"

✅ Repita para adicionar mais jogadores!

---

## 🎮 Como Jogar

| Tecla | Ação |
|-------|------|
| **W** | ⬆️ Cima |
| **A** | ⬅️ Esquerda |
| **S** | ⬇️ Baixo |
| **D** | ➡️ Direita |
| **E** | 🤝 Interagir |
| **ESC** | 🚪 Sair |

---

## 🎨 Elementos do Jogo

| Símbolo | O que é? |
|---------|----------|
| ☺ | **Você** (cinza) |
| ◉ | **Outros jogadores** (ciano) |
| ▤ | Parede (não pode passar) |
| ♣ | Vegetação (decoração) |
| ☠ | Inimigo (decoração) |

---

## 🔥 Demo Automática (Windows)

Quer ver funcionando **AGORA**?

```bash
run_demo.bat
```

Abre automaticamente:
- 1 servidor
- 2 clientes

Jogue imediatamente! 🎮

---

## 🐛 Problemas?

### "Servidor não inicia"
```bash
# Porta 8080 já está em uso?
netstat -ano | findstr :8080

# Mude a porta em cmd/server/main.go
# Linha: porta := ":8080"
```

### "Cliente não conecta"
```bash
# Servidor está rodando?
# Deve mostrar: "Aguardando conexões..."

# Endereço correto?
# Cliente usa: localhost:8080
```

### "Erro ao compilar"
```bash
# Instale dependências
go mod tidy

# Verifique versão do Go
go version
# Precisa: go1.21+
```

### "Interface não aparece"
```bash
# Use terminal moderno:
# Windows: Windows Terminal
# Mac: iTerm2 ou Terminal.app
# Linux: gnome-terminal, konsole, etc.
```

---

## 📚 Quer Saber Mais?

### Documentação Completa
- **`RESUMO_PROJETO.md`** - Visão geral
- **`IMPLEMENTATION.md`** - Arquitetura técnica
- **`TESTING_GUIDE.md`** - Como testar
- **`README_MULTIPLAYER.md`** - Guia detalhado

### Entender o Código
- **`cmd/server/main.go`** - Código do servidor
- **`cmd/client/main.go`** - Código do cliente
- **`pkg/game/`** - Biblioteca compartilhada

---

## ✅ Checklist de Sucesso

Você completou o quick start se conseguiu:

- [x] Compilar sem erros
- [x] Iniciar servidor (vê mensagem de "Aguardando")
- [x] Conectar cliente (vê interface gráfica)
- [x] Mover personagem com WASD
- [x] Ver logs no servidor

**Parabéns! 🎉**

---

## 🎓 O que Aconteceu?

### Servidor
1. Iniciou na porta `:8080`
2. Aguarda conexões TCP
3. Registra jogadores
4. Gerencia estado compartilhado
5. **NÃO tem o mapa** (só posições!)

### Cliente
1. Conectou ao `localhost:8080`
2. Carregou mapa de `mapa.txt`
3. Recebeu ID do servidor
4. Iniciou goroutine de atualização (500ms)
5. Renderiza interface com termbox

### Comunicação
1. Cliente → Servidor: **Conectar()**
2. Servidor → Cliente: **JogadorID**
3. Loop:
   - Cliente → Servidor: **ProcessarComando()** (quando move)
   - Cliente → Servidor: **ObterEstado()** (a cada 500ms)
   - Cliente renderiza outros jogadores

---

## 🔥 Experimente Agora!

### Teste Multiplayer

**Abra 3 terminais:**

**Terminal 1 (Servidor):**
```bash
cd cmd/server && server.exe
```

**Terminal 2 (Cliente 1):**
```bash
cmd\client\client.exe
```

**Terminal 3 (Cliente 2):**
```bash
cmd\client\client.exe
```

**Resultado:**
- Cliente 1 vê Cliente 2 como '◉'
- Cliente 2 vê Cliente 1 como '◉'
- Movimentos são sincronizados!

---

## 🚀 Próximos Passos

Agora que está funcionando:

1. **Teste multiplayer** com 3+ clientes
2. **Leia os logs** do servidor
3. **Tente mapas diferentes**: `client.exe maze.txt`
4. **Leia documentação**: `IMPLEMENTATION.md`
5. **Execute testes**: `TESTING_GUIDE.md`

---

## 💡 Dicas

### Performance
- Servidor suporta 10+ clientes simultâneos
- Atualização a cada 500ms (configurável)
- Retry automático: 3 tentativas

### Mapas
```bash
# Mapa padrão
client.exe

# Labirinto
client.exe maze.txt

# Seu mapa customizado
client.exe seu_mapa.txt
```

### Customização
```go
// cmd/client/main.go

// Mudar intervalo de atualização
cliente.IniciarAtualizacaoPeriodica(200 * time.Millisecond) // 200ms

// Mudar servidor
enderecoServidor := "192.168.1.10:8080" // Servidor remoto
```

---

## 🎯 Características Implementadas

Mesmo simples, o projeto tem:

- ✅ **Exactly-once** (comandos não duplicam)
- ✅ **Retry automático** (3 tentativas)
- ✅ **Goroutine** (atualização periódica)
- ✅ **Thread-safe** (mutexes)
- ✅ **Logging** (debug fácil)
- ✅ **Multiplayer** (N jogadores)

Tudo em ~800 linhas de código Go! 🔥

---

## 📞 Ajuda Rápida

| Problema | Solução |
|----------|---------|
| Não compila | `go mod tidy` |
| Servidor não inicia | Porta 8080 livre? |
| Cliente não conecta | Servidor rodando? |
| Interface bugada | Terminal com UTF-8 |
| Jogador não move | Pressione WASD (não setas) |
| Outros não aparecem | Aguarde 500ms (atualização) |

---

## ✨ Está Funcionando?

**SIM** ✅ → Parabéns! Explore o código.  
**NÃO** ❌ → Veja seção "Problemas?" acima.

---

## 🎉 Parabéns!

Você agora tem um jogo multiplayer funcional em Go!

**Explore**:
- Adicione mais clientes
- Leia o código
- Modifique e experimente
- Leia a documentação completa

**Divirta-se! 🎮**

---

**Tempo total**: ~5 minutos  
**Dificuldade**: ⭐☆☆☆☆ (Muito fácil)  
**Diversão**: ⭐⭐⭐⭐⭐ (Muito divertido!)

---

**Dúvidas?** → Leia `IMPLEMENTATION.md`  
**Bugs?** → Veja `TESTING_GUIDE.md`  
**Visão geral?** → Leia `RESUMO_PROJETO.md`

