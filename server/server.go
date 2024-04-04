package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	defer listen.Close()
	fmt.Println("Server started on localhost:8081")

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		go handleRequest(connection)
	}
}
func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection.")
			} else {
				fmt.Println("Error reading:", err.Error())
			}
			break // Exit the loop onread error
		}

		if strings.TrimSpace(message) == "[get messages]" {

			messages, err := readMessagesFromFile()
			if err != nil {
				fmt.Println("Error reading messages from file:", err)
				continue
			}
			fmt.Println("Messages are read from file")
			formattedMessages := strings.ReplaceAll(messages, "\n", "[NEWLINE]") // reaplce new line with [NEWLINE]
			conn.Write([]byte(formattedMessages + "\n"))
		} else {

			message = strings.TrimSpace(message)
			if message == "quit" {
				fmt.Println("Client requested to close the connection.")
				break
			}

			fmt.Println("Received data:", message)
			conn.Write([]byte("Message received.\n"))

			appendMessageToFile(message)
		}
	}
}

func appendMessageToFile(message string) {
	file, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func readMessagesFromFile() (string, error) {

	file, err := os.Open("messages.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n"), scanner.Err()
}
