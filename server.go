package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
)

// ===== Tipos RPC (DEVEM bater com o cliente, nomes exportados) =====
type CreateUserRequest struct {
	NewPosX int
	NewPosY int
}

type GetUserRequest struct{ ID int }

type User struct {
	ID   int
	PosX int
	PosY int
}

// ===== Serviço =====
type UserService struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func (s *UserService) CreateUser(req *CreateUserRequest, resp *User) error {
	log.Printf("[RPC] CreateUser called: NewPosX=%d NewPosY=%d", req.NewPosX, req.NewPosY)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	u := User{ID: s.nextID, PosX: req.NewPosX, PosY: req.NewPosY}
	s.users[u.ID] = u
	*resp = u
	log.Printf("[RPC] CreateUser ok: id=%d pos=(%d,%d)", u.ID, u.PosX, u.PosY)
	return nil
}

func (s *UserService) GetUser(req *GetUserRequest, resp *User) error {
	log.Printf("[RPC] GetUser called: id=%d", req.ID)
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[req.ID]
	if !ok {
		err := errors.New("usuário não encontrado")
		log.Printf("[RPC] GetUser erro: %v", err)
		return err
	}
	*resp = u
	log.Printf("[RPC] GetUser ok: id=%d pos=(%d,%d)", u.ID, u.PosX, u.PosY)
	return nil
}

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

// ===== Utilidades =====
func setupLogging() {
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("não foi possível abrir server.log: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func primaryIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}

func allIPv4s() []string {
	var ips []string
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
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

// Dump reflexivo pra garantir o que o SERVIDOR compilou
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
	setupLogging()

	// (opcional, mas ajuda a travar/diagnosticar)
	gob.Register(CreateUserRequest{})
	gob.Register(GetUserRequest{})
	gob.Register(User{})

	debugDumpServerTypes()

	svc := &UserService{users: make(map[int]User)}
	if err := rpc.RegisterName("UserService", svc); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":8932")
	if err != nil {
		log.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port

	log.Printf("Servidor RPC iniciando na porta %d ...", port)
	if ip, err := primaryIP(); err == nil {
		log.Printf("IP principal (rota padrão): %s:%d", ip, port)
	}
	for _, ip := range allIPv4s() {
		log.Printf("IP local disponível: %s:%d", ip, port)
	}
	fmt.Printf("Servidor RPC pronto em porta %d (veja server.log para detalhes)\n", port)

	var clientSeq uint64
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("[ACCEPT] erro: %v", err)
			continue
		}
		id := atomic.AddUint64(&clientSeq, 1)
		remote := conn.RemoteAddr().String()
		log.Printf("[CLIENT %d] conectado de %s", id, remote)
		go func(id uint64, c net.Conn) {
			defer func() {
				_ = c.Close()
				log.Printf("[CLIENT %d] desconectado (%s)", id, remote)
			}()
			rpc.ServeConn(c)
		}(id, conn)
	}
}
