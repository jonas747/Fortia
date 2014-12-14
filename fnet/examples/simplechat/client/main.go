package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/fnet"
	"github.com/jonas747/fortia/fnet/examples/simplechat"
	"github.com/jonas747/fortia/fnet/tcp"
	"os"
)

var addr = flag.String("addr", "jonas747.com:7447", "The address to listen on")

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Running simplechat client!")
	engine := fnet.NewEngine()

	hUserMsg, err := fnet.NewHandler(HandleMsg, int32(simplechat.Events_MESSAGE))
	panicErr(err)

	engine.AddHandler(hUserMsg)
	conn, err := tcp.Dial(*addr)
	panicErr(err)

	go engine.ListenChannels()
	go engine.HandleConn(conn)

	fmt.Println("Enter your name:")
	name := ""
	fmt.Scanln(&name)

	msg := &simplechat.User{
		Name: proto.String(name),
	}

	encoded, err := fnet.EncodeMessage(msg, int32(simplechat.Events_USERJOIN))
	panicErr(err)
	err = conn.Send(encoded)
	fmt.Println("message encoded and all!")
	panicErr(err)

	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		panicErr(err)
		line = line[:len(line)-1]
		msg := &simplechat.ChatMsg{
			Msg: proto.String(line),
		}

		encoded, err := fnet.EncodeMessage(msg, int32(simplechat.Events_MESSAGE))
		panicErr(err)
		err = conn.Send(encoded)
		panicErr(err)
	}
}

func HandleMsg(conn fnet.Connection, msg simplechat.ChatMsg) {
	fmt.Printf("[%s]: %s\n", msg.GetFrom(), msg.GetMsg())
}
