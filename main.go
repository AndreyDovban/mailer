package main

import "fmt"

func main() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("work")
}
