package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Port     string
	Listener net.Listener
	Client   map[net.Conn]*Client
	wg       sync.WaitGroup
	Shutdown chan struct{}
	Mutex    sync.Mutex
}

func NewServer(p string) (*Server, error) {

	port := fmt.Sprintf(":%s", p)

	// Listen on a specified TCP port
	listener, err := net.Listen("tcp", port)

	server := &Server{
		Port:     port,
		Listener: listener,
		Client:   make(map[net.Conn]*Client),
		Shutdown: make(chan struct{}),
	}

	return server, err
}

func (server *Server) Stop() {
	close(server.Shutdown)
	defer server.Listener.Close()

	done := make(chan struct{})
	go func() {
		server.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish.")
		return
	}
}

func (server *Server) StartServer() {
	server.wg.Add(1)

	fmt.Println("Server is running on :8080")

	server.acceptConnections()
}

func (server *Server) acceptConnections() {
	defer server.wg.Done()

	server.wg.Add(1)

	for {
		select {
		case <-server.Shutdown:
			fmt.Println("Received shutdown signal, stopping server...")
			return
		default:
			conn, err := server.Listener.Accept()
			if err != nil {
				continue
			}

			go server.handleConnection(conn)
		}
	}

}

func (server *Server) handleConnection(conn net.Conn) {
	defer server.wg.Done()
	// Read incoming messages

	server.Mutex.Lock()
	server.Client[conn] = &Client{
		Conn: conn,
		Name: conn.RemoteAddr().String(),
	}
	server.Mutex.Unlock()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	// Read messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		log.Printf("Received from %s: %s", conn.RemoteAddr().String(), message)
		server.broadcast(fmt.Sprintf("%s: %s", conn.RemoteAddr().String(), message), conn)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from client %s: %v", conn.RemoteAddr().String(), err)
	}

	// Handle client disconnection
	server.Mutex.Lock()
	delete(server.Client, conn)
	server.Mutex.Unlock()

	fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
	server.broadcast(fmt.Sprintf("%s has left the chat\n", conn.RemoteAddr().String()), conn)
}

func (server *Server) broadcast(message string, sender net.Conn) {
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	for conn := range server.Client {
		if conn == sender {
			continue
		}
		_, err := fmt.Fprintln(conn, message)
		if err != nil {
			log.Printf("Error broadcasting to %s: %v", conn.RemoteAddr().String(), err)
			conn.Close()
			delete(server.Client, conn)
		}
	}
}
