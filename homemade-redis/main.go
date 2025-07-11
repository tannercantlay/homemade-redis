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

// handleConnection handles a single connection to the Redis server.
// It reads commands from the client, processes them, and sends back responses.
// It uses the NewResp and NewWriter functions to read and write RESP values.
// It also uses the Handlers map to find the appropriate handler for each command.
func handleConnection(conn net.Conn) {
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

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
