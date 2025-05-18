package network

import (
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

	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)

	assert.Equal(t, n, len(msg))
	assert.Nil(t, err)

	// rpc.payload will be msg
	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, tra.addr)
}
