package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jonas747/fortia/fnet"
	"github.com/jonas747/fortia/fnet/examples/simplechat"
	"github.com/jonas747/fortia/fnet/tcp"
)

var addr = flag.String("addr", ":7447", "The address to listen on")

var engine *fnet.Engine

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	fmt.Println("Running simplechat server!")
	engine = fnet.NewEngine()

	hUserJoin, err := fnet.NewHandler(HandleUserJoin, int32(simplechat.Events_USERJOIN))
	hUserLeave, err2 := fnet.NewHandler(HandleUserLeave, int32(simplechat.Events_USERLEAVE))
	hUserMsg, err3 := fnet.NewHandler(HandleSendMsg, int32(simplechat.Events_MESSAGE))
	panicErr(err)
	panicErr(err2)
	panicErr(err3)

	engine.AddHandler(hUserJoin)
	engine.AddHandler(hUserLeave)
	engine.AddHandler(hUserMsg)

	listener := &tcp.TCPListner{
		Engine: engine,
		Addr:   *addr,
	}
	go engine.ListenChannels()
	go engine.AddListener(listener)

	engine.EmitConnOnClose = true

	for {
		c := <-engine.ConnCloseChan
		name, ok := c.GetSessionData().Get("name")
		if ok {
			nameStr := name.(string)
			chatMsg := fmt.Sprintf("\"%s\" Has left!", nameStr)
			msg := &simplechat.ChatMsg{
				From: proto.String("server"),
				Msg:  proto.String(chatMsg),
			}
			encoded, err := fnet.EncodeMessage(msg, int32(simplechat.Events_MESSAGE))
			panicErr(err)
			engine.Broadcast(encoded)
		}
	}
}

func HandleUserJoin(conn fnet.Connection, user simplechat.User) {
	name := user.GetName()
	conn.GetSessionData().Set("name", name)
	msg := &simplechat.ChatMsg{
		From: proto.String("server"),
		Msg:  proto.String("\"" + name + "\" Joined!"),
	}
	encoded, err := fnet.EncodeMessage(msg, int32(simplechat.Events_MESSAGE))
	panicErr(err)

	engine.Broadcast(encoded)
}

func HandleUserLeave(conn fnet.Connection, user simplechat.User) {
	fmt.Println("UserLeave!")
}

func HandleSendMsg(conn fnet.Connection, msg simplechat.ChatMsg) {
	name, _ := conn.GetSessionData().Get("name")
	nameStr := name.(string)
	message := &simplechat.ChatMsg{
		From: proto.String(nameStr),
		Msg:  proto.String(msg.GetMsg()),
	}
	encoded, err := fnet.EncodeMessage(message, int32(simplechat.Events_MESSAGE))
	panicErr(err)

	engine.Broadcast(encoded)
}
