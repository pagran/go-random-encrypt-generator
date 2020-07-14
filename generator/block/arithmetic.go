package block

import (
	"fmt"
	"go-crypt-generator/generator/provider"
	"go/ast"
	"go/token"
	"math/rand"
)

type arithmeticBlock struct{}

var supportedOp = []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.AND, token.OR, token.XOR}

func calcKey(operatorIdx, key, randKey int) byte {
	switch supportedOp[operatorIdx] {
	case token.ADD:
		return byte(key + randKey)
	case token.SUB:
		return byte(key - randKey)
	case token.MUL:
		return byte(key * randKey)
	case token.QUO:
		return byte(key / randKey)
	case token.REM:
		return byte(key % randKey)
	case token.AND:
		return byte(key & randKey)
	case token.OR:
		return byte(key | randKey)
	case token.XOR:
		return byte(key ^ randKey)
	default:
		panic("unknown operator: " + supportedOp[operatorIdx].String())
	}
}

func (a *arithmeticBlock) Generate(dataFieldName string, keyFieldNames []string, _ provider.NameProvider, rand *rand.Rand) (EncryptFunc, *ast.BlockStmt) {
	index := rand.Int()
	operatorIdx := rand.Intn(len(supportedOp))
	keyIdx := rand.Intn(len(keyFieldNames))
	randKey := rand.Int()

	keyName := keyFieldNames[keyIdx]

	decryptBlock := a.createDecryptBlock(dataFieldName, index, keyName, operatorIdx, randKey)
	encryptFunc := a.createEncryptFunc(index, operatorIdx, keyIdx, randKey)
	return encryptFunc, decryptBlock
}

func (*arithmeticBlock) createEncryptFunc(index int, operatorIdx int, keyIdx int, randKey int) func(bytes []byte, keys []int) []byte {
	return func(bytes []byte, keys []int) []byte {
		bytes[index%len(bytes)] ^= calcKey(operatorIdx, keys[keyIdx], randKey)
		return bytes
	}
}

func (*arithmeticBlock) createDecryptBlock(dataFieldName string, index int, keyName string, operatorIdx int, randKey int) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.IndexExpr{
						X: &ast.Ident{
							Name: dataFieldName,
						},
						Index: &ast.BinaryExpr{
							X: &ast.BasicLit{
								Kind:  token.INT,
								Value: fmt.Sprint(index),
							},
							Op: token.REM,
							Y: &ast.CallExpr{
								Fun: &ast.Ident{
									Name: "len",
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: dataFieldName,
									},
								},
							},
						},
					},
				},
				Tok: token.XOR_ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{
							Name: "byte",
						},
						Args: []ast.Expr{
							&ast.BinaryExpr{
								X: &ast.Ident{
									Name: keyName,
								},
								Op: supportedOp[operatorIdx],
								Y: &ast.BasicLit{
									Kind:  token.INT,
									Value: fmt.Sprint(randKey),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (*arithmeticBlock) Type() Type {
	return ArithmeticType
}

func (*arithmeticBlock) Complexity() int {
	return 1
}
