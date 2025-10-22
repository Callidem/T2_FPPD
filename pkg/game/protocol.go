// protocol.go - Estruturas de protocolo para comunicação cliente-servidor
package game

import "time"

// Representa a posição de um jogador no mapa
type Posicao struct {
	X int
	Y int
}

// Informações sobre um jogador na sessão
type JogadorInfo struct {
	ID      string  // Identificador único do jogador
	Nome    string  // Nome do jogador
	Posicao Posicao // Posição atual no mapa
	Vidas   int     // Número de vidas restantes
	Ativo   bool    // Se o jogador está ativo na sessão
}

// Comando enviado pelo cliente ao servidor
type Comando struct {
	JogadorID      string    // ID do jogador que envia o comando
	SequenceNumber int64     // Número de sequência para garantir exactly-once
	Tipo           string    // Tipo do comando: "mover", "interagir", "conectar", "desconectar"
	Direcao        string    // Direção do movimento: "w", "a", "s", "d"
	Timestamp      time.Time // Timestamp do comando
}

// Resposta do servidor para um comando
type RespostaComando struct {
	Sucesso        bool      // Se o comando foi executado com sucesso
	Mensagem       string    // Mensagem de status/erro
	SequenceNumber int64     // Número de sequência do comando processado
	Timestamp      time.Time // Timestamp da resposta
}

// Requisição para obter o estado atual do jogo
type RequisicaoEstado struct {
	JogadorID string // ID do jogador solicitante
}

// Estado atual do jogo retornado pelo servidor
type EstadoJogo struct {
	Jogadores []JogadorInfo // Lista de todos os jogadores ativos
	Timestamp time.Time     // Timestamp do estado
}

// Requisição para conectar um novo jogador
type RequisicaoConexao struct {
	Nome           string  // Nome do jogador
	PosicaoInicial Posicao // Posição inicial sugerida
}

// Resposta de conexão
type RespostaConexao struct {
	Sucesso   bool        // Se a conexão foi bem-sucedida
	JogadorID string      // ID atribuído ao jogador
	Mensagem  string      // Mensagem de status
	Jogador   JogadorInfo // Informações completas do jogador
}
