package network

import (
	"ProjectX/core"
	"ProjectX/types"
	"sort"
)

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool{
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

type TxMapSorter struct{
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter{
	
	txx := make([]*core.Transaction, len(txMap))

	i := 0
	for _, tx := range txMap{
		txx[i] = tx
		i++
	}

	s := &TxMapSorter{
		transactions : txx,
	}
	sort.Sort(s);
	
	return s
}

func (s *TxMapSorter) Less(i,j int) bool { 
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen();
}

func (s *TxMapSorter) Len() int { 
	return len(s.transactions) 
}

func (s *TxMapSorter) Swap(i, j int)  { 
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

func (tp *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(tp.transactions)
	return s.transactions
}

//Add adds transaction to the pool, caller is responsible to check transaction
// whether that transaction already exists or not.
func (tp *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	tp.transactions[hash] = tx
	return nil
}


func (tp *TxPool) Has(hash types.Hash) bool{
	_, ok := tp.transactions[hash]
	return ok
}

func (tp *TxPool) Len() int{
	return len(tp.transactions)
}

func (tp *TxPool) Flush() {
	tp.transactions = make(map[types.Hash]*core.Transaction)
}