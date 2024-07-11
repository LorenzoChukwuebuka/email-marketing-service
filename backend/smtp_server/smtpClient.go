package smtpserver

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"
)

func SMTPClient() {
	conn, err := net.Dial("tcp", "localhost:1025")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Read server greeting
	readResponse(reader)

	// Send EHLO
	sendCommand(writer, "EHLO localhost\r\n")
	readResponse(reader)

	// Authenticate
	sendCommand(writer, "AUTH PLAIN\r\n")
	readResponse(reader)

	// Send Base64 encoded username and password
	auth := base64.StdEncoding.EncodeToString([]byte("\x00username\x00password"))
	sendCommand(writer, auth+"\r\n")
	readResponse(reader)

	// Keep-alive loop
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send NOOP command
			sendCommand(writer, "NOOP\r\n")
			readResponse(reader)
		}
	}
}

func sendCommand(writer *bufio.Writer, cmd string) {
	_, err := writer.WriteString(cmd)
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func readResponse(reader *bufio.Reader) {
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response)
}
