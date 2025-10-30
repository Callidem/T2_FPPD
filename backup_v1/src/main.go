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

	c, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Digite seu nome de usuário: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	/*
		fmt.Print("Escolha a cor do seu personagem: \n
		1 - Default
		2 - Black
		3 - Red
		4 - Green
		5 - Yellow
		6 - Blue
		7 - Magenta
		8 - Cyan
		9 - White")
		reader = bufio.NewReader(os.Stdin)
		color, err := reader.ReadString('\n')
	*/
	// 1) CreateUser
	req := CreateUserRequest{Username: username, NewPosX: 16, NewPosY: 5, PlayerColor: CorCinzaEscuro}
	var u User
	if err := c.Call("UserService.CreateUser", &req, &u); err != nil {
		log.Fatal("RPC erro(CreateUser):", err)
	}
	fmt.Printf("Criado: Username=%s pos=(%d,%d) ID=%d\n", u.Username, u.PosX, u.PosY, u.ID)
	localPlayerID := u.ID

	// 2) GetUser
	var got User
	if err := c.Call("UserService.GetUser", &GetUserRequest{ID: u.ID}, &got); err != nil {
		log.Fatal("RPC erro(GetUser):", err)
	}
	fmt.Printf("Buscado: Username = %s pos=(%d,%d)\n", got.Username, got.PosX, got.PosY)

	// 3) ListUsers
	var all []User
	if err := c.Call("UserService.ListUsers", &struct{}{}, &all); err != nil {
		log.Fatal("RPC erro(ListUsers):", err)
	}
	fmt.Println("Todos os usuários:")
	for _, it := range all {
		fmt.Printf("  - Username = %s pos=(%d,%d)\n", it.Username, it.PosX, it.PosY)
	}

	// 4) Polling da posição dos jogadores
	/*go func(myUser string) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			var msgs []Message
			if err := c.Call("UserService.GetNewMessages", &GetNewMessagesRequest{Username: myUser}, &msgs); err != nil {
				log.Printf("RPC erro(GetNewMessages): %v", err)
				continue
			}
			for _, m := range msgs {
				notifyLine(fmt.Sprintf("[%s] %s", m.From, m.Body))
			}
		}
	}(me.Username)
	*/

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

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

	// Start polling goroutine to keep local players positions updated from server
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			var all []User
			if err := c.Call("UserService.ListUsers", &struct{}{}, &all); err != nil {
				// log but keep trying
				log.Printf("RPC erro(ListUsers): %v", err)
				continue
			}
			// Merge remote users into local jogo state (map by ID)
			jogo.mu.Lock()
			for _, ru := range all {
				jog := jogo.Jogadores[ru.ID]
				jog.Nome = ru.Username
				jog.PosX = ru.PosX
				jog.PosY = ru.PosY
				jog.Cor = ru.PlayerColor
				jog.Active = ru.Active
				// preserve UltimoVisitado if previously set in local map
				if jog.UltimoVisitado.simbolo == 0 {
					jog.UltimoVisitado = Vazio
				}
				jogo.Jogadores[ru.ID] = jog
			}
			jogo.mu.Unlock()
		}
	}()

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, localPlayerID); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}
