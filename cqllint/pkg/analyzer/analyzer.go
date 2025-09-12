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
	cqlSelect         = "Select"
	cqlFunctions      = []string{"Query", "Update", "Delete", cqlSelect}
	cqlOrderOrGroupBy = []string{"Descending", "Ascending", "GroupBy", "SelectValue", "Having"}
	cqlConnectors     = []string{"And", "Or", "Not"}
	cqlSetMultiple    = "SetMultiple"
	cqlSets           = []string{cqlSetMultiple, "Set"}
	cqlMethods        = append(cqlOrderOrGroupBy, cqlSets...)

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

type Runner struct {
	positionsToReport []Report
	models            []string
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

	runner := Runner{}

	inspectorG.Preorder(nodeFilter, func(node ast.Node) {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()

		callExpr := node.(*ast.CallExpr)

		// methods to be verified have at least one parameter
		if len(callExpr.Args) < 1 {
			return
		}

		selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
		if isSelector && !selectorIsCQLFunction(selectorExpr) {
			runner.findForSelector(callExpr)
		} else {
			runner.findNotConcernedForIndex(callExpr)
		}
	})

	for _, report := range runner.positionsToReport {
		pass.Reportf(
			report.model.Pos,
			report.message,
			report.model.Name,
		)
	}

	return nil, nil //nolint:nilnil // is necessary
}

// Finds NotConcerned and Repeated errors in selector functions: Descending, Ascending, SetMultiple, Set
func (r *Runner) findForSelector(callExpr *ast.CallExpr) {
	selectorExpr := callExpr.Fun.(*ast.SelectorExpr)

	if !pie.Contains(cqlMethods, selectorExpr.Sel.Name) {
		return
	}

	findRepeatedFields(callExpr, selectorExpr)

	r.fieldNotConcerned(callExpr, selectorExpr)
}

// Finds NotConcerned errors in selector functions: Descending, Ascending, SetMultiple, Set
func (r *Runner) fieldNotConcerned(callExpr *ast.CallExpr, selectorExpr *ast.SelectorExpr) {
	r.findNotConcernedForIndex(selectorExpr.X.(*ast.CallExpr))

	methodName := selectorExpr.Sel.Name
	isOrderOrGroupBy := pie.Contains(cqlOrderOrGroupBy, methodName)

	for _, arg := range callExpr.Args {
		if isOrderOrGroupBy {
			r.findForOrderOrGroupBy(arg)
		} else {
			r.findForSet(arg, methodName)
		}
	}
}

func isAppendCall(call *ast.CallExpr) bool {
	if ident, isIdent := call.Fun.(*ast.Ident); isIdent && ident.Name == "append" {
		return true
	}

	return false
}

func (r *Runner) findForSet(set ast.Expr, methodName string) {
	if setCall, isCall := set.(*ast.CallExpr); isCall {
		r.findForSetCall(setCall, methodName)

		return
	}

	if variable, isVar := set.(*ast.Ident); isVar {
		assignments := findVariableAssignments(variable)
		for _, assign := range assignments {
			r.findForSet(assign.Rhs[0], methodName)
		}
	}

	if composite, isComposite := set.(*ast.CompositeLit); isComposite {
		for _, expr := range composite.Elts {
			r.findForSet(expr, methodName)
		}
	}
}

func (r *Runner) findForSetCall(setCall *ast.CallExpr, methodName string) {
	if isAppendCall(setCall) {
		for _, arg := range setCall.Args[1:] { // first argument is the base list
			r.findForSet(arg, methodName)
		}

		return
	}

	newModels := getModelsFromExpr(setCall.Args[0])
	r.addPositionsToReport(newModels)

	if methodName == cqlSetMultiple {
		// set multiple needs more verifications as the model to be set is not compiled
		newModels := getModelFromCall(setCall.Fun.(*ast.SelectorExpr).X.(*ast.CallExpr))
		r.addPositionsToReport(newModels)
	}
}

func (r *Runner) findForOrderOrGroupBy(expr ast.Expr) {
	if exprCall, isCall := expr.(*ast.CallExpr); isCall {
		if isAppendCall(exprCall) {
			for _, arg := range exprCall.Args[1:] { // first argument is the base list
				r.findForOrderOrGroupBy(arg)
			}

			return
		}

		newModels := getModelFromCall(exprCall)
		r.addPositionsToReport(newModels)

		for _, arg := range exprCall.Args {
			r.findForOrderOrGroupBy(arg)
		}

		return
	}

	if exprSelector, isSelector := expr.(*ast.SelectorExpr); isSelector {
		model := getModel(exprSelector.X.(*ast.SelectorExpr))

		r.addPositionsToReport([]ModelWithAppearance{{model, Appearance{selected: false}}})

		return
	}

	if variable, isVar := expr.(*ast.Ident); isVar {
		assignments := findVariableAssignments(variable)
		for _, assign := range assignments {
			r.findForOrderOrGroupBy(assign.Rhs[0])
		}

		return
	}

	if composite, isComposite := expr.(*ast.CompositeLit); isComposite {
		for _, expr := range composite.Elts {
			r.findForOrderOrGroupBy(expr)
		}

		return
	}
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

// Finds NotConcerned errors in index functions: cql.Query, cql.Update, cql.Delete, cql.Select
func (r *Runner) findNotConcernedForIndex(callExpr *ast.CallExpr) {
	indexExpr, isIndex := callExpr.Fun.(*ast.IndexExpr)
	if isIndex {
		if !selectorIsCQLFunction(indexExpr.X) {
			return
		}

		r.models = []string{
			getFirstGenericType(
				passG.TypesInfo.Types[indexExpr].Type.(*types.Signature).Results().At(0).Type().(*types.Pointer).Elem().(*types.Named),
			),
		}

		r.findErrorIsDynamic(callExpr.Args[1:]) // first parameters is ignored as it's the db object

		return
	}

	selectorExpr, isSelector := callExpr.Fun.(*ast.SelectorExpr)
	if isSelector {
		// other functions may be between callExpr and the cql method, example: cql.Query(...).Limit(1).Descending
		internalCallExpr, isCall := selectorExpr.X.(*ast.CallExpr)
		if isCall {
			r.findNotConcernedForIndex(internalCallExpr)

			return
		}

		if selectorIsCQLSelect(selectorExpr) {
			r.findForSelect(callExpr)

			return
		}

		if selectorIsCQLFunction(selectorExpr) {
			r.findErrorIsDynamic(callExpr.Args[1:]) // first parameters is ignored as it's the db object

			return
		}

		return
	}

	indexListExpr, isIndexList := callExpr.Fun.(*ast.IndexListExpr)
	if isIndexList {
		if !selectorIsCQLSelect(indexListExpr.X) {
			return
		}

		r.findForSelect(callExpr)

		return
	}
}

func selectorIsCQLFunction(expr ast.Expr) bool {
	return selectorIs(expr, cqlFunctions)
}

func selectorIsCQLSelect(expr ast.Expr) bool {
	return selectorIs(expr, []string{cqlSelect})
}

func selectorIs(expr ast.Expr, values []string) bool {
	selectorExpr, isSelector := expr.(*ast.SelectorExpr)
	if !isSelector {
		return false
	}

	ident, isIdent := selectorExpr.X.(*ast.Ident)
	if !isIdent {
		return false
	}

	return ident.Name == "cql" && pie.Contains(values, selectorExpr.Sel.Name)
}

func (r *Runner) findForSelect(callExpr *ast.CallExpr) {
	newRunner := Runner{}

	newRunner.findNotConcernedForIndex(callExpr.Args[0].(*ast.CallExpr))
	newRunner.findErrorIsDynamic(callExpr.Args[1:])

	r.positionsToReport = append(r.positionsToReport, newRunner.positionsToReport...)
}

func (r *Runner) findErrorIsDynamic(conditions []ast.Expr) {
	for _, condition := range conditions {
		if conditionCall, isCall := condition.(*ast.CallExpr); isCall {
			r.findErrorIsDynamicForCall(conditionCall)

			continue
		}

		if variable, isVar := condition.(*ast.Ident); isVar {
			assignments := findVariableAssignments(variable)
			for _, assign := range assignments {
				r.findErrorIsDynamic(assign.Rhs)
			}

			continue
		}

		if composite, isComposite := condition.(*ast.CompositeLit); isComposite {
			r.findErrorIsDynamic(composite.Elts)

			continue
		}
	}
}

func (r *Runner) findErrorIsDynamicForCall(conditionCall *ast.CallExpr) {
	if conditionSelector, isSelector := conditionCall.Fun.(*ast.SelectorExpr); isSelector {
		if pie.Contains(cqlConnectors, conditionSelector.Sel.Name) {
			r.findErrorIsDynamic(conditionCall.Args)

			return
		}

		if conditionSelector.Sel.Name == "ValueInto" {
			newModels := getModelsFromExpr(conditionCall.Args[0])
			r.addPositionsToReport(newModels)

			return
		}

		if conditionSelector.Sel.Name == "Preload" {
			conditionCall = conditionSelector.X.(*ast.CallExpr)
			conditionSelector = conditionCall.Fun.(*ast.SelectorExpr)
		}

		if _, isJoinCondition := conditionSelector.X.(*ast.SelectorExpr); isJoinCondition {
			r.models = append(r.models, getModelFromJoinCondition(conditionSelector))

			r.findErrorIsDynamic(conditionCall.Args)

			return
		}

		r.findErrorIsDynamicWhereCondition(conditionCall, conditionSelector)

		return
	}

	if isAppendCall(conditionCall) {
		r.findErrorIsDynamic(conditionCall.Args[1:]) // first argument is the base list

		return
	}
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

func (r *Runner) findErrorIsDynamicWhereCondition(conditionCall *ast.CallExpr, conditionSelector *ast.SelectorExpr) {
	whereCondition, isWhereCondition := conditionSelector.X.(*ast.CallExpr)
	if isWhereCondition {
		conditionFunc := whereCondition.Fun.(*ast.SelectorExpr)

		if conditionFunc.Sel.Name == "Is" {
			r.addFirstModel(conditionFunc)

			if leftSideCall, isCall := conditionFunc.X.(*ast.CallExpr); isCall {
				// check arguments of the functions on the left side
				r.getPositionsToReportFromCall(leftSideCall)
			}

			// check arguments of the conditions
			r.getPositionsToReportFromCall(conditionCall)
		}
	}
}

func (r *Runner) getPositionsToReportFromCall(call *ast.CallExpr) {
	for _, arg := range call.Args {
		r.addPositionsToReport(getModelsFromExpr(arg))
	}
}

// add first conditions model to models in case is was not set by the query index
func (r *Runner) addFirstModel(conditionFunc *ast.SelectorExpr) {
	if len(r.models) == 0 {
		newModels := getModelsFromExpr(conditionFunc.X)
		if len(newModels) > 0 {
			r.models = append(r.models, newModels[0].Model.Name)
		}
	}
}

func (r *Runner) addPositionsToReport(newModels []ModelWithAppearance) {
	for _, newModel := range newModels {
		model := newModel.Model
		appearance := newModel.Appearance

		if !pie.Contains(r.models, model.Name) {
			newReport := Report{
				model:   model,
				message: notJoinedMessage,
			}

			if !pie.Contains(r.positionsToReport, newReport) {
				r.positionsToReport = append(r.positionsToReport, newReport)
			}
		}

		joinedTimes := len(pie.Filter(r.models, func(modelName string) bool {
			return modelName == model.Name
		}))

		if appearance.bypassCheck {
			continue
		}

		if appearance.selected {
			if joinedTimes == 1 {
				r.positionsToReport = append(r.positionsToReport, Report{
					model:   model,
					message: appearanceNotNecessaryMessage,
				})
			} else if joinedTimes > 1 && appearance.number > joinedTimes-1 {
				r.positionsToReport = append(r.positionsToReport, Report{
					model:   model,
					message: appearanceOutOfRangeMessage,
				})
			}
		} else if joinedTimes > 1 {
			r.positionsToReport = append(r.positionsToReport, Report{
				model:   model,
				message: appearanceMoreThanOnceMessage,
			})
		}
	}
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

func getModelsFromExpr(expr ast.Expr) []ModelWithAppearance {
	argSelector, isSelector := expr.(*ast.SelectorExpr)
	if isSelector {
		model, isModel := getModelFromSelector(argSelector)
		if isModel {
			return []ModelWithAppearance{{model, Appearance{bypassCheck: true}}}
		}

		return nil
	}

	argCall, isCall := expr.(*ast.CallExpr)
	if isCall {
		return getModelFromCall(argCall)
	}

	argVar, isVar := expr.(*ast.Ident)
	if isVar {
		return []ModelWithAppearance{getModelFromVar(argVar)}
	}

	return nil
}

// Returns model's package the model name and true if Appearance method is called
func getModelFromVar(variable *ast.Ident) ModelWithAppearance {
	return ModelWithAppearance{
		getModel(variable),
		Appearance{bypassCheck: true}, // Appearance for variables not implemented
	}
}

type ModelWithAppearance struct {
	Model      Model
	Appearance Appearance
}

// Returns model's package the model name and true if Appearance method is called
func getModelFromCall(call *ast.CallExpr) []ModelWithAppearance {
	newModels := getModelFromCallFunction(call)

	for _, arg := range call.Args {
		newModels = append(newModels, getModelsFromExpr(arg)...)
	}

	return newModels
}

func getModelFromCallFunction(call *ast.CallExpr) []ModelWithAppearance {
	if funSelector, isSelector := call.Fun.(*ast.SelectorExpr); isSelector {
		if selectorX, isXSelector := funSelector.X.(*ast.SelectorExpr); isXSelector {
			model, isModel := getModelFromSelector(selectorX)
			if isModel {
				appearance := getAppearance(call, funSelector)

				return []ModelWithAppearance{ModelWithAppearance{model, appearance}}
			}

			return nil
		}

		if xCall, isCall := funSelector.X.(*ast.CallExpr); isCall {
			// x is not a selector, so Appearance method or a function is called
			return getModelFromCall(xCall)
		}

		if argVar, isVar := funSelector.X.(*ast.Ident); isVar && argVar.Name != "cql" {
			return []ModelWithAppearance{getModelFromVar(argVar)}
		}
	}

	return nil
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
