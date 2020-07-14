package main

import (
	"go-crypt-generator/generator"
	"go-crypt-generator/generator/block"
	"go-crypt-generator/generator/provider"
	"go/ast"
	"go/printer"
	"go/token"
	"math/rand"
	"os"
)

func main() {
	r := rand.New(rand.NewSource(422))

	a, _ := generator.NewAlgorithm(&generator.AlgorithmSpec{
		MinKeyCount:      3,
		MaxKeyCount:      10,
		FunctionName:     "decryptFunc",
		Complexity:       1000,
		AllowedBlockType: block.AnyType,
	}, provider.NewDefaultNameProvider(r), r)

	f := &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "main",
		},
		Decls: []ast.Decl{
			a.DecryptFunc(),
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "main",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.Ident{
									Name: "println",
								},
								Args: []ast.Expr{
									a.MakeCall("hello world"),
								},
							},
						},
					},
				},
			},
		},
	}

	_ = printer.Fprint(os.Stdout, token.NewFileSet(), f)
}
