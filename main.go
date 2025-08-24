package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type User struct {
	name    string
	message string
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New client connected:", conn.RemoteAddr())
	file, err2 := os.Open("promopt.txt")
	if err2 != nil {
		return
	}
	promot := bufio.NewScanner(file)
	for promot.Scan() {
		_, err := conn.Write([]byte(promot.Text()))
		conn.Write([]byte("\n"))
		if err != nil {
			fmt.Println("Error sending name prompt:", err)
			return
		}
	}
	scanner := bufio.NewScanner(conn)
	_, err := conn.Write([]byte("[ENTER YOUR NAME]:"))
	if err != nil {
		fmt.Println("Error sending name", err)
		return
	}
	if !scanner.Scan() {
		fmt.Println("No name entered by client")
		return
	}
	user := "[" + scanner.Text() + "]" + ":"
	fmt.Print("User set to:", user)
	currentTime := time.Now().Format("[2006-01-02 15:04:05]")
	conn.Write([]byte(currentTime))
	conn.Write([]byte(user))
	data_user := make(chan string, 10)
	for scanner.Scan() {

		message := scanner.Text()
		data_user <- user + message
		fmt.Println(user, message)
		msg := <-data_user
		conn.Write([]byte(currentTime))
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}
