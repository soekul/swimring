package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	GetCmd    = "get"
	PutCmd    = "put"
	DeleteCmd = "del"
	ExitCmd   = "exit"
)

var client *SwimringClient

func main() {
	var serverAddr string
	var serverPort int
	var readLevel, writeLevel string

	flag.StringVar(&serverAddr, "host", "127.0.0.1", "address of server node")
	flag.IntVar(&serverPort, "port", 7000, "port number of server node")
	flag.StringVar(&readLevel, "rl", ALL, "read consistency level")
	flag.StringVar(&writeLevel, "wl", ALL, "write consistency level")
	flag.Parse()

	client = NewSwimringClient(serverAddr, serverPort)
	client.SetReadLevel(readLevel)
	client.SetWriteLevel(writeLevel)

	err := client.Connect()
	if err != nil {
		fmt.Printf("error: unable to connect to %s:%d\n", serverAddr, serverPort)
		os.Exit(0)
	}
	fmt.Printf("connected to %s:%d\n", serverAddr, serverPort)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		if err := processCommand(strings.Trim(command, " \n")); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func processCommand(line string) error {
	tokens := strings.Split(line, " ")

	if len(tokens) == 0 {
		return nil
	}

	switch tokens[0] {
	case GetCmd:
		processGet(tokens)
	case PutCmd:
		processPut(tokens)
	case DeleteCmd:
		processDelete(tokens)
	case ExitCmd:
		os.Exit(0)
	default:
		return errors.New("unknown command")
	}

	return nil
}

func processGet(tokens []string) {
	if len(tokens) != 2 {
		fmt.Println("usage: get <key>")
		return
	}

	val, err := client.Get(tokens[1])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Println(val)
}

func processPut(tokens []string) {
	if len(tokens) != 3 {
		fmt.Println("usage: put <key> <value>")
		return
	}

	err := client.Put(tokens[1], tokens[2])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Println("ok")
}

func processDelete(tokens []string) {
	if len(tokens) != 2 {
		fmt.Println("usage: del <key>")
		return
	}

	err := client.Delete(tokens[1])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Println("ok")
}