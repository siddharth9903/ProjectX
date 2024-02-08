package types

import (
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)
	for i:=0;i<32;i++ {
		b[i] = h[i]
	}
	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func (h Hash) IsZero() bool {
	for i := 0; i < len(h); i++{
		if h[i] != 0 {
			return false
		}
	}
	return true
}


func HashFromBytes(b []byte) Hash{
	if len(b)!=32 {
		msg := fmt.Sprintf("given bytes with length %d is should be 32",len(b))
		panic(msg)
	}
	
	var value [32]uint8
	for i:=0 ; i<len(b) ; i++ {
		value[i] = b[i]
	}
	return value
}