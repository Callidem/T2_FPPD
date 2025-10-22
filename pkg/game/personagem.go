// personagem.go - Funções para movimentação e ações do personagem
package game

import "fmt"

// PersonagemMover atualiza a posição do personagem com base na tecla pressionada (WASD)
func PersonagemMover(tecla rune, jogo *Jogo) {
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

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if JogoPodeMoverPara(jogo, nx, ny) {
		JogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// PersonagemInteragir define o que ocorre quando o jogador pressiona a tecla de interação
func PersonagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// PersonagemExecutarAcao processa o evento do teclado e executa a ação correspondente
func PersonagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		PersonagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		PersonagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
