package network

import (
	"bytes"
	"os"
	"time"

	"github.com/andantan/go-node/core"
	"github.com/andantan/go-node/crypto"
	"github.com/andantan/go-node/types"
	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport // This is will be container
	Blocktime     time.Duration
	PrivateKey    *crypto.PrivateKey // If has PrivateKey, node is validator
	BlockChain    *core.BlockChain
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	chain       *core.BlockChain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.Blocktime == time.Duration(0) {
		opts.Blocktime = defaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockChain(genesisBlock())

	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		memPool:     newTxPool(),
		chain:       chain,
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}),
	}

	// If server does not have RPCProcessor
	// then, default processor is this server.
	// Looks if we dont got any processor option,
	// we going to use the server as default. 이래도 될라나
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			// Is there some message from rpc channel
			// fmt.Printf("%+v\n", rpc)

			// Somebody send wrong byte or malformed payload
			// then just logging
			msg, err := s.RPCDecodeFunc(rpc)

			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}

		case <-s.quitCh:
			// break -> Break select statement NOT for loop
			break free
		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.Blocktime)

	s.Logger.Log(
		"msg", "Starting validator loop",
		"BlockTime", s.Blocktime,
	)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodeMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}

	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}

	return nil
}

// For private
func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.Logger.Log(
		"msg", "adding new tx to mempool",
		"hash", hash,
		"mempool-length", s.memPool.len(),
	)

	// Broadcast tx before add mempool
	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}

	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
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

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())

	if err != nil {
		return err
	}

	block, err := core.NewBlockFromPrevHeader(currentHeader, nil)

	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		TimeStamp: time.Now().UnixNano(),
		Height:    0,
	}

	// TODO(@andantan): DEMO
	b, _ := core.NewBlock(header, nil)

	return b
}
