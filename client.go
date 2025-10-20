package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os"
)

// MESMOS tipos (nomes/maiusculas) do servidor:
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: ./client <addr:port>")
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

	// 1) CreateUser
	req := CreateUserRequest{NewPosX: rand.Intn(100), NewPosY: rand.Intn(100)}
	var u User
	if err := c.Call("UserService.CreateUser", &req, &u); err != nil {
		log.Fatal("RPC erro(CreateUser):", err)
	}
	fmt.Printf("Criado: ID=%d pos=(%d,%d)\n", u.ID, u.PosX, u.PosY)

	// 2) GetUser
	var got User
	if err := c.Call("UserService.GetUser", &GetUserRequest{ID: u.ID}, &got); err != nil {
		log.Fatal("RPC erro(GetUser):", err)
	}
	fmt.Printf("Buscado: ID=%d pos=(%d,%d)\n", got.ID, got.PosX, got.PosY)

	// 3) ListUsers
	var all []User
	if err := c.Call("UserService.ListUsers", &struct{}{}, &all); err != nil {
		log.Fatal("RPC erro(ListUsers):", err)
	}
	fmt.Println("Todos os usuários:")
	for _, it := range all {
		fmt.Printf("  - ID=%d pos=(%d,%d)\n", it.ID, it.PosX, it.PosY)
	}
}
