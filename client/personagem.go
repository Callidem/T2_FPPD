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

	// Leitura e escrita do jogador a partir do map
	jogo.mu.Lock()
	jogador, ok := jogo.Jogadores[jogo.localID]
	if !ok {
		jogador = Jogador{UltimoVisitado: Vazio}
	}
	nx, ny := jogador.PosX+dx, jogador.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, nx, ny, jogo.localID)
		jogo.seq++
		jogo.Jogadores[jogo.localID] = jogador
		// Notifica o servidor fora do lock para não bloquear a UI
		go func(id, x, y int) {
			if jogo.RPCClient == nil {
				return
			}
			req := UpdatePositionRequest{ClientID: id, PosX: x, PosY: y}
			// ignora erro, apenas loga para debug
			if err := jogo.RPCClient.Call("UserService.UpdatePosition", &req, &struct{}{}); err != nil {
				log.Printf("RPC erro(UpdatePosition): %v", err)
			}
		}(jogo.localID, nx, ny)
	}
	jogo.mu.Unlock()
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
