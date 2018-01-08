package state

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/vechain/thor/acc"
	"github.com/vechain/thor/cry"
	"github.com/vechain/thor/fortest"
	"github.com/vechain/thor/lvldb"
)

var M = fortest.Multi

func TestAccount(t *testing.T) {
	assert.True(t, Account{}.IsEmpty(), "newly constructed account should be empty")

	assert.False(t, Account{Balance: big.NewInt(1)}.IsEmpty())
	assert.False(t, Account{CodeHash: []byte{1}}.IsEmpty())
	assert.True(t, Account{StorageRoot: []byte{1}}.IsEmpty())
}

func newTrie() *trie.SecureTrie {
	kv, _ := lvldb.NewMem()
	trie, _ := trie.NewSecure(common.Hash{}, kv, 0)
	return trie
}
func TestTrie(t *testing.T) {
	trie := newTrie()

	addr := acc.BytesToAddress([]byte("account1"))
	assert.Equal(t,
		M(loadAccount(trie, addr)),
		[]interface{}{Account{Balance: new(big.Int)}, nil},
		"should load an empty account")

	acc1 := Account{
		big.NewInt(1),
		[]byte("code hash"),
		[]byte("storage root"),
	}
	saveAccount(trie, addr, acc1)
	assert.Equal(t,
		M(loadAccount(trie, addr)),
		[]interface{}{acc1, nil})

	saveAccount(trie, addr, Account{})
	assert.Equal(t,
		M(trie.TryGet(addr[:])),
		[]interface{}{[]byte(nil), nil},
		"empty account should be deleted")
}

func TestStorageTrie(t *testing.T) {
	trie := newTrie()

	key := cry.BytesToHash([]byte("key"))
	assert.Equal(t,
		M(loadStorage(trie, key)),
		[]interface{}{cry.Hash{}, nil})

	value := cry.BytesToHash([]byte("value"))
	saveStorage(trie, key, value)
	assert.Equal(t,
		M(loadStorage(trie, key)),
		[]interface{}{value, nil})

	saveStorage(trie, key, cry.Hash{})
	assert.Equal(t,
		M(trie.TryGet(key[:])),
		[]interface{}{[]byte(nil), nil},
		"empty storage value should be deleted")
}