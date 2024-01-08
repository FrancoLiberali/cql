package analyzer

import (
	"go/ast"

	"github.com/elliotchance/pie/v2"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

// TODO ver version para no tener problemas entre el linter y el cql
var Analyzer = &analysis.Analyzer{
	Name:     "cql",
	Doc:      "Checks that cql queries will not generate run-time errors.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var cqlMethods = []string{"Query", "Update", "Delete"}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()

		callExpr := node.(*ast.CallExpr)

		// cql.Query, cql.Update and cql.Delete have at least db as parameter
		if len(callExpr.Args) < 1 {
			return
		}

		selectorExpr, isSelector := callExpr.Fun.(*ast.IndexExpr).X.(*ast.SelectorExpr)
		if !isSelector {
			return
		}

		ident, isIdent := selectorExpr.X.(*ast.Ident)
		if !isIdent {
			return
		}

		if ident.Name != "cql" {
			return
		}

		if !pie.Contains(cqlMethods, selectorExpr.Sel.Name) {
			return
		}

		findErrorIsDynamic(pass, node, []string{}, callExpr.Args[1:]) // first parameters is ignored as it's the db object

		// pass.Reportf(node.Pos(), "printf-like formatting function '%s' should be named '%sf'",
		// 	funcDecl.Name.Name, funcDecl.Name.Name)
	})

	return nil, nil
}

func findErrorIsDynamic(pass *analysis.Pass, node ast.Node, models []string, conditions []ast.Expr) []string {
	addedToModels := false

	for _, condition := range conditions {
		// TODO puede haber problemas con el orden de las condiciones
		conditionCall := condition.(*ast.CallExpr)
		conditionSelector := conditionCall.Fun.(*ast.SelectorExpr)

		// TODO tener en cuenta el preload
		joinCondition, isJoinCondition := conditionSelector.X.(*ast.SelectorExpr)
		if isJoinCondition {
			addedToModels = true
			models = append(models, joinCondition.Sel.Name)
			// TODO que pasa si el join adentro no tiene wheres -> intentar con el nombre de la relacion y sino bueno falso positivo
			models = findErrorIsDynamic(pass, node, models, conditionCall.Args)
		} else {
			whereCondition, isWhereCondition := conditionSelector.X.(*ast.CallExpr)
			if isWhereCondition {
				if !addedToModels {
					addedToModels = true
					_, _, modelName := getFieldAndModel(whereCondition)
					models = append(models, modelName)
				}

				if getFieldIsMethodName(whereCondition) == "IsDynamic" {
					_, modelPkg, isDynamicModelName := getFieldAndModel(conditionCall.Args[0].(*ast.CallExpr))
					if !pie.Contains(models, isDynamicModelName) {
						pass.Reportf(
							node.Pos(),
							"%s is not joined by the query",
							modelPkg+"."+isDynamicModelName,
						)
					}
				}
			}
		}
	}

	return models
}

func getFieldIsMethodName(whereCondition *ast.CallExpr) string {
	return whereCondition.Fun.(*ast.SelectorExpr).Sel.Name
}

// Returns the field name, the model package the model name
func getFieldAndModel(whereCondition *ast.CallExpr) (string, string, string) {
	field := whereCondition.Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr)
	model := field.X.(*ast.SelectorExpr)

	return field.Sel.Name, model.X.(*ast.Ident).Name, model.Sel.Name
}
