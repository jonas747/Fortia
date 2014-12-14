package tcp

import (
	"errors"
	"github.com/jonas747/fortia/fnet"
	"io"
	"net"
	"time"
)

type TCPListner struct {
	Engine    *fnet.Engine
	Addr      string
	Listening bool
}

// Implements fnet.Listener.Listen
func (w *TCPListner) Listen() error {
	listener, err := net.Listen("tcp", w.Addr)
	if err != nil {
		return err
	}

	for {
		con, err := listener.Accept()
		if err != nil {
			return err
		}
		wrappedConn := NewTCPConn(con)
		go w.Engine.HandleConn(wrappedConn)
	}
}

// Implements fnet.Listener.IsListening
func (w *TCPListner) IsListening() bool {
	return w.Listening
}

type TCPConn struct {
	sessionStore *fnet.SessionStore
	conn         net.Conn

	writeChan   chan []byte
	stopWriting chan bool

	isOpen bool
}

func NewTCPConn(c net.Conn) fnet.Connection {
	store := &fnet.SessionStore{make(map[string]interface{})}
	conn := TCPConn{
		sessionStore: store,
		conn:         c,
		writeChan:    make(chan []byte),
		stopWriting:  make(chan bool),
		isOpen:       true,
	}
	return &conn
}

// Implements Connection.Send([]byte)
func (w *TCPConn) Send(b []byte) error {
	if !w.isOpen {
		return errors.New("Cannot call TCPConn.Send() on a closed connection")
	}
	after := time.After(time.Duration(5) * time.Second) // Time out
	select {
	case w.writeChan <- b:
		return nil
	case <-after:
		w.isOpen = false
		w.Close()
		return errors.New("Timed out sending payload to writechan")
	}
}

func (w *TCPConn) Read(buf []byte) error {
	_, err := io.ReadFull(w.conn, buf)
	return err
}

// Implements Connection.Kind() string
func (w *TCPConn) Kind() string {
	return "websocket"
}

// Implements Connection.Close()
func (w *TCPConn) Close() {
	w.isOpen = false
	w.stopWriting <- true
	w.conn.Close()
}

func (w *TCPConn) Open() bool {
	return w.isOpen
}

func (w *TCPConn) GetSessionData() *fnet.SessionStore {
	return w.sessionStore
}

// Implements Connection.Run()
func (w *TCPConn) Run() {
	// Launch the write goroutine
	go w.writer()
}

func (w *TCPConn) writer() {
	for {
		select {
		case m := <-w.writeChan:
			_, err := w.conn.Write(m)
			if err != nil {
				break
			}
		case <-w.stopWriting:
			return
		}
	}
}
