package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//Добавить в приложение чата возможность устанавливать клиентам свой никнейм при
//подключении к серверу

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
