package fnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

// The networking engine. Holds togheter all the connections and handlers
type Engine struct {
	ConnCloseChan   chan Connection
	EmitConnOnClose bool

	registerConn   chan Connection // Channel for registering new connections
	unregisterConn chan Connection // Channel for unregistering connections
	broadcastChan  chan []byte     // Channel for broadcasting messages to all connections

	listeners   []Listener          // Slice Containing all listeners
	handlers    map[int32]Handler   // Map with all the event handlers, their id's as keys
	connections map[Connection]bool // Map containing all conncetions
	ErrChan     chan error
}

func NewEngine() *Engine {
	return &Engine{
		ConnCloseChan:  make(chan Connection),
		registerConn:   make(chan Connection),
		unregisterConn: make(chan Connection),
		broadcastChan:  make(chan []byte),
		listeners:      make([]Listener, 0),
		handlers:       make(map[int32]Handler),
		connections:    make(map[Connection]bool),
		ErrChan:        make(chan error),
	}
}

func (e *Engine) Broadcast(msg []byte) {
	e.broadcastChan <- msg
}

// Adds a listener and make it start listening for incoming connections
func (e *Engine) AddListener(listener Listener) {
	e.listeners = append(e.listeners, listener)
	if !listener.IsListening() {
		go func() {
			err := listener.Listen()
			e.ErrChan <- err
		}()
	}
}

// Handles connections
func (e *Engine) HandleConn(conn Connection) {
	conn.Run()
	e.registerConn <- conn
	for {
		err := e.readMessage(conn)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	e.unregisterConn <- conn
}

func (e *Engine) readMessage(conn Connection) error {
	// start with receving the evt id and payload length
	header := make([]byte, 8)
	err := conn.Read(header)
	if err != nil {
		return err
	}
	evtId, pl, err := readHeader(header)
	if err != nil {
		return err
	}

	payload := make([]byte, pl)
	if pl != 0 {
		err = conn.Read(payload)
		if err != nil {
			return err
		}
	}
	return e.handleMessage(evtId, payload, conn)
}

func readHeader(header []byte) (evtId int32, payloadLength int32, err error) {
	buf := bytes.NewReader(header)

	err = binary.Read(buf, binary.LittleEndian, &evtId)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &payloadLength)
	if err != nil {
		return
	}

	return
}

// Retrieves the event id, decodes the data and calls the callback
func (e *Engine) handleMessage(evtId int32, payload []byte, conn Connection) error {
	handler, found := e.handlers[evtId]
	if !found {
		return errors.New(fmt.Sprintf("No handler found for %d", evtId))
	}

	decoded := reflect.New(handler.DataType).Interface().(proto.Message) // We use reflect to unmarshal the data into the appropiate type
	if len(payload) > 0 {
		err := proto.Unmarshal(payload, decoded)
		if err != nil {
			return err
		}
	}
	// ready the function
	funcVal := reflect.ValueOf(handler.CallBack)
	decVal := reflect.Indirect(reflect.ValueOf(decoded)) // decoded is a pointer, so we get the value it points to
	connVal := reflect.ValueOf(conn)
	resp := funcVal.Call([]reflect.Value{connVal, decVal}) // Call it

	if len(resp) == 0 {
		return nil
	}

	returnVal := resp[0]
	if returnVal.Kind() == reflect.Slice {
		inter := returnVal.Interface()
		responseRaw := inter.([]byte)
		conn.Send(responseRaw)
	}

	return nil
}

// Adds a handler
func (e *Engine) AddHandler(handler Handler) {
	e.handlers[handler.Event] = handler
}

func (e *Engine) ListenChannels() {
	for {
		select {
		case d := <-e.registerConn: //Register a connection
			e.connections[d] = true
		case d := <-e.unregisterConn: //Unregister a connection
			delete(e.connections, d)
			if e.EmitConnOnClose {
				e.ConnCloseChan <- d
			}
		case msg := <-e.broadcastChan: //Broadcast a message to all connections
			for conn := range e.connections {
				conn.Send(msg)
			}
		}
	}
}

func (e *Engine) NumClients() int {
	return len(e.connections)
}
