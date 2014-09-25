package log

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Server struct {
	Out chan LogMsg
}

func NewServer(laddr string) (*Server, error) {
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Out: make(chan LogMsg, 0),
	}
	go server.listen(listener)

	return server, nil
}

func (s *Server) listen(listener net.Listener) {
	fmt.Println("Listening!")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming logger connection", err)
			continue
		}
		go s.handleconn(conn)
	}
}

// TODO: Check if the first 8 bytes is actually the msg length
func (s *Server) handleconn(conn net.Conn) {
	for {
		// Get the length of the message
		msgLengthRaw := make([]byte, 8)
		_, err := conn.Read(msgLengthRaw)
		if err != nil {
			fmt.Println("Error reading log message: ", err)
			return
		}
		lBuffer := bytes.NewBuffer(msgLengthRaw)
		msgLength := int64(0)
		err = binary.Read(lBuffer, binary.LittleEndian, &msgLength)
		if err != nil {
			fmt.Println("Error converting byte array to int:", err)
			continue
		}

		// Read the rest of the message
		msgRaw := make([]byte, msgLength)
		_, err = conn.Read(msgRaw)
		if err != nil {
			fmt.Println("Error reading log message body: ", err)
			return
		}

		// Decode the final json
		var msg LogMsg
		err = json.Unmarshal(msgRaw, &msg)
		if err != nil {
			fmt.Println("Error decoding json: ", err)
			continue
		}

		// Send it out on a channel
		s.Out <- msg
	}
}
