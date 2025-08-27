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

// force ci
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
	selected    bool
	number      int
	bypassCheck bool // bypassCheck allows to bypass the Appearance call check
}

var (
	passG      *analysis.Pass
	inspectorG *inspector.Inspector
)

func run(pass *analysis.Pass) (interface{}, error) {
	passG = pass
	inspectorG = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	positionsToReport := []Report{}

	inspectorG.Preorder(nodeFilter, func(node ast.Node) {
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

func isAppendCall(call *ast.CallExpr) bool {
	if ident, isIdent := call.Fun.(*ast.Ident); isIdent && ident.Name == "append" {
		return true
	}

	return false
}

func findForSet(set ast.Expr, positionsToReport []Report, models []string, methodName string) []Report {
	if setCall, isCall := set.(*ast.CallExpr); isCall {
		return findForSetCall(setCall, positionsToReport, models, methodName)
	}

	if variable, isVar := set.(*ast.Ident); isVar {
		assignments := findVariableAssignments(variable)
		for _, assign := range assignments {
			positionsToReport = findForSet(assign.Rhs[0], positionsToReport, models, methodName)
		}
	}

	if composite, isComposite := set.(*ast.CompositeLit); isComposite {
		for _, expr := range composite.Elts {
			positionsToReport = findForSet(expr, positionsToReport, models, methodName)
		}
	}

	return positionsToReport
}

func findForSetCall(setCall *ast.CallExpr, positionsToReport []Report, models []string, methodName string) []Report {
	if isAppendCall(setCall) {
		for _, arg := range setCall.Args[1:] { // first argument is the base list
			positionsToReport = findForSet(arg, positionsToReport, models, methodName)
		}

		return positionsToReport
	}

	model, appearance, isModel := getModelFromExpr(setCall.Args[0])
	if isModel {
		positionsToReport = addPositionsToReport(positionsToReport, models, model, appearance)
	}

	if methodName == cqlSetMultiple {
		// set multiple needs more verifications as the model to be set is not compiled
		model, appearance, isModel := getModelFromCall(setCall.Fun.(*ast.SelectorExpr).X.(*ast.CallExpr))
		if isModel {
			positionsToReport = addPositionsToReport(positionsToReport, models, model, appearance)
		}
	}

	return positionsToReport
}

func findForOrder(order ast.Expr, positionsToReport []Report, models []string) []Report {
	if orderCall, isCall := order.(*ast.CallExpr); isCall {
		model, appearance, isModel := getModelFromCall(orderCall)
		if isModel {
			return addPositionsToReport(positionsToReport, models, model, appearance)
		}

		return positionsToReport
	}

	if orderSelector, isSelector := order.(*ast.SelectorExpr); isSelector {
		model := getModel(orderSelector.X.(*ast.SelectorExpr))

		return addPositionsToReport(positionsToReport, models, model, Appearance{selected: false})
	}

	if variable, isVar := order.(*ast.Ident); isVar {
		assignments := findVariableAssignments(variable)
		for _, assign := range assignments {
			positionsToReport = findForOrder(assign.Rhs[0], positionsToReport, models)
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
		if argCall, isCall := arg.(*ast.CallExpr); isCall {
			condition := argCall.Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.SelectorExpr)

			fieldName := getFieldName(condition)
			fieldPos := condition.Sel.NamePos

			_, isPresent := fields[fieldName]
			if !isPresent {
				fields[fieldName] = []token.Pos{fieldPos}
			} else {
				fields[fieldName] = append(fields[fieldName], fieldPos)
			}

			for _, internalArg := range argCall.Args {
				comparedField, isSelector := internalArg.(*ast.SelectorExpr)
				if !isSelector {
					// only for selectors, as they are the field.
					// if not selector, it means function is applied to the field so it is not the same value
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
		if conditionCall, isCall := condition.(*ast.CallExpr); isCall {
			positionsToReport, models = findErrorIsDynamicForCall(positionsToReport, models, conditionCall)

			continue
		}

		if variable, isVar := condition.(*ast.Ident); isVar {
			assignments := findVariableAssignments(variable)
			for _, assign := range assignments {
				positionsToReport, models = findErrorIsDynamic(positionsToReport, models, assign.Rhs)
			}

			continue
		}

		if composite, isComposite := condition.(*ast.CompositeLit); isComposite {
			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, composite.Elts)

			continue
		}
	}

	return positionsToReport, models
}

func findErrorIsDynamicForCall(positionsToReport []Report, models []string, conditionCall *ast.CallExpr) ([]Report, []string) {
	if conditionSelector, isSelector := conditionCall.Fun.(*ast.SelectorExpr); isSelector {
		if pie.Contains(cqlConnectors, conditionSelector.Sel.Name) {
			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args)

			return positionsToReport, models
		}

		if conditionSelector.Sel.Name == "Preload" {
			conditionCall = conditionSelector.X.(*ast.CallExpr)
			conditionSelector = conditionCall.Fun.(*ast.SelectorExpr)
		}

		if _, isJoinCondition := conditionSelector.X.(*ast.SelectorExpr); isJoinCondition {
			models = append(models, getModelFromJoinCondition(conditionSelector))

			positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args)

			return positionsToReport, models
		}

		positionsToReport = findErrorIsDynamicWhereCondition(positionsToReport, models, conditionCall, conditionSelector)

		return positionsToReport, models
	}

	if isAppendCall(conditionCall) {
		positionsToReport, models = findErrorIsDynamic(positionsToReport, models, conditionCall.Args[1:]) // first argument is the base list

		return positionsToReport, models
	}

	// cql.True
	return positionsToReport, models
}

func findVariableAssignments(variable *ast.Ident) []*ast.AssignStmt {
	assignments := []*ast.AssignStmt{}

	cursor, found := inspectorG.Root().FindByPos(passG.TypesInfo.ObjectOf(variable).Parent().Pos(), passG.TypesInfo.ObjectOf(variable).Parent().End())
	if found {
		for cursor := range cursor.Preorder((*ast.AssignStmt)(nil)) {
			assign := cursor.Node().(*ast.AssignStmt)
			if assign.Lhs[0].(*ast.Ident).Name == variable.Name {
				assignments = append(assignments, assign)
			}
		}
	}

	return assignments
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
	if isWhereCondition && getFieldIsMethodName(whereCondition) == "Is" {
		for _, arg := range conditionCall.Args {
			model, appearance, isModel := getModelFromExpr(arg)
			if isModel {
				positionsToReport = addPositionsToReport(positionsToReport, models, model, appearance)
			}
		}
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

	if appearance.bypassCheck {
		return positionsToReport
	}

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

// getAppearance returns the Appearance from a model
func getAppearance(call *ast.CallExpr, fun *ast.SelectorExpr) Appearance {
	if fun.Sel.Name == "Appearance" {
		appearanceNumber, err := strconv.Atoi(call.Args[0].(*ast.BasicLit).Value)
		if err != nil {
			panic(err)
		}

		return Appearance{selected: true, number: appearanceNumber}
	}

	return Appearance{selected: false}
}

func getModelFromExpr(expr ast.Expr) (Model, Appearance, bool) {
	argSelector, isSelector := expr.(*ast.SelectorExpr)
	if isSelector {
		model, isModel := getModelFromSelector(argSelector)

		return model, Appearance{}, isModel
	}

	argCall, isCall := expr.(*ast.CallExpr)
	if isCall {
		return getModelFromCall(argCall)
	}

	argVar, isVar := expr.(*ast.Ident)
	if isVar {
		return getModelFromVar(argVar)
	}

	return Model{}, Appearance{}, false
}

// Returns model's package the model name and true if Appearance method is called
func getModelFromVar(variable *ast.Ident) (Model, Appearance, bool) {
	return getModel(variable),
		Appearance{bypassCheck: true}, // Appearance for variables not implemented
		true
}

// Returns model's package the model name and true if Appearance method is called
func getModelFromCall(call *ast.CallExpr) (Model, Appearance, bool) {
	if funSelector, isSelector := call.Fun.(*ast.SelectorExpr); isSelector {
		if selectorX, isXSelector := funSelector.X.(*ast.SelectorExpr); isXSelector {
			model, isModel := getModelFromSelector(selectorX)
			if isModel {
				return model, getAppearance(call, funSelector), isModel
			}

			return Model{}, Appearance{}, false
		}

		if xCall, isCall := funSelector.X.(*ast.CallExpr); isCall {
			// x is not a selector, so Appearance method or a function is called
			return getModelFromCall(xCall)
		}

		if argVar, isVar := funSelector.X.(*ast.Ident); isVar && argVar.Name != "cql" {
			return getModelFromVar(argVar)
		}
	}

	return Model{}, Appearance{}, false
}

// Returns model's package the model name and true if Appearance method is called
func getModelFromSelector(selector *ast.SelectorExpr) (Model, bool) {
	if selectorX, isXSelector := selector.X.(*ast.SelectorExpr); isXSelector {
		return getModel(selectorX), true
	}

	return Model{}, false
}

// Returns model's package the model name
func getModel(e ast.Expr) Model {
	return Model{
		Name: passG.TypesInfo.TypeOf(e).Underlying().(*types.Struct).Field(0).Type().(*types.Named).TypeArgs().At(0).(*types.Named).String(),
		Pos:  e.Pos(),
	}
}
