package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/andantan/go-node/core"
	"github.com/andantan/go-node/crypto"
	"github.com/andantan/go-node/network"
	"github.com/sirupsen/logrus"
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
	trRemote := network.NewLocalTransport("REMOTE") // 34.4.244.33

	// Local chaining -> This will be node chaining
	// which validate TX, PBFT and commit block
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			// trRemote.SendMessage(trLocal.Addr(), []byte("Hello world"))
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}

			// time.Sleep(1 * time.Second)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	privKey := crypto.GeneratePrivateKey()

	opts := network.ServerOpts{
		PrivateKey: &privKey,
		ID:         "LOCAL",
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}

// Just placeholder (demo)
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()

	data := []byte(strconv.FormatInt(int64(rand.Intn(100000000)), 10))
	tx := core.NewTransaction(data)

	tx.Sign(privKey)

	buf := &bytes.Buffer{}

	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
