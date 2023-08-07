package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
	"github.com/elliotchance/pie/v2"
)

const BaseURL = "http://localhost:8000"

func (t *TestContext) requestGet(url string) error {
	return t.request(url, http.MethodGet, nil, nil)
}

func (t *TestContext) requestWithJSON(url, method string, jsonTable *godog.Table) error {
	return t.request(url, method, nil, jsonTable)
}

func (t *TestContext) request(url, method string, query map[string]string, jsonTable *godog.Table) error {
	var payload io.Reader

	var err error

	if jsonTable != nil {
		payload = buildJSONFromTable(jsonTable)
	}

	method, err = checkMethod(method)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(method, BaseURL+url, payload)
	if err != nil {
		return fmt.Errorf("failed to build request ERROR=%s", err.Error())
	}

	q := request.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}

	request.URL.RawQuery = q.Encode()

	response, err := t.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to run request ERROR=%s", err.Error())
	}

	t.storeResponseInContext(response)
	response.Body.Close()

	return nil
}

func (t *TestContext) storeResponseInContext(response *http.Response) {
	t.statusCode = response.StatusCode

	err := json.NewDecoder(response.Body).Decode(&t.json)
	if err != nil {
		t.json = map[string]any{}
	}
}

func (t *TestContext) assertStatusCode(expectedStatusCode int) error {
	if t.statusCode != expectedStatusCode {
		return fmt.Errorf("expect status code %d but is %d", expectedStatusCode, t.statusCode)
	}

	return nil
}

func (t *TestContext) assertResponseFieldIsEquals(field string, expectedValue string) error {
	fields := strings.Split(field, ".")
	jsonMap := t.json.(map[string]any)

	for _, field := range fields[:len(fields)-1] {
		intValue, present := jsonMap[field]
		if !present {
			return fmt.Errorf("expected response field %s to be %s but it is not present", field, expectedValue)
		}

		jsonMap = intValue.(map[string]any)
	}

	lastValue, present := jsonMap[pie.Last(fields)]
	if !present {
		return fmt.Errorf("expected response field %s to be %s but it is not present", field, expectedValue)
	}

	if !assertValue(lastValue, expectedValue) {
		return fmt.Errorf("expected response field %s to be %s but is %v", field, expectedValue, lastValue)
	}

	return nil
}

func assertValue(value any, expectedValue string) bool {
	switch value.(type) {
	case string:
		return expectedValue == value
	case int:
		expectedValueInt, err := strconv.Atoi(expectedValue)
		if err != nil {
			panic(err)
		}

		return expectedValueInt == value
	case float64:
		expectedValueFloat, err := strconv.ParseFloat(expectedValue, 64)
		if err != nil {
			panic(err)
		}

		return expectedValueFloat == value
	default:
		panic("unsupported format")
	}
}

// build a map from a godog.Table
func buildMapFromTable(table *godog.Table) (map[string]any, error) {
	data := make(map[string]any, 0)

	err := verifyHeader(table.Rows[0])
	if err != nil {
		return nil, err
	}

	for _, row := range table.Rows[1:] {
		key := row.Cells[0].Value
		valueAsString := row.Cells[1].Value
		valueType := row.Cells[2].Value

		value, err := getTableValue(key, valueAsString, valueType)
		if err != nil {
			return nil, err
		}

		data[key] = value
	}

	return data, nil
}

// Verifies that the header row of a table has the correct format
func verifyHeader(row *messages.PickleTableRow) error {
	for indexCell, cell := range row.Cells {
		if cell.Value != []string{"key", "value", "type"}[indexCell] {
			return fmt.Errorf("should have %q as first line of the table", "| key | value | type |")
		}
	}

	return nil
}

// Returns the value present in a table casted to the correct type
func getTableValue(key, valueAsString, valueType string) (any, error) {
	switch valueType {
	case stringValueType:
		return valueAsString, nil
	case booleanValueType:
		boolean, err := strconv.ParseBool(valueAsString)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as boolean for key %q", valueAsString, key)
		}

		return boolean, nil
	case integerValueType:
		integer, err := strconv.ParseInt(valueAsString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as integer for key %q", valueAsString, key)
		}

		return integer, nil
	case floatValueType:
		floatingNumber, err := strconv.ParseFloat(valueAsString, 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as float for key %q", valueAsString, key)
		}

		return floatingNumber, nil
	case jsonValueType:
		jsonMap := map[string]string{}

		err := json.Unmarshal([]byte(valueAsString), &jsonMap)
		if err != nil {
			return nil, fmt.Errorf("can't parse %q as json for key %q", valueAsString, key)
		}

		return jsonMap, nil
	default:
		return nil, fmt.Errorf(
			"type %q does not exists, please use %v",
			valueType,
			[]string{stringValueType, booleanValueType, integerValueType, floatValueType},
		)
	}
}

// build a json payload in the form of a reader from a godog.Table
func buildJSONFromTable(table *godog.Table) io.Reader {
	data, err := buildMapFromTable(table)
	if err != nil {
		panic("should not return an error")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		panic("should not return an error")
	}

	return strings.NewReader(string(bytes))
}

const (
	stringValueType  = "string"
	booleanValueType = "boolean"
	integerValueType = "integer"
	floatValueType   = "float"
	nullValueType    = "null"
	jsonValueType    = "json"
)

// check if the method is allowed and sanitize the string
func checkMethod(method string) (string, error) {
	allowedMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	sanitizedMethod := strings.TrimSpace(strings.ToUpper(method))

	if !pie.Contains(
		allowedMethods,
		sanitizedMethod,
	) {
		return "", fmt.Errorf("%q is not a valid HTTP method (please choose between %v)", method, allowedMethods)
	}

	return sanitizedMethod, nil
}
