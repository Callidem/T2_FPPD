// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
	"log"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {

	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1 // Move para cima
	case 'a':
		dx = -1 // Move para a esquerda
	case 's':
		dy = 1 // Move para baixo
	case 'd':
		dx = 1 // Move para a direita
	}

	// Lê posição atual do jogador sob lock curto
	jogo.mu.Lock()
	jogador, ok := jogo.Jogadores[jogo.localID]
	if !ok {
		jogador = Jogador{UltimoVisitado: Vazio}
	}
	currX, currY := jogador.PosX, jogador.PosY
	jogo.mu.Unlock()

	nx, ny := currX+dx, currY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		// Move no estado local (jogoMoverElemento faz locking interno)
		jogoMoverElemento(jogo, nx, ny, jogo.localID)

		// Incrementa o sequence number de forma protegida e captura seu valor
		jogo.mu.Lock()
		jogo.seq++
		seq := jogo.seq
		jogo.mu.Unlock()

		// Notifica o servidor fora do lock para não bloquear a UI
		go func(id, x, y int, s uint64) {
			if jogo.RPCClient == nil {
				return
			}
			req := UpdatePositionRequest{ClientID: id, Seq: s, PosX: x, PosY: y}
			var rep UpdatePositionReply
			if err := jogo.RPCClient.Call("UserService.UpdatePosition", &req, &rep); err != nil {
				log.Printf("RPC erro(UpdatePosition) seq=%d: %v", s, err)
				return
			}
			// Log básico de confirmação para ajudar a debugar
			if !rep.OK {
				log.Printf("UpdatePosition not applied: seq=%d applied=%d OK=%v", s, rep.AppliedSeq, rep.OK)
			} else {
				log.Printf("UpdatePosition applied: seq=%d applied=%d", s, rep.AppliedSeq)
			}
		}(jogo.localID, nx, ny, seq)
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.mu.Lock()
	jogador, ok := jogo.Jogadores[jogo.localID]
	if !ok {
		jogador = Jogador{UltimoVisitado: Vazio}
	}
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogador.PosX, jogador.PosY)
	jogo.mu.Unlock()
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
