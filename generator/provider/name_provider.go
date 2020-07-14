package provider

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync/atomic"
)

type NameProvider interface {
	GenerateName() string
}

type defaultNameProvider struct {
	prefix  string
	counter uint64
}

func (r *defaultNameProvider) GenerateName() string {
	counter := atomic.AddUint64(&r.counter, 1)
	return fmt.Sprintf("n_%s_%d", r.prefix, counter)
}

func NewDefaultNameProvider(rand *rand.Rand) NameProvider {
	rawPrefix := make([]byte, 8)
	_, _ = rand.Read(rawPrefix)

	return &defaultNameProvider{
		prefix:  hex.EncodeToString(rawPrefix),
		counter: 0,
	}
}
