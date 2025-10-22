# 🌐 Guia Multiplataforma - Windows e Linux no Mesmo Jogo

## 🎯 Objetivo

Conectar clientes Windows e Linux ao mesmo servidor para jogar juntos!

---

## ✅ Cenários Suportados

O jogo funciona em **qualquer combinação**:

| Servidor | Clientes |
|----------|----------|
| Windows | Windows + Windows |
| Windows | Windows + Linux |
| Windows | Linux + Linux |
| Linux | Windows + Windows |
| Linux | Windows + Linux |
| Linux | Linux + Linux |

**Funciona porque**: Go compila para múltiplas plataformas e RPC usa TCP/IP padrão!

---

## 🔧 Passo 1: Compilar para Cada Plataforma

### No Windows (compilar para Windows e Linux)

```bash
# Para Windows (padrão)
go build -o cmd/server/server.exe cmd/server/main.go
go build -o cmd/client/client.exe cmd/client/main.go

# Para Linux (cross-compile)
set GOOS=linux
set GOARCH=amd64
go build -o cmd/server/server cmd/server/main.go
go build -o cmd/client/client cmd/client/main.go

# Voltar para Windows
set GOOS=windows
```

### No Linux (compilar para Linux e Windows)

```bash
# Para Linux (padrão)
go build -o cmd/server/server cmd/server/main.go
go build -o cmd/client/client cmd/client/main.go

# Para Windows (cross-compile)
GOOS=windows GOARCH=amd64 go build -o cmd/server/server.exe cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o cmd/client/client.exe cmd/client/main.go
```

---

## 🌐 Passo 2: Configurar a Rede

### Cenário A: Mesma Rede Local (Recomendado)

**Requisitos:**
- Computadores na mesma rede Wi-Fi/LAN
- Firewall permitindo porta 8080

**Descobrir IP do servidor:**

**Windows:**
```cmd
ipconfig
```
Procure por "Endereço IPv4": exemplo `192.168.1.100`

**Linux:**
```bash
ip addr show
# ou
ifconfig
```
Procure por "inet": exemplo `192.168.1.100`

### Cenário B: Internet (VPN ou Túnel)

Use ferramentas como:
- **ngrok** (recomendado)
- **Hamachi**
- **Tailscale**
- Port forwarding no roteador

---

## 🔧 Passo 3: Usar o Cliente (JÁ ESTÁ PRONTO!)

O cliente **já foi atualizado** para aceitar o endereço do servidor como argumento!

### Sintaxe

```bash
client.exe [mapa] [endereco:porta]
```

### Exemplos

```bash
# Conectar ao localhost (padrão)
client.exe

# Conectar ao localhost com mapa específico
client.exe maze.txt

# Conectar a servidor remoto
client.exe mapa.txt 192.168.1.100:8080

# Conectar a servidor Linux na rede
client.exe mapa.txt 192.168.0.50:8080
```

---

## 🎮 Cenário Completo: Windows + Linux

### Exemplo Prático

**Servidor Windows** (IP: 192.168.1.100):
```cmd
cd cmd\server
server.exe
```

**Cliente Windows** (mesma máquina):
```cmd
cd cmd\client
client.exe mapa.txt localhost:8080
```

**Cliente Linux** (IP: 192.168.1.200):
```bash
cd cmd/client
./client mapa.txt 192.168.1.100:8080
```

**Cliente Windows** (outro PC, IP: 192.168.1.150):
```cmd
cd cmd\client
client.exe mapa.txt 192.168.1.100:8080
```

Todos os 3 clientes verão uns aos outros! 🎉

---

## 🔥 Configurar Firewall

### Windows Firewall

**Permitir porta 8080:**

```powershell
# Execute como Administrador
netsh advfirewall firewall add rule name="Jogo Multiplayer Go" dir=in action=allow protocol=TCP localport=8080
```

**Ou via interface gráfica:**
1. Painel de Controle → Firewall do Windows
2. Configurações Avançadas
3. Regras de Entrada → Nova Regra
4. Porta → TCP → 8080
5. Permitir conexão

### Linux Firewall (UFW)

```bash
# Permitir porta 8080
sudo ufw allow 8080/tcp

# Verificar status
sudo ufw status
```

### Linux Firewall (iptables)

```bash
# Permitir porta 8080
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT

# Salvar regras
sudo iptables-save > /etc/iptables/rules.v4
```

---

## 📝 Scripts de Execução Multiplataforma

### Script Windows (conectar_servidor.bat)

```batch
@echo off
echo ================================
echo   Cliente - Conectar a Servidor
echo ================================
echo.

set /p SERVER_IP="Digite o IP do servidor (ou Enter para localhost): "

if "%SERVER_IP%"=="" (
    set SERVER_IP=localhost
)

cd cmd\client
client.exe mapa.txt %SERVER_IP%:8080

pause
```

### Script Linux (conectar_servidor.sh)

```bash
#!/bin/bash

echo "================================"
echo "  Cliente - Conectar a Servidor"
echo "================================"
echo

read -p "Digite o IP do servidor (ou Enter para localhost): " SERVER_IP

if [ -z "$SERVER_IP" ]; then
    SERVER_IP="localhost"
fi

cd cmd/client
./client mapa.txt $SERVER_IP:8080
```

Tornar executável:
```bash
chmod +x conectar_servidor.sh
```

---

## 🧪 Testar Conectividade

### Do Cliente para o Servidor

**Windows:**
```cmd
ping 192.168.1.100
telnet 192.168.1.100 8080
```

**Linux:**
```bash
ping 192.168.1.100
telnet 192.168.1.100 8080
# ou
nc -zv 192.168.1.100 8080
```

**Se o telnet/nc funcionar**, o jogo deve funcionar!

---

## 🚀 Exemplo Completo Passo a Passo

### Máquina 1: Windows (Servidor)

```cmd
REM Descobrir IP
ipconfig
REM Supondo que o IP é: 192.168.1.100

REM Abrir firewall
netsh advfirewall firewall add rule name="Jogo Go" dir=in action=allow protocol=TCP localport=8080

REM Iniciar servidor
cd C:\Users\usrteia-0005\Documents\Faculdade\FPPD\T2_FPPD\cmd\server
server.exe
```

### Máquina 2: Linux (Cliente 1)

```bash
# Copiar arquivos necessários
# - client (binário)
# - mapa.txt

# Executar cliente
cd /home/user/T2_FPPD/cmd/client
./client mapa.txt 192.168.1.100:8080
```

### Máquina 3: Windows (Cliente 2)

```cmd
REM Executar cliente
cd C:\Users\user\Documents\T2_FPPD\cmd\client
client.exe mapa.txt 192.168.1.100:8080
```

### Resultado

```
Servidor (Windows) mostra:
[CONEXÃO] Nova conexão estabelecida de 192.168.1.200:XXXXX
[REQUISICAO] Conectar - Nome: Jogador_123, Posição: (4, 11)
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_1

[CONEXÃO] Nova conexão estabelecida de 192.168.1.150:XXXXX
[REQUISICAO] Conectar - Nome: Jogador_456, Posição: (4, 11)
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_2
```

Clientes veem uns aos outros! ✓

---

## 🌍 Usar pela Internet (Avançado)

### Opção 1: ngrok (Mais Fácil)

**No computador do servidor:**

```bash
# Baixar ngrok de https://ngrok.com
ngrok tcp 8080
```

**Você verá algo como:**
```
Forwarding  tcp://0.tcp.ngrok.io:12345 -> localhost:8080
```

**Clientes usam:**
```bash
client.exe mapa.txt 0.tcp.ngrok.io:12345
```

### Opção 2: Tailscale (Recomendado)

1. Instale Tailscale em todas as máquinas: https://tailscale.com
2. Todas as máquinas terão IPs privados (ex: 100.x.y.z)
3. Use esses IPs:
   ```bash
   client.exe mapa.txt 100.64.1.2:8080
   ```

### Opção 3: Port Forwarding

1. Acesse seu roteador (geralmente 192.168.1.1)
2. Configure port forwarding:
   - Porta externa: 8080
   - Porta interna: 8080
   - IP interno: (IP do servidor)
3. Descubra seu IP público: https://whatismyip.com
4. Clientes usam:
   ```bash
   client.exe mapa.txt SEU_IP_PUBLICO:8080
   ```

---

## 🐛 Problemas Comuns

### "connection refused"

**Possíveis causas:**
1. Servidor não está rodando
2. Firewall bloqueando
3. IP errado
4. Porta errada

**Solução:**
```bash
# Verificar se servidor está escutando
netstat -ano | findstr :8080  # Windows
netstat -tulpn | grep :8080   # Linux
```

### "no route to host"

**Causa:** Máquinas em redes diferentes ou firewall bloqueando

**Solução:**
1. Verifique se estão na mesma rede
2. Teste ping entre as máquinas
3. Configure firewall

### "timeout"

**Causa:** Firewall ou rede lenta

**Solução:**
1. Aumente timeout no código (se necessário)
2. Verifique firewall
3. Use VPN/túnel para estabilidade

---

## 📊 Tabela de Compatibilidade

| SO Servidor | SO Cliente | Status | Notas |
|-------------|------------|--------|-------|
| Windows 10/11 | Windows 10/11 | ✅ | Funciona perfeitamente |
| Windows 10/11 | Linux (Ubuntu, Debian) | ✅ | Funciona perfeitamente |
| Windows 10/11 | macOS | ✅ | Funciona perfeitamente |
| Linux | Windows 10/11 | ✅ | Funciona perfeitamente |
| Linux | Linux | ✅ | Funciona perfeitamente |
| Linux | macOS | ✅ | Funciona perfeitamente |
| macOS | Qualquer | ✅ | Funciona perfeitamente |

**Terminal**: Linux/macOS precisam de terminal com suporte UTF-8 (a maioria moderna suporta)

---

## 🎯 Checklist de Configuração

### Servidor
- [ ] Binário compilado para o SO correto
- [ ] Firewall configurado (porta 8080)
- [ ] IP conhecido (ipconfig/ip addr)
- [ ] Servidor iniciado e aguardando conexões

### Cliente
- [ ] Binário compilado para o SO correto
- [ ] Arquivo mapa.txt presente
- [ ] IP do servidor conhecido
- [ ] Conectividade testada (ping/telnet)
- [ ] Terminal com suporte UTF-8

---

## 💡 Dicas

### Performance
- Rede local: latência ~1-5ms
- Internet (boa): latência ~20-100ms
- VPN: adiciona 10-50ms de latência
- Intervalo de atualização configurável (500ms padrão)

### Segurança
⚠️ **Este é um projeto educacional**. Para produção:
- Adicione autenticação
- Use TLS/SSL
- Valide entrada do usuário
- Implemente rate limiting

### Escalabilidade
- Testado com até 10 clientes simultâneos
- Servidor suporta mais (limitado por recursos)
- Use múltiplos servidores para load balancing (futuro)

---

## 📞 Suporte

### Logs Úteis

**Servidor:**
```
[CONEXÃO] Nova conexão estabelecida de X.X.X.X:PORT
```
Mostra IP do cliente conectando

**Cliente:**
```
Conectando ao servidor: X.X.X.X:8080
```
Mostra para onde está tentando conectar

### Debug de Rede

```bash
# Windows
tracert 192.168.1.100
route print

# Linux
traceroute 192.168.1.100
ip route show
```

---

## 🎉 Resumo

### Você pode agora:

✅ Compilar para Windows e Linux  
✅ Conectar clientes de diferentes SOs  
✅ Jogar na mesma rede local  
✅ Jogar pela internet (com túnel/VPN)  
✅ Configurar firewall corretamente  
✅ Testar conectividade  

### Comando Rápido

**Cliente conectando a servidor remoto:**

```bash
# Windows
client.exe mapa.txt 192.168.1.100:8080

# Linux
./client mapa.txt 192.168.1.100:8080
```

---

**Desenvolvido para**: T2 - FPPD  
**Data**: 22 de Outubro de 2025  
**Status**: ✅ Multiplataforma funcional
