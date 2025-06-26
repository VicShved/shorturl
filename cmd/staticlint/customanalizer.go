package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitCheckAnalyzer Анализатор на os.Exit
var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check os.Exit in func main package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var packageName string
	var funcName string

	analyzerFunc := func(node ast.Node) bool {

		switch x := node.(type) {
		case *ast.File:
			{
				packageName = x.Name.Name
				fmt.Println("packageName = " + packageName)
			}
		case *ast.FuncDecl:
			{
				funcName = x.Name.Name
				fmt.Println("funcName = " + funcName)
			}
		case *ast.SelectorExpr:
			{
				if (x.X.(*ast.Ident).Name == "os") && (x.Sel.Name == "Exit") && (funcName == "main") && (packageName == "main") {
					pass.Reportf(x.Pos(), "os.Exit call error")
				}

			}
		}
		return true
	}

	for _, file := range pass.Files {
		ast.Inspect(file, analyzerFunc)
	}
	return nil, nil
}
