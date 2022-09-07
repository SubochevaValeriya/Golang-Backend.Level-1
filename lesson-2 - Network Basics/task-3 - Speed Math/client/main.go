package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//Реализовать игру “Математика на скорость”: сервер генерирует случайное
//выражение с двумя операндами, сохраняет ответ, а затем отправляет выражение всем
//клиентам. Первый клиент, отправивший правильный ответ - побеждает, затем
//генерируется следующее выражение и так далее.

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
