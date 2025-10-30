package main

// MESMOS tipos (nomes/maiusculas) do servidor:

// CreateUserRequest: payload da chamada RPC para criar um novo usuário.
// Os campos são exportados (inicial maiúscula) para o gob enxergar.
type CreateUserRequest struct {
	Username    string
	NewPosX     int
	NewPosY     int
	PlayerColor Cor
	Avatar      rune
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
	Avatar      rune
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
