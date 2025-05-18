package network

type NetAddr string

type Transport interface {
	// This will be module of server
	// Access to all messages from transport layers
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
}
