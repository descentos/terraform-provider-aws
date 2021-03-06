package analysisutils

import (
	"go/ast"
	"go/types"

	"github.com/bflad/tfproviderlint/passes/commentignore"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// DeprecatedWithReplacementSelectorExprRunner returns an Analyzer runner for deprecated *ast.SelectorExpr with replacement
func DeprecatedWithReplacementSelectorExprRunner(analyzerName string, selectorExprAnalyzer *analysis.Analyzer, oldPackageName, oldSelectorName, newPackageName, newSelectorName string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		selectorExprs := pass.ResultOf[selectorExprAnalyzer].([]*ast.SelectorExpr)
		ignorer := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)

		for _, selectorExpr := range selectorExprs {
			if ignorer.ShouldIgnore(analyzerName, selectorExpr) {
				continue
			}

			pass.Reportf(selectorExpr.Pos(), "%s: deprecated %s.%s should be replaced with %s.%s", analyzerName, oldPackageName, oldSelectorName, newPackageName, newSelectorName)
		}

		return nil, nil
	}
}

// FunctionCallExprRunner returns an Analyzer runner for function *ast.CallExpr
func FunctionCallExprRunner(packageFunc func(ast.Expr, *types.Info, string) bool, functionName string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		nodeFilter := []ast.Node{
			(*ast.CallExpr)(nil),
		}
		var result []*ast.CallExpr

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			callExpr := n.(*ast.CallExpr)

			if !packageFunc(callExpr.Fun, pass.TypesInfo, functionName) {
				return
			}

			result = append(result, callExpr)
		})

		return result, nil
	}
}

// ReceiverMethodCallExprRunner returns an Analyzer runner for receiver method *ast.CallExpr
func ReceiverMethodCallExprRunner(packageReceiverMethodFunc func(ast.Expr, *types.Info, string, string) bool, receiverName string, methodName string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		nodeFilter := []ast.Node{
			(*ast.CallExpr)(nil),
		}
		var result []*ast.CallExpr

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			callExpr := n.(*ast.CallExpr)

			if !packageReceiverMethodFunc(callExpr.Fun, pass.TypesInfo, receiverName, methodName) {
				return
			}

			result = append(result, callExpr)
		})

		return result, nil
	}
}

// SelectorExprRunner returns an Analyzer runner for *ast.SelectorExpr
func SelectorExprRunner(packageFunc func(ast.Expr, *types.Info, string) bool, selectorName string) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		nodeFilter := []ast.Node{
			(*ast.SelectorExpr)(nil),
		}
		var result []*ast.SelectorExpr

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			selectorExpr := n.(*ast.SelectorExpr)

			if !packageFunc(selectorExpr, pass.TypesInfo, selectorName) {
				return
			}

			result = append(result, selectorExpr)
		})

		return result, nil
	}
}
