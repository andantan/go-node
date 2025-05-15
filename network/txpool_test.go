package network

import (
	"math/rand"
	"strconv"
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

func TestSortTransactions(t *testing.T) {
	p := newTxPool()
	txLen := 1000

	for i := range txLen {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.len())

	txx := p.Transactions()

	// Sort test
	for i := range len(txx) - 1 {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
}
