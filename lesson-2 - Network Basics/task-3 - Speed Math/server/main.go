package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type client chan<- string

var (
	entering    = make(chan client)
	leaving     = make(chan client)
	messages    = make(chan string)
	expressions = make(chan string)
	results     = make(chan int)
	expr        string
	result      int
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	winTable := map[string]int{}
	expr, result = randomMathExpression()
	go func() {
		expressions <- expr
		results <- result
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, winTable)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn, winTable map[string]int) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	ch <- "Please input your name"
	var who string
	inputName := bufio.NewScanner(conn)
	for inputName.Scan() {
		who = inputName.Text()
		break
	}
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch
	speedMath(who, conn, winTable)
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func randomMathExpression() (string, int) {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(10)
	b := rand.Intn(10)

	var operand string
	var result int
	operandR := rand.Intn(4)

	switch operandR {
	case 0:
		operand = "+"
		result = a + b
	case 1:
		operand = "-"
		result = a - b
	case 2:
		operand = "*"
		result = a * b
	case 3:
		operand = "/"
		if b == 0 {
			b = 1
		}
		result = a / b
	}

	return strconv.Itoa(a) + operand + strconv.Itoa(b), result
}

func speedMath(who string, conn net.Conn, winTable map[string]int) {
	var m sync.Mutex
	go distributeExprs(expressions, results)

	messages <- fmt.Sprintf("Enter an answer faster than others: %s=?", expr)
	input := bufio.NewScanner(conn)
	for input.Scan() {
		answer, _ := strconv.Atoi(input.Text())
		if answer == result {
			messages <- fmt.Sprintf("%s wins!", who)
			winTable[who]++
			messages <- "Win table at the moment:"
			for s, i := range winTable {
				messages <- fmt.Sprintf("[%s]: %d", s, i)
			}
			expr, result = randomMathExpression()

			fmt.Println(expr, result)

			go func() {
				m.Lock()
				expressions <- expr
				results <- result
				m.Unlock()
			}()
			speedMath(who, conn, winTable)
		} else {
			messages <- fmt.Sprintf("%s, please try again", who)
		}
	}
}

func distributeExprs(exprs <-chan string, results <-chan int) {
	for expr := range exprs {
		expr = expr
	}

	for result := range results {
		result = result
	}
}
