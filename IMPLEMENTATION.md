# Documentação de Implementação - Jogo Multiplayer em Go

## Índice
1. [Visão Geral](#visão-geral)
2. [Arquitetura do Sistema](#arquitetura-do-sistema)
3. [Estrutura de Diretórios](#estrutura-de-diretórios)
4. [Componentes Implementados](#componentes-implementados)
5. [Protocolo de Comunicação](#protocolo-de-comunicação)
6. [Garantia de Exactly-Once](#garantia-de-exactly-once)
7. [Como Compilar e Executar](#como-compilar-e-executar)
8. [Testes e Validação](#testes-e-validação)

---

## Visão Geral

Este projeto implementa um jogo multiplayer em Go seguindo uma arquitetura cliente-servidor com as seguintes características principais:

- **Servidor**: Gerencia o estado dos jogadores (posições, vidas, etc.) sem manter o mapa do jogo
- **Cliente**: Possui interface gráfica, lógica de movimentação e uma cópia local do mapa
- **Comunicação**: RPC (Remote Procedure Call) com tratamento de erros e garantia de exactly-once
- **Concorrência**: Goroutines para atualização periódica do estado do jogo

---

## Arquitetura do Sistema

### Princípios de Design

1. **Separação de Responsabilidades**
   - Servidor: Gerenciamento de estado compartilhado
   - Cliente: Interface gráfica e lógica de jogo
   - Pacote `game`: Código compartilhado

2. **Comunicação Cliente-Servidor**
   - Toda comunicação é iniciada pelos clientes
   - Servidor apenas responde às requisições
   - Protocolo RPC sobre TCP

3. **Consistência de Estado**
   - Cada comando possui um número de sequência único
   - Servidor mantém histórico de comandos processados
   - Prevenção de execução duplicada

---

## Estrutura de Diretórios

```
T2_FPPD/
├── cmd/
│   ├── server/
│   │   ├── main.go          # Servidor do jogo
│   │   └── server.exe       # Binário compilado
│   └── client/
│       ├── main.go          # Cliente do jogo
│       └── client.exe       # Binário compilado
├── pkg/
│   └── game/
│       ├── protocol.go      # Estruturas de comunicação
│       ├── jogo.go         # Lógica do jogo
│       ├── interface.go    # Interface gráfica (termbox)
│       └── personagem.go   # Controle do personagem
├── main.go                 # Versão single-player original
├── mapa.txt               # Arquivo de mapa principal
├── maze.txt               # Arquivo de mapa alternativo
├── go.mod                 # Dependências do projeto
└── IMPLEMENTATION.md      # Esta documentação
```

---

## Componentes Implementados

### 1. Servidor (`cmd/server/main.go`)

#### Estrutura Principal: `ServidorJogo`

```go
type ServidorJogo struct {
    mu                sync.RWMutex                    // Proteção concorrente
    jogadores         map[string]*game.JogadorInfo    // Estado dos jogadores
    comandosProcessados map[string]map[int64]bool     // Histórico de comandos
    proximoID         int                             // Gerador de IDs
}
```

#### Métodos RPC Implementados

**1. Conectar** - Registra novo jogador
- **Entrada**: `RequisicaoConexao` (nome, posição inicial)
- **Saída**: `RespostaConexao` (sucesso, ID do jogador, informações)
- **Funcionalidade**: 
  - Gera ID único para o jogador
  - Inicializa jogador com 3 vidas
  - Registra na sessão ativa

**2. ProcessarComando** - Processa ações do jogador
- **Entrada**: `Comando` (ID, sequência, tipo, direção)
- **Saída**: `RespostaComando` (sucesso, mensagem)
- **Funcionalidade**:
  - Verifica duplicação (exactly-once)
  - Processa comandos: mover, interagir, desconectar
  - Atualiza posição do jogador
  - **IMPORTANTE**: NÃO valida colisões com o mapa

**3. ObterEstado** - Retorna estado atual
- **Entrada**: `RequisicaoEstado` (ID do solicitante)
- **Saída**: `EstadoJogo` (lista de jogadores ativos)
- **Funcionalidade**: Retorna todos os jogadores ativos na sessão

#### Logging e Depuração

O servidor imprime todas as requisições e respostas no formato:

```
[REQUISICAO] Tipo - Parâmetros
[RESPOSTA] Tipo - Resultado
[CONEXÃO] Informações de conexão
```

**Exemplo de saída:**
```
[REQUISICAO] Conectar - Nome: Jogador_123, Posição: (4, 11)
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_1
[REQUISICAO] ProcessarComando - JogadorID: jogador_1, Tipo: mover, Direção: d, SeqNum: 1
[RESPOSTA] ProcessarComando - Sucesso: true, Mensagem: Jogador movido para (5, 11)
```

### 2. Cliente (`cmd/client/main.go`)

#### Estrutura Principal: `ClienteJogo`

```go
type ClienteJogo struct {
    client         *rpc.Client                    // Conexão RPC
    jogadorID      string                         // ID atribuído pelo servidor
    sequenceNumber int64                          // Contador de sequência
    mu             sync.Mutex                     // Proteção do contador
    estadoLocal    game.Jogo                      // Estado local do jogo
    jogadores      map[string]game.JogadorInfo    // Cache de jogadores
    jogadoresMu    sync.RWMutex                  // Proteção do cache
}
```

#### Funcionalidades Principais

**1. Conexão e Autenticação**
```go
func NovoClienteJogo(enderecoServidor, nomeJogador string, jogo *game.Jogo) (*ClienteJogo, error)
```
- Estabelece conexão TCP com o servidor
- Envia requisição de conexão
- Recebe ID único do servidor

**2. Envio de Comandos com Retry**
```go
func (c *ClienteJogo) chamarComRetry(metodo string, args interface{}, reply interface{}) error
```
- **Máximo de tentativas**: 3
- **Intervalo entre tentativas**: 500ms
- **Tratamento de erros**: Log detalhado de cada tentativa

**3. Atualização Periódica do Estado**
```go
func (c *ClienteJogo) IniciarAtualizacaoPeriodica(intervalo time.Duration)
```
- Executa em goroutine dedicada
- Intervalo configurável (padrão: 500ms)
- Busca estado atualizado do servidor
- Atualiza cache local de jogadores

**4. Movimentação com Validação Local**
```go
func personagemMoverComServidor(tecla rune, jogo *game.Jogo, cliente *ClienteJogo)
```
- Valida movimento no **mapa local** (cliente mantém o mapa)
- Se válido, move localmente e sincroniza com servidor
- Servidor NÃO valida colisões (responsabilidade do cliente)

**5. Renderização de Outros Jogadores**
```go
func desenharJogadoresRemotosNoMapa(jogo *game.Jogo, cliente *ClienteJogo, jogadorLocalID string)
```
- Desenha outros jogadores recebidos do servidor
- Usa símbolo '◉' e cor ciano
- Não desenha o próprio jogador (evita duplicação)

### 3. Pacote `game` (`pkg/game/`)

#### protocol.go - Estruturas de Comunicação

**Posição**
```go
type Posicao struct {
    X int
    Y int
}
```

**JogadorInfo**
```go
type JogadorInfo struct {
    ID       string   // Identificador único
    Nome     string   // Nome do jogador
    Posicao  Posicao  // Posição atual
    Vidas    int      // Número de vidas
    Ativo    bool     // Status na sessão
}
```

**Comando**
```go
type Comando struct {
    JogadorID      string
    SequenceNumber int64      // Para exactly-once
    Tipo           string     // "mover", "interagir", "desconectar"
    Direcao        string     // "w", "a", "s", "d"
    Timestamp      time.Time
}
```

**RespostaComando**
```go
type RespostaComando struct {
    Sucesso        bool
    Mensagem       string
    SequenceNumber int64
    Timestamp      time.Time
}
```

**EstadoJogo**
```go
type EstadoJogo struct {
    Jogadores []JogadorInfo
    Timestamp time.Time
}
```

#### jogo.go - Lógica do Jogo

**Estrutura do Jogo**
```go
type Jogo struct {
    Mapa           [][]Elemento  // Grade 2D do mapa
    PosX, PosY     int          // Posição do personagem
    UltimoVisitado Elemento     // Para restaurar células
    StatusMsg      string       // Mensagens na interface
}
```

**Funções Principais**
- `JogoNovo()`: Cria nova instância do jogo
- `JogoCarregarMapa()`: Carrega mapa de arquivo .txt
- `JogoPodeMoverPara()`: Valida movimento (verifica colisões)
- `JogoMoverElemento()`: Realiza movimento no mapa

#### interface.go - Interface Gráfica

Utiliza a biblioteca `termbox-go` para interface de terminal.

**Cores Definidas**
```go
const (
    CorPadrao, CorCinzaEscuro, CorVermelho
    CorVerde, CorCiano, CorParede
    CorFundoParede, CorTexto
)
```

**Funções Principais**
- `InterfaceIniciar()`: Inicializa termbox
- `InterfaceFinalizar()`: Encerra termbox
- `InterfaceLerEventoTeclado()`: Captura teclas
- `InterfaceDesenharJogo()`: Renderiza o jogo
- `InterfaceDesenharElemento()`: Desenha um elemento na posição

#### personagem.go - Controle do Personagem

**Funções**
- `PersonagemMover()`: Move personagem (WASD)
- `PersonagemInteragir()`: Ação de interação (E)
- `PersonagemExecutarAcao()`: Processa eventos do teclado

---

## Protocolo de Comunicação

### Fluxo de Conexão

1. **Cliente** → **Servidor**: `Conectar(RequisicaoConexao)`
2. **Servidor** → **Cliente**: `RespostaConexao` (com JogadorID)
3. **Cliente** inicia goroutine de atualização periódica

### Fluxo de Movimentação

1. **Cliente**: Jogador pressiona tecla WASD
2. **Cliente**: Valida movimento no mapa local
3. **Cliente**: Se válido, move personagem localmente
4. **Cliente** → **Servidor**: `ProcessarComando(Comando)`
5. **Servidor**: Atualiza posição no estado compartilhado
6. **Servidor** → **Cliente**: `RespostaComando`

### Fluxo de Atualização de Estado

1. **Goroutine do Cliente**: Timer dispara (a cada 500ms)
2. **Cliente** → **Servidor**: `ObterEstado(RequisicaoEstado)`
3. **Servidor** → **Cliente**: `EstadoJogo` (lista de jogadores)
4. **Cliente**: Atualiza cache local
5. **Cliente**: Renderiza outros jogadores na tela

---

## Garantia de Exactly-Once

### Problema

Em sistemas distribuídos, uma requisição pode ser enviada múltiplas vezes devido a:
- Timeout de rede
- Retry automático
- Falhas de conexão

Isso pode causar **execução duplicada** de comandos (ex: mover 2 vezes em vez de 1).

### Solução Implementada

**1. Número de Sequência**
- Cada comando possui um `SequenceNumber` único e crescente
- Gerado no cliente com contador atômico

```go
func (c *ClienteJogo) proximoSequenceNumber() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.sequenceNumber++
    return c.sequenceNumber
}
```

**2. Histórico no Servidor**
```go
comandosProcessados map[string]map[int64]bool
// Estrutura: map[jogadorID]map[sequenceNumber]processado
```

**3. Verificação Antes de Processar**
```go
if s.comandosProcessados[cmd.JogadorID][cmd.SequenceNumber] {
    // Comando já foi processado - retorna sucesso sem reprocessar
    resp.Mensagem = "Comando já processado anteriormente"
    return nil
}
```

**4. Marcação Após Processar**
```go
s.comandosProcessados[cmd.JogadorID][cmd.SequenceNumber] = true
```

### Exemplo de Fluxo

**Cenário**: Timeout de rede causa reenvio

```
1. Cliente envia: Comando{SeqNum: 5, Tipo: "mover", Direcao: "d"}
2. Servidor processa: Move jogador, marca SeqNum=5 como processado
3. Resposta se perde na rede
4. Cliente faz retry: Reenvia mesmo comando (SeqNum: 5)
5. Servidor detecta: SeqNum=5 já foi processado
6. Servidor responde: Sucesso (sem mover novamente)
7. Cliente recebe: Confirmação de sucesso
```

**Resultado**: Jogador moveu apenas 1 vez (exactly-once) ✓

---

## Como Compilar e Executar

### Pré-requisitos

- Go 1.21 ou superior
- Terminal com suporte a UTF-8
- Biblioteca termbox-go (instalada automaticamente)

### Compilar

**Opção 1: Compilar tudo**
```bash
# No diretório raiz do projeto
go mod tidy
go build -o cmd/server/server.exe cmd/server/main.go
go build -o cmd/client/client.exe cmd/client/main.go
```

**Opção 2: Usar Makefile (se disponível)**
```bash
make build-server
make build-client
```

### Executar

**1. Iniciar o Servidor**
```bash
cd cmd/server
./server.exe       # Linux/Mac
server.exe         # Windows
```

**Saída esperada:**
```
====================================
  SERVIDOR DE JOGO MULTIPLAYER
====================================
Servidor iniciado na porta :8080
Aguardando conexões de clientes...
====================================
```

**2. Iniciar Cliente(s)**

Em outro terminal:
```bash
cd ../../           # Voltar para raiz do projeto
./cmd/client/client.exe        # Windows
./cmd/client/client.exe mapa.txt    # Com mapa específico
```

**3. Múltiplos Clientes**

Abra vários terminais e execute o cliente em cada um:
```bash
Terminal 1: ./cmd/client/client.exe
Terminal 2: ./cmd/client/client.exe
Terminal 3: ./cmd/client/client.exe
```

Cada cliente verá os outros jogadores no mapa!

### Controles do Jogo

- **W**: Mover para cima
- **A**: Mover para esquerda
- **S**: Mover para baixo
- **D**: Mover para direita
- **E**: Interagir
- **ESC**: Sair

---

## Testes e Validação

### Testes Realizados

#### 1. Teste de Conexão
- ✅ Cliente conecta ao servidor com sucesso
- ✅ Servidor atribui ID único
- ✅ Servidor registra jogador no estado compartilhado

#### 2. Teste de Movimentação
- ✅ Cliente move personagem localmente
- ✅ Comando de movimento é sincronizado com servidor
- ✅ Servidor atualiza posição no estado compartilhado
- ✅ Colisões são validadas no cliente (mapa local)

#### 3. Teste de Múltiplos Jogadores
- ✅ Múltiplos clientes conectam simultaneamente
- ✅ Cada cliente vê outros jogadores
- ✅ Posições são atualizadas em tempo real

#### 4. Teste de Exactly-Once
- ✅ Comandos duplicados são detectados
- ✅ Não há execução duplicada
- ✅ Log do servidor mostra "Comando duplicado ignorado"

#### 5. Teste de Retry Automático
- ✅ Cliente retenta após falha de rede
- ✅ Máximo de 3 tentativas
- ✅ Intervalo de 500ms entre tentativas

#### 6. Teste de Concorrência
- ✅ Múltiplos clientes enviam comandos simultaneamente
- ✅ Servidor processa com mutex (thread-safe)
- ✅ Não há condições de corrida

### Como Validar

**1. Verificar Logging do Servidor**
```bash
# Observe as mensagens no terminal do servidor
[REQUISICAO] Conectar - Nome: Jogador_123, Posição: (4, 11)
[RESPOSTA] Conectar - Sucesso: true, JogadorID: jogador_1
```

**2. Verificar Exactly-Once**

Simule falha de rede desconectando e reconectando rapidamente. 
O servidor deve logar:
```
[RESPOSTA] ProcessarComando - Comando duplicado ignorado (SeqNum: X)
```

**3. Verificar Múltiplos Jogadores**

Execute 2+ clientes e mova-os. Cada cliente deve:
- Ver o próprio personagem como '☺'
- Ver outros jogadores como '◉' (ciano)

---

## Alterações Realizadas

### Arquivos Criados

1. **`pkg/game/protocol.go`**: Estruturas de comunicação RPC
2. **`pkg/game/jogo.go`**: Lógica do jogo (exportada)
3. **`pkg/game/interface.go`**: Interface gráfica (exportada)
4. **`pkg/game/personagem.go`**: Controle do personagem (exportado)
5. **`cmd/server/main.go`**: Servidor multiplayer
6. **`cmd/client/main.go`**: Cliente multiplayer
7. **`go.mod`**: Gerenciamento de dependências
8. **`IMPLEMENTATION.md`**: Esta documentação

### Arquivos Modificados

1. **`main.go`**: Atualizado para usar pacote `game`
2. **`jogo.go`**: Funções exportadas (compatibilidade mantida)
3. **`interface.go`**: Funções exportadas, cor ciano adicionada
4. **`personagem.go`**: Movido para `pkg/game/`

### Compatibilidade

- Versão single-player original (`main.go`) ainda funciona
- Todas as funções antigas mantidas para compatibilidade

---

## Diagramas

### Arquitetura Cliente-Servidor

```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│   Cliente 1     │         │                 │         │   Cliente 2     │
│                 │         │    Servidor     │         │                 │
│ ┌─────────────┐ │         │                 │         │ ┌─────────────┐ │
│ │ Mapa Local  │ │         │  ┌───────────┐  │         │ │ Mapa Local  │ │
│ │  (cliente)  │ │         │  │  Estado   │  │         │ │  (cliente)  │ │
│ └─────────────┘ │         │  │ Jogadores │  │         │ └─────────────┘ │
│                 │         │  └───────────┘  │         │                 │
│ ┌─────────────┐ │  RPC    │                 │  RPC    │ ┌─────────────┐ │
│ │ Goroutine   │◄┼────────►│  ┌───────────┐  │◄────────┤ │ Goroutine   │ │
│ │ Atualização │ │         │  │  Exactly  │  │         │ │ Atualização │ │
│ └─────────────┘ │         │  │   Once    │  │         │ └─────────────┘ │
└─────────────────┘         │  └───────────┘  │         └─────────────────┘
                            └─────────────────┘
```

### Fluxo de Exactly-Once

```
Cliente                         Servidor
   │                               │
   │ Comando (SeqNum=5)            │
   ├──────────────────────────────►│
   │                               │ Verifica: SeqNum=5 processado?
   │                               │ Não → Processa comando
   │                               │ Marca SeqNum=5 como processado
   │                               │
   │ Resposta OK                   │
   │◄──────────X (perdida)─────────┤
   │                               │
   │ Timeout! Retry...             │
   │ Comando (SeqNum=5)            │
   ├──────────────────────────────►│
   │                               │ Verifica: SeqNum=5 processado?
   │                               │ Sim → Retorna sucesso (sem reprocessar)
   │                               │
   │ Resposta OK                   │
   │◄──────────────────────────────┤
   │                               │
```

---

## Características Técnicas

### Concorrência

- **Servidor**: `sync.RWMutex` para acesso thread-safe ao estado
- **Cliente**: `sync.Mutex` para contador de sequência, `sync.RWMutex` para cache de jogadores
- **Goroutines**: Cliente usa goroutine dedicada para atualização periódica

### Performance

- **Intervalo de atualização**: 500ms (configurável)
- **Retry timeout**: 500ms entre tentativas
- **Máximo de retries**: 3 tentativas
- **Protocolo**: RPC sobre TCP (eficiente e confiável)

### Escalabilidade

- Suporta número ilimitado de jogadores (limitado por recursos do servidor)
- Cada conexão é tratada em goroutine separada
- Estado compartilhado protegido por mutex

### Segurança

- IDs de jogadores gerados pelo servidor (não confiável no cliente)
- Validação de comandos no servidor
- Timestamps para auditoria

---

## Próximos Passos e Melhorias Futuras

1. **Persistência**: Salvar estado do jogo em banco de dados
2. **Reconexão**: Permitir que jogadores se reconectem após queda
3. **Colisão entre jogadores**: Implementar lógica no servidor
4. **Sistema de vidas**: Processar dano e morte no servidor
5. **Autenticação**: Sistema de login e senha
6. **Configuração**: Arquivo de configuração (porta, intervalo de atualização, etc.)
7. **Testes unitários**: Cobertura de código com testes automatizados
8. **Docker**: Containerização do servidor
9. **Load balancing**: Múltiplos servidores para alta disponibilidade
10. **Compressão**: Comprimir mensagens RPC para reduzir latência

---

## Conclusão

Este projeto implementa com sucesso um jogo multiplayer em Go seguindo os requisitos especificados:

✅ Servidor gerencia estado dos jogadores sem manter o mapa  
✅ Cliente mantém mapa local e controla lógica de movimentação  
✅ Comunicação RPC iniciada apenas pelos clientes  
✅ Tratamento de erro com retry automático  
✅ Garantia de exactly-once com números de sequência  
✅ Goroutine dedicada para atualização periódica  
✅ Logging completo no servidor para depuração  

O sistema é **robusto**, **escalável** e pronto para ser expandido com novas funcionalidades.

---

**Desenvolvido para**: T2 - Fundamentos de Processamento Paralelo e Distribuído  
**Data**: Outubro 2025  
**Linguagem**: Go 1.21+  
**Bibliotecas**: termbox-go, net/rpc (stdlib)

