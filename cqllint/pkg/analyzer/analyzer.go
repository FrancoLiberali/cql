package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/elliotchance/pie/v2"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/FrancoLiberali/cql/cqllint/pkg/version"
)

var doc = "v" + version.Version + "\nChecks that cql queries will not generate run-time errors."

var Analyzer = &analysis.Analyzer{
	Name:     "cqllint",
	Doc:      doc,
	URL:      "compiledquerylenguage.readthedocs.io",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	cqlMethods     = []string{"Query", "Update", "Delete"}
	cqlOrder       = []string{"Descending", "Ascending"}
	cqlSetMultiple = "SetMultiple"
	cqlSets        = []string{cqlSetMultiple, "Set"}
	cqlSelectors   = append(cqlOrder, cqlSets...)
)

type Model struct {
	Pos  token.Pos
	Name string
}

var passG *analysis.Pass

func run(pass *analysis.Pass) (interface{}, error) {
	passG = pass
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	positionsToReport := []Model{}

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
			positionsToReport, _ = findNotConcernedForIndex(callExpr, positionsToReport)
		}
	})

	for _, position := range positionsToReport {
		pass.Reportf(
			position.Pos,
			"%s is not joined by the query",
			position.Name,
		)
	}

	return nil, nil //nolint:nilnil // is necessary
}

// Finds NotConcerned and Repeated errors in selector functions: Descending, Ascending, SetMultiple, Set
func findForSelector(callExpr *ast.CallExpr, positionsToReport []Model) []Model {
	selectorExpr := callExpr.Fun.(*ast.SelectorExpr)

	if !pie.Contains(cqlSelectors, selectorExpr.Sel.Name) {
		return positionsToReport
	}

	findRepeatedFields(callExpr, selectorExpr)

	return fieldNotConcerned(callExpr, selectorExpr, positionsToReport)
}

// Finds NotConcerned errors in selector functions: Descending, Ascending, SetMultiple, Set
func fieldNotConcerned(callExpr *ast.CallExpr, selectorExpr *ast.SelectorExpr, positionsToReport []Model) []Model {
	_, models := findNotConcernedForIndex(selectorExpr.X.(*ast.CallExpr), positionsToReport)

	for _, arg := range callExpr.Args {
		var model Model

		methodName := selectorExpr.Sel.Name

		if pie.Contains(cqlOrder, methodName) {
			model = getModel(arg.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
			positionsToReport = addPositionsToReport(positionsToReport, models, model)
		} else {
			argCallExpr := arg.(*ast.CallExpr)

			setFunction := argCallExpr.Fun.(*ast.SelectorExpr).Sel.Name

			if setFunction == "Dynamic" {
				model = getModel(argCallExpr.Args[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
				positionsToReport = addPositionsToReport(positionsToReport, models, model)
			}

			if methodName == cqlSetMultiple {
				model = getModel(argCallExpr.Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
				positionsToReport = addPositionsToReport(positionsToReport, models, model)
			}
		}
	}

	return positionsToReport
}

func findRepeatedFields(call *ast.CallExpr, selectorExpr *ast.SelectorExpr) {
	if !pie.Contains(cqlSets, selectorExpr.Sel.Name) {
		return
	}

	fields := map[string][]token.Pos{}

	for _, arg := range call.Args {
		condition := arg.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr)
		conditionModel := condition.X.(*ast.SelectorExpr)

		fieldName := conditionModel.X.(*ast.Ident).Name + "." + conditionModel.Sel.Name + "." + condition.Sel.Name
		fieldPos := condition.Sel.NamePos

		_, isPresent := fields[fieldName]
		if !isPresent {
			fields[fieldName] = []token.Pos{fieldPos}
		} else {
			fields[fieldName] = append(fields[fieldName], fieldPos)
		}
	}

	for fieldName, positions := range fields {
		if len(positions) > 1 {
			for _, pos := range positions {
				passG.Reportf(
					pos,
					"%s is repeated",
					fieldName,
				)
			}
		}
	}
}

// Finds NotConcerned errors in index functions: cql.Query, cql.Update, cql.Delete
func findNotConcernedForIndex(callExpr *ast.CallExpr, positionsToReport []Model) ([]Model, []string) {
	indexExpr, isIndex := callExpr.Fun.(*ast.IndexExpr)
	if !isIndex {
		// other functions may be between callExpr and the cql method, example: cql.Query(...).Limit(1).Descending
		selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
		if isSelector {
			internalCallExpr, isCall := selectorExpr.X.(*ast.CallExpr)
			if isCall {
				return findNotConcernedForIndex(internalCallExpr, positionsToReport)
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

	models := []string{
		getFirstGenericType(
			passG.TypesInfo.Types[indexExpr].Type.(*types.Signature).Results().At(0).Type().(*types.Pointer).Elem().(*types.Named),
		),
	}

	return findErrorIsDynamic(positionsToReport, models, callExpr.Args[1:]) // first parameters is ignored as it's the db object
}

func findErrorIsDynamic(positionsToReport []Model, models []string, conditions []ast.Expr) ([]Model, []string) {
	for _, condition := range conditions {
		conditionCall := condition.(*ast.CallExpr)
		conditionSelector := conditionCall.Fun.(*ast.SelectorExpr)

		if conditionSelector.Sel.Name == "Preload" {
			conditionCall = conditionSelector.X.(*ast.CallExpr)
			conditionSelector = conditionCall.Fun.(*ast.SelectorExpr)
		}

		_, isJoinCondition := conditionSelector.X.(*ast.SelectorExpr)

		if isJoinCondition {
			models = append(models, getModelFromJoinCondition(conditionSelector))

			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args)
		} else {
			positionsToReport = findErrorIsDynamicWhereCondition(positionsToReport, models, conditionCall, conditionSelector)
		}
	}

	return positionsToReport, models
}

// conditions.Phone.Brand -> Brand (or its correct type if not the same)
func getModelFromJoinCondition(conditionSelector *ast.SelectorExpr) string {
	return getFirstGenericType(
		passG.TypesInfo.Types[conditionSelector].Type.(*types.Signature).Params().At(0).Type().(*types.Slice).Elem().(*types.Named),
	)
}

func getFirstGenericType(parent *types.Named) string {
	return parent.TypeArgs().At(
		0,
	).(*types.Named).String()
}

func findErrorIsDynamicWhereCondition(
	positionsToReport []Model, models []string,
	conditionCall *ast.CallExpr,
	conditionSelector *ast.SelectorExpr,
) []Model {
	whereCondition, isWhereCondition := conditionSelector.X.(*ast.CallExpr)
	if isWhereCondition && getFieldIsMethodName(whereCondition) == "IsDynamic" {
		isDynamicModel := getModelFromWhereCondition(conditionCall.Args[0].(*ast.CallExpr))
		return addPositionsToReport(positionsToReport, models, isDynamicModel)
	}

	return positionsToReport
}

func addPositionsToReport(positionsToReport []Model, models []string, model Model) []Model {
	if !pie.Contains(models, model.Name) {
		return append(positionsToReport, Model{
			Pos:  model.Pos,
			Name: model.Name,
		})
	}

	return positionsToReport
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
		Name: passG.TypesInfo.Types[selExpr].Type.Underlying().(*types.Struct).Field(0).Type().(*types.Named).TypeArgs().At(0).(*types.Named).String(),
		Pos:  selExpr.Pos(),
	}
}
