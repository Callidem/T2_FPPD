package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type CreateUserRequest struct {
	Username string
	NewPosX  int
	NewPosY  int
}

type UpdatePositionRequest struct {
	ID   int
	PosX int
	PosY int
}

type User struct {
	ID       int
	Username string
	PosX     int
	PosY     int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: client <server:port>")
		return
	}
	addr := os.Args[1]
	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	var u User
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite seu nome de usuario: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if err := c.Call("UserService.CreateUser", &CreateUserRequest{Username: name, NewPosX: 2, NewPosY: 2}, &u); err != nil {
		log.Fatalf("CreateUser RPC erro: %v", err)
	}
	fmt.Printf("Criado user ID=%d name=%s pos=(%d,%d)\n", u.ID, u.Username, u.PosX, u.PosY)

	fmt.Println("Use WASD para mover (tecla única, sem Enter). Pressione q para sair.")
	posX, posY := u.PosX, u.PosY
	// helper para imprimir a posição sobrescrevendo a mesma linha
	var printMu sync.Mutex
	prevLen := 0
	printPos := func(x, y int) {
		s := fmt.Sprintf("pos=(%d,%d) > ", x, y)
		// se a string atual for menor que a anterior, preenche com espaços (por segurança)
		if prevLen > len(s) {
			s = s + strings.Repeat(" ", prevLen-len(s))
		}
		// limpa a linha e escreve (ESC[2K limpa a linha, mas ESC[K limpa até o fim)
		fmt.Print("\r\x1b[K" + s)
		prevLen = len(s)
	}

	// Goroutine para polling de usuários e impressão da lista
	lastUsersCount := 0
	lastLineWidth := 0
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			var users []User
			if err := c.Call("UserService.ListUsers", &struct{}{}, &users); err != nil {
				// não interrompe a goroutine; apenas registra
				log.Printf("ListUsers erro: %v", err)
				continue
			}

			// Ordena por ID para manter ordem fixa
			sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })

			// Prepara linhas formatadas e calcula largura máxima
			lines := make([]string, len(users))
			curMax := 0
			for i, uu := range users {
				s := fmt.Sprintf("ID [%d] - pos = (%d,%d)", uu.ID, uu.PosX, uu.PosY)
				lines[i] = s
				if len(s) > curMax {
					curMax = len(s)
				}
			}

			printMu.Lock()
			// Reimprime a linha de posição primeiro (mantém cursor alinhado)
			printPos(posX, posY)
			// move para a próxima linha
			fmt.Print("\n")

			// largura usada: garante que sempre escrevemos pelo menos lastLineWidth
			width := curMax
			if width < lastLineWidth {
				width = lastLineWidth
			}

			// Imprime cada usuário em sua própria linha, preenchendo até width
			for _, s := range lines {
				// limpa a linha antes de imprimir para evitar restos
				fmt.Printf("\r\x1b[K%-*s\n", width, s)
			}

			// Limpa linhas remanescentes quando a lista diminui
			printedLines := len(lines)
			if printedLines < lastUsersCount {
				for i := 0; i < (lastUsersCount - printedLines); i++ {
					fmt.Printf("\r\x1b[K%-*s\n", width, "")
				}
			}

			// atualiza contadores
			lastUsersCount = len(users)
			lastLineWidth = width

			// Depois de imprimir a lista, reposiciona o cursor de volta à linha da posição
			// Precisamos subir 1 (linha da posição) + número de linhas impressas
			up := 1 + printedLines
			if up > 0 {
				fmt.Printf("\x1b[%dA", up)
			}
			printMu.Unlock()
		}
	}()

	// Tenta colocar o terminal em modo raw para leitura de uma tecla sem Enter.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		// fallback para leitura por linha se não for possível
		log.Printf("warning: não foi possível ativar modo raw, usando leitura com Enter: %v", err)
		for {
			printMu.Lock()
			printPos(posX, posY)
			printMu.Unlock()
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("read error: %v", err)
			}
			cmd := strings.TrimSpace(line)
			if cmd == "q" {
				fmt.Println("\nsaindo")
				return
			}
			var dx, dy int
			switch cmd {
			case "w":
				dy = -1
			case "s":
				dy = 1
			case "a":
				dx = -1
			case "d":
				dx = 1
			default:
				fmt.Println("tecla invalida (use w/a/s/d ou q)")
				continue
			}
			posX += dx
			posY += dy
			req := UpdatePositionRequest{ID: u.ID, PosX: posX, PosY: posY}
			if err := c.Call("UserService.UpdatePosition", &req, &struct{}{}); err != nil {
				log.Printf("UpdatePosition erro: %v", err)
			}
			// reimprime posicao atualizada
			printMu.Lock()
			printPos(posX, posY)
			printMu.Unlock()
		}
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Leitura byte-a-byte
	buf := make([]byte, 1)
	for {
		printMu.Lock()
		printPos(posX, posY)
		printMu.Unlock()
		n, err := os.Stdin.Read(buf)
		if err != nil {
			log.Fatalf("read raw error: %v", err)
		}
		if n == 0 {
			continue
		}
		b := buf[0]
		var dx, dy int
		switch b {
		case 'w', 'W':
			dy = -1
		case 's', 'S':
			dy = 1
		case 'a', 'A':
			dx = -1
		case 'd', 'D':
			dx = 1
		case 'q', 'Q':
			fmt.Println("\nsaindo")
			return
		default:
			// ignora outras teclas
			continue
		}
		posX += dx
		posY += dy
		req := UpdatePositionRequest{ID: u.ID, PosX: posX, PosY: posY}
		if err := c.Call("UserService.UpdatePosition", &req, &struct{}{}); err != nil {
			log.Printf("UpdatePosition erro: %v", err)
		}
	}
}
