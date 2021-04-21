package main

import (
	"fmt"
)

func callerA(c chan string) {
	c <- "Hello World!"
	close(c)
}

func callerB(c chan string) {
	c <- "Hola Mundo!"
	close(c)
}

func main() {
	a, b := make(chan string), make(chan string)
	go callerA(a)
	go callerB(b)
	var msg string
	openA, openB := true, true
	for openA || openB {
		select {
		case msg, openA = <-a:
			if openA {
				fmt.Printf("%s from A\n", msg)
			}
		case msg, openB = <-b:
			if openB {
				fmt.Printf("%s from B\n", msg)
			}
		}
	}
}
