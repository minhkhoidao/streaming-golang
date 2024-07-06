package main

import (
	"fmt"
	"gin-tcp/tcp"
	"io"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

func sendFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Sprintf("Failed to open file: %v", err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		return fmt.Sprintf("Failed to dial: %v", err)
	}

	buf := make([]byte, 4000)
	totalBytes := 0
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Sprintf("Failed to read file: %v", err)
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			return fmt.Sprintf("Failed to write to connection: %v", err)
		}

		totalBytes += n
	}

	return fmt.Sprintf("Written %d bytes over the network", totalBytes)
}

func main() {
	tcpServer := tcp.FileServer{}
	go tcpServer.StartTCPServer()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		result := sendFile("test.wav")
		c.JSON(200, gin.H{
			"message": result,
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
