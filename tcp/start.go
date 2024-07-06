package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type FileServer struct {
}

func (fs *FileServer) StartTCPServer() {
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.ReadLoop(con)
	}
}

func (fs *FileServer) ReadLoop(con net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(con, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, con, 4000)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.Bytes())
		fmt.Printf("received %d bytes over the network\n", n)
	}
}
