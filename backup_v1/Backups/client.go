package main

// MESMOS tipos (nomes/maiusculas) do servidor:
type CreateUserRequest struct {
	Username    string
	NewPosX     int
	NewPosY     int
	PlayerColor Cor
}
type GetUserRequest struct{ ID int }

type User struct {
	Username    string
	ID          int
	PosX        int
	PosY        int
	PlayerColor Cor
	Active      bool
}
