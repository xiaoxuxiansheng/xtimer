package hash

import (
	"crypto/sha1"
	"encoding/base32"
	"math/big"
	"strings"
)

type SHA1Encryptor struct {
}

func NewSHA1Encryptor() *SHA1Encryptor {
	return &SHA1Encryptor{}
}

func (s *SHA1Encryptor) Encrypt(origin string) uint64 {
	hashWorker := sha1.New()
	hashWorker.Write([]byte(origin))
	var i big.Int
	res, _ := i.SetString(strings.ToLower(base32.HexEncoding.EncodeToString(hashWorker.Sum(nil))), 32)
	return res.Uint64()
}
