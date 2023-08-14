package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

type ChatServer struct {
	clients     map[net.Conn]struct{}
	clientsLock sync.RWMutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[net.Conn]struct{}),
	}
}
func (cs *ChatServer) Start() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Chat server started on :12345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting client:", err)
			continue
		}

		cs.clientsLock.Lock()
		cs.clients[conn] = struct{}{}
		cs.clientsLock.Unlock()

		go cs.handleConnection(conn)
	}
}
func (cs *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	cs.broadcastMessage("Client connected: " + conn.RemoteAddr().String() + "\n")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		cs.broadcastMessage("[" + conn.RemoteAddr().String() + "]: " + message + "\n")
	}

	cs.clientsLock.Lock()
	delete(cs.clients, conn)
	cs.clientsLock.Unlock()

	cs.broadcastMessage("Client disconnected: " + conn.RemoteAddr().String() + "\n")

	if scanner.Err() != nil {
		fmt.Println("Error reading from client:", scanner.Err())
	}
}
func (cs *ChatServer) broadcastMessage(message string) {
	cs.clientsLock.RLock()
	defer cs.clientsLock.RUnlock()

	for conn := range cs.clients {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
func main() {
	chatServer := NewChatServer()
	chatServer.Start()
}
