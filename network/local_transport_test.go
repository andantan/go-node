package network

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	// Interconnect tra and trb
	tra.Connect(trb)
	trb.Connect(tra)

	// Test mapping address -> This will be IP? WS?
	assert.Equal(t, tra.peers[trb.Addr()], trb)
	assert.Equal(t, trb.peers[tra.Addr()], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A").(*LocalTransport)
	trb := NewLocalTransport("B").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("Hello world!")

	// Sending message msg to trb return not error(nil)
	assert.Nil(t, tra.SendMessage(trb.addr, msg))

	// rpc.(type) is <- chan RPC
	rpc := <-trb.Consume()

	b, err := io.ReadAll(rpc.Payload)

	assert.Nil(t, err)

	// rpc.payload will be msg
	assert.Equal(t, msg, b)
	assert.Equal(t, rpc.From, tra.addr)
}

func TestBroadast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("foo")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, b)

	rpcc := <-trc.Consume()
	c, err := io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, c)
}
