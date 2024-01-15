package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

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
	cqlConnectors  = []string{"And", "Or", "Not"}
	cqlSetMultiple = "SetMultiple"
	cqlSets        = []string{cqlSetMultiple, "Set"}
	cqlSelectors   = append(cqlOrder, cqlSets...)
	dynamicMethod  = "Dynamic"

	notJoinedMessage              = "%s is not joined by the query"
	appearanceNotNecessaryMessage = "Appearance call not necessary, %s appears only once"
	appearanceMoreThanOnceMessage = "%s appears more than once, select which one you want to use with Appearance"
	appearanceOutOfRangeMessage   = "selected appearance is bigger than %s's number of appearances"
)

type Model struct {
	Pos  token.Pos
	Name string
}

type Report struct {
	message string
	model   Model
}

type Appearance struct {
	selected bool
	number   int
}

var passG *analysis.Pass

func run(pass *analysis.Pass) (interface{}, error) {
	passG = pass
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	positionsToReport := []Report{}

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

	for _, report := range positionsToReport {
		pass.Reportf(
			report.model.Pos,
			report.message,
			report.model.Name,
		)
	}

	return nil, nil //nolint:nilnil // is necessary
}

// Finds NotConcerned and Repeated errors in selector functions: Descending, Ascending, SetMultiple, Set
func findForSelector(callExpr *ast.CallExpr, positionsToReport []Report) []Report {
	selectorExpr := callExpr.Fun.(*ast.SelectorExpr)

	if !pie.Contains(cqlSelectors, selectorExpr.Sel.Name) {
		return positionsToReport
	}

	findRepeatedFields(callExpr, selectorExpr)

	return fieldNotConcerned(callExpr, selectorExpr, positionsToReport)
}

// Finds NotConcerned errors in selector functions: Descending, Ascending, SetMultiple, Set
func fieldNotConcerned(callExpr *ast.CallExpr, selectorExpr *ast.SelectorExpr, positionsToReport []Report) []Report {
	_, models := findNotConcernedForIndex(selectorExpr.X.(*ast.CallExpr), positionsToReport)

	for _, arg := range callExpr.Args {
		methodName := selectorExpr.Sel.Name

		if pie.Contains(cqlOrder, methodName) {
			positionsToReport = findForOrder(arg, positionsToReport, models)
		} else {
			positionsToReport = findForSet(arg, positionsToReport, models, methodName)
		}
	}

	return positionsToReport
}

func findForSet(set ast.Expr, positionsToReport []Report, models []string, methodName string) []Report {
	setCall := set.(*ast.CallExpr)

	setFunction := setCall.Fun.(*ast.SelectorExpr).Sel.Name

	if setFunction == dynamicMethod {
		model, appearance := getModelFromCall(setCall.Args[0].(*ast.CallExpr))
		positionsToReport = addPositionsToReport(positionsToReport, models, model, appearance)
	}

	if methodName == cqlSetMultiple {
		model, appearance := getModelFromCall(setCall.Fun.(*ast.SelectorExpr).X.(*ast.CallExpr))
		positionsToReport = addPositionsToReport(positionsToReport, models, model, appearance)
	}

	return positionsToReport
}

func findForOrder(order ast.Expr, positionsToReport []Report, models []string) []Report {
	var model Model

	appearance := Appearance{selected: false}

	orderCall, isCall := order.(*ast.CallExpr)
	if isCall {
		model, appearance = getModelFromCall(orderCall)
	} else {
		model = getModel(order.(*ast.SelectorExpr).X.(*ast.SelectorExpr))
	}

	return addPositionsToReport(positionsToReport, models, model, appearance)
}

func findRepeatedFields(call *ast.CallExpr, selectorExpr *ast.SelectorExpr) {
	if !pie.Contains(cqlSets, selectorExpr.Sel.Name) {
		return
	}

	fields := map[string][]token.Pos{}

	for _, arg := range call.Args {
		argCall := arg.(*ast.CallExpr)
		argSelector := argCall.Fun.(*ast.SelectorExpr)
		condition := argSelector.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr)

		fieldName := getFieldName(condition)
		fieldPos := condition.Sel.NamePos

		_, isPresent := fields[fieldName]
		if !isPresent {
			fields[fieldName] = []token.Pos{fieldPos}
		} else {
			fields[fieldName] = append(fields[fieldName], fieldPos)
		}

		if argSelector.Sel.Name == dynamicMethod && len(argCall.Args) == 1 {
			comparedField, isSelector := argCall.Args[0].(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr)
			if !isSelector {
				continue
			}

			comparedFieldName := getFieldName(comparedField)

			if comparedFieldName == fieldName {
				passG.Reportf(
					comparedField.Sel.NamePos,
					"%s is set to itself",
					comparedFieldName,
				)
			}
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

func getFieldName(condition *ast.SelectorExpr) string {
	conditionModel := condition.X.(*ast.SelectorExpr)

	return conditionModel.X.(*ast.Ident).Name + "." + conditionModel.Sel.Name + "." + condition.Sel.Name
}

// Finds NotConcerned errors in index functions: cql.Query, cql.Update, cql.Delete
func findNotConcernedForIndex(callExpr *ast.CallExpr, positionsToReport []Report) ([]Report, []string) {
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

func findErrorIsDynamic(positionsToReport []Report, models []string, conditions []ast.Expr) ([]Report, []string) {
	for _, condition := range conditions {
		conditionCall := condition.(*ast.CallExpr)
		conditionSelector, isSelector := conditionCall.Fun.(*ast.SelectorExpr)

		if !isSelector {
			// cql.True
			continue
		}

		if pie.Contains(cqlConnectors, conditionSelector.Sel.Name) {
			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args)
		}

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
	positionsToReport []Report, models []string,
	conditionCall *ast.CallExpr,
	conditionSelector *ast.SelectorExpr,
) []Report {
	whereCondition, isWhereCondition := conditionSelector.X.(*ast.CallExpr)
	if isWhereCondition && getFieldIsMethodName(whereCondition) == "IsDynamic" {
		isDynamicModel, appearance := getModelFromCall(conditionCall.Args[0].(*ast.CallExpr))
		return addPositionsToReport(positionsToReport, models, isDynamicModel, appearance)
	}

	return positionsToReport
}

func addPositionsToReport(positionsToReport []Report, models []string, model Model, appearance Appearance) []Report {
	if !pie.Contains(models, model.Name) {
		return append(positionsToReport, Report{
			model:   model,
			message: notJoinedMessage,
		})
	}

	joinedTimes := len(pie.Filter(models, func(modelName string) bool {
		return modelName == model.Name
	}))

	if appearance.selected {
		if joinedTimes == 1 {
			return append(positionsToReport, Report{
				model:   model,
				message: appearanceNotNecessaryMessage,
			})
		}

		if appearance.number > joinedTimes-1 {
			return append(positionsToReport, Report{
				model:   model,
				message: appearanceOutOfRangeMessage,
			})
		}
	} else if joinedTimes > 1 {
		return append(positionsToReport, Report{
			model:   model,
			message: appearanceMoreThanOnceMessage,
		})
	}

	return positionsToReport
}

func getFieldIsMethodName(whereCondition *ast.CallExpr) string {
	return whereCondition.Fun.(*ast.SelectorExpr).Sel.Name
}

// Returns model's package the model name and true if Appearance method is called
func getModelFromCall(call *ast.CallExpr) (Model, Appearance) {
	fun := call.Fun.(*ast.SelectorExpr)

	funX, isXSelector := fun.X.(*ast.SelectorExpr)
	if isXSelector {
		model := getModel(funX.X.(*ast.SelectorExpr))

		if fun.Sel.Name == "Appearance" {
			appearanceNumber, err := strconv.Atoi(call.Args[0].(*ast.BasicLit).Value)
			if err != nil {
				panic(err)
			}

			return model, Appearance{selected: true, number: appearanceNumber}
		}

		return model, Appearance{selected: false}
	}

	// x is not a selector, so Appearance method or a function is called
	return getModelFromCall(fun.X.(*ast.CallExpr))
}

// Returns model's package the model name
func getModel(selExpr *ast.SelectorExpr) Model {
	return Model{
		Name: passG.TypesInfo.Types[selExpr].Type.Underlying().(*types.Struct).Field(0).Type().(*types.Named).TypeArgs().At(0).(*types.Named).String(),
		Pos:  selExpr.Pos(),
	}
}
