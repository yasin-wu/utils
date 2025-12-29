package blowfish

import (
	"encoding/hex"
	"errors"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

type Blowfish struct {
	initialVector string
	secKey        string
}

func New(initialVector, secKey string) (*Blowfish, error) {
	if len([]byte(initialVector)) != 16 {
		return nil, errors.New("initialVector must be 16 bytes")
	}
	if len([]byte(secKey)) != 16 {
		return nil, errors.New("secKey must be 16 bytes")
	}
	return &Blowfish{
		initialVector: initialVector,
		secKey:        secKey,
	}, nil
}

func (b *Blowfish) Encrypt(data []byte) (string, error) {
	cryptoCli := crypto.New().FromBytes(data).SetKey(b.secKey).SetIv(b.initialVector).
		Aes().Blowfish().PKCS5Padding().Encrypt()
	return cryptoCli.ToHexString(), cryptoCli.Error()
}

func (b *Blowfish) Decrypt(data string) ([]byte, error) {
	dataHex, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	cryptoCli := crypto.New().FromBytes(dataHex).SetKey(b.secKey).SetIv(b.initialVector).
		Aes().Blowfish().PKCS5Padding().Decrypt()
	return cryptoCli.ToBytes(), cryptoCli.Error()
}
