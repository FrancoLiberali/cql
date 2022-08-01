package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

const BaseUrl = "http://localhost:8000"

func (t *TestContext) requestGET(url string) error {
	response, err := t.httpClient.Get(fmt.Sprintf("%s%s", BaseUrl, url))
	if err != nil {
		return err
	}

	t.storeResponseInContext(response)
	return nil
}

func (t *TestContext) storeResponseInContext(response *http.Response) {
	t.statusCode = response.StatusCode

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	response.Body.Close()
	json.Unmarshal(buffer, &t.json)
}

func (t *TestContext) assertStatusCode(_ context.Context, expectedStatusCode int) error {
	if t.statusCode != expectedStatusCode {
		return fmt.Errorf("expect status code %d but is %d", expectedStatusCode, t.statusCode)
	}
	return nil
}
func (t *TestContext) assertResponseFieldIsEquals(field string, expectedValue string) error {
	value := t.json[field].(string)
	if !assertValue(value, expectedValue) {
		return fmt.Errorf("expect response field %s is %s but is %s", field, expectedValue, value)
	}
	return nil
}

func assertValue(value string, expectedValue string) bool {
	return expectedValue == value
}

func (t *TestContext) requestWithJson(url, method string, jsonTable *godog.Table) error {
	payload, err := buildJSONFromTable(jsonTable)
	if err != nil {
		return err
	}

	method, err = checkMethod(method)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, BaseUrl+url, payload)
	if err != nil {
		return fmt.Errorf("failed to build request ERROR=%s", err.Error())
	}
	response, err := t.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to run request ERROR=%s", err.Error())
	}
	t.storeResponseInContext(response)
	return nil
}

// build a json payload in the form of a reader from a godog.Table
func buildJSONFromTable(table *godog.Table) (io.Reader, error) {
	data := make(map[string]any, 0)
	for indexRow, row := range table.Rows {
		if indexRow == 0 {
			for indexCell, cell := range row.Cells {
				if cell.Value != []string{"key", "value", "type"}[indexCell] {
					return nil, fmt.Errorf("should have %q as first line of the table", "| key | value | type |")
				}
			}
		} else {
			key := row.Cells[0].Value
			valueAsString := row.Cells[1].Value
			valueType := row.Cells[2].Value

			switch valueType {
			case stringValueType:
				data[key] = valueAsString
			case booleanValueType:
				boolean, err := strconv.ParseBool(valueAsString)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as boolean for key %q", valueAsString, key)
				}
				data[key] = boolean
			case integerValueType:
				integer, err := strconv.ParseInt(valueAsString, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as integer for key %q", valueAsString, key)
				}
				data[key] = integer
			case floatValueType:
				floatingNumber, err := strconv.ParseFloat(valueAsString, 64)
				if err != nil {
					return nil, fmt.Errorf("can't parse %q as float for key %q", valueAsString, key)
				}
				data[key] = floatingNumber
			case nullValueType:
				data[key] = nil
			default:
				return nil, fmt.Errorf("type %q does not exists, please use %v", valueType, []string{stringValueType, booleanValueType, integerValueType, floatValueType, nullValueType})
			}

		}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		panic("should not return an error")
	}
	return strings.NewReader(string(bytes)), nil
}

const (
	stringValueType  = "string"
	booleanValueType = "boolean"
	integerValueType = "integer"
	floatValueType   = "float"
	nullValueType    = "null"
)

// check if the method is allowed and sanitize the string
func checkMethod(method string) (string, error) {
	allowedMethods := []string{http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace}
	sanitizedMethod := strings.TrimSpace(strings.ToUpper(method))
	if !contains(
		allowedMethods,
		sanitizedMethod,
	) {
		return "", fmt.Errorf("%q is not a valid HTTP method (please choose between %v)", method, allowedMethods)
	}
	return sanitizedMethod, nil

}

// return true if the set contains the target
func contains[T comparable](set []T, target T) bool {
	for _, elem := range set {
		if target == elem {
			return true
		}
	}
	return false
}
