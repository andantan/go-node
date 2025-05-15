package network

import (
	"testing"

	"github.com/andantan/go-node/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := newTxPool()

	assert.Equal(t, p.len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := newTxPool()
	tx := core.NewTransaction([]byte("f00"))

	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.len(), 1)

	_ = core.NewTransaction([]byte("f00"))

	assert.Equal(t, p.len(), 1)

	p.Flush()

	assert.Equal(t, p.len(), 0)
}
