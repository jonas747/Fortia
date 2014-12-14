package fnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"reflect"
)

// Listener is a interface for listening for incoming connections
type Listener interface {
	Listen() error     // Listens for incoming connections
	IsListening() bool // Returns wether this listener is listening or not
}

// Struct which represents an event
type Event struct {
	Id   int32
	Data reflect.Value
}

// Struct which represents a event handler
type Handler struct {
	CallBack interface{}
	Event    int32
	DataType reflect.Type
}

func NewHandler(callback interface{}, evt int32) (Handler, error) {
	err := validateCallback(callback)
	if err != nil {
		return Handler{}, err
	}

	dType := reflect.TypeOf(callback).In(1) // Get the type of the first parameter for later use

	return Handler{
		CallBack: callback,
		Event:    evt,
		DataType: dType,
	}, nil
}

func validateCallback(callback interface{}) error {
	t := reflect.TypeOf(callback)
	if t.Kind() != reflect.Func {
		return errors.New("Callback not a function")
	}
	if t.NumIn() != 2 {
		return errors.New("Callback does not have 2 parameters")
	}
	return nil
}

// Encodes a protocol buffer message with event id and payload length header
func EncodeMessage(msg proto.Message, evtId int32) ([]byte, error) {
	// Encode the protobuf message itself
	encoded, err := proto.Marshal(msg)
	if err != nil {
		return make([]byte, 0), err
	}

	// Create a new buffer, stuff the event id and the encoded message in it
	buffer := new(bytes.Buffer)
	err = binary.Write(buffer, binary.LittleEndian, evtId)
	if err != nil {
		return make([]byte, 0), err
	}

	// Add the length to the buffer
	length := len(encoded)
	err = binary.Write(buffer, binary.LittleEndian, int32(length))
	if err != nil {
		return make([]byte, 0), err
	}

	// Then the actual payload
	_, err = buffer.Write(encoded)
	if err != nil {
		return make([]byte, 0), err
	}

	unread := buffer.Bytes()
	return unread, nil
}
