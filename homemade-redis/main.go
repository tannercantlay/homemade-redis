package main

import (
    "fmt"
    "io"
    "net"
)

func main() {
    fmt.Println("Listening on port: 6379")

    //creating a new server
    l, err := net.Listen("tcp", ":6379")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    // Listen for connections
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    for {
        buf := make([]byte, 1024)
        n, err := conn.Read(buf)
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client disconnected")
                break
            }
            fmt.Println("Error reading from connection:", err)
            return
        }
        fmt.Printf("Received: %s\n", string(buf[:n]))
        conn.Write([]byte("+OK\r\n"))
    }
}