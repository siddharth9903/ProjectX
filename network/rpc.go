package network

import (
	"ProjectX/core"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type MessageType byte

const (
	MessageTypeTx  MessageType = 0x1
	MessageTypeBlock MessageType = 0x2
)

type Message struct {
	Header MessageType
	Data []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{Header: t, Data: data}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type RPCHandler interface{
	HandleRPC(rpc RPC) error
}

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(rpc RPC) (*DecodedMessage, error)


func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error){
		msg := &Message{}

		dec := gob.NewDecoder(rpc.Payload)
		err := dec.Decode(&msg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
		}

		logrus.WithFields(logrus.Fields{
			"from": rpc.From,
			"type": msg.Header,
		}).Debug("New incoming message")

		switch msg.Header{
			case MessageTypeTx:
				tx := new(core.Transaction)
				if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
					return nil,err 
				}
			return &DecodedMessage{
				From :rpc.From,
				Data : tx,
			}, nil
		default:
			return nil,fmt.Errorf("invalid message header %x", msg.Header)
	}

}

type RPCProcessor interface{
	ProcessMessage(msg *DecodedMessage) error
}

