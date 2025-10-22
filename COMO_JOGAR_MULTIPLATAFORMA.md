# 🎮 Como Jogar Multiplataforma - Guia Rápido

## ⚡ Início Rápido (3 Passos)

### 1️⃣ Recompilar o Cliente (necessário!)

```bash
# Windows
go build -o cmd\client\client.exe cmd\client\main.go

# Linux
go build -o cmd/client/client cmd/client/main.go
```

**Por quê?** O cliente foi atualizado para aceitar IP do servidor!

---

### 2️⃣ Iniciar Servidor

**Windows:**
```cmd
cd cmd\server
server.exe
```

**Linux:**
```bash
cd cmd/server
./server
```

**Anote o IP:**
- Windows: `ipconfig` → procure "IPv4"
- Linux: `ip addr` → procure "inet"

Exemplo: `192.168.1.100`

---

### 3️⃣ Conectar Clientes

**Mesmo PC (localhost):**
```bash
# Windows
cd cmd\client
client.exe

# Linux
cd cmd/client
./client
```

**Outro PC na rede:**
```bash
# Windows
cd cmd\client
client.exe mapa.txt 192.168.1.100:8080

# Linux
cd cmd/client
./client mapa.txt 192.168.1.100:8080
```

Substitua `192.168.1.100` pelo IP do servidor!

---

## 🎯 Exemplos Práticos

### Exemplo 1: Servidor Windows + Cliente Linux

**PC 1 (Windows - Servidor):**
```cmd
ipconfig
REM Digamos que o IP é 192.168.1.100

cd cmd\server
server.exe
```

**PC 2 (Linux - Cliente):**
```bash
cd cmd/client
./client mapa.txt 192.168.1.100:8080
```

✅ Funciona!

---

### Exemplo 2: Servidor Linux + 2 Clientes Windows

**PC 1 (Linux - Servidor):**
```bash
ip addr
# Digamos que o IP é 192.168.0.50

cd cmd/server
./server
```

**PC 2 (Windows - Cliente 1):**
```cmd
cd cmd\client
client.exe mapa.txt 192.168.0.50:8080
```

**PC 3 (Windows - Cliente 2):**
```cmd
cd cmd\client
client.exe mapa.txt 192.168.0.50:8080
```

✅ Os 2 clientes Windows veem um ao outro no servidor Linux!

---

## 🛠️ Scripts Prontos

### Windows: Conectar Facilmente

```cmd
conectar_servidor.bat
```

Digite o IP quando perguntado!

### Linux: Conectar Facilmente

```bash
chmod +x conectar_servidor.sh
./conectar_servidor.sh
```

Digite o IP quando perguntado!

---

## 🔥 Firewall (IMPORTANTE!)

### Windows

```powershell
# Execute como Administrador
netsh advfirewall firewall add rule name="Jogo Go" dir=in action=allow protocol=TCP localport=8080
```

### Linux

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow 8080/tcp

# CentOS/RHEL (firewalld)
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

---

## 🧪 Testar Conexão

**Do cliente para o servidor:**

```bash
# Windows
ping 192.168.1.100
telnet 192.168.1.100 8080

# Linux
ping 192.168.1.100
nc -zv 192.168.1.100 8080
```

Se funcionar → jogo vai funcionar! ✓

---

## 📝 Sintaxe Completa

```bash
client [mapa] [endereco:porta]
```

**Exemplos:**

```bash
# Localhost (padrão)
client.exe

# Mapa diferente no localhost
client.exe maze.txt

# Servidor remoto
client.exe mapa.txt 192.168.1.100:8080

# Servidor com nome de host
client.exe mapa.txt servidor.local:8080
```

---

## 🐛 Problemas?

### "connection refused"

✅ **Solução:**
1. Servidor está rodando?
2. IP está correto?
3. Firewall está aberto?

### "no such host"

✅ **Solução:**
- Verifique o IP
- Use IP numérico ao invés de nome

### "timeout"

✅ **Solução:**
1. Ping funciona?
2. Mesma rede?
3. Firewall configurado?

---

## 📚 Documentação Completa

Para mais detalhes, veja:
- **`MULTIPLATAFORMA.md`** - Guia completo
- **`IMPLEMENTATION.md`** - Arquitetura técnica

---

## ✅ Checklist Rápido

Antes de jogar multiplataforma:

- [ ] Cliente recompilado com suporte a IP remoto
- [ ] Servidor rodando
- [ ] IP do servidor anotado
- [ ] Firewall configurado
- [ ] Teste de ping funcionando
- [ ] Arquivos de mapa copiados para cmd/client/

---

## 🎉 Resumo

**Cliente atualizado** para aceitar IP do servidor!

**Uso:**
```bash
client.exe mapa.txt IP_DO_SERVIDOR:8080
```

**Funciona em qualquer combinação:**
- Windows ↔ Windows
- Windows ↔ Linux
- Linux ↔ Linux
- Windows ↔ macOS
- Linux ↔ macOS

**Divirta-se! 🎮**

---

**Atualizado:** 22/10/2025  
**Status:** ✅ Multiplataforma funcionando

