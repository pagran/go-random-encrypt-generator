package generator

import (
	"fmt"
	"go-crypt-generator/generator/block"
	"go/ast"
	"go/token"
	"math/rand"
)

type Algorithm struct {
	decryptFunc      *ast.FuncDecl
	encryptFunctions []block.EncryptFunc

	keyCount       int
	dataFieldIndex int

	rand *rand.Rand
}

func (a *Algorithm) Encrypt(data []byte, keys []int) []byte {
	for i := len(a.encryptFunctions) - 1; i >= 0; i-- {
		data = a.encryptFunctions[i](data, keys)
	}
	return data
}

func (a *Algorithm) generateKeys() []int {
	keys := make([]int, a.keyCount)
	for i := range keys {
		keys[i] = a.rand.Int()
	}
	return keys
}

func dataToByteSlice(data []byte) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.ArrayType{
			Elt: &ast.Ident{Name: "byte"},
		},
		Args: []ast.Expr{&ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%q", data),
		}},
	}
}

func keysToIntLits(keys []int) []ast.Expr {
	keyLits := make([]ast.Expr, len(keys))
	for i := range keyLits {
		keyLits[i] = &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprint(keys[i]),
		}
	}
	return keyLits
}

func (a *Algorithm) MakeCall(text string) *ast.CallExpr {
	keys := a.generateKeys()
	encryptedData := a.Encrypt([]byte(text), keys)

	args := keysToIntLits(keys)
	encryptedDataStmt := dataToByteSlice(encryptedData)

	args = append(args, nil)
	copy(args[a.dataFieldIndex+1:], args[a.dataFieldIndex:])
	args[a.dataFieldIndex] = encryptedDataStmt

	return &ast.CallExpr{
		Fun: &ast.Ident{Name: "string"},
		Args: []ast.Expr{
			&ast.CallExpr{
				Fun:  a.decryptFunc.Name,
				Args: args,
			},
		},
	}
}

func (a *Algorithm) DecryptFunc() *ast.FuncDecl {
	return a.decryptFunc
}
