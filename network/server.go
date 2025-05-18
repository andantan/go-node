package network

import (
	"fmt"
	"time"

	"github.com/andantan/go-node/core"
	"github.com/andantan/go-node/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	RPCHandler RPCHandler
	// This is will be container
	Transports []Transport
	Blocktime  time.Duration
	PrivateKey *crypto.PrivateKey // If has PrivateKey, node is validator
}

type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.Blocktime == time.Duration(0) {
		opts.Blocktime = defaultBlockTime
	}

	s := &Server{
		blockTime:   opts.Blocktime,
		memPool:     newTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}),
	}

	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}

	s.ServerOpts = opts

	return s
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			// Is there some message from rpc channel
			// fmt.Printf("%+v\n", rpc)

			// Somebody send wrong byte or malformed payload
			// then just logging
			if err := s.RPCHandler.HandlerRPC(rpc); err != nil {
				logrus.Error(err)
			}
		case <-s.quitCh:
			// break -> Break select statement NOT for loop
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")

		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash,
		"mempool length": s.memPool.len(),
	}).Info("adding new tx to mempool")

	// TODO(@andantan): broadcast this tx to peers

	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// Handle
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
