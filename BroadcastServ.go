package main

import (
    "net"
    "log"
    "sync"
    "io"
)

const (
    PACKET_SIZE = 19
)

// Array of tcp socket clients
// RW protection
var (
    mu sync.RWMutex 
    listOfClients []*net.TCPConn
)


func main() {
    // Host TCP server in localhost with port 8100

    server, _ := net.Listen("tcp", "0.0.0.0:8100")

    // waiting for new connections
    for {
        conn, _ := server.Accept()

        tcp := conn.( * net.TCPConn)

        tcp.SetNoDelay(true)
        tcp.SetKeepAlive(true)
        tcp.SetLinger(1)
        tcp.SetReadBuffer(10000)
        tcp.SetWriteBuffer(10000)

        // appending the client socket to array/slice
        AddClient(tcp)

        // calling handleRequest using goroutines to run in seperate thread
        go HandleRequest( *tcp)
    }
}

func AddClient(tcp *net.TCPConn) {
    mu.Lock()
    defer mu.Unlock();
    log.Printf("New Connection: %s", tcp.RemoteAddr().String())
    listOfClients = append(listOfClients, tcp);
}

func RemoveClient(tcp * net.TCPConn) {
    mu.Lock()
    defer mu.Unlock()
        // Why do slices not have a helper function for this ?!
    for index := 0;
    index < len(listOfClients);
    index++{
        if listOfClients[index].RemoteAddr() == tcp.RemoteAddr() {
            // Close connection
            listOfClients[index].Close()
                // Slice madness for fast removal
            listOfClients[index] = listOfClients[len(listOfClients) - 1]
            listOfClients = listOfClients[: len(listOfClients) - 1]
            break
        }
    }
}

func HandleRequest(conn net.TCPConn) {

    // waiting for new messages

    for {

        // reading the message from socket buffer

        data := make([] byte, PACKET_SIZE)
        _,
        err := conn.Read(data)

        if err != nil {
            if err == io.EOF {
                log.Printf("Connection Closed by: %s", conn.RemoteAddr().String())
                RemoveClient( & conn)
                return
            } else {
                log.Printf("Error: %+v on Connection", err.Error(), conn.RemoteAddr().String())
            }
            return
        }

        // broadcasting to all clients
        go broadCast(data)
    }
}


func broadCast(data[] byte) {
    // iterating over the listOfClients
    mu.Lock()
    defer mu.Unlock()
    for i := range listOfClients {
        // sending message to all clients in the socket
        listOfClients[i].Write(data)
    }
}
