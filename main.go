package main

import (
	"time"

	"github.com/andantan/go-node/network"
)

// Server - node
// Transport Layer => TCP, UDP
// Block
// TX
// Keypair

func main() {
	// Running machine which is a Local note
	// transport for own note server...?
	// Then we have peers which our remote-peers
	// which will represent servers in the network
	// which are not our machine...
	trLocal := network.NewLocalTransport("LOCAL")
	// Instead of way of commiunate TCP/UDP
	trRemote := network.NewLocalTransport("REMOTE")

	// Local chaining -> This will be node chaining
	// which validate TX, PBFT and commit block
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("Hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}
