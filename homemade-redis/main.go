package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port: 6379")

	//creating a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	// Create a new AOF instance
	// This will create the AOF file if it does not exist
	// and start a goroutine to sync it to disk every second.
	// The AOF file will be named "database.aof".
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// Read existing AOF data and register handlers
	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	// Listen for connections
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, aof)
	}
}

// handleConnection handles a single connection to the Redis server.
// It reads commands from the client, processes them, and sends back responses.

func handleConnection(conn net.Conn, aof *Aof) {
	defer conn.Close()
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		// Create a new writer to send responses back to the client
		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		// If the command is SET or HSET, write it to the AOF file
		// This ensures that the command is persisted to disk.
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}
		// Call the handler for the command
		// and write the result back to the client.
		result := handler(args)
		writer.Write(result)
	}
}
