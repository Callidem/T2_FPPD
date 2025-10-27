// main.go - Loop principal do jogo
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: ./jogo <addr:port>")
		return
	}
	addr := os.Args[1]

	// (opcional, ajuda a diagnosticar cedo)
	gob.Register(CreateUserRequest{})
	gob.Register(GetUserRequest{})
	gob.Register(User{})
	gob.Register(ListUsersRequest{})
	gob.Register(ListUsersReply{})
	gob.Register(UpdatePositionRequest{})
	gob.Register(UpdatePositionReply{})

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt" //mapa é fixo neste exemplo
	/*if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}*/

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Conecta ao servidor RPC
	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	jogo.RPCClient = c

	// Solicita ao usuário seu nome de usuário
	fmt.Print("Digite seu nome de usuário: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')

	fmt.Print("1 - Default \n 2 - Black \n 3 - Red \n 4 - Green \n 5 - Yellow \n 6 - Blue \n 7 - Magenta \n 8 - Cyan \n 9 - White\nEscolha a cor do seu personagem: ")
	reader = bufio.NewReader(os.Stdin)
	textColor, err := reader.ReadString('\n')
	intColor := termbox.ColorDefault
	switch textColor {
	case "1\n":
		intColor = termbox.ColorDefault
	case "2\n":
		intColor = termbox.ColorBlack
	case "3\n":
		intColor = termbox.ColorRed
	case "4\n":
		intColor = termbox.ColorGreen
	case "5\n":
		intColor = termbox.ColorYellow
	case "6\n":
		intColor = termbox.ColorBlue
	case "7\n":
		intColor = termbox.ColorMagenta
	case "8\n":
		intColor = termbox.ColorCyan
	case "9\n":
		intColor = termbox.ColorWhite
	}

	//Criação do personagem no servidor

	// 1) CreateUser
	req := CreateUserRequest{Username: username, NewPosX: 16, NewPosY: 5, Avatar: '☺', PlayerColor: intColor}
	var u User
	if err := jogo.RPCClient.Call("UserService.CreateUser", &req, &u); err != nil {
		log.Fatal("RPC erro(CreateUser):", err)
	}
	fmt.Printf("Criado: Username=%s pos=(%d,%d) ID=%d\n", u.Username, u.PosX, u.PosY, u.ID)
	jogo.localID = u.ID

	// 2) GetUser
	var got User
	if err := jogo.RPCClient.Call("UserService.GetUser", &GetUserRequest{Username: u.Username}, &got); err != nil {
		log.Fatal("RPC erro(GetUser):", err)
	}
	fmt.Printf("Buscado: Username = %s pos=(%d,%d)\n", got.Username, got.PosX, got.PosY)

	// 3) ListUsers
	var all []User
	if err := jogo.RPCClient.Call("UserService.ListUsers", &struct{}{}, &all); err != nil {
		log.Fatal("RPC erro(ListUsers):", err)
	}
	fmt.Println("Todos os usuários:")
	for _, it := range all {
		fmt.Printf("  - Username = %s pos=(%d,%d)\n", it.Username, it.PosX, it.PosY)
	}

	// 4) Polling da posição dos jogadores
	go func(jogo *Jogo) {
		for {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				var all []User
				if err := jogo.RPCClient.Call("UserService.ListUsers", &struct{}{}, &all); err != nil {
					// log but keep trying
					log.Printf("RPC erro(ListUsers): %v", err)
					continue
				}
				// Merge remote users into local jogo state (map by ID)
				for _, ru := range all {
					nx, ny := ru.PosX, ru.PosY
					jogoMoverElemento(jogo, nx, ny, ru.ID)
				}
			}
		}
	}(&jogo)

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
