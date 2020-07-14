package generator

import (
	"go-crypt-generator/generator/block"
	"go-crypt-generator/generator/provider"
	"go/ast"
	"math/rand"
)

type AlgorithmSpec struct {
	MinKeyCount, MaxKeyCount int

	FunctionName string
	Complexity   int

	AllowedBlockType block.Type
}

func generateFields(count int, nameProvider provider.NameProvider) (*ast.Field, string, []*ast.Field, []string) {
	fields := make([]*ast.Field, count)
	fieldNames := make([]string, count)

	dataName := nameProvider.GenerateName()
	dataField := &ast.Field{
		Names: []*ast.Ident{
			{
				Name: dataName,
			},
		},
		Type: &ast.ArrayType{
			Elt: &ast.Ident{
				Name: "byte",
			},
		},
	}

	for i := range fields {
		name := nameProvider.GenerateName()
		fieldNames[i] = name
		fields[i] = &ast.Field{
			Names: []*ast.Ident{
				{
					Name: name,
				},
			},
			Type: &ast.Ident{
				Name: "int",
			},
		}
	}

	return dataField, dataName, fields, fieldNames
}

func newDecryptFunc(name string, dataField *ast.Field, dataFieldIndex int, keys []*ast.Field) *ast.FuncDecl {
	params := make([]*ast.Field, len(keys)+1)
	copy(params[:dataFieldIndex], keys[:dataFieldIndex])
	params[dataFieldIndex] = dataField
	copy(params[dataFieldIndex+1:], keys[dataFieldIndex:])

	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: name,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{List: params},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								Name: "byte",
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{},
	}
}

func NewAlgorithm(spec *AlgorithmSpec, nameProvider provider.NameProvider, rand *rand.Rand) (*Algorithm, error) {
	keyCount := spec.MinKeyCount + rand.Intn(spec.MaxKeyCount-spec.MinKeyCount)
	dataFieldIndex := rand.Intn(keyCount)
	dataField, dataName, fields, fieldNames := generateFields(keyCount, nameProvider)

	decryptFunc := newDecryptFunc(spec.FunctionName, dataField, dataFieldIndex, fields)

	remainComplexity := spec.Complexity

	a := &Algorithm{
		encryptFunctions: make([]block.EncryptFunc, 0),
		decryptFunc:      decryptFunc,
		rand:             rand,
		keyCount:         keyCount,
		dataFieldIndex:   dataFieldIndex,
	}

	for {
		if remainComplexity <= 0 {
			break
		}

		sourceBlocks := block.GetBlocks(spec.AllowedBlockType, remainComplexity)

		b := sourceBlocks[rand.Intn(len(sourceBlocks))]
		encryptFunc, decryptBlock := b.Generate(dataName, fieldNames, nameProvider, rand)

		a.encryptFunctions = append(a.encryptFunctions, encryptFunc)
		if len(decryptBlock.List) > 1 {
			decryptFunc.Body.List = append(decryptFunc.Body.List, decryptBlock)
		} else {
			decryptFunc.Body.List = append(decryptFunc.Body.List, decryptBlock.List[0])
		}

		remainComplexity -= b.Complexity()
	}
	decryptFunc.Body.List = append(decryptFunc.Body.List, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.Ident{
				Name: dataName,
			},
		},
	})
	return a, nil
}
