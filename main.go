package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type data_user struct {
	user         string
	message      string
	user_conn    string
	name         string
	chat_history []string
}

var (
	data_message = make(chan data_user)
	clients      = make(map[string]net.Conn)
	mutex        sync.Mutex
	new_user     string

	chat_history []string

	join bool
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	go to_chat()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		mutex.Lock()
		clients[conn.RemoteAddr().String()] = conn
		mutex.Unlock()
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
	name_user := scanner.Text()
	join = true
	new_user = name_user
	user := "[" + name_user + "]" + ":"

	fmt.Print("User set to:", user)
	currentTime := time.Now().Format("[2006-01-02 15:04:05]")
	for name_user == "" {

		_, err := conn.Write([]byte("[YOU HAVE TO ENTER YOUR NAME TO START CHAT]:"))
		if err != nil {
			fmt.Println("Error sending name", err)
			return
		}
		if !scanner.Scan() {
			fmt.Println("No name entered by client")
			return
		}
		name_user = scanner.Text()
		new_user = name_user
		user = "[" + name_user + "]" + ":"

		fmt.Print("User set to:", user)
	}
	data := data_user{user: currentTime + user, user_conn: conn.RemoteAddr().String(), name: name_user}

	data_message <- data
	for scanner.Scan() {
		mutex.Lock()
		message := scanner.Text()
		if !join {
			chat_history = append(chat_history, data.user, data.message+"\n")
		}
		data = data_user{user: currentTime + user, message: message, user_conn: conn.RemoteAddr().String(), chat_history: chat_history}

		data_message <- data
		mutex.Unlock()
		// fmt.Println(user, message)
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}

func to_chat() {
	for data := range data_message {
		mutex.Lock()
		for add, user := range clients {
			if add != data.user_conn {
				if new_user != "" && join {
					_, err := user.Write([]byte("\n" + new_user + " has joined our chat..."))
					if err != nil {
						continue
					}
				} else {
					_, err := user.Write([]byte("\n" + data.user + data.message))
					if err != nil {
						continue
					}
				}
			} else if add == data.user_conn {
				if new_user != "" && join {
					for _, mesg := range chat_history {
						_, err := user.Write([]byte(mesg))
						if err != nil {
							continue
						}
					}
					new_user = ""
					join = false
				} else {
					_, err := user.Write([]byte(data.user))
					if err != nil {
						continue
					}
				}
			}
		}
		mutex.Unlock()
	}
}
