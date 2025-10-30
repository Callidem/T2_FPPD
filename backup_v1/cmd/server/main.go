package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// Types used by RPC (must match client types)
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

// UserService minimal implementation
type UserService struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func (s *UserService) CreateUser(req *CreateUserRequest, resp *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	u := User{ID: s.nextID, Username: req.Username, PosX: req.NewPosX, PosY: req.NewPosY}
	s.users[u.ID] = u
	*resp = u
	log.Printf("[RPC] CreateUser: %+v", u)
	return nil
}

func (s *UserService) UpdatePosition(req *UpdatePositionRequest, _ *struct{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[req.ID]
	if !ok {
		return fmt.Errorf("user id %d not found", req.ID)
	}
	u.PosX = req.PosX
	u.PosY = req.PosY
	s.users[req.ID] = u
	log.Printf("[RPC] UpdatePosition: id=%d pos=(%d,%d)", req.ID, req.PosX, req.PosY)
	return nil
}

func (s *UserService) ListUsers(_ *struct{}, resp *[]User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	*resp = out
	return nil
}

func main() {
	gob.Register(CreateUserRequest{})
	gob.Register(UpdatePositionRequest{})
	gob.Register(User{})

	svc := &UserService{users: make(map[int]User)}
	if err := rpc.RegisterName("UserService", svc); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":8932")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on %s", l.Addr().String())
	fmt.Printf("Server ready on %s\n", l.Addr().String())
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
