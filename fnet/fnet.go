package fnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

// Listener is a interface for listening for incoming connections
type Listener interface {
	Listen() error     // Listens for incoming connections
	IsListening() bool // Returns wether this listener is listening or not
	Stop() error       // Stops listening for incoming connections and
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

type Message struct {
	EvtId int32
	PB    proto.Message
}

func (m *Message) Encode() ([]byte, error) {
	// Encode the protobuf message itself
	encoded, err := proto.Marshal(m.PB)
	if err != nil {
		return make([]byte, 0), err
	}

	// Create a new buffer, stuff the event id and the encoded message in it
	buffer := new(bytes.Buffer)
	err = binary.Write(buffer, binary.LittleEndian, m.EvtId)
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

var (
	ErrCantStopListener   = errors.New("Unable to stop listener")
	ErrSliceLengthsDiffer = errors.New("Slice lengths differ")
)

func NewHandlers(callbacks []interface{}, events []int32) ([]Handler, error) {
	if len(callbacks) != len(events) {
		return []Handler{}, ErrSliceLengthsDiffer
	}
	outHandlers := make([]Handler, 0)
	var outErr error
	for k, cb := range callbacks {
		evt := events[k]
		handler, err := NewHandler(cb, evt)
		if err != nil {
			if outErr == nil {
				outErr = err
			} else {
				outErr = errors.New(fmt.Sprintf("%s\nHandler for event %d invalid, skipping: %s", outErr.Error(), evt, err.Error()))
			}
			continue
		}
		outHandlers = append(outHandlers, handler)
	}
	return outHandlers, outErr
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
