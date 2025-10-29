// main.go - Loop principal do jogo
package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"strings"
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

	// Redireciona logs para arquivo para não poluir a tela do termbox
	// Ao usar termbox a escrita em stdout/stderr sobrescreve a interface,
	// por isso é melhor gravar logs em arquivo.
	if f, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
		log.SetOutput(f)
		// fechamos no término do main; defer aqui é seguro
		defer f.Close()
	} else {
		// Se não conseguir criar arquivo, continuamos com o logger padrão
		log.Printf("Aviso: não foi possível abrir client.log: %v", err)
	}

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
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("1 - White \n 2 - Cyan \n 3 - Red \n 4 - Green \n 5 - Yellow \n 6 - Blue \n 7 - Magenta\nEscolha a cor do seu personagem: ")
	reader = bufio.NewReader(os.Stdin)
	textColor, err := reader.ReadString('\n')
	intColor := termbox.ColorDefault
	switch textColor {
	case "1\n":
		intColor = termbox.ColorWhite
	case "2\n":
		intColor = termbox.ColorCyan
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
	}

	//Criação do personagem no servidor

	// 1) CreateUser
	req := CreateUserRequest{Username: username, NewPosX: 16, NewPosY: 5, Avatar: '☺', PlayerColor: intColor}
	var u User
	if err := jogo.RPCClient.Call("UserService.CreateUser", &req, &u); err != nil {
		log.Fatal("RPC erro(CreateUser):", err)
	}
	log.Printf("Criado: Username=%s pos=(%d,%d) ID=%d\n", u.Username, u.PosX, u.PosY, u.ID)
	jogo.localID = u.ID

	// Inicializa jogador localmente para aparecer imediatamente na UI
	avatarElem := Elemento{u.Avatar, u.PlayerColor, CorPadrao, true}
	jogo.mu.Lock()
	// salvamos o que havia naquela célula como UltimoVisitado
	if u.PosY >= 0 && u.PosY < len(jogo.Mapa) && u.PosX >= 0 && u.PosX < len(jogo.Mapa[u.PosY]) {
		ultimo := jogo.Mapa[u.PosY][u.PosX]
		jogo.Jogadores[u.ID] = Jogador{
			Nome:           u.Username,
			PosX:           u.PosX,
			PosY:           u.PosY,
			Cor:            u.PlayerColor,
			UltimoVisitado: ultimo,
			Active:         true,
			Avatar:         avatarElem,
		}
		jogo.Mapa[u.PosY][u.PosX] = avatarElem
	} else {
		// posição inválida no mapa: ainda assim criamos jogador com defaults
		jogo.Jogadores[u.ID] = Jogador{Nome: u.Username, PosX: u.PosX, PosY: u.PosY, Cor: u.PlayerColor, Active: true, Avatar: avatarElem, UltimoVisitado: Vazio}
	}
	jogo.mu.Unlock()

	// 2) GetUser
	var got User
	if err := jogo.RPCClient.Call("UserService.GetUser", &GetUserRequest{Username: u.Username}, &got); err != nil {
		log.Fatal("RPC erro(GetUser):", err)
	}
	log.Printf("Buscado: Username = %s pos=(%d,%d)\n", got.Username, got.PosX, got.PosY)

	// 3) ListUsers
	var all []User
	if err := jogo.RPCClient.Call("UserService.ListUsers", &struct{}{}, &all); err != nil {
		log.Fatal("RPC erro(ListUsers):", err)
	}
	log.Println("Todos os usuários:")
	for _, it := range all {
		log.Printf("  - Username = %s pos=(%d,%d)\n", it.Username, it.PosX, it.PosY)
	}

	// 4) Polling da posição dos jogadores
	go func(jogo *Jogo) {
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
				// Se o jogador ainda não existe localmente, inicializa corretamente
				jogo.mu.Lock()
				jogador, ok := jogo.Jogadores[ru.ID]
				jogo.mu.Unlock()
				if !ok {
					// Cria um Elemento de avatar a partir dos dados do servidor
					avatarElem := Elemento{ru.Avatar, ru.PlayerColor, CorPadrao, true}
					// Protege a escrita no mapa e no map de jogadores
					jogo.mu.Lock()
					// Guarda o que estava naquela posição para restaurar depois
					// (verifica limites do mapa por segurança)
					if ny >= 0 && ny < len(jogo.Mapa) && nx >= 0 && nx < len(jogo.Mapa[ny]) {
						ultimo := jogo.Mapa[ny][nx]
						jogo.Jogadores[ru.ID] = Jogador{
							Nome:           ru.Username,
							PosX:           nx,
							PosY:           ny,
							Cor:            ru.PlayerColor,
							UltimoVisitado: ultimo,
							Active:         true,
							Avatar:         avatarElem,
						}
						// Coloca o avatar no mapa
						jogo.Mapa[ny][nx] = avatarElem
					} else {
						jogo.Jogadores[ru.ID] = Jogador{Nome: ru.Username, PosX: nx, PosY: ny, Cor: ru.PlayerColor, Active: true, Avatar: avatarElem, UltimoVisitado: Vazio}
					}
					jogo.mu.Unlock()
				} else {
					// Jogador já existente: move apenas se a posição mudou
					if jogador.PosX != nx || jogador.PosY != ny {
						jogoMoverElemento(jogo, nx, ny, ru.ID)
					} else {
						// Atualiza flags simples (cor/active) se necessário
						jogo.mu.Lock()
						j := jogo.Jogadores[ru.ID]
						j.Active = ru.Active
						j.Cor = ru.PlayerColor
						jogo.Jogadores[ru.ID] = j
						jogo.mu.Unlock()
					}
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
