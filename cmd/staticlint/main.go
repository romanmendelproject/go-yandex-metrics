// Static code analysis package.
//
// customchecker - custom os.Exit analyzers
//
// fieldtype - finds confliction type of field
//
// assign - Package assign defines an Analyzer that detects useless assignments.
//
// findcall - Package findcall defines an Analyzer that serves as a trivial example and test of the Analysis API.
//
// inspect - Package inspect defines an Analyzer that provides an AST inspector
//
// printf - Package printf defines an Analyzer that checks consistency of Printf format strings and arguments.
//
// shadow - Package shadow defines an Analyzer that checks for shadowed variables.
//
// shift - Package shift defines an Analyzer that checks for shifts that exceed the width of an integer.
//
// structtag - Package structtag defines an Analyzer that checks struct field tags are well formed.
//
// Usage:
//
//	go run cmd/staticlint/main.go ./...
package main

import (
	"github.com/gostaticanalysis/funcstat"
	"github.com/gostaticanalysis/zapvet/passes/fieldtype"
	"github.com/romanmendelproject/go-yandex-metrics/pkg/customchecker"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

// Register analyzers
func main() {
	var analyzers []*analysis.Analyzer

	// Staticcheck analyzers.
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}
	// Quickfix analyzers.
	for _, qf := range quickfix.Analyzers {
		if qf.Analyzer.Name == "QF1006" {
			analyzers = append(analyzers, qf.Analyzer)
		}
	}

	// Custom os.Exit analyzers.
	analyzers = append(analyzers, customchecker.Analyzer)

	// External analyzers.
	analyzers = append(analyzers, fieldtype.Analyzer)
	analyzers = append(analyzers, funcstat.Analyzer)

	// Passes analyzers.
	analyzers = append(analyzers, assign.Analyzer)
	analyzers = append(analyzers, findcall.Analyzer)
	analyzers = append(analyzers, inspect.Analyzer)
	analyzers = append(analyzers, printf.Analyzer)
	analyzers = append(analyzers, shadow.Analyzer)
	analyzers = append(analyzers, shift.Analyzer)
	analyzers = append(analyzers, structtag.Analyzer)

	// Run all analyzers.
	multichecker.Main(analyzers...)
}
