package block

import (
	"go-crypt-generator/generator/provider"
	"go/ast"
	"go/token"
	"math/rand"
)

type loopBlock struct{}

func (l *loopBlock) Generate(dataFieldName string, keyFieldNames []string, nameProvider provider.NameProvider, rand *rand.Rand) (EncryptFunc, *ast.BlockStmt) {
	keyIdx := rand.Intn(len(keyFieldNames))
	indexFieldName := nameProvider.GenerateName()
	keyName := keyFieldNames[keyIdx]

	return l.createEncryptFunc(keyIdx), l.createDecryptBlock(indexFieldName, dataFieldName, keyName)
}

func (l *loopBlock) createDecryptBlock(indexFieldName string, dataFieldName string, keyName string) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.RangeStmt{
				Key: &ast.Ident{
					Name: indexFieldName,
				},
				Tok: token.DEFINE,
				X: &ast.Ident{
					Name: dataFieldName,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.IndexExpr{
									X: &ast.Ident{
										Name: dataFieldName,
									},
									Index: &ast.Ident{
										Name: indexFieldName,
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
												Name: indexFieldName,
											},
											Op: token.XOR,
											Y: &ast.Ident{
												Name: keyName,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (*loopBlock) createEncryptFunc(keyIdx int) func(bytes []byte, keys []int) []byte {
	return func(bytes []byte, keys []int) []byte {
		for i := range bytes {
			bytes[i] ^= byte(i ^ keys[keyIdx])
		}
		return bytes
	}
}

func (*loopBlock) Type() Type {
	return LoopType
}

func (*loopBlock) Complexity() int {
	return 5
}
