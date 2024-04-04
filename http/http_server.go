package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func sendMessageToTCPServer(message string) (string, error) {

	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	fmt.Fprintln(conn, message)

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func main() {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			message := r.FormValue("message")

			response, err := sendMessageToTCPServer(message)
			if err != nil {
				http.Error(w, "Failed to send message to TCP server", http.StatusInternalServerError)
				log.Println("Error sending message to TCP server:", err)
				return
			}

			fmt.Fprintf(w, "TCP Server said: %s", response)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/retrieve-messages", func(w http.ResponseWriter, r *http.Request) {
		response, err := sendMessageToTCPServer("[get messages]")
		if err != nil {
			http.Error(w, "Failed to retrieve messages from TCP server", http.StatusInternalServerError)
			return
		}
		response = strings.ReplaceAll(response, "[NEWLINE]", "\n") // Replace the [NEWLINE] placeholder with actual newlines

		fmt.Fprintf(w, "Retrieved Messages:\n%s", response)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("HTTP server started on localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
