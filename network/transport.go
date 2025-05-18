package network

type NetAddr string

type Transport interface {
	// Return Receive-only RPC channel
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Broadcast([]byte) error
	Addr() NetAddr
}
