// Servidor RPC em Go para gerenciar "usuários" (players) com posição e cor.
// Expõe métodos RPC para criar, buscar e listar usuários.
// Também registra logs detalhados (em stdout e server.log) e mostra os IPs locais.

package main

import (
	"encoding/gob" // Registro/serialização de tipos usados em RPC (gob)
	"errors"
	"fmt" // Para o tipo color.RGBA no struct User
	"io"
	"log"
	"net"
	"net/rpc" // Biblioteca RPC padrão do Go
	"os"
	"reflect" // Inspeção reflexiva dos tipos (debug)
	"sync"
	"sync/atomic" // Contador atômico para sequenciar clientes conectados
)

// ===== Tipos RPC (DEVEM bater com o cliente, nomes exportados) =====

// CreateUserRequest: payload da chamada RPC para criar um novo usuário.
// Os campos são exportados (inicial maiúscula) para o gob enxergar.
type CreateUserRequest struct {
	Username string
	NewPosX  int
	NewPosY  int
}

// GetUserRequest: payload para consultar um usuário por ID.
type GetUserRequest struct{ Username string }

// User: objeto de domínio retornado/armazenado pelo serviço.
// Inclui posição, IP (pode ser populado depois) e cor do jogador.
type User struct {
	Username    string
	ID          int
	PosX        int
	PosY        int
	PlayerColor Cor
	Active      bool
}

type ListUsersRequest struct {
}
type ListUsersReply struct {
	Users []User
}

type UpdatePositionRequest struct {
	ClientID int //ou Username/
	Seq      uint64
	PosX     int
	PosY     int
}

type UpdatePositionReply struct {
	OK         bool
	AppliedSeq uint64
}

//Seq = sequence number monotónico gerado pelo cliente para cada comando que modifica estado.
//ClientID pode ser ID numérico retornado no CreateUser ou o Username (desde que seja único).

// ===== Serviço =====

// UserService encapsula o estado (map de usuários) e o próximo ID.
type UserService struct {
	mu           sync.Mutex     // Protege acesso concorrente ao mapa/nextID
	users        map[int]User   // "Banco" em memória dos usuários
	nextID       int            // Autoincremento de IDs
	usernameToID map[string]int // Mapeia username -> user ID para recuperar sessão
	processed    map[int]uint64 // Último Seq processado por ClientID
}

// CreateUser: método RPC para criar usuário.
// Recebe as coordenadas iniciais e devolve o struct User criado.
func (s *UserService) CreateUser(req *CreateUserRequest, resp *User) error {
	log.Printf("[RPC] CreateUser called: username=%s NewPosX=%d NewPosY=%d", req.Username, req.NewPosX, req.NewPosY)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Se já existe sessão para este username, reutiliza o ID e atualiza posição.
	if id, ok := s.usernameToID[req.Username]; ok {
		u := s.users[id]
		s.users[id] = u
		*resp = u
		log.Printf("[RPC] Recovered session for username=%s id=%d", req.Username, id)
		return nil
	}

	// Senão cria nova sessão/usuário e associa ao username.
	s.nextID++
	u := User{ID: s.nextID, PosX: req.NewPosX, PosY: req.NewPosY}
	s.users[u.ID] = u
	s.usernameToID[req.Username] = u.ID

	*resp = u
	log.Printf("[RPC] CreateUser ok: username=%s id=%d pos=(%d,%d)", req.Username, u.ID, u.PosX, u.PosY)
	return nil
}

// GetUser: método RPC para retornar um usuário por ID.
// Se não existir, retorna erro (que chega como erro RPC no cliente).
func (s *UserService) GetUser(req *GetUserRequest, resp *User) error {
	log.Printf("[RPC] GetUser called: Username=%s", req.Username)

	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.usernameToID[req.Username]
	u, ok := s.users[id]
	if !ok {
		err := errors.New("usuário não encontrado")
		log.Printf("[RPC] GetUser erro: %v", err)
		return err
	}

	*resp = u
	log.Printf("[RPC] GetUser ok: id=%d Username=%s pos=(%d,%d)", u.ID, u.PosX, u.PosY)
	return nil
}

// ListUsers: método RPC para listar todos os usuários cadastrados.
// Retorna um slice com cópias dos Users.
func (s *UserService) ListUsers(_ *struct{}, resp *[]User) error {
	log.Printf("[RPC] ListUsers called")

	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	*resp = out

	log.Printf("[RPC] ListUsers ok: count=%d", len(out))
	return nil
}

func (s *UserService) UpdatePosition(req *UpdatePositionRequest, resp *UpdatePositionReply) error {
	log.Printf("[RPC] UpdatePosition called: ClientID=%d Seq=%d PosX=%d PosY=%d", req.ClientID, req.Seq, req.PosX, req.PosY)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Verifica se o ClientID é válido
	if _, ok := s.users[req.ClientID]; !ok {
		err := errors.New("usuário não encontrado")
		log.Printf("[RPC] UpdatePosition erro: %v", err)
		return err
	}

	// Aplica a atualização de posição
	s.users[req.ClientID] = User{
		ID:          req.ClientID,
		PosX:        req.PosX,
		PosY:        req.PosY,
		PlayerColor: s.users[req.ClientID].PlayerColor,
		Active:      s.users[req.ClientID].Active,
	}

	resp.OK = true
	resp.AppliedSeq = req.Seq
	log.Printf("[RPC] UpdatePosition ok: ClientID=%d Seq=%d", req.ClientID, req.Seq)
	return nil
}

// ===== Utilidades =====

// setupLogging configura o log para ir ao mesmo tempo para stdout e para server.log.
// Adiciona prefixos úteis: data, hora com microssegundos e arquivo:linha.
func setupLogging() {
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("não foi possível abrir server.log: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, f) // Escreve nos dois destinos
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

// primaryIP tenta descobrir o IP local usado para sair à Internet (rota padrão)
// abrindo um "dial" UDP para 8.8.8.8:80 (não estabelece conexão real TCP).
func primaryIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}

// allIPv4s enumera os IPv4 das interfaces ativas e não-loopback.
// Útil para mostrar todas as formas de acessar o servidor na LAN.
func allIPv4s() []string {
	var ips []string

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		// Ignora interfaces down e loopback
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			var ip net.IP
			// Extrai IP de IPNet/IPAddr
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// Filtra inválidos/loopback/não IPv4
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip = ip.To4(); ip == nil {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips
}

// debugDumpServerTypes imprime, via reflexão, a estrutura dos tipos principais.
// Serve como "prova" de que o binário do servidor compilou os tipos esperados
// (nome dos campos, se são exportados, e tipos) — muito útil para depurar
// erros de gob no RPC (mismatch de structs entre cliente/servidor).
func debugDumpServerTypes() {
	dump := func(v any) {
		t := reflect.TypeOf(v)
		fmt.Printf("[SERVER] Type %s has %d fields:\n", t.String(), t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Printf("  - %s (exported=%v) type=%v\n", f.Name, f.PkgPath == "", f.Type)
		}
	}

	fmt.Println("[SERVER] ==== Verificando tipos compilados no servidor ====")
	dump(CreateUserRequest{})
	dump(GetUserRequest{})
	dump(User{})
	fmt.Println("[SERVER] ===================================================")
}

func main() {
	// 1) Configura destino/formato dos logs.
	setupLogging()

	// 2) (Opcional, porém recomendado) registra explicitamente os tipos usados em RPC.
	// Isso ajuda o gob a conhecer os tipos antes do tráfego, evitando surpresas.
	gob.Register(CreateUserRequest{})
	gob.Register(GetUserRequest{})
	gob.Register(User{})
	gob.Register(SendMessageRequest{})
	gob.Register(SendMessageReply{})
	gob.Register()

	// 3) Mostra no stdout a "fotografia" dos tipos compilados (debug).
	debugDumpServerTypes()

	// 4) Instancia o serviço com mapas inicializados.
	svc := &UserService{
		users:        make(map[int]User),
		usernameToID: make(map[string]int),
	}
	// 5) Registra o serviço no servidor RPC sob o nome "UserService".
	//    Os métodos exportados com assinatura adequada viram endpoints RPC.
	if err := rpc.RegisterName("UserService", svc); err != nil {
		log.Fatal(err)
	}

	// 6) Abre um listener TCP na porta 8932 (pode ajustar aqui se quiser outra).
	l, err := net.Listen("tcp", ":8932")
	if err != nil {
		log.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port

	// 7) Logs de inicialização e IPs úteis para o cliente conectar.
	log.Printf("Servidor RPC iniciando na porta %d ...", port)

	if ip, err := primaryIP(); err == nil {
		log.Printf("IP principal (rota padrão): %s:%d", ip, port)
	}
	for _, ip := range allIPv4s() {
		log.Printf("IP local disponível: %s:%d", ip, port)
	}

	// Mensagem amigável no stdout (também loga em server.log por conta do MultiWriter).
	fmt.Printf("Servidor RPC pronto em porta %d (veja server.log para detalhes)\n", port)

	// 8) Loop de aceitação de conexões. Cada cliente recebe um ID sequencial.
	var clientSeq uint64
	for {
		conn, err := l.Accept()
		if err != nil {
			// Em caso de erro ao aceitar, apenas registra e continua.
			log.Printf("[ACCEPT] erro: %v", err)
			continue
		}

		// Gera um identificador atômico para a conexão, e captura o endereço remoto.
		id := atomic.AddUint64(&clientSeq, 1)
		remote := conn.RemoteAddr().String()
		log.Printf("[CLIENT %d] conectado de %s", id, remote)

		// 9) Atende o cliente em goroutine separada.
		//    rpc.ServeConn faz o dispatch dos métodos RPC nesta conexão.

		go func(id uint64, c net.Conn) {
			defer func() {
				_ = c.Close()
				log.Printf("[CLIENT %d] desconectado (%s)", id, remote)
			}()
			rpc.ServeConn(c)
		}(id, conn)
	}
}
