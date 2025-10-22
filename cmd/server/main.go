// cmd/server/main.go - Servidor do jogo multiplayer
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/usrteia-0005/T2_FPPD/pkg/game"
)

// ServidorJogo gerencia o estado da sessão de jogo
type ServidorJogo struct {
	mu                  sync.RWMutex                 // Mutex para proteger acesso concorrente
	jogadores           map[string]*game.JogadorInfo // Mapa de jogadores ativos
	comandosProcessados map[string]map[int64]bool    // Controle de comandos já processados (exactly-once)
	proximoID           int                          // Próximo ID numérico para jogadores
}

// NovoServidorJogo cria uma nova instância do servidor
func NovoServidorJogo() *ServidorJogo {
	return &ServidorJogo{
		jogadores:           make(map[string]*game.JogadorInfo),
		comandosProcessados: make(map[string]map[int64]bool),
		proximoID:           1,
	}
}

// Conectar registra um novo jogador na sessão
func (s *ServidorJogo) Conectar(req *game.RequisicaoConexao, resp *game.RespostaConexao) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("[REQUISICAO] Conectar - Nome: %s, Posição: (%d, %d)\n",
		req.Nome, req.PosicaoInicial.X, req.PosicaoInicial.Y)

	// Gera um ID único para o jogador
	jogadorID := fmt.Sprintf("jogador_%d", s.proximoID)
	s.proximoID++

	// Cria informações do jogador
	jogador := &game.JogadorInfo{
		ID:      jogadorID,
		Nome:    req.Nome,
		Posicao: req.PosicaoInicial,
		Vidas:   3, // Inicia com 3 vidas
		Ativo:   true,
	}

	// Registra o jogador
	s.jogadores[jogadorID] = jogador
	s.comandosProcessados[jogadorID] = make(map[int64]bool)

	// Prepara resposta
	resp.Sucesso = true
	resp.JogadorID = jogadorID
	resp.Mensagem = fmt.Sprintf("Jogador %s conectado com sucesso!", req.Nome)
	resp.Jogador = *jogador

	fmt.Printf("[RESPOSTA] Conectar - Sucesso: %v, JogadorID: %s\n", resp.Sucesso, resp.JogadorID)
	return nil
}

// ProcessarComando processa um comando enviado por um cliente
func (s *ServidorJogo) ProcessarComando(cmd *game.Comando, resp *game.RespostaComando) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("[REQUISICAO] ProcessarComando - JogadorID: %s, Tipo: %s, Direção: %s, SeqNum: %d\n",
		cmd.JogadorID, cmd.Tipo, cmd.Direcao, cmd.SequenceNumber)

	// Verifica se o comando já foi processado (exactly-once)
	if s.comandosProcessados[cmd.JogadorID][cmd.SequenceNumber] {
		resp.Sucesso = true
		resp.Mensagem = "Comando já processado anteriormente"
		resp.SequenceNumber = cmd.SequenceNumber
		resp.Timestamp = time.Now()

		fmt.Printf("[RESPOSTA] ProcessarComando - Comando duplicado ignorado (SeqNum: %d)\n", cmd.SequenceNumber)
		return nil
	}

	// Verifica se o jogador existe
	jogador, existe := s.jogadores[cmd.JogadorID]
	if !existe {
		resp.Sucesso = false
		resp.Mensagem = "Jogador não encontrado"
		resp.SequenceNumber = cmd.SequenceNumber
		resp.Timestamp = time.Now()

		fmt.Printf("[RESPOSTA] ProcessarComando - Erro: Jogador não encontrado\n")
		return nil
	}

	// Processa o comando de acordo com o tipo
	switch cmd.Tipo {
	case "mover":
		// Atualiza a posição do jogador com base na direção
		novaPosicao := jogador.Posicao
		switch cmd.Direcao {
		case "w":
			novaPosicao.Y--
		case "s":
			novaPosicao.Y++
		case "a":
			novaPosicao.X--
		case "d":
			novaPosicao.X++
		}

		// IMPORTANTE: O servidor NÃO valida se a posição é válida no mapa
		// Isso é responsabilidade do cliente, que possui o mapa
		jogador.Posicao = novaPosicao

		resp.Sucesso = true
		resp.Mensagem = fmt.Sprintf("Jogador movido para (%d, %d)", novaPosicao.X, novaPosicao.Y)

	case "interagir":
		resp.Sucesso = true
		resp.Mensagem = fmt.Sprintf("Interação registrada em (%d, %d)", jogador.Posicao.X, jogador.Posicao.Y)

	case "desconectar":
		jogador.Ativo = false
		resp.Sucesso = true
		resp.Mensagem = "Jogador desconectado"

	default:
		resp.Sucesso = false
		resp.Mensagem = "Tipo de comando desconhecido"
	}

	// Marca o comando como processado
	s.comandosProcessados[cmd.JogadorID][cmd.SequenceNumber] = true

	resp.SequenceNumber = cmd.SequenceNumber
	resp.Timestamp = time.Now()

	fmt.Printf("[RESPOSTA] ProcessarComando - Sucesso: %v, Mensagem: %s\n", resp.Sucesso, resp.Mensagem)
	return nil
}

// ObterEstado retorna o estado atual do jogo
func (s *ServidorJogo) ObterEstado(req *game.RequisicaoEstado, resp *game.EstadoJogo) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Printf("[REQUISICAO] ObterEstado - JogadorID: %s\n", req.JogadorID)

	// Coleta todos os jogadores ativos
	resp.Jogadores = make([]game.JogadorInfo, 0, len(s.jogadores))
	for _, jogador := range s.jogadores {
		if jogador.Ativo {
			resp.Jogadores = append(resp.Jogadores, *jogador)
		}
	}

	resp.Timestamp = time.Now()

	fmt.Printf("[RESPOSTA] ObterEstado - Total de jogadores ativos: %d\n", len(resp.Jogadores))
	return nil
}

func main() {
	// Cria o servidor de jogo
	servidor := NovoServidorJogo()

	// Registra o servidor RPC
	if err := rpc.Register(servidor); err != nil {
		log.Fatalf("Erro ao registrar servidor RPC: %v", err)
	}

	// Configura o listener TCP
	porta := ":8080"
	listener, err := net.Listen("tcp", porta)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor na porta %s: %v", porta, err)
	}
	defer listener.Close()

	fmt.Println("====================================")
	fmt.Println("  SERVIDOR DE JOGO MULTIPLAYER")
	fmt.Println("====================================")
	fmt.Printf("Servidor iniciado na porta %s\n", porta)
	fmt.Println("Aguardando conexões de clientes...")
	fmt.Println("====================================\n")

	// Loop principal: aceita conexões
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}

		fmt.Printf("[CONEXÃO] Nova conexão estabelecida de %s\n", conn.RemoteAddr())

		// Processa cada conexão em uma goroutine separada
		go rpc.ServeConn(conn)
	}
}
