package network

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	message := []byte("hello world")

	assert.Nil(t, tra.SendMessage(trb.addr, message))
	rpc := <-trb.Consume()

	// buf := &bytes.Buffer{}
	buf := make([]byte, len(message))
	n, err := rpc.Payload.Read(buf)

	assert.Nil(t, err)
	assert.Equal(t, n, len(message))
	assert.Equal(t, buf, message)
	assert.Equal(t, rpc.From, tra.addr)
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)
	trc := NewLocalTransport("C").(*LocalTransport)

	assert.Nil(t, tra.Connect(trb))
	assert.Nil(t, tra.Connect(trc))

	message := []byte("hello world")

	assert.Nil(t, tra.Broadcast(message))

	rpcb := <-trb.Consume()
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, message)

	rpcc := <-trc.Consume()
	b, err = io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, message)

}
