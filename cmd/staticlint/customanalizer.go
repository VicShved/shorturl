package main

import "golang.org/x/tools/go/analysis"

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check os.Exit in func main package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	return nil, nil
}
