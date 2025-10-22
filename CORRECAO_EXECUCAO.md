# 🔧 Correção - Como Executar o Cliente

## ❗ Problema Encontrado

O cliente precisa acessar o arquivo `mapa.txt`, que está no diretório raiz do projeto.

## ✅ Solução Aplicada

Os arquivos `mapa.txt` e `maze.txt` foram copiados para `cmd/client/`.

---

## 🚀 Como Executar Corretamente

### Opção 1: Do Diretório do Cliente

```bash
cd cmd\client
client.exe
```

### Opção 2: Do Diretório Raiz (Recomendado)

```bash
# Certifique-se de estar no diretório raiz
cd C:\Users\usrteia-0005\Documents\Faculdade\FPPD\T2_FPPD

# Execute o cliente
cd cmd\client
client.exe
```

### Opção 3: Demo Automática

```bash
# Do diretório raiz
run_demo.bat
```

Isso abre automaticamente servidor + 2 clientes!

---

## 🎮 Passo a Passo Completo

### 1. Iniciar Servidor

**Terminal 1:**
```bash
cd C:\Users\usrteia-0005\Documents\Faculdade\FPPD\T2_FPPD\cmd\server
server.exe
```

**Deve aparecer:**
```
====================================
  SERVIDOR DE JOGO MULTIPLAYER
====================================
Servidor iniciado na porta :8080
Aguardando conexões de clientes...
====================================
```

### 2. Iniciar Cliente(s)

**Terminal 2:**
```bash
cd C:\Users\usrteia-0005\Documents\Faculdade\FPPD\T2_FPPD\cmd\client
client.exe
```

**Terminal 3 (opcional):**
```bash
cd C:\Users\usrteia-0005\Documents\Faculdade\FPPD\T2_FPPD\cmd\client
client.exe
```

### 3. Jogar!

Use **WASD** para mover e **ESC** para sair.

---

## 📁 Estrutura de Arquivos Atualizada

```
T2_FPPD/
├── mapa.txt                    # Mapa original (raiz)
├── maze.txt                    # Labirinto original (raiz)
│
├── cmd/
│   ├── server/
│   │   └── server.exe          # Executar: cd cmd\server && server.exe
│   │
│   └── client/
│       ├── client.exe          # Executar: cd cmd\client && client.exe
│       ├── mapa.txt            # ✅ COPIADO
│       └── maze.txt            # ✅ COPIADO
│
└── run_demo.bat                # Demo automática (executa da raiz)
```

---

## 🐛 Erros Comuns

### Erro: "open mapa.txt: O sistema não pode encontrar o arquivo"

**Causa:** Executando do diretório errado.

**Solução:**
```bash
# Vá para o diretório do cliente
cd cmd\client

# Execute
client.exe
```

### Erro: "panic: open 8080: O sistema não pode encontrar o arquivo"

**Causa:** Esse erro aparece se o servidor não estiver rodando.

**Solução:**
1. Abra outro terminal
2. Execute o servidor:
   ```bash
   cd cmd\server
   server.exe
   ```
3. Aguarde a mensagem "Aguardando conexões..."
4. Então execute o cliente

---

## ✅ Teste Rápido

Execute este teste para validar:

```bash
# Terminal 1 - Servidor
cd cmd\server
server.exe

# Terminal 2 - Cliente 1
cd cmd\client
client.exe

# Terminal 3 - Cliente 2
cd cmd\client
client.exe
```

**Resultado esperado:**
- ✅ Servidor mostra conexões
- ✅ Clientes mostram interface gráfica
- ✅ Personagem '☺' visível
- ✅ Outros jogadores como '◉'

---

## 🎯 Scripts Atualizados

### run_demo.bat (Atualizado!)

Agora executa corretamente:
```bash
run_demo.bat
```

Abre automaticamente:
1. Servidor em `cmd/server`
2. Cliente 1 em `cmd/client`
3. Cliente 2 em `cmd/client`

---

## 📞 Suporte

### Ainda não funciona?

1. **Verifique o diretório atual:**
   ```bash
   cd
   ```
   Deve mostrar: `C:\Users\usrteia-0005\Documents\Faculdade\FPPD\T2_FPPD\cmd\client`

2. **Verifique se os mapas existem:**
   ```bash
   dir mapa.txt
   ```
   Deve mostrar o arquivo.

3. **Execute o demo automático:**
   ```bash
   cd ..\..
   run_demo.bat
   ```

---

**Data da correção:** 22/10/2025  
**Problema:** Arquivo mapa.txt não encontrado  
**Solução:** Arquivos copiados para cmd/client/  
**Status:** ✅ Corrigido

