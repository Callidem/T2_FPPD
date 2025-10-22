// cmd/client/main.go - Cliente do jogo multiplayer
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/usrteia-0005/T2_FPPD/pkg/game"
)

// ClienteJogo gerencia a conexão e comunicação com o servidor
type ClienteJogo struct {
	client         *rpc.Client
	jogadorID      string
	sequenceNumber int64
	mu             sync.Mutex
	estadoLocal    game.Jogo
	jogadores      map[string]game.JogadorInfo
	jogadoresMu    sync.RWMutex
}

// NovoClienteJogo cria e conecta um novo cliente ao servidor
func NovoClienteJogo(enderecoServidor, nomeJogador string, jogo *game.Jogo) (*ClienteJogo, error) {
	// Conecta ao servidor RPC
	client, err := rpc.Dial("tcp", enderecoServidor)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao servidor: %v", err)
	}

	c := &ClienteJogo{
		client:         client,
		sequenceNumber: 0,
		estadoLocal:    *jogo,
		jogadores:      make(map[string]game.JogadorInfo),
	}

	// Envia requisição de conexão
	reqConexao := &game.RequisicaoConexao{
		Nome: nomeJogador,
		PosicaoInicial: game.Posicao{
			X: jogo.PosX,
			Y: jogo.PosY,
		},
	}

	var respConexao game.RespostaConexao
	if err := c.chamarComRetry("ServidorJogo.Conectar", reqConexao, &respConexao); err != nil {
		client.Close()
		return nil, fmt.Errorf("erro ao conectar jogador: %v", err)
	}

	if !respConexao.Sucesso {
		client.Close()
		return nil, fmt.Errorf("falha ao conectar: %s", respConexao.Mensagem)
	}

	c.jogadorID = respConexao.JogadorID
	log.Printf("Conectado ao servidor como %s (ID: %s)", nomeJogador, c.jogadorID)
	jogo.StatusMsg = respConexao.Mensagem

	return c, nil
}

// chamarComRetry executa uma chamada RPC com retry automático em caso de falha
func (c *ClienteJogo) chamarComRetry(metodo string, args interface{}, reply interface{}) error {
	maxTentativas := 3
	intervaloRetry := 500 * time.Millisecond

	var ultimoErro error
	for tentativa := 0; tentativa < maxTentativas; tentativa++ {
		if tentativa > 0 {
			log.Printf("Tentativa %d/%d para %s...", tentativa+1, maxTentativas, metodo)
			time.Sleep(intervaloRetry)
		}

		err := c.client.Call(metodo, args, reply)
		if err == nil {
			return nil // Sucesso
		}

		ultimoErro = err
		log.Printf("Erro na chamada RPC %s: %v", metodo, err)
	}

	return fmt.Errorf("falha após %d tentativas: %v", maxTentativas, ultimoErro)
}

// proximoSequenceNumber gera o próximo número de sequência
func (c *ClienteJogo) proximoSequenceNumber() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sequenceNumber++
	return c.sequenceNumber
}

// EnviarComando envia um comando ao servidor
func (c *ClienteJogo) EnviarComando(tipo, direcao string) error {
	cmd := &game.Comando{
		JogadorID:      c.jogadorID,
		SequenceNumber: c.proximoSequenceNumber(),
		Tipo:           tipo,
		Direcao:        direcao,
		Timestamp:      time.Now(),
	}

	var resp game.RespostaComando
	if err := c.chamarComRetry("ServidorJogo.ProcessarComando", cmd, &resp); err != nil {
		return err
	}

	if !resp.Sucesso {
		return fmt.Errorf("comando falhou: %s", resp.Mensagem)
	}

	c.estadoLocal.StatusMsg = resp.Mensagem
	return nil
}

// AtualizarEstadoLocal busca o estado atual do jogo no servidor
func (c *ClienteJogo) AtualizarEstadoLocal() error {
	req := &game.RequisicaoEstado{
		JogadorID: c.jogadorID,
	}

	var resp game.EstadoJogo
	if err := c.chamarComRetry("ServidorJogo.ObterEstado", req, &resp); err != nil {
		return err
	}

	// Atualiza a lista de jogadores
	c.jogadoresMu.Lock()
	c.jogadores = make(map[string]game.JogadorInfo)
	for _, j := range resp.Jogadores {
		c.jogadores[j.ID] = j
	}
	c.jogadoresMu.Unlock()

	return nil
}

// IniciarAtualizacaoPeriodica inicia uma goroutine que busca atualizações periodicamente
func (c *ClienteJogo) IniciarAtualizacaoPeriodica(intervalo time.Duration) {
	go func() {
		ticker := time.NewTicker(intervalo)
		defer ticker.Stop()

		for range ticker.C {
			if err := c.AtualizarEstadoLocal(); err != nil {
				log.Printf("Erro ao atualizar estado: %v", err)
			}
		}
	}()
}

// ObterJogadores retorna a lista atual de jogadores
func (c *ClienteJogo) ObterJogadores() []game.JogadorInfo {
	c.jogadoresMu.RLock()
	defer c.jogadoresMu.RUnlock()

	jogadores := make([]game.JogadorInfo, 0, len(c.jogadores))
	for _, j := range c.jogadores {
		jogadores = append(jogadores, j)
	}
	return jogadores
}

// Fechar encerra a conexão com o servidor
func (c *ClienteJogo) Fechar() {
	// Envia comando de desconexão
	_ = c.EnviarComando("desconectar", "")

	// Fecha a conexão
	if c.client != nil {
		c.client.Close()
	}
}

// personagemMoverComServidor move o personagem e sincroniza com o servidor
func personagemMoverComServidor(tecla rune, jogo *game.Jogo, cliente *ClienteJogo) {
	dx, dy := 0, 0
	direcao := ""

	switch tecla {
	case 'w':
		dy = -1
		direcao = "w"
	case 'a':
		dx = -1
		direcao = "a"
	case 's':
		dy = 1
		direcao = "s"
	case 'd':
		dx = 1
		direcao = "d"
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy

	// Verifica se o movimento é válido no mapa LOCAL (o cliente mantém o mapa)
	if game.JogoPodeMoverPara(jogo, nx, ny) {
		// Realiza o movimento localmente
		game.JogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny

		// Sincroniza com o servidor
		if err := cliente.EnviarComando("mover", direcao); err != nil {
			log.Printf("Erro ao enviar movimento ao servidor: %v", err)
			jogo.StatusMsg = "Erro ao sincronizar movimento"
		}
	}
}

// personagemInteragirComServidor processa interação e sincroniza com o servidor
func personagemInteragirComServidor(jogo *game.Jogo, cliente *ClienteJogo) {
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)

	// Sincroniza com o servidor
	if err := cliente.EnviarComando("interagir", ""); err != nil {
		log.Printf("Erro ao enviar interação ao servidor: %v", err)
		jogo.StatusMsg = "Erro ao sincronizar interação"
	}
}

// desenharJogadoresRemotosNoMapa desenha outros jogadores no mapa
func desenharJogadoresRemotosNoMapa(jogo *game.Jogo, cliente *ClienteJogo, jogadorLocalID string) {
	jogadores := cliente.ObterJogadores()

	for _, j := range jogadores {
		// Não desenha o jogador local (já desenhado pela interface)
		if j.ID == jogadorLocalID {
			continue
		}

		// Desenha outros jogadores como elementos "Inimigo" temporariamente
		// (pode ser customizado para um elemento específico)
		if j.Posicao.Y >= 0 && j.Posicao.Y < len(jogo.Mapa) &&
			j.Posicao.X >= 0 && j.Posicao.X < len(jogo.Mapa[j.Posicao.Y]) {

			outroJogador := game.Elemento{
				Simbolo:  '◉', // Símbolo para outros jogadores
				Cor:      game.CorCiano,
				CorFundo: game.CorPadrao,
				Tangivel: true,
			}
			game.InterfaceDesenharElemento(j.Posicao.X, j.Posicao.Y, outroJogador)
		}
	}
}

func main() {
	// Inicializa a interface (termbox)
	game.InterfaceIniciar()
	defer game.InterfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo LOCAL (cada cliente mantém seu próprio mapa)
	jogo := game.JogoNovo()
	if err := game.JogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Conecta ao servidor
	// Permite especificar endereço do servidor como segundo argumento
	// Uso: client.exe [mapa.txt] [endereco:porta]
	// Exemplos:
	//   client.exe
	//   client.exe mapa.txt
	//   client.exe mapa.txt 192.168.1.100:8080
	//   client.exe mapa.txt servidor.local:8080
	enderecoServidor := "localhost:8080"
	if len(os.Args) > 2 {
		enderecoServidor = os.Args[2]
		log.Printf("Conectando ao servidor: %s", enderecoServidor)
	}
	nomeJogador := fmt.Sprintf("Jogador_%d", time.Now().Unix()%1000)

	cliente, err := NovoClienteJogo(enderecoServidor, nomeJogador, &jogo)
	if err != nil {
		log.Fatalf("Erro ao conectar ao servidor: %v", err)
	}
	defer cliente.Fechar()

	// Inicia atualização periódica do estado (busca a cada 500ms)
	cliente.IniciarAtualizacaoPeriodica(500 * time.Millisecond)

	// Desenha o estado inicial do jogo
	game.InterfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := game.InterfaceLerEventoTeclado()

		switch evento.Tipo {
		case "sair":
			return
		case "interagir":
			personagemInteragirComServidor(&jogo, cliente)
		case "mover":
			personagemMoverComServidor(evento.Tecla, &jogo, cliente)
		}

		// Redesenha o jogo
		game.InterfaceDesenharJogo(&jogo)

		// Desenha outros jogadores
		desenharJogadoresRemotosNoMapa(&jogo, cliente, cliente.jogadorID)

		game.InterfaceAtualizarTela()
	}
}
