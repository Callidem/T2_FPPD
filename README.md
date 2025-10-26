# Trabalho 2 da cadeira de Fundamentos de Processamento Paralelo e Distribuído
## Decisões do grupo
- X Cliente tem um mapa de structs com ID e posição de cada jogador
- No cliente, o jogador pede o mapa com as posições de todos, mas bloqueia a atualização da sua própria
- O cliente envia somente sua posição para o  servidor
- O servidor recebe de todos e atualiza para todos, conforme solicitado
- Depois de logar, a pessoa fica com uma posição do vetor fixa, atrelada ao seu nome de usuário

## Arquitetura
### Servidor
- Gerência da sessão de jogo
- Mantem estado atual do jogo
- Não mantem uma cópia do mapa do jogo 
- Não controla lógica de movimentação
- Não contem interface gráfica
- Requisições e respostas devem ser impressas no terminal

### Cliente
- Possui interface de interação com o jogo
- Possui lógica de movimentação e de todo o funcionamento do jogo
- Se conecta ao servidor para obter estado atual do jogo
- Goroutine dedicada a buscar atualizações do jogo e atualizar estado local

## Requisitos de Comunicação e Consistência
- TODA comunicação é iniciada pelos clientes
- Servidor apenas responde
- Chamadas de procedimento remoto devem ter tratamento de erro e reexecução em caso de falha
- Garantia de execução única (exactly-once) de comandos que modifiquem estado do servidor
    - Cada comando tem um sequenceNumber
    - Servidor tem o controle de comandos processados por cliente para evitar reexecução em caso de retransmissão

## Jogo de Terminal em Go

Este projeto é um pequeno jogo desenvolvido em Go que roda no terminal usando a biblioteca [termbox-go](https://github.com/nsf/termbox-go). O jogador controla um personagem que pode se mover por um mapa carregado de um arquivo de texto.

### Como funciona

- O mapa é carregado de um arquivo `.txt` contendo caracteres que representam diferentes elementos do jogo.
- O personagem se move com as teclas **W**, **A**, **S**, **D**.
- Pressione **E** para interagir com o ambiente.
- Pressione **ESC** para sair do jogo.

#### Controles

| Tecla | Ação              |
|-------|-------------------|
| W     | Mover para cima   |
| A     | Mover para esquerda |
| S     | Mover para baixo  |
| D     | Mover para direita |
| E     | Interagir         |
| ESC   | Sair do jogo      |

### Como compilar

1. Instale o Go e clone este repositório.
2. Inicialize um novo módulo "jogo":

```bash
go mod init jogo
go get -u github.com/nsf/termbox-go
```

3. Compile o programa:

Linux:

```bash
go build -o jogo
```

Windows:

```bash
go build -o jogo.exe
```

Também é possivel compilar o projeto usando o comando `make` no Linux ou o script `build.bat` no Windows.

### Como executar

1. Certifique-se de ter o arquivo `mapa.txt` com um mapa válido.
2. Execute o programa no termimal:

```bash
./jogo
```

### Estrutura do projeto

- main.go — Ponto de entrada e loop principal
- interface.go — Entrada, saída e renderização com termbox
- jogo.go — Estruturas e lógica do estado do jogo
- personagem.go — Ações do jogador


