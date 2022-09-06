package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

//Добавить в приложение рассылки даты/времени возможность отправлять клиентам
//произвольные сообщения из консоли сервера

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		handleConn(conn)
	}
}
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		var msg string
		fmt.Scan(&msg)
		timeNow := time.Now().Format("15:04:05")
		_, err := io.WriteString(c, fmt.Sprintf("Time: %s, Message: %s\n\r", timeNow, msg))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
