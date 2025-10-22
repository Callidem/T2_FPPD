// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	Simbolo  rune // Exportado para ser usado pelo cliente
	Cor      Cor
	CorFundo Cor
	Tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
)

// JogoNovo cria e retorna uma nova instância do jogo (exportado)
func JogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// jogoNovo mantido para compatibilidade
func jogoNovo() Jogo {
	return JogoNovo()
}

// JogoCarregarMapa lê um arquivo texto linha por linha e constrói o mapa do jogo (exportado)
func JogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.Simbolo:
				e = Parede
			case Inimigo.Simbolo:
				e = Inimigo
			case Vegetacao.Simbolo:
				e = Vegetacao
			case Personagem.Simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// jogoCarregarMapa mantido para compatibilidade
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	return JogoCarregarMapa(nome, jogo)
}

// JogoPodeMoverPara verifica se o personagem pode se mover para a posição (x, y) (exportado)
func JogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].Tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// jogoPodeMoverPara mantido para compatibilidade
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	return JogoPodeMoverPara(jogo, x, y)
}

// JogoMoverElemento move um elemento para a nova posição (exportado)
func JogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}

// jogoMoverElemento mantido para compatibilidade
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	JogoMoverElemento(jogo, x, y, dx, dy)
}
