// Модуль кастомного анализатора

package customchecker

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// OsExit проверка на наличие os.Exit в main функции
func OsExit(file ast.Node) *token.Pos {
	var pos *token.Pos
	ast.Inspect(file, func(n1 ast.Node) bool {
		pkg, okFile := n1.(*ast.File)
		if !okFile || pkg.Name.Name != "main" {
			return true
		}

		ast.Inspect(pkg, func(n2 ast.Node) bool {
			a, okFuncDecl := n2.(*ast.FuncDecl)
			if !okFuncDecl || a.Name.Name != "main" {
				return true
			}

			ast.Inspect(a.Body, func(bodyNode ast.Node) bool {
				callExpr, ok := bodyNode.(*ast.CallExpr)
				if !ok {
					return true
				}

				selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := selectorExpr.X.(*ast.Ident)
				if !ok || ident.Name != "os" || selectorExpr.Sel.Name != "Exit" {
					return true
				}

				pos = &ident.NamePos
				return false
			})
			return false
		})
		return false
	})

	return pos
}

var Analyzer = &analysis.Analyzer{
	Name: "osexitchecker",
	Doc:  "check is exist in package 'main'",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		foundPos := OsExit(file)
		if foundPos != nil {
			pos := *foundPos
			pass.Reportf(pos, "detected os.Exit()")
		}
	}
	return nil, nil
}
