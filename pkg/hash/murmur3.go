package hash

import (
	"github.com/spaolacci/murmur3"
)

type Murmur3Encyptor struct {
}

func NewMurmur3Encryptor() *Murmur3Encyptor {
	return &Murmur3Encyptor{}
}

func (m *Murmur3Encyptor) Encrypt(origin string) uint64 {
	return uint64(murmur3.Sum32([]byte(origin)))
}
