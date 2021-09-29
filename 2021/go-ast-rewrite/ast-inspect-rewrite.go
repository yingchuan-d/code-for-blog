package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", os.Stdin, 0)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			id, ok := x.Fun.(*ast.Ident)
			if ok {
				if id.Name == "pred" {
					id.Name += "2"
				}
			}
		case *ast.FuncDecl:
			if x.Name.Name == "pred" {
				x.Name.Name += "2"
			}

			body := x.Body
			newCallStmt := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "fmt",
						},
						Sel: &ast.Ident{
							Name: "Println",
						},
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"instrumentation"`,
						},
					},
				},
			}

			body.List = append([]ast.Stmt{newCallStmt}, body.List...)
		}

		return true
	})

	fmt.Println("Modified AST:")
	printer.Fprint(os.Stdout, fset, file)
}