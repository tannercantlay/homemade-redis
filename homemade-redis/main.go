package main

import (
    "fmt"
    "io"
    "os"
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
    conn, err := l.Accept()
    if err != nil {
        fmt.Println("Error accepting connection:", err)
        return
    }

    defer conn.Close()

    for{
        buf := make([]byte, 1024)

        //read message from client
        n, err := conn.Read(buf)
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client disconnected")
                break
            }
            fmt.Println("Error reading from connection:", err)
            os.Exit(1)
        }
        fmt.Printf("Received: %s\n", string(buf[:n]))

        //ignore request and send back a PONG
        conn.Write([]byte("+OK\r\n"))

    }
}
