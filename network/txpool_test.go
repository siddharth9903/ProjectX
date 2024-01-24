package network

import (
	"ProjectX/core"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T){
	tp := NewTxPool()
	assert.Equal(t, tp.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T){
	tp := NewTxPool()
	assert.Equal(t, tp.Len(), 0)
	
	tx := core.NewTransaction([]byte("hello new tx"))
	err := tp.Add(tx)
	assert.Nil(t,err)
	assert.Equal(t,tp.Len(),1)

	tx2 := core.NewTransaction([]byte("hello new txxx"))
	err2 := tp.Add(tx2)
	assert.Nil(t,err2)
	assert.Equal(t,tp.Len(),2)

	tp.Flush()
	assert.Equal(t,tp.Len(),0)
}

func TestTxSort(t *testing.T){
	tp := NewTxPool()

	end := 1000
	for i := 0; i < end; i++ { 
		newTx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		newTx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, tp.Add(newTx))
	}

	assert.Equal(t,tp.Len(),end)

	sortedTx := tp.Transactions()
	assert.Equal(t, len(sortedTx), end)

	for i := 0; i < len(sortedTx) - 1; i++ { 
		assert.True(t,  sortedTx[i].FirstSeen() < sortedTx[i + 1].FirstSeen())  
	}
}