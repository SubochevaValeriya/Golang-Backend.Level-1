package main

import "fmt"

func main() {
	expressions := make(chan string)
	results := make(chan int)
	go func() {
		expressions <- "s"
		results <- 1
	}()
	expr := <-expressions
	res := <-results
	fmt.Println(expr, res)
}
