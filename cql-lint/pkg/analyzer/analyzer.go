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

var (
	cqlMethods   = []string{"Query", "Update", "Delete"}
	cqlOrder     = []string{"Descending", "Ascending"}
	cqlSet       = []string{"SetMultiple"}
	cqlSelectors = append(cqlOrder, cqlSet...)
)

type Position struct {
	Number token.Pos
	Model  string
}

type Model struct {
	pkg  string
	name string
	pos  token.Pos
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

		// methods to be verified have at least one parameter (cql.Query, cql.Update, cql.Delete, Descending, Ascending)
		if len(callExpr.Args) < 1 {
			return
		}

		_, isSelector := callExpr.Fun.(*ast.SelectorExpr)
		if isSelector {
			positionsToReport = findForSelector(callExpr, positionsToReport)
		} else {
			positionsToReport, _ = findForIndex(callExpr, positionsToReport)
		}
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

// TODO set dynamic

// Finds errors in selector functions: Descending, Ascending, SetMultiple
func findForSelector(callExpr *ast.CallExpr, positionsToReport []Position) []Position {
	selectorExpr := callExpr.Fun.(*ast.SelectorExpr)

	if !pie.Contains(cqlSelectors, selectorExpr.Sel.Name) {
		return positionsToReport
	}

	_, models := findForIndex(selectorExpr.X.(*ast.CallExpr), positionsToReport)

	for _, arg := range callExpr.Args {
		var model Model

		if pie.Contains(cqlOrder, selectorExpr.Sel.Name) {
			model = getModel(arg.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
		} else {
			// TODO tambien podria ser ser set dynamic que es distitno, hay que unificar para que ambos sean un solo llamado
			model = getModel(arg.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
		}

		positionsToReport = addPositionsToReport(positionsToReport, models, model)
	}

	return positionsToReport
}

// Finds errors in index functions: cql.Query, cql.Update, cql.Delete
func findForIndex(callExpr *ast.CallExpr, positionsToReport []Position) ([]Position, []string) {
	indexExpr, isIndex := callExpr.Fun.(*ast.IndexExpr)
	if !isIndex {
		// other functions may be between callExpr and the cql method, example: cql.Query(...).Limit(1).Descending
		selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
		if isSelector {
			callExpr, isCall := selectorExpr.X.(*ast.CallExpr)
			if isCall {
				return findForIndex(callExpr, positionsToReport)
			}
		}

		return positionsToReport, []string{}
	}

	selectorExpr, isSelector := indexExpr.X.(*ast.SelectorExpr)
	if !isSelector {
		return positionsToReport, []string{}
	}

	ident, isIdent := selectorExpr.X.(*ast.Ident)
	if !isIdent {
		return positionsToReport, []string{}
	}

	if ident.Name != "cql" || !pie.Contains(cqlMethods, selectorExpr.Sel.Name) {
		return positionsToReport, []string{}
	}

	var models []string

	positionsToReport, models = findErrorIsDynamic(positionsToReport, []string{}, callExpr.Args[1:]) // first parameters is ignored as it's the db object

	if len(models) == 0 {
		// add method generic in case there where not conditions
		models = append(models, indexExpr.Index.(*ast.SelectorExpr).Sel.Name)
	}

	return positionsToReport, models
}

func findErrorIsDynamic(positionsToReport []Position, models []string, conditions []ast.Expr) ([]Position, []string) {
	for _, condition := range conditions {
		conditionCall := condition.(*ast.CallExpr)
		conditionSelector := conditionCall.Fun.(*ast.SelectorExpr)

		if conditionSelector.Sel.Name == "Preload" {
			conditionCall = conditionSelector.X.(*ast.CallExpr)
			conditionSelector = conditionCall.Fun.(*ast.SelectorExpr)
		}

		joinCondition, isJoinCondition := conditionSelector.X.(*ast.SelectorExpr)

		if isJoinCondition {
			models = addUnique(models, joinCondition.Sel.Name) // conditions.Phone.Brand -> Phone

			oldLen := len(models)

			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args)

			if len(models) == oldLen {
				// only add the joined model if no model was added inside the join
				// this is because maybe the joined model is called different that the relation
				// so we prioritize conditions.Brand.Name over conditions.Phone.Brand
				models = addUnique(models, conditionSelector.Sel.Name) // conditions.Phone.Brand -> Brand
			}
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
		model := getModelFromWhereCondition(whereCondition)

		models = addUnique(models, model.name)

		if getFieldIsMethodName(whereCondition) == "IsDynamic" {
			isDynamicModel := getModelFromWhereCondition(conditionCall.Args[0].(*ast.CallExpr))
			positionsToReport = addPositionsToReport(positionsToReport, models, isDynamicModel)
		}
	}

	return positionsToReport, models
}

func addPositionsToReport(positionsToReport []Position, models []string, model Model) []Position {
	if !pie.Contains(models, model.name) {
		return append(positionsToReport, Position{
			Number: model.pos,
			Model:  model.pkg + "." + model.name,
		})
	}

	return positionsToReport
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
func getModelFromWhereCondition(whereCondition *ast.CallExpr) Model {
	return getModel(whereCondition.Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
}

// Returns model's package the model name
func getModel(selExpr *ast.SelectorExpr) Model {
	return Model{
		pkg:  selExpr.X.(*ast.Ident).Name,
		name: selExpr.Sel.Name,
		pos:  selExpr.Pos(),
	}
}
