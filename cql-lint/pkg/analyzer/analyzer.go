package analyzer

import (
	"go/ast"
	"go/token"

	"github.com/elliotchance/pie/v2"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TODO ver version para no tener problems entre el linter y el cql

var Analyzer = &analysis.Analyzer{
	Name:     "cql",
	Doc:      "Checks that cql queries will not generate run-time errors.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var cqlMethods = []string{"Query", "Update", "Delete"}

type Position struct {
	Number token.Pos
	Model  string
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	positionsToReport := []Position{}

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

		positionsToReport, _ = findErrorIsDynamic(positionsToReport, []string{}, callExpr.Args[1:]) // first parameters is ignored as it's the db object
	})

	for _, position := range positionsToReport {
		pass.Reportf(
			position.Number,
			"%s is not joined by the query",
			position.Model,
		)
	}

	return nil, nil //nolint:nilnil // is necessary
}

func findErrorIsDynamic(positionsToReport []Position, models []string, conditions []ast.Expr) ([]Position, []string) {
	for _, condition := range conditions {
		// TODO puede haber problems con el orden de las condiciones
		conditionCall := condition.(*ast.CallExpr)
		conditionSelector := conditionCall.Fun.(*ast.SelectorExpr)

		joinCondition, isJoinCondition := conditionSelector.X.(*ast.SelectorExpr)
		if isJoinCondition {
			models = addUnique(models, joinCondition.Sel.Name)
			// TODO que pasa si el join adentro no tiene wheres -> intentar con el nombre de la relacion y sino bueno falso positivo
			// func testJoinedWithJoinedWithoutCondition() {
			// 	cql.Query[models.Phone](
			// 		db,
			// 		conditions.Phone.Brand(),
			// 		conditions.Phone.Name.IsDynamic().Eq(conditions.Brand.Name.Value()),
			// 	).Find()
			// }
			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args)
		} else {
			positionsToReport, models = findErrorIsDynamicWhereCondition(positionsToReport, models, conditionCall, conditionSelector)
		}
	}

	return positionsToReport, models
}

func findErrorIsDynamicWhereCondition(
	positionsToReport []Position, models []string,
	conditionCall *ast.CallExpr,
	conditionSelector *ast.SelectorExpr,
) ([]Position, []string) {
	whereCondition, isWhereCondition := conditionSelector.X.(*ast.CallExpr)
	if isWhereCondition {
		_, modelName, _ := getModel(whereCondition)
		models = addUnique(models, modelName)

		if getFieldIsMethodName(whereCondition) == "IsDynamic" {
			modelPkg, isDynamicModelName, modelPos := getModel(conditionCall.Args[0].(*ast.CallExpr))
			if !pie.Contains(models, isDynamicModelName) {
				positionsToReport = append(positionsToReport, Position{
					Number: modelPos,
					Model:  modelPkg + "." + isDynamicModelName,
				})
			}
		}
	}

	return positionsToReport, models
}

func addUnique(list []string, elem string) []string {
	if !pie.Contains(list, elem) {
		return append(list, elem)
	}

	return list
}

func getFieldIsMethodName(whereCondition *ast.CallExpr) string {
	return whereCondition.Fun.(*ast.SelectorExpr).Sel.Name
}

// Returns model's package the model name
func getModel(whereCondition *ast.CallExpr) (string, string, token.Pos) {
	model := whereCondition.Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr).X.(*ast.SelectorExpr)

	return model.X.(*ast.Ident).Name, model.Sel.Name, model.Pos()
}
