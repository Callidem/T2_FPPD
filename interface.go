// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorVerde           = termbox.ColorGreen
	CorCiano           = termbox.ColorCyan
	CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// InterfaceIniciar inicializa a interface gráfica usando termbox (exportado)
func InterfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// interfaceIniciar mantido para compatibilidade
func interfaceIniciar() {
	InterfaceIniciar()
}

// InterfaceFinalizar encerra o uso da interface termbox (exportado)
func InterfaceFinalizar() {
	termbox.Close()
}

// interfaceFinalizar mantido para compatibilidade
func interfaceFinalizar() {
	InterfaceFinalizar()
}

// InterfaceLerEventoTeclado lê um evento do teclado e o traduz para um EventoTeclado (exportado)
func InterfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// interfaceLerEventoTeclado mantido para compatibilidade
func interfaceLerEventoTeclado() EventoTeclado {
	return InterfaceLerEventoTeclado()
}

// InterfaceDesenharJogo renderiza todo o estado atual do jogo na tela (exportado)
func InterfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)

	// Força a atualização do terminal
	interfaceAtualizarTela()
}

// interfaceDesenharJogo mantido para compatibilidade
func interfaceDesenharJogo(jogo *Jogo) {
	InterfaceDesenharJogo(jogo)
}

// InterfaceLimparTela limpa a tela do terminal (exportado)
func InterfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// interfaceLimparTela mantido para compatibilidade
func interfaceLimparTela() {
	InterfaceLimparTela()
}

// InterfaceAtualizarTela força a atualização da tela do terminal com os dados desenhados (exportado)
func InterfaceAtualizarTela() {
	termbox.Flush()
}

// interfaceAtualizarTela mantido para compatibilidade
func interfaceAtualizarTela() {
	InterfaceAtualizarTela()
}

// InterfaceDesenharElemento desenha um elemento na posição (x, y) (exportado)
func InterfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.Simbolo, elem.Cor, elem.CorFundo)
}

// interfaceDesenharElemento mantido para compatibilidade
func interfaceDesenharElemento(x, y int, elem Elemento) {
	InterfaceDesenharElemento(x, y, elem)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover e E para interagir. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}
