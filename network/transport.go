package network

type NetAddr string

// MESSAGE TYPE
type RPC struct {
	// Message sent to transport layers
	// String or []bytes ?

	From    NetAddr
	Payload []byte
}

type Transport interface {
	// This will be module of server
	// Access to all messages from transport layers
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
}
