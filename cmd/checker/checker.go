package main

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/gofix"
	"golang.org/x/tools/go/analysis/passes/hostport"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
	"honnef.co/go/tools/quickfix/qf1001"
	"honnef.co/go/tools/simple/s1000"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck/st1000"
)

func addPassesAnalizers(analizersList []*analysis.Analyzer) []*analysis.Analyzer {
	analizersList = append(analizersList, appends.Analyzer)
	analizersList = append(analizersList, asmdecl.Analyzer)
	analizersList = append(analizersList, assign.Analyzer)
	analizersList = append(analizersList, atomic.Analyzer)
	analizersList = append(analizersList, atomicalign.Analyzer)
	analizersList = append(analizersList, bools.Analyzer)
	analizersList = append(analizersList, buildssa.Analyzer)
	analizersList = append(analizersList, buildtag.Analyzer)
	analizersList = append(analizersList, cgocall.Analyzer)
	analizersList = append(analizersList, composite.Analyzer)
	analizersList = append(analizersList, copylock.Analyzer)
	analizersList = append(analizersList, ctrlflow.Analyzer)
	analizersList = append(analizersList, deepequalerrors.Analyzer)
	analizersList = append(analizersList, defers.Analyzer)
	analizersList = append(analizersList, directive.Analyzer)
	analizersList = append(analizersList, errorsas.Analyzer)
	analizersList = append(analizersList, fieldalignment.Analyzer)
	analizersList = append(analizersList, findcall.Analyzer)
	analizersList = append(analizersList, framepointer.Analyzer)
	analizersList = append(analizersList, gofix.Analyzer)
	analizersList = append(analizersList, hostport.Analyzer)
	analizersList = append(analizersList, httpmux.Analyzer)
	analizersList = append(analizersList, httpresponse.Analyzer)
	analizersList = append(analizersList, ifaceassert.Analyzer)
	analizersList = append(analizersList, inspect.Analyzer)
	analizersList = append(analizersList, loopclosure.Analyzer)
	analizersList = append(analizersList, lostcancel.Analyzer)
	analizersList = append(analizersList, nilfunc.Analyzer)
	analizersList = append(analizersList, nilness.Analyzer)
	analizersList = append(analizersList, pkgfact.Analyzer)
	analizersList = append(analizersList, printf.Analyzer)
	analizersList = append(analizersList, reflectvaluecompare.Analyzer)
	analizersList = append(analizersList, shadow.Analyzer)
	analizersList = append(analizersList, shift.Analyzer)
	analizersList = append(analizersList, sigchanyzer.Analyzer)
	analizersList = append(analizersList, slog.Analyzer)
	analizersList = append(analizersList, sortslice.Analyzer)
	analizersList = append(analizersList, stdmethods.Analyzer)
	analizersList = append(analizersList, stdversion.Analyzer)
	analizersList = append(analizersList, stringintconv.Analyzer)
	analizersList = append(analizersList, structtag.Analyzer)
	analizersList = append(analizersList, testinggoroutine.Analyzer)
	analizersList = append(analizersList, tests.Analyzer)
	analizersList = append(analizersList, timeformat.Analyzer)
	analizersList = append(analizersList, unmarshal.Analyzer)
	analizersList = append(analizersList, unreachable.Analyzer)
	analizersList = append(analizersList, unsafeptr.Analyzer)
	analizersList = append(analizersList, unusedresult.Analyzer)
	analizersList = append(analizersList, unusedwrite.Analyzer)
	analizersList = append(analizersList, usesgenerics.Analyzer)
	analizersList = append(analizersList, waitgroup.Analyzer)
	return analizersList
}

func addOtherSCAnalizers(analizersList []*analysis.Analyzer) []*analysis.Analyzer {
	analizersList = append(analizersList, s1000.Analyzer)
	analizersList = append(analizersList, st1000.Analyzer)
	analizersList = append(analizersList, qf1001.Analyzer)
	return analizersList
}

func main() {
	fmt.Println("start checker")
	var analizersList []*analysis.Analyzer
	analizersList = addPassesAnalizers(analizersList)
	for _, v := range staticcheck.Analyzers {
		if strings.Contains(v.Analyzer.Name, "SA") {
			analizersList = append(analizersList, v.Analyzer)
		}
	}

	multichecker.Main(
		analizersList...,
	)
}
