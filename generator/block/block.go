package block

import (
	"go-crypt-generator/generator/provider"
	"go/ast"
	"math/rand"
)

type Type uint64

const (
	ArithmeticType Type = 1 << iota
	LoopType
	ChannelType   // not implement
	GoroutineType // not implement

	AnyType = ArithmeticType | LoopType | ChannelType | GoroutineType
)

var allBlocks = []Block{
	&arithmeticBlock{},
	&loopBlock{},
}

type EncryptFunc func([]byte, []int) []byte

type Block interface {
	Generate(dataFieldName string, keyFieldNames []string, nameProvider provider.NameProvider, rand *rand.Rand) (EncryptFunc, *ast.BlockStmt)

	Type() Type
	Complexity() int
}

func GetBlocks(allowedTypes Type, maxComplexity int) []Block {
	blocks := make([]Block, 0)
	for _, block := range allBlocks {
		if block.Complexity() <= maxComplexity && block.Type()&allowedTypes == block.Type() {
			blocks = append(blocks, block)
		}
	}
	return blocks
}
