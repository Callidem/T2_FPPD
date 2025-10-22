// main.go - Loop principal do jogo (versão single-player original)
package main

import (
	"os"

	"github.com/usrteia-0005/T2_FPPD/pkg/game"
)

func main() {
	// Inicializa a interface (termbox)
	game.InterfaceIniciar()
	defer game.InterfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := game.JogoNovo()
	if err := game.JogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	game.InterfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := game.InterfaceLerEventoTeclado()
		if continuar := game.PersonagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		game.InterfaceDesenharJogo(&jogo)
	}
}
