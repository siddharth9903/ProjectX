package crypto

import (
	"ProjectX/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func GeneratePrivateKey() PrivateKey{
	key, err := ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
	if err != nil{
		panic(err)
	}
	return PrivateKey{
		key: key,
	}
}

func (p PrivateKey) PublicKey() PublicKey{
	return PublicKey{
		Key: &p.key.PublicKey,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) ToSlice() []byte{
	return elliptic.MarshalCompressed(k.Key,k.Key.X,k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h:= sha256.Sum256(k.ToSlice())
	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	R,S *big.Int
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r,s, err := ecdsa.Sign(rand.Reader, k.key,data)

	if err != nil {
		return nil, err
	}

	return &Signature{
		R:r,
		S:s,
	}, nil
}

func (sig Signature) Verify(pubKey PublicKey,data []byte) bool{
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}