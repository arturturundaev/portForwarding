package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

/**
Read file with config and forwarding
*/
func main() {
	data := getDataFromFile()
	fmt.Println("Start listening: " + data["FROM_ADDRESS"] + ":" + data["FROM_PORT"])
	ln, err := net.Listen("tcp", data["FROM_ADDRESS"]+":"+data["FROM_PORT"])
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(conn, data["TO_ADDRESS"], data["TO_PORT"])
	}

}

/**
Return data from .config file like key => value
*/
func getDataFromFile() map[string]string {
	file, err := os.Open(".config.txt")
	if err != nil {
		log.Fatal(err)
	}
	var data = make(map[string]string, 4)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		keyAndValue := strings.Split(line, ":")
		data[keyAndValue[0]] = keyAndValue[1]
	}

	return data
}

func handleRequest(conn net.Conn, address string, port string) {
	fmt.Println("new client")

	proxy, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		panic(err)
	}

	fmt.Println("proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}

/**
PANIC! IT'S A ... ERROR

func fail(error error) {
	if error != nil {
		panic(error)
	}
}*/
