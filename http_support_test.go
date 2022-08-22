package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	buffer, err := ioutil.ReadAll(response.Body)
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
